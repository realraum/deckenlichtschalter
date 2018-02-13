// (c) Bernhard Tittelbach, 2016
package main

import bbhw "github.com/btittelbach/go-bbhw"

func CeilinglightsGPIO_FakeGPIOinit() CeilingLightsSwitchGPIO {
	LogMain_.Print("FAKE GPIO/PWM init start")
	bbhw.FakeGPIODefaultLogTarget_ = LogGPIO_
	gpios_ceiling_lights := []bbhw.GPIOControllablePin{
		bbhw.NewFakeGPIO(23, bbhw.OUT),
		bbhw.NewFakeGPIO(17, bbhw.OUT),
		bbhw.NewFakeGPIO(21, bbhw.OUT),
		bbhw.NewFakeGPIO(22, bbhw.OUT),
		bbhw.NewFakeGPIO(18, bbhw.OUT),
		bbhw.NewFakeGPIO(4, bbhw.OUT),
	}
	for _, gpio := range gpios_ceiling_lights {
		gpio.SetActiveLow(true)
	}
	LogMain_.Print("FAKE GPIO init done")
	return gpios_ceiling_lights
}

func CeilinglightsGPIO_GPIOinit() CeilingLightsSwitchGPIO {
	LogMain_.Print("GPIO/PWM init start")
	gpios_ceiling_lights := []bbhw.GPIOControllablePin{
		bbhw.NewSysfsGPIOOrPanic(23, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(17, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(21, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(22, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(18, bbhw.OUT),
		bbhw.NewSysfsGPIOOrPanic(4, bbhw.OUT),
	}
	for _, gpio := range gpios_ceiling_lights {
		gpio.SetActiveLow(true)
	}
	LogMain_.Print("GPIO init done")
	return gpios_ceiling_lights
}

func (ceiling_lights CeilingLightsSwitchGPIO) GetCeilingLightsStates() []bool {
	rv := make([]bool, len(ceiling_lights))
	for i, gpio := range ceiling_lights {
		rv[i] = bbhw.GetStateOrPanic(gpio)
	}
	return rv
}

func (ceiling_lights CeilingLightsSwitchGPIO) SetCeilingLightsState(ceiling_light_number int, onoff bool) {
	if ceiling_light_number < 0 || ceiling_light_number >= len(ceiling_lights) {
		return
	}
	ceiling_lights[ceiling_light_number].SetState(onoff)
}

func (ceiling_lights CeilingLightsSwitchGPIO) SetCeilingLightsStates(states []bool) {
	for i, state := range states {
		ceiling_lights.SetCeilingLightsState(i, state)
	}
}
