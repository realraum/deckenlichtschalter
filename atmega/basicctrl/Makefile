##
##  basicctrl
##
##
##  Copyright (C) 2017 Christian Pointner <equinox@realraum.at>
##
##  This file is part of basicctrl.
##
##  basicctrl is free software: you can redistribute it and/or modify
##  it under the terms of the GNU General Public License as published by
##  the Free Software Foundation, either version 3 of the License, or
##  any later version.
##
##  basicctrl is distributed in the hope that it will be useful,
##  but WITHOUT ANY WARRANTY; without even the implied warranty of
##  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
##  GNU General Public License for more details.
##
##  You should have received a copy of the GNU General Public License
##  along with basicctrl. If not, see <http://www.gnu.org/licenses/>.
##

NAME := basicctrl
BOARD_TYPE := arduinoProMicro
OBJ := $(NAME).o relay.o keypad.o fancy.o
LIBS := util lufa-descriptor-usbserial usbio
EXTERNAL_LIBS := lufa
SPREADAVR_PATH := ../contrib/avr-utils

LUFA_PATH := $(SPREADAVR_PATH)/contrib/lufa-LUFA-151115
LUFA_OPTS  = -D USB_DEVICE_ONLY
LUFA_OPTS += -D DEVICE_STATE_AS_GPIOR=0
LUFA_OPTS += -D ORDERED_EP_CONFIG
LUFA_OPTS += -D FIXED_CONTROL_ENDPOINT_SIZE=8
LUFA_OPTS += -D FIXED_NUM_CONFIGURATIONS=1
LUFA_OPTS += -D USE_FLASH_DESCRIPTORS
LUFA_OPTS += -D USE_STATIC_OPTIONS="(USB_DEVICE_OPT_FULLSPEED | USB_OPT_REG_ENABLED | USB_OPT_AUTO_PLL)"
LUFA_OPTS += -D INTERRUPT_CONTROL_ENDPOINT

LUFA_OPTS += -D USB_MANUFACTURER="L\"realraum\""
LUFA_OPTS += -D USB_PRODUCT="L\"basicctrl\""

LUFA_COMPONENTS := USB USBCLASS

include $(SPREADAVR_PATH)/include.mk
