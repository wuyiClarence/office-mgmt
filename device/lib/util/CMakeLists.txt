# CMake 最低版本要求
cmake_minimum_required(VERSION 3.10)

# 项目名称和语言
project(libLogLibrary C)

# 设置 C 标准
set(CMAKE_C_STANDARD 99)

# 包含头文件目录
include_directories(${CMAKE_SOURCE_DIR}/src)

# 将 liblog.c 编译成静态库或动态库
# 生成静态库
add_library(liblog_static STATIC liblog.c)

# 生成动态库
add_library(liblog_shared SHARED liblog.c)

# 设置输出的库名称
set_target_properties(liblog_static PROPERTIES OUTPUT_NAME "liblog")
set_target_properties(liblog_shared PROPERTIES OUTPUT_NAME "liblog")

# 安装库和头文件
install(TARGETS liblog_static liblog_shared
    ARCHIVE DESTINATION lib
    LIBRARY DESTINATION lib
    RUNTIME DESTINATION bin
)

# 安装头文件
install(FILES liblog.h DESTINATION include)
