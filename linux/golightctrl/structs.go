// (c) Bernhard Tittelbach 2017
package main

import "github.com/realraum/door_and_sensors/r3events"

type jsonButtonUsed struct {
	Name string `json:"name"`
}

type ActionNameHandler struct {
	handler     string
	codeon      []byte
	codeoff     []byte
	codedefault []byte
	metaaction  []string
	pipepattern *r3events.SetPipeLEDsPattern
}

type RFCmdToSend struct {
	handler string
	code    []byte
}
