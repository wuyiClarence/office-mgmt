#ifndef __MQTT_H__
#define __MQTT_H__
#include "MQTTAsync.h"
#include "config.h"

typedef struct MqttConfig_t
{
    char serverAddress[MAX_STR + 1];
    char clientId[MAX_STR + 1];
    char userName[MAX_STR + 1];
    char password[MAX_STR + 1];
    char keepaliveTopic[MAX_STR + 1];
    char shutdownHostTopic[MAX_STR + 1];
	char shutdownHostResTopic[MAX_STR + 1];
	char shutdownHostVirTopic[MAX_STR + 1];
	char wakeHostTopic[MAX_STR + 1];
	char wakeVirHostTopic[MAX_STR + 1];
    int  configChanged;
} MqttConfig;

void *mqtt_thread(void *arg);
extern MQTTAsync client;
extern volatile int connected;
extern MqttConfig *g_mqttConfig;
#endif