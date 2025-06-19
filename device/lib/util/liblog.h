#ifndef __LIBLOG_H__
#define __LIBLOG_H__

#ifdef _WIN32
#include <Windows.h>
#include <stdio.h>

#define LOG_EMERG   0
#define LOG_ALERT   1
#define LOG_CRIT    2
#define LOG_ERR     3
#define LOG_WARNING 4
#define LOG_NOTICE  5
#define LOG_INFO    6
#define LOG_DEBUG   7

#define LOG_PRI(prio) ((prio) & 0x07)  // 提取优先级部分
#else
#include <syslog.h>
#endif
#include <stdarg.h>

extern int liblog_prio;
extern int liblog_init(const char *program_name, char *log_dir);
extern int liblog_conf_reload(void);
extern void liblog_pre(char *buf, int buf_len, int prio);
extern void liblog(const char* format, ...);
extern void vliblog(const char* format, va_list ap);
extern const char *prio_str[];

#define LIBLOG(prio, fmt, ...) \
	do { \
		char log_pre[256];                              \
		if (LOG_PRI(prio) <= LOG_PRI(liblog_prio)) {    \
			liblog_pre(log_pre, sizeof(log_pre), prio); \
			syslog(prio, "<%s> [%s,%d] "fmt"\n",  prio_str[LOG_PRI(prio)], __FUNCTION__, __LINE__, ##__VA_ARGS__); \
			liblog("%s [%s,%d] "fmt"\n", log_pre, __FUNCTION__, __LINE__, ##__VA_ARGS__); \
		} \
	} while(0)

/*
 * simple tutorial:
 * 1. Init. example:
 *         liblog_init("test", "/zwx/var/log");
 *
 * 2. Use LIBLOG to write log. example:
 *         LIBLOG(LOG_ERR, "test log");
 *
 * 3. reload configure alternately:
 *         liblog_conf_reload();
 *
 * 4. Cat log file [LIBLOG_TMP_PATH"/liblog.program name]. example:
 *        tail -f /zwx/var/log/liblog/liblog.portald
 *
 * 5. Log level is 0~7, map with syslog LOG_EMERG, LOG_CRIT, ...
 *    Log level is LOG_ERR defaultly, change log level, example:
 *        echo 7 > /zwx/var/log/liblog/.liblog.prio.portald
 *
 * 6  Log file max size is 1MB defaultly, change log file max size to 100KB, exapmle:
 *        echo 100 > /zwx/var/log/liblog/.liblog.size.portald
*/

#endif

