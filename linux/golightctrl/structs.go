// (c) Bernhard Tittelbach 2017
package main

import "github.com/realraum/door_and_sensors/r3events"

type wsMessage struct {
	Ctx  string                 `json:"ctx"`
	Data map[string]interface{} `json:"data"`
}

type wsMessageOut struct {
	Ctx  string      `json:"ctx"`
	Data interface{} `json:"data"`
}

type wsMsgFancyLight struct {
	Name    string              `json:"name"`
	Setting r3events.FancyLight `json:"setting"`
}

type jsonButtonUsed struct {
	Name string `json:"name"`
}

typedef RetainRecallID uint32


type JsonFuture struct {
	future   chan [][]byte
	what  []RetainRecallID
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
