// (c) Bernhard Tittelbach, 2016, 2017
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/realraum/door_and_sensors/r3events"
)

type NameAction struct {
	name  string
	onoff bool
}

const button_reset_timeout_ = time.Duration(7 * time.Second)

var (
	onaction  = []byte("{\"Action\":\"on\"}")
	offaction = []byte("{\"Action\":\"off\"}")
)

var payload_off []byte = []byte("{\"r\":0,\"g\":0,\"b\":0,\"cw\":0,\"ww\":0}")
var payload_ww1 []byte = []byte("{\"r\":0,\"g\":0,\"b\":0,\"cw\":0,\"ww\":50}")
var payload_ww2 []byte = []byte("{\"r\":0,\"g\":0,\"b\":0,\"cw\":0,\"ww\":800}")
var payload_cw1 []byte = []byte("{\"r\":0,\"g\":0,\"b\":0,\"cw\":100,\"ww\":0}")
var payload_cw2 []byte = []byte("{\"r\":0,\"g\":0,\"b\":0,\"cw\":650,\"ww\":0}")
var payload_ww3 []byte = []byte("{\"r\":0,\"g\":0,\"b\":0,\"cw\":0,\"ww\":1000}")
var payload_wwcw []byte = []byte("{\"r\":800,\"g\":0,\"b\":0,\"cw\":1000,\"ww\":1000}")
var payload_ww4 []byte = []byte("{\"r\":1000,\"g\":0,\"b\":0,\"cw\":200,\"ww\":1000}")
var payload_wwcwfade []byte = []byte("{\"r\":200,\"g\":0,\"b\":0,\"cw\":1000,\"ww\":1000,\"fade\":{\"duration\":2200}}")
var payload_c1 []byte = []byte("{\"r\":1000,\"g\":0,\"b\":0,\"cw\":0,\"ww\":0}")
var payload_c2 []byte = []byte("{\"r\":400,\"g\":0,\"b\":40,\"cw\":0,\"ww\":0}")
var payload_c3 []byte = []byte("{\"r\":0,\"g\":500,\"b\":0,\"cw\":0,\"ww\":0}")
var payload_c4 []byte = []byte("{\"r\":0,\"g\":0,\"b\":300,\"cw\":0,\"ww\":0}")

func fancytopic(light int) string {
	return r3events.TOPIC_ACTIONS + fmt.Sprintf("ceiling%d/", light) + r3events.TYPE_LIGHT
}

func payload_script(script string, value float64) []byte {
	return []byte(fmt.Sprintf("{\"script\":\"%s\",\"value\":%f}", script, value))
}

func sonofftopic(name string) string {
	return r3events.TOPIC_ACTIONS + name + "/power"
}

var fancytopic_all string = r3events.TOPIC_ACTIONS + r3events.CLIENTID_CEILINGALL + "/" + r3events.TYPE_LIGHT

type ButtonAction []ActionMQTTMsg

var name_actions_ [][]ButtonAction = [][]ButtonAction{
	//button0 up
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_ww2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_ww4}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_wwcw}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", onaction}, ActionMQTTMsg{fancytopic(3), payload_off}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_c1}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_c2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_c3}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_c4}},
	},
	//button0 down
	[]ButtonAction{ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", offaction}, ActionMQTTMsg{fancytopic(3), payload_off}}},
	//button1 up
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_ww2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_ww4}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_wwcw}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", onaction}, ActionMQTTMsg{fancytopic(2), payload_off}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_c1}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_c2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_c3}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_c4}},
	},
	//buttion1 down
	[]ButtonAction{ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction}, ActionMQTTMsg{fancytopic(2), payload_off}}},
	//button2 up
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_ww2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_ww4}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_wwcw}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", onaction}, ActionMQTTMsg{fancytopic(1), payload_off}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_c1}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_c2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_c3}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_c4}},
	},

	//button2 down
	[]ButtonAction{ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction}, ActionMQTTMsg{fancytopic(1), payload_off}}},
	//button3 up
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_ww2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_ww4}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_wwcw}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", onaction}, ActionMQTTMsg{fancytopic(4), payload_off}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_c1}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_c2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_c3}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_c4}},
	},
	//button3 down
	[]ButtonAction{ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction}, ActionMQTTMsg{fancytopic(4), payload_off}}},
	//button4 up
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_ww2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_ww4}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_wwcw}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", onaction}, ActionMQTTMsg{fancytopic(5), payload_off}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_c1}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_c2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_c3}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_c4}},
	},

	//button4 down
	[]ButtonAction{ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction}, ActionMQTTMsg{fancytopic(5), payload_off}}},
	//button5 up
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_ww2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_ww4}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_wwcw}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", onaction}, ActionMQTTMsg{fancytopic(6), payload_off}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_c1}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_c2}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_c3}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_c4}},
	},
	//button5 down
	[]ButtonAction{ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", offaction}, ActionMQTTMsg{fancytopic(6), payload_off}}},
	//button6
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "basiclightAll", offaction}, ActionMQTTMsg{fancytopic_all, payload_off}, ActionMQTTMsg{topic_lightctrl_pre_ + "regalleinwand", offaction}, ActionMQTTMsg{r3events.ACT_ACTIVATE_SCRIPT, payload_script("off", 0.0)}, ActionMQTTMsg{topic_lightctrl_pre_ + "floodtesla", offaction}, ActionMQTTMsg{sonofftopic("subtable"), []byte("off")}},
	},
	//button7
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "cxleds", onaction}},
		ButtonAction{ActionMQTTMsg{topic_lightctrl_pre_ + "cxleds", offaction}},
	},
	//button8
	[]ButtonAction{
		ButtonAction{ActionMQTTMsg{r3events.ACT_ACTIVATE_SCRIPT, payload_script("off", 0.0)}, ActionMQTTMsg{topic_lightctrl_pre_ + "basiclightAll", onaction}, ActionMQTTMsg{fancytopic_all, payload_off}},
		ButtonAction{ActionMQTTMsg{r3events.ACT_ACTIVATE_SCRIPT, payload_script("off", 0.0)},
			ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling1", offaction},
			ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling2", offaction},
			ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling3", onaction},
			ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling4", offaction},
			ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling5", offaction},
			ActionMQTTMsg{topic_lightctrl_pre_ + "ceiling6", onaction},
			ActionMQTTMsg{fancytopic(1), payload_off},
			ActionMQTTMsg{fancytopic(2), payload_cw2},
			ActionMQTTMsg{fancytopic(3), payload_off},
			ActionMQTTMsg{fancytopic(4), payload_cw2},
			ActionMQTTMsg{fancytopic(5), payload_off},
			ActionMQTTMsg{fancytopic(6), payload_off},
		},
		ButtonAction{ActionMQTTMsg{r3events.ACT_ACTIVATE_SCRIPT, payload_script("off", 0.0)}, ActionMQTTMsg{topic_lightctrl_pre_ + "basiclightAll", offaction}, ActionMQTTMsg{topic_lightctrl_pre_ + "fancyvortrag", offaction}, ActionMQTTMsg{sonofftopic("subtable"), []byte("on")}},
		ButtonAction{ActionMQTTMsg{r3events.ACT_ACTIVATE_SCRIPT, payload_script("off", 0.0)},
			ActionMQTTMsg{topic_lightctrl_pre_ + "basiclightAll", offaction},
			ActionMQTTMsg{fancytopic(1), payload_off},
			ActionMQTTMsg{fancytopic(2), payload_wwcwfade},
			ActionMQTTMsg{fancytopic(3), payload_wwcwfade},
			ActionMQTTMsg{fancytopic(4), payload_wwcwfade},
			ActionMQTTMsg{fancytopic(5), payload_wwcwfade},
			ActionMQTTMsg{fancytopic(6), payload_off},
		},
	},
}
var corresponding_btn [15]int = [15]int{-1, 0, -1, 2, -1, 4, -1, 6, -1, 8, -1, 10, -1, -1, -1}

func goListenForButtons(buttonchange_chan <-chan SerialLine) {
	action_index := make([]int, len(name_actions_))
	last_button_press := time.Now()
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
			continue
		}
		// reset button press index to 0 if long time since last button press
		if time.Now().Sub(last_button_press) > button_reset_timeout_ {
			for idx, _ := range action_index {
				action_index[idx] = 0
			}
		}
		last_button_press = time.Now()
		//reset up botton press index to 0 if corresponding down button was pressend
		for bidx, reset_other_button_press_index := range corresponding_btn {
			if buttonbits&(1<<uint(bidx)) > 0 && reset_other_button_press_index >= 0 && reset_other_button_press_index < len(action_index) {
				action_index[reset_other_button_press_index] = 0
			}
		}
		//read buttons pressed and handle
		LogBTN_.Printf("Button State received: 0x%x", buttonbits)
		for bidx, btn_action_list := range name_actions_ {
			if buttonbits&(1<<uint(bidx)) > 0 {
				LogBTN_.Printf("Button %d pressed", bidx)
				if action_index[bidx] >= len(btn_action_list) {
					action_index[bidx] = 0
				}
				for _, na := range btn_action_list[action_index[bidx]] {
					MQTT_chan_ <- na
				}
				action_index[bidx]++
			}
		}
	}
}
