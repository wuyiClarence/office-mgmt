#ifndef __VIRHOST_H__
#define __VIRHOST_H__
#include "cJSON.h"
#include "zwx_log.h"
cJSON *getvirhostlist();
int shutdownhost(char *data);
int opvirhost(char *data, int op);
#endif