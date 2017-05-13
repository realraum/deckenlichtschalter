// (c) Bernhard Tittelbach, 2016
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pubsub "github.com/btittelbach/pubsub"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/realraum/door_and_sensors/r3events"
)

const (
	DEFAULT_GOMQTTWEBFRONT_MQTTBROKER     string = "tcp://mqtt.realraum.at:1883"
	DEFAULT_GOMQTTWEBFRONT_HTTP_INTERFACE string = ":80"
	DEFAULT_GOMQTTWEBFRONT_RF433TTYDEV    string = "/dev/ttyACM0"
	DEFAULT_GOMQTTWEBFRONT_BUTTONTTYDEV   string = "/dev/ttyACM1"
)

type SerialLine []byte

var (
	DebugFlags_ string
	ps_         *pubsub.PubSub
)

func init() {
	flag.StringVar(&DebugFlags_, "debug", "", "List of DebugFlags separated by ,")
	ps_ = pubsub.New(50)
}

//Connect and keep trying to connect to MQTT Broker
//And while we cannot, still provide as much functionality as possible
func goConnectToMQTTBrokerAndFunctionWithoutInTheMeantime() {
	//Start Channel Gobblers and Functionality that works without mqtt
	//These shut down once we send PS_SHUTDOWN_CONSUMER
	//consume MQTT_ledpattern_chan_ and MQTT_ir_chan_
	go func() {
		shutdown_c := ps_.SubOnce(PS_SHUTDOWN_CONSUMER)
		for {
			select {
			case <-shutdown_c:
				return
			case <-MQTT_sendmsg_chan_:
				//drop msg
			}
		}
	}()
	//
	//
	//now try connect to mqtt daemon until it works once
	for {
		mqttc := ConnectMQTTBroker(EnvironOrDefault("GOMQTTWEBFRONT_MQTTBROKER", DEFAULT_GOMQTTWEBFRONT_MQTTBROKER), EnvironOrDefault("GOMQTTWEBFRONT_CLIENTID", r3events.CLIENTID_WEBFRONT))
		//start real goroutines after mqtt connected
		if mqttc != nil {
			topic_in_chan := SubscribeMultipleAndForwardToChannel(mqttc, ws_allowed_ctx_all)
			go func(c mqtt.Client, msg_in_chan chan mqtt.Message) {
				// if msg.Retained() {
				// 	return
				// }
				for msg := range msg_in_chan {
					lp := make(map[string]interface{}, 10)
					//Error check, then forward
					if err := json.Unmarshal(msg.Payload(), &lp); err == nil {
						webmsg := wsMessage{Ctx: msg.Topic(), Data: lp}
						ps_.Pub(webmsg, PS_WEBSOCK_ALL)
					} else {
						webmsg := wsMessage{Ctx: msg.Topic(), Data: string(msg.Payload())}
						ps_.Pub(webmsg, PS_WEBSOCK_ALL)
					}
				}
			}(mqttc, topic_in_chan)
			ps_.Pub(true, PS_SHUTDOWN_CONSUMER) //shutdown all chan consumers for mqttc == nil
			time.Sleep(5 * time.Second)         //avoid goLinearizeRFSender that we start below to shutdown right away
			go goSendMQTTMsgToBroker(mqttc, MQTT_sendmsg_chan_)
			//and LAST but not least:
			RequestStatusFromAllFancyLightsMQTT(mqttc)
			return // no need to keep on trying, mqtt-auto-reconnect will do the rest now
		} else {
			time.Sleep(5 * time.Minute)
		}
	}
}

func main() {
	flag.Parse()
	if len(DebugFlags_) > 0 {
		LogEnable(strings.Split(DebugFlags_, ",")...)
	}

	go goConnectToMQTTBrokerAndFunctionWithoutInTheMeantime()
	go goRunMartini()

	// wait on Ctrl-C or sigInt or sigKill
	func() {
		ctrlc_c := make(chan os.Signal, 1)
		signal.Notify(ctrlc_c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-ctrlc_c //block until ctrl+c is pressed || we receive SIGINT aka kill -1 || kill
		fmt.Println("SIGINT received, exiting gracefully ...")
	}()

}
