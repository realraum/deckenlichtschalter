#ifndef APPLICATION_H
#define APPLICATION_H

void setupPWM();
void flashMeNow();
void wifiConnectOk();
void wifiConnectFail();
void connectToWifi();
void startTelnetServer();
void publishMessage();
void startMqttClient();
void stopMqttClient();
void ready();
void init();
#endif