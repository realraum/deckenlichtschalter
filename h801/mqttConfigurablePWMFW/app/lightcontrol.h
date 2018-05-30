#ifndef LIGHTCONTROL_H
#define LIGHTCONTROL_H

#define MAX_ALLOWED_EFFECT_DURATION 120000
#define MIN_ALLOWED_EFFECT_DURATION 5
#define MAX_ALLOWED_EFFECT_REPETITIONS 40
#define MIN_ALLOWED_EFFECT_PERIOD   5
#define MAX_ALLOWED_EFFECT_PERIOD 1500

#define DEFAULT_EFFECT_REPETITIONS 1
#define DEFAULT_EFFECT_DURATION 3200
const uint32_t DEFAULT_FLASH_PERIOD_ = 800; //ms

enum FLASHFLAGS {FLASH_INTERMED_USERSET, FLASH_INTERMED_DARK, FLASH_INTERMED_ORIG};

const uint32_t pwm_period	 = 5000; // * 200ns ^= 1 kHz

extern uint32_t effect_target_values_[PWM_CHANNELS];
extern uint32_t effect_intermid_values_[PWM_CHANNELS];
extern uint32_t button_on_values_[PWM_CHANNELS];
extern String mqtt_forward_to_;
extern String mqtt_payload_;

void enableFan(bool en);
void checkFanNeeded();
void applyValues(uint32_t values[PWM_CHANNELS]);
void startFlash(uint8_t repetitions, FLASHFLAGS intermed, uint32_t flash_period);
void flashSingleChannel(uint8_t repetitions, uint8_t channel);
void startFade(uint32_t duration_ms);
void stopAndRestoreValues(bool abort=false);

#endif