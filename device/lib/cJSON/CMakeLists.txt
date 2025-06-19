# CMake 最低版本要求
cmake_minimum_required(VERSION 3.10)

# 项目名称和语言
project(cJSONLibrary C)

# 设置 C 标准
set(CMAKE_C_STANDARD 99)

# 包含头文件目录
include_directories(${CMAKE_SOURCE_DIR}/src)

# 将 cJSON.c 编译成静态库或动态库
# 生成静态库
add_library(cjson_static STATIC cJSON.c)

# 生成动态库
add_library(cjson_shared SHARED cJSON.c)

# 设置输出的库名称
set_target_properties(cjson_static PROPERTIES OUTPUT_NAME "cjson")
set_target_properties(cjson_shared PROPERTIES OUTPUT_NAME "cjson")

# 安装库和头文件
install(TARGETS cjson_static cjson_shared
    ARCHIVE DESTINATION lib
    LIBRARY DESTINATION lib
    RUNTIME DESTINATION bin
)

# 安装头文件
install(FILES cJSON.h DESTINATION include)
