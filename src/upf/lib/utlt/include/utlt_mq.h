#ifndef __UTLT_MQ_H__
#define __UTLT_MQ_H__

#include <stdint.h>
#include <mqueue.h>

#include "utlt_debug.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef uintptr_t MQId;

/**
 * @return MQId or NULL on error.
 */
MQId MQCreate(int oflag);

Status MQDelete(MQId mqId);

long MQGetMsgSize(MQId mqId);

/**
 * @return STATUS_OK or STATUS_EAGAIN if the queue is full and the oflag O_NONBLOCK was set.
 */
Status MQSend(MQId mqId, const char *msg, int msgLen);

/**
 * @return STATUS_OK or STATUS_EAGAIN if the queue is empty and the oflag O_NONBLOCK was set.
 */
Status MQRecv(MQId mqId, char *msg, int msgLen);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* #ifndef __UTLT_MQ_H__ */
