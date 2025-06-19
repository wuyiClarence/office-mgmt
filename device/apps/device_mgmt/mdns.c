#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <libmdns.h>
#include "mdns.h"
#ifdef _WIN32
#define _CRT_SECURE_NO_WARNINGS 1
#endif
#include <pthread.h>
#include <stdio.h>

#include <errno.h>
#include <signal.h>

#ifdef _WIN32
#include <winsock2.h>
#include <iphlpapi.h>
#define sleep(x) Sleep(x * 1000)
#else
#include <netdb.h>
#include <ifaddrs.h>
#include <net/if.h>
#include <sys/time.h>
#endif
#include "mqtt.h"
#include "cJSON.h"
#include "zwx_log.h"
static char addrbuffer[64];
static char entrybuffer[256];
static char namebuffer[256];
static char sendbuffer[1024];
static mdns_record_txt_t txtbuffer[128];

static struct sockaddr_in service_address_ipv4;
static struct sockaddr_in6 service_address_ipv6;

static int has_ipv4;
static int has_ipv6;

char *g_base64_str = NULL;
char *g_arecord_str = NULL;

static mdns_string_t
ipv4_address_to_string(char *buffer, size_t capacity, const struct sockaddr_in *addr,
                       size_t addrlen)
{
    char host[NI_MAXHOST] = {0};
    char service[NI_MAXSERV] = {0};
    int ret = getnameinfo((const struct sockaddr *)addr, (socklen_t)addrlen, host, NI_MAXHOST,
                          service, NI_MAXSERV, NI_NUMERICSERV | NI_NUMERICHOST);
    int len = 0;
    if (ret == 0)
    {
        if (addr->sin_port != 0)
            len = snprintf(buffer, capacity, "%s:%s", host, service);
        else
            len = snprintf(buffer, capacity, "%s", host);
    }
    if (len >= (int)capacity)
        len = (int)capacity - 1;
    mdns_string_t str;
    str.str = buffer;
    str.length = len;
    return str;
}

static mdns_string_t
ipv6_address_to_string(char *buffer, size_t capacity, const struct sockaddr_in6 *addr,
                       size_t addrlen)
{
    char host[NI_MAXHOST] = {0};
    char service[NI_MAXSERV] = {0};
    int ret = getnameinfo((const struct sockaddr *)addr, (socklen_t)addrlen, host, NI_MAXHOST,
                          service, NI_MAXSERV, NI_NUMERICSERV | NI_NUMERICHOST);
    int len = 0;
    if (ret == 0)
    {
        if (addr->sin6_port != 0)
            len = snprintf(buffer, capacity, "[%s]:%s", host, service);
        else
            len = snprintf(buffer, capacity, "%s", host);
    }
    if (len >= (int)capacity)
        len = (int)capacity - 1;
    mdns_string_t str;
    str.str = buffer;
    str.length = len;
    return str;
}

static mdns_string_t
ip_address_to_string(char *buffer, size_t capacity, const struct sockaddr *addr, size_t addrlen)
{
    if (addr->sa_family == AF_INET6)
        return ipv6_address_to_string(buffer, capacity, (const struct sockaddr_in6 *)addr, addrlen);
    return ipv4_address_to_string(buffer, capacity, (const struct sockaddr_in *)addr, addrlen);
}

// Open sockets for sending one-shot multicast queries from an ephemeral port
static int
open_client_sockets(int *sockets, int max_sockets, int port)
{
    // When sending, each socket can only send to one network interface
    // Thus we need to open one socket for each interface and address family
    int num_sockets = 0;

#ifdef _WIN32

    IP_ADAPTER_ADDRESSES *adapter_address = 0;
    ULONG address_size = 8000;
    unsigned int ret;
    unsigned int num_retries = 4;
    do
    {
        adapter_address = (IP_ADAPTER_ADDRESSES *)malloc(address_size);
        ret = GetAdaptersAddresses(AF_UNSPEC, GAA_FLAG_SKIP_MULTICAST | GAA_FLAG_SKIP_ANYCAST, 0,
                                   adapter_address, &address_size);
        if (ret == ERROR_BUFFER_OVERFLOW)
        {
            free(adapter_address);
            adapter_address = 0;
            address_size *= 2;
        }
        else
        {
            break;
        }
    } while (num_retries-- > 0);

    if (!adapter_address || (ret != NO_ERROR))
    {
        free(adapter_address);
        ZWX_LOG(LOG_ERR,"Failed to get network adapter addresses\n");
        return num_sockets;
    }

    int first_ipv4 = 1;
    int first_ipv6 = 1;
    for (PIP_ADAPTER_ADDRESSES adapter = adapter_address; adapter; adapter = adapter->Next)
    {
        if (adapter->TunnelType == TUNNEL_TYPE_TEREDO)
            continue;
        if (adapter->OperStatus != IfOperStatusUp)
            continue;

        for (IP_ADAPTER_UNICAST_ADDRESS *unicast = adapter->FirstUnicastAddress; unicast;
             unicast = unicast->Next)
        {
            if (unicast->Address.lpSockaddr->sa_family == AF_INET)
            {
                struct sockaddr_in *saddr = (struct sockaddr_in *)unicast->Address.lpSockaddr;
                if ((saddr->sin_addr.S_un.S_un_b.s_b1 != 127) ||
                    (saddr->sin_addr.S_un.S_un_b.s_b2 != 0) ||
                    (saddr->sin_addr.S_un.S_un_b.s_b3 != 0) ||
                    (saddr->sin_addr.S_un.S_un_b.s_b4 != 1))
                {
                    int log_addr = 0;
                    if (first_ipv4)
                    {
                        service_address_ipv4 = *saddr;
                        first_ipv4 = 0;
                        log_addr = 1;
                    }
                    has_ipv4 = 1;
                    if (num_sockets < max_sockets)
                    {
                        saddr->sin_port = htons((unsigned short)port);
                        int sock = mdns_socket_open_ipv4(saddr);
                        if (sock >= 0)
                        {
                            sockets[num_sockets++] = sock;
                            log_addr = 1;
                        }
                        else
                        {
                            log_addr = 0;
                        }
                    }
                    if (log_addr)
                    {
                        char buffer[128];
                        mdns_string_t addr = ipv4_address_to_string(buffer, sizeof(buffer), saddr,
                                                                    sizeof(struct sockaddr_in));
                        ZWX_LOG(LOG_DEBUG,"Local IPv4 address: %.*s\n", MDNS_STRING_FORMAT(addr));
                    }
                }
            }
            else if (unicast->Address.lpSockaddr->sa_family == AF_INET6)
            {
                struct sockaddr_in6 *saddr = (struct sockaddr_in6 *)unicast->Address.lpSockaddr;
                // Ignore link-local addresses
                if (saddr->sin6_scope_id)
                    continue;
                static const unsigned char localhost[] = {0, 0, 0, 0, 0, 0, 0, 0,
                                                          0, 0, 0, 0, 0, 0, 0, 1};
                static const unsigned char localhost_mapped[] = {0, 0, 0, 0, 0, 0, 0, 0,
                                                                 0, 0, 0xff, 0xff, 0x7f, 0, 0, 1};
                if ((unicast->DadState == NldsPreferred) &&
                    memcmp(saddr->sin6_addr.s6_addr, localhost, 16) &&
                    memcmp(saddr->sin6_addr.s6_addr, localhost_mapped, 16))
                {
                    int log_addr = 0;
                    if (first_ipv6)
                    {
                        service_address_ipv6 = *saddr;
                        first_ipv6 = 0;
                        log_addr = 1;
                    }
                    has_ipv6 = 1;
                    if (num_sockets < max_sockets)
                    {
                        saddr->sin6_port = htons((unsigned short)port);
                        int sock = mdns_socket_open_ipv6(saddr);
                        if (sock >= 0)
                        {
                            sockets[num_sockets++] = sock;
                            log_addr = 1;
                        }
                        else
                        {
                            log_addr = 0;
                        }
                    }
                    if (log_addr)
                    {
                        char buffer[128];
                        mdns_string_t addr = ipv6_address_to_string(buffer, sizeof(buffer), saddr,
                                                                    sizeof(struct sockaddr_in6));
                        ZWX_LOG(LOG_DEBUG,"Local IPv6 address: %.*s\n", MDNS_STRING_FORMAT(addr));
                    }
                }
            }
        }
    }

    free(adapter_address);

#else

    struct ifaddrs *ifaddr = 0;
    struct ifaddrs *ifa = 0;

    if (getifaddrs(&ifaddr) < 0)
        ZWX_LOG(LOG_ERR,"Unable to get interface addresses\n");

    int first_ipv4 = 1;
    int first_ipv6 = 1;
    for (ifa = ifaddr; ifa; ifa = ifa->ifa_next)
    {
        if (!ifa->ifa_addr)
            continue;
        if (!(ifa->ifa_flags & IFF_UP) || !(ifa->ifa_flags & IFF_MULTICAST))
            continue;
        if ((ifa->ifa_flags & IFF_LOOPBACK) || (ifa->ifa_flags & IFF_POINTOPOINT))
            continue;

        if (ifa->ifa_addr->sa_family == AF_INET)
        {
            struct sockaddr_in *saddr = (struct sockaddr_in *)ifa->ifa_addr;
            if (saddr->sin_addr.s_addr != htonl(INADDR_LOOPBACK))
            {
                int log_addr = 0;
                if (first_ipv4)
                {
                    service_address_ipv4 = *saddr;
                    first_ipv4 = 0;
                    log_addr = 1;
                }
                has_ipv4 = 1;
                if (num_sockets < max_sockets)
                {
                    saddr->sin_port = htons(port);
                    int sock = mdns_socket_open_ipv4(saddr);
                    if (sock >= 0)
                    {
                        sockets[num_sockets++] = sock;
                        log_addr = 1;
                    }
                    else
                    {
                        log_addr = 0;
                    }
                }
                if (log_addr)
                {
                    char buffer[128];
                    mdns_string_t addr = ipv4_address_to_string(buffer, sizeof(buffer), saddr,
                                                                sizeof(struct sockaddr_in));
                    ZWX_LOG(LOG_DEBUG,"Local IPv4 address: %.*s\n", MDNS_STRING_FORMAT(addr));
                }
            }
        }
        else if (ifa->ifa_addr->sa_family == AF_INET6)
        {
            struct sockaddr_in6 *saddr = (struct sockaddr_in6 *)ifa->ifa_addr;
            // Ignore link-local addresses
            if (saddr->sin6_scope_id)
                continue;
            static const unsigned char localhost[] = {0, 0, 0, 0, 0, 0, 0, 0,
                                                      0, 0, 0, 0, 0, 0, 0, 1};
            static const unsigned char localhost_mapped[] = {0, 0, 0, 0, 0, 0, 0, 0,
                                                             0, 0, 0xff, 0xff, 0x7f, 0, 0, 1};
            if (memcmp(saddr->sin6_addr.s6_addr, localhost, 16) &&
                memcmp(saddr->sin6_addr.s6_addr, localhost_mapped, 16))
            {
                int log_addr = 0;
                if (first_ipv6)
                {
                    service_address_ipv6 = *saddr;
                    first_ipv6 = 0;
                    log_addr = 1;
                }
                has_ipv6 = 1;
                if (num_sockets < max_sockets)
                {
                    saddr->sin6_port = htons(port);
                    int sock = mdns_socket_open_ipv6(saddr);
                    if (sock >= 0)
                    {
                        sockets[num_sockets++] = sock;
                        log_addr = 1;
                    }
                    else
                    {
                        log_addr = 0;
                    }
                }
                if (log_addr)
                {
                    char buffer[128];
                    mdns_string_t addr = ipv6_address_to_string(buffer, sizeof(buffer), saddr,
                                                                sizeof(struct sockaddr_in6));
                    ZWX_LOG(LOG_DEBUG,"Local IPv6 address: %.*s\n", MDNS_STRING_FORMAT(addr));
                }
            }
        }
    }

    freeifaddrs(ifaddr);

#endif

    return num_sockets;
}

// Callback handling parsing answers to queries sent
static int
query_callback(int sock, const struct sockaddr *from, size_t addrlen, mdns_entry_type_t entry,
               uint16_t query_id, uint16_t rtype, uint16_t rclass, uint32_t ttl, const void *data,
               size_t size, size_t name_offset, size_t name_length, size_t record_offset,
               size_t record_length, void *user_data)
{
    (void)sizeof(sock);
    (void)sizeof(query_id);
    (void)sizeof(name_length);
    (void)sizeof(user_data);
    mdns_string_t fromaddrstr = ip_address_to_string(addrbuffer, sizeof(addrbuffer), from, addrlen);
    const char *entrytype = (entry == MDNS_ENTRYTYPE_ANSWER) ? "answer" : ((entry == MDNS_ENTRYTYPE_AUTHORITY) ? "authority" : "additional");
    mdns_string_t entrystr =
        mdns_string_extract(data, size, &name_offset, entrybuffer, sizeof(entrybuffer));
    if (rtype == MDNS_RECORDTYPE_PTR)
    {
        mdns_string_t namestr = mdns_record_parse_ptr(data, size, record_offset, record_length,
                                                      namebuffer, sizeof(namebuffer));
        ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s PTR %.*s rclass 0x%x ttl %u length %d\n",
               MDNS_STRING_FORMAT(fromaddrstr), entrytype, MDNS_STRING_FORMAT(entrystr),
               MDNS_STRING_FORMAT(namestr), rclass, ttl, (int)record_length);
    }
    else if (rtype == MDNS_RECORDTYPE_SRV)
    {
        mdns_record_srv_t srv = mdns_record_parse_srv(data, size, record_offset, record_length,
                                                      namebuffer, sizeof(namebuffer));
        ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s SRV %.*s priority %d weight %d port %d\n",
               MDNS_STRING_FORMAT(fromaddrstr), entrytype, MDNS_STRING_FORMAT(entrystr),
               MDNS_STRING_FORMAT(srv.name), srv.priority, srv.weight, srv.port);
    }
    else if (rtype == MDNS_RECORDTYPE_A)
    {
        struct sockaddr_in addr;
        mdns_record_parse_a(data, size, record_offset, record_length, &addr);
        mdns_string_t addrstr =
            ipv4_address_to_string(namebuffer, sizeof(namebuffer), &addr, sizeof(addr));
        ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s A %.*s\n", MDNS_STRING_FORMAT(fromaddrstr), entrytype,
               MDNS_STRING_FORMAT(entrystr), MDNS_STRING_FORMAT(addrstr));
        if (g_arecord_str != NULL)
        {
            free(g_arecord_str);
        }
        g_arecord_str = (char *)malloc(addrstr.length + 1);
        if (g_arecord_str == NULL)
        {
            ZWX_LOG(LOG_ERR,"malloc %zu buffer error\n", addrstr.length + 1);
            return 0;
        }
        snprintf(g_arecord_str, addrstr.length + 1, "%.*s", MDNS_STRING_FORMAT(addrstr));
    }
    else if (rtype == MDNS_RECORDTYPE_AAAA)
    {
        struct sockaddr_in6 addr;
        mdns_record_parse_aaaa(data, size, record_offset, record_length, &addr);
        mdns_string_t addrstr =
            ipv6_address_to_string(namebuffer, sizeof(namebuffer), &addr, sizeof(addr));
        ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s AAAA %.*s\n", MDNS_STRING_FORMAT(fromaddrstr), entrytype,
               MDNS_STRING_FORMAT(entrystr), MDNS_STRING_FORMAT(addrstr));
    }
    else if (rtype == MDNS_RECORDTYPE_TXT)
    {
        size_t parsed = mdns_record_parse_txt(data, size, record_offset, record_length, txtbuffer,
                                              sizeof(txtbuffer) / sizeof(mdns_record_txt_t));
        size_t itxt = 0;
        for (itxt = 0; itxt < parsed; ++itxt)
        {
            if (txtbuffer[itxt].value.length)
            {
                ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s TXT %.*s = %.*s\n", MDNS_STRING_FORMAT(fromaddrstr),
                       entrytype, MDNS_STRING_FORMAT(entrystr),
                       MDNS_STRING_FORMAT(txtbuffer[itxt].key),
                       MDNS_STRING_FORMAT(txtbuffer[itxt].value));

                if (strncmp(txtbuffer[itxt].key.str, "base64", txtbuffer[itxt].key.length) == 0)
                {
                    if (g_base64_str != NULL)
                    {
                        free(g_base64_str);
                    }
                    g_base64_str = (char *)malloc(txtbuffer[itxt].value.length + 1);
                    if (g_base64_str == NULL)
                    {
                        ZWX_LOG(LOG_DEBUG,"malloc %zu buffer error\n", txtbuffer[itxt].value.length + 1);
                        continue;
                    }
                    snprintf(g_base64_str, txtbuffer[itxt].value.length + 1, "%.*s", MDNS_STRING_FORMAT(txtbuffer[itxt].value));
                }
            }
            else
            {
                ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s TXT %.*s\n", MDNS_STRING_FORMAT(fromaddrstr), entrytype,
                       MDNS_STRING_FORMAT(entrystr), MDNS_STRING_FORMAT(txtbuffer[itxt].key));
            }
        }
    }
    else
    {
        ZWX_LOG(LOG_DEBUG,"%.*s : %s %.*s type %u rclass 0x%x ttl %u length %d\n",
               MDNS_STRING_FORMAT(fromaddrstr), entrytype, MDNS_STRING_FORMAT(entrystr), rtype,
               rclass, ttl, (int)record_length);
    }
    return 0;
}

static int send_mdns_query(mdns_query_t *query, size_t count)
{
    int sockets[32];
    int query_id[32];
    int num_sockets = open_client_sockets(sockets, sizeof(sockets) / sizeof(sockets[0]), 0);
    if (num_sockets <= 0)
    {
        ZWX_LOG(LOG_ERR,"Failed to open any client sockets\n");
        return -1;
    }
    ZWX_LOG(LOG_DEBUG,"Opened %d socket%s for mDNS query\n", num_sockets, num_sockets ? "s" : "");

    size_t capacity = 2048;
    void *buffer = malloc(capacity);
    void *user_data = 0;
    size_t iq = 0;
    ZWX_LOG(LOG_DEBUG,"Sending mDNS query");
    for (iq = 0; iq < count; ++iq)
    {
        const char *record_name = "PTR";
        if (query[iq].type == MDNS_RECORDTYPE_SRV)
            record_name = "SRV";
        else if (query[iq].type == MDNS_RECORDTYPE_A)
            record_name = "A";
        else if (query[iq].type == MDNS_RECORDTYPE_AAAA)
            record_name = "AAAA";
        else
            query[iq].type = MDNS_RECORDTYPE_PTR;
        ZWX_LOG(LOG_DEBUG," : %s %s", query[iq].name, record_name);
    }
    ZWX_LOG(LOG_DEBUG,"\n");
    int isock = 0;
    for (isock = 0; isock < num_sockets; ++isock)
    {
        query_id[isock] =
            mdns_multiquery_send(sockets[isock], query, count, buffer, capacity, 0);
        if (query_id[isock] < 0)
            ZWX_LOG(LOG_ERR,"Failed to send mDNS query: %s\n", strerror(errno));
    }

    // This is a simple implementation that loops for 5 seconds or as long as we get replies
    int res;
    ZWX_LOG(LOG_DEBUG,"Reading mDNS query replies\n");
    int records = 0;
    do
    {
        struct timeval timeout;
        timeout.tv_sec = 10;
        timeout.tv_usec = 0;

        int nfds = 0;
        fd_set readfs;
        FD_ZERO(&readfs);

        for (isock = 0; isock < num_sockets; ++isock)
        {
            if (sockets[isock] >= nfds)
                nfds = sockets[isock] + 1;
            FD_SET(sockets[isock], &readfs);
        }

        res = select(nfds, &readfs, 0, 0, &timeout);
        if (res > 0)
        {
            for (isock = 0; isock < num_sockets; ++isock)
            {
                if (FD_ISSET(sockets[isock], &readfs))
                {
                    size_t rec = mdns_query_recv(sockets[isock], buffer, capacity, query_callback,
                                                 user_data, query_id[isock]);
                    if (rec > 0)
                        records += rec;
                }
                FD_SET(sockets[isock], &readfs);
            }
        }
    } while (res > 0);

    ZWX_LOG(LOG_DEBUG,"Read %d records\n", records);

    free(buffer);

    for (isock = 0; isock < num_sockets; ++isock)
        mdns_socket_close(sockets[isock]);
    ZWX_LOG(LOG_DEBUG,"Closed socket%s\n", num_sockets ? "s" : "");

    return 0;
}
const char base64_table[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

int base64_char_to_value(char c)
{
    if ('A' <= c && c <= 'Z')
        return c - 'A';
    if ('a' <= c && c <= 'z')
        return c - 'a' + 26;
    if ('0' <= c && c <= '9')
        return c - '0' + 52;
    if (c == '+')
        return 62;
    if (c == '/')
        return 63;
    return -1; // Invalid character for Base64
}
int base64_decode(const char *input, uint8_t *output)
{
    int len = strlen(input);
    int output_len = 0;
    int i = 0, j = 0;
    for (i = 0; i < len; i += 4)
    {
        int val = 0;
        int padding = 0;

        for (j = 0; j < 4; j++)
        {
            val <<= 6;
            if (input[i + j] == '=')
            {
                padding++;
            }
            else
            {
                int v = base64_char_to_value(input[i + j]);
                if (v == -1)
                    return -1; // Invalid Base64 character
                val |= v;
            }
        }

        output[output_len++] = (val >> 16) & 0xFF;
        if (padding < 2)
            output[output_len++] = (val >> 8) & 0xFF;
        if (padding < 1)
            output[output_len++] = val & 0xFF;
    }

    return output_len;
}
unsigned char *xor_encrypt_decrypt(const char *input, int input_len, const char *key)
{

    int key_len = strlen(key);
    int i = 0;
    unsigned char *output = (unsigned char *)malloc(input_len + 1); // 为结果分配内存

    for (i = 0; i < input_len; i++)
    {
        output[i] = input[i] ^ key[i % key_len]; // 使用密钥字符进行循环异或
    }
    output[input_len] = '\0'; // 添加字符串结束符

    return output;
}

extern pthread_rwlock_t mqtt_rwlock;

void *mdns_thread(void *arg)
{
    int ret;
    char userName[MAX_STR + 1] = "";
    char password[MAX_STR + 1] = "";
    char brokerAddress[MAX_STR + 1] = "";
    mdns_query_t query[1];
    size_t query_count = 0;
    char query_name[128] = "_officemgmt._tcp.local.";
    char xorkey[64] = "zwxlink";
    query[query_count].name = query_name;
    query[query_count].type = MDNS_RECORDTYPE_PTR;
    query[query_count].length = strlen(query[query_count].name);
    query_count++;
    while (1)
    {
        ret = send_mdns_query(query, query_count);

        if (g_base64_str != NULL)
        {
            int input_len = strlen(g_base64_str);
            uint8_t *decoded_output = malloc((input_len * 3) / 4);
            if (decoded_output != NULL)
            {
                int decoded_len = base64_decode(g_base64_str, decoded_output);

                if (decoded_len == -1)
                {
                    ZWX_LOG(LOG_ERR, "Invalid Base64 input");
                }
                else
                {
                    char *data = xor_encrypt_decrypt(decoded_output, decoded_len, xorkey);
                    if (data)
                    {
                        cJSON *root = cJSON_Parse(data);
                        if (!root)
                        {
                            ZWX_LOG(LOG_ERR, "Failed to parse JSON");
                            fprintf(stderr, "Failed to parse JSON\n");
                        }
                        else
                        {
                            cJSON *mqtt_broker = cJSON_GetObjectItem(root, "mqtt_broker");
                            cJSON *mqtt_username = cJSON_GetObjectItem(root, "mqtt_username");
                            cJSON *mqtt_password = cJSON_GetObjectItem(root, "mqtt_password");

                            if (mqtt_broker && (mqtt_broker->type == cJSON_String) && (mqtt_broker->valuestring != NULL))
                            {
                                if (mqtt_username && (mqtt_username->type == cJSON_String) && (mqtt_username->valuestring != NULL))
                                {
                                    if (mqtt_password && (mqtt_password->type == cJSON_String) && (mqtt_password->valuestring != NULL))
                                    {
                                        if (strcmp(mqtt_broker->valuestring, "tcp://127.0.0.1:1883") == 0)
                                        {
                                            if (g_arecord_str != NULL)
                                            {
                                                snprintf(brokerAddress, MAX_STR, "tcp://%s:1883", g_arecord_str);
                                            }
                                            else
                                            {
                                            }
                                        }
                                        else
                                        {
                                            snprintf(brokerAddress, MAX_STR, "%s", mqtt_broker->valuestring);
                                        }
                                        pthread_rwlock_wrlock(&mqtt_rwlock);
                                        if (strcmp(g_mqttConfig->serverAddress, brokerAddress) != 0)
                                        {
                                            ZWX_LOG(LOG_INFO,"new brocker address discovery %s", brokerAddress);
                                            snprintf(g_mqttConfig->serverAddress, MAX_STR, "%s", brokerAddress);
                                            g_mqttConfig->configChanged = 1;
                                        }
                                        if (strcmp(g_mqttConfig->userName, mqtt_username->valuestring) != 0)
                                        {
                                            ZWX_LOG(LOG_INFO,"new mqtt userName  %s", mqtt_username->valuestring);
                                            snprintf(g_mqttConfig->userName, MAX_STR, "%s", mqtt_username->valuestring);
                                            g_mqttConfig->configChanged = 1;
                                        }
                                        if (strcmp(g_mqttConfig->password, mqtt_password->valuestring) != 0)
                                        {
                                            ZWX_LOG(LOG_INFO,"new mqtt password  %s", mqtt_password->valuestring);
                                            snprintf(g_mqttConfig->password, MAX_STR, "%s", mqtt_password->valuestring);
                                            g_mqttConfig->configChanged = 1;
                                        }
                                        pthread_rwlock_unlock(&mqtt_rwlock);
                                    }
                                }
                            }
                            cJSON_Delete(root);
                        }
                        free(data);
                    }
                }
                free(decoded_output);
            }
            free(g_base64_str);
            g_base64_str = NULL;
        }
        if (g_arecord_str != NULL)
        {
            ZWX_LOG(LOG_DEBUG,"g_arecord_str:%s\n", g_arecord_str);
            free(g_arecord_str);
            g_arecord_str = NULL;
        }
        sleep(30);
    }
}