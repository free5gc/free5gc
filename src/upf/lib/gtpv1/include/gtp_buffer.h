#ifndef __GTP_BUFFER_H__
#define __GTP_BUFFER_H__

/* This header file is only for buffering functions.
 * Include or modify it if buffering would be maintained
 * in this process instead of other space, likes kernel
 */

#include "utlt_debug.h"
#include "utlt_network.h"
#include "utlt_buff.h"

#define MAX_OF_BUFFER_PACKET_SIZE 1600

Sock *BufferServerCreate(int type, const char *path, SockHandler handler, void *data);
Status BufferServerFree(Sock *sock);

int BufferRecv(Sock *sock, Bufblk *pktbuf, uint16_t *pdrId, uint8_t *farAction);

Status BufferEpollRegister(int epfd, Sock *sock);
Status BufferEpollDeregister(int epfd, Sock *sock);

#endif /* __GTP_BUFFER_H__ */