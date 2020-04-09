#ifndef __POOL_H__
#define __POOL_H__

#include <pthread.h>

#define PoolDeclare(__name, __type, __cap) \
    typedef struct { \
        int qFront, qEnd, qCap; \
        __type *queueAvail[__cap + 1], pool[__cap + 1];\
        pthread_mutex_t lock; \
    } pool##__name##_t; \
    pool##__name##_t __name

// The number of available space in this pool
#define PoolSize(__nameptr) \
    ((((__nameptr)->qEnd + (__nameptr)->qCap + 1 - (__nameptr)->qFront)) % ((__nameptr)->qCap + 1))

// Total space of this pool, including used and unused
#define PoolCap(__nameptr) ((__nameptr)->qCap)

#define PoolInit(__nameptr, __cap) do { \
    (__nameptr)->qFront = 0; \
    (__nameptr)->qEnd = __cap; \
    (__nameptr)->qCap = __cap; \
    for (int __i = 0; __i < __cap; __i++) \
        (__nameptr)->queueAvail[__i] = &((__nameptr)->pool[__i]); \
    pthread_mutex_init(&(__nameptr)->lock, 0); \
    UTLT_Trace("Pool Init Finish: %d", PoolSize(__nameptr)); \
} while(0)

#define PoolTerminate(__nameptr) pthread_mutex_destroy(&(__nameptr)->lock)

#define PoolAlloc(__nameptr, __assignedPtr) do { \
    pthread_mutex_lock(&(__nameptr)->lock); \
    if (PoolSize(__nameptr) > 0) { \
        (__assignedPtr) = (__nameptr)->queueAvail[(__nameptr)->qFront]; \
        (__nameptr)->qFront = ((__nameptr)->qFront + 1) % ((__nameptr)->qCap + 1); \
        UTLT_Debug("Pool alloc successful, total capacity[%d], available[%d]" \
        , PoolCap(__nameptr), PoolSize(__nameptr)); \
    } else { \
        (__assignedPtr) = NULL; \
        UTLT_Warning("Pool is empty"); \
    } \
    pthread_mutex_unlock(&(__nameptr)->lock); \
} while(0)

#define PoolFree(__nameptr, __assignedPtr) do { \
    pthread_mutex_lock(&(__nameptr)->lock); \
    if (PoolSize(__nameptr) < (__nameptr)->qCap) { \
        (__nameptr)->queueAvail[(__nameptr)->qEnd] = (__assignedPtr); \
        (__nameptr)->qEnd = ((__nameptr)->qEnd + 1) % ((__nameptr)->qCap + 1); \
        UTLT_Debug("Pool Free successful, total capacity[%d], available[%d]" \
        , PoolCap(__nameptr), PoolSize(__nameptr)); \
    } else { \
        UTLT_Error("Pool is full, it may not belong to this pool"); \
    } \
    pthread_mutex_unlock(&(__nameptr)->lock); \
} while(0)

#define PoolUsedCheck(__pname) (PoolCap(__pname) - PoolSize(__pname))

#define PoolAvailable(__pname) (PoolSize(__pname) > 0)

#endif /* __POOL_H__ */
