// (c) Bernhard Tittelbach, 2016
package main

import "bytes"

const num_basicctrl_relays_ = 8

func NewBasicCtrl(rd, wr chan SerialLine) *BasicCtrlBox {
	//create and init new struct
	bcb := &BasicCtrlBox{
		rd:    rd,
		wr:    wr,
		state: make([]bool, num_basicctrl_relays_),
	}
	//start thread that updates state from tty
	go func() {
		//update state from rd
		for line := range rd {
			newstate := make([]bool, num_basicctrl_relays_)
			for i, chr := range line {
				if i >= num_basicctrl_relays_ {
					break
				}
				switch chr {
				case '0':
					newstate[i] = true
				case '1':
					newstate[i] = false
				default:
					bcb.state_mutex.RLock()
					newstate[i] = bcb.state[i]
					bcb.state_mutex.RUnlock()
				}
			}
			bcb.state_mutex.Lock()
			bcb.state = newstate
			bcb.state_mutex.Unlock()
		}
	}()
	//query current state on startup
	wr <- append(bytes.Repeat([]byte{'-'}, num_basicctrl_relays_), '\r', '\n')
	//return new struct
	return bcb
}

func (ceiling_lights *BasicCtrlBox) GetCeilingLightsStates() []bool {
	state := make([]bool, num_basicctrl_relays_)
	ceiling_lights.state_mutex.RLock()
	copy(state, ceiling_lights.state)
	ceiling_lights.state_mutex.RUnlock()
	return state
}

func (ceiling_lights *BasicCtrlBox) SetCeilingLightsState(ceiling_light_number int, onoff bool) {
	if ceiling_light_number < 0 || ceiling_light_number >= num_basicctrl_relays_ {
		return
	}
	bs := bytes.Repeat([]byte{'-'}, num_basicctrl_relays_)
	if onoff {
		bs[ceiling_light_number] = '1'
	} else {
		bs[ceiling_light_number] = '0'
	}
	ceiling_lights.wr <- SerialLine(append(bs, '\r', '\n'))
}

func (ceiling_lights *BasicCtrlBox) SetCeilingLightsStates(state []bool) {
	if len(state) >= num_basicctrl_relays_ {
		return
	}
	bs := make([]byte, len(state))
	for i, s := range state {
		if s {
			bs[i] = '1'
		} else {
			bs[i] = '0'
		}
	}
	ceiling_lights.wr <- SerialLine(append(bs, '\r', '\n'))
}
