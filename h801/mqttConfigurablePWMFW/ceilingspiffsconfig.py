#!/usr/bin/python3

from ipaddress import ip_address
import os
import struct

def writeConfig(ip,nm,gw,wifi0_ssid,wifi0_pass,mqtt_broker,mqtt_clientid,mqtt_user,mqtt_pass,authtoken,dhcp=True,mqtt_port=1883,fan_threshold=2000,simulate_cw_with_rgb=False,chan_ranges=[1000,1000,1000,1000,1000],wifi1_ssid="",wifi1_pass="",wifi2_ssid="",wifi2_pass="",debounce_interval=30,debounce_interval_longpress=700,debounce_button_timer_interval=800):
    NET_SETTINGS_FILE = "net.conf"
    WIFISSID0_SETTINGS_FILE = "wifi0.ssid"
    WIFIPASS0_SETTINGS_FILE = "wifi0.pass"
    WIFISSID1_SETTINGS_FILE = "wifi1.ssid"
    WIFIPASS1_SETTINGS_FILE = "wifi1.pass"
    WIFISSID2_SETTINGS_FILE = "wifi2.ssid"
    WIFIPASS2_SETTINGS_FILE = "wifi2.pass"
    MQTTCLIENT_SETTINGS_FILE = "mqtt.client"
    MQTTUSER_SETTINGS_FILE = "mqtt.user"
    MQTTPASS_SETTINGS_FILE = "mqtt.pass"
    MQTTBROKER_SETTINGS_FILE = "mqttbrkr.conf"
    AUTHTOKEN_SETTINGS_FILE = "authtoken"
    USEDHCP_SETTINGS_FILE = "dhcp.flag"
    SIMULATE_CW_SETTINGS_FILE = "simcw.flag"
    FAN_SETTINGS_FILE = "fan.conf"
    CHAN_RANGE_SETTINGS_FILE = "chanranges.conf"
    BUTTON_SETTINGS_FILE = "btn.conf"
    DIR="./files/"
    with open(os.path.join(DIR, NET_SETTINGS_FILE),"wb") as fh:
        fh.write(struct.pack(">III", int(ip_address(ip)), int(ip_address(nm)), int(ip_address(gw))))
        fh.write(struct.pack("<I",  int(mqtt_port)))
    with open(os.path.join(DIR, CHAN_RANGE_SETTINGS_FILE),"wb") as fh:
        fh.write(struct.pack("<IIIII", *map(int,chan_ranges)))
    with open(os.path.join(DIR, FAN_SETTINGS_FILE),"wb") as fh:
        fh.write(struct.pack("<I",  int(fan_threshold)))
    with open(os.path.join(DIR, WIFISSID0_SETTINGS_FILE),"wb") as fh:
        fh.write(wifi0_ssid)
    with open(os.path.join(DIR, WIFIPASS0_SETTINGS_FILE),"wb") as fh:
        fh.write(wifi0_pass)
    with open(os.path.join(DIR, WIFISSID1_SETTINGS_FILE),"wb") as fh:
        fh.write(wifi1_ssid)
    with open(os.path.join(DIR, WIFIPASS1_SETTINGS_FILE),"wb") as fh:
        fh.write(wifi1_pass)
    with open(os.path.join(DIR, WIFISSID2_SETTINGS_FILE),"wb") as fh:
        fh.write(wifi2_ssid)
    with open(os.path.join(DIR, WIFIPASS2_SETTINGS_FILE),"wb") as fh:
        fh.write(wifi2_pass)
    with open(os.path.join(DIR, MQTTBROKER_SETTINGS_FILE),"wb") as fh:
        fh.write(mqtt_broker)
    with open(os.path.join(DIR, MQTTCLIENT_SETTINGS_FILE),"wb") as fh:
        fh.write(mqtt_clientid)
    with open(os.path.join(DIR, MQTTUSER_SETTINGS_FILE),"wb") as fh:
        fh.write(mqtt_user)
    with open(os.path.join(DIR, MQTTPASS_SETTINGS_FILE),"wb") as fh:
        fh.write(mqtt_pass)
    with open(os.path.join(DIR, BUTTON_SETTINGS_FILE),"wb") as fh:
        fh.write(struct.pack("<III", int(debounce_interval), int(debounce_interval_longpress), int(debounce_button_timer_interval)))
    with open(os.path.join(DIR, AUTHTOKEN_SETTINGS_FILE),"wb") as fh:
        fh.write(authtoken)
    if dhcp:
        with open(os.path.join(DIR, USEDHCP_SETTINGS_FILE),"wb") as fh:
            fh.write(b"true")
    else:
        try:
            os.unlink(os.path.join(DIR, USEDHCP_SETTINGS_FILE))
        except:
            pass
    if simulate_cw_with_rgb:
        with open(os.path.join(DIR, SIMULATE_CW_SETTINGS_FILE),"wb") as fh:
            fh.write(b"true")
    else:
        try:
            os.unlink(os.path.join(DIR, SIMULATE_CW_SETTINGS_FILE))
        except:
            pass

## Example Use
# writeConfig(
#     ip="",
#     nm="",
#     gw="",
#     dhcp=True,
#     wifi0_ssid=b"",
#     wifi0_pass=b"",
#     mqtt_port=1883,
#     mqtt_broker=b"",
#     mqtt_clientid=b"",
#     mqtt_user=b"",
#     mqtt_pass=b"",
#     authtoken=b"")
