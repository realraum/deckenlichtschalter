
#include <SmingCore/SmingCore.h>
#include <defaultlightconfig.h>
#include <lightcontrol.h>
#include "mqtt.h"
#include "telnet.h"

///////////////////////////////////////
///// Telnet Backup command interface
///////////////////////////////////////

TelnetServer telnetServer;
HttpFirmwareUpdate ota_updater;
String ota_update_url_0, ota_update_url_9;
uint32_t auth_ip;
uint16_t auth_port=0;

void telnetCmdNetSettings(String commandLine  ,CommandOutput* commandOutput)
{
	Vector<String> commandToken;
	int numToken = splitString(commandLine, ' ' , commandToken);
	if (((uint32_t) telnetServer.getRemoteIp()) != auth_ip || telnetServer.getRemotePort() != auth_port)
	{
		commandOutput->println("Prevent Mistakes, give auth token");
		return;
	}
	if (numToken != 3)
	{
		commandOutput->println("Usage set ip|nm|gw|dhcp|wifissid|wifipass|mqttbroker|mqttport|mqttclientid|mqttuser|mqttpass <value>");
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
		uint32_t newport = commandToken[2].toInt();
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
	if (((uint32_t) telnetServer.getRemoteIp()) != auth_ip || telnetServer.getRemotePort() != auth_port)
	{
		commandOutput->println("Prevent Mistakes, give auth token");
		return;
	}
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
	if (((uint32_t) telnetServer.getRemoteIp()) != auth_ip || telnetServer.getRemotePort() != auth_port)
	{
		commandOutput->println("Prevent Mistakes, give auth token");
		return;
	}
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

void telnetAirUpdate(String commandLine  ,CommandOutput* commandOutput)
{
	Vector<String> commandToken;
	int numToken = splitString(commandLine, ' ' , commandToken);

	if (((uint32_t) telnetServer.getRemoteIp()) != auth_ip || telnetServer.getRemotePort() != auth_port)
	{
		commandOutput->println("Prevent Mistakes, give auth token");
		return;
	}
	if (2 != numToken)
	{
		commandOutput->println("Usage: update <url>|godoit");
		return;
	} else if (String("godoit") == commandToken[1] && ota_update_url_0.length() > 0 && ota_update_url_9.length() > 0)
	{
		stopMqttClient(); //disconnect MQTT
		stopAndRestoreValues(); // stop effects
		//disable lights
		for (uint8_t i=0;i<PWM_CHANNELS;i++)
			pwm_set_duty(0,i);
		pwm_start();
		ota_updater.addItem(0x0000, ota_update_url_0);
		ota_updater.addItem(0x9000, ota_update_url_9);
		commandOutput->println("OK, updating now");
		ota_updater.start();
	} else {
		ota_update_url_0 = commandToken[1] + "0x00000.bin";
		ota_update_url_9 = commandToken[1] + "0x09000.bin";
		commandOutput->println("Update URLs set, please check");
		commandOutput->println(ota_update_url_0);
		commandOutput->println(ota_update_url_9);
	}

}

void telnetAuth(String commandLine  ,CommandOutput* commandOutput)
{
	if (commandLine != "auth prevents mistakes "+NetConfig.authtoken)
		return;
	auth_ip = telnetServer.getRemoteIp();
	auth_port = telnetServer.getRemotePort();
	commandOutput->println("go ahead, but if you break it, you fix it");
}

void startTelnetServer()
{
	telnetServer.listen(TELNET_PORT_);
	telnetServer.enableCommand(true);
	//TODO: use encryption and client authentification
#ifdef ENABLE_SSL
	telnetServer.addSslOptions(SSL_SERVER_VERIFY_LATER);
	telnetServer.setSslClientKeyCert(default_private_key, default_private_key_len,
							  default_certificate, default_certificate_len, NULL, true);
	telnetServer.useSsl = true;
#endif
}

void telnetRegisterCmdsWithCommandHandler()
{
	commandHandler.registerCommand(CommandDelegate("set","Change network settings","configGroup", telnetCmdNetSettings));
	commandHandler.registerCommand(CommandDelegate("save","Save network settings","configGroup", telnetCmdSave));
	commandHandler.registerCommand(CommandDelegate("load","Save network settings","configGroup", telnetCmdLoad));
	commandHandler.registerCommand(CommandDelegate("show","Show network settings","configGroup", telnetCmdPrint));
	commandHandler.registerCommand(CommandDelegate("ls","List files","configGroup", telnetCmdLs));
	commandHandler.registerCommand(CommandDelegate("cat","List files","configGroup", telnetCmdCatFile));
	commandHandler.registerCommand(CommandDelegate("light","Test light","systemGroup", telnetCmdLight));
	commandHandler.registerCommand(CommandDelegate("restart","restart ESP8266","systemGroup", telnetCmdReboot));
	commandHandler.registerCommand(CommandDelegate("update","OTA Firmware update","systemGroup", telnetAirUpdate));
	commandHandler.registerCommand(CommandDelegate("auth","auth token","systemGroup", telnetAuth));
}
