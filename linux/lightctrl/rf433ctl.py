#!/usr/bin/python3
# -*- coding: utf-8 -*-

import sys,os
import json
import time
import paho.mqtt.client as mqtt
import Pyro4

class SwitchARealSwitch():
    def __init__(self, ttydev, mqtthost, mqttport):
        self.ttydev = ttydev
        self.mqtthost = mqtthost
        self.mqttport = mqttport
        self.ttyfh=None
        self.mqttcon=None

        self.rfcodes = {
          "regalleinwand":{"on":b"\xa2\xa0\xa8","off":b"\xa2\xa0\x28", "transmitter":self.rfcode2TTY}, #white remote B 1
          "bluebar":{"on":b"\xa8\xa0\xa8","off":b"\xa8\xa0\x28", "transmitter":self.rfcode2TTY}, #white remote C 1
          "labortisch":{"on":b"\xa2\xa2\xaa","off":b"\xa2\xa2\x2a", "transmitter":self.rfcode2TTY},
          "couchred":{"on":b"\x8a\xa0\x8a","off":b"\x8a\xa0\x2a", "transmitter":self.rfcode2TTY}, #pollin 00101 a
          "couchwhite":{"on":b"\x8a\xa8\x88","off":b"\x8a\xa8\x28", "transmitter":self.rfcode2TTY}, #pollin 00101 d
          "cxleds":{"on":b"\x8a\x88\x8a","off":b"\x8a\x88\x2a", "transmitter":self.rfcode2TTY}, #pollin 00101 b
          "mashadecke":{"on":b"\x8a\x28\x8a","off":b"\x8a\x28\x2a", "transmitter":self.rfcode2TTY}, #pollin 00101 c
          "boiler":{"on":b"\xa0\xa2\xa8","off":b"\xa0\xa2\x28", "transmitter":self.rfcode2BOTH}, #white remote A 2
          "spots":{"on":b"\x00\xaa\x88","off":b"\x00\xaa\x28", "transmitter":self.rfcode2TTY}, #polling 11110 d
          "lichtpi":{"on":b"\x00\xa2\x8a","off":b"\x00\xa2\x2a", "transmitter":self.rfcode2TTY}, # Funksteckdose an welcher der RaspberryPi / Kiosk / Deckenlicht hängt
          "abwasch":{"on":b"\xaa\xa2\xa8","off":b"\xaa\xa2\x28", "transmitter":self.rfcode2MQTT},  #alte jk16 decke vorne

          "ymhpoweroff":{"on":"ymhpoweroff", "transmitter":self.ymhircmd2MQTT},
          "ymhpower":{"on":"ymhpower", "off":"ymhpoweroff", "transmitter":self.ymhircmd2MQTT},
          "ymhpoweron":{"on":"ymhpoweron", "transmitter":self.ymhircmd2MQTT},
          "ymhcd":{"on":"ymhcd", "transmitter":self.ymhircmd2MQTT},
          "ymhtuner":{"on":"ymhtuner", "transmitter":self.ymhircmd2MQTT},
          "ymhtape":{"on":"ymhtape", "transmitter":self.ymhircmd2MQTT},
          "ymhwdtv":{"on":"ymhwdtv", "transmitter":self.ymhircmd2MQTT},
          "ymhsattv":{"on":"ymhsattv", "transmitter":self.ymhircmd2MQTT},
          "ymhvcr":{"on":"ymhvcr", "transmitter":self.ymhircmd2MQTT},
          "ymh7":{"on":"ymh7", "transmitter":self.ymhircmd2MQTT},
          "ymhaux":{"on":"ymhaux", "transmitter":self.ymhircmd2MQTT},
          "ymhextdec":{"on":"ymhextdec", "transmitter":self.ymhircmd2MQTT},
          "ymhtest":{"on":"ymhtest", "transmitter":self.ymhircmd2MQTT},
          "ymhtunabcde":{"on":"ymhtunabcde", "transmitter":self.ymhircmd2MQTT},
          "ymheffect":{"on":"ymheffect", "transmitter":self.ymhircmd2MQTT},
          "ymhtunplus":{"on":"ymhtunplus", "transmitter":self.ymhircmd2MQTT},
          "ymhtunminus":{"on":"ymhtunminus", "transmitter":self.ymhircmd2MQTT},
          "ymhvolup":{"on":"ymhvolup", "transmitter":self.ymhircmd2MQTT},
          "ymhvoldown":{"on":"ymhvoldown", "transmitter":self.ymhircmd2MQTT},
          "ymhvolmute":{"on":"ymhvolmute", "transmitter":self.ymhircmd2MQTT},
          "ymhmenu":{"on":"ymhmenu", "transmitter":self.ymhircmd2MQTT},
          "ymhplus":{"on":"ymhplus", "transmitter":self.ymhircmd2MQTT},
          "ymhminus":{"on":"ymhminus", "transmitter":self.ymhircmd2MQTT},
          "ymhtimelevel":{"on":"ymhtimelevel", "transmitter":self.ymhircmd2MQTT},
          "ymhprgdown":{"on":"ymhprgdown", "transmitter":self.ymhircmd2MQTT},
          "ymhprgup":{"on":"ymhprgup", "transmitter":self.ymhircmd2MQTT},
          "ymhsleep":{"on":"ymhsleep", "transmitter":self.ymhircmd2MQTT},
          "ymhp5":{"on":"ymhp5", "transmitter":self.ymhircmd2MQTT},
        }
        #  "jk16decke":{"on:":"\xaa\xa0\xa8","off":"\xaa\xa0\x28"},

    def tty(self):
        if self.ttyfh is None:
            self.ttyfh=open(self.ttydev,"wb")
        return self.ttyfh

    def mqttc(self):
        if self.mqttcon is None:
            self.mqttcon = mqtt.Client(client_id="rf433ctl@licht")
            self.mqttcon.connect(self.mqtthost, self.mqttport, 60)
            self.mqttcon.loop_start()
        return self.mqttcon

    def cleanup(self):
        if not self.ttyfh is None:
            self.ttyfh.close()
        if not self.mqttcon is None:
            self.mqttcon.loop_stop()
            self.mqttcon.disconnect()

    def rfcode2TTY(self, code):
        self.tty().write(b">"+code)
        self.tty().flush()

    def rfcode2MQTT(self, code, presend_delay_ms=0):
        if presend_delay_ms > 0:
            self.mqttc().publish("action/rf433/setdelay", json.dumps({"Location":"pillar", "DelayMs":presend_delay}, retain=True),qos=1)
        self.mqttc().publish("action/rf433/sendcode3byte", json.dumps({"Code":list(code), "Ts":time.time()}),qos=1)

    def ymhircmd2MQTT(self, cmd):
        self.mqttc().publish("action/yamahastereo/ircmd", json.dumps({"Cmd":cmd, "Ts":time.time()}),qos=1)

    def rfcode2BOTH(self, code):
        self.rfcode2TTY(code)
        time.sleep(2000)
        self.rfcode2MQTT(code)

    def toggleSwitch(self, onoff, rfname):
        try:
            code = self.rfcodes[rfname][onoff]
        except KeyError:
            return
        try:
            self.rfcodes[rfname]["transmitter"](code)
        except KeyError:
            self.rfcode2TTY(code)


class MultiSwitcherQueue():
    def __init__(self, switcher):
        self.switcher=switcher

        self.multinames = {
          "ambientlights":["bluebar","couchred","couchwhite","regalleinwand","cxleds","abwasch"],
          "all":list(set(switcher.rfcodes.keys()) - set(["lichtpi"]))
        }

    def toggleSwitch(self, onoff, rfname):
        if rfname in self.multinames:
            for rfname2 in self.multinames[rfname]:
                self.switcher.toggleSwitch(onoff, rfname2)
        else:
            self.switcher.toggleSwitch(onoff, rfname)



switcher = SwitchARealSwitch("/dev/ttyACM0","mqtt.realraum.at",1883)
multiswitcher = MultiSwitcherQueue(switcher)


if len(sys.argv) > 2:
    if sys.argv[1] == "1":
        sys.argv[1] = "on"
    elif sys.argv[1] == "0":
        sys.argv[1] = "off"
    multiswitcher.toggleSwitch(sys.argv[1], sys.argv[2])
    switcher.cleanup()
