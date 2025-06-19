#include "config.h"
#include "netface.h"
#include "md5.h"
#include "cJSON.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "zwx_log.h"

// 读取文件内容
char *read_file()
{
    size_t readLen = 0;
    FILE *file = fopen("config.json", "r");
    if ( file == NULL)
    {
        file = fopen("/etc/devicemgmt/config.json", "r");
    }
    if ( file == NULL)
    {
        ZWX_LOG(LOG_ERR, "Failed to open file %s", "/etc/devicemgmt/config.json");
        perror("Failed to open file");
        return NULL;
    }

    fseek(file, 0, SEEK_END);
    long file_size = ftell(file);
    fseek(file, 0, SEEK_SET);

    char *content = (char *)malloc(file_size + 1);
    if (!content)
    {
        ZWX_LOG(LOG_ERR, "Failed to allocate memory");
        perror("Failed to allocate memory");
        fclose(file);
        return NULL;
    }

    readLen = fread(content, 1, file_size, file);
    content[readLen] = '\0';

    fclose(file);
    return content;
}

// 解析 JSON 配置文件
void parse_config(const char *config_data)
{
    cJSON *json = cJSON_Parse(config_data);
    if (!json)
    {
        ZWX_LOG(LOG_ERR, "Failed to parse JSON");
        fprintf(stderr, "Failed to parse JSON\n");
        return;
    }

    const cJSON *hostname = cJSON_GetObjectItem(json, "hostname");
    const cJSON *serverAddress = cJSON_GetObjectItem(json, "serverAddress");
    const cJSON *userName = cJSON_GetObjectItem(json, "userName");
    const cJSON *password = cJSON_GetObjectItem(json, "password");
    const cJSON *interface = cJSON_GetObjectItem(json, "interface");
    const cJSON *uniqueid = cJSON_GetObjectItem(json, "uniqueid");
    const cJSON *port = cJSON_GetObjectItem(json,"port");

    if (hostname && (hostname->type == cJSON_String) && (hostname->valuestring != NULL) && strlen(hostname->valuestring) > 0)
    {
        snprintf(g_appConfig.hostname,MAX_STR,"%s", hostname->valuestring);
        ZWX_LOG(LOG_INFO, "hostname: %s", hostname->valuestring);
    }
    else
    {
        if (gethostname(g_appConfig.hostname, sizeof(g_appConfig.hostname)) == 0)
        {
            ZWX_LOG(LOG_INFO, "hostname: %s", g_appConfig.hostname);
        }
        else
        {
            ZWX_LOG(LOG_ERR, "gethostname");
            perror("gethostname");
            snprintf(g_appConfig.hostname, MAX_STR, "device");
        }
    }

    if (serverAddress && (serverAddress->type == cJSON_String) && (serverAddress->valuestring != NULL))
    {
        ZWX_LOG(LOG_INFO, "serverAddress: %s", serverAddress->valuestring);
        snprintf(g_appConfig.serverAddress, MAX_STR,"%s", serverAddress->valuestring);
    }
    if (port && (port->type == cJSON_Number))
    {
        ZWX_LOG(LOG_INFO, "port: %d", port->valueint );
        g_appConfig.port = port->valueint;
    }
    if (userName && (userName->type == cJSON_String) && (userName->valuestring != NULL))
    {
        ZWX_LOG(LOG_INFO, "userName: %s", userName->valuestring);
        snprintf(g_appConfig.userName, MAX_STR,"%s",userName->valuestring);
    }

    if (password && (password->type == cJSON_String) && (password->valuestring != NULL))
    {
        ZWX_LOG(LOG_INFO, "Password: %s", password->valuestring);
        snprintf(g_appConfig.password, MAX_STR,"%s",password->valuestring);
    }

    if (interface && (interface->type == cJSON_String) && (interface->valuestring != NULL))
    {
        ZWX_LOG(LOG_INFO, "interface: %s", interface->valuestring);
        snprintf(g_appConfig.devInterface, IFNAMSIZ,"%s",interface->valuestring);
    }

    if (uniqueid && (uniqueid->type == cJSON_String) && (uniqueid->valuestring != NULL))
    {
        ZWX_LOG(LOG_INFO, "uniqueid: %s", uniqueid->valuestring);
        snprintf(g_appConfig.uniqueid, MAX_STR,"%s",uniqueid->valuestring);
    }

    cJSON_Delete(json);
}

int InitConfig()
{
    int ret = -1;
    char encrypt_source[MAX_STR] = {0};
    unsigned char md5_sign[16] = {0};

    memset(&g_appConfig, 0, sizeof(g_appConfig));

    
    char *config_data = read_file();

    if (config_data)
    {
        parse_config(config_data);
        free(config_data);
    }

    if (strlen(g_appConfig.devInterface) > 0)
    {
        ret = get_mac_address(g_appConfig.devInterface, g_appConfig.macAddress);
    }
    if (ret == -1)
    {
        ret = get_default_interface(g_appConfig.devInterface);
        if (ret == 0)
        {
            ret = get_mac_address(g_appConfig.devInterface, g_appConfig.macAddress);
        }
    }
    if (ret == 0)
    {
        snprintf(g_appConfig.macAddressStr, 64, "%02X:%02X:%02X:%02X:%02X:%02X", g_appConfig.macAddress[0],
                 g_appConfig.macAddress[1],
                 g_appConfig.macAddress[2],
                 g_appConfig.macAddress[3],
                 g_appConfig.macAddress[4],
                 g_appConfig.macAddress[5]);
        snprintf(encrypt_source, MAX_STR, "zwx-%s", g_appConfig.macAddressStr);

        MD5_CTX md5ctx;
        MD5_Init(&md5ctx);
        MD5_Update(&md5ctx, (unsigned char *)encrypt_source, strlen(encrypt_source));
        MD5_Final(md5_sign, &md5ctx);

        snprintf(g_appConfig.uniqueid, MAX_STR, "%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X%02X",
                 md5_sign[0],
                 md5_sign[1],
                 md5_sign[2],
                 md5_sign[3],
                 md5_sign[4],
                 md5_sign[5],
                 md5_sign[6],
                 md5_sign[7],
                 md5_sign[8],
                 md5_sign[9],
                 md5_sign[10],
                 md5_sign[11],
                 md5_sign[12],
                 md5_sign[13],
                 md5_sign[14],
                 md5_sign[15]);

        ZWX_LOG(LOG_INFO, "uniqueid %s\n", g_appConfig.uniqueid);
    }

    return ret;
}