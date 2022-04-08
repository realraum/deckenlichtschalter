#!/usr/bin/python3
# -*- coding: utf-8 -*-

import importlib
import traceback
from interface import CeilingClass

def loadScripts(ceilingobj):
    scripts = importlib.import_module("scripts.__init__")
    for module_name in scripts.__all__:
        try:
            imported_module = importlib.import_module("scripts."+module_name)
            scr = ceilingobj.newScript(module_name)
            imported_module.init(scr)
        except:
            traceback.print_exc()
            ceilingobj.removeScript(module_name)


if __name__ == "__main__":
    ## init Ceiling and mqtt
    ceiling = CeilingClass()
    ## register functions
    loadScripts(ceiling)
    ## run main loop
    ceiling.mqttrun("mqtt.realraum.at", 1883, keepalive=31)
