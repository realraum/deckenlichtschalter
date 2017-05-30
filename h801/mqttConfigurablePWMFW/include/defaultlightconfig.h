#include <SmingCore/SmingCore.h>
#include "pwmchannels.h"
#ifndef INCLUDE_DEFAULTCONFIG_H_
#define INCLUDE_DEFAULTCONFIG_H_

class DefaultLightConfigStorage
{
public:
	DefaultLightConfigStorage(String fn) : filename(fn) {};

	void load(uint32_t values[PWM_CHANNELS])
	{
		if (exist())
		{
			file_t f = fileOpen(filename, eFO_ReadOnly);
			fileRead(f, (void*) values, PWM_CHANNELS*sizeof(uint32_t));
			fileClose(f);
		}
	}

	void save(uint32_t values[PWM_CHANNELS])
	{
		file_t f = fileOpen(filename, eFO_WriteOnly | eFO_CreateNewAlways);
		fileWrite(f, (void*) values, PWM_CHANNELS*sizeof(uint32_t));
		fileClose(f);
	}

	bool exist() { return fileExist(filename)	&& fileGetSize(filename) >= PWM_CHANNELS*sizeof(uint32_t); }
private:
	String filename;
};

extern DefaultLightConfigStorage DefaultLightConfig;
extern DefaultLightConfigStorage ButtonLightConfig;


#endif /* INCLUDE_DEFAULTCONFIG_H_ */
