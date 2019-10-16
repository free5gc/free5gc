#ifndef __UPF_H__
#define __UPF_H__

#include "utlt_debug.h"
#include "utlt_thread.h"
#include "utlt_event.h"

Status UpfInit(char *configPath);
Status UpfTerminate();

Status UtltLibInit();
Status UtltLibTerminate();

#endif /* __UPF_H__ */
