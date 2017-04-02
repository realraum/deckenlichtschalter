#!/usr/bin/python3
# -*- coding: utf-8 -*-

import traceback
import random
import sys
import colorsys
import time

mytrigger_ = "continue"
hsvvalue_="random"
fade_duration_=600
interval_=60.0
participating_targets_ = list(range(1, 9))


def activate(scr, newsettings):
    global hsvvalue_, fade_duration_, participating_targets_
    if "value" in newsettings and isinstance(newsettings["value"],float) and newsettings["value"] >= 0.0 and newsettings["value"] <= 1.0:
        hsvvalue_ = newsettings["value"]
    else:
        hsvvalue_ = "random"
    if "fadeduration" in newsettings and isinstance(newsettings["fadeduration"], int):
        fade_duration_= min(120000,max(600,newsettings["fadeduration"]))
    else:
        fade_duration_ = 600
    if "interval" in newsettings and isinstance(newsettings["interval"], int):
        interval_= min(1200,max(fade_duration_/1000,newsettings["interval"]))
    else:
        interval_ = 60.0
    if "participating" in newsettings and isinstance(newsettings["participating"],list) and all([isinstance(x,int) and x >=src.light_min for x in newsettings["participating"]]):
    	participating_targets_ = newsettings["participating"]
    else:
    	participating_targets_ = list(range(src.light_min, 9))
    colorAllLights(scr)

def deactivate(scr):
    pass

def colorALight(scr, targets, duration=1000):
    r,g,b = colorsys.hsv_to_rgb(random.randint(0,1000)/1000.0,1,random.randint(300,1000)/1000.0 if hsvvalue_ == "random" else hsvvalue_)
    r *= 1000.0
    g *= 1000.0
    b *= 1000.0
    r = int(min(r,1000))
    g = int(min(g*4.0/5.0,800))
    b = int(min(b*2.0/3.0,666))
    scr.setLight(targets[0],r=r,g=g,b=b,cw=0,ww=0,fade_duration=duration,cc=targets)

def colorAllLights(scr):
    lst = participating_targets_
    for l in lst:
        colorALight(scr, [l], fade_duration_)

last_run_ = 0
def loop(scr):
    global last_run_
    if time.time() - last_run_ > interval_:
        last_run_ = time.time()
        colorAllLights(scr)

def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    scr.registerLoop(loop) 
