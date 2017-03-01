#!/usr/bin/python3
# -*- coding: utf-8 -*-

import time
import math
import copy

triggername_ = "continue"

phase_=0
cs_default_={
    "r":{"offset":680,"amplitude":200,"lst":[]},
    "g":{"offset":0,"amplitude":0,"lst":[]},
    "b":{"offset":0,"amplitude":0,"lst":[]},
    "ww":{"offset":270,"amplitude":90,"lst":[]},
    "cc":{"offset":0,"amplitude":0,"lst":[]},
}
cs_=copy.deepcopy(cs_default_)
fade_duration_=20000

def updateList(scr):
    global cs_
    o=0
    for k in cs_.values():
        k["lst"] = list([k["offset"] + int(k["amplitude"]*math.sin(2*math.pi/scr.light_num*(i+o)%scr.light_num)) for i in range(0,scr.light_num)])
        o+=1

def activate(scr, newsettings):
    global cs_, fade_duration_
    cs_=copy.deepcopy(cs_default_)
    if "fadeduration" in newsettings and isinstance(newsettings["fadeduration"], int):
        fade_duration_= newsettings["fadeduration"]
    else:
        fade_duration_ = 20000
    for k,v in cs_.items():
        if k in newsettings:
            for kk in v.keys():
                if kk in newsettings[k]:
                    if isinstance(newsettings[k][kk],int) and newsettings[k][kk] >= 0 and newsettings[k][kk] <= 1000:
                        v[kk] = newsettings[k][kk]
    updateList(scr)
    animateAllLights(scr)

def deactivate(scr):
    pass

def loop(scr):
    pass

def triggerMeToContinue(scr):
    animateAllLights(scr)

def animateAllLights(scr):
    global phase_
    for i in reversed(range(0, scr.light_num)):
        kwargs = {}
        kwargs["fade_duration"]=fade_duration_
        if i == 0:
            kwargs["trigger_on_complete"]=[triggername_]
        idx = (i+phase_)%scr.light_num
        for k in cs_.keys():
            kwargs[k] = cs_[k]["lst"][idx]
        scr.setLight(i+1, **kwargs)
    phase_+=1
	
def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    #scr.registerLoop(loop)
    scr.registerTrigger(triggername_, triggerMeToContinue)    
