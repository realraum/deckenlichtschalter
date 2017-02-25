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

#include <stdio.h>
#include <avr/wdt.h>
#include <avr/interrupt.h>

#include "util.h"

#include "relay.h"
#include "keypad.h"
#include "fancy.h"

static void print_relay_state(void)
{
  for(uint8_t i = 0; i < RELAY_NUM; i++) {
    putchar(relay_get(i));
  }
  putchar('\r');
  putchar('\n');
}

int main(void)
{
  MCUSR &= ~(1 << WDRF);
  wdt_disable();
  jtag_disable();

  cpu_init();
  relay_init();
  keypad_init();
  fancy_init();
  sei();

  for(;;) {
    if(keypad_task()) {
      print_relay_state();
    }
    if(fancy_task()) {
      print_relay_state();
    }
  }
}
