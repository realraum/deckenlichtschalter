#ifndef MQTT__H
#define MQTT__H

void publishMessage();
inline void setArrayFromKey(JsonObject& root, uint32_t a[5], String key, uint8_t pwm_channel);
inline void setPWMDutyFromKey(JsonObject& root, String key, uint8_t pwm_channel);
void checkForwardInJsonAndSetCC(JsonObject& root, JsonObject& checkme);
void onMessageReceived(String topic, String message);
void startMqttClient();
void stopMqttClient();
void instantinateMQTT();

#endif