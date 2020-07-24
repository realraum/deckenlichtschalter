// (c) Bernhard Tittelbach, 2016
package main

import (
	"fmt"
	"time"

	"github.com/btittelbach/pubsub"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	IRCmd2MQTT            = "IRCmd2MQTT"
	RFCode2TTY            = "RFCode2TTY"
	RFCode2BOTH           = "qRFCode2BOTH"
	RFCode2MQTT           = "RFCode2MQTT"
	LEDPattern2MQTT       = "LEDPattern2MQTT"
	CeilingLightByteState = "CeilingLightByteState"
	MetaAction            = "MetaAction"
	POST_RF433_MQTT_DELAY = 600 * time.Millisecond
	POST_RF433_TTY_DELAY  = 400 * time.Millisecond
)

var actionname_map_ map[string]ActionNameHandler = map[string]ActionNameHandler{
	//RF Power Outlets
	"regalleinwand": ActionRFCode{codeon: []byte{0xa2, 0xa0, 0xa8}, codeoff: []byte{0xa2, 0xa0, 0x28}, handler: RFCode2TTY},  //white remote B 1
	"bluebar":       ActionRFCode{codeon: []byte{0xa8, 0xa0, 0xa8}, codeoff: []byte{0xa8, 0xa0, 0x28}, handler: RFCode2TTY},  //white remote C 1
	"laserball":     ActionRFCode{codeon: []byte{0xa2, 0xa2, 0x8a}, codeoff: []byte{0xa2, 0xa2, 0x2a}, handler: RFCode2BOTH}, //polling 01000 a
	// "??":            ActionNameHandler{codeon: []byte{0xa2, 0x2a, 0x8a}, codeoff: []byte{0xa2, 0x2a, 0x2a}, handler: RFCode2BOTH}, //polling 01000 c
	"floodtesla": ActionRFCode{codeon: []byte{0xa2, 0xaa, 0x88}, codeoff: []byte{0xa2, 0xaa, 0x28}, handler: RFCode2BOTH}, //polling 01000 b
	"logo":       ActionRFCode{codeon: []byte{0x8a, 0xa0, 0x8a}, codeoff: []byte{0x8a, 0xa0, 0x2a}, handler: RFCode2TTY},  //pollin 00101 a
	"couchred":   ActionMQTTMsg{sonofftopic("couchred"), []byte("off")},
	"cxleds":     ActionRFCode{codeon: []byte{0x8a, 0x88, 0x8a}, codeoff: []byte{0x8a, 0x88, 0x2a}, handler: RFCode2TTY}, //pollin 00101 b
	"couchwhite": ActionRFCode{codeon: []byte{0x8a, 0xa8, 0x88}, codeoff: []byte{0x8a, 0xa8, 0x28}, handler: RFCode2TTY}, //pollin 00101 d
	// "mashadecke": ActionRFCode{codeon: []byte{0x8a, 0x28, 0x8a}, codeoff: []byte{0x8a, 0x28, 0x2a}, handler: RFCode2BOTH}, //pollin 00101 c
	"boilerolga": ActionRFCode{codeon: []byte{0xa0, 0xa2, 0xa8}, codeoff: []byte{0xa0, 0xa2, 0x28}, handler: RFCode2BOTH}, //white remote A 2
	"spots":      ActionRFCode{codeon: []byte{0x00, 0xaa, 0x88}, codeoff: []byte{0x00, 0xaa, 0x28}, handler: RFCode2TTY},  //polling 11110 d
	"abwasch":    ActionRFCode{codeon: []byte{0xaa, 0xa2, 0xa8}, codeoff: []byte{0xaa, 0xa2, 0x28}, handler: RFCode2BOTH}, //alte jk16 decke vorne
	//rf not to be included in any, just for resetting POEarduino
	"olgatemp": ActionRFCode{codeon: []byte{0x00, 0xa2, 0x8a}, codeoff: []byte{0x00, 0xa2, 0x2a}, handler: RFCode2TTY}, // Funksteckdose an welcher olgafreezer.realraum.at hängt

	//Yamaha IR codes
	"ymhpoweroff":  ActionIRCmdMQTT{ircmd: "ymhpoweroff"},
	"ymhpower":     ActionIRCmdMQTT{ircmd: "ymhpower"},
	"ymhpoweron":   ActionIRCmdMQTT{ircmd: "ymhpoweron"},
	"ymhcd":        ActionIRCmdMQTT{ircmd: "ymhcd"},
	"ymhtuner":     ActionIRCmdMQTT{ircmd: "ymhtuner"},
	"ymhtape":      ActionIRCmdMQTT{ircmd: "ymhtape"},
	"ymhwdtv":      ActionIRCmdMQTT{ircmd: "ymhwdtv"},
	"ymhsattv":     ActionIRCmdMQTT{ircmd: "ymhsattv"},
	"ymhvcr":       ActionIRCmdMQTT{ircmd: "ymhvcr"},
	"ymh7":         ActionIRCmdMQTT{ircmd: "ymh7"},
	"ymhaux":       ActionIRCmdMQTT{ircmd: "ymhaux"},
	"ymhextdec":    ActionIRCmdMQTT{ircmd: "ymhextdec"},
	"ymhtest":      ActionIRCmdMQTT{ircmd: "ymhtest"},
	"ymhtunabcde":  ActionIRCmdMQTT{ircmd: "ymhtunabcde"},
	"ymheffect":    ActionIRCmdMQTT{ircmd: "ymheffect"},
	"ymhtunplus":   ActionIRCmdMQTT{ircmd: "ymhtunplus"},
	"ymhtunminus":  ActionIRCmdMQTT{ircmd: "ymhtunminus"},
	"ymhvolup":     ActionIRCmdMQTT{ircmd: "ymhvolup"},
	"ymhvoldown":   ActionIRCmdMQTT{ircmd: "ymhvoldown"},
	"ymhvolmute":   ActionIRCmdMQTT{ircmd: "ymhvolmute"},
	"ymhmenu":      ActionIRCmdMQTT{ircmd: "ymhmenu"},
	"ymhplus":      ActionIRCmdMQTT{ircmd: "ymhplus"},
	"ymhminus":     ActionIRCmdMQTT{ircmd: "ymhminus"},
	"ymhtimelevel": ActionIRCmdMQTT{ircmd: "ymhtimelevel"},
	"ymhprgdown":   ActionIRCmdMQTT{ircmd: "ymhprgdown"},
	"ymhprgup":     ActionIRCmdMQTT{ircmd: "ymhprgup"},
	"ymhsleep":     ActionIRCmdMQTT{ircmd: "ymhsleep"},
	"ymhp5":        ActionIRCmdMQTT{ircmd: "ymhp5"},

	//Basic Light
	"ceiling1":    ActionBasicLight{0},
	"ceiling2":    ActionBasicLight{1},
	"ceiling3":    ActionBasicLight{2},
	"ceiling4":    ActionBasicLight{3},
	"ceiling5":    ActionBasicLight{4},
	"ceiling6":    ActionBasicLight{5},
	"basiclight1": ActionBasicLight{0},
	"basiclight2": ActionBasicLight{1},
	"basiclight3": ActionBasicLight{2},
	"basiclight4": ActionBasicLight{3},
	"basiclight5": ActionBasicLight{4},
	"basiclight6": ActionBasicLight{5},

	//Fancy Light
	"fancyalloff": ActionMQTTMsg{fancytopic_all, payload_off},
	"fancy1off":   ActionMQTTMsg{fancytopic(1), payload_off},
	"fancy2off":   ActionMQTTMsg{fancytopic(2), payload_off},
	"fancy3off":   ActionMQTTMsg{fancytopic(3), payload_off},
	"fancy4off":   ActionMQTTMsg{fancytopic(4), payload_off},
	"fancy5off":   ActionMQTTMsg{fancytopic(5), payload_off},
	"fancy6off":   ActionMQTTMsg{fancytopic(6), payload_off},
	"fancy7off":   ActionMQTTMsg{fancytopic(7), payload_off},
	"fancy8off":   ActionMQTTMsg{fancytopic(8), payload_off},
	"fancy9off":   ActionMQTTMsg{fancytopic(9), payload_off},
	"fancy1cw1":   ActionMQTTMsg{fancytopic(1), payload_cw1},
	"fancy2cw1":   ActionMQTTMsg{fancytopic(2), payload_cw1},
	"fancy3cw1":   ActionMQTTMsg{fancytopic(3), payload_cw1},
	"fancy4cw1":   ActionMQTTMsg{fancytopic(4), payload_cw1},
	"fancy5cw1":   ActionMQTTMsg{fancytopic(5), payload_cw1},
	"fancy6cw1":   ActionMQTTMsg{fancytopic(6), payload_cw1},
	"fancy7cw1":   ActionMQTTMsg{fancytopic(7), payload_cw1},
	"fancy8cw1":   ActionMQTTMsg{fancytopic(8), payload_cw1},
	"fancy9cw1":   ActionMQTTMsg{fancytopic(9), payload_cw1},
	"fancy1cw2":   ActionMQTTMsg{fancytopic(1), payload_cw2},
	"fancy2cw2":   ActionMQTTMsg{fancytopic(2), payload_cw2},
	"fancy3cw2":   ActionMQTTMsg{fancytopic(3), payload_cw2},
	"fancy4cw2":   ActionMQTTMsg{fancytopic(4), payload_cw2},
	"fancy5cw2":   ActionMQTTMsg{fancytopic(5), payload_cw2},
	"fancy6cw2":   ActionMQTTMsg{fancytopic(6), payload_cw2},
	"fancy7cw2":   ActionMQTTMsg{fancytopic(7), payload_cw2},
	"fancy8cw2":   ActionMQTTMsg{fancytopic(8), payload_cw2},
	"fancy9cw2":   ActionMQTTMsg{fancytopic(9), payload_cw2},

	//Meta Events
	"ambientlights": ActionMeta{metaaction: []string{"regalleinwand", "bluebar", "couchred", "couchwhite", "abwasch", "floodtesla"}},
	"ceilingAll":    ActionMeta{metaaction: []string{"basiclightAll"}},
	"basiclightAll": ActionMeta{metaaction: []string{"ceiling1", "ceiling2", "ceiling3", "ceiling4", "ceiling5", "ceiling6"}},
	"fancyvortrag":  ActionMeta{metaaction: []string{"fancy1off", "fancy6off", "fancy2cw1", "fancy5cw1", "fancy3cw2", "fancy4cw2", "floodtesla"}},
	"allrf":         ActionMeta{metaaction: []string{"regalleinwand", "bluebar", "couchred", "couchwhite", "abwasch", "labortisch", "boiler", "boilerolga", "cxleds", "laserball", "logo", "ymhpower", "floodtesla"}},
	"all":           ActionMeta{metaaction: []string{"regalleinwand", "bluebar", "couchred", "couchwhite", "abwasch", "labortisch", "boiler", "boilerolga", "cxleds", "ymhpower", "floodtesla", "laserball", "logo", "subtable", "ceiling1", "ceiling2", "ceiling3", "ceiling4", "ceiling5", "ceiling6"}},
}

func ConvertCeilingLightsStateTomap(states []bool, offset int) CeilingLightStateMap {
	rv := make(map[string]bool, 6)
	for i, st := range states {
		lightname := fmt.Sprintf("ceiling%d", i+offset)
		rv[lightname] = st
	}
	return rv
}

func GoSwitchNameAsync() {
FORLOOP:
	for snc := range switch_name_chan_ {
		var onoff bool
		switch snc.Action {
		case "1", "on", "send", "{\"Action\":1}", "{\"Action\":\"on\"}", "{\"Action\":\"send\"}":
			onoff = true
		case "0", "off", "{\"Action\":0}", "{\"Action\":\"off\"}":
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
	nmi, inmap := actionname_map_[name]
	if !inmap {
		LogRF433_.Printf("Name %s does not exist in actionname_map_", name)
		return fmt.Errorf("Name does not exist")
	}
	LogRF433_.Printf("SwitchName(%s,%t)", name, onoff)

	switch nm := nmi.(type) {
	case ActionRFCode:
		var code []byte
		if onoff && nm.codeon != nil {
			code = nm.codeon
		} else if onoff == false && nm.codeoff != nil {
			code = nm.codeoff
		}
		if code == nil || len(code) != 3 {
			return fmt.Errorf("No valid code for action (%s,%t)", name, onoff)
		}
		RF433_linearize_chan_ <- RFCmdToSend{handler: nm.handler, code: code}

	case ActionMQTTMsg:
		MQTT_chan_ <- nm

	case ActionBasicLight:
		CeilingLightsSwitch_.SetCeilingLightsState(nm.light, onoff)

	case ActionMeta:
		if len(nm.metaaction) == 0 {
			return fmt.Errorf("Could not do anything, no metaaction defined")
		}
		for _, metaname := range nm.metaaction {
			err = SwitchName(metaname, onoff)
		}

	case ActionIRCmdMQTT:
		MQTT_ir_chan_ <- nm.ircmd
	}

	if err != nil {
		return
	}

	//notify Everyone
	switch nmi.(type) {
	case ActionMeta, ActionIRCmdMQTT, ActionRFCode:
		ps_.PubNonBlocking(jsonButtonUsed{name}, PS_IRRF433_CHANGED)
	case ActionBasicLight:
		ps_.PubNonBlocking(ConvertCeilingLightsStateTomap(CeilingLightsSwitch_.GetCeilingLightsStates(), 1), PS_LIGHTS_CHANGED)
	}

	return
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
