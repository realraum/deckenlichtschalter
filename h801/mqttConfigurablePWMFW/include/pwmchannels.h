#ifndef USERCONFIG_H
#define USERCONFIG_H

#define PWM_CHANNELS 5

#define CHAN_RED 0    //R- on GPIO15  (WemosD1: D8)
#define CHAN_GREEN 1  //G- on GPIO13  (WemosD1: D7)
#define CHAN_BLUE 2   //B- on GPIO12  (WemosD1: D6)
#define CHAN_CW 3     //W1 on GPIO14  (WemosD1: D5)
#define CHAN_WW 4     //W2 on GPIO4   (WemosD1: D2)

#define FAN_GPIO 16

#define JSONKEY_RED "r"
#define JSONKEY_GREEN "g"
#define JSONKEY_BLUE "b"
#define JSONKEY_CW "cw"
#define JSONKEY_WW "ww"

#define JSONKEY_FLASH "flash"
#define JSONKEY_FADE "fade"
#define JSONKEY_REPETITIONS "repetitions"
#define JSONKEY_DURATION "duration"
#define JSONKEY_PERIOD "period"
#define JSONKEY_FORWARD "cc"

const String JSON_TOPIC1 = "action/";
const String JSON_TOPIC2_ALL ="ceilingAll";
const String JSON_TOPIC3_LIGHT = "/light";
const String JSON_TOPIC3_DEFAULTLIGHT = "/defaultlight";
const String JSON_TOPIC3_PLEASEREPEAT = "/pleaserepeat";

const uint32_t TELNET_PORT_ = 2323;

#endif // USERCONFIG_H