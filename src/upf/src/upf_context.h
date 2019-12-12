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

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _UpfUeIp     UpfUeIp;
typedef struct _UpfDev      UpfDev;
typedef struct _UpfFar      UpfFar;
typedef struct _UpfBar      UpfBar;
typedef struct _UpfQer      UpfQer;
typedef struct _UpfUrr      UpfUrr;

typedef enum _UpfEvent {

    UPF_EVENT_N4_MESSAGE,
    UPF_EVENT_SESSION_REPORT,
    UPF_EVENT_N4_T3_RESPONSE,
    UPF_EVENT_N4_T3_HOLDING,

    UPF_EVENT_TOP,

} UpfEvent;

typedef struct {
    const char      *gtpDevNamePrefix;   // Default : "upfgtp"

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

    // DNS
#define MAX_NUM_OF_DNS          2
    const char      *dns[MAX_NUM_OF_DNS];
    const char      *dns6[MAX_NUM_OF_DNS];

    // Add other context here
    ListNode        ranS1uList;         // RAN List connected to UPF
    ListNode        upfN4List;          // UPF PFCP Node List
    ListNode        apnList;

    // Different list of policy rule
    ListNode        pdrList;
    ListNode        farList;
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
    ListNode        dlPdrList;
    ListNode        ulPdrList;

    /* Buff the un-tunnel packet */
#define MAX_NUM_OF_PACKET_BUFFER_SIZE 0xff
    int             pktBufIdx;
    Bufblk          *packetBuffer[MAX_NUM_OF_PACKET_BUFFER_SIZE];
    pthread_mutex_t bufLock;
} UpfSession;

typedef struct _UpfPdr {
    ListNode        node;               // Node List
    int             index;

    uint32_t        upfGtpUTeid;
    uint8_t         ulDl;               // UL(0) or DL(1) PDR
    uint16_t        pdrId;

#define SMF_TEID_IP_DESC_IPV4   1
#define SMF_TEID_IP_DESC_IPV6   2
    // Upf Ip (with upfGtpUTeid conbine to F-TEID)
    union {
        /* IPV4 */
        struct in_addr      addr4;
        /* IPV6 */
        struct in6_addr     addr6;
        /* BOTH */
        struct {
            struct in_addr  addr4;
            struct in6_addr addr6;
        } dualStack;
    };

    UpfUeIp         ueIp;

    uint32_t        precedence;
    uint8_t         outerHeaderRemove;
    uint8_t         sourceInterface;

    UpfFar          *far;
    UpfQer          *qer;
    UpfUrr          *urr;

    PfcpNode        *pfcpNode;
    UpfSession      *session;
} UpfPdr;

typedef struct _UpfFar {
    ListNode        node;
    int             index;

    uint32_t        farId;
    uint8_t         applyAction;
    uint8_t         destinationInterface;

    uint16_t        referenceCount; // for reported usage
    uint8_t         created;

    uint32_t        upfN3Teid;
    Ip              ranIp;

    UpfBar          *bar;
    PfcpNode        *pfcpNode;
    //SockNode        *gtpNode; // TODO: check if can used
} UpfFar;

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
UpfPdr *UpfPdrAdd(UpfSession *session);
Status UpfPdrRemove(UpfPdr *pdr);
UpfPdr *UpfPdrFindByPdrId(uint16_t pdrId);
UpfPdr *UpfPdrFindByFarId(uint32_t farId);
UpfPdr *UpfPdrFidByUpfGtpUTeid(uint32_t teid);
UpfFar *UpfFarAdd();
Status UpfFarRemove(UpfFar *far);
UpfFar *UpfFarFindByFarId(uint32_t farId);

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
