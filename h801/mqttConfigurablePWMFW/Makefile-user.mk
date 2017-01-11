## Local build configuration
## Parameters configured here will override default and ENV values.
## Uncomment and change examples:

## Add your source directories here separated by space
# MODULES = app
# EXTRA_INCDIR = include

## ESP_HOME sets the path where ESP tools and SDK are located.
## Windows:
# ESP_HOME = c:/Espressif

## MacOS / Linux:
ESP_HOME = /home/bernhard/source/esp-open-sdk

## SMING_HOME sets the path where Sming framework is located.
## Windows:
# SMING_HOME = c:/tools/sming/Sming 

## MacOS / Linux
SMING_HOME = /home/bernhard/source/Sming/Sming

## COM port parameter is reqruied to flash firmware correctly.
## Windows: 
# COM_PORT = COM3

## MacOS / Linux:
# COM_PORT = /dev/tty.usbserial

## Com port speed
# COM_SPEED	= 115200

## Configure flash parameters (for ESP12-E and other new boards):
# SPI_MODE = dio
#SPI_SIZE=512K
SPI_SIZE=1024K

#RBOOT_BIG_FLASH=0

## SPIFFS options
# DISABLE_SPIFFS = 1
SPIFF_FILES = files
#SPIFF_SIZE = 131072
#SPIFF_SIZE = 262144


ENABLE_CUSTOM_PWM=1
ENABLE_SSL=0
