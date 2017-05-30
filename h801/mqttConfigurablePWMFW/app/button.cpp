#include <user_config.h>
#include <SmingCore/SmingCore.h>
#include "button.h"

DebouncedButton::DebouncedButton(int8_t button_pin, uint32_t debounce_duration_ms, uint32_t longpress_duration_ms, bool pulls_to_ground)
{
	button_debounce_duration_ = debounce_duration_ms;
	longpress_duration_ms_ = max(3*debounce_duration_ms,longpress_duration_ms);
	button_pin_ = button_pin;
	button_pulls_to_ground_ = pulls_to_ground;
	if (pulls_to_ground) {
		pinMode(button_pin, INPUT_PULLUP);
	} else {
		pinMode(button_pin, INPUT);
	}
	attachInterrupt(button_pin, ButtonInterruptDelegate(&DebouncedButton::buttonInterruptHandler,this), CHANGE);
}

DebouncedButton::~DebouncedButton()
{
	detachInterrupt(button_pin_);
}

bool DebouncedButton::isPressed()
{
	return (digitalRead(button_pin_) == HIGH) ^ button_pulls_to_ground_;
}

bool DebouncedButton::isLongPressed()
{
	if (((digitalRead(button_pin_) == HIGH) ^ button_pulls_to_ground_) && millis() - last_button_event_time_ >= longpress_duration_ms_)
	{
		return true;
	}
	return false;
}

void IRAM_ATTR DebouncedButton::buttonInterruptHandler()
{
	if (isPressed())
		last_button_event_time_ = millis();
	else {
		//debugf("pulsedur: %dms",millis() - last_button_event_time_);
		uint32_t btn_dur = millis() - last_button_event_time_;
		if (btn_dur > button_debounce_duration_ && btn_dur < longpress_duration_ms_) {
			button_event_ctr_++;
		}
	}
}

bool DebouncedButton::wasPressed()
{
	/*uint32_t btn_dur = millis() - last_button_event_time;
	if (isPressed() && button_event_ctr_ == last_read_button_event_ctr_ && btn_dur > button_debounce_duration_ && btn_dur < longpress_duration_ms_)
	{
		last_read_button_event_ctr_ = button_event_ctr_ + 1; //ignore next button press, we reported the one currently in progress
		return true;
	} else */
	if (button_event_ctr_ > last_read_button_event_ctr_)
	{
		last_read_button_event_ctr_ = button_event_ctr_;
		return true;
	}
	return false;
}

