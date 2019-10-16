set(CMAKE_C_COMPILER "/usr/bin/gcc")
cmake_minimum_required(VERSION 3.5)

project(free5GC_UPF C)

# Build destination
set(BUILD_BIN_DIR "${CMAKE_BINARY_DIR}/bin")
set(BUILD_CONFIG_DIR "${CMAKE_BINARY_DIR}/config")

# Build paths
set(LIBGTPNL_SRC "${CMAKE_SOURCE_DIR}/lib/libgtpnl-1.2.1")
set(LIBGTPNL_DST "${CMAKE_BINARY_DIR}/libgtpnl")
set(LOGGER_SRC "${CMAKE_SOURCE_DIR}/lib/utlt/logger")
set(LOGGER_DST "${CMAKE_BINARY_DIR}/utlt_logger")

# Build environment
file(MAKE_DIRECTORY ${BUILD_BIN_DIR})

set(CONFIG_SRC "${CMAKE_SOURCE_DIR}/config")
file(GLOB CONFIG_FILES "${CONFIG_SRC}/**")
file(COPY ${CONFIG_FILES} DESTINATION ${BUILD_CONFIG_DIR})

# Submodules
add_subdirectory(src)
add_subdirectory(lib/utlt)
add_subdirectory(lib/gtpv1)
