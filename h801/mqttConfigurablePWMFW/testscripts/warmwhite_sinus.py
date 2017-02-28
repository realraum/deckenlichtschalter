#!/usr/bin/python3
# -*- coding: utf-8 -*-

import json
import paho.mqtt.client as mqtt
import traceback
import random
import sys
import colorsys
import time
import math
import signal

keep_running=True
def signal_handler(signal, frame):
    global keep_running, client
    print('You pressed Ctrl+C!')
    keep_running=False
    sys.exit(0)

signal.signal(signal.SIGINT, signal_handler)

myclientid_ = "ceilinganimator3"
mytopic_ = "action/"+myclientid_+"/continue"
alltopic_ = "action/ceilingAll/light"

ceiling_clientids_= list(["action/ceiling%d/light" % x for x in range(1,7)])

def sendR3Message(client, topic, datadict, qos=0, retain=False):
    client.publish(topic, json.dumps(datadict), qos, retain)

def decodeR3Payload(payload):
    try:
        return json.loads(payload.decode("utf-8"))
    except Exception as e:
        print("Error decodeR3Payload:" + str(e))
        return {}

phase_=0
ww_offset=270
r_offset=680
ww_amp=90
r_amp=200
ww_lst = list([ww_offset + int(ww_amp*math.sin(2*math.pi/len(ceiling_clientids_)*i)) for i in range(0,len(ceiling_clientids_))])
r_lst = list([r_offset + int(r_amp*math.sin(2*math.pi/len(ceiling_clientids_)*i)) for i in range(0,len(ceiling_clientids_))])
print(ww_lst)
print(r_lst)
def animateSomeLightsOnce(client, phase, duration=20000):
    targets=[]
    for i in reversed(range(0,len(ceiling_clientids_))):
        if i == 0:
            targets=[mytopic_]
        msg = {"r":r_lst[(i+phase)%len(r_lst)],"g":0,"b":0,"cw":0, "ww":ww_lst[(i+3+phase)%len(ww_lst)],"fade":{"duration":duration, "cc":targets}}
        print(ceiling_clientids_[i], msg)
        sendR3Message(client, ceiling_clientids_[i], msg)

def animateAllLights(client):
    global phase_
    animateSomeLightsOnce(client, phase_)
    phase_+=1

def onMQTTMessage(client, userdata, msg):
    animateAllLights(client)

def onMQTTDisconnect(mqttc, userdata, rc):
    if rc != 0:
        print("Unexpected disconnection.")
        while True:
            time.sleep(5)
            print("Attempting reconnect")
            try:
                mqttc.reconnect()
                break
            except ConnectionRefusedError:
                continue
    else:
        print("Clean disconnect.")
        sys.exit()


def initMQTT():
    client = mqtt.Client(client_id=myclientid_)
    client.on_connect = lambda client, userdata, flags, rc: client.subscribe(
        [(mytopic_, 2)]
    )
    client.connect("mqtt.realraum.at", 1883, keepalive=31)
    client.on_message = onMQTTMessage
    client.on_disconnect = onMQTTDisconnect
    return client


if __name__ == '__main__':
    if len(sys.argv) > 1:
        hsvvalue_ = float(sys.argv[1])
    client = None
    try:
        client = initMQTT()

        ## kick things off
        animateAllLights(client)

        ## now keep them running
        while keep_running:
            client.loop()
            time.sleep(0.1)

    except Exception as e:
        traceback.print_exc()
    finally:
        print("Exiting ... ")
        sendR3Message(client, alltopic_ ,{"r":0,"g":0,"b":0,"cw":0,"ww":0,"fade":{"duration":1000,"cc":[]}})
        time.sleep(1.0)
        sendR3Message(client, alltopic_ ,{"r":0,"g":0,"b":0,"cw":0,"ww":0,"fade":{"duration":500,"cc":[alltopic_,alltopic_]}})
        if isinstance(client, mqtt.Client):
            client.disconnect()
	
