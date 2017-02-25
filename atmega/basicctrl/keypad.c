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
#include <string.h>

#include "keypad.h"
#include "relay.h"

#define KEYPAD_LP_MAX 800

uint8_t keypad_state[KEYPAD_NUM];
uint16_t keypad_cnt[KEYPAD_NUM];

static uint8_t keypad_get_raw(uint8_t num)
{
  switch(num) {
  case 0: return (PINF & 1<<PF0) ? 1 : 0;
  case 1: return (PINF & 1<<PF1) ? 1 : 0;
  case 2: return (PINF & 1<<PF4) ? 1 : 0;
  case 3: return (PINF & 1<<PF5) ? 1 : 0;
  case 4: return (PINF & 1<<PF6) ? 1 : 0;
  case 5: return (PINF & 1<<PF7) ? 1 : 0;
  case 6: return (PINB & 1<<PB5) ? 1 : 0;
  case 7: return (PINB & 1<<PB4) ? 1 : 0;
  }
  return 0;
}

void keypad_init(void)
{
      //         #0         #1         #2         #3         #4         #5
  DDRF &= ~( (1<<PF0) | (1<<PF1) | (1<<PF4) | (1<<PF5) | (1<<PF6) | (1<<PF7) );
  PORTF |= (1<<PF0) | (1<<PF1) | (1<<PF4) | (1<<PF5) | (1<<PF6) | (1<<PF7);
      //         #7         #6
  DDRB &= ~( (1<<PB4) | (1<<PB5) );
  PORTB |= (1<<PB4) | (1<<PB5);

  memset(keypad_cnt, 0, sizeof(keypad_cnt));
  for(uint8_t i = 0; i < KEYPAD_NUM; i++) {
    keypad_state[i] = keypad_get_raw(i);
  }
}

uint8_t keypad_task(void)
{
  uint8_t changed = 0;
  for(uint8_t i = 0; i < KEYPAD_NUM; i++) {
    uint8_t state = keypad_get_raw(i);
    if(state != keypad_state[i])
      keypad_cnt[i]++;
    else
      keypad_cnt[i] += keypad_cnt[i] ? -1 : 0;

    if(keypad_cnt[i] >= KEYPAD_LP_MAX) {
      if(!state) {
        relay_toggle(i);
        changed = 1;
      }
      keypad_state[i] = state;
      keypad_cnt[i] = 0;
    }
  }
  return changed;
}
