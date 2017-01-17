#ifndef TELNET__H
#define TELNET__H

#include <SmingCore/Network/TelnetServer.h>

extern TelnetServer telnetServer;

void telnetCmdNetSettings(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdPrint(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdLight(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdSave(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdLs(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdCatFile(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdLoad(String commandLine  ,CommandOutput* commandOutput);
void telnetCmdReboot(String commandLine  ,CommandOutput* commandOutput);
void telnetAirUpdate(String commandLine  ,CommandOutput* commandOutput);
void startTelnetServer();
void telnetRegisterCmdsWithCommandHandler();


#endif