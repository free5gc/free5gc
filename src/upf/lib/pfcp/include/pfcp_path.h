#ifndef __PFCP_PATH_H__
#define __PFCP_PATH_H__

#include "utlt_debug.h"
#include "utlt_3gppTypes.h"
#include "utlt_network.h"
#include "utlt_buff.h"

#include "pfcp_node.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

Status PfcpServer(SockNode *snode, SockHandler handler);
Status PfcpReceiveFrom(Sock *sock, Bufblk **bufBlk, SockAddr *from);
Status PfcpClient(PfcpNode *node);
Status PfcpConnect(Sock *sockv4, Sock *sockv6, PfcpNode *node);
Status PfcpServerList(ListNode *list, SockHandler handler, int epfd);
Sock *PfcpLocalSockFirst(ListNode *list);
SockAddr *PfcpLocalAddrFirst(ListNode *list);
Status PfcpReceive(Sock *sock, Bufblk **bufBlk);
Status PfcpReceiveFrom(Sock *sock, Bufblk **bufBlk, SockAddr *from);
Status PfcpSend(PfcpNode *node, Bufblk *bufBlk);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif