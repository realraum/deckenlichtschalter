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

#include "keypad.h"

void keypad_init(void)
{
      //         #6         #7
  DDRB &= ~( (1<<PB4) | (1<<PB5) );
      //       #3
  DDRC &= ~(1<<PC6);
      //         #1         #0         #2         #4
  DDRD &= ~( (1<<PD0) | (1<<PD1) | (1<<PD4) | (1<<PD7) );
      //       #5
  DDRE &= ~(1<<PE6);
}
