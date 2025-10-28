#!/usr/bin/python3
# -*- coding: utf-8 -*-

import traceback
import random
import colorsys
import time

flash_period_ms_=120
interval_s_=10.0


def activate(scr, newsettings):
    global flash_period_ms_, interval_s_
    if "interval" in newsettings and (isinstance(newsettings["interval"], float) or isinstance(newsettings["interval"], int)):
        interval_s_= max(0.1,newsettings["interval"])
    else:
        interval_s_ = 7.0

def deactivate(scr):
    pass

def sparkleALight(scr, target):
    r,g,b = colorsys.hsv_to_rgb(random.randint(0,1000)/1000.0,1,1)
    r *= 1000.0
    g *= 1000.0
    b *= 1000.0
    r = int(min(r,1000))
    g = int(min(g*4.0/5.0,800))
    b = int(min(b*2.0/3.0,666))
    cw=random.randint(0,30)
    scr.setLight(target,r=r,g=g,b=b,cw=cw,ww=0,flash_period=flash_period_ms_,flash_repetitions=1)

last_run_ = 0
def loop(scr):
    global last_run_
    if time.time() - last_run_ > interval_s_:
        last_run_ = time.time()
        sparkleALight(scr,random.choice(scr.participating))

def init(scr):
    global participating_targets_
    scr.setDefaultParticipating(scr.lightids)
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    scr.registerLoop(loop)
