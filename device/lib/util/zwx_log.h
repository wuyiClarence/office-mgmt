#ifndef __ZWX_LOG_H__
#define __ZWX_LOG_H__
	
#include <unistd.h>
#include "liblog.h"

#define LOG_ALLOC(P, SIZE)\
	do { \
		P = (typeof(*(P))*)malloc(SIZE); \
		if (!P) { \
			ZWX_LOG(LOG_ERR, "Alloc memory fail, length=%d, P="#P, (int)(SIZE)); \
		} \
	} while(0)

#define LOG_ZALLOC(P, SIZE)\
	do { \
		P = (typeof(*(P))*)malloc(SIZE); \
		if (!P) { \
			ZWX_LOG(LOG_ERR, "Alloc memory fail, length=%d, P="#P, (int)(SIZE)); \
		} else { \
			memset(P, 0, SIZE); \
		} \
	} while(0)
	
#define LOG_RALLOC(P, SIZE)\
		do { \
			P = (typeof(*(P))*)realloc(P, SIZE); \
			if (!P) { \
				ZWX_LOG(LOG_ERR, "Realloc memory fail, length=%d, P="#P, (int)(SIZE)); \
			} \
		} while(0)

#define LOG_FREE(P) \
	do { \
		if (P != NULL) { \
			free(P); \
			P = NULL; \
		} \
	} while(0)

#define LOG_SYSTEM(ret, cmd)  \
	do {                       \
		ret = cmd[0] ? system(cmd) : 0; \
		if (ret) {             \
			ZWX_LOG(LOG_ERR, "cmd exec err, ret=%d, cmd=%s", ret, cmd); \
		} else {               \
			ZWX_LOG(LOG_DEBUG, "cmd exec, cmd=%s", cmd); \
		}                      \
	} while (0);
	
#define LOG_SYSTEM_RETRY(ret, cmd) { \
	int i = 0;                 \
	do {                       \
		ret = cmd[0] ? system(cmd) : 0; \
		if (ret) {             \
			ZWX_LOG(LOG_ERR, "cmd exec err, retry=%d, ret=%d, cmd=%s", i, ret, cmd); \
		} else {               \
			ZWX_LOG(LOG_DEBUG, "cmd exec, cmd=%s", cmd); \
		}                      \
	} while (ret && (i++ < 5) && (!usleep(1000)));\
}

#define LOG_POPEN(pp, cmd)   \
	do {                      \
		pp = popen(cmd, "r"); \
		if (!pp) {            \
			ZWX_LOG(LOG_ERR, "popen err, cmd=%s", cmd); \
		} else {              \
			ZWX_LOG(LOG_DEBUG, "popen, cmd=%s", cmd); \
		}                     \
	} while (0);

#ifndef MAC2STR
#define MAC2STR(a) (a)[0], (a)[1], (a)[2], (a)[3], (a)[4], (a)[5]
#define MACSTR "%02hhX:%02hhX:%02hhX:%02hhX:%02hhX:%02hhX"
#define MAC2STRP(a, p) p((a)[0]), p((a)[1]), p((a)[2]), p((a)[3]), p((a)[4]), p((a)[5])
#endif

#ifndef IPSTR
#define IPSTR "%u.%u.%u.%u"
#define IP2STR(a) (((a) >> 24) & 0xff), (((a) >> 16) & 0xff), (((a) >> 8) & 0xff), ((a) & 0xff)
#endif

#define MAC_IS_ZERO(mac) (\
		   0 == ((char *)(mac))[0] \
		&& 0 == ((char *)(mac))[1] \
		&& 0 == ((char *)(mac))[2] \
		&& 0 == ((char *)(mac))[3] \
		&& 0 == ((char *)(mac))[4] \
		&& 0 == ((char *)(mac))[5])
	
#define MAC_SET_ZERO(mac) {\
		((char *)(mac))[0] = 0; \
		((char *)(mac))[1] = 0; \
		((char *)(mac))[2] = 0; \
		((char *)(mac))[3] = 0; \
		((char *)(mac))[4] = 0; \
		((char *)(mac))[5] = 0; }

#define SYSTEM_NO_LOG_POST " >/dev/null 2>&1 "

#if defined __GNUC__
#define LOG_VALUE_TRUE(log_lvl, value) ({  \
    int _log_value_true = __builtin_expect(!!(value), 0); \
    if (_log_value_true) {                 \
        ZWX_LOG(log_lvl, #value" is true"); \
    }                                      \
    _log_value_true;                       \
})
#else
#define LOG_VALUE_TRUE(log_lvl, value) ({  \
    int _log_value_true = !!(value);       \
    if (_log_value_true) {                 \
        ZWX_LOG(log_lvl, #value" is true"); \
    }                                      \
    _log_value_true;                       \
})
#endif

#define ZWX_LOG LIBLOG
#define ZWX_LOG_INIT liblog_init

#endif
