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

#include "usbio.h"

#include "keypad.h"
#include "relay.h"

static uint8_t fancy_buf[16];

void fancy_init(void)
{
  usbio_init();
  memset(fancy_buf, 0, sizeof(fancy_buf));
}

static int16_t fancy_read(int16_t bytes_received, uint8_t* changed)
{
  *changed = 0;
  return bytes_received;
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
