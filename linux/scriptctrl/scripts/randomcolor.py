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
interval_s_=60.0


def activate(scr, newsettings):
    global hsvvalue_, fade_duration_, interval_s_, last_run_
    if "value" in newsettings and isinstance(newsettings["value"],(int,float)) and newsettings["value"] >= 0.0 and newsettings["value"] <= 1.0:
        hsvvalue_ = newsettings["value"]
    else:
        hsvvalue_ = "random"
    if "fadeduration" in newsettings and isinstance(newsettings["fadeduration"], int):
        fade_duration_= min(120000,max(600,newsettings["fadeduration"]))
    else:
        fade_duration_ = 600
    if "interval" in newsettings and isinstance(newsettings["interval"], int):
        interval_s_= min(1200,max(fade_duration_/1000,newsettings["interval"]))
    else:
        interval_s_ = 60.0
    colorAllLights(scr)
    last_run_ = time.time()

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
    scr.setLight(targets[0],r=r,g=g,b=b,cw=0,ww=0,fade_duration=duration,cc=targets[1:])

def colorAllLights(scr):
    for l in scr.participating:
        colorALight(scr, [l], fade_duration_)

last_run_ = 0
def loop(scr):
    global last_run_
    if time.time() - last_run_ > interval_s_:
        last_run_ = time.time()
        colorAllLights(scr)

def init(scr):
    global participating_targets_
    scr.setDefaultParticipating(scr.lightids)
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    scr.registerLoop(loop)
