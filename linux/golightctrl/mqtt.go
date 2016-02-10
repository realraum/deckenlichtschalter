// (c) Bernhard Tittelbach, 2015
package main

import (
	"sync"
	"time"

	"github.com/realraum/door_and_sensors/r3events"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
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

func mqttOnConnectionHandler(mqttc *mqtt.Client) {
	LogMQTT_.Print("MQTT connection to broker established. (re)subscribing topics")
	mqtt_topics_we_subscribed_lock_.RLock()
	defer mqtt_topics_we_subscribed_lock_.RUnlock()
	if len(mqtt_topics_we_subscribed_) > 0 {
		tk := mqttc.SubscribeMultiple(mqtt_topics_we_subscribed_, nil)
		tk.Wait()
		if tk.Error() != nil {
			LogMQTT_.Fatalf("Error resubscribing on connect", tk.Error())
		}
	}
}

func ConnectMQTTBroker(brocker_addr, clientid string) *mqtt.Client {
	options := mqtt.NewClientOptions().AddBroker(brocker_addr).SetAutoReconnect(true).SetKeepAlive(30 * time.Second).SetMaxReconnectInterval(2 * time.Minute)
	options = options.SetClientID(clientid).SetConnectionLostHandler(func(c *mqtt.Client, err error) { LogMQTT_.Print("ERROR MQTT connection lost:", err) })
	options = options.SetOnConnectHandler(mqttOnConnectionHandler)
	c := mqtt.NewClient(options)
	tk := c.Connect()
	tk.Wait()
	if tk.Error() != nil {
		LogMQTT_.Fatal("Error connecting to mqtt broker", tk.Error())
	}
	return c
}

func goSendCodeToMQTT(mqttc *mqtt.Client, code_chan chan []byte) {
	for code := range code_chan {
		LogMQTT_.Printf("goSendToMQTT(%+v)", code)
		if len(code) == 3 {
			r3evt := r3events.SendRF433Code{Code: [3]byte{code[0], code[1], code[2]}, Ts: time.Now().Unix()}
			LogMQTT_.Printf("goSendToMQTT: %+v", r3evt)
			mqttc.Publish(r3events.ACT_RF433_SEND, MQTT_QOS_REQCONFIRMATION, false, r3events.MarshalEvent2ByteOrPanic(r3evt))
		}
	}
}

func goSendIRCmdToMQTT(mqttc *mqtt.Client, ir_chan chan string) {
	for cmd := range ir_chan {
		r3evt := r3events.YamahaIRCmd{Cmd: cmd, Ts: time.Now().Unix()}
		LogMQTT_.Printf("goSendIRCmdToMQTT: %+v", r3evt)
		mqttc.Publish(r3events.ACT_YAMAHA_SEND, MQTT_QOS_REQCONFIRMATION, false, r3events.MarshalEvent2ByteOrPanic(r3evt))
	}
}
