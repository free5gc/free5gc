#ifndef __UTLT_LIB_H__
#define __UTLT_LIB_H__

#include "utlt_debug.h"
#include "logger.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

#define MAX_FILE_PATH_STRLEN  128
#define MAX_IFNAME_STRLEN     40

GoString UTLT_CStr2GoStr(const char *str);
Status GetAbsPath(char *str);

// For win32
#if WORDS_BIGENDIAN
#define ENDIAN2(x1, x2) x1 x2
#define ENDIAN3(x1, x2, x3) x1 x2 x3
#define ENDIAN4(x1, x2, x3, x4) x1 x2 x3 x4
#define ENDIAN5(x1, x2, x3, x4, x5) x1 x2 x3 x4 x5
#define ENDIAN6(x1, x2, x3, x4, x5, x6) x1 x2 x3 x4 x5 x6
#define ENDIAN7(x1, x2, x3, x4, x5, x6, x7) x1 x2 x3 x4 x5 x6 x7
#define ENDIAN8(x1, x2, x3, x4, x5, x6, x7, x8) x1 x2 x3 x4 x5 x6 x7 x8
#else
#define ENDIAN2(x1, x2) x2 x1
#define ENDIAN3(x1, x2, x3) x3 x2 x1
#define ENDIAN4(x1, x2, x3, x4) x4 x3 x2 x1
#define ENDIAN5(x1, x2, x3, x4, x5) x5 x4 x3 x2 x1
#define ENDIAN6(x1, x2, x3, x4, x5, x6) x6 x5 x4 x3 x2 x1
#define ENDIAN7(x1, x2, x3, x4, x5, x6, x7) x7 x6 x5 x4 x3 x2 x1
#define ENDIAN8(x1, x2, x3, x4, x5, x6, x7, x8) x8 x7 x6 x5 x4 x3 x2 x1
#endif


#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif    /* #ifndef __UTLT_LIB_H__ */

