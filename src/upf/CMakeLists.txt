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

# Config files
set(CONFIG_SRC "${CMAKE_SOURCE_DIR}/config")
file(GLOB CONFIG_FILES "${CONFIG_SRC}/**")
add_custom_command(OUTPUT ${BUILD_CONFIG_DIR}
    COMMENT "Copying configs"
    WORKING_DIRECTORY ${CONFIG_SRC}

    COMMAND mkdir -p ${BUILD_CONFIG_DIR} >/dev/null 2>&1
    COMMAND cp "${CONFIG_SRC}/upfcfg.example.yaml" "${BUILD_CONFIG_DIR}/upfcfg.yaml"
    COMMAND cp "${CONFIG_SRC}/upfcfg.test.example.yaml" "${BUILD_CONFIG_DIR}/upfcfg.test.yaml"
)
add_custom_target(configs ALL DEPENDS ${BUILD_CONFIG_DIR} VERBATIM)

# Submodules
add_subdirectory(src)
add_subdirectory(lib/gtpv1)
add_subdirectory(lib/utlt)
