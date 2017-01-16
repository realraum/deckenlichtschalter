#include <user_config.h>
#include <SmingCore/SmingCore.h>
#include <SmingCore/Network/TelnetServer.h>
#include <defaultlightconfig.h>
#include <SmingCore/Debug.h>
#include <pwmchannels.h>
#include "application.h"
#include "lightcontrol.h"
#ifdef ENABLE_SSL
	#include <ssl/private_key.h>
	#include <ssl/cert.h>
#endif


Timer procMQTTTimer;

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
	procMQTTTimer.initializeMs(20 * 1000, publishMessage).start(); // every 20 seconds
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

///////////////////////////////////////
///// Telnet Backup command interface
///////////////////////////////////////

void telnetCmdNetSettings(String commandLine  ,CommandOutput* commandOutput)
{
	Vector<String> commandToken;
	int numToken = splitString(commandLine, ' ' , commandToken);
	if (numToken != 3)
	{
		commandOutput->printf("Usage set ip|nm|gw|dhcp|wifissid|wifipass|mqttbroker|mqttport|mqttclientid|mqttuser|mqttpass <value>\r\n");
	}
	else if (commandToken[1] == "ip")
	{
		IPAddress newip(commandToken[2]);
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),newip.toString().c_str());
		if (!newip.isNull())
			NetConfig.ip = newip;
	}
	else if (commandToken[1] == "nm")
	{
		IPAddress newip(commandToken[2]);
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),newip.toString().c_str());
		if (!newip.isNull())
			NetConfig.netmask = newip;
	}
	else if (commandToken[1] == "gw")
	{
		IPAddress newip(commandToken[2]);
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),newip.toString().c_str());
		if (!newip.isNull())
			NetConfig.gw = newip;
	}
	else if (commandToken[1] == "wifissid")
	{
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),commandToken[2].c_str());
		NetConfig.wifi_ssid = commandToken[2];
	}
	else if (commandToken[1] == "wifipass")
	{
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),commandToken[2].c_str());
		NetConfig.wifi_pass = commandToken[2];
	}
	else if (commandToken[1] == "mqttbroker")
	{
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),commandToken[2].c_str());
		NetConfig.mqtt_broker = commandToken[2];
	}
	else if (commandToken[1] == "mqttport")
	{
		uint32_t newport = atoi(commandToken[2].c_str());
		commandOutput->printf("%s: '%d'\r\n",commandToken[1].c_str(),newport);
		if (newport > 0 && newport < 65536)
			NetConfig.mqtt_port = newport;
	}
	else if (commandToken[1] == "mqttclientid")
	{
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),commandToken[2].c_str());
		NetConfig.mqtt_clientid = commandToken[2];
	}
	else if (commandToken[1] == "mqttuser")
	{
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),commandToken[2].c_str());
		NetConfig.mqtt_user = commandToken[2];
	}
	else if (commandToken[1] == "mqttpass")
	{
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),commandToken[2].c_str());
		NetConfig.mqtt_pass = commandToken[2];
	}
	else if (commandToken[1] == "dhcp")
	{
		NetConfig.enabledhcp = commandToken[2] == "1" || commandToken[2] == "true" || commandToken[2] == "yes" || commandToken[2] == "on";
		commandOutput->printf("%s: '%s'\r\n",commandToken[1].c_str(),(NetConfig.enabledhcp)?"on":"off");
	} else {
		commandOutput->printf("Invalid subcommand. Try %s list\r\n", commandToken[0].c_str());
	}
}

void telnetCmdPrint(String commandLine  ,CommandOutput* commandOutput)
{
	commandOutput->println("Dumping Configuration");
	commandOutput->println("WiFi SSID: " + NetConfig.wifi_ssid + " actualSSID: "+WifiStation.getSSID());
	commandOutput->println("WiFi Pass: " + NetConfig.wifi_pass + " actualPASS: "+WifiStation.getPassword());
	commandOutput->println("Hostname: " + WifiStation.getHostname());
	commandOutput->println("MAC: " + WifiStation.getMAC());
	commandOutput->println("IP: " + NetConfig.ip.toString() + " actualIP: "+WifiStation.getIP().toString());
	commandOutput->println("NM: " + NetConfig.netmask.toString()+ " actualGW: "+WifiStation.getNetworkMask().toString());
	commandOutput->println("GW: " + NetConfig.gw.toString()+ " actualGW: "+WifiStation.getNetworkGateway().toString());
	commandOutput->println((NetConfig.enabledhcp)?"DHCP: on":"DHCP: off");
	commandOutput->println((WifiStation.isEnabledDHCP())?"actual DHCP: on":"DHCP: off");
	commandOutput->println("MQTT Broker: " + NetConfig.mqtt_broker + ":" + String(NetConfig.mqtt_port));
	commandOutput->println("MQTT ClientID: " + NetConfig.mqtt_clientid);
	commandOutput->println("MQTT Login: " + NetConfig.mqtt_user +"/"+ NetConfig.mqtt_pass);
}

void telnetCmdLight(String commandLine  ,CommandOutput* commandOutput)
{
	Vector<String> commandToken;
	int numToken = splitString(commandLine, ' ' , commandToken);
	if (numToken != 2)
	{
		commandOutput->println("Usage light on|off|info|half|flash0,flash1,flash2,flash3");
	}
	else if (commandToken[1] == "on")
	{
		uint32_t deflightconf[PWM_CHANNELS]={0,0,0,0,0};
		DefaultLightConfig.load(deflightconf);
		for (uint8_t i=0;i<PWM_CHANNELS;i++)
			pwm_set_duty(deflightconf[i],i);
		pwm_start();
	}
	else if (commandToken[1] == "off")
	{
		for (uint8_t i=0;i<PWM_CHANNELS;i++)
			pwm_set_duty(0,i);
		pwm_start();
	}
	else if (commandToken[1] == "half")
	{
		for (uint8_t i=0;i<PWM_CHANNELS;i++)
			pwm_set_duty(pwm_period/2,i);
		pwm_start();
	}
	else if (commandToken[1] == "info")
	{
		uint32_t deflightconf[PWM_CHANNELS]={0,0,0,0,0};
		DefaultLightConfig.load(deflightconf);
		commandOutput->println("Current: r:"+String(pwm_get_duty(CHAN_RED))+" g:"+String(pwm_get_duty(CHAN_GREEN))+" b:"+String(pwm_get_duty(CHAN_BLUE))+" cw:"+String(pwm_get_duty(CHAN_CW))+" ww:"+String(pwm_get_duty(CHAN_WW)));
		commandOutput->println("Default: r:"+String(deflightconf[CHAN_RED])+" g:"+String(deflightconf[CHAN_GREEN])+" b:"+String(deflightconf[CHAN_BLUE])+" cw:"+String(deflightconf[CHAN_CW])+" ww:"+String(deflightconf[CHAN_WW]));
		commandOutput->println("effect_target_values_: r:"+String(effect_target_values_[CHAN_RED])+" g:"+String(effect_target_values_[CHAN_GREEN])+" b:"+String(effect_target_values_[CHAN_BLUE])+" cw:"+String(effect_target_values_[CHAN_CW])+" ww:"+String(effect_target_values_[CHAN_WW]));
		commandOutput->println("effect_intermid_values_: r:"+String(effect_intermid_values_[CHAN_RED])+" g:"+String(effect_intermid_values_[CHAN_GREEN])+" b:"+String(effect_intermid_values_[CHAN_BLUE])+" cw:"+String(effect_intermid_values_[CHAN_CW])+" ww:"+String(effect_intermid_values_[CHAN_WW]));
	}
	else if (commandToken[1] == "flash0")
	{
		flashSingleChannel(3,0);
	}
	else if (commandToken[1] == "flash1")
	{
		effect_target_values_[0]=pwm_period/3;
		effect_target_values_[1]=0;
		effect_target_values_[2]=pwm_period/3;
		effect_target_values_[3]=0;
		effect_target_values_[4]=10;
		effect_intermid_values_[0]=0;
		effect_intermid_values_[1]=pwm_get_duty(1);
		effect_intermid_values_[2]=0;
		effect_intermid_values_[3]=pwm_get_duty(3);
		effect_intermid_values_[4]=pwm_get_duty(4);
		startFlash(2,FLASH_INTERMED_USERSET);
	}
	else if (commandToken[1] == "flash2")
	{
		flashSingleChannel(3,2);
	}
	else if (commandToken[1] == "flash3")
	{
		effect_target_values_[0]=pwm_period/3;
		effect_target_values_[1]=0;
		effect_target_values_[2]=pwm_period/3;
		effect_target_values_[3]=0;
		effect_target_values_[4]=10;
		startFlash(2,FLASH_INTERMED_DARK);
	}

}

void telnetCmdSave(String commandLine  ,CommandOutput* commandOutput)
{
	commandOutput->println("OK, saving values...");
	NetConfig.save();
}

void telnetCmdLs(String commandLine  ,CommandOutput* commandOutput)
{
	Vector<String> list = fileList();
	for (int i = 0; i < list.count(); i++)
		commandOutput->println(String(fileGetSize(list[i])) + " " + list[i]);
}

void telnetCmdCatFile(String commandLine  ,CommandOutput* commandOutput)
{
	Vector<String> commandToken;
	int numToken = splitString(commandLine, ' ' , commandToken);

	if (numToken != 2)
	{
		commandOutput->println("Usage: cat <file>");
		return;
	}
	if (fileExist(commandToken[1]))
	{
		commandOutput->println("Contents of "+commandToken[1]);
		commandOutput->println(fileGetContent(commandToken[1]));
	} else {
		commandOutput->println("File '"+commandToken[1]+"' does not exist");
	}
}

void telnetCmdLoad(String commandLine  ,CommandOutput* commandOutput)
{
	commandOutput->printf("OK, reloading values...\r\n");
	NetConfig.load();
}

void telnetCmdReboot(String commandLine  ,CommandOutput* commandOutput)
{
	commandOutput->printf("OK, restarting...\r\n");
	telnetServer.flush();
	telnetServer.close();
	System.restart();
}

void startTelnetServer()
{
	telnetServer.listen(2323);
	telnetServer.enableCommand(true);
	//TODO: use encryption and client authentification
#ifdef ENABLE_SSL
	telnetServer.addSslOptions(SSL_SERVER_VERIFY_LATER);
	telnetServer.setSslClientKeyCert(default_private_key, default_private_key_len,
							  default_certificate, default_certificate_len, NULL, true);
	telnetServer.useSsl = true;
#endif
}

//////////////////////////////////
/////// MQTT Stuff ///////////////
//////////////////////////////////


// Check for MQTT Disconnection
void checkMQTTDisconnect(TcpClient& client, bool flag){

	// Called whenever MQTT connection is failed.
	if (flag == true)
	{
		//Serial.println("MQTT Broker Disconnected!!");
		flashSingleChannel(2,CHAN_RED);
	}
	else
	{
		//Serial.println("MQTT Broker Unreachable!!");
		flashSingleChannel(3,CHAN_RED);
	}

	// Restart connection attempt after few seconds
	// changes procMQTTTimer callback function
	procMQTTTimer.initializeMs(2 * 1000, startMqttClient).start(); // every 2 seconds
}

void onMessageDelivered(uint16_t msgId, int type) {
	//Serial.printf("Message with id %d and QoS %d was delivered successfully.", msgId, (type==MQTT_MSG_PUBREC? 2: 1));
}

// Publish our message
void publishMessage()
{
	if (mqtt->getConnectionState() != eTCS_Connected)
		startMqttClient(); // Auto reconnect

	//Serial.println("Let's publish message now!");
	//mqtt->publish("main/frameworks/sming", "Hello friends, from Internet of things :)");

	//mqtt->publishWithQoS("important/frameworks/sming", "Request Return Delivery", 1, false, onMessageDelivered); // or publishWithQoS
}

inline void setArrayFromKey(JsonObject& root, uint32_t a[5], String key, uint8_t pwm_channel)
{
	//BUG: can't check type because .is is buggy and does not compile in Sming 3.0.1.
	//if (root.containsKey(key) && root[key].is<unsigned int>())
	if (root.containsKey(key))
	{
		uint32_t value = (uint32_t)root[key];
		if (value > 1000)
		{
			return;
		}
		a[pwm_channel] = value * pwm_period / 1000;
		if (a[pwm_channel] > pwm_period)
			a[pwm_channel] = pwm_period;
	} else
	{
		a[pwm_channel] = pwm_get_duty(pwm_channel);
	}
}

inline void setPWMDutyFromKey(JsonObject& root, String key, uint8_t pwm_channel)
{
	//BUG: can't check type because .is is buggy and does not compile in Sming 3.0.1
	// if (root.containsKey(key) && root[key].is<unsigned int>())
	if (root.containsKey(key))
	{
		uint32_t value = (uint32_t)root[key];
		if (value > 1000)
		{
			return;
		}
		value = value * pwm_period / 1000;
		if (value > pwm_period)
			value = pwm_period;
		pwm_set_duty(value, pwm_channel);
	}
}

// Callback for messages, arrived from MQTT server
void onMessageReceived(String topic, String message)
{
	debugf("topic: %s",topic.c_str());
	debugf("msg: %s",message.c_str());
	//GRML BUG :-( It would be really nice to filter out retained messages,
	//             to avoid the light powering up, going into defaultlight settings, then getting wifi and switching to a retained /light setting
	//GRML :-( unfortunately we can't distinguish between retained and fresh messages here

	DynamicJsonBuffer jsonBuffer;

	// pleaserepeat does not care about message content, thus it is checked before using the jsonBuffer for parsing
	// (as JsonBuffer should not be reused) This allows us to use that buffer for sending a message ourselves
	if (topic.endsWith(JSON_TOPIC3_PLEASEREPEAT))
	{
		JsonObject& root = jsonBuffer.createObject();
		root[JSONKEY_RED] = effect_target_values_[CHAN_RED];
		root[JSONKEY_GREEN] = effect_target_values_[CHAN_GREEN];
		root[JSONKEY_BLUE] = effect_target_values_[CHAN_BLUE];
		root[JSONKEY_CW] = effect_target_values_[CHAN_CW];
		root[JSONKEY_WW] = effect_target_values_[CHAN_WW];
		root.printTo(message);
		//publish to myself (where presumably everybody else also listens), the current settings
		mqtt->publish(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT), message, false);
		return; //return so we don't reuse the now used jsonBuffer
	}

	//use jsonBuffer to parse message
	JsonObject& root = jsonBuffer.parseObject(message);

	if (!root.success())
	{
	  //Serial.println("JSON parseObject() failed");
	  return;
	}

	if (topic.endsWith(JSON_TOPIC3_LIGHT))
	{
		setArrayFromKey(root, effect_target_values_, JSONKEY_RED, CHAN_RED);
		setArrayFromKey(root, effect_target_values_, JSONKEY_GREEN, CHAN_GREEN);
		setArrayFromKey(root, effect_target_values_, JSONKEY_BLUE, CHAN_BLUE);
		setArrayFromKey(root, effect_target_values_, JSONKEY_CW, CHAN_CW);
		setArrayFromKey(root, effect_target_values_, JSONKEY_WW, CHAN_WW);

		if (root.containsKey(JSONKEY_FLASH))
		{
			JsonObject& effectobj = root[JSONKEY_FLASH];
			uint32_t repetitions = DEFAULT_EFFECT_REPETITIONS;
			if (effectobj.containsKey(JSONKEY_REPETITIONS))
				repetitions = effectobj[JSONKEY_REPETITIONS];
			startFlash(repetitions, FLASH_INTERMED_ORIG);
		}
		else if (root.containsKey(JSONKEY_FADE))
		{
			JsonObject& effectobj = root[JSONKEY_FADE];
			uint32_t duration = DEFAULT_EFFECT_REPETITIONS;
			if (effectobj.containsKey(JSONKEY_DURATION))
				duration = effectobj[JSONKEY_DURATION];
			startFade(duration);
		} else
		{
			//apply Values right now
			applyValues(effect_target_values_);
		}
	} else if (topic.endsWith(JSON_TOPIC3_DEFAULTLIGHT))
	{
		uint32_t pwm_duty_default[PWM_CHANNELS] = {0,0,0,0,0};
		setArrayFromKey(root, pwm_duty_default, JSONKEY_RED, CHAN_RED);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_GREEN, CHAN_GREEN);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_BLUE, CHAN_BLUE);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_CW, CHAN_CW);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_WW, CHAN_WW);
		DefaultLightConfig.save(pwm_duty_default);
		flashSingleChannel(1,CHAN_BLUE);
	}
}

// Run MQTT client, connect to server, subscribe topics
void startMqttClient()
{
	procMQTTTimer.stop();
/*	if(!mqtt->setWill("last/will","The connection from this device is lost:(", 1, true)) {
		debugf("Unable to set the last will and testament. Most probably there is not enough memory on the device.");
	}
*/
	mqtt->connect(NetConfig.mqtt_clientid, NetConfig.mqtt_user, NetConfig.mqtt_pass, true);
	mqtt->setKeepAlive(42);
	mqtt->setPingRepeatTime(21);
#ifdef ENABLE_SSL
	mqtt->addSslOptions(SSL_SERVER_VERIFY_LATER);

	mqtt->setSslClientKeyCert(default_private_key, default_private_key_len,
							  default_certificate, default_certificate_len, NULL, true);

#endif
	// Assign a disconnect callback function
	mqtt->setCompleteDelegate(checkMQTTDisconnect);
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT,true));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_DEFAULTLIGHT,true));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_PLEASEREPEAT,true));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT,false));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_DEFAULTLIGHT,false));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_PLEASEREPEAT,false));
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
	mqtt = new MqttClient(NetConfig.mqtt_broker, NetConfig.mqtt_port, onMessageReceived);
	connectToWifi();
}

void init()
{
	//Serial.begin(115200);
	//Serial.systemDebugOutput(true); // Allow debug print to serial
	//Serial.commandProcessing(true);
	spiffs_mount(); // Mount file system, in order to work with files
	setupPWM(); //Init PWM with spiffs saved default settings

	commandHandler.registerCommand(CommandDelegate("set","Change network settings","configGroup", telnetCmdNetSettings));
	commandHandler.registerCommand(CommandDelegate("save","Save network settings","configGroup", telnetCmdSave));
	commandHandler.registerCommand(CommandDelegate("load","Save network settings","configGroup", telnetCmdSave));
	commandHandler.registerCommand(CommandDelegate("show","Show network settings","configGroup", telnetCmdPrint));
	commandHandler.registerCommand(CommandDelegate("ls","List files","configGroup", telnetCmdLs));
	commandHandler.registerCommand(CommandDelegate("cat","List files","configGroup", telnetCmdCatFile));
	commandHandler.registerCommand(CommandDelegate("light","Test light","systemGroup", telnetCmdLight));
	commandHandler.registerCommand(CommandDelegate("restart","restart ESP8266","systemGroup", telnetCmdReboot));
	commandHandler.registerSystemCommands();
	// Set system ready callback method
	System.onReady(ready);
}
