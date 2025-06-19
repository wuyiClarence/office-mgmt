CC = $(cc_prefix)gcc

all: libzcjson.a libzcjson.so install

clean:
	rm -f *.so *.o *.a

cJSON.o: cJSON.c
	$(CC) $(CFLAGS) -fPIC -c cJSON.c -o cJSON.o $(LDFLAGS)

libzcjson.a: cJSON.o
	ar rcs libzcjson.a cJSON.o

libzcjson.so: cJSON.o
	$(CC) $(CFLAGS) -fPIC -shared cJSON.o -o libzcjson.so -lm $(LDFLAGS)

install:
	cp libzcjson.a $(ROOTFS)/lib/
	cp libzcjson.so $(ROOTFS)/lib/
	cp cJSON.h $(ROOTFS)/include
