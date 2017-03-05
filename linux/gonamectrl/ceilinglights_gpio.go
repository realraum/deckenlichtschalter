// (c) Bernhard Tittelbach, 2016
package main

import (
	"fmt"

	bbhw "github.com/btittelbach/go-bbhw"
)

type CeilingLightStateMap map[string]bool

var (
	gpios_ceiling_lights_ []bbhw.GPIOControllablePin
)

func FakeGPIOinit() {
	LogMain_.Print("FAKE GPIO/PWM init start")
	bbhw.FakeGPIODefaultLogTarget_ = LogGPIO_
	gpios_ceiling_lights_ = []bbhw.GPIOControllablePin{
		bbhw.NewFakeGPIO(23, bbhw.OUT),
		bbhw.NewFakeGPIO(22, bbhw.OUT),
		bbhw.NewFakeGPIO(21, bbhw.OUT),
		bbhw.NewFakeGPIO(18, bbhw.OUT),
		bbhw.NewFakeGPIO(17, bbhw.OUT),
		bbhw.NewFakeGPIO(4, bbhw.OUT),
	}
	for _, gpio := range gpios_ceiling_lights_ {
		gpio.SetActiveLow(true)
	}
	LogMain_.Print("FAKE GPIO init done")
}

func GPIOinit() {
	LogMain_.Print("GPIO/PWM init start")
	gpios_ceiling_lights_ = []bbhw.GPIOControllablePin{
		bbhw.NewSysfsGPIOOrPanic(23, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(22, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(21, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(18, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(17, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(4, bbhw.OUT),
	}
	for _, gpio := range gpios_ceiling_lights_ {
		gpio.SetActiveLow(true)
	}
	LogMain_.Print("GPIO init done")
}

func GetCeilingLightsState() []bool {
	rv := make([]bool, len(gpios_ceiling_lights_))
	for i, gpio := range gpios_ceiling_lights_ {
		rv[i] = bbhw.GetStateOrPanic(gpio)
	}
	return rv
}

func ConvertCeilingLightsStateTomap(states []bool, offset int) CeilingLightStateMap {
	rv := make(map[string]bool, 6)
	for i, st := range states {
		lightname := fmt.Sprintf("ceiling%d", i+offset)
		rv[lightname] = st
	}
	return rv
}

func SetCeilingLightsState(ceiling_light_number int, onoff bool) {
	if ceiling_light_number < 0 || ceiling_light_number >= len(gpios_ceiling_lights_) {
		return
	}
	gpios_ceiling_lights_[ceiling_light_number].SetState(onoff)
}
