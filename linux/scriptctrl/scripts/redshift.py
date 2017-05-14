#!/usr/bin/python3
# -*- coding: utf-8 -*-

from Pysolar import solar
import datetime

hsvvalue_=0.5
fade_duration_=120000
latitude_ = 47.065554
longitude_ = 15.450435
transition_high_ = 1.0  # SolarAltitude when light should start fading into warmwhite
transition_middle_ = -12.0   # Solar Altitude when light should be fully warm white
transition_low_ =  -20.0       # Solar Altitude when light should be warmwhite and red

def activate(scr, newsettings):
    global hsvvalue_, fade_duration_, latitude_, longitude_
    if "value" in newsettings and isinstance(newsettings["value"],(int,float)) and newsettings["value"] >= 0.0 and newsettings["value"] <= 1.0:
        hsvvalue_ = newsettings["value"]
    else:
        hsvvalue_ = 0.5
    if "fadeduration" in newsettings and isinstance(newsettings["fadeduration"], int):
        fade_duration_= min(120000,max(600,newsettings["fadeduration"]))
    else:
        fade_duration_ = 120000
    if "latitude" in newsettings and isinstance(newsettings["latitude"],(int,float)):
        latitude_ = newsettings["latitude"]
    if "longitude" in newsettings and isinstance(newsettings["longitude"],(int,float)):
        longitude_ = newsettings["longitude"]
    if "transition_high" in newsettings and isinstance(newsettings["transition_high"],(int,float)):
        transition_high_ = newsettings["transition_high"]
    if "transition_middle" in newsettings and isinstance(newsettings["transition_middle"],(int,float)):
        transition_middle_ = newsettings["transition_middle"]
    if "transition_low" in newsettings and isinstance(newsettings["transition_low"],(int,float)):
        transition_low_ = newsettings["transition_low"]

    for t in scr.participating:
        redshiftLight(scr, t, True)
        redshiftLight(scr, t)


def deactivate(scr):
    pass

### normal color temp ranges are
### 3000-4000K (3700K default) at night
### 5500-6500K (5500K default and higher during overcast days) during day
###
### @arg day_factor 1.0 .. full day, sun in zenith
### @arg day_factor 0.0 .. sun below horizon
### @arg day_factor -1.0 .. sun in negative zenith
### @arg value .. light intensity
### on full day we have full cold-white light
### on sundown we have full warm-white light
### further into the night, we add even more red
def calcColorFromDayLevel(day_factor, value):
    day_factor = min(1.0,max(-1.0,day_factor))
    r = 1000 * value * max(0.0, -1.0 * day_factor)
    b = 0
    cw = 1000 * value * max(0.0, day_factor)
    ww = max(0,1000 * value - cw - (r/3))
    return int(r), int(b), int(cw), int(ww)

def redshiftLight(scr, lightid, initial=False):
    solar_altitude = solar.GetAltitudeFast(latitude_, longitude_, datetime.datetime.utcnow())
    daylevel = 1.0
    if solar_altitude >= transition_high_:
        daylevel = 1.0
    elif solar_altitude <= transition_low_:
        daylevel = -1.0
    elif solar_altitude < transition_high_ and solar_altitude >= transition_middle_:
        daylevel = (solar_altitude - transition_middle_) / (transition_high_ - transition_middle_)
    elif solar_altitude > transition_low_ and solar_altitude < transition_middle_:
        daylevel = -1 * (solar_altitude - transition_middle_) / (transition_low_ - transition_middle_)
    r,b,cw,ww = calcColorFromDayLevel(daylevel, hsvvalue_)
    scr.setLight(lightid,r=r,g=0,b=b,cw=cw,ww=ww,
        fade_duration=None if initial else fade_duration_,
        trigger_on_complete=[] if initial else [lightid]
        )

def redshiftLightOnTrigger(scr, lightid):
    if lightid in scr.participating:
        redshiftLight(scr, lightid)

def mkTriggerClosure(lightid):
    return lambda scr: redshiftLightOnTrigger(scr, lightid)

def init(scr):
    scr.setDefaultParticipating(scr.lightidsceiling)
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    for t in scr.lightidsceiling:
        scr.registerTrigger(t,mkTriggerClosure(t))
