#!/usr/bin/python3
# -*- coding: utf-8 -*-

import traceback
import random
import sys
import colorsys
import time

mytrigger_ = "continue"
hsvvalue_="random"

def activate(scr, newsettings):
    global hsvvalue_
    if "value" in newsettings and isinstance(newsettings["value"],float) and newsettings["value"] >= 0.0 and newsettings["value"] <= 1.0:
        hsvvalue_ = newsettings["value"]
    else:
        hsvvalue_ = "random"
    animateAllLights(scr)

def deactivate(scr):
    pass

def animateSomeLightsOnce(scr, targets, duration=1000, triggerme=False):
    triggercc=[]
    if triggerme:
        triggercc=[mytrigger_]    
    r,g,b = colorsys.hsv_to_rgb(random.randint(0,1000)/1000.0,1,random.randint(300,1000)/1000.0 if hsvvalue_ == "random" else hsvvalue_)
    r *= 1000.0
    g *= 1000.0
    b *= 1000.0
    r = int(min(r,1000))
    g = int(min(g*4.0/5.0,800))
    b = int(min(b*2.0/3.0,666))
    scr.setLight(targets[0],r=r,g=g,b=b,cw=0,ww=0,fade_duration=duration,cc=targets[1:],trigger_on_complete=triggercc)

def animateAllLights(scr):
    lst = list(range(scr.light_min, scr.light_max+1))
    random.shuffle(lst)
    lsthalf = int(len(lst)/2)
    duration=random.randint(6,50)*100
    animateSomeLightsOnce(scr,lst[:lsthalf],duration)
    animateSomeLightsOnce(scr,lst[lsthalf:],duration,True)

def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    #scr.registerLoop(loop)
    scr.registerTrigger(mytrigger_, animateAllLights)    
