// (c) Bernhard Tittelbach 2017
package main

import "github.com/realraum/door_and_sensors/r3events"

type jsonButtonUsed struct {
	Name string `json:"name"`
}

type ActionRFCode struct {
	codeon  []byte
	codeoff []byte
	handler string
}

type ActionPipePattern *r3events.SetPipeLEDsPattern

type ActionMQTTMsg struct {
	topic   string
	payload []byte
}

type ActionMeta struct {
	metaaction []string
}

type ActionIRCmdMQTT struct {
	ircmd string
}

type ActionBasicLight struct {
	light int
}

type ActionNameHandler interface{}

type RFCmdToSend struct {
	handler string
	code    []byte
}
