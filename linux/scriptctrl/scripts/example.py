#!/usr/bin/python3
# -*- coding: utf-8 -*-

#### REMEMBER
#### TO ADD 
#### YOUR SCRIPT
#### TO __all__
#### in __init__.py

import traceback
import random
import colorsys
import math

mytrigger_ = "nextanimation"

def animateCeiling(scr):
    whichlights=scr.lightidsceiling
    scr.setLight(whichlights[0],
            r=random.randint(0,1000),
            g=random.randint(0,1000),
            b=random.randint(0,1000),
            cw=random.randint(0,1000),
            ww=random.randint(0,1000),
            fade_duration=1000,
            cc=whichlights[1:],
            trigger_on_complete=[mytrigger_])

def activate(scr, newsettings):
    animateCeiling(scr)

def deactivate(scr):
    pass

def loop(scr):
    pass

def triggerMe(scr):
    animateCeiling(scr)

def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    scr.registerLoop(loop)
    scr.registerTrigger(mytrigger_, triggerMe)
