cmake_minimum_required(VERSION 3.10)
project(device_mgmt)

set(TOP_DIR "${CMAKE_SOURCE_DIR}/../..")



# 指定头文件目录（例如：包含 src 和 include 文件夹）
include_directories(${TOP_DIR}/lib/cJSON
                    ${TOP_DIR}/lib/paho.mqtt.c-1.3.13/src
)

# 设置 C 标准
set(CMAKE_C_STANDARD 99)

link_directories(${TOP_DIR}/lib/cJSON/build/Release
                 ${TOP_DIR}/lib/paho.mqtt.c-1.3.13/build/src/Release
                 ${TOP_DIR}/lib/openssl-1.1.1/_install/x64/Release/lib)

# 指定源文件
add_executable(device_mgmt 
main.c
mqtt.c)


target_link_libraries(device_mgmt cjson paho-mqtt3as-static libssl libcrypto ws2_32)

# Debug 构建输出目录
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY_DEBUG ${CMAKE_BINARY_DIR}/bin/Debug)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY_DEBUG ${CMAKE_BINARY_DIR}/lib/Debug)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${CMAKE_BINARY_DIR}/lib/Debug)

# Release 构建输出目录
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY_RELEASE ${CMAKE_BINARY_DIR}/bin/Release)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY_RELEASE ${CMAKE_BINARY_DIR}/lib/Release)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${CMAKE_BINARY_DIR}/lib/Release)
