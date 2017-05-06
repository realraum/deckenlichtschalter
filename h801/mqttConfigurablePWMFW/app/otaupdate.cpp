#include <SmingCore/SmingCore.h>
#include <defaultlightconfig.h>
#include <lightcontrol.h>
#include "mqtt.h"
#include "otaupdate.h"

rBootHttpUpdate* otaUpdater = 0;

void OtaUpdate_prepareSystem()
{
	stopMqttClient(); //disconnect MQTT
	stopAndRestoreValues(); // stop effects
	// set light to same state they will be in
	// once GPIOs switch to INPUT
	// so there won't be a sudden power drop during flash when all LEDs switch on
	uint32_t light_during_flash[PWM_CHANNELS] = {0,pwm_period,pwm_period,pwm_period,0};
	applyValues(light_during_flash);
	enableFan(true);
	//start firmware update
}

void OtaUpdate_CallBack(rBootHttpUpdate& client, bool result) {

	//Serial.println("In callback...");
	if(result == true) {
		// success
		uint8 slot;
		slot = rboot_get_current_rom();
		if (slot == 0) slot = 1; else slot = 0;
		// set to boot new rom and then reboot
		// Serial.printf("Firmware updated, rebooting to rom %d...\r\n", slot);
		rboot_set_current_rom(slot);
	} else {
		// fail
		// Serial.println("Firmware update failed!");
	}
	System.restart();
}

void OtaUpdate(String rom0url, String rom1url, String spiffsurl) {

	uint8 slot;
	rboot_config bootconf;

	// need a clean object, otherwise if run before and failed will not run again
	if (otaUpdater) delete otaUpdater;
	otaUpdater = new rBootHttpUpdate();

	//prepare System
	OtaUpdate_prepareSystem();

	// select rom slot to flash
	bootconf = rboot_get_config();
	slot = bootconf.current_rom;
	if (slot == 0) slot = 1; else slot = 0;

#ifndef RBOOT_TWO_ROMS
	// flash rom to position indicated in the rBoot config rom table
	otaUpdater->addItem(bootconf.roms[slot], rom0url);
#else
	// flash appropriate rom
	if (slot == 0) {
		otaUpdater->addItem(bootconf.roms[slot], rom0url);
	} else {
		otaUpdater->addItem(bootconf.roms[slot], rom1url);
	}
#endif

#ifndef DISABLE_SPIFFS
	// use user supplied values (defaults for 4mb flash in makefile)
	if (slot == 0) {
		otaUpdater->addItem(RBOOT_SPIFFS_0, spiffsurl);
	} else {
		otaUpdater->addItem(RBOOT_SPIFFS_1, spiffsurl);
	}
#endif

	// request switch and reboot on success
	//otaUpdater->switchToRom(slot);
	// and/or set a callback (called on failure or success without switching requested)
	otaUpdater->setCallback(OtaUpdate_CallBack);

	// start update
	otaUpdater->start();
}