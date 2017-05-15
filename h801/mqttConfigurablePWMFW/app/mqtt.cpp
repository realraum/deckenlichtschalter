#include <SmingCore/SmingCore.h>
#include <defaultlightconfig.h>
#include "pwmchannels.h"
#include "lightcontrol.h"
#include "mqtt.h"

//////////////////////////////////
/////// MQTT Stuff ///////////////
//////////////////////////////////


Timer procMQTTTimer;
MqttClient *mqtt = 0;


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

inline void setArrayFromKey(JsonObject& root, uint32_t a[PWM_CHANNELS], String key, uint8_t pwm_channel)
{
	//BUG: can't check type because .is is buggy and does not compile in Sming 3.0.1.
	//if (root.containsKey(key) && root[key].is<unsigned int>())
	if (root.containsKey(key))
	{
		uint32_t value = (uint32_t)root[key];
		if (value > NetConfig.chan_range[pwm_channel])
		{
			return;
		}
		a[pwm_channel] = value * pwm_period / NetConfig.chan_range[pwm_channel];
		if (a[pwm_channel] > pwm_period)
			a[pwm_channel] = pwm_period;
	} else
	{
		//this allows us to set the default to current light values by sending an empty json object to /defaultlight
		a[pwm_channel] = effect_target_values_[pwm_channel];
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

//if present and an array of strings, will forward entire json to next element in list
//minus that element in the list
void checkForwardInJsonAndSetCC(JsonObject& root, JsonObject& checkme)
{
	if (checkme.containsKey(JSONKEY_FORWARD))
	{
		JsonArray& cc_list = checkme[JSONKEY_FORWARD];
		if (cc_list.success() && cc_list.size() > 0)
		{
			mqtt_forward_to_ = cc_list.get<String>(0);
			cc_list.removeAt(0);
			root.printTo(mqtt_payload_);
		}
	}
}

void simulateCWwithRGB(uint32_t a[PWM_CHANNELS])
{
	if (NetConfig.simulatecw_w_rgb)
	{
		a[CHAN_RED] = max(a[CHAN_RED],a[CHAN_CW]);
		a[CHAN_GREEN] = max(a[CHAN_GREEN],a[CHAN_CW]);
		a[CHAN_BLUE] = max(a[CHAN_BLUE],a[CHAN_CW]);
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

	StaticJsonBuffer<1024> jsonBuffer;

	// pleaserepeat does not care about message content, thus it is checked before using the jsonBuffer for parsing
	// (as JsonBuffer should not be reused) This allows us to use that buffer for sending a message ourselves
	if (topic.endsWith(JSON_TOPIC3_PLEASEREPEAT))
	{
		JsonObject& root = jsonBuffer.createObject();
		root[JSONKEY_RED] = effect_target_values_[CHAN_RED] * NetConfig.chan_range[CHAN_RED] / pwm_period;
		root[JSONKEY_GREEN] = effect_target_values_[CHAN_GREEN] * NetConfig.chan_range[CHAN_GREEN] / pwm_period;
		root[JSONKEY_BLUE] = effect_target_values_[CHAN_BLUE] * NetConfig.chan_range[CHAN_BLUE] / pwm_period;
		root[JSONKEY_CW] = effect_target_values_[CHAN_CW] * NetConfig.chan_range[CHAN_CW] / pwm_period;
		root[JSONKEY_WW] = effect_target_values_[CHAN_WW] * NetConfig.chan_range[CHAN_WW] / pwm_period;
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
		simulateCWwithRGB(effect_target_values_);

		//-----
		if (root.containsKey(JSONKEY_FLASH))
		{
			JsonObject& effectobj = root[JSONKEY_FLASH];
			uint32_t repetitions = DEFAULT_EFFECT_REPETITIONS;
			if (effectobj.containsKey(JSONKEY_REPETITIONS))
				repetitions = effectobj[JSONKEY_REPETITIONS];
			uint32_t period = DEFAULT_FLASH_PERIOD_;
			if (effectobj.containsKey(JSONKEY_PERIOD))
				period = effectobj[JSONKEY_PERIOD];
			checkForwardInJsonAndSetCC(root, effectobj);
			startFlash(repetitions, FLASH_INTERMED_ORIG, period);
		}
		else if (root.containsKey(JSONKEY_FADE))
		{
			JsonObject& effectobj = root[JSONKEY_FADE];
			uint32_t duration = DEFAULT_EFFECT_DURATION;
			if (effectobj.containsKey(JSONKEY_DURATION))
				duration = effectobj[JSONKEY_DURATION];
			checkForwardInJsonAndSetCC(root, effectobj);
			startFade(duration);
		} else
		{
			//apply Values right now
			stopAndRestoreValues(true); //disable any Effects
			applyValues(effect_target_values_); //apply light
		}
	} else if (topic.endsWith(JSON_TOPIC3_DEFAULTLIGHT))
	{
		uint32_t pwm_duty_default[PWM_CHANNELS] = {0,0,0,0,0};
		setArrayFromKey(root, pwm_duty_default, JSONKEY_RED, CHAN_RED);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_GREEN, CHAN_GREEN);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_BLUE, CHAN_BLUE);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_CW, CHAN_CW);
		setArrayFromKey(root, pwm_duty_default, JSONKEY_WW, CHAN_WW);
		simulateCWwithRGB(pwm_duty_default);
		DefaultLightConfig.save(pwm_duty_default);
		flashSingleChannel(1,CHAN_BLUE);
	}
}

// Run MQTT client, connect to server, subscribe topics
void startMqttClient()
{
	procMQTTTimer.stop();

	if (0 == mqtt)
		mqtt = new MqttClient(NetConfig.mqtt_broker, NetConfig.mqtt_port, onMessageReceived);

/*	if(!mqtt->setWill("last/will","The connection from this device is lost:(", 1, true)) {
		debugf("Unable to set the last will and testament. Most probably there is not enough memory on the device.");
	}
*/
	mqtt->setKeepAlive(42);
	mqtt->setPingRepeatTime(21);
	bool usessl=false;
#ifdef ENABLE_SSL
	usessl=true;
	mqtt->addSslOptions(SSL_SERVER_VERIFY_LATER);

	mqtt->setSslClientKeyCert(default_private_key, default_private_key_len,
							  default_certificate, default_certificate_len, NULL, true);
#endif

	// Assign a disconnect callback function
	mqtt->setCompleteDelegate(checkMQTTDisconnect);
	debugf("connecting to to MQTT broker");
	mqtt->connect(NetConfig.mqtt_clientid, NetConfig.mqtt_user, NetConfig.mqtt_pass, true);
	debugf("connected to MQTT broker");
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT,true));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_DEFAULTLIGHT,true));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_PLEASEREPEAT,true));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT,false));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_DEFAULTLIGHT,false));
	mqtt->subscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_PLEASEREPEAT,false));

	procMQTTTimer.initializeMs(20 * 1000, publishMessage).start(); // every 20 seconds
}

void stopMqttClient()
{
	mqtt->unsubscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT,true));
	mqtt->unsubscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_DEFAULTLIGHT,true));
	mqtt->unsubscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_PLEASEREPEAT,true));
	mqtt->unsubscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_LIGHT,false));
	mqtt->unsubscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_DEFAULTLIGHT,false));
	mqtt->unsubscribe(NetConfig.getMQTTTopic(JSON_TOPIC3_PLEASEREPEAT,false));
	mqtt->setKeepAlive(0);
	mqtt->setPingRepeatTime(0);
	procMQTTTimer.stop();
}

