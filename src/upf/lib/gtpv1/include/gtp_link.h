#ifndef __GTP_LINK_H__
#define __GTP_LINK_H__

#include <netinet/in.h>

#include "utlt_lib.h"
#include "utlt_debug.h"
#include "utlt_list.h"
#include "utlt_network.h"

#define GTP_V1_PORT 2152
#define GTP_PRIME_PORT 3386

typedef struct {
    ListNode node;
    char ip[INET6_ADDRSTRLEN], ifname[MAX_IFNAME_STRLEN];
    Sock *sock1, *sockPrime;
} Gtpv1TunDevNode;

Status GtpLinkCreate(Gtpv1TunDevNode *node);
Status GtpLinkFree(Gtpv1TunDevNode *node);

Status GtpTunnelAdd(const char *ifname, int iteid, int oteid, const char *destIP, const char *tunIP);
Status GtpTunnelDel(const char *ifname, int iteid);
Status GtpTunnelList();

Status Gtpv1DevPoolInit();
Status Gtpv1DevPoolFinal();
Gtpv1TunDevNode *Gtpv1DevListAdd(ListNode *list, const char *ip, const char *ifname);
Status Gtpv1DevListFree(ListNode *list);

#endif /* __GTP_LINK_H__ */
