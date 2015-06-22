##
##  lightctrl - teensy firmware code
##
##
##  Copyright (C) 2013 Bernhard Tittelbach <xro@realraum.at>
##   uses avr-utils, anyio & co by Christian Pointner <equinox@spreadspace.org>
##
##  This file is part of lightctrl.
##
##  lightctrl is free software: you can redistribute it and/or modify
##  it under the terms of the GNU General Public License as published by
##  the Free Software Foundation, either version 3 of the License, or
##  any later version.
##
##  lightctrl is distributed in the hope that it will be useful,
##  but WITHOUT ANY WARRANTY; without even the implied warranty of
##  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
##  GNU General Public License for more details.
##
##  You should have received a copy of the GNU General Public License
##  along with lightctrl. If not, see <http://www.gnu.org/licenses/>.
##

NAME := lightctrl
BOARD_TYPE := teensy2
OBJ := $(NAME).o rf433.o
LIBS := util led lufa-descriptor-usbdualserial dualusbio
EXTERNAL_LIBS := lufa
SPREADAVR_PATH := ../contrib/avr-utils

LUFA_PATH := $(SPREADAVR_PATH)/contrib/LUFA-120219
LUFA_OPTS  = -D USB_DEVICE_ONLY
LUFA_OPTS += -D DEVICE_STATE_AS_GPIOR=0
LUFA_OPTS += -D ORDERED_EP_CONFIG
LUFA_OPTS += -D FIXED_CONTROL_ENDPOINT_SIZE=8
LUFA_OPTS += -D FIXED_NUM_CONFIGURATIONS=1
LUFA_OPTS += -D USE_FLASH_DESCRIPTORS
LUFA_OPTS += -D USE_STATIC_OPTIONS="(USB_DEVICE_OPT_FULLSPEED | USB_OPT_REG_ENABLED | USB_OPT_AUTO_PLL)"
LUFA_OPTS += -D INTERRUPT_CONTROL_ENDPOINT

LUFA_OPTS += -D USB_MANUFACTURER="L\"xro\"" -D USB_MANUFACTURER_LEN=3
LUFA_OPTS += -D USB_PRODUCT="L\"lightctrl\"" -D USB_PRODUCT_LEN=9

LUFA_COMPONENTS := USB USBCLASS SERIAL

include $(SPREADAVR_PATH)/include.mk
