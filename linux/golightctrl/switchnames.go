// (c) Bernhard Tittelbach, 2016
package main

import (
	"fmt"
	"time"

	"github.com/btittelbach/pubsub"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/realraum/door_and_sensors/r3events"
)

const (
	IRCmd2MQTT            = "IRCmd2MQTT"
	RFCode2TTY            = "RFCode2TTY"
	RFCode2BOTH           = "RFCode2BOTH"
	RFCode2MQTT           = "RFCode2MQTT"
	LEDPattern2MQTT       = "LEDPattern2MQTT"
	CeilingLightByteState = "CeilingLightByteState"
	MetaAction            = "MetaAction"
	POST_RF433_MQTT_DELAY = 600 * time.Millisecond
	POST_RF433_TTY_DELAY  = 400 * time.Millisecond
)

var actionname_map_ map[string]ActionNameHandler = map[string]ActionNameHandler{
	//RF Power Outlets
	"regalleinwand": ActionNameHandler{codeon: []byte{0xa2, 0xa0, 0xa8}, codeoff: []byte{0xa2, 0xa0, 0x28}, handler: RFCode2TTY},  //white remote B 1
	"bluebar":       ActionNameHandler{codeon: []byte{0xa8, 0xa0, 0xa8}, codeoff: []byte{0xa8, 0xa0, 0x28}, handler: RFCode2TTY},  //white remote C 1
	"labortisch":    ActionNameHandler{codeon: []byte{0xa2, 0xa2, 0x8a}, codeoff: []byte{0xa2, 0xa2, 0x2a}, handler: RFCode2BOTH}, //polling 01000 a
	"boilerolga":    ActionNameHandler{codeon: []byte{0xa2, 0x8a, 0x8a}, codeoff: []byte{0xa2, 0x8a, 0x2a}, handler: RFCode2BOTH}, //polling 01000 b
	// "??":            ActionNameHandler{codeon: []byte{0xa2, 0x2a, 0x8a}, codeoff: []byte{0xa2, 0x2a, 0x2a}, handler: RFCode2BOTH}, //polling 01000 c
	"floodtesla": ActionNameHandler{codeon: []byte{0xa2, 0xaa, 0x88}, codeoff: []byte{0xa2, 0xaa, 0x28}, handler: RFCode2BOTH}, //polling 01000 b
	"couchred":   ActionNameHandler{codeon: []byte{0x8a, 0xa0, 0x8a}, codeoff: []byte{0x8a, 0xa0, 0x2a}, handler: RFCode2TTY},  //pollin 00101 a
	"cxleds":     ActionNameHandler{codeon: []byte{0x8a, 0x88, 0x8a}, codeoff: []byte{0x8a, 0x88, 0x2a}, handler: RFCode2TTY},  //pollin 00101 b
	"couchwhite": ActionNameHandler{codeon: []byte{0x8a, 0xa8, 0x88}, codeoff: []byte{0x8a, 0xa8, 0x28}, handler: RFCode2TTY},  //pollin 00101 d
	"mashadecke": ActionNameHandler{codeon: []byte{0x8a, 0x28, 0x8a}, codeoff: []byte{0x8a, 0x28, 0x2a}, handler: RFCode2BOTH}, //pollin 00101 c
	"boiler":     ActionNameHandler{codeon: []byte{0xa0, 0xa2, 0xa8}, codeoff: []byte{0xa0, 0xa2, 0x28}, handler: RFCode2BOTH}, //white remote A 2
	"spots":      ActionNameHandler{codeon: []byte{0x00, 0xaa, 0x88}, codeoff: []byte{0x00, 0xaa, 0x28}, handler: RFCode2TTY},  //polling 11110 d
	"abwasch":    ActionNameHandler{codeon: []byte{0xaa, 0xa2, 0xa8}, codeoff: []byte{0xaa, 0xa2, 0x28}, handler: RFCode2MQTT}, //alte jk16 decke vorne
	//rf not to be included in any, just for resetting POEarduino
	"olgatemp": ActionNameHandler{codeon: []byte{0x00, 0xa2, 0x8a}, codeoff: []byte{0x00, 0xa2, 0x2a}, handler: RFCode2TTY}, // Funksteckdose an welcher olgafreezer.realraum.at hängt

	//Yamaha IR codes
	"ymhpoweroff":  ActionNameHandler{codedefault: []byte("ymhpoweroff"), handler: IRCmd2MQTT},
	"ymhpower":     ActionNameHandler{codedefault: []byte("ymhpower"), codeoff: []byte("ymhpoweroff"), handler: IRCmd2MQTT},
	"ymhpoweron":   ActionNameHandler{codedefault: []byte("ymhpoweron"), handler: IRCmd2MQTT},
	"ymhcd":        ActionNameHandler{codedefault: []byte("ymhcd"), handler: IRCmd2MQTT},
	"ymhtuner":     ActionNameHandler{codedefault: []byte("ymhtuner"), handler: IRCmd2MQTT},
	"ymhtape":      ActionNameHandler{codedefault: []byte("ymhtape"), handler: IRCmd2MQTT},
	"ymhwdtv":      ActionNameHandler{codedefault: []byte("ymhwdtv"), handler: IRCmd2MQTT},
	"ymhsattv":     ActionNameHandler{codedefault: []byte("ymhsattv"), handler: IRCmd2MQTT},
	"ymhvcr":       ActionNameHandler{codedefault: []byte("ymhvcr"), handler: IRCmd2MQTT},
	"ymh7":         ActionNameHandler{codedefault: []byte("ymh7"), handler: IRCmd2MQTT},
	"ymhaux":       ActionNameHandler{codedefault: []byte("ymhaux"), handler: IRCmd2MQTT},
	"ymhextdec":    ActionNameHandler{codedefault: []byte("ymhextdec"), handler: IRCmd2MQTT},
	"ymhtest":      ActionNameHandler{codedefault: []byte("ymhtest"), handler: IRCmd2MQTT},
	"ymhtunabcde":  ActionNameHandler{codedefault: []byte("ymhtunabcde"), handler: IRCmd2MQTT},
	"ymheffect":    ActionNameHandler{codedefault: []byte("ymheffect"), handler: IRCmd2MQTT},
	"ymhtunplus":   ActionNameHandler{codedefault: []byte("ymhtunplus"), handler: IRCmd2MQTT},
	"ymhtunminus":  ActionNameHandler{codedefault: []byte("ymhtunminus"), handler: IRCmd2MQTT},
	"ymhvolup":     ActionNameHandler{codedefault: []byte("ymhvolup"), handler: IRCmd2MQTT},
	"ymhvoldown":   ActionNameHandler{codedefault: []byte("ymhvoldown"), handler: IRCmd2MQTT},
	"ymhvolmute":   ActionNameHandler{codedefault: []byte("ymhvolmute"), handler: IRCmd2MQTT},
	"ymhmenu":      ActionNameHandler{codedefault: []byte("ymhmenu"), handler: IRCmd2MQTT},
	"ymhplus":      ActionNameHandler{codedefault: []byte("ymhplus"), handler: IRCmd2MQTT},
	"ymhminus":     ActionNameHandler{codedefault: []byte("ymhminus"), handler: IRCmd2MQTT},
	"ymhtimelevel": ActionNameHandler{codedefault: []byte("ymhtimelevel"), handler: IRCmd2MQTT},
	"ymhprgdown":   ActionNameHandler{codedefault: []byte("ymhprgdown"), handler: IRCmd2MQTT},
	"ymhprgup":     ActionNameHandler{codedefault: []byte("ymhprgup"), handler: IRCmd2MQTT},
	"ymhsleep":     ActionNameHandler{codedefault: []byte("ymhsleep"), handler: IRCmd2MQTT},
	"ymhp5":        ActionNameHandler{codedefault: []byte("ymhp5"), handler: IRCmd2MQTT},

	//LED Pipe
	"piperainbow10":    ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "rainbow", Arg: 5, Brightness: 10, Speed: 150}, handler: LEDPattern2MQTT},
	"piperainbow30":    ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "rainbow", Arg: 5, Brightness: 30, Speed: 150}, handler: LEDPattern2MQTT},
	"piperainbow50":    ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "rainbow", Arg: 5, Brightness: 50, Speed: 150}, handler: LEDPattern2MQTT},
	"piperainbow80":    ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "rainbow", Arg: 3, Brightness: 80, Speed: 150}, handler: LEDPattern2MQTT},
	"pipeplasma":       ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "plasma", Speed: 150}, handler: LEDPattern2MQTT},
	"pipecircles":      ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "circles"}, handler: LEDPattern2MQTT},
	"pipeuspolice":     ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "uspol"}, handler: LEDPattern2MQTT},
	"pipemovingspots1": ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "movingspots", Arg: 1}, handler: LEDPattern2MQTT},
	"pipemovingspots3": ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "movingspots", Arg: 3}, handler: LEDPattern2MQTT},
	"pipemovingspots5": ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "movingspots", Arg: 4}, handler: LEDPattern2MQTT},
	"white50":          ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "white", Brightness: 50}, handler: LEDPattern2MQTT},
	"white100":         ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "white", Brightness: 100}, handler: LEDPattern2MQTT},
	"red50":            ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "hue", Hue: 0, Brightness: 50, Speed: 150}, handler: LEDPattern2MQTT},
	"green50":          ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "hue", Hue: 60, Brightness: 50, Speed: 150}, handler: LEDPattern2MQTT},
	"blue50":           ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "hue", Hue: 180, Brightness: 50, Speed: 150}, handler: LEDPattern2MQTT},
	"purple50":         ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "hue", Hue: 205, Brightness: 50, Speed: 150}, handler: LEDPattern2MQTT},
	"orange50":         ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "hue", Hue: 20, Brightness: 50, Speed: 150}, handler: LEDPattern2MQTT},
	"huefadeSS30":      ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "huefade", Speed: 20, Brightness: 30}, handler: LEDPattern2MQTT},
	"huefadeSS70":      ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "huefade", Speed: 20, Brightness: 70}, handler: LEDPattern2MQTT},
	"huefadeS30":       ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "huefade", Speed: 150, Brightness: 30}, handler: LEDPattern2MQTT},
	"huefadeS70":       ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "huefade", Speed: 150, Brightness: 70}, handler: LEDPattern2MQTT},
	"huefadeF30":       ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "huefade", Speed: 230, Brightness: 30}, handler: LEDPattern2MQTT},
	"huefadeF70":       ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "huefade", Speed: 230, Brightness: 70}, handler: LEDPattern2MQTT},
	"rstrobo":          ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "rstrobo"}, handler: LEDPattern2MQTT},
	"pipeoff":          ActionNameHandler{pipepattern: &r3events.SetPipeLEDsPattern{Pattern: "off"}, handler: LEDPattern2MQTT},
	"ceiling1":         ActionNameHandler{codeon: []byte{0, 1}, codeoff: []byte{0, 0}, handler: CeilingLightByteState},
	"ceiling2":         ActionNameHandler{codeon: []byte{1, 1}, codeoff: []byte{1, 0}, handler: CeilingLightByteState},
	"ceiling3":         ActionNameHandler{codeon: []byte{2, 1}, codeoff: []byte{2, 0}, handler: CeilingLightByteState},
	"ceiling4":         ActionNameHandler{codeon: []byte{3, 1}, codeoff: []byte{3, 0}, handler: CeilingLightByteState},
	"ceiling5":         ActionNameHandler{codeon: []byte{4, 1}, codeoff: []byte{4, 0}, handler: CeilingLightByteState},
	"ceiling6":         ActionNameHandler{codeon: []byte{5, 1}, codeoff: []byte{5, 0}, handler: CeilingLightByteState},

	//Meta Events
	"ambientlights": ActionNameHandler{handler: MetaAction, metaaction: []string{"regalleinwand", "bluebar", "couchred", "couchwhite", "abwasch", "floodtesla"}},
	"allrf":         ActionNameHandler{handler: MetaAction, metaaction: []string{"regalleinwand", "bluebar", "couchred", "couchwhite", "abwasch", "labortisch", "boiler", "boilerolga", "cxleds", "ymhpower", "floodtesla"}},
	"all":           ActionNameHandler{handler: MetaAction, metaaction: []string{"regalleinwand", "bluebar", "couchred", "couchwhite", "abwasch", "labortisch", "boiler", "boilerolga", "cxleds", "ymhpower", "floodtesla", "ceiling1", "ceiling2", "ceiling3", "ceiling4", "ceiling5", "ceiling6"}},
}

func GoSwitchNameAsync() {
FORLOOP:
	for snc := range switch_name_chan_ {
		var onoff bool
		switch snc.Action {
		case "1", "on", "send":
			onoff = true
		case "0", "off":
			onoff = false
		case "toggle":
			//TODO
		default:
			continue FORLOOP
		}
		if err := SwitchName(snc.Name, onoff); err != nil {
			LogRF433_.Println(err)
		}
	}
}

func SwitchName(name string, onoff bool) (err error) {
	nm, inmap := actionname_map_[name]
	if !inmap {
		LogRF433_.Printf("Name %s does not exist in actionname_map_", name)
		return fmt.Errorf("Name does not exist")
	}
	LogRF433_.Printf("SwitchName(%s,%t)", name, onoff)
	var code []byte
	if onoff && nm.codeon != nil {
		code = nm.codeon
	} else if onoff == false && nm.codeoff != nil {
		code = nm.codeoff
	} else if nm.codedefault != nil {
		code = nm.codedefault
	}

	switch nm.handler {
	case MetaAction:
		if len(nm.metaaction) == 0 {
			return fmt.Errorf("Could not do anything, no metaaction defined")
		}
		for _, metaname := range nm.metaaction {
			err = SwitchName(metaname, onoff)
		}
	case IRCmd2MQTT:
		if code == nil {
			return fmt.Errorf("No code for IR defined in ActionNameHandler")
		}
		err = sendIRCmd2MQTT(code)
	case RFCode2TTY, RFCode2BOTH, RFCode2MQTT:
		if code == nil {
			return fmt.Errorf("No code for RF433 defined in ActionNameHandler")
		}
		RF433_linearize_chan_ <- RFCmdToSend{handler: nm.handler, code: code}
	case CeilingLightByteState:
		if code == nil {
			return fmt.Errorf("No code for Ceiling defined in ActionNameHandler")
		}
		err = setCeilingLightByteState(code)
	case LEDPattern2MQTT:
		if nm.pipepattern == nil {
			return fmt.Errorf("No pattern for PipeLEDs defined in ActionNameHandler")
		}
		MQTT_ledpattern_chan_ <- nm.pipepattern
	default:
		return fmt.Errorf("Unknown handler %s", nm.handler)
	}

	if err != nil {
		return
	}

	//notify Everyone
	switch nm.handler {
	case MetaAction, IRCmd2MQTT, RFCode2TTY, RFCode2BOTH, RFCode2MQTT:
		ps_.PubNonBlocking(jsonButtonUsed{name}, PS_IRRF433_CHANGED)
	case CeilingLightByteState:
		ps_.PubNonBlocking(ConvertCeilingLightsStateTomap(GetCeilingLightsState(), 1), PS_LIGHTS_CHANGED)
	}

	return
}

func sendIRCmd2MQTT(code []byte) error {
	LogRF433_.Printf("IRCmd2MQTT(%s)", string(code))
	MQTT_ir_chan_ <- string(code)
	return nil
}

func setCeilingLightByteState(code []byte) error {
	if len(code) != 2 {
		LogRF433_.Printf("Invalid Code %s for setCeilingLightByteState", code)
		return fmt.Errorf("Invalid Code for setCeilingLightByteState")
	}
	SetCeilingLightsState(int(code[0]), code[1] == 1)
	return nil
}

func goLinearizeRFSenders(ps *pubsub.PubSub, rfchan <-chan RFCmdToSend, rf433_tty_chan_ chan SerialLine, mqttc mqtt.Client) {
	shutdown1_c := ps.SubOnce(PS_SHUTDOWN)
	shutdown2_c := ps.SubOnce(PS_SHUTDOWN_CONSUMER)

	for {
		select {
		case rfcmd := <-rfchan:
			switch rfcmd.handler {
			case RFCode2TTY:
				rf433_tty_chan_ <- append([]byte(">"), rfcmd.code...)
				time.Sleep(POST_RF433_TTY_DELAY)
			case RFCode2BOTH:
				sendCodeToMQTT(mqttc, rfcmd.code)
				time.Sleep(POST_RF433_MQTT_DELAY)
				rf433_tty_chan_ <- append([]byte(">"), rfcmd.code...)
				time.Sleep(POST_RF433_TTY_DELAY)
			case RFCode2MQTT:
				sendCodeToMQTT(mqttc, rfcmd.code)
				time.Sleep(POST_RF433_MQTT_DELAY)
			}
		case <-shutdown1_c:
			return
		case <-shutdown2_c:
			return
		}
	}
}
