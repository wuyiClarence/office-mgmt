#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <unistd.h>
#include "wakeonlan.h"
#include "zwx_log.h"
#include "config.h"
#include "cJSON.h"
#include <ifaddrs.h>

#define BROADCAST_IP "255.255.255.255"
#define WOL_PORT 9 // 默认的 WOL 端口
// 将 MAC 地址字符串转换为二进制格式
void mac_str_to_bytes(const char *mac_str, unsigned char *mac_bytes)
{
    int i;
    unsigned int temp;
    for (i = 0; i < 6; i++)
    {
        sscanf(mac_str + (i * 3), "%2x", &temp); // 每两个字符转为一个字节
        mac_bytes[i] = (unsigned char)temp;
    }
}

// 枚举网络接口，根据条件发送 Magic Packet
int send_wol_packet(const char *mac_str)
{
    struct ifaddrs *ifaddr, *ifa;
    int ret = 0;

    if (getifaddrs(&ifaddr) == -1)
    {
        perror("getifaddrs");
        return -1;
    }

    // 遍历所有接口
    for (ifa = ifaddr; ifa != NULL; ifa = ifa->ifa_next)
    {
        // 跳过没有地址或是回环接口
        if (ifa->ifa_addr == NULL || !(ifa->ifa_flags & IFF_UP) || (ifa->ifa_flags & IFF_LOOPBACK))
            continue;

        // 检查是否支持广播
        if (!(ifa->ifa_flags & IFF_BROADCAST))
            continue;

        // 获取该接口的广播地址，如果不为空则使用
        char bcast_ip[INET_ADDRSTRLEN] = {0};
        if (ifa->ifa_broadaddr && ifa->ifa_broadaddr->sa_family == AF_INET)
        {
            struct sockaddr_in *bcast = (struct sockaddr_in *)ifa->ifa_broadaddr;
            inet_ntop(AF_INET, &bcast->sin_addr, bcast_ip, sizeof(bcast_ip));
        }
        else
        {
            // 没有有效的广播地址时，使用默认值（可能适用于某些场景）
            strncpy(bcast_ip, BROADCAST_IP, sizeof(bcast_ip)-1);
        }

        // 发送 Magic Packet 到当前接口的广播地址
        if (send_wol_packet_via_interface(mac_str, ifa->ifa_name, bcast_ip) != 0)
        {
            fprintf(stderr, "Failed to send on interface: %s\n", ifa->ifa_name);
            ret = -1; // 标记失败，但可以继续尝试其他接口
        }
    }

    freeifaddrs(ifaddr);
    return ret;
}

// 发送 Wake-on-LAN Magic Packet
int send_wol_packet_via_interface(const char *mac_str, const char *if_name, const char *bcast_ip)
{
    unsigned char packet[102]; // Magic Packet的最大长度
    struct sockaddr_in dest_addr;
    int sockfd;
    int ret = -1;
    int i = 0;

    // 填充 Magic Packet
    memset(packet, 0xFF, 6); // 前 6 个字节是 0xFF
    unsigned char mac_bytes[6];
    mac_str_to_bytes(mac_str, mac_bytes);

    for (i = 6; i < sizeof(packet); i += 6)
    {
        memcpy(packet + i, mac_bytes, 6); // 后续部分是 MAC 地址，重复 16 次
    }

    // 创建 UDP 套接字
    sockfd = socket(AF_INET, SOCK_DGRAM, 0);
    if (sockfd < 0)
    {
        ZWX_LOG(LOG_ERR, "socket");
        perror("socket");
        return ret;
    }


    // 允许套接字发送广播
    int broadcast_enable = 1;
    if (setsockopt(sockfd, SOL_SOCKET, SO_BROADCAST, &broadcast_enable, sizeof(broadcast_enable)) < 0)
    {
        ZWX_LOG(LOG_ERR, "setsockopt");
        perror("setsockopt");
        close(sockfd);
        return ret;
    }

     // 绑定到指定的网络接口
    if (if_name && strlen(if_name) > 0) {
        if (setsockopt(sockfd, SOL_SOCKET, SO_BINDTODEVICE, if_name, strlen(if_name)) < 0) {
            ZWX_LOG(LOG_ERR, "setsockopt SO_BINDTODEVICE");
            perror("setsockopt SO_BINDTODEVICE");
            close(sockfd);
            return ret;
        }
    }


    // 设置广播地址
    memset(&dest_addr, 0, sizeof(dest_addr));
    dest_addr.sin_family = AF_INET;
    dest_addr.sin_port = htons(WOL_PORT);
    dest_addr.sin_addr.s_addr = inet_addr(bcast_ip);

    // 发送 Magic Packet
    if (sendto(sockfd, packet, sizeof(packet), 0, (struct sockaddr *)&dest_addr, sizeof(dest_addr)) < 0)
    {
        ZWX_LOG(LOG_ERR, "sendto");
        perror("sendto");
        close(sockfd);
        return ret;
    }

    ret = 0;
    ZWX_LOG(LOG_INFO, "Magic Packet sent to %s\n", mac_str);

    close(sockfd);
    return ret;
}

int wakehost(char *data)
{
    int ret = -1;
    cJSON *root = cJSON_Parse(data);
    if (root == NULL)
    {
        ZWX_LOG(LOG_ERR, "Parse wakehost message error");
        return ret;
    }
    cJSON *mac = cJSON_GetObjectItem(root, "mac");
    if (mac && (mac->type == cJSON_String) && (mac->valuestring != NULL))
    {
        if (strcmp(mac->valuestring, g_appConfig.macAddressStr) != 0)
        {
            ZWX_LOG(LOG_ERR, "wakehost message mac error,req %s, our %s", mac->valuestring, g_appConfig.macAddressStr);
            cJSON_Delete(root);
            return ret;
        }
    }
    cJSON *wakemac = cJSON_GetObjectItem(root, "wakemac");
    if (wakemac && (wakemac->type == cJSON_String) && (wakemac->valuestring != NULL))
    {

        ZWX_LOG(LOG_INFO, "wake host: %s\n", wakemac->valuestring);
        ret = send_wol_packet(wakemac->valuestring);
    }
    cJSON_Delete(root);
    return ret;
}