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

void configureDNSServer(uint8_t n, IPAddress a)
{
	ip_addr_t adr;
	ip4_addr_set_u32(&adr, (uint32_t) a);
	dns_setserver(n,&adr);
}

void configureWifi()
{
	// ip_addr_t ipdnsr3_;
	// ip_addr_t ipdnsffgraz_;
	// ip4_addr_set_u32(&ipdnsr3_, 0x21D36A59); //89.106.211.33
	// ip4_addr_set_u32(&ipdnsffgraz_,0x0A000B0A); //10.12.0.10
	// dns_setserver(0,&ipdnsr3_);
	// dns_setserver(1,&ipdnsffgraz_);
	for (uint8_t d=0; d<DNS_MAX_SERVERS;d++)
		configureDNSServer(d,NetConfig.dns[d]);
	WifiAccessPoint.enable(false);
	WifiStation.enable(true);
	// Serial.println("clientid: "+NetConfig.mqtt_clientid);
	// Serial.println("SSID: "+NetConfig.getWifiSSID());
	// Serial.println("WifiPass: "+NetConfig.getWifiPASS());
	WifiStation.setHostname(NetConfig.mqtt_clientid+".realraum.at");
	WifiStation.config(NetConfig.getWifiSSID(), NetConfig.getWifiPASS()); // Put you SSID and Password here
	WifiStation.enableDHCP(NetConfig.getEnableDHCP());
	if (!NetConfig.getEnableDHCP())
		WifiStation.setIP(NetConfig.ip,NetConfig.netmask,NetConfig.gw);
}

// Will be called when WiFi station was connected to AP
void wifiConnectOk(IPAddress ip, IPAddress mask, IPAddress gateway)
{
	// debugf("WiFi CONNECTED");
	// Serial.println(ip.toString());
	Serial.println(WifiStation.getIP().toString());
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
		uint32_t any_channel_on=0;
		for (uint32_t i=0;i<PWM_CHANNELS;i++)
		{
			any_channel_on |= effect_target_values_[i];
		}
		if (any_channel_on > 0)
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
	Serial.begin(SERIAL_BAUD_RATE);
	Serial.systemDebugOutput(true); // Allow debug print to serial
	Serial.commandProcessing(true);
	// Mount file system, in order to work with files
	uint8_t slot = rboot_get_current_rom();
	debugf("\r\nrunning rom %d.\r\n", slot);
#ifndef DISABLE_SPIFFS
	if (slot == 0) {
#ifdef RBOOT_SPIFFS_0
		// debugf("trying to mount spiffs at 0x%08x, length %d", RBOOT_SPIFFS_0, SPIFF_SIZE);
		spiffs_mount_manual(RBOOT_SPIFFS_0, SPIFF_SIZE);
#else
		// debugf("trying to mount spiffs at 0x%08x, length %d", 0x100000, SPIFF_SIZE);
		spiffs_mount_manual(0x100000, SPIFF_SIZE);
#endif
	} else {
#ifdef RBOOT_SPIFFS_1
		// debugf("trying to mount spiffs at 0x%08x, length %d", RBOOT_SPIFFS_1, SPIFF_SIZE);
		spiffs_mount_manual(RBOOT_SPIFFS_1, SPIFF_SIZE);
#else
		// debugf("trying to mount spiffs at 0x%08x, length %d", 0x300000, SPIFF_SIZE);
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
