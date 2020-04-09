#ifndef __GTP_PATH_H__
#define __GTP_PATH_H__

#include "utlt_debug.h"
#include "utlt_list.h"
#include "utlt_buff.h"
#include "utlt_network.h"
#include "gtp_link.h"

#define MAX_OF_GTPV1_PACKET_SIZE 8192

Sock *GtpServerCreate(int domain, const char *addr, int port, SockHandler handler, void *data);
Status GtpServerFree(Sock *sock);

Sock *GtpClientCreate(int domain, const char *addr, int port);
Status GtpClientFree(Sock *sock);

Status GtpTunCreate(Gtpv1TunDevNode *node, SockHandler handler, void *data);
Status GtpTunFree(Gtpv1TunDevNode *node);

int GtpRecv(Sock *sock, Bufblk *pktbuf);
Status GtpSend(Sock *sock, Bufblk *pktbuf);

Status GtpEpollRegister(int epfd, Sock *sock);
Status GtpEpollDeregister(int epfd, Sock *sock);

Status GtpServerListCreate(int epfd, int domain, ListNode *sockList, int port, SockHandler handler, void *data);
Status GtpServerListFree(int epfd, ListNode *sockList);

Status GtpDevListCreate(int epfd, int domain, ListNode *sockList, SockHandler handler, void *data);
Status GtpDevListFree(int epfd, ListNode *sockList);

SockNode *GtpFindSockNodeByIp(ListNode *list, Ip *ip);
SockNode *GtpAddSockNodeWithIp(ListNode *list, Ip *ip, int port);

Status GtpBuildEchoRequest(Bufblk *pktbuf, int teid, int seq);
Status GtpBuildEndMark(Bufblk *pktbuf, int teid);

#endif /* __GTP_PATH_H__ */
