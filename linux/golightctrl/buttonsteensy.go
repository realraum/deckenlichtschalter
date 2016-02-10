// (c) Bernhard Tittelbach, 2016
package main

import "encoding/binary"

type NameAction struct {
	name  string
	onoff bool
}

type ButtonAction []NameAction

var name_actions_ []ButtonAction = []ButtonAction{
	ButtonAction{NameAction{"ceiling6", true}}, ButtonAction{NameAction{"ceiling6", false}},
	ButtonAction{NameAction{"ceiling4", true}}, ButtonAction{NameAction{"ceiling4", false}},
	ButtonAction{NameAction{"ceiling2", true}, NameAction{"regalleinwand", true}}, ButtonAction{NameAction{"ceiling2", false}, NameAction{"regalleinwand", false}},
	ButtonAction{NameAction{"ceiling5", true}}, ButtonAction{NameAction{"ceiling5", false}},
	ButtonAction{NameAction{"ceiling3", true}}, ButtonAction{NameAction{"ceiling3", false}},
	ButtonAction{NameAction{"ceiling1", true}}, ButtonAction{NameAction{"ceiling1", false}},
	ButtonAction{NameAction{"ceiling1", false}, NameAction{"ceiling2", false}, NameAction{"ceiling3", false}, NameAction{"ceiling4", false}, NameAction{"ceiling5", false}, NameAction{"ceiling6", false}, NameAction{"cxleds", false}, NameAction{"regalleinwand", false}},
	ButtonAction{NameAction{"ceiling1", false}, NameAction{"ceiling2", false}, NameAction{"ceiling3", true}, NameAction{"ceiling4", true}, NameAction{"ceiling5", false}, NameAction{"ceiling6", false}, NameAction{"cxleds", true}, NameAction{"regalleinwand", false}},
	ButtonAction{NameAction{"ceiling1", true}, NameAction{"ceiling2", true}, NameAction{"ceiling3", true}, NameAction{"ceiling4", true}, NameAction{"ceiling5", true}, NameAction{"ceiling6", true}, NameAction{"cxleds", true}, NameAction{"regalleinwand", true}},
}

func goListenForButtons(buttonchange_chan <-chan SerialLine) {
	for btnchange := range buttonchange_chan {
		if len(btnchange) < 3 {
			continue
		}
		buttonbits, _ := binary.Uvarint(btnchange[1:2])

		for bidx, btn_action := range name_actions_ {
			if buttonbits&(1<<uint(bidx)) > 0 {
				for _, na := range btn_action {
					SwitchName(na.name, na.onoff)
				}
			}
		}
	}
}
