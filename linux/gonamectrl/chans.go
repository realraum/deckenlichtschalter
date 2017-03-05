// (c) Bernhard Tittelbach 2017
package main

import "github.com/realraum/door_and_sensors/r3events"

const (
	PS_LIGHTS_CHANGED    = "light_state_changed"
	PS_IRRF433_CHANGED   = "stateless_button_send_event"
	PS_SHUTDOWN          = "shutdown"
	PS_SHUTDOWN_CONSUMER = "shutdownindiscriminateconsumer"
)

var switch_name_chan_ chan r3events.LightCtrlActionOnName
var MQTT_ir_chan_ chan string
var RF433_linearize_chan_ chan RFCmdToSend

func init() {
	switch_name_chan_ = make(chan r3events.LightCtrlActionOnName, 50)
	RF433_linearize_chan_ = make(chan RFCmdToSend, 10)
	MQTT_ir_chan_ = make(chan string, 10)
}
