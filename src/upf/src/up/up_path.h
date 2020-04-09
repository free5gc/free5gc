#ifndef __UP_PATH_H_
#define __UP_PATH_H_

#include "utlt_debug.h"
#include "utlt_network.h"
#include "upf_context.h"

Status UpRouteInit();
Status UpRouteTerminate();

Status GTPv1ServerInit();
Status GTPv1ServerTerminate();

Status GtpHandler(Sock *sock, void *data);

Status GtpHandleEchoRequest(Sock *sock, void *data);
Status GtpHandleEchoResponse(void *data);
Status GtpHandleEndMark(Sock *sock, void *data);

Status BufferServerInit();
Status BufferServerTerminate();
Status BufferHandler(Sock *sock, void *data);

Status UpSendPacketByPdrFar(UpfPdr *pdr, UpfFar *far, Sock *sock);

#endif /* __UP_PATH_H_ */
