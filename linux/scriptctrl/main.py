#!/usr/bin/python3
# -*- coding: utf-8 -*-

from interface import CeilingClass, myclientid
import importlib
import traceback

def loadScripts(ceiling):
    scripts = importlib.import_module("scripts.__init__")
    for module_name in scripts.__all__:
        try:
            foo = importlib.import_module("scripts."+module_name)
            scr = ceiling.newScript(module_name)
            foo.init(scr)
        except:
            traceback.print_exc()
            ceiling.removeScript(module_name)


if __name__ == "__main__":
    ## init Ceiling and mqtt
    ceiling = CeilingClass()
    ## register functions
    loadScripts(ceiling)
    ## run main loop
    ceiling.mqttrun("mqtt.realraum.at", 1883, keepalive=31)
