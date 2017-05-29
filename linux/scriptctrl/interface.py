#!/usr/bin/python3
# -*- coding: utf-8 -*-

import paho.mqtt.client as mqtt
import signal
import traceback
import json
import time
import sys
from collections import defaultdict

topic_presence = "realraum/metaevt/presence"
topic_action = "action/"
myclientid = "ceilingscripts"
format_ceiling_topic = "action/%s/light"
topic_base_scripts_ = topic_action+myclientid+"/script/"


########################################################
class CeilingScriptClass():
    def __init__(self, ceiling, scriptname):
        self.scriptname = scriptname
        self.ceiling = ceiling
        self.mybasetopic = topic_base_scripts_+self.scriptname+"/"
        self.triggers = {}
        self.deactivatefunc = None
        self.activatefunc = None
        self.loopfunc = None
        self.trigger_seq_num = 0
        self.trigger_expected_seq_num = defaultdict(int)
        self.default__participating = self.lightids
        self._participating = []

    @property
    def lightidall(self):
        return "ceilingAll"

    @property
    def lightids(self):
        return list(["ceiling%d" % x for x in range(1,9)] + ["abwasch","flooddoor"])

    @property
    def lightidsceiling(self):
        return list(["ceiling%d" % x for x in range(1,5)]+["flooddoor"]+["ceiling%d" % x for x in range(5,7)])

    @property
    def participating(self):
        return list(self._participating)

    def callcallback(self, client, trigger, msg):
        if trigger in self.triggers:
            payload={}
            try:
                payload = json.loads(msg.payload.decode("utf-8"))
            except Exception as e:
                 print("Exception in callcallback: ",e,file=sys.stderr)
            if "sq" in payload and self.trigger_expected_seq_num[trigger] == payload["sq"]:
                try:
                    self.triggers[trigger](self)
                except:
                    traceback.print_exc()
                    self.ceiling.removeScript(self.scriptname)

    def activate(self, newsettings):
        if not isinstance(newsettings, dict):
            return
        if "participating" in newsettings and isinstance(newsettings["participating"],list):
            self._participating = list([x for x in newsettings["participating"] if x in self.lightids])
        else:
            self._participating = self.default__participating
        if self.activatefunc:
            try:
                self.activatefunc(self, newsettings)
            except:
                traceback.print_exc()
                self.ceiling.removeScript(self.scriptname)

    def deactivate(self):
        if self.deactivatefunc:
            try:
                self.deactivatefunc(self)
            except:
                traceback.print_exc()
                self.ceiling.removeScript(self.scriptname)

    def loop(self):
        if self.loopfunc:
            try:
                self.loopfunc(self)
            except:
                traceback.print_exc()
                self.ceiling.removeScript(self.scriptname)

    def registerTrigger(self, triggername, callback):
        topic = self.mybasetopic + triggername
        self.ceiling.subscribe(topic)
        self.triggers[triggername] = callback
        return self

    def unregisterTrigger(self, triggername):
        topic = self.mybasetopic + triggername
        self.ceiling.unsubscribe(topic)
        del self.triggers[triggername]
        return self

    def registerLoop(self, callback):
        self.loopfunc = callback
        return self

    def registerActivate(self, activatefunc):
        self.activatefunc = activatefunc
        return self

    def registerDeactivate(self, deactivatefunc):
        self.deactivatefunc = deactivatefunc
        return self

    def setDefaultParticipating(self, lst):
        self.default__participating = lst
        self._participating = lst

    def setLight(self, light,r=None,g=None,b=None,cw=None,ww=None,fade_duration=None,flash_repetitions=None,cc=[],trigger_on_complete=[], include_scriptname=True):
        if not (light == self.lightidall) and not light in self._participating:
            return
        msg = {"r":r, "g":g, "b":b, "cw":cw, "ww": ww}
        ## sanity check
        for k,v in list(msg.items()):
            if not isinstance(v,int):
                del msg[k]
            else:
                msg[k] = max(min(v,1000),0)
        if isinstance(cc, list) and len(cc) < 9:
            cc = filter(lambda ccitem: ccitem == self.lightidall or (ccitem in self.lightids), cc)
            cc = map(str, cc)
            cc = list([format_ceiling_topic % ccitem for ccitem in cc])
        else:
            cc = []
        if isinstance(trigger_on_complete, list) and len(trigger_on_complete) < 3:
            trigger_on_complete = list(filter(lambda ccitem: ccitem in self.triggers.keys(), trigger_on_complete))
            for trigger in trigger_on_complete:
                self.trigger_expected_seq_num[trigger] = self.trigger_seq_num
            trigger_on_complete = list([self.mybasetopic+ccitem for ccitem in trigger_on_complete])
            cc += trigger_on_complete
            msg["sq"] = self.trigger_seq_num
            self.trigger_seq_num = (self.trigger_seq_num + 1) % (1<<30) # ensure seq number fits in signed int32
        if include_scriptname:
            msg["s"]=self.scriptname
        if light == self.lightidall:
            cc=None  # we don't want to be triggerd by x lights at once
        if fade_duration != None and fade_duration >= 100 and fade_duration <= 120000:
            msg["fade"]={"duration":fade_duration, "cc":cc}
        elif flash_repetitions != None and flash_repetitions >= 1 and flash_repetitions <= 10:
            msg["fade"]={"repetitions":flash_repetitions, "cc":cc}
        self.ceiling.client.publish(format_ceiling_topic % light, json.dumps(msg), 0, False)
        return self



########################################################
class CeilingClass():
    def __init__(self):
        self.client = mqtt.Client(client_id=myclientid)
        self.mybasetopic = topic_action+myclientid+"/"
        self._scripts = {}
        self._active_script = None
        self._subscribed_topics = {self.mybasetopic+"activatescript":True, topic_presence:True}
        self.keep_running = True

    def onmqttmsg(self, client, userdata, msg):
        if msg.topic == topic_presence:
            try:
                payload = json.loads(msg.payload.decode("utf-8"))
                if payload["Present"] == False:
                    self.deactivateCurrentScript()
            except Exception as e:
                print("onmqttmsg json2",e,file=sys.stderr)
            return
        if msg.topic == self.mybasetopic+"activatescript":
            self.deactivateCurrentScript()
            try:
                payload = json.loads(msg.payload.decode("utf-8"))
                self.activateScript(payload["script"], payload)
            except Exception as e:
                print("onmqttmsg json1",e,file=sys.stderr)
            return
        if msg.topic.startswith(topic_base_scripts_):
            script, trigger = msg.topic[len(topic_base_scripts_):].split("/")
            if script != self._active_script:
                return
            if script in self._scripts:
                self._scripts[script].callcallback(client, trigger, msg)

    def onmqttconnect(self, client, userdata, flags, rc):
        client.subscribe(list([(t,2) for t in self._subscribed_topics.keys()]))

    def onmqttdisconnect(self, mqttc, userdata, rc):
        if rc != 0:
            print("Unexpected disconnection.",file=sys.stderr)
            while True:
                time.sleep(5)
                print("Attempting reconnect",file=sys.stderr)
                try:
                    mqttc.reconnect()
                    break
                except ConnectionRefusedError:
                    continue
        else:
            print("Clean disconnect.",file=sys.stderr)
            sys.exit()

    def signal_handler(self, signal, frame):
        print('You pressed Ctrl+C!',file=sys.stderr)
        self.keep_running=False

    def mqttrun(self,*args,**kwargs):
        signal.signal(signal.SIGINT, self.signal_handler)
        self.client.on_message=self.onmqttmsg
        self.client.on_connect=self.onmqttconnect
        self.client.on_disconnect=self.onmqttdisconnect
        self.client.connect(*args, **kwargs)
        while self.keep_running:
            self.client.loop()
            if self._active_script:
                self._scripts[self._active_script].loop()
        self.deactivateCurrentScript()
        offscript = CeilingScriptClass(self,"off")
        offscript.setLight(offscript.lightidall,r=0,g=0,b=0,cw=0,ww=0,include_scriptname=False)

    def subscribe(self, topic):
        self._subscribed_topics[topic] = True
        if self.client:
            self.client.subscribe([(topic,2)])

    def unsubscribe(self, topic):
        del self._subscribed_topics[topic]
        if self.client:
            self.client.unsubscribe([(topic,2)])

    def deactivateCurrentScript(self):
        if self._active_script:
            script = self._active_script
            self._active_script = None
            self._scripts[script].deactivate()
            time.sleep(0.1)
            for l in self._scripts[script].participating:
                self._scripts[script].setLight(l,r=0,g=0,b=0,cw=0,ww=0,fade_duration=1000,include_scriptname=False)
            time.sleep(0.4)
            for l in self._scripts[script].participating:
                self._scripts[script].setLight(l,r=0,g=0,b=0,cw=0,ww=0,include_scriptname=False)

    def activateScript(self, script, newsettings):
        if script in self._scripts:
            self._active_script = script
            self._scripts[script].activate(newsettings)

    def newScript(self, scriptname):
        newscript = CeilingScriptClass(self, scriptname)
        self._scripts[scriptname] = newscript
        return newscript

    def removeScript(self, scriptname):
        if scriptname == self._active_script:
            self.deactivateCurrentScript()
        try:
            del self._scripts[scriptname]
        except:
            pass
