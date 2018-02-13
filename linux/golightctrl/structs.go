// (c) Bernhard Tittelbach 2017
package main

import (
	"sync"

	bbhw "github.com/btittelbach/go-bbhw"
	"github.com/realraum/door_and_sensors/r3events"
)

type CeilingLightsSwitch interface {
	GetCeilingLightsStates() []bool
	SetCeilingLightsState(ceiling_light_number int, onoff bool)
	SetCeilingLightsStates([]bool)
}

type SerialLine []byte

type CeilingLightStateMap map[string]bool

type CeilingLightsSwitchGPIO []bbhw.GPIOControllablePin
type CeilingLightsSwitchBasicCtrl BasicCtrlBox

type BasicCtrlBox struct {
	rd          chan SerialLine
	wr          chan SerialLine
	state       []bool
	state_mutex sync.RWMutex
}

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
