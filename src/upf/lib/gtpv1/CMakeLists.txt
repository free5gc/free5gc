cmake_minimum_required(VERSION 3.5)

project(free5GC_gtpv1 C)

set(LIBGTPNL_DST_SO "${LIBGTPNL_DST}/lib/libgtpnl.so")

# Build libgtpnl
# Add_custom_command does not create a new target. You have to define targets explicitly
# by add_executable, add_library or add_custom_target in order to make them visible to make.
add_custom_command(OUTPUT ${LIBGTPNL_DST_SO}
    # Display the given message before the commands are executed at build time
    COMMENT "Building libgtpnl"
    WORKING_DIRECTORY ${LIBGTPNL_SRC}
    
    COMMAND chmod +x git-version-gen
    COMMAND autoreconf -iv
    COMMAND ./configure --prefix=${LIBGTPNL_DST}
    COMMAND make -j`nproc`
    COMMAND make install
)
add_custom_target(libgtpnl ALL
    # This is ALL target 'libgtpnl', and it depends on ${LIBGTPNL_MAKEFILE}"
    # If the file exists, then commands related to that file won't be executed.
    # DONOT let other target depends on the same OUTPUT as current target,
    #   or it may be bad when doing parallel make.
    DEPENDS ${LIBGTPNL_DST_SO}

    # To make quotes printable
    VERBATIM
)

link_directories("${CMAKE_SOURCE_DIR}/lib" "${LIBGTPNL_DST}/lib" ${LOGGER_DST})

# Test
add_executable("${PROJECT_NAME}_test" test.c)
set_target_properties("${PROJECT_NAME}_test" PROPERTIES 
    OUTPUT_NAME "${BUILD_BIN_DIR}/testgtpv1"
)

target_link_libraries("${PROJECT_NAME}_test" free5GC_lib gtpnl logger)
target_include_directories("${PROJECT_NAME}_test" PRIVATE
    include
    ${LOGGER_DST}
    "${LIBGTPNL_DST}/include"
    "${CMAKE_SOURCE_DIR}/lib/utlt/include"
)
target_compile_options("${PROJECT_NAME}_test" PRIVATE -Wall -Werror)
add_dependencies("${PROJECT_NAME}_test" utlt_logger libgtpnl)
