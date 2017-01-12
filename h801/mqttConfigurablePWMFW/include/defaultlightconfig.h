#include <SmingCore/SmingCore.h>

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
const String USEDHCP_SETTINGS_FILE = "dhcp.flag";

struct DefaultLightConfigStorage
{
	void load(uint32_t values[5])
	{
		if (exist())
		{
			file_t f = fileOpen(DEFAULTLIGHT_SETTINGS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) values, 5*sizeof(uint32_t));
			fileClose(f);
		}
	}

	void save(uint32_t values[5])
	{
		file_t f = fileOpen(DEFAULTLIGHT_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
		fileWrite(f, (void*) values, 5*sizeof(uint32_t));
		fileClose(f);
	}

	bool exist() { return fileExist(DEFAULTLIGHT_SETTINGS_FILE)	&& fileGetSize(DEFAULTLIGHT_SETTINGS_FILE) >= 5*sizeof(uint32_t); }
};

static DefaultLightConfigStorage DefaultLightConfig;

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
	uint32_t mqtt_port=1883;  //8883 for ssl
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
	}

	bool exist() { return fileExist(NET_SETTINGS_FILE); }
};

static NetConfigStorage NetConfig;

#endif /* INCLUDE_DEFAULTCONFIG_H_ */
