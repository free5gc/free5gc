#ifndef __UTLT_BUFF_H__
#define __UTLT_BUFF_H__

#include <stdint.h>

#include "utlt_debug.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _Bufblk {
    void *buf;
    int size; // The capacity of the buf, it may be larger than alloc size
    int len; // The used length of buf, maintain the next position to be written
} Bufblk;

/********************************************************************
 * Please call BufblkAlloc(...) to new a Bufblk,
 * do not declare a bufblk as global/local variable at outsize func.
 *
 * (O) Bufblk *bufblk = BufblkAlloc(0x40, sizeof(char));
 * (X) Bufblk bufblk; ...
 ********************************************************************/

Status BufblkPoolInit();
Status BufblkPoolFinal();
void BufblkPoolCheck(const char *showInfo);

Bufblk *BufblkAlloc(uint32_t num, uint32_t size);
Status BufblkResize(Bufblk *bufblk, uint32_t num, uint32_t size);
Status BufblkClear(Bufblk *bufblk);
Status BufblkFree(Bufblk *bufblk);

Status BufblkStr(Bufblk *bufblk, const char *str);
Status BufblkBuf(Bufblk *bufblk, const Bufblk *bufblk2);
Status BufblkFmt(Bufblk *bufblk, const char *fmt, ...)
        __attribute__((format(printf, 2, 3)));
Status BufblkBytes(Bufblk *bufblk, const char *str, uint32_t size);
Status BufblkAppend(Bufblk *bufblk, uint32_t num, uint32_t size);

void *UTLT_Malloc(uint32_t size);
void *UTLT_Calloc(uint32_t num, uint32_t size);
Status UTLT_Free(void *buf);
Status UTLT_Resize(void *buf, uint32_t size);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __UTLT_bufblk_H__ */
