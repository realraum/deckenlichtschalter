#!/usr/bin/python3

import requests
import json
from urllib.parse import urlencode, quote_plus
import time
import paho.mqtt.client as mqtt

polling_intervall_s_=0.1
flash_update_intervall_s_=2.5
empty_epsilon_http_server="http://localhost:8080/get.lua?"

def _did_request_succeed(r):
    if "error" in r.__dict__:
        return r.error is None
    elif "status_code" in r.__dict__:
        return r.status_code in [requests.codes.ok, 302]
    else:
        assert False

def sendR3Message(client, topic, datadict, qos=0, retain=False):
    client.publish(topic, json.dumps(datadict), qos, retain)

def queryEmptyEpsilon(url,fields):
    r = requests.get(url+urlencode(fields, quote_via=quote_plus), allow_redirects=False)
    if _did_request_succeed(r):
        return r.json()
    else:
        return {}

alert_level_=""
def setAlert(alerttext):
    global alert_level_
    if alert_level_ == alerttext:
        return
    alert_level_=alerttext
    flashoptions={"duration":900,"repetitions":5}
    if alert_level_.startswith("YELLOW"):
        sendR3Message(mqtt_client_, "action/ceilingAll/light",{"b":0,"ww":80,"cw":0,"g":185,"r":1000})
        sendR3Message(mqtt_client_, "action/ceilingAll/light",{"b":0,"ww":80*0.7,"cw":0,"g":185*0.7,"r":1000*0.7,"flash":flashoptions})
    elif alert_level_.startswith("RED"):
        sendR3Message(mqtt_client_, "action/ceilingAll/light",{"r":1000,"b":0,"g":0,"ww":0,"cw":0})
        sendR3Message(mqtt_client_, "action/ceilingAll/light",{"r":1000*0.7,"b":0,"g":0,"ww":0,"cw":0,"flash":flashoptions})
    else:
        sendR3Message(mqtt_client_, "action/ceilingAll/light",{"r":0,"b":0,"g":0,"ww":200,"cw":200,"fade":{}})

def updateLights():
    pass

def actOnEmptyEpsilon():
    alertlevel = queryEmptyEpsilon(empty_epsilon_http_server,{"alertlevel":"getAlertLevel()"})
    print(alertlevel)
    if "alertlevel" in alertlevel:
        setAlert(alertlevel["alertlevel"])
    else:
        setAlert("")

# Start zmq connection to publish / forward sensor data
def initMQTT():
    client = mqtt.Client()
    client.connect("mqtt.realraum.at", 1883, keepalive=31)
    return client

if __name__ == '__main__':
    mqtt_client_ = None
    last_get_time=time.time()
    last_get_time2=time.time()
    try:
        mqtt_client_ = initMQTT()
        # mqtt_client_.start_loop()
        while True:
            if time.time() - last_get_time > polling_intervall_s_:
                actOnEmptyEpsilon()
                last_get_time = time.time()
            if time.time() - last_get_time2 > flash_update_intervall_s_:
                updateLights()
                last_get_time2 = time.time()
            mqtt_client_.loop()

    except Exception as e:
        traceback.print_exc()
    finally:
        if isinstance(mqtt_client_, mqtt.Client):
            # client_stop_loop()
            mqtt_client_.disconnect()


