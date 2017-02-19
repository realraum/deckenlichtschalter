#!/usr/bin/python3
# -*- coding: utf-8 -*-

import json
import paho.mqtt.client as mqtt
import traceback
import random
import sys

myclientid_ = "ceilinganimator"
mytopic_ = "action/"+myclientid_+"/continue"

ceiling_clientids_= list(["ceiling%d" % x for x in range(1,7)])

def sendR3Message(client, topic, datadict, qos=0, retain=False):
    client.publish(topic, json.dumps(datadict), qos, retain)


def decodeR3Payload(payload):
    try:
        return json.loads(payload.decode("utf-8"))
    except Exception as e:
        print("Error decodeR3Payload:" + str(e))
        return {}


def animateSixLightsOnce(client):
    targets = list(["action/%s/light" % cid for cid in ceiling_clientids_])
    random.shuffle(targets)
    targets.append(mytopic_)
    duration=random.randint(5,15)*1000
    duration=1000
    msg = {"r":random.randint(0,1000),"g":random.randint(0,1000),"b":random.randint(0,1000),"cw":0, "ww":random.randint(0,5),"fade":{"duration":duration, "cc":targets[1:]}}
    sendR3Message(client, targets[0],msg)


def onMQTTMessage(client, userdata, msg):
    animateSixLightsOnce(client) #again


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
    client = None
    try:
        client = initMQTT()

        ## kick things off
        animateSixLightsOnce(client)

        ## now keep them running
        while True:
            client.loop()

    except Exception as e:
        traceback.print_exc()
    finally:
        if isinstance(client, mqtt.Client):
            client.disconnect()
