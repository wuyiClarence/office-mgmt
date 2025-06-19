#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mqtt.h"
#include "zwx_log.h"
#include "cJSON.h"
#include "wakeonlan.h"
#include "virhost.h"

#if !defined(_WIN32)
#include <unistd.h>
#else
#include <windows.h>
#endif

MqttConfig *g_mqttConfig;

#define QOS 1
#define TIMEOUT 10000L

MQTTAsync client = NULL;
// 默认 MQTT 服务器地址

volatile int connected = 0;

void connlost(void *context, char *cause)
{
	MQTTAsync cl = (MQTTAsync)context;
	connected = 0;
	ZWX_LOG(LOG_ERR, "Connection lost");
	if (cause)
		ZWX_LOG(LOG_ERR, "     cause: %s", cause);
}

int msgarrvd(void *context, char *topicName, int topicLen, MQTTAsync_message *message)
{
	char *payload = (char *)message->payload;
	int ret = -1;

	if (payload == NULL)
	{
		ZWX_LOG(LOG_DEBUG, "Message arrived null payload");
		return 1;
	}

	ZWX_LOG(LOG_DEBUG, "Message arrived on topic %s: %s", topicName, (char *)message->payload);
	if (strcmp(g_mqttConfig->wakeHostTopic, topicName) == 0)
	{
		ret = wakehost(payload);
	}
	else if (strcmp(g_mqttConfig->shutdownHostTopic, topicName) == 0)
	{
		shutdownhost(payload);
	}
	else if (strcmp(g_mqttConfig->shutdownHostVirTopic, topicName) == 0)
	{
		ret = opvirhost(payload, 1);
	}
	else if (strcmp(g_mqttConfig->wakeVirHostTopic, topicName) == 0)
	{
		ret = opvirhost(payload, 2);
	}

	MQTTAsync_freeMessage(&message);
	MQTTAsync_free(topicName);
	return 1;
}
void onDisconnectFailure(void *context, MQTTAsync_failureData *response)
{
	MQTTAsync cl = (MQTTAsync)context;
	connected = 0;
	ZWX_LOG(LOG_ERR, "Disconnect failed, rc %d", response->code);
}
void onDisconnect(void *context, MQTTAsync_successData *response)
{
	MQTTAsync cl = (MQTTAsync)context;
	connected = 0;
	ZWX_LOG(LOG_ERR, "Successful disconnection");
}
void onSubscribe(void *context, MQTTAsync_successData *response)
{
	ZWX_LOG(LOG_DEBUG, "Subscribe succeeded");
}

void onSubscribeFailure(void *context, MQTTAsync_failureData *response)
{
	ZWX_LOG(LOG_ERR, "Subscribe failed, rc %d", response->code);
}
void onConnectFailure(void *context, MQTTAsync_failureData *response)
{
	ZWX_LOG(LOG_ERR, "Connect failed, rc %d", response->code);
	MQTTAsync cl = (MQTTAsync)context;
	connected = 0;
}
void onConnect(void *context, MQTTAsync_successData *response)
{
	MQTTAsync client = (MQTTAsync)context;
	MQTTAsync_responseOptions opts = MQTTAsync_responseOptions_initializer;
	int rc;

	ZWX_LOG(LOG_DEBUG, "Successful connection");
	connected = 2;
	opts.onSuccess = onSubscribe;
	opts.onFailure = onSubscribeFailure;
	opts.context = client;
	if ((rc = MQTTAsync_subscribe(client, g_mqttConfig->shutdownHostTopic, QOS, &opts)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to start subscribe %s, return code %d", g_mqttConfig->shutdownHostTopic, rc);
	} else {
		ZWX_LOG(LOG_ERR, "start subscribe %s", g_mqttConfig->shutdownHostTopic);
	}
	if ((rc = MQTTAsync_subscribe(client, g_mqttConfig->shutdownHostVirTopic, QOS, &opts)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to start subscribe %s, return code %d", g_mqttConfig->shutdownHostVirTopic, rc);
	} else {
		ZWX_LOG(LOG_ERR, "start subscribe %s", g_mqttConfig->shutdownHostVirTopic);
	}
	if ((rc = MQTTAsync_subscribe(client, g_mqttConfig->wakeHostTopic, QOS, &opts)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to start subscribe %s, return code %d", g_mqttConfig->wakeHostTopic, rc);
	} else {
		ZWX_LOG(LOG_ERR, "start subscribe %s", g_mqttConfig->wakeHostTopic);
	}
	if ((rc = MQTTAsync_subscribe(client, g_mqttConfig->wakeVirHostTopic, QOS, &opts)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to start subscribe %s, return code %d", g_mqttConfig->wakeVirHostTopic, rc);
	} else {
		ZWX_LOG(LOG_ERR, "start subscribe %s", g_mqttConfig->wakeVirHostTopic);
	}
}

void mqtt_start(char *broker_address, char *username, char *password)
{
	int rc;

	if (client != NULL)
	{
		MQTTAsync_destroy(&client);
	}

	ZWX_LOG(LOG_INFO, "broker_address %s client_id %s", broker_address, g_mqttConfig->clientId);
	if ((rc = MQTTAsync_create(&client, broker_address, g_mqttConfig->clientId, MQTTCLIENT_PERSISTENCE_NONE, NULL)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to create client, return code %d", rc);
		MQTTAsync_destroy(&client);
		return;
	}

	if ((rc = MQTTAsync_setCallbacks(client, client, connlost, msgarrvd, NULL)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to set callbacks, return code %d", rc);
		MQTTAsync_destroy(&client);
		return;
	}

	MQTTAsync_connectOptions conn_opts = MQTTAsync_connectOptions_initializer;
	conn_opts.keepAliveInterval = 20;
	conn_opts.cleansession = 1;
	conn_opts.onSuccess = onConnect;
	conn_opts.onFailure = onConnectFailure;
	conn_opts.username = g_mqttConfig->userName;
	conn_opts.password = g_mqttConfig->password;
	conn_opts.context = client;

	if ((rc = MQTTAsync_connect(client, &conn_opts)) != MQTTASYNC_SUCCESS)
	{
		ZWX_LOG(LOG_ERR, "Failed to start connect, return code %d", rc);
		MQTTAsync_destroy(&client);
		return;
	}
	connected = 1;
}

void stop_mqtt()
{
	int rc;
	MQTTAsync_disconnectOptions disc_opts = MQTTAsync_disconnectOptions_initializer;
	disc_opts.onSuccess = onDisconnect;
	disc_opts.onFailure = onDisconnectFailure;

	if (client)
	{
		if ((rc = MQTTAsync_disconnect(client, &disc_opts)) != MQTTASYNC_SUCCESS)
		{
			ZWX_LOG(LOG_ERR, "Failed to start disconnect, return code %d", rc);
			MQTTAsync_destroy(&client);
			connected = 0;
		}
	}
}
extern pthread_rwlock_t mqtt_rwlock;

void *mqtt_thread(void *arg)
{
	char userName[MAX_STR + 1] = "";
	char password[MAX_STR + 1] = "";
	char brokerAddress[MAX_STR + 1] = "";

	g_mqttConfig = (MqttConfig *)malloc(sizeof(MqttConfig));
	if (g_mqttConfig == NULL)
	{
		ZWX_LOG(LOG_ERR, "Failed to allocate memory for MqttConfig");
		exit(EXIT_FAILURE);
	}
	pthread_rwlock_wrlock(&mqtt_rwlock);
	memset(g_mqttConfig, 0, sizeof(g_mqttConfig));

	if (strlen(g_appConfig.serverAddress) > 0)
	{
		if(g_appConfig.port != 0 ) {
			snprintf(g_mqttConfig->serverAddress, MAX_STR, "tcp://%s:%d", g_appConfig.serverAddress,g_appConfig.port);
		} else {
			snprintf(g_mqttConfig->serverAddress, MAX_STR, "tcp://%s:1883", g_appConfig.serverAddress);
		}
		
		snprintf(g_mqttConfig->userName, MAX_STR, "%s", g_appConfig.userName);
		snprintf(g_mqttConfig->password, MAX_STR, "%s", g_appConfig.password);

		snprintf(brokerAddress, MAX_STR, "%s", g_mqttConfig->serverAddress);
		snprintf(userName, MAX_STR, "%s", g_appConfig.userName);
		snprintf(password, MAX_STR, "%s", g_appConfig.password);
	}
	snprintf(g_mqttConfig->clientId, MAX_STR, "device-%s", g_appConfig.macAddressStr);
	snprintf(g_mqttConfig->keepaliveTopic, MAX_STR, "/keepalive");
	snprintf(g_mqttConfig->shutdownHostTopic, MAX_STR, "/%s/shutdownhost", g_appConfig.uniqueid);
	snprintf(g_mqttConfig->shutdownHostResTopic, MAX_STR, "/shutdownhost/response");
	snprintf(g_mqttConfig->shutdownHostVirTopic, MAX_STR, "/%s/shutdownvirhost", g_appConfig.uniqueid);
	snprintf(g_mqttConfig->wakeHostTopic, MAX_STR, "/%s/wakehost", g_appConfig.uniqueid);
	snprintf(g_mqttConfig->wakeVirHostTopic, MAX_STR, "/%s/wakevirhost", g_appConfig.uniqueid);
	pthread_rwlock_unlock(&mqtt_rwlock);

	while (1)
	{
		pthread_rwlock_wrlock(&mqtt_rwlock);
		if (g_mqttConfig->configChanged == 1)
		{
			snprintf(brokerAddress, MAX_STR, "%s", g_mqttConfig->serverAddress);
			snprintf(userName, MAX_STR, "%s", g_mqttConfig->userName);
			snprintf(password, MAX_STR, "%s", g_mqttConfig->password);
			g_mqttConfig->configChanged = 0;
			stop_mqtt();
		}
		pthread_rwlock_unlock(&mqtt_rwlock);
		if (strlen(brokerAddress) > 0)
		{
			if (connected == 0)
			{
				sleep(5);
				mqtt_start(brokerAddress, userName, password);
			}
			else
			{
				sleep(5);
			}
		}
		else
		{
			sleep(5);
		}
	}
}