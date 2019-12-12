#ifndef __UP_GTP_PATH_H__
#define __UP_GTP_PATH_H__

#include "utlt_debug.h"
#include "utlt_buff.h"
#include "utlt_network.h"
#include "upf_context.h"

Status GtpRouteInit();
Status GtpRouteTerminate();

Status GTPv1ServerInit();
Status GTPv1ServerTerminate();

Status GtpHandler(Sock *sock, void *data);

Status GtpHandleEchoRequest(Sock *sock, void *data);
Status GtpHandleEchoResponse(void *data);
Status GtpHandleEndMark(Sock *sock, void *data);
Status GtpHandleTPDU(Sock *sock, Bufblk *data);

Status UpfSessionPacketSend(UpfSession *session, Sock *sock);
Status UpfSessionPacketRecv(UpfSession *session, Bufblk *pktBuf);
Status UpfSessionPacketClear(UpfSession *session);

#endif /* __UP_GTP_PATH_H__ */
