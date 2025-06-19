#ifndef __WAKEONLAN_H__
#define __WAKEONLAN_H__
int send_wol_packet(const char *mac_str);
int wakehost(char *data);
#endif