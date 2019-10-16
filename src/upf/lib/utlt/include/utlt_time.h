#ifndef __TIME_H__
#define __TIME_H__

#include <stdint.h>
#include <time.h>
#include <sys/time.h>

#include "utlt_debug.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

static const char month_snames[12][4] = {
    "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
};
static const char day_snames[7][4] = {
    "Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"
};

#define TIME_USE_GMT     0
#define TIME_USE_LOCAL   1

// number of microseconds since 00:00:00 january 1, 1970 UTC
typedef int64_t utime_t;

#define TIME_C(val) INT64_C(val)

#define USEC_PER_SEC TIME_C(1000000)

#define TimeSec(time) ((time) / USEC_PER_SEC)   // return as sec

#define TimeUsec(time) ((time) % USEC_PER_SEC)  // return as usec

#define TimeMsec(time) ((time) / 1000)  // return as msec

#define TimeMsecToUsec(msec) ((utime_t)(msec) * 1000)

#define TimeSecToUsec(sec) ((utime_t)(sec) * USEC_PER_SEC)

typedef struct _TimeTM {
    int32_t tm_usec;    // microseconds
    struct tm *tm_ptr;
} TimeTM;

// return current time in microseconds (usec)
utime_t TimeNow(void);

// convert time to human readable components in GMT / local timezone
Status TimeConvert(TimeTM *timeTM, utime_t ut, int timeType);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif  /* __TIME_H__ */