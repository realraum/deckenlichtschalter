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
#ESP_HOME = 

## SMING_HOME sets the path where Sming framework is located.
## Windows:
# SMING_HOME = c:/tools/sming/Sming 

## MacOS / Linux
#SMING_HOME = 

## COM port parameter is reqruied to flash firmware correctly.
## Windows: 
# COM_PORT = COM3

## MacOS / Linux:
# COM_PORT = /dev/tty.usbserial

## Com port speed
# COM_SPEED	= 115200

## Configure flash parameters (for ESP12-E and other new boards):
# SPI_MODE = dio

## SPIFFS options
# DISABLE_SPIFFS = 1
SPIFF_FILES = files

#### overridable rBoot options ####
## use rboot build mode
RBOOT_ENABLED ?= 1
## enable big flash support (for multiple roms, each in separate 1mb block of flash)
RBOOT_BIG_FLASH ?= 1
## two rom mode (where two roms sit in the same 1mb block of flash)
#RBOOT_TWO_ROMS  ?= 1
## size of the flash chip
#SPI_SIZE        ?= 4M
## output file for first rom (.bin will be appended)
#RBOOT_ROM_0     ?= rom0
## input linker file for first rom
#RBOOT_LD_0      ?= rom0.ld
## these next options only needed when using two rom mode
#RBOOT_ROM_1     ?= rom1
#RBOOT_LD_1      ?= rom1.ld
## size of the spiffs to create
#SPIFF_SIZE      ?= 65536
## option to completely disable spiffs
#DISABLE_SPIFFS  = 1
## flash offsets for spiffs, set if using two rom mode or not on a 4mb flash
## (spiffs location defaults to the mb after the rom slot on 4mb flash)
#RBOOT_SPIFFS_0  ?= 0x100000
#RBOOT_SPIFFS_1  ?= 0x300000
## esptool2 path
#ESPTOOL2        ?= esptool2



SPI_SIZE=1M
RBOOT_BIG_FLASH=0
RBOOT_TWO_ROMS=1
SPIFF_SIZE=0x15000
SHELL=/bin/zsh
## H801 flash is 1MiB big
## flash ends at    0x100000
## half of flash is 0x080000
## our rboot .ld scripts write rom0 to 0x2010 and rom1 to 0x82000
## expressiv writes data to last 16KiB of Flash at 0x0fc000 in our case
## see https://github.com/raburton/rboot for explanation
## max size for code and spiffs together is min(0x80000-0x2000-0x10,0x80000-0x2000-0x4000)=0x7a000
## thus if SPIFF is 0x15000 then max code len is 0x80000-0x2000-0x4000-0x15000=0x65000
## thus we set len=0x65000 in .lds
## thus our flash layout will look like this on 1MB H801 Flash
## 
## 0x000000 rboot
## 0x001000 
## 0x002000 rboot settings 10 byte
## 0x002010 rom0  <-- this is where rboot expects to find the image at virtual 0x40202010
## 0x06b000 spiffs0  <-- this is where our app will load spiffs from
## 0x080000
## 0x082000 rom1  <-- this is where rboot expects to find the image at virtual 0x40282010
## 0x0e7000 spiffs1  <-- this is where our app will load spiffs from
## 0x0fc000 expressif ssid settings, etc
## 0x100000 end
RBOOT_SPIFFS_0   ?= $(shell echo $$((  0x80000 - $(SPIFF_SIZE) )) )
RBOOT_SPIFFS_1   ?= $(shell echo $$((  0xfc000 - $(SPIFF_SIZE) )) )

ENABLE_CUSTOM_PWM=1
ENABLE_SSL=0

## Compile Options
#USER_CFLAGS += -DENABLE_BUTTON=1
#USER_CFLAGS += -DREPLACE_CW_WITH_UV=1
USER_CFLAGS += -DTELNET_CMD_LIGHTTEST=1
