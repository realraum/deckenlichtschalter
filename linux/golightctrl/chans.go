// (c) Bernhard Tittelbach 2017
package main

const (
	PS_WEBSOCK_ALL_JSON   = "websock_toall_json"
	PS_WEBSOCK_ALL        = "websock_toall"
	PS_LIGHTS_CHANGED     = "light_state_changed"
	PS_IRRF433_CHANGED    = "stateless_button_send_event"
	PS_FANCYLIGHT_CHANGED = "fancylight_update"
	PS_LEDPIPE_CHANGED    = "ledpipe_update"
	PS_SHUTDOWN           = "shutdown"
	PS_SHUTDOWN_CONSUMER  = "shutdownindiscriminateconsumer"
)

var MQTT_sendmsg_chan_ chan MQTTOutboundMsg

// var switch_name_chan_ chan r3events.LightCtrlActionOnName
// var MQTT_ir_chan_ chan string
// var MQTT_ledpattern_chan_ chan *r3events.SetPipeLEDsPattern
// var MQTT_fancylight_chan_ chan *wsMsgFancyLight
// var RF433_linearize_chan_ chan RFCmdToSend
// var AdvFancyLight_chan_ chan *wsMsgFancyLight

func init() {
	MQTT_sendmsg_chan_ = make(chan MQTTOutboundMsg, 50)
	// switch_name_chan_ = make(chan r3events.LightCtrlActionOnName, 50)
	// RF433_linearize_chan_ = make(chan RFCmdToSend, 10)
	// MQTT_ir_chan_ = make(chan string, 10)
	// MQTT_ledpattern_chan_ = make(chan *r3events.SetPipeLEDsPattern, 5)
	// MQTT_fancylight_chan_ = make(chan *wsMsgFancyLight, 3)
	// AdvFancyLight_chan_ = make(chan *wsMsgFancyLight, 3)
}
