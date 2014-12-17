##
##  spreadspace avr utils
##
##
##  Copyright (C) 2013 Christian Pointner <equinox@spreadspace.org>
##
##  This file is part of spreadspace avr utils.
##
##  spreadspace avr utils is free software: you can redistribute it and/or modify
##  it under the terms of the GNU General Public License as published by
##  the Free Software Foundation, either version 3 of the License, or
##  any later version.
##
##  spreadspace avr utils is distributed in the hope that it will be useful,
##  but WITHOUT ANY WARRANTY; without even the implied warranty of
##  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
##  GNU General Public License for more details.
##
##  You should have received a copy of the GNU General Public License
##  along with spreadspace avr utils. If not, see <http://www.gnu.org/licenses/>.
##

NAME := lightctrl
BOARD_TYPE := teensy2
OBJ := $(NAME).o pwm.o pid_control.o temp_curve.o cmd_queue.o
LIBS := util led lufa-descriptor-usbdualserial anyio
EXTERNAL_LIBS := lufa
RESET_FUNC := ../tools/reset_lufa_cdc
RESET_PARAM := 'R'

LUFA_PATH := ../contrib/LUFA-120219
LUFA_OPTS  = -D USB_DEVICE_ONLY
LUFA_OPTS += -D DEVICE_STATE_AS_GPIOR=0
LUFA_OPTS += -D ORDERED_EP_CONFIG
LUFA_OPTS += -D FIXED_CONTROL_ENDPOINT_SIZE=8
LUFA_OPTS += -D FIXED_NUM_CONFIGURATIONS=1
LUFA_OPTS += -D USE_FLASH_DESCRIPTORS
LUFA_OPTS += -D USE_STATIC_OPTIONS="(USB_DEVICE_OPT_FULLSPEED | USB_OPT_REG_ENABLED | USB_OPT_AUTO_PLL)"
LUFA_OPTS += -D INTERRUPT_CONTROL_ENDPOINT

LUFA_OPTS += -D USB_MANUFACTURER="L\"xro\""
LUFA_OPTS += -D USB_PRODUCT="L\"$(NAME)\""

LUFA_COMPONENTS := USB USBCLASS SERIAL

include ../include.mk
