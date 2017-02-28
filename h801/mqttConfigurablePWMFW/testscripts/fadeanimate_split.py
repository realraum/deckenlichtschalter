#!/usr/bin/python3
# -*- coding: utf-8 -*-

import json
import paho.mqtt.client as mqtt
import traceback
import random
import sys
import colorsys
import time

myclientid_ = "ceilinganimator"
mytopic_ = "action/"+myclientid_+"/continue"
alltopic_ = "action/ceilingAll/light"
hsvvalue_="random"

ceiling_clientids_= list(["ceiling%d" % x for x in range(1,7)])

def sendR3Message(client, topic, datadict, qos=0, retain=False):
    client.publish(topic, json.dumps(datadict), qos, retain)


def decodeR3Payload(payload):
    try:
        return json.loads(payload.decode("utf-8"))
    except Exception as e:
        print("Error decodeR3Payload:" + str(e))
        return {}


def animateSomeLightsOnce(client, clientlist, duration=1000, triggerme=False):
    targets = list(["action/%s/light" % cid for cid in clientlist])
    if triggerme:
        targets.append(mytopic_)
    r,g,b = colorsys.hsv_to_rgb(random.randint(0,1000)/1000.0,1,random.randint(300,1000)/1000.0 if hsvvalue_ == "random" else hsvvalue_)
    r *= 1000.0
    g *= 1000.0
    b *= 1000.0
    r = int(min(r,1000))
    g = int(min(g*4.0/5.0,800))
    b = int(min(b*2.0/3.0,666))
    msg = {"r":r,"g":g,"b":b,"cw":0, "ww":0,"fade":{"duration":duration, "cc":targets[1:]}}
    sendR3Message(client, targets[0],msg)

def animateAllLights():
    lst = list(ceiling_clientids_)
    random.shuffle(lst)
    lsthalf = int(len(lst)/2)
    duration=random.randint(6,50)*100
    animateSomeLightsOnce(client,lst[:lsthalf],duration)
    animateSomeLightsOnce(client,lst[lsthalf:],duration,True)

def onMQTTMessage(client, userdata, msg):
    animateAllLights()

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
        animateAllLights()

        ## now keep them running
        while True:
            client.loop()

    except Exception as e:
        traceback.print_exc()
    finally:
        print("Exiting ... ")
        sendR3Message(client, alltopic_ ,{"r":0,"g":0,"b":0,"cw":0,"ww":0,"fade":{"duration":1000,"cc":[]}})
        time.sleep(1.0)
        sendR3Message(client, alltopic_ ,{"r":0,"g":0,"b":0,"cw":0,"ww":0,"fade":{"duration":500,"cc":[alltopic_,alltopic_]}})
        if isinstance(client, mqtt.Client):
            client.disconnect()
	
