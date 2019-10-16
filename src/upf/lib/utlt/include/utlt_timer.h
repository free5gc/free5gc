#ifndef __TIMER_H__
#define __TIMER_H__

#include <stdint.h>

#include "utlt_debug.h"
#include "utlt_list.h"
#include "utlt_time.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

#define MAX_NUM_OF_TIMER        1024

// paramID for setting timer parameter
#define PARAM1   0
#define PARAM2   1
#define PARAM3   2
#define PARAM4   3
#define PARAM5   4
#define PARAM6   5

#define TIMER_TYPE_PERIOD   0
#define TIMER_TYPE_ONCE     1

typedef struct _TimerList {
    ListNode active;
    ListNode idle;
} TimerList;

typedef uintptr_t TimerBlkID;

typedef void (*ExpireFunc)(uintptr_t data, uintptr_t param[]);

Status TimerPoolInit(void);
Status TimerFinal(void);
uint32_t TimerGetPoolSize(void);

void TimerListInit(TimerList *tmList);
Status TimerExpireCheck(TimerList *tmList, uintptr_t data);

Status TimerStart(TimerBlkID id);
Status TimerStop(TimerBlkID id);

TimerBlkID TimerCreate(TimerList *tmList, int type, uint32_t duration, ExpireFunc expireFunc);
void TimerDelete(TimerBlkID id);
Status TimerSet(int paramID, TimerBlkID id, uintptr_t param);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif
