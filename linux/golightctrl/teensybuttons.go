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

var name_actions_ [][]ButtonAction = [][]ButtonAction{
	//button0 up
	[]ButtonAction{
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3ww2", true}},
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3ww4", true}},
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3wwcw", true}},
		ButtonAction{NameAction{"ceiling3", true}, NameAction{"fancy3off", true}},
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3c1", true}},
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3c2", true}},
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3c3", true}},
		ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3c4", true}},
	},
	//button0 down
	[]ButtonAction{ButtonAction{NameAction{"ceiling3", false}, NameAction{"fancy3off", true}}},
	//button1 up
	[]ButtonAction{
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2ww2", true}},
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2ww4", true}},
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2wwcw", true}},
		ButtonAction{NameAction{"ceiling2", true}, NameAction{"fancy2off", true}},
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2c1", true}},
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2c2", true}},
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2c3", true}},
		ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2c4", true}},
	},
	//buttion1 down
	[]ButtonAction{ButtonAction{NameAction{"ceiling2", false}, NameAction{"fancy2off", true}}},
	//button2 up
	[]ButtonAction{
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1ww2", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1ww4", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1wwcw", true}},
		ButtonAction{NameAction{"ceiling1", true}, NameAction{"fancy1off", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1c1", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1c2", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1c3", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1c4", true}},
	},

	//button2 down
	[]ButtonAction{ButtonAction{NameAction{"ceiling1", false}, NameAction{"fancy1off", true}}},
	//button3 up
	[]ButtonAction{
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4ww2", true}},
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4ww4", true}},
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4wwcw", true}},
		ButtonAction{NameAction{"ceiling4", true}, NameAction{"fancy4off", true}},
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4c1", true}},
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4c2", true}},
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4c3", true}},
		ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4c4", true}},
	},
	//button3 down
	[]ButtonAction{ButtonAction{NameAction{"ceiling4", false}, NameAction{"fancy4off", true}}},
	//button4 up
	[]ButtonAction{
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5ww2", true}},
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5ww4", true}},
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5wwcw", true}},
		ButtonAction{NameAction{"ceiling5", true}, NameAction{"fancy5off", true}},
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5c1", true}},
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5c2", true}},
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5c3", true}},
		ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5c4", true}},
	},

	//button4 down
	[]ButtonAction{ButtonAction{NameAction{"ceiling5", false}, NameAction{"fancy5off", true}}},
	//button5 up
	[]ButtonAction{
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6ww2", true}},
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6ww4", true}},
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6wwcw", true}},
		ButtonAction{NameAction{"ceiling6", true}, NameAction{"fancy6off", true}},
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6c1", true}},
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6c2", true}},
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6c3", true}},
		ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6c4", true}},
	},
	//button5 down
	[]ButtonAction{ButtonAction{NameAction{"ceiling6", false}, NameAction{"fancy6off", true}}},
	//button6
	[]ButtonAction{ButtonAction{NameAction{"ceilingAll", false}, NameAction{"fancyalloff", true}, NameAction{"regalleinwand", false}}},
	//button7
	[]ButtonAction{
		ButtonAction{NameAction{"cxleds", true}},
		ButtonAction{NameAction{"cxleds", false}},
	},
	//button8
	[]ButtonAction{
		ButtonAction{NameAction{"ceilingAll", true}, NameAction{"fancyalloff", true}},
		ButtonAction{NameAction{"ceiling1", false}, NameAction{"ceiling2", false}, NameAction{"ceiling3", true}, NameAction{"ceiling4", false}, NameAction{"ceiling5", false}, NameAction{"ceiling6", true}, NameAction{"fancyalloff", true}, NameAction{"fancy4c2", true}, NameAction{"fancy2c2", true}},
		ButtonAction{NameAction{"ceilingAll", false}, NameAction{"fancyvortrag", true}},
	},
}

func goListenForButtons(buttonchange_chan <-chan SerialLine) {
	action_index := make([]int, len(name_actions_))
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
		for bidx, btn_action_list := range name_actions_ {
			if buttonbits&(1<<uint(bidx)) > 0 {
				LogBTN_.Printf("Button %d pressed", bidx)
				if action_index[bidx] >= len(btn_action_list) {
					action_index[bidx] = 0
				}
				for _, na := range btn_action_list[action_index[bidx]] {
					SwitchName(na.name, na.onoff)
				}
				action_index[bidx]++
			}
		}
	}
}
