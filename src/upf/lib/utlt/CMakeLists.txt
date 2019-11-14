cmake_minimum_required(VERSION 3.5)

project(free5GC_utlt_logger C)

link_directories(${LOGGER_DST})

# Logger
add_custom_command(OUTPUT ${LOGGER_DST}
    COMMENT "Building utlt_logger"
    
    WORKING_DIRECTORY ${LOGGER_SRC}
    COMMAND go build -o ${LOGGER_DST}/liblogger.so -buildmode=c-shared
    COMMAND mv ${LOGGER_DST}/liblogger.h ${LOGGER_DST}/logger.h
)
add_custom_target(utlt_logger ALL
    DEPENDS ${LOGGER_DST}
    VERBATIM
)
link_directories(${LOGGER_DST})
