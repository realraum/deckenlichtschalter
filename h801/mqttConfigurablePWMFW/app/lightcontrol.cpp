#include <SmingCore/SmingCore.h>
#include <defaultlightconfig.h>
#include <pwmchannels.h>
#include "lightcontrol.h"


enum effect_id_type {EFF_NO, EFF_FLASH, EFF_FADE};

Timer flashTimer;
effect_id_type effect = EFF_NO;
uint32_t apply_last_values[PWM_CHANNELS];
uint32_t active_values[PWM_CHANNELS];
uint32_t steps_left=0;
int32_t fade_diff_values[PWM_CHANNELS];
uint32_t effect_target_values[PWM_CHANNELS]={0,0,0,0,0};
uint32_t effect_intermid_values[PWM_CHANNELS]={0,0,0,0,0};
uint32_t effect_interval=800;
String mqtt_forward_to;
String mqtt_payload;
////////////////////
//// PWM Stuff ////
///////////////////

//init PWM and restore stored pwm values
void setupPWM()
{
	// PWM setup
	uint32 io_info[PWM_CHANNELS][3] = {
		// MUX, FUNC, PIN
		{PERIPHS_IO_MUX_MTDO_U,  FUNC_GPIO15, 15}, //R-
		{PERIPHS_IO_MUX_MTCK_U,  FUNC_GPIO13, 13}, //G-
		{PERIPHS_IO_MUX_MTDI_U,  FUNC_GPIO12, 12}, //B-
		{PERIPHS_IO_MUX_MTMS_U,  FUNC_GPIO14, 14}, //W1-
		{PERIPHS_IO_MUX_GPIO4_U, FUNC_GPIO4 ,  4}, //W2-
	};
	// uint32 pwm_duty_initial[PWM_CHANNELS] = {0, 0, 0, 0, 0};
	uint32 pwm_duty_initial[PWM_CHANNELS] = {0, 0, 0, 0, 0};

	DefaultLightConfig.load(pwm_duty_initial); //load initial default values

	pwm_init(pwm_period, pwm_duty_initial, PWM_CHANNELS, io_info);
	pwm_start();
}


/////////////////////
//// Light Stuff ////
/////////////////////

void saveCurrentValues()
{
	for (uint8_t i=0;i<PWM_CHANNELS;i++)
		apply_last_values[i] = pwm_get_duty(i);
}

void applyValues(uint32_t values[PWM_CHANNELS])
{
	for (uint8_t i=0;i<PWM_CHANNELS;i++)
		pwm_set_duty(values[i],i);
	pwm_start();
}

void stopAndRestoreValues()
{
	flashTimer.stop();
	applyValues(apply_last_values);
	effect = EFF_NO;
}

void timerFuncShowFlashEffect()
{
	if (steps_left > 0)
	{
		if (steps_left % 2 == 0)
		{
			applyValues(active_values);
		} else {
			applyValues(effect_intermid_values);
		}
		steps_left--;
	} else {
		stopAndRestoreValues();
	}
}

void timerFuncShowFadeEffect()
{
	if (steps_left > 0)
	{
		for (uint8_t i=0; i<PWM_CHANNELS; i++)
			active_values[i] += fade_diff_values[i];
		applyValues(active_values);			
		steps_left--;
	} else {
		stopAndRestoreValues();
	}
}

void timerFuncShowEffect()
{
	switch (effect)
	{
		default:
		case EFF_NO:
			flashTimer.stop();
			break;
		case EFF_FLASH:
			timerFuncShowFlashEffect();
			break;
		case EFF_FADE:
			timerFuncShowFadeEffect();
			break;
	}
}

void startFlash(uint8_t repetitions)
{
	if (effect != EFF_NO)
		stopAndRestoreValues();
	saveCurrentValues();
	memcpy(active_values, effect_target_values, PWM_CHANNELS*sizeof(uint32_t));
	applyValues(effect_intermid_values);
	steps_left = repetitions * 2;
	effect = EFF_FLASH;
	flashTimer.initializeMs(effect_interval, timerFuncShowFlashEffect).start();
}

void flashSingleChannel(uint8_t times, uint8_t channel)
{
	if (channel >= PWM_CHANNELS)
		return;
	if (effect != EFF_NO)
		stopAndRestoreValues();
	saveCurrentValues();
	memcpy(effect_target_values, apply_last_values, PWM_CHANNELS*sizeof(uint32_t));
	memcpy(effect_intermid_values, apply_last_values, PWM_CHANNELS*sizeof(uint32_t));
	effect_target_values[channel] = pwm_period/3;
	effect_intermid_values[channel] = 0;
	startFlash(times);
}

void startFade(uint16_t steps)
{
	if (effect != EFF_NO)
		stopAndRestoreValues();
	saveCurrentValues();		
	for (uint8_t i=0; i<PWM_CHANNELS; i++)
	{
		fade_diff_values[i] = ((int32_t)effect_target_values[i] - (int32_t)apply_last_values[i]) / (int32_t)steps;
	}
	steps_left = steps;

	//copy active_values to apply_last_values so active_values get applied on stop or interrupt of effect
	memcpy(apply_last_values, effect_target_values, PWM_CHANNELS*sizeof(uint32_t));
	memcpy(active_values, effect_target_values, PWM_CHANNELS*sizeof(uint32_t));
	effect = EFF_FADE;
	flashTimer.initializeMs(effect_interval, timerFuncShowFadeEffect).start();
}
