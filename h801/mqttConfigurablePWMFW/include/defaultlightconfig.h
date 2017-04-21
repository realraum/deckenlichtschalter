#include <SmingCore/SmingCore.h>
#include "pwmchannels.h"
#ifndef INCLUDE_DEFAULTCONFIG_H_
#define INCLUDE_DEFAULTCONFIG_H_

const String DEFAULTLIGHT_SETTINGS_FILE = "defaultlight.conf";
const String NET_SETTINGS_FILE = "net.conf";
const String WIFISSID_SETTINGS_FILE = "wifi.ssid.conf";
const String WIFIPASS_SETTINGS_FILE = "wifi.pass.conf";
const String MQTTCLIENT_SETTINGS_FILE = "mqtt.clientid.conf";
const String MQTTUSER_SETTINGS_FILE = "mqtt.user.conf";
const String MQTTPASS_SETTINGS_FILE = "mqtt.pass.conf";
const String MQTTBROKER_SETTINGS_FILE = "mqttbroker.conf";
const String AUTHTOKEN_SETTINGS_FILE = "authtoken.conf";
const String FAN_SETTINGS_FILE = "fan.conf";
const String USEDHCP_SETTINGS_FILE = "dhcp.flag";
const String SIMULATE_CW_SETTINGS_FILE = "simulatecw.flag";

struct DefaultLightConfigStorage
{
	void load(uint32_t values[PWM_CHANNELS])
	{
		if (exist())
		{
			file_t f = fileOpen(DEFAULTLIGHT_SETTINGS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) values, PWM_CHANNELS*sizeof(uint32_t));
			fileClose(f);
		}
	}

	void save(uint32_t values[PWM_CHANNELS])
	{
		file_t f = fileOpen(DEFAULTLIGHT_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
		fileWrite(f, (void*) values, PWM_CHANNELS*sizeof(uint32_t));
		fileClose(f);
	}

	bool exist() { return fileExist(DEFAULTLIGHT_SETTINGS_FILE)	&& fileGetSize(DEFAULTLIGHT_SETTINGS_FILE) >= PWM_CHANNELS*sizeof(uint32_t); }
};

extern DefaultLightConfigStorage DefaultLightConfig;

struct NetConfigStorage
{
	IPAddress ip = IPAddress(192, 168, 127, 247);
	IPAddress netmask = IPAddress(255,255,255,0);
	IPAddress gw = IPAddress(192, 168, 127, 254);
	String wifi_ssid="realraum";
	String wifi_pass="";
	String mqtt_broker="mqtt.realraum.at";
	String mqtt_clientid="ceiling1";
	String mqtt_user;
	String mqtt_pass;
	bool enabledhcp=true;
	bool simulatecw_w_rgb=false;
	uint32_t mqtt_port=1883;  //8883 for ssl
	uint32_t fan_threshold=2500;
	String authtoken;

	void load()
	{
		if (exist())
		{
			uint32_t netsettings[4];
			file_t f = fileOpen(NET_SETTINGS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) netsettings, 4*sizeof(uint32_t));
			fileClose(f);
			ip = IPAddress(netsettings[0]);
			netmask = IPAddress(netsettings[1]);
			gw = IPAddress(netsettings[2]);
			mqtt_port = netsettings[3];
			wifi_ssid = fileGetContent(WIFISSID_SETTINGS_FILE);
			wifi_pass = fileGetContent(WIFIPASS_SETTINGS_FILE);
			mqtt_broker = fileGetContent(MQTTBROKER_SETTINGS_FILE);
			mqtt_clientid = fileGetContent(MQTTCLIENT_SETTINGS_FILE);
			mqtt_user = fileGetContent(MQTTUSER_SETTINGS_FILE);
			mqtt_pass = fileGetContent(MQTTPASS_SETTINGS_FILE);
			authtoken = fileGetContent(AUTHTOKEN_SETTINGS_FILE);
			enabledhcp = fileExist(USEDHCP_SETTINGS_FILE);
			simulatecw_w_rgb = fileExist(SIMULATE_CW_SETTINGS_FILE);
			f = fileOpen(FAN_SETTINGS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) &fan_threshold, sizeof(uint32_t));
			fileClose(f);
		}
	}

	void save()
	{
		uint32_t netsettings[4] = {ip,netmask,gw,mqtt_port};
		file_t f = fileOpen(NET_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
		fileWrite(f, (void*) netsettings, 4*sizeof(uint32_t));
		fileClose(f);
		fileSetContent(WIFISSID_SETTINGS_FILE, wifi_ssid);
		fileSetContent(WIFIPASS_SETTINGS_FILE, wifi_pass);
		fileSetContent(MQTTBROKER_SETTINGS_FILE, mqtt_broker);
		fileSetContent(MQTTCLIENT_SETTINGS_FILE, mqtt_clientid);
		fileSetContent(MQTTUSER_SETTINGS_FILE, mqtt_user);
		fileSetContent(MQTTPASS_SETTINGS_FILE, mqtt_pass);
		fileSetContent(AUTHTOKEN_SETTINGS_FILE, authtoken);
		if (enabledhcp)
			fileSetContent(USEDHCP_SETTINGS_FILE, "true");
		else
			fileDelete(USEDHCP_SETTINGS_FILE);
		if (simulatecw_w_rgb)
			fileSetContent(SIMULATE_CW_SETTINGS_FILE, "true");
		else
			fileDelete(SIMULATE_CW_SETTINGS_FILE);
		f = fileOpen(FAN_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
		fileWrite(f, (void*) &fan_threshold, sizeof(uint32_t));
		fileClose(f);
	}

	bool exist() { return fileExist(NET_SETTINGS_FILE); }

	String getMQTTTopic(String topic3, bool all=false)
	{
		return JSON_TOPIC1+((all) ? JSON_TOPIC2_ALL : mqtt_clientid)+topic3;
	}
};

extern NetConfigStorage NetConfig;

#endif /* INCLUDE_DEFAULTCONFIG_H_ */
