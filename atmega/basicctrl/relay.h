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

#ifndef BASICCTRL_relay_h_INCLUDED
#define BASICCTRL_relay_h_INCLUDED

#define RELAY_NUM 8

#define RELAY_OFF      '0'
#define RELAY_ON       '1'
#define RELAY_TOGGLE   't'
#define RELAY_NOCHANGE '-'

void relay_init(void);

void relay_on(uint8_t num);
void relay_off(uint8_t num);
void relay_toggle(uint8_t num);

uint8_t relay_get(uint8_t num);
void relay_set(uint8_t num, uint8_t action);

#endif
