CP=cp
CC = $(cc_prefix)gcc
DESTBIN=${ROOTFS}/bin

CFLAGS += -I${ROOTFS}/include  -O2
LDFLAGS += -L${ROOTFS}/lib
LIBS= ${ROOTFS}/lib/libzcjson.a ${ROOTFS}/lib/libzlog.a ${ROOTFS}/lib/libpaho-mqtt3as.a ${ROOTFS}/lib/libmd5.a -lm -lpthread

EXEC1 = devicemgmt

OBJS1 = main.o mqtt.o netface.o config.c keepalive.c wakeonlan.c virhost.c mdns.c

.PHONY: $(EXEC1) install clean
all: $(EXEC1) install

$(EXEC1): $(OBJS1)
	$(CC) $(CFLAGS) -o $@ $(OBJS1) $(LDFLAGS) $(LIBS)

install: $(EXEC1)
	$(CP) -rf ${EXEC1} $(DESTBIN)
	
clean:
	-rm -f $(EXEC1)  *.elf *.gdb *.o

