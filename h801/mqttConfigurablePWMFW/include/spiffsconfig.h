#ifndef INCLUDE_SPIFFSCONFIG_H_
#define INCLUDE_SPIFFSCONFIG_H_
#include "pwmchannels.h"

#define MAX_WIFI_SETS 3

class SpiffsConfigStorage
{
public:
	IPAddress ip = IPAddress(192, 168, 127, 242);
	IPAddress netmask = IPAddress(255,255,255,0);
	IPAddress gw = IPAddress(192, 168, 127, 254);
	uint32_t wifi_settings_num = 1;
	String wifi_ssid[MAX_WIFI_SETS]={"realraum","realstuff",""};
	String wifi_pass[MAX_WIFI_SETS]={"","",""};
	String mqtt_broker="mqtt.realraum.at";
	String mqtt_clientid="ceiling9";
	String mqtt_user;
	String mqtt_pass;
	bool enabledhcp=true;
	uint16_t mqtt_port=1883;  //8883 for ssl
	uint32_t publish_interval=8000;
	String authtoken;
	
	uint32_t debounce_interval=300;

	uint32_t fan_threshold=2500;
	bool simulatecw_w_rgb=false;
	uint32_t chan_range[PWM_CHANNELS] = {1000,1000,1000,1000,1000};

	void load();
	void save();
	bool exist();

	String getWifiSSID() {return wifi_ssid[wifi_settings_idx];}
	String getWifiPASS() {return wifi_pass[wifi_settings_idx];}
	void nextWifi() {wifi_settings_idx++; if (wifi_ssid[wifi_settings_idx].length() == 0) {wifi_settings_idx=0;}}

private:
	uint32_t wifi_settings_idx = 0;

};

extern SpiffsConfigStorage NetConfig;

#endif /* INCLUDE_DEFAULTCONFIG_H_ */
