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

#include <string.h>
#include <stdio.h>

#include "usbio.h"

#include "keypad.h"
#include "relay.h"

static uint8_t fancy_buf_idx;
static uint8_t fancy_buf[RELAY_NUM+2];

void fancy_init(void)
{
  usbio_init();
  fancy_buf_idx = 0;
  memset(fancy_buf, 0, sizeof(fancy_buf));
}

static int16_t fancy_read(int16_t bytes_received, uint8_t* changed)
{
  uint8_t bytes_consumed = 0;

  for(;;) {
    fancy_buf[fancy_buf_idx] = (uint8_t)getchar();
    bytes_consumed++;
    if(fancy_buf[fancy_buf_idx] == '\n' || fancy_buf[fancy_buf_idx] == '\r') {
      if(fancy_buf_idx == RELAY_NUM) {
        for(uint8_t i = 0; i < RELAY_NUM; ++i) {
          relay_set(i, fancy_buf[i]);
        }
        *changed = 1;
      }
      fancy_buf_idx = 0;
    } else {
      if((++fancy_buf_idx) >= sizeof(fancy_buf)) {
        fancy_buf_idx = sizeof(fancy_buf) - 1;
      }
    }
    if(bytes_consumed == bytes_received) {
      break;
    }
  }
  return bytes_consumed;
}

uint8_t fancy_task(void)
{
  usbio_task();

  uint8_t changed = 0;
  int16_t bytes_received = usbio_bytes_received();
  while(bytes_received > 0) {
    bytes_received -= fancy_read(bytes_received, &changed);
  }
  return changed;
}
