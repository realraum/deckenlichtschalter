/*
 *  basicctrl
 *
 *
 *  Copyright (C) 2017 Christian Pointner <equinox@realraum.at>
 *
 *  This file is part of basicctrl.
 *
 *  basicctrl is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  any later version.
 *
 *  basicctrl is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with basicctrl. If not, see <http://www.gnu.org/licenses/>.
 */

#include <avr/io.h>

#include "relay.h"

void relay_init(void)
{
      //       #0         #1         #2         #3
  DDRB  |= (1<<PB0) | (1<<PB1) | (1<<PB2) | (1<<PB3);
  PORTB &= ~( (1<<PB0) | (1<<PB1) | (1<<PB2) | (1<<PB3) );
      //       #6         #7
  DDRC  |= (1<<PC6) | (1<<PC7);
  PORTC &= ~( (1<<PC6) | (1<<PC7) );
      //       #4         #5
  DDRD  |= (1<<PD0) | (1<<PD1);
  PORTD &= ~( (1<<PD0) | (1<<PD1) );
}

void relay_off(uint8_t num)
{
  switch(num) {
  case 0: PORTB &= ~(1<<PB0); return;
  case 1: PORTB &= ~(1<<PB1); return;
  case 2: PORTB &= ~(1<<PB2); return;
  case 3: PORTB &= ~(1<<PB3); return;
  case 4: PORTD &= ~(1<<PD0); return;
  case 5: PORTD &= ~(1<<PD1); return;
  case 6: PORTC &= ~(1<<PC6); return;
  case 7: PORTC &= ~(1<<PC7); return;
  }
}

void relay_on(uint8_t num)
{
  switch(num) {
  case 0: PORTB |= 1<<PB0; return;
  case 1: PORTB |= 1<<PB1; return;
  case 2: PORTB |= 1<<PB2; return;
  case 3: PORTB |= 1<<PB3; return;
  case 4: PORTD |= 1<<PD0; return;
  case 5: PORTD |= 1<<PD1; return;
  case 6: PORTC |= 1<<PC6; return;
  case 7: PORTC |= 1<<PC7; return;
  }
}

void relay_toggle(uint8_t num)
{
  switch(num) {
  case 0: PORTB ^= 1<<PB0; return;
  case 1: PORTB ^= 1<<PB1; return;
  case 2: PORTB ^= 1<<PB2; return;
  case 3: PORTB ^= 1<<PB3; return;
  case 4: PORTD ^= 1<<PD0; return;
  case 5: PORTD ^= 1<<PD1; return;
  case 6: PORTC ^= 1<<PC6; return;
  case 7: PORTC ^= 1<<PC7; return;
  }
}

uint8_t relay_get(uint8_t num)
{
  switch(num) {
  case 0: return (PORTB & 1<<PB0) ? RELAY_ON : RELAY_OFF;
  case 1: return (PORTB & 1<<PB1) ? RELAY_ON : RELAY_OFF;
  case 2: return (PORTB & 1<<PB2) ? RELAY_ON : RELAY_OFF;
  case 3: return (PORTB & 1<<PB3) ? RELAY_ON : RELAY_OFF;
  case 4: return (PORTD & 1<<PD0) ? RELAY_ON : RELAY_OFF;
  case 5: return (PORTD & 1<<PD1) ? RELAY_ON : RELAY_OFF;
  case 6: return (PORTC & 1<<PC6) ? RELAY_ON : RELAY_OFF;
  case 7: return (PORTC & 1<<PC7) ? RELAY_ON : RELAY_OFF;
  }
  return 0;
}

void relay_set(uint8_t num, uint8_t action)
{
  switch(action) {
  case RELAY_OFF: relay_off(num); return;
  case RELAY_ON: relay_on(num); return;
  case RELAY_TOGGLE: relay_toggle(num); return;
  }
}
