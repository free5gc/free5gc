#ifndef __UTLT_DEBUG_H__
#define __UTLT_DEBUG_H__

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef int Status;
#define STATUS_ERROR -1
#define STATUS_OK 0
#define STATUS_EAGAIN 1

Status UTLT_SetLogLevel(const char *level);
const char *UTLT_StrStatus(Status status);

int UTLT_LogPrint(int level, const char *filename, const int line,
                  const char *funcname, const char *fmt, ...);

#define UTLT_Panic(fmt, ...) \
    UTLT_LogPrint(0, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)
#define UTLT_Fatal(fmt, ...) \
    UTLT_LogPrint(1, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)
#define UTLT_Error(fmt, ...) \
    UTLT_LogPrint(2, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)
#define UTLT_Warning(fmt, ...) \
    UTLT_LogPrint(3, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)
#define UTLT_Info(fmt, ...) \
    UTLT_LogPrint(4, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)
#define UTLT_Debug(fmt, ...) \
    UTLT_LogPrint(5, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)
#define UTLT_Trace(fmt, ...) \
    UTLT_LogPrint(6, __FILE__, __LINE__, __func__, fmt, ## __VA_ARGS__)

#define UTLT_Assert(cond, expr, fmt, ...) \
    if (!(cond)) { \
        UTLT_Error(fmt, ## __VA_ARGS__); \
        expr; \
    }

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* #ifndef __UTLT_DEBUG_H__ */
