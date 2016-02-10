// (c) Bernhard Tittelbach, 2016
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	DEFAULT_GOLIGHTCTRL_MQTTBROKER     string = "tcp://mqtt.realraum.at:1883"
	DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE string = ":80"
	DEFAULT_GOLIGHTCTRL_RF433TTYDEV    string = "/dev/ttyACM0"
	DEFAULT_GOLIGHTCTRL_BUTTONTTYDEV   string = "/dev/ttyACM1"
	DEFAULT_GOLIGHTCTRL_MQTTCLIENTID   string = "GoLightCtrl"
)

type SerialLine []byte

var (
	UseFakeGPIO_ bool
	DebugFlags_  string
)

func init() {
	flag.BoolVar(&UseFakeGPIO_, "fakegpio", false, "For testing")
	flag.StringVar(&DebugFlags_, "debug", "", "List of DebugFlags separated by ,")
}
func main() {
	flag.Parse()
	if len(DebugFlags_) > 0 {
		LogEnable(strings.Split(DebugFlags_, ",")...)
	}

	if UseFakeGPIO_ {
		FakeGPIOinit()
	} else {
		GPIOinit()
	}

	var err error
	var tty_button_chan chan SerialLine
	if UseFakeGPIO_ {
		RF433_chan_ = make(chan []byte, 10)
		go func() {
			for str := range RF433_chan_ {
				LogRF433_.Println(str)
			}
		}()
		tty_button_chan = make(chan SerialLine, 1)
	} else {
		RF433_chan_, _, err = OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_RF433TTYDEV", DEFAULT_GOLIGHTCTRL_RF433TTYDEV), 9600)
		if err != nil {
			panic("can't open GOLIGHTCTRL_RF433TTYDEV")
		}

		_, tty_button_chan, err = OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_BUTTONTTYDEV", DEFAULT_GOLIGHTCTRL_BUTTONTTYDEV), 9600)
		if err != nil {
			panic("can't open GOLIGHTCTRL_BUTTONTTYDEV")
		}
	}
	mqttc := ConnectMQTTBroker(EnvironOrDefault("GOLIGHTCTRL_MQTTBROKER", DEFAULT_GOLIGHTCTRL_MQTTBROKER), EnvironOrDefault("GOLIGHTCTRL_MQTTCLIENTID", DEFAULT_GOLIGHTCTRL_MQTTCLIENTID))

	MQTT_rf_chan_ = make(chan []byte, 10)
	MQTT_ir_chan_ = make(chan string, 10)
	go goSendCodeToMQTT(mqttc, MQTT_rf_chan_)
	go goSendIRCmdToMQTT(mqttc, MQTT_ir_chan_)

	go goListenForButtons(tty_button_chan)
	go goRunMartini()

	// wait on Ctrl-C or sigInt or sigKill
	func() {
		ctrlc_c := make(chan os.Signal, 1)
		signal.Notify(ctrlc_c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-ctrlc_c //block until ctrl+c is pressed || we receive SIGINT aka kill -1 || kill
		fmt.Println("SIGINT received, exiting gracefully ...")
	}()

}
