#include <SmingCore/SmingCore.h>
#include "spiffsconfig.h"

SpiffsConfigStorage NetConfig;

const String NET_SETTINGS_FILE = "net.conf";
const String WIFISSID_SETTINGS_FILES[MAX_WIFI_SETS] = {"wifi0.ssid","wifi1.ssid","wifi2.ssid"};
const String WIFIPASS_SETTINGS_FILES[MAX_WIFI_SETS] = {"wifi0.pass","wifi1.pass","wifi2.pass"};
const String MQTTCLIENT_SETTINGS_FILE = "mqtt.client";
const String MQTTUSER_SETTINGS_FILE = "mqtt.user";
const String MQTTPASS_SETTINGS_FILE = "mqtt.pass";
const String MQTTBROKER_SETTINGS_FILE = "mqttbrkr.conf";
const String AUTHTOKEN_SETTINGS_FILE = "authtoken";
const String BUTTON_SETTINGS_FILE = "btn.conf";
const String USEDHCP_SETTINGS_FILE[MAX_WIFI_SETS] = {"dhcp0.flag","dhcp1.flag","dhcp2.flag"};
const String DNS_SERVERS_FILE = "dns";
const String FAN_SETTINGS_FILE = "fan.conf";
const String SIMULATE_CW_SETTINGS_FILE = "simcw.flag";
const String CHAN_RANGE_SETTINGS_FILE = "chanranges.conf";


void SpiffsConfigStorage::load()
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
		mqtt_port = (uint16_t)(netsettings[3]);
		if (fileExist(DNS_SERVERS_FILE))
		{
			uint32_t dnssettings[DNS_MAX_SERVERS];
			f = fileOpen(DNS_SERVERS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) dnssettings, DNS_MAX_SERVERS*sizeof(uint32_t));
			fileClose(f);
			for (uint32_t d=0; d<DNS_MAX_SERVERS; d++)
				dns[d] = IPAddress(dnssettings[d]);
		}
		for (uint32_t wifi_settings_num=0; wifi_settings_num<MAX_WIFI_SETS; wifi_settings_num++)
		{
			wifi_ssid[wifi_settings_num] = fileGetContent(WIFISSID_SETTINGS_FILES[wifi_settings_num]);
			wifi_pass[wifi_settings_num] = fileGetContent(WIFIPASS_SETTINGS_FILES[wifi_settings_num]);
			enabledhcp[wifi_settings_num] = fileExist(USEDHCP_SETTINGS_FILE[wifi_settings_num]);
		}
		mqtt_broker = fileGetContent(MQTTBROKER_SETTINGS_FILE);
		mqtt_clientid = fileGetContent(MQTTCLIENT_SETTINGS_FILE);
		mqtt_user = fileGetContent(MQTTUSER_SETTINGS_FILE);
		mqtt_pass = fileGetContent(MQTTPASS_SETTINGS_FILE);
		authtoken = fileGetContent(AUTHTOKEN_SETTINGS_FILE);

		f = fileOpen(BUTTON_SETTINGS_FILE, eFO_ReadOnly);
		fileRead(f, (void*) &debounce_interval, sizeof(uint32_t));
		fileRead(f, (void*) &debounce_interval_longpress, sizeof(uint32_t));
		fileRead(f, (void*) &debounce_button_timer_interval, sizeof(uint32_t));
		fileClose(f);

		simulatecw_w_rgb = fileExist(SIMULATE_CW_SETTINGS_FILE);
		f = fileOpen(FAN_SETTINGS_FILE, eFO_ReadOnly);
		fileRead(f, (void*) &fan_threshold, sizeof(uint32_t));
		fileClose(f);
		if (fileExist(CHAN_RANGE_SETTINGS_FILE))
		{
			f = fileOpen(CHAN_RANGE_SETTINGS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) chan_range, PWM_CHANNELS*sizeof(uint32_t));
			fileClose(f);
		}
	}
}

void SpiffsConfigStorage::save()
{
	uint32_t netsettings[4] = {ip,netmask,gw, (uint32_t)mqtt_port};
	file_t f = fileOpen(NET_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
	fileWrite(f, (void*) netsettings, 4*sizeof(uint32_t));
	fileClose(f);
	uint32_t dnssettings[DNS_MAX_SERVERS];
	for (uint32_t d=0; d<DNS_MAX_SERVERS; d++)
		dnssettings[d]=(uint32_t)dns[d];
	f = fileOpen(DNS_SERVERS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
	fileWrite(f, (void*) dnssettings, DNS_MAX_SERVERS*sizeof(uint32_t));
	fileClose(f);
	for (uint32_t ws=0; ws<MAX_WIFI_SETS; ws++)
	{
		if (wifi_ssid[ws].length() > 0) {
			fileSetContent(WIFISSID_SETTINGS_FILES[ws], wifi_ssid[ws]);
			fileSetContent(WIFIPASS_SETTINGS_FILES[ws], wifi_pass[ws]);
		}
		if (enabledhcp[ws])
			fileSetContent(USEDHCP_SETTINGS_FILE[ws], "true");
		else
			fileDelete(USEDHCP_SETTINGS_FILE[ws]);
	}
	fileSetContent(MQTTBROKER_SETTINGS_FILE, mqtt_broker);
	fileSetContent(MQTTCLIENT_SETTINGS_FILE, mqtt_clientid);
	fileSetContent(MQTTUSER_SETTINGS_FILE, mqtt_user);
	fileSetContent(MQTTPASS_SETTINGS_FILE, mqtt_pass);
	fileSetContent(AUTHTOKEN_SETTINGS_FILE, authtoken);
	f = fileOpen(BUTTON_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
	fileWrite(f, (void*) &debounce_interval, sizeof(uint32_t));
	fileWrite(f, (void*) &debounce_interval_longpress, sizeof(uint32_t));
	fileWrite(f, (void*) &debounce_button_timer_interval, sizeof(uint32_t));
	fileClose(f);

	if (simulatecw_w_rgb)
		fileSetContent(SIMULATE_CW_SETTINGS_FILE, "true");
	else
		fileDelete(SIMULATE_CW_SETTINGS_FILE);
	f = fileOpen(FAN_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
	fileWrite(f, (void*) &fan_threshold, sizeof(uint32_t));
	fileClose(f);
	f = fileOpen(CHAN_RANGE_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
	fileWrite(f, (void*) chan_range, PWM_CHANNELS*sizeof(uint32_t));
	fileClose(f);


}

bool SpiffsConfigStorage::exist() { return fileExist(NET_SETTINGS_FILE); }