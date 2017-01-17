/*
 *  lightctrl - teensy firmware code
 *
 *
 *  Copyright (C) 2013 Bernhard Tittelbach <xro@realraum.at>
 *   uses avr-utils, anyio & co by Christian Pointner <equinox@spreadspace.org>
 *
 *  This file is part of lightctrl.
 *
 *  lightctrl is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  any later version.
 *
 *  lightctrl is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with lightctrl. If not, see <http://www.gnu.org/licenses/>.
 */

#include "mypins.h"
#include "rf433.h"

#include <avr/interrupt.h>

#define TIMER_RUNNING (TIMSK1 & (1<<OCIE1A))

void rf433_start_timer(void)
{
  // timer 1: 2 ms
  TCCR1A = 0;                    // prescaler 1:8, WGM = 4 (CTC)
  TCCR1B = 1<<WGM12 | 1<<CS11;   //
//  OCR1A = 39;        // (1+39)*8 = 320 -> 0.02ms @ 16 MHz -> 1*alpha
//default: alpha=0.08
  OCR1A = 159;        // (1+159)*8 = 1280 -> 0.08ms @ 16 MHz -> 1*alpha
//  OCR1A = 154;        // (1+154)*8 = 1240 -> 0.0775ms @ 16 MHz -> 1*alpha
//  OCR1A = 207;        // (1+207)*8 = 1664 -> 0.104ms @ 16 MHz -> 1*alpha
  TCNT1 = 0;          // reseting timer
  TIMSK1 = 1<<OCIE1A; // enable Interrupt
}

void rf433_stop_timer(void) // stop the timer
{
  // timer1
  TCCR1B = 0; // no clock source
  TIMSK1 = 0; // disable timer interrupt
}

#define NUM_REPEAT_SIGNAL 8
#define RF_SIGNAL_BYTES 3
#define RF_SIGNAL_BITS RF_SIGNAL_BYTES * 8

typedef struct {
  uint8_t duration_short_pulse;  //mulitple of 0.08ms, should be === 0 (mod 4)
  uint8_t short_mult;
  uint8_t long_mult;
  uint8_t sync_mult;
  uint8_t signal[RF_SIGNAL_BYTES];  //24bit signal info, excluding sync signal (short 1 followed by long pause (~128*0.08ms))
                            //for each bit: 0 means 1/4 Tau high followed by 3/4 Tau low;    1 means 3/4 Tau high followed by 1/4 Tau low
} rf_signal;

rf_signal current_signal = {6, 1, 3, 31, {0,0,0}};

typedef struct {
  uint8_t atime; // time counter
  uint8_t bit;  //index for current bit
  uint8_t repeatc; //downward couner of repetition
  uint8_t state; // current output to RF Pin (position within the bit)
} rf_state;

rf_state current_state = { 0, 0, 0, 0};
int rf_num_transmissions_to_acknowledge = 0;

#define CURRENT_BIT_CNT (RF_SIGNAL_BITS - current_state.bit - 1)
#define CURRENT_BIT (( current_signal.signal[ CURRENT_BIT_CNT/8] >> (CURRENT_BIT_CNT % 8)  )& 1)
#define RF_TIME_SHORT (current_signal.short_mult * current_signal.duration_short_pulse)
#define RF_TIME_LONG (current_signal.long_mult * current_signal.duration_short_pulse)
#define RF_TIME_SNYC (current_signal.sync_mult * current_signal.duration_short_pulse)
#define RF_OFF PIN_LOW(RF_DATA_OUT_PORT, RF_DATA_OUT_PIN)
#define RF_ON PIN_HIGH(RF_DATA_OUT_PORT, RF_DATA_OUT_PIN)

ISR(TIMER1_COMPA_vect)
{
  if ( current_state.state || current_state.bit || current_state.repeatc || current_state.atime)
  {
    if (current_state.atime)
    {
       current_state.atime--;
    }
    //atime ran out
    else if (current_state.state) //was in state 1 or 2
    {
      RF_OFF;  //stop sending
      if (current_state.state == 2) //aka sync
        current_state.atime=RF_TIME_SNYC;
      else
        current_state.atime=CURRENT_BIT?
           RF_TIME_SHORT
          :RF_TIME_LONG;
      current_state.state=0;
    }
    else if  (current_state.bit)  //still more than 0 bits to do
    {
      current_state.bit--;
      current_state.state=1;
      current_state.atime=CURRENT_BIT?
           RF_TIME_LONG
          :RF_TIME_SHORT;
      RF_ON;  //start sending
    }
    else if (current_state.repeatc)
    {
      current_state.bit=RF_SIGNAL_BITS;
      current_state.repeatc--;
      current_state.state=2;
      //start sync (short pulse followed by long pause)
      RF_ON;
      current_state.atime=RF_TIME_SHORT;
    }
  }
  else
  {
    rf433_stop_timer();
    RF_OFF;
    rf_num_transmissions_to_acknowledge++;
  }
}
//********************************************************************//

void rf433_send_rf_cmd(uint8_t sr[])
{
  while (TIMER_RUNNING)
  {}
  for (uint8_t chr=0; chr < 3; chr++)
  {
    current_signal.signal[chr]=sr[chr];
  }
  current_state.repeatc=NUM_REPEAT_SIGNAL;
  rf433_start_timer();
}

void rf433_check_frame_done(void)
{
  while (rf_num_transmissions_to_acknowledge > 0)
  {
    rf_num_transmissions_to_acknowledge--;
  }
}
