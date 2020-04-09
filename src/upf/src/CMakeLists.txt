cmake_minimum_required(VERSION 3.5)

project(free5GC_UPF_main C)

link_directories("${CMAKE_SOURCE_DIR}/lib" "${LIBGTP5GNL_DST}/lib" ${LOGGER_DST})

# Sources
file(GLOB SRC_FILES
    "*.c"
    "n4/*.c"
    "up/*.c"
)
add_executable(${PROJECT_NAME} ${SRC_FILES})
set_target_properties(
    ${PROJECT_NAME}
    PROPERTIES 
        OUTPUT_NAME "${BUILD_BIN_DIR}/free5gc-upfd"
        SUFFIX ""
)

target_include_directories(${PROJECT_NAME} PRIVATE
    ${LOGGER_DST}
    "${LIBGTP5GNL_DST}/include"
    "${CMAKE_SOURCE_DIR}/src"
    "${CMAKE_SOURCE_DIR}/lib/gtpv1/include"
    "${CMAKE_SOURCE_DIR}/lib/knet/include"
    "${CMAKE_SOURCE_DIR}/lib/pfcp/include"
    "${CMAKE_SOURCE_DIR}/lib/utlt/include"
    "${CMAKE_SOURCE_DIR}/lib/utlt/logger/include"
)
target_link_libraries(${PROJECT_NAME} PRIVATE 
    free5GC_lib rt pthread gtp5gnl mnl logger yaml
)
target_compile_options(${PROJECT_NAME} PRIVATE -Wall -Werror)
add_dependencies(${PROJECT_NAME} utlt_logger libgtp5gnl)
