#include <SmingCore/SmingCore.h>
#include "pwmchannels.h"
#ifndef INCLUDE_DEFAULTCONFIG_H_
#define INCLUDE_DEFAULTCONFIG_H_

const String DEFAULTLIGHT_SETTINGS_FILE = "defaultlight.conf";

struct DefaultLightConfigStorage
{
	void load(uint32_t values[PWM_CHANNELS])
	{
		if (exist())
		{
			file_t f = fileOpen(DEFAULTLIGHT_SETTINGS_FILE, eFO_ReadOnly);
			fileRead(f, (void*) values, PWM_CHANNELS*sizeof(uint32_t));
			fileClose(f);
		}
	}

	void save(uint32_t values[PWM_CHANNELS])
	{
		file_t f = fileOpen(DEFAULTLIGHT_SETTINGS_FILE, eFO_WriteOnly | eFO_CreateNewAlways);
		fileWrite(f, (void*) values, PWM_CHANNELS*sizeof(uint32_t));
		fileClose(f);
	}

	bool exist() { return fileExist(DEFAULTLIGHT_SETTINGS_FILE)	&& fileGetSize(DEFAULTLIGHT_SETTINGS_FILE) >= PWM_CHANNELS*sizeof(uint32_t); }

};

extern DefaultLightConfigStorage DefaultLightConfig;


#endif /* INCLUDE_DEFAULTCONFIG_H_ */
