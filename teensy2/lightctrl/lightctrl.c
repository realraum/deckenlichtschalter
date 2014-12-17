/*
 *  r3PCR Teensy Controller Code
 *
 *
 *  Copyright (C) 2013 Bernhard Tittelbach <xro@realraum.at>
*   uses avr-utils, anyio & co by Christian Pointner <equinox@spreadspace.org>
 *
 *  This file is part of spreadspace avr utils.
 *
 *  spreadspace avr utils is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  any later version.
 *
 *  spreadspace avr utils is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with spreadspace avr utils. If not, see <http://www.gnu.org/licenses/>.
 */


#include <avr/wdt.h>
#include <avr/interrupt.h>
#include <avr/power.h>
#include <stdio.h>

#include "util.h"
#include "led.h"
#include "anyio.h"
#include "mypins.h"

uint8_t relais_state_ = 0;
uint16_t buttons_pressed_ = 0;
// at f_system_clk = 10Hz, system_clk_ will not overrun for at least 13 years. PCR won't run that long
volatile uint32_t system_clk_ = 0;

//with F_CPU = 16MHz and TIMER3 Prescaler set to /1024, TIMER3 increments with f = 16KHz. Thus if TIMER3 reaches 16, 1ms has passed.
#define T3_MS     *16
//set TICK_TIME to 1/10 of a second
#define SYSCLKTICK_DURATION_IN_MS 100
#define	TICK_TIME (SYSCLKTICK_DURATION_IN_MS T3_MS)

ISR(TIMER3_COMPA_vect)
{
  //increment system_clk every TIME_TICK (aka 100ms)
	system_clk_++;
  //set up "clock" comparator for next tick
  OCR3A = (OCR3A + TICK_TIME) & 0xFFFF;
  if (debug_)
    led_toggle();
}

void initSysClkTimer3(void)
{
  system_clk_ = 0;
  // set counter to 0
  TCNT3 = 0x0000;
	// no outputs
	TCCR3A = 0;
	// Prescaler for Timer3: F_CPU / 1024 -> counts with f= 16KHz ms
	TCCR3B = _BV(CS32) | _BV(CS30);
	// set up "clock" comparator for first tick
	OCR3A = TICK_TIME & 0xFFFF;
	// enable interrupt
	TIMSK3 = _BV(OCIE3A);
}

void printStatus(void)
{
  printf("%c%c%c\n",relais_state_, buttons_pressed_>>8, buttons_pressed_&0xff);
}

void readButtons(uint16_t *buttons)
{
  *buttons = 0;
  for (uint8_t c=0; c<6; c++)
  {
    if (TEST_BTN_ON(c))
    {
      *buttons |= (1<<(c*2));
    } else if (TEST_BTN_OFF(c))
    {
      *buttons |= (1<<(c*2+1));
    }
  }
  for (uint8_t c=0; c<3; c++)
  {
    if (TEST_BTN_SIG(c))
    {
      *buttons |= (1<<c+12);
    }
  }
}

void buttonsToNewState()
{
  if (buttons_pressed_ & (1<<12)) {
    relais_state_ = 0;
  } else if (buttons_pressed_ & (1<<13)) {
    relais_state_ = 0x0C;
  } else if (buttons_pressed_ & (1<<14)) {
    relais_state_ = RELAIS_MASK;
  }
  for (uint8_t c=0; c<6; c++)
  {
    if (buttons_pressed_ & (1<<(c*2)))
    {
      //on
      relais_state_ |= (1 << c);
    } else if (buttons_pressed_ & (1<<(c*2+1)))
    {
      //off
      relais_state_ &= ~(1 << c);
    }
  }
}

void applyRelaisState(uint8_t new_state)
{
  relais_state_ = new_state & RELAIS_MASK;
  RELAIS_PORT = (RELAIS_PORT & ~RELAIS_MASK) | (relais_state_ ^ RELAIS_INVERT_MASK);
}

//reads exactly bufflen-1 chars
//does not start to fill buffer until start_escape char is seen
//further occurances of start_escape have to be escaped by an additional start_escape char
void readFixedLenSeqIntoBufferWStartEscapeSymbol(char *buffer, uint8_t buflen, char start_escape, uint8_t start_with_startescape_seen)
{
  while (anyio_bytes_received() == 0);
  int ReceivedByte=0;
  uint8_t index=0;
  uint8_t startescape_seen = start_with_startescape_seen & 0x01;
  do {
    ReceivedByte = fgetc(stdin);
    if (ReceivedByte != EOF)
    {
      if ((char) ReceivedByte == start_escape)
      {
        if (startescape_seen == 0)
        {
          startescape_seen=1;
          continue;
        }
      } else {
        if (startescape_seen)
        {
          index=0;
        }
      }
      startescape_seen=0;
      buffer[index++] = (char) ReceivedByte;
    }
  } while (index < buflen-1);
  buffer[index] = 0;
}

int main(void)
{
  /* Disable watchdog if enabled by bootloader/fuses */
  MCUSR &= ~(1 << WDRF);
  wdt_disable();

  cpu_init();
  led_init();
  anyio_init(115200, 0);
  sei();

  led_off();

  //set default all to INPUT
  DDRB = 0;
  DDRC = 0;
  DDRD = 0;
  DDRF = 0;
  //set pins for Relais to OUTPUT
  RELAIS_DDR |= RELAIS_MASK;
  //set pin for RF to OUTPUT
  PINMODE_OUTPUT(RF_DATA_OUT_DDR, RF_DATA_OUT_PIN);

  //set PULL-UP on button pins
  //PORTB = ~DDRB & PULLUP_DDRB;
  //PORTC = ~DDRC & PULLUP_DDRC;
  //PORTD = ~DDRD & PULLUP_DDRD;
  //PORTF = ~DDRF & PULLUP_DDRF;
  //set PULL-UP on everything that is not an OUTPUT to have well defined levels on unconnected and button pins
  PORTB = ~DDRB;
  PORTC = ~DDRC;
  PORTD = ~DDRD;
  PORTF = ~DDRF;

//  pwm_init();
//  pwm_b5_set(0);

  initSysClkTimer3(); //start system clock

  applyRelaisState(0);

  for(;;)
  {
    int16_t BytesReceived = anyio_bytes_received();
    while(BytesReceived > 0)
    {
      int ReceivedByte = fgetc(stdin);
      if (ReceivedByte != EOF)
      {
        applyRelaisState((uint8_t) ReceivedByte);
        printStatus();
      }
      BytesReceived--;
    }

    //if got rf sequence
    //rf433_start_timer();

    readButtons(&buttons_pressed_);

    if (buttons_pressed_) {
      buttonsToNewState();
      printStatus();
      buttons_pressed_ = 0;
    } else {
    }

    anyio_task();
  }
}
