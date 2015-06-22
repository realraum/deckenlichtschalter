/*
 *  lightctrl - teensy firmware code
 *
 *
 *  Copyright (C) 2013 Bernhard Tittelbach <xro@realraum.at>
 *   uses avr-utils, anyio & co by Christian Pointner <equinox@spreadspace.org>
 *
 *  This file is part of lightctrl.
 *
 *  lightctrl is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  any later version.
 *
 *  lightctrl is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with lightctrl. If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef LIGHTCTRL_rf433_h_INCLUDED
#define LIGHTCTRL_rf433_h_INCLUDED

void rf433_send_rf_cmd(uint8_t sr[]);
void rf433_check_frame_done(void);

#endif
