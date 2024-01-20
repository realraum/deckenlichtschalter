#!/usr/bin/python3
# -*- coding: utf-8 -*-

import copy

triggername_ = "continue"

phase_=0
colourlist_default_ = [{"r":200,"ww":400,"cw":0,"b":0,"g":0},{"r":0,"ww":200,"cw":200,"b":0,"g":0},{"r":0,"ww":0,"cw":400,"b":20,"g":0},{"r":0,"ww":200,"cw":200,"b":0,"g":0}]
wavegroups_default_ = [["ceiling1","ceiling6", "memberregal"],["ceiling2","ceiling5","flooddoor"],["ceiling3","ceiling4","abwasch"]]
wavegroups_=copy.deepcopy(wavegroups_default_)
colourlist_=copy.deepcopy(colourlist_default_)
fade_duration_=20000

def validateColorlist(inputlist):
    try:
        return isinstance(inputlist, list) and all(
            map(
                lambda x: isinstance(x,dict) and all(map(lambda y: isinstance(y,int),x.values())) and all(map(lambda z:z in ["r","g","b","cw","ww","uv"],x.keys())),inputlist
                )
            )
    except:
        return False

def activate(scr, newsettings):
    global colourlist_, wavegroups_, fade_duration_
    colourlist_=copy.deepcopy(colourlist_default_)
    wavegroups_=copy.deepcopy(wavegroups_default_)
    if "fadeduration" in newsettings and isinstance(newsettings["fadeduration"], int):
        fade_duration_= min(60000,max(100,newsettings["fadeduration"]))
    elif "speed" in newsettings and isinstance(newsettings["speed"], int):
        fade_duration_ = 60000 - int(59.9*min(1000,max(0,newsettings["speed"])))
    else:
        fade_duration_ = 20000
    if "colourlist" in newsettings and validateColorlist(newsettings["colourlist"]):
        colourlist_= newsettings["colourlist"]
    if "colorlist" in newsettings and validateColorlist(newsettings["colorlist"]):
        colourlist_= newsettings["colorlist"]
    if "reversed" in newsettings:
        wavegroups_= list(reversed(wavegroups_))
    animateAllLights(scr,initial=True)

def deactivate(scr):
    pass

def loop(scr):
    pass

def triggerMeToContinue(scr):
    animateAllLights(scr)

def animateAllLights(scr, initial=False):
    global phase_
    trigger_set = False
    for grpidx in range(0,len(wavegroups_)):
        kwargs = {}
        if initial:
            kwargs["fade_duration"]=300
        else:
            kwargs["fade_duration"]=fade_duration_
        kwargs.update(colourlist_[(grpidx + phase_) % len(colourlist_)])
        for lname in wavegroups_[grpidx]:
            if trigger_set:
                kwargs["trigger_on_complete"]=[]
            else:
                kwargs["trigger_on_complete"]=[triggername_]
                trigger_set = True
            print(lname,kwargs)
            scr.setLight(lname, **kwargs)
    phase_=(phase_+1) % len(colourlist_)

def init(scr):
    scr.registerActivate(activate)
    scr.registerDeactivate(deactivate)
    #scr.registerLoop(loop)
    scr.setDefaultParticipating(scr.lightids)
    scr.registerTrigger(triggername_, triggerMeToContinue)
