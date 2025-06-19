#!/bin/bash

set -euo pipefail

# 显示脚本使用说明
usage() {
  echo "Usage: $0 [--host=<host>] [--target=<target>] [--gcc=<path_to_gcc>] [--help]"
  echo
  echo "Options:"
  echo "  --host=<host>         Optional. Target host for cross-compilation (default: current system's host)"
  echo "  --target=<target>         Optional. Target compile (default: linux-generic32) linux-x86_64 linux-armv4、linux-armv5、linux-armv7、linux-aarch64"
  echo "  --gcc=<path_to_gcc>    Optional. Specify a custom GCC compiler path (default: system gcc)"
  echo "  --help, -h             Show this help message and exit"
  exit 0
}

# 初始化参数默认值
HOST=""            # 默认使用当前系统的 host
GCC="gcc"                 # 默认使用系统的 gcc
TARGET="linux-x86_64"

# 解析传入的参数
for arg in "$@"; do
  case $arg in
    --host=*)
      HOST="${arg#*=}"
      shift
      ;;
    --gcc=*)
      GCC="${arg#*=}"
      shift
      ;;
    --target=*)
      TARGET="${arg#*=}"
      shift
      ;;
    --help|-h)
      usage
      ;;
    *)
      echo "Unknown option: $arg"
      usage
      ;;
  esac
done

export TOPDIR=`pwd`
export ROOTFS=${TOPDIR}/install

mkdir -p $ROOTFS
mkdir -p $ROOTFS/bin
mkdir -p $ROOTFS/lib
mkdir -p $ROOTFS/include
# 配置编译选项
echo "Configuring with host=$HOST, install_dir=$ROOTFS, and gcc=$GCC..."

if [ -n $HOST  ]; then
  export cc_prefix="${HOST}-"
else
  export cc_prefix=""
fi

cd ${TOPDIR}/lib/cJSON
make 
make install

echo "Start build util"

cd ${TOPDIR}/lib/util
make 
make install


if [ ! -d ${TOPDIR}/lib/openssl-1.1.1 ]; then
    mkdir -p ${TOPDIR}/lib/openssl-1.1.1
    tar xzvf ${TOPDIR}/lib/openssl-OpenSSL_1_1_1w.tar.gz -C ${TOPDIR}/lib/openssl-1.1.1 --strip-components=1
fi

if [ ! -f ${TOPDIR}/lib/openssl-1.1.1/.build_ok ]; then
    cd ${TOPDIR}/lib/openssl-1.1.1
    if [ -n "$HOST" ]; then
        ./Configure ${TARGET} --cross-compile-prefix="${HOST}-" --prefix="$ROOTFS" no-shared no-tests
    else
        ./Configure ${TARGET} --prefix="$ROOTFS" no-shared no-tests
    fi
    make
    make install
    touch ${TOPDIR}/lib/openssl-1.1.1/.build_ok
fi

if [ ! -d ${TOPDIR}/lib/paho.mqtt.c-1.3.13 ]; then
    mkdir -p ${TOPDIR}/lib/paho.mqtt.c-1.3.13
    tar xzvf ${TOPDIR}/lib/paho.mqtt.c-1.3.13.tar.gz -C ${TOPDIR}/lib/paho.mqtt.c-1.3.13 --strip-components=1
fi

if [ ! -f ${TOPDIR}/lib/paho.mqtt.c-1.3.13/.build_ok ]; then
    cd ${TOPDIR}/lib/paho.mqtt.c-1.3.13
    if [ -n "$HOST" ]; then
        make LDFLAGS="-L${ROOTFS}/lib/" CFLAGS="-I$ROOTFS/include/" CC=${cc_prefix}
    else
        make LDFLAGS="-L${ROOTFS}/lib/" CFLAGS="-I${ROOTFS}/include/"
    fi
    make install prefix="$ROOTFS"
    touch ${TOPDIR}/lib/paho.mqtt.c-1.3.13/.build_ok
fi

cd ${TOPDIR}/apps/device_mgmt
make 