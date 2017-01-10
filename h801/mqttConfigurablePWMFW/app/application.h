#ifndef APPLICATION_H
#define APPLICATION_H
void setupPWM();
void flashMeNow();
void publishMessage();
void startMqttClient();
void wifiConnectOk();
void wifiConnectFail();
void ready();
void connectToWifi();
void startUDPServer();
#endif
