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
#include <avr/wdt.h>
#include <avr/interrupt.h>
#include <avr/power.h>
#include <stdio.h>

#include "util.h"
#include "led.h"
#include "usbio.h"

#include "relay.h"
#include "keypad.h"

uint8_t current_relay = 'a';

static void print_relay_state(void)
{
  for(uint8_t i = 0; i < RELAY_NUM; i++) {
    putchar(relay_get(i));
  }
  putchar('\r');
  putchar('\n');
}

static void handle_cmd(uint8_t cmd)
{
  switch(cmd) {
  case 'a':
  case 'b':
  case 'c':
  case 'd':
  case 'e':
  case 'f':
  case 'g':
  case 'h':
    current_relay = cmd; printf("output %c selected\r\n", current_relay); return;
  case '0':
  case '1':
  case 't':
  case '-':
    relay_set(current_relay - 'a', cmd); break;
  default: printf("????????\r\n"); return;
  }
  print_relay_state();
}

int main(void)
{
  MCUSR &= ~(1 << WDRF);
  wdt_disable();

  cpu_init();
  led_init();
  usbio_init();
  relay_init();
  keypad_init();
  sei();

  for(;;) {
    int16_t BytesReceived = usbio_bytes_received();
    while(BytesReceived > 0) {
      int ReceivedByte = fgetc(stdin);
      if(ReceivedByte != EOF) {
        handle_cmd(ReceivedByte);
      }
      BytesReceived--;
    }

    if(keypad_task()) {
      print_relay_state();
    }
    usbio_task();
  }
}
