#ifndef APPLICATION_H
#define APPLICATION_H

MqttClient *mqtt = 0;
TelnetServer telnetServer;

void setupPWM();
void flashMeNow();
void wifiConnectOk();
void wifiConnectFail();
void connectToWifi();
void startTelnetServer();
void publishMessage();
void startMqttClient();
void ready();
void init();
#endif