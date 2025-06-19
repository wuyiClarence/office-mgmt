#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <cJSON.h>
#include <pthread.h>
#include "zwx_log.h"

#include "config.h"
#include "mqtt.h"
#include "keepalive.h"
#include "mdns.h"

pthread_t mqtt_tid;
pthread_t keepalive_tid;
pthread_t mdns_tid;

pthread_rwlock_t mqtt_rwlock; // 定义读写锁

AppConfig g_appConfig;

int main(int argc, char *argv[])
{
    ZWX_LOG_INIT("devicemgmt", "/tmp/zwx/log");
    int ret = 0;

    if (pthread_rwlock_init(&mqtt_rwlock, NULL) != 0)
    {
        ZWX_LOG(LOG_ERR, "Failed to initialize rwlock");
        perror("Failed to initialize rwlock");
        return EXIT_FAILURE;
    }

    ret = InitConfig();
    if (ret !=0) {
        ZWX_LOG(LOG_ERR, "Init Config Error");
        exit(EXIT_FAILURE);
    }
    if (pthread_create(&mqtt_tid, NULL, mqtt_thread, NULL) != 0)
    {
        ZWX_LOG(LOG_ERR, "Failed to create MQTT thread");
        perror("Failed to create MQTT thread");
        exit(EXIT_FAILURE);
    }
    if (pthread_create(&keepalive_tid, NULL, keepalive_thread, NULL) != 0)
    {
        ZWX_LOG(LOG_ERR, "Failed to create Keepalive thread");
        perror("Failed to create Keepalive thread");
        exit(EXIT_FAILURE);
    }
    if (pthread_create(&mdns_tid, NULL, mdns_thread, NULL) != 0)
    {
        ZWX_LOG(LOG_ERR, "Failed to create Mdns thread");
        perror("Failed to create Mdns thread");
        exit(EXIT_FAILURE);
    }

    pthread_join(mqtt_tid, NULL);
    pthread_join(keepalive_tid, NULL);
    pthread_join(mdns_tid, NULL);

    pthread_rwlock_destroy(&mqtt_rwlock);
    return 0;
}
