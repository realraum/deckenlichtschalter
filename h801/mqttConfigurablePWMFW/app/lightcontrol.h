#ifndef LIGHTCONTROL_H
#define LIGHTCONTROL_H

const uint32_t pwm_period = 5000; // * 200ns ^= 1 kHz

extern uint32_t effect_target_values[PWM_CHANNELS];
extern uint32_t effect_intermid_values[PWM_CHANNELS];
extern uint32_t effect_interval;
extern String mqtt_forward_to;
extern String mqtt_payload;

void applyValues(uint32_t values[PWM_CHANNELS]);
void startFlash(uint8_t repetitions);
void flashSingleChannel(uint8_t repetitions, uint8_t channel);
void startFade(uint16_t steps);

#endif