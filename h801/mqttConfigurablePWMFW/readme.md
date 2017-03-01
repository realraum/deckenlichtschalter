H801 MQTT controllable light
============================

Blink Codes & Userfeedback
==========================

Blink GREEN Once: Wifi connected successfully

Blink RED Once: WiFi connection error

Blink RED Twice: MQTT connection error

Blink RED Thrice: MQTT broker unreachable

Blink Blue Once: New default power-on setting accepted


MQTT Usage
==========

all mqtt strings are configurable in ```pwmchannels.h``` or with the config files but here the defaults are used.

e.g. for a H801 with clientid set to ```ceiling1``` which will also respond to ```ceilingAll```:

Topic ```action/ceilingAll/defaultlight```
----------------------------------------

Set default light values, applied on power-up.

Unset values will be taken from current light level. E.g. sending an empty object to ```+/defaultlight``` will save
current light values as default.

    {
	"description": "Schema for message to topic action/+/defaultlight",
	"type": "object",
	"required": [],
	"properties": {
		"r": {
			"description": "set red LED light default intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"g": {
			"description": "set green LED light default intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"b": {
			"description": "set blue LED light default intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"cw": {
			"description": "set cool white LED light default intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"ww": {
			"description": "set warm white LED light default intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
	}

### Example: Saving current light values as default

    mosquitto_pub -t action/ceiling2/defaultlight -m '{}'

### Example: Have lights off after power up

    mosquitto_pub -t action/ceiling2/defaultlight -m '{"r":0,"g":0,"b":0,"cw":0,"ww":0}'


Topic ```action/ceilingAll/light```
---------------------------------

Change current light values. Optionally use an effect and chain effects.

Unset values will remain unchanged.


    {
	"description": "Schema for message to topic action/+/light",
	"type": "object",
	"required": [],
	"properties": {
		"r": {
			"description": "change red LED light intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"g": {
			"description": "change green LED light intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"b": {
			"description": "change blue LED light intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"cw": {
			"description": "change cool white LED light intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"ww": {
			"description": "change warm white LED light intensity",
			"type": "integer",
			"minimum": 0,
			"maximum": 1000,
		},
		"fade": {
			"title": "Fade Effect",
			"description": "If present, changes will fade in over time. If object flash is present, fade will be ignored.",
			"type": "object",
			"required": [],
			"properties": {
				"duration":	{
					"description": "effect duration in ms",
					"type": "integer",
					"minimum": 100,
					"maximum": 60000,
				},
				"cc": {
					"title": "carbon copy",
					"description": "After effect has finished or been aborted, forward this MQTT message to first topic in this array. The array is modified and the first element removed on each forward step. Can be used to chain lights or report that effect has finished or stupidly even to increase repetitions.",
					"type": "array",
					"items": {
            			"type": "string"
        			},
        			"minItems": 0,
        			"uniqueItems": false,
				},
			},
		},
		"flash": {
			"title": "Flash Effect"
			"description": "If present, changes will be applied temporarily for 800ms at a time. If present, fade object will be ignored."
			"type": "object"
			"required": [],
			"properties": {
				"repetitions":	{
					"description": "number of repetitions of flash effect",
					"type": "integer",
					"minimum": 1,
					"maximum": 10,
				},
				"cc": {
					"title": "carbon copy",
					"description": "After effect has finished or been aborted, forward this MQTT message to first topic in this array. The array is modified and the first element removed on each forward step. Can be used to chain lights or report that effect has finished or stupidly even to increase repetitions.",
					"type": "array",
					"items": {
            			"type": "string"
        			},
        			"minItems": 0,
        			"uniqueItems": false,
				},					
			},
		},			
	}


### Example: Change WarmWhite to 500%%

    mosquitto_pub -t action/ceiling1/light -m '{"ww":500}'

### Example: Fade to Black

    mosquitto_pub -t action/ceiling1/light -m '{"r":0,"g":0,"b":0,"cw":0,"ww":0,"fade":{"duration":3200}}'

### Example: Fade one light after another to intense blue

    mosquitto_pub -t action/ceiling1/light -m '{"r":0,"g":0,"b":1000,"cw":0,"ww":0,"fade":{"duration":1500, "cc":["action/ceiling2/light","action/ceiling3/light","action/ceiling4/light"]}}'

### Example: Flash Red 3 times

    mosquitto_pub -t action/ceiling1/light -m '{"r":1000,"g":0,"b":0,"cw":0,"ww":0,"flash":{"repetitions":3}}'

### Example: Flash Green 5 times, then report back

    mosquitto_pub -t action/ceiling1/light -m '{"r":0,"g":1000,"b":0,"cw":0,"ww":0,"flash":{"repetitions":5, "cc":["action/ceilingcoordinator/ceiling1/flashdone"]}}'


Topic ```action/ceilingAll/pleaserepeat```
----------------------------------------

Useful if you want to discover current light values. Subscribe to the same topic as the light listens on, e.g. ```action/ceiling1/light``` and trigger all lights to send a change command to themselves by sending an empty object to ```action/ceilingAll/pleaserepeat```


    {
	"description": "Triggers sending MQTT msg with light values to itself",
	"type": "object",
	"required": [],
	"properties": {}
	}



### Example: Tell me your current light settings

    mosquitto_sub -t action/ceiling1/light &
    mosquitto_pub -t action/ceiling1/pleaserepeat -m '{}'


Advanced Usage Examples
=======================

see ```testscripts/fadeanime.py``` for an example script that one after another changes each light to a new random colour. Choosing a new colour after all lights have changed.


Telnet Interface
================

Telnet-Port is 2323

Changing Settings
-----------------

* ```load``` command  
  reads SPIFF files and provides their contents via NetConfig structs. Largely useless except to undo mistakes with ```set```.
* ```cat``` and ```ls``` commands  
  provide a peek into what would be loaded
* ```show``` command  
  lists the current settings and some status. Actual stuff are currently applied settings, the rest is the content of the NetConfig struct.
* ```set``` command  
  can be used to change the settings in the NetConfig structs.
* ```save``` command  
  writes NetConfig struct data to SPIFF, making changes permanent. Call ```restart``` to use new settings.
* ```restart``` command  
  restart's the H801. Afterwards new settings from SPIFF will be used.
* ```auth``` command  
  prevents accidential mistakes
* ```update``` command  
  can update the firmware over-the-air. See next section.


Over-The-Air Update Notes
=========================

Compile using Sming

Update Procedure
----------------

1. configure SPIFFS using ```ceilingspiffsconfig.py```
2. ```make clean; make```
3. ```cd out/firmware```
4. start Webserver: ```python -m SimpleHTTPServer 8080```
5. connect via telnet to H801, e.g. ```telnet ceiling1.lan 2323```
6. provide configured auth string: ```auth prevents mistakes <...>```
7. start OTA update: e.g. ```update http://mypc.lan/```
8. terminate telnet session
9. wait
10. power-cycle H801


Developement Notes
------------------

 * H801 has 1MiB flash 
     * ```SPI_SIZE=1M```
 * We divide the 1MiB into two 512KiB slots and boot one of them using rBoot bootloader
     * ``RBOOT_TWO_ROMS=1```
 * On OTA-Update, the currently unused slot is written, then the slot ist switched and bootet into after restart.
 * Each slot's SPIFFS is set to size 64KiB and put at end of the 512KiB.
     * ``` SPIFF_SIZE=65536```
     * ``` RBOOT_SPIFFS_0   ?= $(shell echo $$((  0x80000 - $(SPIFF_SIZE) )) ) ```
     * ``` RBOOT_SPIFFS_1   ?= $(shell echo $$(( 0x100000 - $(SPIFF_SIZE) )) ) ```



