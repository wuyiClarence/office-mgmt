#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ifaddrs.h>
#include <sys/ioctl.h>
#include <net/if.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <linux/if_link.h>
#include <netinet/in.h>
#include "netface.h"
#include "zwx_log.h"

int get_mac_address(const char *interf, unsigned char *macAddress)
{
    int fd;
    struct ifreq ifr;
    unsigned char *mac;
    int ret = -1;

    // 创建一个套接字
    fd = socket(AF_INET, SOCK_DGRAM, 0);
    if (fd == -1)
    {
        ZWX_LOG(LOG_ERR, "socket SOCK_DGRAM fail %s", interf);
        perror("socket");
        return ret;
    }

    // 将接口名称复制到 ifr 结构中
    strncpy(ifr.ifr_name, interf, IFNAMSIZ - 1);

    // 使用 ioctl 获取 MAC 地址
    if (ioctl(fd, SIOCGIFHWADDR, &ifr) == -1)
    {
        ZWX_LOG(LOG_ERR, "ioctl SIOCGIFHWADDR error %s", interf);
        perror("ioctl");
        close(fd);
        return ret;
    }

    // 关闭套接字
    close(fd);

    // 从 ifr 中提取 MAC 地址
    mac = (unsigned char *)ifr.ifr_hwaddr.sa_data;

    macAddress[0] = mac[0];
    macAddress[1] = mac[1];
    macAddress[2] = mac[2];
    macAddress[3] = mac[3];
    macAddress[4] = mac[4];
    macAddress[5] = mac[5];
    ret = 0;
    // 打印 MAC 地址

    ZWX_LOG(LOG_INFO, "MAC address of %s: %02x:%02x:%02x:%02x:%02x:%02x\n",
            interf, mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]);
    return ret;
}

void get_ip_address(const char *interf, char *ipaddr)
{
    struct ifaddrs *ifaddr, *ifa;

    int found = 0;

    // 获取网络接口地址信息
    if (getifaddrs(&ifaddr) == -1)
    {
        ZWX_LOG(LOG_ERR, "getifaddrs %s", interf);
        perror("getifaddrs");
        exit(EXIT_FAILURE);
    }

    // 遍历每个接口
    for (ifa = ifaddr; ifa != NULL; ifa = ifa->ifa_next)
    {
        if (ifa->ifa_addr == NULL)
            continue;

        // 检查接口名称是否匹配
        if (strcmp(ifa->ifa_name, interf) == 0)
        {
            // 检查是否为 IPv4 地址
            if (ifa->ifa_addr->sa_family == AF_INET)
            {
                struct sockaddr_in *sa = (struct sockaddr_in *)ifa->ifa_addr;
                inet_ntop(AF_INET, &sa->sin_addr, ipaddr, INET_ADDRSTRLEN);
                ZWX_LOG(LOG_DEBUG, "Interface: %s\tIPv4 Address: %s\n", ifa->ifa_name, ipaddr);
                found = 1;
                break; // 找到后退出循环
            }
            // 检查是否为 IPv6 地址
            // else if (ifa->ifa_addr->sa_family == AF_INET6) {
            //     struct sockaddr_in6 *sa = (struct sockaddr_in6 *)ifa->ifa_addr;
            //     inet_ntop(AF_INET6, &sa->sin6_addr, addr, sizeof(addr));
            //     printf("Interface: %s\tIPv6 Address: %s\n", ifa->ifa_name, addr);
            //     found = 1;
            //     break;  // 找到后退出循环
            // }
        }
    }

    if (!found)
    {
        ZWX_LOG(LOG_DEBUG, "No IP address found for interface: %s\n", interf);
    }

    // 释放内存
    freeifaddrs(ifaddr);
}

int get_default_interface(char *default_interface)
{
    int ret = -1;
    FILE *fp = popen("ip route | grep default | awk '{print $5}'", "r");
    if (fp == NULL)
    {
        ZWX_LOG(LOG_ERR, "popen");
        perror("popen");
        return ret;
    }

    if (fgets(default_interface, IFNAMSIZ, fp) != NULL)
    {
        // 去掉换行符
        default_interface[strcspn(default_interface, "\n")] = 0;
        pclose(fp);
        ret = 0;
        return ret;
    }
    else
    {
        ZWX_LOG(LOG_ERR, "No default interface found.");
        pclose(fp);
        return ret;
    }
}