#!/usr/bin/python3
# -*- coding: utf-8 -*-

import sys,os
import json
import time
import paho.mqtt.client as mqtt

ttydev="/dev/ttyACM0"

ttyfh=None
def localtty(code):
    global ttyfh
    if ttyfh is None:
      ttyfh=open(ttydev,"wb")
    ttyfh.write(b">"+code)
    ttyfh.flush()

mqttc=None
def mqtt_connect():
  global mqttc
  if mqttc is None:
    mqttc = mqtt.Client(client_id="rf433ctl@licht")
    mqttc.connect("mqtt.realraum.at", 1883, 60)

def mqttsendcode(code, presend_delay_ms=0):
  mqtt_connect()
  if presend_delay_ms > 0:
    mqttc.publish("action/rf433/setdelay", json.dumps({"Location":"pillar", "DelayMs":presend_delay}, retain=True),qos=1)
  mqttc.publish("action/rf433/sendcode3byte", json.dumps({"Code":list(code), "Ts":time.time()}),qos=1)

def mqttsendymhircmd(cmd):
  mqtt_connect()
  mqttc.publish("action/yamahastereo/ircmd", json.dumps({"Cmd":cmd, "Ts":time.time()}),qos=1)


def both(code):
  localtty(code)
  #mqttsendcode(code,2000)
  time.sleep(2)
  mqttsendcode(code)

rfcodes = {
  "regalleinwand":{"on":b"\xa2\xa0\xa8","off":b"\xa2\xa0\x28", "transmitter":localtty}, #white remote B 1
  "bluebar":{"on":b"\xa8\xa0\xa8","off":b"\xa8\xa0\x28", "transmitter":localtty}, #white remote C 1
  "labortisch":{"on":b"\xa2\xa2\xaa","off":b"\xa2\xa2\x2a", "transmitter":localtty},
  "couchred":{"on":b"\x8a\xa0\x8a","off":b"\x8a\xa0\x2a", "transmitter":localtty}, #pollin 00101 a
  "couchwhite":{"on":b"\x8a\xa8\x88","off":b"\x8a\xa8\x28", "transmitter":localtty}, #pollin 00101 d
  "cxleds":{"on":b"\x8a\x88\x8a","off":b"\x8a\x88\x2a", "transmitter":localtty}, #pollin 00101 b
  "mashadecke":{"on":b"\x8a\x28\x8a","off":b"\x8a\x28\x2a", "transmitter":localtty}, #pollin 00101 c
  "boiler":{"on":b"\xa0\xa2\xa8","off":b"\xa0\xa2\x28", "transmitter":both}, #white remote A 2
  "spots":{"on":b"\x00\xaa\x88","off":b"\x00\xaa\x28", "transmitter":localtty}, #polling 11110 d
  "lichtpi":{"on":b"\x00\xa2\x8a","off":b"\x00\xa2\x2a", "transmitter":localtty}, # Funksteckdose an welcher der RaspberryPi / Kiosk / Deckenlicht hängt
  "abwasch":{"on":b"\xaa\xa2\xa8","off":b"\xaa\xa2\x28", "transmitter":mqttsendcode},  #alte jk16 decke vorne

  "ymhpoweroff":{"on":"ymhpoweroff", "transmitter":mqttsendymhircmd},
  "ymhpower":{"on":"ymhpower", "off":"ymhpoweroff", "transmitter":mqttsendymhircmd},
  "ymhpoweron":{"on":"ymhpoweron", "transmitter":mqttsendymhircmd},
  "ymhcd":{"on":"ymhcd", "transmitter":mqttsendymhircmd},
  "ymhtuner":{"on":"ymhtuner", "transmitter":mqttsendymhircmd},
  "ymhtape":{"on":"ymhtape", "transmitter":mqttsendymhircmd},
  "ymhwdtv":{"on":"ymhwdtv", "transmitter":mqttsendymhircmd},
  "ymhsattv":{"on":"ymhsattv", "transmitter":mqttsendymhircmd},
  "ymhvcr":{"on":"ymhvcr", "transmitter":mqttsendymhircmd},
  "ymh7":{"on":"ymh7", "transmitter":mqttsendymhircmd},
  "ymhaux":{"on":"ymhaux", "transmitter":mqttsendymhircmd},
  "ymhextdec":{"on":"ymhextdec", "transmitter":mqttsendymhircmd},
  "ymhtest":{"on":"ymhtest", "transmitter":mqttsendymhircmd},
  "ymhtunabcde":{"on":"ymhtunabcde", "transmitter":mqttsendymhircmd},
  "ymheffect":{"on":"ymheffect", "transmitter":mqttsendymhircmd},
  "ymhtunplus":{"on":"ymhtunplus", "transmitter":mqttsendymhircmd},
  "ymhtunminus":{"on":"ymhtunminus", "transmitter":mqttsendymhircmd},
  "ymhvolup":{"on":"ymhvolup", "transmitter":mqttsendymhircmd},
  "ymhvoldown":{"on":"ymhvoldown", "transmitter":mqttsendymhircmd},
  "ymhvolmute":{"on":"ymhvolmute", "transmitter":mqttsendymhircmd},
  "ymhmenu":{"on":"ymhmenu", "transmitter":mqttsendymhircmd},
  "ymhplus":{"on":"ymhplus", "transmitter":mqttsendymhircmd},
  "ymhminus":{"on":"ymhminus", "transmitter":mqttsendymhircmd},
  "ymhtimelevel":{"on":"ymhtimelevel", "transmitter":mqttsendymhircmd},
  "ymhprgdown":{"on":"ymhprgdown", "transmitter":mqttsendymhircmd},
  "ymhprgup":{"on":"ymhprgup", "transmitter":mqttsendymhircmd},
  "ymhsleep":{"on":"ymhsleep", "transmitter":mqttsendymhircmd},
  "ymhp5":{"on":"ymhp5", "transmitter":mqttsendymhircmd},
}
#  "jk16decke":{"on:":"\xaa\xa0\xa8","off":"\xaa\xa0\x28"},

multinames = {
  "ambientlights":["bluebar","couchred","couchwhite","regalleinwand","cxleds","abwasch"],
  "all":list(set(rfcodes.keys()) - set(["lichtpi"]))
}

namestoswitch=[]

if len(sys.argv) > 2:
  if sys.argv[1] == "1":
    sys.argv[1] = "on"
  elif sys.argv[1] == "0":
    sys.argv[1] = "off"
  if sys.argv[2] in multinames:
    namestoswitch = multinames[sys.argv[2]]
  elif sys.argv[2] in rfcodes:
    namestoswitch = [sys.argv[2]]
  else:
    sys.exit(1)
  for rfname in namestoswitch:
    if sys.argv[1] in rfcodes[rfname]:
      code = rfcodes[rfname][sys.argv[1]]
      try:
        rfcodes[rfname]["transmitter"](code)
      except KeyError:
        localtty(code)

if not ttyfh is None:
  ttyfh.close()
if not mqttc is None:
  mqttc.loop()
  mqttc.disconnect()
