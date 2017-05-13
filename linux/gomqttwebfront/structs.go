// (c) Bernhard Tittelbach 2017
package main

import "github.com/realraum/door_and_sensors/r3events"

type wsMessage struct {
	Ctx  string      `json:"ctx"`
	Data interface{} `json:"data"`
}

type HSV struct {
	H float64
	S float64
	V float64
}

type AdvFancyLightSettings struct {
	FollowDawnDusk *bool  `json:"followdawndusk"`
	WIntensity     *int64 `json:"wintensity"`
	WBalance       *int64 `json:"wbalance"`
	HSV            *HSV   `json:"hsv"`
}

type wsMsgFancyLight struct {
	Name       string                 `json:"name"`
	Setting    *r3events.FancyLight   `json:"setting,omitempty"`
	AdvSetting *AdvFancyLightSettings `json:"advsetting,omitempty"`
}

type jsonButtonUsed struct {
	Name string `json:"name"`
}

type OurFutures [][]byte

type JsonFuture struct {
	future    chan OurFutures
	omitempty bool
	what      []string
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

type MQTTOutboundMsg struct {
	topic string
	msg   interface{}
}
