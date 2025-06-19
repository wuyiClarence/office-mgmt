#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mqtt.h"
#include "virhost.h"
#include "config.h"

cJSON *getvirhostlist()
{
    FILE *fp;
    char output[1024];
    char id_str[16] = {0};
    char name[256], state[256];
    cJSON *root = cJSON_CreateArray();
    if (root == NULL)
    {
        ZWX_LOG(LOG_ERR, "cJSON Create Array Error");
        return root;
    }
    // 执行 'virsh list --all' 命令，获取所有虚拟机的信息
    fp = popen("virsh list --all", "r");
    if (fp == NULL)
    {
        perror("popen");
        return root;
    }
    // 跳过输出的表头（第一行）
    if (fgets(output, sizeof(output), fp) == NULL) {
        // 读取失败，处理错误（比如直接返回或报错）
        return root;
    }

    // 逐行读取命令输出并解析虚拟机名称和状态
    while (fgets(output, sizeof(output), fp) != NULL)
    {
        // 假设每一行的格式是 " ID   Name        State "
        // 需要跳过空白字符并提取名称和状态

        // 使用 sscanf 提取名称和状态

        int num_items = sscanf(output, "%15s %63s %63s", id_str, name, state);

        // 确保读取成功
        if (num_items == 3)
        {
            cJSON *object = cJSON_CreateObject();
            if (object == NULL)
            {
                ZWX_LOG(LOG_ERR, "cJSON Create Object Error\n");
                return root;
            }
            ZWX_LOG(LOG_DEBUG, "VM Name: %s, Status: %s\n", name, state);
            cJSON_AddStringToObject(object, "name", name);
            cJSON_AddStringToObject(object, "virtype", "kvm");
            cJSON_AddStringToObject(object, "state", state);
            cJSON_AddItemToArray(root, object);
        }
    }

    // 关闭文件指针
    fclose(fp);

    return root;
}

/*
    {
		"mac": "00:01:20:33:44:55",
		"virhost": [{
			"name": "win10-wudaoyuan",
			"virtype": "kvm"
		}]
    }
*/
int opvirhost(char *data, int op)
{
    int ret = -1;
    char command[256];
    cJSON *root = NULL;
    cJSON *mac = NULL;
    cJSON *virhosts = NULL;
    cJSON *object = NULL;

    root = cJSON_Parse(data);
    if (root == NULL)
    {
        ZWX_LOG(LOG_ERR, "Parse opvirhost message error");
        return ret;
    }
    mac = cJSON_GetObjectItem(root, "mac");
    if (mac && (mac->type == cJSON_String) && (mac->valuestring != NULL))
    {

        if (strcmp(mac->valuestring, g_appConfig.macAddressStr) != 0)
        {
            ZWX_LOG(LOG_ERR, "opvirhost message mac error,req %s, our %s", mac->valuestring, g_appConfig.macAddressStr);
            cJSON_Delete(root);
            return ret;
        }
    }

    virhosts = cJSON_GetObjectItem(root, "virhost");

    if (virhosts == NULL || (virhosts->type != cJSON_Array))
    {
        cJSON_Delete(root);
        return ret;
    }

    cJSON_ArrayForEach(object, virhosts)
    {
        cJSON *virtype = cJSON_GetObjectItem(object, "virtype");
        if (virtype && (virtype->type == cJSON_String) && (virtype->valuestring != NULL))
        {
            if (strcmp(virtype->valuestring, "kvm") == 0)
            {
                cJSON *name = cJSON_GetObjectItem(object, "name");
                if (name && (name->type == cJSON_String) && (name->valuestring != NULL))
                {

                    if (op == 1)
                    {
                        // 构建 virsh shutdown 命令
                        snprintf(command, sizeof(command), "virsh shutdown %s", name->valuestring);
                    }
                    else if (op == 2)
                    {
                        // 构建 virsh start 命令
                        snprintf(command, sizeof(command), "virsh start %s", name->valuestring);
                    }

                    // 执行命令
                    int result = system(command);
                    if (result == -1)
                    {
                        ZWX_LOG(LOG_ERR, "system command error %s", command);
                        perror("system");
                    }
                    else
                    {
                        ZWX_LOG(LOG_INFO, "command sent to VM %s.\n", command);
                    }
                }
            }
        }
    }
    cJSON_Delete(root);
    return 0;
}

int shutdownhost(char *data)
{
    cJSON *root = NULL;
    char command[256];
    cJSON *mac = NULL;
    int ret = -1;

    root = cJSON_Parse(data);
    if (root == NULL)
    {
        ZWX_LOG(LOG_ERR, "Parse opvirhost message error");
        return ret;
    }
    mac = cJSON_GetObjectItem(root, "mac");
    if (mac && (mac->type == cJSON_String) && (mac->valuestring != NULL))
    {

        if (strcmp(mac->valuestring, g_appConfig.macAddressStr) != 0)
        {
            ZWX_LOG(LOG_ERR, "opvirhost message mac error,req %s, our %s", mac->valuestring, g_appConfig.macAddressStr);
            cJSON_Delete(root);
            return ret;
        }
    }

    snprintf(command, sizeof(command), "shutdown -h now");
    // 执行命令
    int result = system(command);
    if (result == -1)
    {
        ZWX_LOG(LOG_ERR, "system command error %s", command);
        perror("system");
    }
    else
    {
        char *responsemsg = NULL;

        MQTTAsync_responseOptions opts = MQTTAsync_responseOptions_initializer;
        MQTTAsync_message pubmsg = MQTTAsync_message_initializer;

        cJSON_AddNumberToObject(root, "status", 1);
        responsemsg = cJSON_Print(root);

        // opts.onSuccess = onSend;
        // opts.onFailure = onSendFailure;
        opts.context = client;
        pubmsg.payload = responsemsg;
        pubmsg.payloadlen = (int)strlen(responsemsg);
        pubmsg.qos = 1;
        pubmsg.retained = 0;

        MQTTAsync_sendMessage(client, g_mqttConfig->shutdownHostResTopic, &pubmsg, &opts);

        free(responsemsg);
    }
    cJSON_Delete(root);
}