#!/usr/bin/python3
# -*- coding: utf-8 -*-

from Pysolar import solar
import datetime

hsvvalue_=0.5
fade_duration_=120000
participating_targets_ = list(range(1, 7))
latitude_ = 47.065554
longitude_ = 15.450435
transition_high_ = 1.0  # SolarAltitude when light should start fading into warmwhite
transition_middle_ = -12.0   # Solar Altitude when light should be fully warm white
transition_low_ =  -20.0       # Solar Altitude when light should be warmwhite and red

def activate(scr, newsettings):
    global hsvvalue_, fade_duration_, participating_targets_, latitude_, longitude_
    if "value" in newsettings and isinstance(newsettings["value"],float) and newsettings["value"] >= 0.0 and newsettings["value"] <= 1.0:
        hsvvalue_ = newsettings["value"]
    else:
        hsvvalue_ = 0.5
    if "fadeduration" in newsettings and isinstance(newsettings["fadeduration"], int):
        fade_duration_= min(120000,max(600,newsettings["fadeduration"]))
    else:
        fade_duration_ = 120000
    if "participating" in newsettings and isinstance(newsettings["participating"],list) and all([isinstance(x,int) and x >=scr.light_min for x in newsettings["participating"]]):
        participating_targets_ = newsettings["participating"]
    else:
        participating_targets_ = list(range(scr.light_min, scr.light_max+1))
    if "latitude" in newsettings and isinstance(newsettings["latitude"],float):
        latitude_ = newsettings["latitude"]
    if "longitude" in newsettings and isinstance(newsettings["longitude"],float):
        longitude_ = newsettings["longitude"]
    if "transition_high" in newsettings and isinstance(newsettings["transition_high"],float):
        transition_high_ = newsettings["transition_high"]
    if "transition_middle" in newsettings and isinstance(newsettings["transition_middle"],float):
        transition_middle_ = newsettings["transition_middle"]
    if "transition_low" in newsettings and isinstance(newsettings["transition_low"],float):
        transition_low_ = newsettings["transition_low"]

    for t in participating_targets_:
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

def redshiftLight(scr, lightnum, initial=False):
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
    scr.setLight(lightnum,r=r,g=0,b=b,cw=cw,ww=ww,
        fade_duration=None if initial else fade_duration_,
        trigger_on_complete=[] if initial else ["c%d" % lightnum] )

def redshiftLightOnTrigger(scr, lightnum):
	if lightnum in participating_targets_:
		redshiftLight(scr, lightnum)

def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    for t in participating_targets_:
        scr.registerTrigger("c%d" % t,lambda scr: redshiftLightOnTrigger(scr, t))
