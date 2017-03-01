ScriptCtrl
==========

* Runs and manages Scripts that animate the ceiling lights
* controlled via mqtt

Usage
-----

* Send a MQTT message to topic ```action/ceilingscripts/activatescript```
* let the payload be in JSON format and contain an element ```script```
* which is a string-value contaning the name of the script to activate.
* or an empty or invalid string-value to deactivate the running script.


Scripts currently in ScriptCtrl
===============================

ceilingsinus
------------

Circle a sinusodial colorchange around the room.

Can be used to create warm and uniform lighting that changes almost inperceptable (the default)
or to creat a noticable effect.

### Arguments

send to topic ```action/ceilingscripts/activatescript```

```
{
"description": "Arguments schema to run ceilingsinus",
"type": "object",
"required": ["script"],
"properties": {
    "script":{
        "description": "name of script to run",
        "type": "string",
        "value": "ceilingsinus"    
    },
    "r": {
        "description": "sinus parameters for red",
        "type": "object",
        "properties": {
          "amplitude":{
            "description": "amplitude of sinus curve"
            "minimum": 0,
            "maximum": 1000,
          }
          "offset":{
            "description": "zero offset of sinus curve. i.e. light level on sin(0)"
            "minimum": 0,
            "maximum": 1000,
          }
          "phase":{
            "description": "TODO: phase shift to other colors"
            "minimum": 0,
            "maximum": NumCeilingLights,
          }
        },
    },
    .... repeat for "g", "b", "cw", "ww" ...
}
```

### Algorithm

1. For each LED strip (r,g,b,ww,cw) calculate NumberOfLights Values from term ```value=offset + amplitude*sin(2*pi/NumLights*(LightIndex+extraoffset)%NumLights)```
2. Tell each ceilinglight to fade to


colorfade
---------

Makes the ceiling lights fade between different colors.

### Arguments

* ```value```: Illumination or rather HSV-V-Value of ceiling light colors in range 0.0 to 1.0.
  Set to 0.2 for low-key lightning or 1.0 for intensive colors.

send to topic ```action/ceilingscripts/activatescript```

```
{
"description": "Arguments schema to run ceilingsinus",
"type": "object",
"required": ["script"],
"properties": {
    "script":{
        "description": "name of script to run",
        "type": "string",
        "value": "colorfade"    
    },
    "value": {
        "description": "illumination level",
        "type": "float",
        "minimum": 0.0,
        "maximum": 1.0,
    },
}
```

### Algorithm

1. choose two sets of random lights and put them in a random order list
2. choose two random colorhues
3. fade first light in each list to fade to choosen color
4. fade next light in each list to same color and repeat until lists are empty
5. repeat from 1.




TODO
====

* howto add new scripts
* write systemd.service


EXAMPLES
========

```
mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"colorfade","value":0.5}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"colorfade","value":0.2}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":""}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"ceilingsinus"}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"ceilingsinus","b":{"amplitude":200,"offset":400}}'

```
