#ifndef __UTLT_EVENT_H__
#define __UTLT_EVENT_H__

#include <stdint.h>
#include <stdarg.h>
#include <fcntl.h>

#include "utlt_debug.h"
#include "utlt_mq.h"
#include "utlt_timer.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef uintptr_t EvtQId;

typedef struct {
    uintptr_t type;
    int argc;
    uintptr_t arg0;
    uintptr_t arg1;
    uintptr_t arg2;
    uintptr_t arg3;
    uintptr_t arg4;
    uintptr_t arg5;
    uintptr_t arg6;
    uintptr_t arg7;
} __attribute__ ((packed)) Event;

#define EVTQ_O_BLOCK        0
#define EVTQ_O_NONBLOCK     O_NONBLOCK

/**
 * @param  option: either EVTQ_O_BLOCK or EVTQ_O_NONBLOCK.
 * @return eqId or NULL on error.
 */
EvtQId EventQueueCreate(int option);

Status EventQueueDelete(EvtQId eqId);

/**
 * Push an event with parameters(0 to 8) into event queue.
 * 
 * @param  eqId:
 * @param  eventType:
 * @param  argc:        number of event parameters (0 <= argc <= 8)
 * @param  (uintptr_t): event parameter1
 * @param  (uintptr_t): event parameter2 ...
 * @return STATUS_OK or STATUS_EAGAIN if the queue is full and the oflag EVTQ_O_NONBLOCK was set.
 */
Status EventSend(EvtQId eqId, uintptr_t eventType, int argc, ...);

/**
 * @return  STATUS_OK or STATUS_EAGAIN if the queue is empty and the oflag O_NONBLOCK was set.
 */
Status EventRecv(EvtQId eqId, Event *event);

TimerBlkID EventTimerCreate(TimerList *timerList, int type, uint32_t duration, uintptr_t event);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* #ifndef __UTLT_EVENT_H__ */
