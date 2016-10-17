// (c) Bernhard Tittelbach, 2016
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	pubsub "github.com/btittelbach/pubsub"
	"github.com/realraum/door_and_sensors/r3events"
)

const (
	DEFAULT_GOLIGHTCTRL_MQTTBROKER     string = "tcp://mqtt.realraum.at:1883"
	DEFAULT_GOLIGHTCTRL_HTTP_INTERFACE string = ":80"
	DEFAULT_GOLIGHTCTRL_RF433TTYDEV    string = "/dev/ttyACM0"
	DEFAULT_GOLIGHTCTRL_BUTTONTTYDEV   string = "/dev/ttyACM1"
)

type SerialLine []byte

var (
	UseFakeGPIO_ bool
	DebugFlags_  string
	ps_          *pubsub.PubSub
)

const (
	PS_WEBSOCK_ALL_JSON = "websock_toall_json"
	PS_WEBSOCK_ALL      = "websock_toall"
	PS_LIGHTS_CHANGED   = "light_state_changed"
	PS_IRRF433_CHANGED  = "stateless_button_send_event"
	PS_SHUTDOWN         = "shutdown"
)

func init() {
	flag.BoolVar(&UseFakeGPIO_, "fakegpio", false, "For testing")
	flag.StringVar(&DebugFlags_, "debug", "", "List of DebugFlags separated by ,")
	ps_ = pubsub.New(50)
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
	var tty_rf433_chan chan SerialLine
	if UseFakeGPIO_ {
		tty_rf433_chan = make(chan SerialLine, 10)
		go func() {
			for str := range tty_rf433_chan {
				LogRF433_.Println(str)
			}
		}()
		tty_button_chan = make(chan SerialLine, 1)
	} else {
		tty_rf433_chan, _, err = OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_RF433TTYDEV", DEFAULT_GOLIGHTCTRL_RF433TTYDEV), 9600)
		if err != nil {
			panic("can't open GOLIGHTCTRL_RF433TTYDEV")
		}

		_, tty_button_chan, err = OpenAndHandleSerial(EnvironOrDefault("GOLIGHTCTRL_BUTTONTTYDEV", DEFAULT_GOLIGHTCTRL_BUTTONTTYDEV), 9600)
		if err != nil {
			panic("can't open GOLIGHTCTRL_BUTTONTTYDEV")
		}
	}
	mqttc := ConnectMQTTBroker(EnvironOrDefault("GOLIGHTCTRL_MQTTBROKER", DEFAULT_GOLIGHTCTRL_MQTTBROKER), r3events.CLIENTID_LIGHTCTRL)

	go goLinearizeRFSenders(RF433_linearize_chan_, tty_rf433_chan, mqttc)
	go goSendIRCmdToMQTT(mqttc, MQTT_ir_chan_)
	go goSetLEDPipePatternViaMQTT(mqttc, MQTT_ledpattern_chan_)

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
