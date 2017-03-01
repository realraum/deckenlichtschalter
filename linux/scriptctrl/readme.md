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


### EXAMPLES

```
mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"colorfade","value":0.5}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"colorfade","value":0.2}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":""}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"ceilingsinus"}'

mosquitto_pub -h mqtt.realraum.at -t action/ceilingscripts/activatescript -m '{"script":"ceilingsinus","b":{"amplitude":200,"offset":400}}'

```


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
            "description": "phase shift between colours"
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



Adding New Scripts
==================

* put new python script in folder ```./scripts```
* the file name will became the scriptname used to activate the script
* see ```./scripts/example.py```


Contents of your script
-----------------------

Each of your functions get passed an object of the ```ScriptClass``` which you can use 
to register callbacks or set the lights.

Your script **must** contain the ```init(scr)``` function, which is called on startup and which you should use to register your callbacks.

### ```def init(scr):```

define this function in your script and use it to register your callback funtions.

### Trigger yourself

The prefered method is to start a fade or flash animation on a ceiling light and have that light trigger you when the animation has finished.

To do this,
1. choose one or several triggernames
2. register them with your ScriptObject.
3. include your triggername in the list ```trigger_on_complete``` when calling ```setLight```


### or Loop

Alternatively you can register a callback function with ```registerLoop(callback)``` to be called in a loop.

* Avoid using time.sleep too much and make sure your function returns in a timely fashion.
* Avoid updating the ceiling lights on each loop as this will most assuredly overload the MQTT broker


Interface of ```ScriptClass```
------------------------------

### ```scr.registerActivate(callback)```

register a function ```def mycallback(scr, newsettings):```  to be called when the script is being activated

Note that this callback gets an additional argument ```newsettings``` which is the dict object of the
JSON object sent to ```action/ceilingscripts/activatescript```. Use this to customize each activation of your script
with different settings.

### ```scr.registerDeactivate(callback)```

register a function ```def mycallback(scr):``` to be called when the script is being deactivated. E.g. to fade to black.

### ```scr.registerLoop(callback)```

register a function ```def mycallback(scr):``` to be called in fast loop

### ```scr.registerTrigger(triggername, callback)```

register a function ```def mycallback(scr):``` to be called when a ceiling light animation is finished that was started with that trigger included in the ```trigger_on_complete``` list

### ```scr.setLight(light,r,g,b,ww,cw,fade_duration,flash_repetitions,cc,triger_on_complete)```

set a ceiling light or start a ceiling light animation

* ```light```: which light to set. (int [1..6] || string "All")
* ```r```: value of red (int [0...1000])
* ```g```: value of green (int [0...1000])
* ```b```: value of blue (int [0...1000])
* ```cw```: value of cold-white (int [0...1000])
* ```ww```: value of warm-white (int [0...1000])
* ```fade_duration```: fade to this colorvalues within this duration of milliseconds (int [100..60000])
* ```flash_repetitions```: flash this color in 600ms intervals for this many repetitions (int [1..10])
* ```cc```: when light is set, forward setting to next ceiling light in list (list of values valid for ```light```)
* ```trigger_on_complete```: when light is set, last light on cc list, will raise this script trigger (list of strings)


TODO
====

* write systemd.service


