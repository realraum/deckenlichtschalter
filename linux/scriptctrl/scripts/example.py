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

def activate(scr, newsettings):
    pass

def deactivate(scr):
    pass

def loop(scr):
    pass

def triggerMe1(scr):
    pass

def triggerMe2(scr):
    pass

def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    scr.registerLoop(loop)
    scr.registerTrigger("FirstHalf", triggerMe1)
    scr.registerTrigger("SecondHalf", triggerMe2)
