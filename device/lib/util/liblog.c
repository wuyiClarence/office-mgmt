#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#ifdef _WIN32
#include <Windows.h>
#else
#include <syslog.h>
#include <unistd.h>
#endif
#include <stdarg.h>
#include <time.h>
#include <signal.h>
#include "liblog.h"

#define FILE_NAME_SIZE 256
#define FILE_NAME_PRE "liblog."
#define PRIO_CONF_NAME_PRE ".liblog.prio."
#define FSMAX_CONF_NAME_PRE ".liblog.size."

static int file_size_max = 1024 * 1024;

int liblog_prio = LOG_ERR;

const char *prio_str[] = {
  "EMERG",   /* system is unusable */
  "ALERT",   /* action must be taken immediately */
  "CRIT",    /* critical conditions */
  "ERR",     /* error conditions */
  "WARNING", /* warning conditions */
  "NOTICE",  /* normal but significant condition */
  "INFO",    /* informational */
  "DEBUG",   /* debug-level messages */
};

struct liblog_conf_s {
	FILE *file;
	pthread_mutex_t file_lock;
	char program_name[128];
	char path[FILE_NAME_SIZE];
	int priority;
	char prio_conf_name[FILE_NAME_SIZE];
	int file_size;
	int file_size_max_kb;
	char fsmax_conf_name[FILE_NAME_SIZE];
};

static struct liblog_conf_s log = {
	.file = NULL,
	.program_name = {0},
	.priority = LOG_ERR,
	.file_size = 0,
	.file_size_max_kb = 1024
};

static int liblog_f_read_int(char *path, int *value)
{
	FILE *fp = NULL;
	int ret = 0;

	if (!(fp = fopen(path, "r"))) {
		return -1;
	}
	ret = fscanf(fp, "%d", value);

	fclose(fp);
	return ret;
}

static int liblog_f_write_int(char *path, int value)
{
	FILE *fp = NULL;

	if (!(fp = fopen(path, "w"))) {
		return -1;
	}
	fprintf(fp, "%d", value);

	fclose(fp);
	return 0;
}

static int liglog_conf_reload_int(char *conf_name, int *value)
{
	if (access(conf_name, 0)) {
		return liblog_f_write_int(conf_name, *value);
	} else {
		return liblog_f_read_int(conf_name, value);
	}
}

void liblog_prio_reload(void)
{
	liglog_conf_reload_int(log.prio_conf_name, &log.priority);
	liblog_prio = log.priority;
}

void liblog_fsmax_reload(void)
{
	liglog_conf_reload_int(log.fsmax_conf_name, &log.file_size_max_kb);
}

static int liblog_conf_init(void)
{
	snprintf(log.prio_conf_name, sizeof(log.prio_conf_name), "%s/"PRIO_CONF_NAME_PRE"%s", log.path, log.program_name);
	snprintf(log.fsmax_conf_name, sizeof(log.fsmax_conf_name), "%s/"FSMAX_CONF_NAME_PRE"%s", log.path, log.program_name);
	liblog_conf_reload();
	pthread_mutex_init(&log.file_lock,NULL);

	return 0;
}

int liblog_conf_reload(void)
{
	liblog_prio_reload();
	liblog_fsmax_reload();

	return 0;
}

static int liblog_open(void)
{
	char file_name[FILE_NAME_SIZE];

	if (!log.program_name[0]) {
		return -1;
	}

	snprintf(file_name, FILE_NAME_SIZE, "%s/"FILE_NAME_PRE"%s", log.path, log.program_name);

	if (log.file) {
		fclose(log.file);
	}
	log.file = fopen(file_name, "r+");
	if (!log.file) {
		log.file = fopen(file_name, "w+");
	}
	if (!log.file) {
		return -1;
	}
	fseek(log.file, 0, SEEK_END);
	log.file_size = ftell(log.file);

	return 0;
}


void liblog_handle_loglvl_signal(int sig)
{
	LIBLOG(LOG_ERR, "Moudle: %s old_level = %d, old_maxsize=%d", log.program_name, log.priority, log.file_size_max_kb);
	liblog_conf_reload();
	LIBLOG(LOG_ERR, "Moudle: %s new_level = %d, new_maxsize=%d", log.program_name, log.priority, log.file_size_max_kb);
}

void liblog_reg_loglevel_signal(void)
{
	struct sigaction act;
	act.sa_handler = liblog_handle_loglvl_signal;
	act.sa_flags = SA_RESTART;                                                                                                                                             
	sigemptyset(&act.sa_mask);
	sigaddset(&act.sa_mask, 49);
	sigaction(49, &act, NULL);
}


int liblog_init(const char *program_name, char *log_dir)
{
	int i = 0;
	char cmd[256] = {0};

	if (!program_name) {
		return -1;
	}

	if (!log_dir) {
		log_dir = ".";
	}

	snprintf(log.path, sizeof(log.path), "%s/liblog", log_dir);
	snprintf(log.program_name, sizeof(log.program_name), "%s", program_name);

	snprintf(cmd, sizeof(cmd), "mkdir -p %s && chmod a+w %s", log.path, log.path);
	system(cmd);

	liblog_reg_loglevel_signal();
	liblog_conf_init();
	liblog_open();

#ifndef _WIN32
	/* syslog init */
	openlog(program_name, LOG_PID|LOG_NDELAY, LOG_LOCAL5);
#endif
	return 0;
}

void liblog_pre(char *buf, int buf_len, int prio)
{
	time_t now;
	char timbuf[26];

	time(&now);

	snprintf(buf, buf_len, "%.15s [%s]: <%s>", 
		ctime_r(&now, timbuf) + 4, 
		log.program_name, 
		prio_str[LOG_PRI(prio)]);
}

void vliblog(const char* format, va_list ap)
{
	if (!log.file) {
		liblog_open();
		if (!log.file) {
			log.file_size += vprintf(format, ap);
			return;
		}
	}

	if (log.file_size > log.file_size_max_kb * 1024) {
		char file_name[FILE_NAME_SIZE];
		fclose(log.file);
		snprintf(file_name, FILE_NAME_SIZE, "%s/"FILE_NAME_PRE"%s", log.path, log.program_name);
		log.file = fopen(file_name, "w+");
		log.file_size = 0;
	}

	log.file_size += vfprintf(log.file, format, ap);

	fflush(log.file);
}

void liblog(const char* format, ...)
{
    va_list ap;
    va_start(ap, format);
	pthread_mutex_lock(&log.file_lock);
	vliblog(format, ap);
	pthread_mutex_unlock(&log.file_lock);
	va_end(ap);
	return;
}

