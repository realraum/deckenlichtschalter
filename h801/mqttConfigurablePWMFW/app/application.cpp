#include <user_config.h>
#include <SmingCore/SmingCore.h>
#include <defaultlightconfig.h>
#include <spiffsconfig.h>
#include <SmingCore/Debug.h>
#include <pwmchannels.h>
#include "application.h"
#include "lightcontrol.h"
#include "telnet.h"
#include "mqtt.h"
#ifdef ENABLE_BUTTON
#include "button.h"
#endif
#ifdef ENABLE_SSL
	#include <ssl/private_key.h>
	#include <ssl/cert.h>
#endif

DefaultLightConfigStorage DefaultLightConfig("defaultlight.conf");
DefaultLightConfigStorage ButtonLightConfig("btnlight.conf");

#ifdef ENABLE_BUTTON
Timer BtnTimer;
DebouncedButton *button = nullptr;
bool button_used_ = false;
#endif
uint8_t wifi_fail_count_ = 0;


///////////////////////////////////////
///// WIFI Stuff
///////////////////////////////////////

void configureWifi()
{
	WifiAccessPoint.enable(false);
	WifiStation.enable(true);
	// Serial.println("clientid: "+NetConfig.mqtt_clientid);
	// Serial.println("SSID: "+NetConfig.getWifiSSID());
	// Serial.println("WifiPass: "+NetConfig.getWifiPASS());
	WifiStation.setHostname(NetConfig.mqtt_clientid+".realraum.at");
	WifiStation.config(NetConfig.getWifiSSID(), NetConfig.getWifiPASS()); // Put you SSID and Password here
	WifiStation.enableDHCP(NetConfig.enabledhcp);
	if (!NetConfig.enabledhcp)
		WifiStation.setIP(NetConfig.ip,NetConfig.netmask,NetConfig.gw);
}

// Will be called when WiFi station was connected to AP
void wifiConnectOk(IPAddress ip, IPAddress mask, IPAddress gateway)
{
	// debugf("WiFi CONNECTED");
	Serial.println(ip.toString());
	// Serial.println(WifiStation.getIP().toString());
	startTelnetServer();
	startMqttClient();
	// Start publishing loop (also needed for mqtt reconnect)
	flashSingleChannel(1,CHAN_GREEN);
}

// Will be called when WiFi station timeout was reached
void wifiConnectFail(String ssid, uint8_t ssidLength, uint8_t *bssid, uint8_t reason)
{
	// The different reason codes can be found in user_interface.h. in your SDK.
	// debugf("Disconnected from %s. Reason: %d", ssid.c_str(), reason);

#ifdef ENABLE_BUTTON
	if (!button_used_)
		flashSingleChannel(1,CHAN_RED);
#endif
	wifi_fail_count_++;
	if (wifi_fail_count_ > 2)
	{
		NetConfig.nextWifi();
		configureWifi();
		wifi_fail_count_ = 0;
	}
}

#ifdef ENABLE_BUTTON
//////////////////////////////////////
////// Button             ////////////
//////////////////////////////////////

uint32_t button_color_ = 0;
bool button_longpress_inprogress_ = false;

void handleButton()
{
	const uint32_t levels_pro_color = 5;
	if (button == nullptr)
		return;
	if (button->isLongPressed())
	{
		for (uint8_t i=0;i<PWM_CHANNELS;i++)
		{
			effect_target_values_[i]=0;
		}
		effect_target_values_[button_color_/levels_pro_color] = (button_color_%levels_pro_color) * pwm_period/ levels_pro_color;
		applyValues(effect_target_values_);
		button_color_=(button_color_+1)%(PWM_CHANNELS*levels_pro_color);
		button_longpress_inprogress_ = true; //publish results later
	} else if (button->wasPressed()) {
		button_used_ = true;
		if (effect_target_values_[CHAN_RED] | effect_target_values_[CHAN_GREEN] | effect_target_values_[CHAN_BLUE] | effect_target_values_[CHAN_WW] | effect_target_values_[CHAN_UV] > 0)
		{
			//Switch OFF: save current values
			for (uint8_t i=0;i<PWM_CHANNELS;i++)
			{
				effect_target_values_[i]=0;
			}
		} else {
			//Switch ON: restore last saved values
			for (uint8_t i=0;i<PWM_CHANNELS;i++)
			{
				effect_target_values_[i]=button_on_values_[i];
			}			
		}
		applyValues(effect_target_values_);
		mqttPublishCurrentLightSetting();
	} else if (button_longpress_inprogress_) {
		button_longpress_inprogress_ = false;
		mqttPublishCurrentLightSetting();
	}
}
#endif

//////////////////////////////////////
////// Base System Stuff  ////////////
//////////////////////////////////////

void init()
{
	Serial.begin(115200);
	Serial.systemDebugOutput(true); // Allow debug print to serial
	Serial.commandProcessing(true);
	// Mount file system, in order to work with files
	int slot = rboot_get_current_rom();
#ifndef DISABLE_SPIFFS
	if (slot == 0) {
#ifdef RBOOT_SPIFFS_0
		// debugf("trying to mount spiffs at %x, length %d", RBOOT_SPIFFS_0, SPIFF_SIZE);
		spiffs_mount_manual(RBOOT_SPIFFS_0, SPIFF_SIZE);
#else
		// debugf("trying to mount spiffs at %x, length %d", 0x100000, SPIFF_SIZE);
		spiffs_mount_manual(0x100000, SPIFF_SIZE);
#endif
	} else {
#ifdef RBOOT_SPIFFS_1
		// debugf("trying to mount spiffs at %x, length %d", RBOOT_SPIFFS_1, SPIFF_SIZE);
		spiffs_mount_manual(RBOOT_SPIFFS_1, SPIFF_SIZE);
#else
		// debugf("trying to mount spiffs at %x, length %d", 0x300000, SPIFF_SIZE);
		spiffs_mount_manual(0x300000, SPIFF_SIZE);
#endif
	}
#else
	// debugf("spiffs disabled");
#endif
	setupPWM(); //Init PWM with spiffs saved default settings
	telnetRegisterCmdsWithCommandHandler();
	//commandHandler.registerSystemCommands();
	// configure stuff that needs to be done before system is ready
	NetConfig.load(); //loads netsettings from fs

#ifdef ENABLE_BUTTON
	//INIT Button
	button_on_values_[CHAN_WW] = pwm_period/2;
	ButtonLightConfig.load(button_on_values_);
	button = new DebouncedButton(FUNC_GPIO0, NetConfig.debounce_interval, NetConfig.debounce_interval_longpress, true);
	BtnTimer.initializeMs(NetConfig.debounce_button_timer_interval, handleButton).start();
#endif

	//INIT WIFI
	configureWifi();
	WifiEvents.onStationGotIP(wifiConnectOk);
	WifiEvents.onStationDisconnect(wifiConnectFail);
}
