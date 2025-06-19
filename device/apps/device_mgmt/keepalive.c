#include <stdio.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <time.h>
#include <string.h>
#include "mqtt.h"
#include "cJSON.h"
#include "config.h"
#include "MQTTAsync.h"
#include "zwx_log.h"
#include "virhost.h"
#include "netface.h"

char *gen_keepalive_mesg()
{
    char *keepalivemsg = NULL;
    cJSON *root = NULL;
    cJSON *virhost = NULL;
    char ipAddress[INET6_ADDRSTRLEN];
    char time_buffer[20];
    time_t now;
    struct tm *time_info;
    time(&now);
    // 将时间转为本地时间
    time_info = localtime(&now);
    strftime(time_buffer, sizeof(time_buffer), "%Y-%m-%d %H:%M:%S", time_info);

    get_ip_address(g_appConfig.devInterface, ipAddress);



    root = cJSON_CreateObject();
    if (root)
    {
        cJSON_AddStringToObject(root, "mac", g_appConfig.macAddressStr);
        cJSON_AddStringToObject(root, "ip", ipAddress);
        cJSON_AddStringToObject(root, "name", g_appConfig.hostname);
        cJSON_AddStringToObject(root, "ostype", "linux");
        cJSON_AddStringToObject(root, "uniqueId", g_appConfig.uniqueid);
        cJSON_AddStringToObject(root, "curTime", time_buffer);
        virhost = getvirhostlist();
        if (virhost)
        {
            cJSON_AddItemToObject(root, "virhost", virhost);
        }
        keepalivemsg = cJSON_Print(root);
        cJSON_Delete(root);
    }

    return keepalivemsg;
}

void *keepalive_thread(void *arg)
{
    char *keepalivemsg = NULL;
    while (1)
    {
        if (connected == 2)
        {
            keepalivemsg = gen_keepalive_mesg();
            if (keepalivemsg != NULL)
            {
                MQTTAsync_responseOptions opts = MQTTAsync_responseOptions_initializer;
                MQTTAsync_message pubmsg = MQTTAsync_message_initializer;

                // opts.onSuccess = onSend;
                // opts.onFailure = onSendFailure;
                opts.context = client;
                pubmsg.payload = keepalivemsg;
                pubmsg.payloadlen = (int)strlen(keepalivemsg);
                pubmsg.qos = 1;
                pubmsg.retained = 0;

                MQTTAsync_sendMessage(client, g_mqttConfig->keepaliveTopic, &pubmsg, &opts);
                free(keepalivemsg);
            }
            ZWX_LOG(LOG_DEBUG, "Message keepalive published.");
        }
        sleep(30);
    }
}