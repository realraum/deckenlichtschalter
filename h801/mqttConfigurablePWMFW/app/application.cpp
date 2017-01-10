#include <user_config.h>
#include <SmingCore/SmingCore.h>
#include <Services/ArduinoJson/ArduinoJson.h>
#include <defaultlightconfig.h>
#include <pwmchannels.h>
#include "application.h"

const uint32_t period = 5000; // * 200ns ^= 1 kHz

Timer flashTimer;
Timer procMQTTTimer;

////////////////////
//// PWM Stuff ////
///////////////////

//init PWM and restore stored pwm values
void setupPWM()
{
	// PWM setup
	uint32 io_info[PWM_CHANNELS][3] = {
		// MUX, FUNC, PIN
		{PERIPHS_IO_MUX_MTDI_U,  FUNC_GPIO12, 12},
		{PERIPHS_IO_MUX_MTDO_U,  FUNC_GPIO15, 15},
		{PERIPHS_IO_MUX_MTCK_U,  FUNC_GPIO13, 13},
		{PERIPHS_IO_MUX_MTMS_U,  FUNC_GPIO14, 14},
		{PERIPHS_IO_MUX_GPIO5_U, FUNC_GPIO5 ,  5},
	};
	// uint32 pwm_duty_initial[PWM_CHANNELS] = {0, 0, 0, 0, 0};
	uint32 pwm_duty_initial[PWM_CHANNELS] = {0, 0, 0, 0, 0};

	DefaultLightConfig.load(pwm_duty_initial); //load initial default values

	pwm_init(period, pwm_duty_initial, PWM_CHANNELS, io_info);
	pwm_start();
}

///////////////////////////////////////////
//// Userfeedback and Flash-LEDs Stuff ////
///////////////////////////////////////////


uint32_t flashme_num = 0;
uint8_t flashme_channel = 0;
uint32_t flashme_origvalue = 0;

void flashMeNow()
{
	if (flashme_num > 0)
	{
		if (flashme_num % 2 == 0)
		{
			pwm_set_duty(period, flashme_channel);
			pwm_start();				
		} else {
			pwm_set_duty(0, flashme_channel);
			pwm_start();				
		}
		flashme_num--;	
	} else {
		// stop timer
		flashTimer.stop();
		// restore values
		pwm_set_duty(flashme_origvalue, flashme_channel);
		pwm_start();		
	}
}

void flashChannel(uint8_t times, uint8_t channel)
{
	flashme_channel = channel;
	flashme_origvalue = pwm_get_duty(channel);
	flashme_num = times*2;
	pwm_set_duty(0,channel);
	pwm_start();
	flashTimer.initializeMs(500, flashMeNow).start(); // every 500ms
}


///////////////////////////////////////
///// WIFI Stuff
///////////////////////////////////////

void listNetworks(bool succeeded, BssList list)
{
	if (!succeeded)
	{
		Serial.println("Failed to scan networks");
		return;
	}

	for (int i = 0; i < list.count(); i++)
	{
		Serial.print("\tWiFi: ");
		Serial.print(list[i].ssid);
		Serial.print(", ");
		Serial.print(list[i].getAuthorizationMethodName());
		if (list[i].hidden) Serial.print(" (hidden)");
		Serial.println();
	}
}


// Will be called when WiFi station was connected to AP
void wifiConnectOk()
{
	debugf("WiFi CONNECTED");
	Serial.println(WifiStation.getIP().toString());
	startMqttClient();
	startUDPServer();
	// Start publishing loop (also needed for mqtt reconnect)
	procMQTTTimer.initializeMs(20 * 1000, publishMessage).start(); // every 20 seconds
	flashChannel(1,CHAN_GREEN);
}

// Will be called when WiFi station timeout was reached
void wifiConnectFail()
{
	debugf("WiFi NOT CONNECTED!");

	flashChannel(1,CHAN_RED);

	WifiStation.waitConnection(wifiConnectOk, 10, wifiConnectFail); // Repeat and check again
}

void connectToWifi()
{
	// Station - WiFi client
	WifiStation.enable(true);
	WifiStation.config(NetConfig.wifi_ssid, NetConfig.wifi_pass); // Put you SSID and Password here	
	WifiStation.setIP(NetConfig.ip,NetConfig.netmask,NetConfig.gw);


	// Print available access points
	WifiStation.startScan(listNetworks); // In Sming we can start network scan from init method without additional code
	// Run our method when station was connected to AP (or not connected)
	WifiStation.waitConnection(wifiConnectOk, 30, wifiConnectFail); // We recommend 20+ seconds at start
}

///////////////////////////////////////
/////// UDP Backup command interface
///////////////////////////////////////

void startUDPServer()
{
	//TODO
}

//////////////////////////////////
/////// MQTT Stuff ///////////////
//////////////////////////////////

// MQTT client
// For quick check you can use: http://www.hivemq.com/demos/websocket-client/ (Connection= test.mosquitto.org:8080)
MqttClient *mqtt;

// Check for MQTT Disconnection
void checkMQTTDisconnect(TcpClient& client, bool flag){
	
	// Called whenever MQTT connection is failed.
	if (flag == true)
		Serial.println("MQTT Broker Disconnected!!");
	else
		Serial.println("MQTT Broker Unreachable!!");
	
	flashChannel(2,CHAN_RED);

	// Restart connection attempt after few seconds
	// changes procMQTTTimer callback function
	procMQTTTimer.initializeMs(2 * 1000, startMqttClient).start(); // every 2 seconds
}

void onMessageDelivered(uint16_t msgId, int type) {
	Serial.printf("Message with id %d and QoS %d was delivered successfully.", msgId, (type==MQTT_MSG_PUBREC? 2: 1));
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
	if (root.containsKey(key) && root[key].is<unsigned int>())
	{
		a[pwm_channel] = ((uint32_t)root[key]) * period/1000;
		if (a[pwm_channel] > period)
			a[pwm_channel] = period;
	}
}

inline void setPWMDutyFromKey(JsonObject& root, String key, uint8_t pwm_channel)
{
	if (root.containsKey(key) && root[key].is<unsigned int>())
	{
		uint32_t newvalue = ((uint32_t)root[key])*period/1000;
		if (newvalue > period)
			newvalue = period;
		pwm_set_duty(newvalue, pwm_channel);
	}
}

// Callback for messages, arrived from MQTT server
void onMessageReceived(String topic, String message)
{
	Serial.print(topic);
	Serial.print(":\r\n\t"); // Pretify alignment for printing
	Serial.println(message);

	StaticJsonBuffer<200> jsonBuffer;

	JsonObject& root = jsonBuffer.parseObject(message);

	if (!root.success())
	{
	  Serial.println("JSON parseObject() failed");
	  return;
	}

	if (topic.endsWith("/light"))
	{
		setPWMDutyFromKey(root, JSONKEY_RED, CHAN_RED);
		setPWMDutyFromKey(root, JSONKEY_GREEN, CHAN_GREEN);
		setPWMDutyFromKey(root, JSONKEY_BLUE, CHAN_BLUE);
		setPWMDutyFromKey(root, JSONKEY_CW, CHAN_CW);
		setPWMDutyFromKey(root, JSONKEY_WW, CHAN_WW);
		pwm_start();
	} else if (topic.endsWith("/defaultlight"))
	{
		uint32_t pwm_duty_default[5] = {0,0,0,0,0};
		setArrayFromKey(root, pwm_duty_default, JSONKEY_RED, CHAN_RED);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_GREEN, CHAN_GREEN);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_BLUE, CHAN_BLUE);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_CW, CHAN_CW);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_WW, CHAN_WW);
		DefaultLightConfig.save(pwm_duty_default);
		flashChannel(1,CHAN_BLUE);
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
#ifdef ENABLE_SSL
	mqtt->addSslOptions(SSL_SERVER_VERIFY_LATER);

	#include <ssl/private_key.h>
	#include <ssl/cert.h>

	mqtt->setSslClientKeyCert(default_private_key, default_private_key_len,
							  default_certificate, default_certificate_len, NULL, true);

#endif
	// Assign a disconnect callback function
	mqtt->setCompleteDelegate(checkMQTTDisconnect);
	mqtt->subscribe("action/ceilingAll/light");
	mqtt->subscribe("action/ceilingAll/defaultlight");
	mqtt->subscribe(String("action/")+NetConfig.mqtt_clientid+"/light");
	mqtt->subscribe(String("action/")+NetConfig.mqtt_clientid+"/defaultlight");
}

//////////////////////////////////////
////// Base System Stuff  ////////////
//////////////////////////////////////


// Will be called when WiFi hardware and software initialization was finished
// And system initialization was completed
void ready()
{
	debugf("READY!");
	NetConfig.load();
	mqtt = new MqttClient(NetConfig.mqtt_broker, NetConfig.mqtt_port);
	connectToWifi();
}

void init()
{
	Serial.begin(SERIAL_BAUD_RATE);
	Serial.systemDebugOutput(true); // Allow debug print to serial

	setupPWM(); //also loads previously saved default settings

	// Set system ready callback method
	System.onReady(ready);
}
