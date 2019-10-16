#ifndef __THREAD_H__
#define __THREAD_H__

#include <stdint.h>
#include <pthread.h>
#include <fcntl.h>      /* For O_* constants */
#include <sys/stat.h>   /* For mode constants */
#include <semaphore.h>

#include "utlt_debug.h"
#include "utlt_pool.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

#define MAX_NUM_OF_THREAD       128

typedef uintptr_t ThreadID;

typedef void (*ThreadFuncType)(ThreadID id, void*);

typedef struct _Thread {
    pthread_t tid;
    ThreadFuncType func;
    void *data;
    sem_t *semaphore;
} Thread;

int ThreadStop();
Status ThreadInit();
Status ThreadFinal();
Status ThreadCreate(ThreadID *id, ThreadFuncType func, void *data);
Status ThreadDelete(ThreadID id);  // Delete a thread of execution
Status ThreadJoin(ThreadID id);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif  /* __THREAD_H__ */