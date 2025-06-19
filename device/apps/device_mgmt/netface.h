#ifndef __NETFACE_H__
#define __NETFACE_H__
int get_default_interface(char *default_interface);
int get_mac_address(const char *interface, unsigned char *macAddress);
void get_ip_address(const char *interface, char *ipaddr);
#endif