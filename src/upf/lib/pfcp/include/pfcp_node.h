#ifndef __PFCP_NODE_H__
#define __PFCP_NODE_H__

#include "utlt_list.h"
#include "utlt_network.h"
#include "utlt_timer.h"

#include "pfcp_types.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _PfcpNode {
    ListNode        node;           /* List of node for PFCP */
    SockAddr        *saList;               /* Socket Address list */
    Sock            *sock;
    Ip              ip;

    ListNode        localList;
    ListNode        remoteList;

#define PFCP_NODE_ST_NULL           0
#define PFCP_NODE_ST_ASSOCIATED     1    
    uint8_t         state;          /* Association complete or not */
    TimerBlkID      timeHeartbeat;    /* no timer lib */
    PfcpNodeId      nodeId;
    union {
        uint8_t     upFunctionFeatures;
        uint8_t     cpFunctionFeatures;
    };
    PfcpUserPlaneIpResourceInformation  userPlaneInfo;
} PfcpNode;

Status PfcpNodeInit(void);
Status PfcpNodeTerminate(void);
Status PfcpAddNode(ListNode *list, PfcpNode **node,
                   const SockAddr *allList, _Bool noIpv4, _Bool noIpv6, _Bool preferIpv4);
PfcpNode *PfcpAddNodeWithSeid(ListNode *list, PfcpFSeid *fSeid,
        uint16_t port, _Bool noIpv4, _Bool noIpv6, _Bool preferIpv4);
Status PfcpRemoveNode(ListNode *list, PfcpNode *node);
Status PfcpRemoveAllNodes(ListNode *list);
PfcpNode *PfcpFindNode(ListNode *list, PfcpFSeid *fSeid);
PfcpNode *PfcpFindNodeSockAddr(ListNode *list, SockAddr *sock);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __PFCP_NODE_H__ */
