#ifndef __CONFIG_H__
#define __CONFIG_H__
#include <net/if.h>
#define MAX_STR 512
typedef struct AppConfig_t
{
    char hostname[MAX_STR + 1];
    char serverAddress[MAX_STR + 1];
    int  port;
    char userName[MAX_STR + 1];
    char password[MAX_STR + 1];
    char devInterface[IFNAMSIZ + 1];
    char macAddressStr[64];
    unsigned char macAddress[6];
    char uniqueid[MAX_STR + 1];
} AppConfig;

extern AppConfig g_appConfig;
int InitConfig();

#endif