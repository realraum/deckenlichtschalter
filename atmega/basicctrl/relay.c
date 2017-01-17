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
      //       #4         #6         #5         #7
  DDRB  |= (1<<PB1) | (1<<PB2) | (1<<PB3) | (1<<PB6);
  PORTB |= (1<<PB1) | (1<<PB2) | (1<<PB3) | (1<<PB6);
      //       #0         #1         #2         #3
  DDRF  |= (1<<PF4) | (1<<PF5) | (1<<PF6) | (1<<PF7);
  PORTF |= (1<<PF4) | (1<<PF5) | (1<<PF6) | (1<<PF7);
}

void relay_on(uint8_t num)
{
  switch(num) {
  case 0: PORTF &= ~(1<<PF4); return;
  case 1: PORTF &= ~(1<<PF5); return;
  case 2: PORTF &= ~(1<<PF6); return;
  case 3: PORTF &= ~(1<<PF7); return;
  case 4: PORTB &= ~(1<<PB1); return;
  case 5: PORTB &= ~(1<<PB3); return;
  case 6: PORTB &= ~(1<<PB2); return;
  case 7: PORTB &= ~(1<<PB6); return;
  }
}

void relay_off(uint8_t num)
{
  switch(num) {
  case 0: PORTF |= 1<<PF4; return;
  case 1: PORTF |= 1<<PF5; return;
  case 2: PORTF |= 1<<PF6; return;
  case 3: PORTF |= 1<<PF7; return;
  case 4: PORTB |= 1<<PB1; return;
  case 5: PORTB |= 1<<PB3; return;
  case 6: PORTB |= 1<<PB2; return;
  case 7: PORTB |= 1<<PB6; return;
  }
}

void relay_toggle(uint8_t num)
{
  switch(num) {
  case 0: PORTF ^= 1<<PF4; return;
  case 1: PORTF ^= 1<<PF5; return;
  case 2: PORTF ^= 1<<PF6; return;
  case 3: PORTF ^= 1<<PF7; return;
  case 4: PORTB ^= 1<<PB1; return;
  case 5: PORTB ^= 1<<PB3; return;
  case 6: PORTB ^= 1<<PB2; return;
  case 7: PORTB ^= 1<<PB6; return;
  }
}

uint8_t relay_get(uint8_t num)
{
  switch(num) {
  case 0: return (PINF & 1<<PF4) ? 0 : 1;
  case 1: return (PINF & 1<<PF5) ? 0 : 1;
  case 2: return (PINF & 1<<PF6) ? 0 : 1;
  case 3: return (PINF & 1<<PF7) ? 0 : 1;
  case 4: return (PINB & 1<<PB1) ? 0 : 1;
  case 5: return (PINB & 1<<PB3) ? 0 : 1;
  case 6: return (PINB & 1<<PB2) ? 0 : 1;
  case 7: return (PINB & 1<<PB6) ? 0 : 1;
  }
  return 0;
}
