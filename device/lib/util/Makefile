CC = $(cc_prefix)gcc

all: libzlog.a libzlog.so libmd5.so libmd5.a mdns install

clean:
	rm -f *.so *.o *.a

libmd5.o: md5.c
	$(CC) $(CFLAGS) -fPIC -c md5.c -o libmd5.o $(LDFLAGS)

libmd5.a: libmd5.o
	ar rcs libmd5.a libmd5.o

libmd5.so: libmd5.o
	$(CC) $(CFLAGS) -fPIC -shared libmd5.o -o libmd5.so -lm $(LDFLAGS)

liblog.o: liblog.c
	$(CC) $(CFLAGS) -fPIC -c liblog.c -o liblog.o $(LDFLAGS)

libzlog.a: liblog.o
	ar rcs libzlog.a liblog.o

libzlog.so: liblog.o
	$(CC) $(CFLAGS) -fPIC -shared liblog.o -o libzlog.so -lm $(LDFLAGS)

mdns.o: mdns.c
	$(CC) $(CFLAGS) -fPIC -c mdns.c -o mdns.o $(LDFLAGS)

mdns: mdns.o
	$(CC) $(CFLAGS) mdns.o -o mdns $(LDFLAGS)

install:
	cp libzlog.a $(ROOTFS)/lib/
	cp libzlog.so $(ROOTFS)/lib/
	cp liblog.h $(ROOTFS)/include
	cp libmd5.a $(ROOTFS)/lib/
	cp libmd5.so $(ROOTFS)/lib/
	cp md5.h $(ROOTFS)/include
	cp zwx_log.h $(ROOTFS)/include
	cp mdns $(ROOTFS)/bin/
	cp libmdns.h $(ROOTFS)/include
