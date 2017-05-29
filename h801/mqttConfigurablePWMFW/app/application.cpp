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
#include "button.h"
#ifdef ENABLE_SSL
	#include <ssl/private_key.h>
	#include <ssl/cert.h>
#endif

DefaultLightConfigStorage DefaultLightConfig;

Timer BtnTimer;
DebouncedButton *button = nullptr;
bool button_used_ = false;
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
	debugf("Disconnected from %s. Reason: %d", ssid.c_str(), reason);

	if (!button_used_)
		flashSingleChannel(1,CHAN_RED);
	wifi_fail_count_++;
	if (wifi_fail_count_ > 2)
	{
		NetConfig.nextWifi();
		configureWifi();
		wifi_fail_count_ = 0;
	}
}

//////////////////////////////////////
////// Button             ////////////
//////////////////////////////////////

uint32_t button_color = 0;
uint32_t button_off_values[PWM_CHANNELS] = {0,0,0,0,0};
uint32_t button_on_values[PWM_CHANNELS] = {0,0,0,pwm_period/2,pwm_period/2};

void setupButton()
{
	button = new DebouncedButton(FUNC_GPIO0, 100, true);
}

void handleButton()
{
	const uint32_t levels_pro_color = 3;
	if (button == nullptr)
		return;
	if (button->isLongPressed())
	{
		debugf("btn longpress");
		for (uint8_t i=0;i<PWM_CHANNELS;i++)
		{
			button_on_values[i]=0;
		}
		button_on_values[button_color/levels_pro_color] = pwm_period/(levels_pro_color-button_color%levels_pro_color);
		applyValues(button_on_values);
		button_color=(button_color+1)%(PWM_CHANNELS*levels_pro_color);
	} else if (button->wasPressed()) {
		debugf("btn press");
		button_used_ = true;
		if (pwm_get_duty(CHAN_RED)+pwm_get_duty(CHAN_GREEN)+pwm_get_duty(CHAN_BLUE)+pwm_get_duty(CHAN_WW)+pwm_get_duty(CHAN_CW) > 0)
		{
			applyValues(button_off_values);
		} else {
			applyValues(button_on_values);
		}
	}
}

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
		debugf("trying to mount spiffs at %x, length %d", RBOOT_SPIFFS_0, SPIFF_SIZE);
		spiffs_mount_manual(RBOOT_SPIFFS_0, SPIFF_SIZE);
#else
		debugf("trying to mount spiffs at %x, length %d", 0x100000, SPIFF_SIZE);
		spiffs_mount_manual(0x100000, SPIFF_SIZE);
#endif
	} else {
#ifdef RBOOT_SPIFFS_1
		debugf("trying to mount spiffs at %x, length %d", RBOOT_SPIFFS_1, SPIFF_SIZE);
		spiffs_mount_manual(RBOOT_SPIFFS_1, SPIFF_SIZE);
#else
		debugf("trying to mount spiffs at %x, length %d", 0x300000, SPIFF_SIZE);
		spiffs_mount_manual(0x300000, SPIFF_SIZE);
#endif
	}
#else
	debugf("spiffs disabled");
#endif
	setupPWM(); //Init PWM with spiffs saved default settings
	telnetRegisterCmdsWithCommandHandler();
	//commandHandler.registerSystemCommands();
	// configure stuff that needs to be done before system is ready
	NetConfig.load(); //loads netsettings from fs

	setupButton();
	BtnTimer.initializeMs(600, handleButton).start();

	configureWifi();
	// Set system ready callback method
	WifiEvents.onStationGotIP(wifiConnectOk);
	WifiEvents.onStationDisconnect(wifiConnectFail);
}
