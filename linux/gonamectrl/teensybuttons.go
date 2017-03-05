// (c) Bernhard Tittelbach, 2016, 2017
package main

import (
	"bytes"
	"encoding/binary"
)

type NameAction struct {
	name  string
	onoff bool
}

type ButtonAction []NameAction

var name_actions_ []ButtonAction = []ButtonAction{
	ButtonAction{NameAction{"ceiling3", true}, NameAction{"couchred", true}}, ButtonAction{NameAction{"ceiling3", false}, NameAction{"couchred", false}},
	ButtonAction{NameAction{"ceiling5", true}}, ButtonAction{NameAction{"ceiling5", false}},
	ButtonAction{NameAction{"ceiling1", true}}, ButtonAction{NameAction{"ceiling1", false}},
	ButtonAction{NameAction{"ceiling2", true}, NameAction{"couchwhite", true}}, ButtonAction{NameAction{"ceiling2", false}, NameAction{"couchwhite", false}},
	ButtonAction{NameAction{"ceiling4", true}}, ButtonAction{NameAction{"ceiling4", false}},
	ButtonAction{NameAction{"ceiling6", true}}, ButtonAction{NameAction{"ceiling6", false}},
	ButtonAction{NameAction{"ceilingAll", false}, NameAction{"cxleds", false}, NameAction{"regalleinwand", false}},
	ButtonAction{NameAction{"ceiling1", false}, NameAction{"ceiling2", true}, NameAction{"ceiling3", false}, NameAction{"ceiling4", false}, NameAction{"ceiling5", true}, NameAction{"ceiling6", false}, NameAction{"cxleds", true}, NameAction{"regalleinwand", false}},
	ButtonAction{NameAction{"ceilingAll", true}, NameAction{"cxleds", true}, NameAction{"regalleinwand", true}},
}

func goListenForButtons(buttonchange_chan <-chan SerialLine) {
	for btnchange := range buttonchange_chan {
		if len(btnchange) < 3 {
			LogBTN_.Println("Did not get enought bytes from SerialLine: ", btnchange)
			continue
		}
		var buttonbits uint16
		buf := bytes.NewReader(btnchange[1:len(btnchange)]) //last 2 bytes contain button pressed information as bitfield
		err := binary.Read(buf, binary.BigEndian, &buttonbits)
		if err != nil {
			LogBTN_.Println("binary.Read failed:", err)
		}
		LogBTN_.Printf("Button State received: 0x%x", buttonbits)
		for bidx, btn_action := range name_actions_ {
			if buttonbits&(1<<uint(bidx)) > 0 {
				LogBTN_.Printf("Button %d pressed", bidx)
				for _, na := range btn_action {
					SwitchName(na.name, na.onoff)
				}
			}
		}
	}
}
