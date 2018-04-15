#!/usr/bin/python3

from ceilingspiffsconfig import writeConfig


writeConfig(
    ip="192.168.0.23",  #only used if dhcp is set to False or as fallback
    nm="255.255.255.0",
    gw="192.168.0.1",
#    dns=["8.8.8.8","1.1.1.1"], # you may want to use your local LANs DNS
#    dhcp0=False,
#    dhcp1=True,
    wifi0_ssid=b"", #first WIFI to try
    wifi0_pass=b"",
#    wifi1_ssid=b"", #second WIFI to try
#    wifi1_pass=b"",
#    wifi2_ssid=b"", #third WIFI to try
#    wifi2_pass=b"",
#    mqtt_port=1883,
    mqtt_broker=b"mqtt.realraum.at",
    mqtt_clientid=b"samplestrip",
    mqtt_user=b"",
    mqtt_pass=b"",
    authtoken=b"telnetpassword",
#    fan_threshold=30000,
#    simulate_cw_with_rgb=True,
    )
