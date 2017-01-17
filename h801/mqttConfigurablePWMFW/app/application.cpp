#include <user_config.h>
#include <SmingCore/SmingCore.h>
#include <defaultlightconfig.h>
#include <SmingCore/Debug.h>
#include <pwmchannels.h>
#include "application.h"
#include "lightcontrol.h"
#include "telnet.h"
#include "mqtt.h"
#ifdef ENABLE_SSL
	#include <ssl/private_key.h>
	#include <ssl/cert.h>
#endif

NetConfigStorage NetConfig;
DefaultLightConfigStorage DefaultLightConfig;

///////////////////////////////////////
///// WIFI Stuff
///////////////////////////////////////

// Will be called when WiFi station was connected to AP
void wifiConnectOk()
{
	debugf("WiFi CONNECTED");
	//Serial.println(WifiStation.getIP().toString());
	startTelnetServer();
	startMqttClient();
	// Start publishing loop (also needed for mqtt reconnect)
	flashSingleChannel(1,CHAN_GREEN);
}

// Will be called when WiFi station timeout was reached
void wifiConnectFail()
{
	debugf("WiFi NOT CONNECTED!");

	flashSingleChannel(1,CHAN_RED);

	WifiStation.waitConnection(wifiConnectOk, 10, wifiConnectFail); // Repeat and check again
}

void connectToWifi()
{
	debugf("connecting 2 WiFi");
	WifiAccessPoint.enable(false);
	WifiStation.enable(true);
	WifiStation.enableDHCP(NetConfig.enabledhcp);
	WifiStation.setHostname(NetConfig.mqtt_clientid+".realraum.at");
	WifiStation.config(NetConfig.wifi_ssid, NetConfig.wifi_pass); // Put you SSID and Password here
	WifiStation.setIP(NetConfig.ip,NetConfig.netmask,NetConfig.gw);

	// Run our method when station was connected to AP (or not connected)
	WifiStation.waitConnection(wifiConnectOk, 30, wifiConnectFail); // We recommend 20+ seconds at start
}

//////////////////////////////////////
////// Base System Stuff  ////////////
//////////////////////////////////////


// Will be called when WiFi hardware and software initialization was finished
// And system initialization was completed
void ready()
{
	NetConfig.load(); //loads netsettings from fs
	//Serial.println(NetConfig.wifi_ssid);
	//Serial.println(NetConfig.wifi_pass);
	//Serial.println(NetConfig.ip.toString());
	instantinateMQTT();
	connectToWifi();
}

void init()
{
	//Serial.begin(115200);
	//Serial.systemDebugOutput(true); // Allow debug print to serial
	//Serial.commandProcessing(true);
	spiffs_mount(); // Mount file system, in order to work with files
	setupPWM(); //Init PWM with spiffs saved default settings

	telnetRegisterCmdsWithCommandHandler();
	commandHandler.registerSystemCommands();
	// Set system ready callback method
	System.onReady(ready);
}
