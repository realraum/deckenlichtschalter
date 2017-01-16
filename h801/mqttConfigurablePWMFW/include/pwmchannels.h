#ifndef USERCONFIG_H
#define USERCONFIG_H

#define PWM_CHANNELS 5

#define CHAN_RED 0    //R- on GPIO15
#define CHAN_GREEN 1  //G- on GPIO13
#define CHAN_BLUE 2   //B- on GPIO12
#define CHAN_CW 3     //W1 on GPIO14
#define CHAN_WW 4     //W2 on GPIO4

#define JSONKEY_RED "r"
#define JSONKEY_GREEN "g"
#define JSONKEY_BLUE "b"
#define JSONKEY_CW "cw"
#define JSONKEY_WW "ww"

#define JSONKEY_FLASH "flash"
#define JSONKEY_FADE "fade"
#define JSONKEY_REPETITIONS "repetitions"
#define JSONKEY_DURATION "duration"

const String JSON_TOPIC1 = "action/";
const String JSON_TOPIC2_ALL ="ceilingAll";
const String JSON_TOPIC3_LIGHT = "/light";
const String JSON_TOPIC3_DEFAULTLIGHT = "/defaultlight";
const String JSON_TOPIC3_PLEASEREPEAT = "/pleaserepeat";


#endif // USERCONFIG_H