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
	DEFAULT_GOLIGHTCTRL_MQTTBROKER      string = "tcp://mqtt.realraum.at:1883"
	DEFAULT_GOLIGHTCTRL_RF433TTYDEV     string = "/dev/ttyACM0"
	DEFAULT_GOLIGHTCTRL_BUTTONTTYDEV    string = "/dev/ttyACM1"
	DEFAULT_GOLIGHTCTRL_BASICCTRLTTYDEV string = "/dev/ttyBasicCtrl"
)

var (
	UseFakeGPIO_         bool
	DebugFlags_          string
	ps_                  *pubsub.PubSub
	topic_lightctrl_pre_ string = r3events.TOPIC_ACTIONS + r3events.CLIENTID_LIGHTCTRL + "/"
	CeilingLightsSwitch_ CeilingLightsSwitch
)

func init() {
	flag.BoolVar(&UseFakeGPIO_, "fakegpio", false, "For testing")
	flag.StringVar(&DebugFlags_, "debug", "", "List of DebugFlags separated by ,")
	ps_ = pubsub.New(50)
}

//Connect and keep trying to connect to MQTT Broker
//And while we cannot, still provide as much functionality as possible
func goConnectToMQTTBrokerAndFunctionWithoutInTheMeantime(tty_rf433_chan chan SerialLine) {
	//Start Channel Gobblers and Functionality that works without mqtt
	//These shut down once we send PS_SHUTDOWN_CONSUMER
	go goLinearizeRFSenders(ps_, RF433_linearize_chan_, tty_rf433_chan, nil)
	//consume MQTT_ledpattern_chan_ and MQTT_ir_chan_
	go func() {
		shutdown_c := ps_.SubOnce(PS_SHUTDOWN_CONSUMER)
		for {
			select {
			case <-shutdown_c:
				return
			case <-MQTT_ir_chan_:
				//drop
			case <-MQTT_chan_:
				//drop
			}
		}
	}()
	//
	//
	//now try connect to mqtt daemon until it works once
	for {
		mqttc := ConnectMQTTBroker(EnvironOrDefault("GOLIGHTCTRL_MQTTBROKER", DEFAULT_GOLIGHTCTRL_MQTTBROKER), EnvironOrDefault("GOLIGHTCTRL_MQTTCLIENTID", r3events.CLIENTID_LIGHTCTRL))
		//start real goroutines after mqtt connected
		if mqttc != nil {
			SubscribeAndAttachCallback(mqttc, r3events.ACT_LIGHTCTRL_NAME, func(c mqtt.Client, msg mqtt.Message) {
				var aon r3events.LightCtrlActionOnName
				if msg.Retained() {
					return
				}
				if err := json.Unmarshal(msg.Payload(), &aon); err != nil {
					LogMain_.Printf("Main:LightCtrlMain: Error: %s", err)
					return
				}
				LogMain_.Printf("Main:LightCtrlMain: %+v", aon)
				switch_name_chan_ <- aon
			})
			for name, _ := range actionname_map_ {
				SubscribeAndAttachCallback(mqttc, topic_lightctrl_pre_+name, func(c mqtt.Client, msg mqtt.Message) {
					var aon r3events.LightCtrlActionOnName
					if msg.Retained() {
						return
					}
					aon.Name = msg.Topic()[len(topic_lightctrl_pre_):]
					aon.Action = string(msg.Payload())
					LogMain_.Printf("Main:LightCtrlMain: %+v", aon)
					switch_name_chan_ <- aon
				})
			}
			ps_.Pub(true, PS_SHUTDOWN_CONSUMER) //shutdown all chan consumers for mqttc == nil
			time.Sleep(5 * time.Second)         //avoid goLinearizeRFSender that we start below to shutdown right away
			go goSendIRCmdToMQTT(mqttc, MQTT_ir_chan_)
			go goSendMQTTMsg(mqttc, MQTT_chan_)
			go goLinearizeRFSenders(ps_, RF433_linearize_chan_, tty_rf433_chan, mqttc)
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

	var err error
	var tty_button_chan chan SerialLine
	var tty_rf433_chan chan SerialLine
	if UseFakeGPIO_ {
		tty_rf433_chan = make(chan SerialLine, 10)
		go func() {
			for str := range tty_rf433_chan {
				LogRF433_.Println(str)
			}
		}()
		tty_button_chan = make(chan SerialLine, 1)
		CeilingLightsSwitch_ = CeilinglightsGPIO_FakeGPIOinit()
	} else {
		tty_rf433_chan, _, err = OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_RF433TTYDEV", DEFAULT_GOLIGHTCTRL_RF433TTYDEV), 9600)
		if err != nil {
			panic("can't open GOLIGHTCTRL_RF433TTYDEV")
		}

		_, tty_button_chan, err = OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_BUTTONTTYDEV", DEFAULT_GOLIGHTCTRL_BUTTONTTYDEV), 9600)
		if err != nil {
			panic("can't open GOLIGHTCTRL_BUTTONTTYDEV")
		}

		tty_basicctrl_write_chan, tty_basicctrl_read_chan, err := OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_BASICCTRLTTYDEV", DEFAULT_GOLIGHTCTRL_BASICCTRLTTYDEV), 9600)
		if err != nil {
			LogMain_.Println("can't open GOLIGHTCTRL_BASICCTRLTTYDEV")
			LogMain_.Println("switching CeilingLights via RPi GPIO")
			CeilingLightsSwitch_ = CeilinglightsGPIO_GPIOinit()
		} else {
			LogMain_.Println("switching CeilingLights via BasicCtrl tty")
			CeilingLightsSwitch_ = NewBasicCtrl(tty_basicctrl_read_chan, tty_basicctrl_write_chan)
		}

	}

	go GoSwitchNameAsync()
	go goConnectToMQTTBrokerAndFunctionWithoutInTheMeantime(tty_rf433_chan)
	go goListenForButtons(tty_button_chan)

	// wait on Ctrl-C or sigInt or sigKill
	func() {
		ctrlc_c := make(chan os.Signal, 1)
		signal.Notify(ctrlc_c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-ctrlc_c //block until ctrl+c is pressed || we receive SIGINT aka kill -1 || kill
		fmt.Println("SIGINT received, exiting gracefully ...")
	}()

}
