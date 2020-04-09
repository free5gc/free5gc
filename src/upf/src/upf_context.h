#ifndef __UPF_CONTEXT_H__
#define __UPF_CONTEXT_H__

#include <stdint.h>
#include <netinet/in.h>
#include <net/if.h>
#include <pthread.h>

#include "utlt_list.h"
#include "utlt_buff.h"
#include "utlt_event.h"
#include "utlt_thread.h"
#include "utlt_network.h"
#include "utlt_hash.h"
#include "utlt_3gppTypes.h"
#include "utlt_timer.h"

#include "pfcp_node.h"
#include "gtp_path.h"
#include "pfcp_message.h"
#include "libgtp5gnl/gtp5g.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _UpfUeIp      UpfUeIp;
typedef struct _UpfDev       UpfDev;
typedef struct gtp5g_pdr     UpfPdr;
typedef struct gtp5g_far     UpfFar;
typedef struct _UpfBufPacket UpfBufPacket;
typedef struct _UpfPdrId     UpfPdrId;
typedef struct _UpfBar       UpfBar;
typedef struct _UpfQer       UpfQer;
typedef struct _UpfUrr       UpfUrr;

typedef enum _UpfEvent {

    UPF_EVENT_N4_MESSAGE,
    UPF_EVENT_SESSION_REPORT,
    UPF_EVENT_N4_T3_RESPONSE,
    UPF_EVENT_N4_T3_HOLDING,

    UPF_EVENT_TOP,

} UpfEvent;

typedef struct {
    const char      *gtpDevNamePrefix;   // Default : "upfgtp"

    ListNode        gtpInterfaceList;    // name of interface (char*)
    // Add context related to GTP-U here
    uint16_t        gtpv1Port;           // Default : GTP_V1_PORT
    int             gtpv1DevSN;          // Serial number for naming gtpv1Dev, gtpv1v6Dev
    ListNode        gtpv1DevList;        // GTPv1 IPv4 Device Socket (Gtpv1TunDevNode)
    ListNode        gtpv1v6DevList;      // GTPv1 IPv6 Device Socket (Gtpv1TunDevNode)
    SockAddr        *gtpv1Addr;          // GTPv1 IPv4 Address
    SockAddr        *gtpv1Addr6;         // GTPv1 IPv6 Address

    // Add context related to PFCP here
    uint16_t        pfcpPort;            // Default : PFCP_PORT
    ListNode        pfcpIPList;          // PFCP IPv4 Server List (SockNode)
    ListNode        pfcpIPv6List;        // PFCP IPv6 Server List (SockNode)
    Sock            *pfcpSock;           // IPv4 Socket
    Sock            *pfcpSock6;          // IPv6 Socket
    SockAddr        *pfcpAddr;           // IPv4 Address
    SockAddr        *pfcpAddr6;          // IPv6 Address

    /* Use Array or Hash for better performance
     * Because max size of the list is 65536 due to the max of PDR ID
     * We can use array for O(1) search instead of O(N) search in list
     * Trade off of speed and memory size
     */
    //ListNode        bufPacketList;       // save pdrId and buffer here

    // DNS
#define MAX_NUM_OF_DNS          2
    const char      *dns[MAX_NUM_OF_DNS];
    const char      *dns6[MAX_NUM_OF_DNS];

    // Add other context here
    ListNode        ranS1uList;         // RAN List connected to UPF
    ListNode        upfN4List;          // UPF PFCP Node List
    ListNode        apnList;

    // Different list of policy rule
    // TODO: if implementing QER in kernel, remove these list
    ListNode        qerList;
    ListNode        urrList;

    uint32_t        recoveryTime;       // UTC time
    TimerList       timerServiceList;

    // Add some self library structure here
    int             epfd;               // Epoll fd
    EvtQId          eventQ;             // Event queue communicate between UP and CP
    ThreadID        pktRecvThread;      // Receive packet thread

    // Session : hash(IMSI+APN)
    Hash            *sessionHash;
    // Save buffer packet here
    Hash            *bufPacketHash;
    // Use spin lock to protect data write
    pthread_spinlock_t buffLock;
    // TODO: read from config
    // no reason, just want to bigger than /tmp/free5gc_unix_sock
#define MAX_SOCK_PATH_LEN 64
    char            buffSockPath[MAX_SOCK_PATH_LEN];
    // Buffering socket for recv packet from kernel
    Sock            *buffSock;


    // Config file
    const char      *configFilePath;
} UpfContext;

typedef struct _UpfUeIp {
    union {
        struct in_addr addr4;
        struct in6_addr addr6;
    };
} UpfUeIp;

typedef struct _UpfSession {
    int             index;

    uint64_t        upfSeid;
    uint64_t        smfSeid;

    /* APN Config */
    Pdn             pdn;
    UpfUeIp         ueIpv4;
    UpfUeIp         ueIpv6;

    /* User location */
    Tai             tai;
    //ECgi          eCgi; // For LTE E-UTRA Cell ID
    //NCgi          nCgi; // For 5GC NR Cell ID

    /* Hashed key: hash(IMSI+APN) */
    uint8_t         hashKey[MAX_IMSI_LEN+MAX_APN_LEN];
    int             hashKeylen;

    /* GTP, PFCP context */
    //SockNode        *gtpNode;
    PfcpNode        *pfcpNode;
    ListNode        pdrIdList;

} UpfSession;

// Used for buffering, Index type for each PDR
typedef struct _UpfBufPacket {
    //ListNode        node;
    int             index;

    // If sessionPtr == NULL, this PDR don't exist
    // TS 29.244 5.2.1 shows that PDR won't cross session
    const UpfSession *sessionPtr;
    uint16_t        pdrId;
    Bufblk          *packetBuffer;
} UpfBufPakcet;

typedef struct _UpfPdrId {
    ListNode        node;
    int             index;

    uint16_t        pdrId;
} UpfPdrId;

typedef struct _UpfUrr {
    ListNode        node;
    int             index;

    uint16_t        urrId;
    uint16_t        referenceCount;

    PfcpNode        *pfcpNode;
} UpfUrr;

typedef struct _UpfQer {
    ListNode        node;
    int             index;

    uint16_t        qerId;
    uint16_t        referenceCount;

    PfcpNode        *pfcpNode;
} UpfQer;

typedef struct _UpfBar {
    ListNode        node;
    int             index;

    uint16_t        barId;
    uint16_t        referenceId;

    PfcpNode        *pfcpNode;
} UpfBar;

typedef struct _ApnNode {
    ListNode      node;
    char          apn[MAX_APN_LEN + 1];
    char          subnetIP[INET6_ADDRSTRLEN];
    char          natifname[IF_NAMESIZE];
    uint8_t       subnetPrefix;
} ApnNode;

UpfContext *Self();
Status UpfContextInit();
Status UpfContextTerminate();

// APN / PDR / FAR
ApnNode *UpfApnAdd(const char *apnName, const char *ip, const char *prefix, const char *natifname);
Status UpfApnRemoveAll();

// BufPacket
HashIndex *UpfBufPacketFirst();
HashIndex *UpfBufPacketNext(HashIndex *hashIdx);
UpfPdrId *UpfPdrIdAdd(uint16_t pdrId);
Status UpfPdrIdRemove(UpfPdrId *pdrIdPtr);
UpfBufPacket *UpfBufPacketThis(HashIndex *hashIdx);
UpfBufPacket *UpfBufPacketFindByPdrId(uint16_t pdrId);
UpfBufPacket *UpfBufPacketAdd(const UpfSession * const session,
                              const uint16_t pdrId);
Status UpfBufPacketRemove(UpfBufPacket *bufPacket);
Status UpfBufPacketRemoveAll();
// Session
HashIndex *UpfSessionFirst();
HashIndex *UpfSessionNext(HashIndex *hashIdx);
UpfSession *UpfSessionThis(HashIndex *hashIdx);
void SessionHashKeygen(uint8_t *out, int *outLen, uint8_t *imsi, int imsiLen, uint8_t *apn);
UpfSession *UpfSessionAdd(PfcpUeIpAddr *ueIp, uint8_t *apn, uint8_t pdnType);
Status UpfSessionRemove(UpfSession *session);
Status UpfSessionRemoveAll();
UpfSession *UpfSessionFind(uint32_t idx);
UpfSession *UpfSessionFindBySeid(uint64_t seid);
UpfSession *UpfSessionAddByMessage(PfcpMessage *message);
UpfSession *UpfSessionFindByPdrTeid(uint32_t teid);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __UPF_CONTEXT_H__ */
