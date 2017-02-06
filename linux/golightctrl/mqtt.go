// (c) Bernhard Tittelbach, 2015
package main

import (
	"sync"
	"time"

	"github.com/realraum/door_and_sensors/r3events"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const MQTT_QOS_NOCONFIRMATION byte = 0
const MQTT_QOS_REQCONFIRMATION byte = 1
const MQTT_QOS_4STPHANDSHAKE byte = 2

var mqtt_topics_we_subscribed_ map[string]byte
var mqtt_topics_we_subscribed_lock_ sync.RWMutex

func init() {
	mqtt_topics_we_subscribed_ = make(map[string]byte, 1)
}

func addSubscribedTopics(subresult map[string]byte) {
	mqtt_topics_we_subscribed_lock_.Lock()
	defer mqtt_topics_we_subscribed_lock_.Unlock()
	for topic, qos := range subresult {
		if qos < 0 || qos > 2 {
			LogMQTT_.Printf("addSubscribedTopics: not remembering topic since we didn't subscribe it successfully: %s (qos: %d)", topic, qos)
			continue
		}
		LogMQTT_.Printf("addSubscribedTopics: remembering subscribed topic: %s (qos: %d)", topic, qos)
		mqtt_topics_we_subscribed_[topic] = qos
	}
}

func removeSubscribedTopic(topic string) {
	mqtt_topics_we_subscribed_lock_.Lock()
	defer mqtt_topics_we_subscribed_lock_.Unlock()
	delete(mqtt_topics_we_subscribed_, topic)
	LogMQTT_.Printf("removeSubscribedTopics: %s ", topic)
}

func mqttOnConnectionHandler(mqttc mqtt.Client) {
	if mqttc == nil {
		return
	}
	LogMQTT_.Print("MQTT connection to broker established. (re)subscribing topics")
	mqtt_topics_we_subscribed_lock_.RLock()
	defer mqtt_topics_we_subscribed_lock_.RUnlock()
	if len(mqtt_topics_we_subscribed_) > 0 {
		tk := mqttc.SubscribeMultiple(mqtt_topics_we_subscribed_, nil)
		tk.Wait()
		if tk.Error() != nil {
			LogMQTT_.Fatal("Error resubscribing on connect", tk.Error())
		}
	}
}

func ConnectMQTTBroker(brocker_addr, clientid string) mqtt.Client {
	options := mqtt.NewClientOptions().AddBroker(brocker_addr).SetAutoReconnect(true).SetKeepAlive(30 * time.Second).SetMaxReconnectInterval(2 * time.Minute)
	options = options.SetClientID(clientid).SetConnectionLostHandler(func(c mqtt.Client, err error) { LogMQTT_.Print("ERROR MQTT connection lost:", err) })
	options = options.SetOnConnectHandler(mqttOnConnectionHandler)
	c := mqtt.NewClient(options)
	tk := c.Connect()
	tk.Wait()
	if tk.Error() != nil {
		LogMQTT_.Println("Error connecting to mqtt broker", tk.Error())
		return nil
	}
	return c
}

func RequestStatusFromAllFancyLightsMQTT(mqttc mqtt.Client) {
	if mqttc == nil {
		return
	}
	mqttc.Publish(r3events.ACT_ALLFANCYLIGHT_PLEASEREPEAT, MQTT_QOS_NOCONFIRMATION, false, []byte{})
}

func sendCodeToMQTT(mqttc mqtt.Client, code []byte) {
	if mqttc == nil {
		return
	}
	LogMQTT_.Printf("sendCodeToMQTT(%+v)", code)
	if len(code) == 3 {
		r3evt := r3events.SendRF433Code{Code: [3]byte{code[0], code[1], code[2]}, Ts: time.Now().Unix()}
		LogMQTT_.Printf("goSendToMQTT: %+v", r3evt)
		mqttc.Publish(r3events.ACT_RF433_SEND, MQTT_QOS_REQCONFIRMATION, false, r3events.MarshalEvent2ByteOrPanic(r3evt))
	}
}

func goSendIRCmdToMQTT(mqttc mqtt.Client, ir_chan chan string) {
	if mqttc == nil {
		return
	}
	for cmd := range ir_chan {
		r3evt := r3events.YamahaIRCmd{Cmd: cmd, Ts: time.Now().Unix()}
		LogMQTT_.Printf("goSendIRCmdToMQTT: %+v", r3evt)
		mqttc.Publish(r3events.ACT_YAMAHA_SEND, MQTT_QOS_REQCONFIRMATION, false, r3events.MarshalEvent2ByteOrPanic(r3evt))
	}
}

func goSetLEDPipePatternViaMQTT(mqttc mqtt.Client, pipepattern_chan chan *r3events.SetPipeLEDsPattern) {
	if mqttc == nil {
		return
	}
	for r3evtptr := range pipepattern_chan {
		LogMQTT_.Printf("SetPipeLEDsPattern: %+v", *r3evtptr)
		mqttc.Publish(r3events.ACT_PIPELEDS_PATTERN, MQTT_QOS_REQCONFIRMATION, false, r3events.MarshalEvent2ByteOrPanic(*r3evtptr))
	}
}

func goSetFancyLightsViaMQTT(mqttc mqtt.Client, fancylights_chan chan *wsMsgFancyLight) {
	if mqttc == nil {
		return
	}
	for r3evtptr := range fancylights_chan {
		json := r3events.MarshalEvent2ByteOrPanic(r3evtptr.Setting)
		LogMQTT_.Printf("goSetFancyLightsViaMQTT: %+v", *r3evtptr)
		LogMQTT_.Printf("goSetFancyLightsViaMQTT: %s", json)
		mqttc.Publish(r3events.TOPIC_ACTIONS+r3evtptr.Name+"/"+r3events.TYPE_LIGHT, MQTT_QOS_NOCONFIRMATION, false, json)
	}
}

func SubscribeAndAttachCallback(mqttc mqtt.Client, filter string, callback mqtt.MessageHandler) {
	tk := mqttc.Subscribe(filter, 0, callback)
	tk.Wait()
	if tk.Error() != nil {
		LogMQTT_.Fatalf("Error subscribing to %s:%s", filter, tk.Error())
	} else {
		LogMQTT_.Printf("SubscribeAndForwardToChannel successfull")
		addSubscribedTopics(tk.(*mqtt.SubscribeToken).Result())
	}
	return
}

func SubscribeAndForwardToChannel(mqttc mqtt.Client, filter string) (channel chan mqtt.Message) {
	channel = make(chan mqtt.Message, 100)
	tk := mqttc.Subscribe(filter, 0, func(mqttc mqtt.Client, msg mqtt.Message) { channel <- msg })
	tk.Wait()
	if tk.Error() != nil {
		LogMQTT_.Fatalf("Error subscribing to %s:%s", filter, tk.Error())
	} else {
		LogMQTT_.Printf("SubscribeAndForwardToChannel successfull")
		addSubscribedTopics(tk.(*mqtt.SubscribeToken).Result())
	}
	return
}

func SubscribeMultipleAndForwardToChannel(mqttc mqtt.Client, filters []string) (channel chan mqtt.Message) {
	channel = make(chan mqtt.Message, 100)
	filtermap := make(map[string]byte, len(filters))
	for _, topicfilter := range filters {
		filtermap[topicfilter] = 0 //qos == 0
	}
	tk := mqttc.SubscribeMultiple(filtermap, func(mqttc mqtt.Client, msg mqtt.Message) {
		LogMQTT_.Printf("forwarding mqtt message to channel %+v", msg)
		channel <- msg
	})
	tk.Wait()
	if tk.Error() != nil {
		LogMQTT_.Fatalf("Error subscribing to %s:%s", filters, tk.Error())
	} else {
		LogMQTT_.Printf("SubscribeMultipleAndForwardToChannel successfull")
		addSubscribedTopics(tk.(*mqtt.SubscribeToken).Result())
	}
	return
}

func UnsubscribeMultiple(mqttc mqtt.Client, topics ...string) {
	mqttc.Unsubscribe(topics...)
	for _, topic := range topics {
		removeSubscribedTopic(topic)
	}
}
