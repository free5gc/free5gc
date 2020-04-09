#define TRACE_MODULE _upf_context

#include "upf_context.h"

#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <pthread.h>
#include <netinet/in.h>
#include <net/if.h>

#include "utlt_debug.h"
#include "utlt_index.h"
#include "utlt_hash.h"
#include "utlt_network.h"
#include "gtp_header.h"
#include "gtp_link.h"
#include "upf.h"
#include "pfcp_message.h"
#include "pfcp_types.h"
#include "pfcp_xact.h"
#include "libgtp5gnl/gtp5g.h"
#include "libgtp5gnl/gtp5gnl.h"
#include "gtp_tunnel.h"


#define MAX_POOL_OF_PDRID (MAX_POOL_OF_BEARER * 2)
#define MAX_POOL_OF_QER (MAX_POOL_OF_SESS * 2)
#define MAX_POOL_OF_URR (MAX_POOL_OF_UE)
#define MAX_POOL_OF_BAR (MAX_POOL_OF_UE)

#define MAX_NUM_OF_SUBNET       16

IndexDeclare(upfPdrIdPool, UpfPdrId, MAX_POOL_OF_PDRID);
IndexDeclare(upfSessionPool, UpfSession, MAX_POOL_OF_SESS);
IndexDeclare(upfQerPool, UpfQer, MAX_POOL_OF_QER);
IndexDeclare(upfUrrPool, UpfUrr, MAX_POOL_OF_URR);
IndexDeclare(upfBarPool, UpfBar, MAX_POOL_OF_BAR);

static UpfContext self;
static _Bool upfContextInitialized = 0;

UpfContext *Self() {
    return &self;
}

Status UpfContextInit() {
    UTLT_Assert(upfContextInitialized == 0, return STATUS_ERROR,
                "UPF context has been initialized!");

    memset(&self, 0, sizeof(UpfContext));

    // TODO : Add GTPv1 init here
    ListInit(&self.gtpv1DevList);
    ListInit(&self.gtpv1v6DevList);

    // TODO : Add PFCP init here
    //ListInit(&self.pfcpIPList);
    //ListInit(&self.pfcpIPv6List);

    // TODO : Add by self if context has been updated
    // TODO: check if gtp node need to init?
    //ListInit(&self.gtpv1DevList);
    //ListInit(&self.gtpv1v6DevList);
    //ListInit(&self.pfcpIPList);
    //ListInit(&self.pfcpIPv6List);
    ListInit(&self.ranS1uList);
    ListInit(&self.upfN4List);
    ListInit(&self.apnList);
    ListInit(&self.qerList);
    ListInit(&self.urrList);

    self.recoveryTime = htonl(time((time_t *)NULL));

    // Set Default Value
    self.gtpDevNamePrefix = "upfgtp";
    // defined in utlt_3gpptypes instead of GTP_V1_PORT defined in GTP_PATH;
    self.gtpv1Port = GTPV1_U_UDP_PORT;
    self.pfcpPort = PFCP_UDP_PORT;
    self.gtpv1DevSN = 0;

    // Init Resource
    IndexInit(&upfSessionPool, MAX_POOL_OF_SESS);
    IndexInit(&upfQerPool, MAX_POOL_OF_QER);
    IndexInit(&upfPdrIdPool, MAX_POOL_OF_PDRID);
    IndexInit(&upfUrrPool, MAX_POOL_OF_URR);
    IndexInit(&upfBarPool, MAX_POOL_OF_BAR);

    PfcpNodeInit(); // init pfcp node for upfN4List (it will used pfcp node)
    TimerListInit(&self.timerServiceList);

    // TODO: Read from config
    strncpy(self.buffSockPath, "/tmp/free5gc_unix_sock", MAX_SOCK_PATH_LEN);
    self.sessionHash = HashMake();
    self.bufPacketHash = HashMake();
    // spin lock protect write data instead of mutex protect code block
    int ret = pthread_spin_init(&self.buffLock, PTHREAD_PROCESS_PRIVATE);
    UTLT_Assert(ret == 0, , "buffLock cannot create: %s", strerror(ret));

    upfContextInitialized = 1;

    return STATUS_OK;
}

// TODO : Need to Remove List Members iterativelyatively
Status UpfContextTerminate() {
    UTLT_Assert(upfContextInitialized == 1, return STATUS_ERROR,
                "UPF context has been terminated!");

    Status status = STATUS_OK;

    int ret = pthread_spin_destroy(&self.buffLock);
    UTLT_Assert(ret == 0, , "buffLock cannot destroy: %s", strerror(ret));
    UTLT_Assert(self.bufPacketHash, , "Buffer Hash Table missing?!");
    HashDestroy(self.bufPacketHash);

    UpfSessionRemoveAll();
    UTLT_Assert(self.sessionHash, , "Session Hash Table missing?!");
    HashDestroy(self.sessionHash);

    // Terminate resource
    IndexTerminate(&upfBarPool);
    IndexTerminate(&upfUrrPool);
    IndexTerminate(&upfQerPool);
    IndexTerminate(&upfPdrIdPool);
    IndexTerminate(&upfSessionPool);

    PfcpRemoveAllNodes(&self.upfN4List);
    PfcpNodeTerminate();

    // TODO: remove gtpv1TunnelList, ranS1uList, upfN4LIst, apnList,
    // pdrList, farList, qerList, urrLIist
    Gtpv1DevListFree(&self.gtpv1DevList);
    Gtpv1DevListFree(&self.gtpv1v6DevList);
    SockNodeListFree(&self.pfcpIPList);
    SockNodeListFree(&self.pfcpIPv6List);

    //UpfBufPacketRemoveAll();
    UpfApnRemoveAll();

    upfContextInitialized = 0;

    return status;
}

UpfPdrId *UpfPdrIdAdd(uint16_t pdrId) {
    UTLT_Assert(pdrId, return NULL, "PDR ID cannot be 0");
    UpfPdrId *pdrIdPtr;

    IndexAlloc(&upfPdrIdPool, pdrIdPtr);
    pdrIdPtr->pdrId = pdrId;

    return pdrIdPtr;
}

Status UpfPdrIdRemove(UpfPdrId *pdrIdPtr) {
    UTLT_Assert(pdrIdPtr, return STATUS_ERROR, "PDR error");

    IndexFree(&upfPdrIdPool, pdrIdPtr);

    return STATUS_OK;
}

HashIndex *UpfBufPacketFirst() {
    UTLT_Assert(self.bufPacketHash, return NULL, "");
    return HashFirst(self.bufPacketHash);
}

HashIndex *UpfBufPacketNext(HashIndex *hashIdx) {
    UTLT_Assert(hashIdx, return NULL, "");
    return HashNext(hashIdx);
}

UpfBufPacket *UpfBufPacketThis(HashIndex *hashIdx) {
    UTLT_Assert(hashIdx, return NULL, "");
    return (UpfBufPacket *)HashThisKey(hashIdx);
}

UpfBufPacket *UpfBufPacketFindByPdrId(uint16_t pdrId) {
    return (UpfBufPacket*)HashGet(self.bufPacketHash,
                                  &pdrId, sizeof(uint16_t));
}

UpfBufPacket *UpfBufPacketAdd(const UpfSession * const session,
                              const uint16_t pdrId) {
    UTLT_Assert(session, return NULL, "No session");
    UTLT_Assert(pdrId, return NULL, "PDR ID cannot be 0");

    UpfBufPacket *newBufPacket = UTLT_Malloc(sizeof(UpfBufPacket));
    UTLT_Assert(newBufPacket, return NULL, "Allocate new slot error");
    newBufPacket->sessionPtr = session;
    newBufPacket->pdrId = pdrId;
    newBufPacket->packetBuffer = NULL;

    HashSet(self.bufPacketHash, &newBufPacket->pdrId,
            sizeof(uint16_t), newBufPacket);

    //ListAppend(&Self()->bufPacketList, newBufPacket);
    return newBufPacket;
}

Status UpfBufPacketRemove(UpfBufPacket *bufPacket) {
    UTLT_Assert(bufPacket, return STATUS_ERROR,
                "Input bufPacket error");
    Status status;

    bufPacket->sessionPtr = NULL;
    bufPacket->pdrId = 0;
    if (bufPacket->packetBuffer) {
        status = BufblkFree(bufPacket->packetBuffer);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "packet in bufPacket free error");
    }

    HashSet(self.bufPacketHash, &bufPacket->pdrId,
            sizeof(uint16_t), NULL);
    //ListRemove(&Self()->bufPacketList, bufPacket);
    status = UTLT_Free(bufPacket);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "bufPacket free error");

    return STATUS_OK;
}

Status UpfBufPacketRemoveAll() {
    HashIndex *hashIdx = NULL;
    UpfBufPacket *bufPacket = NULL;

    for (hashIdx = UpfBufPacketFirst(); hashIdx;
         hashIdx = UpfBufPacketNext(hashIdx)) {
        bufPacket = UpfBufPacketThis(hashIdx);
        UpfBufPacketRemove(bufPacket);
    }
    // List version
    //UpfBufPdr *node, *nextNode;
    //
    //node = ListFirst(&self.bufPacketList);
    //while (node) {
    //  nextNode = (UpfBufPacket *)ListNext(node);
    //  UpfBufPacketRemove(node);
    //  node = nextNode;
    //}

    return STATUS_OK;
}

/**
 * @param  *natifname: nullable
 */
ApnNode *UpfApnAdd(const char *apnName, const char *ip,
                   const char *prefix, const char *natifname) {
    UTLT_Assert(strlen(apnName) <= MAX_APN_LEN, return NULL,
                "apn name should not longer than %d", MAX_APN_LEN);

    ApnNode *newApnNode = UTLT_Malloc(sizeof(ApnNode));
    strcpy(newApnNode->apn, apnName);
    strcpy(newApnNode->subnetIP, ip);
    newApnNode->subnetPrefix = atoi(prefix);
    if (natifname)
        strcpy(newApnNode->natifname, natifname);

    ListAppend(&Self()->apnList, newApnNode);
    return newApnNode;
}

Status UpfApnRemoveAll() {
    ApnNode *node, *nextNode;

    node = ListFirst(&self.apnList);
    while (node) {
        nextNode = (ApnNode *)ListNext(node);
        ListRemove(&self.apnList, node);
        UTLT_Free(node);
        node = nextNode;
    }

    return STATUS_OK;
}

HashIndex *UpfSessionFirst() {
    UTLT_Assert(self.sessionHash, return NULL, "");
    return HashFirst(self.sessionHash);
}

HashIndex *UpfSessionNext(HashIndex *hashIdx) {
    UTLT_Assert(hashIdx, return NULL, "");
    return HashNext(hashIdx);
}

UpfSession *UpfSessionThis(HashIndex *hashIdx) {
    UTLT_Assert(hashIdx, return NULL, "");
    return (UpfSession *)HashThisKey(hashIdx);
}

void SessionHashKeygen(uint8_t *out, int *outLen, uint8_t *imsi,
                       int imsiLen, uint8_t *apn) {
    memcpy(out, imsi, imsiLen);
    strncpy((char *)(out + imsiLen), (char*)apn, MAX_APN_LEN + 1);
    *outLen = imsiLen + strlen((char *)(out + imsiLen));

    return;
}

UpfSession *UpfSessionAdd(PfcpUeIpAddr *ueIp, uint8_t *apn,
                          uint8_t pdnType) {
    UpfSession *session = NULL;

    IndexAlloc(&upfSessionPool, session);
    UTLT_Assert(session, return NULL, "session alloc error");

    //session->gtpNode = NULL;

    if (self.pfcpAddr) {
        session->upfSeid =
          ((uint64_t)self.pfcpAddr->s4.sin_addr.s_addr << 32)
          | session->index;
    } else if (self.pfcpAddr6) {
        uint32_t *ptr =
          (uint32_t *)self.pfcpAddr6->s6.sin6_addr.s6_addr;
        session->upfSeid =
          (((uint64_t)(*ptr)) << 32) | session->index;
        // TODO: check if correct
    }
    session->upfSeid = htobe64(session->upfSeid);
    //UTLT_Info()
    session->upfSeid = 0; // TODO: check why

    /* IMSI APN Hash */
    /* APN */
    strncpy((char*)session->pdn.apn, (char*)apn, MAX_APN_LEN + 1);

    ListInit(&session->pdrIdList);

    session->pdn.paa.pdnType = pdnType;
    if (pdnType == PFCP_PDN_TYPE_IPV4) {
        session->ueIpv4.addr4 = ueIp->addr4;
        //session->pdn.paa.addr4 = ueIp->addr4;
    } else if (pdnType == PFCP_PDN_TYPE_IPV6) {
        session->ueIpv6.addr6 = ueIp->addr6;
        //session->pdn.paa.addr6 = ueIp->addr6;
    } else if (pdnType == PFCP_PDN_TYPE_IPV4V6) {
        // TODO
        // session->ueIpv4 = UpfUeIPAlloc(AF_INET, apn);
        // UTLT_Assert(session->ueIpv4,
        //   UpfSessionRemove(session); return NULL,
        //   "Cannot allocate IPv4");

        // session->ueIpv6 = UpfUeIPAlloc(AF_INET6, apn);
        // UTLT_Assert(session->ueIpv6,
        //   UpfSessionRemove(session); return NULL,
        //   "Cannot allocate IPv6");

        // session->pdn.paa.dualStack.addr4 = session->ueIpv4->addr4;
        // session->pdn.paa.dualStack.addr6 = session->ueIpv6->addr6;
    } else {
        UTLT_Assert(0, return NULL, "UnSupported PDN Type(%d)", pdnType);
    }

    /* Generate Hash Key: IP + APN */
    if (pdnType == PFCP_PDN_TYPE_IPV4) {
        SessionHashKeygen(session->hashKey,
                          &session->hashKeylen,
                          (uint8_t *)&session->ueIpv4.addr4, 4, apn);
    } else {
        SessionHashKeygen(session->hashKey,
                          &session->hashKeylen,
                          (uint8_t *)&session->ueIpv6.addr6,
                          IPV6_LEN, apn);
    }

    HashSet(self.sessionHash, session->hashKey,
            session->hashKeylen, session);

    return session;
}

Status UpfSessionRemove(UpfSession *session) {
    UTLT_Assert(self.sessionHash, return STATUS_ERROR,
                "sessionHash error");
    UTLT_Assert(session, return STATUS_ERROR, "session error");

    HashSet(self.sessionHash, session->hashKey,
            session->hashKeylen, NULL);

    // if (session->ueIpv4) {
    //     UpfUeIPFree(session->ueIpv4);
    // }
    // if (session->ueIpv6) {
    //     UpfUeIPFree(session->ueIpv6);
    // }

    UpfPdrId *pdrIdPtr = ListFirst(&session->pdrIdList);
    while (pdrIdPtr) {
        Gtpv1TunDevNode *gtpv1Dev4 =
          (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
        UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
        Status status = GtpTunnelDelPdr(gtpv1Dev4->ifname, pdrIdPtr->pdrId);
        UTLT_Assert(status == STATUS_OK, ,
                    "Remove PDR[%u] failed", pdrIdPtr->pdrId);
        // TODO: remove FAR of PDR if need
        ListRemove(&session->pdrIdList, pdrIdPtr);
        UTLT_Assert(UpfPdrIdRemove(pdrIdPtr) == STATUS_OK, , "Pdr id remove error");
        pdrIdPtr = (UpfPdrId *)ListFirst(&session->pdrIdList);
    }

    IndexFree(&upfSessionPool, session);

    return STATUS_OK;
}

Status UpfSessionRemoveAll() {
    HashIndex *hashIdx = NULL;
    UpfSession *session = NULL;

    for (hashIdx = UpfSessionFirst(); hashIdx;
         hashIdx = UpfSessionNext(hashIdx)) {
        session = UpfSessionThis(hashIdx);
        UpfSessionRemove(session);
    }

    return STATUS_OK;
}

UpfSession *UpfSessionFind(uint32_t idx) {
    //UTLT_Assert(idx, return NULL, "index error");
    return IndexFind(&upfSessionPool, idx);
}

UpfSession *UpfSessionFindBySeid(uint64_t seid) {
    return UpfSessionFind((seid-1) & 0xFFFFFFFF);
}

UpfSession *UpfSessionAddByMessage(PfcpMessage *message) {
    UpfSession *session;

    PFCPSessionEstablishmentRequest *request =
      &message->pFCPSessionEstablishmentRequest;

    if (!request->nodeID.presence) {
        UTLT_Error("no NodeID");
        return NULL;
    }
    if (!request->cPFSEID.presence) {
        UTLT_Error("No cp F-SEID");
        return NULL;
    }
    if (!request->createPDR[0].presence) {
        UTLT_Error("No PDR");
        return NULL;
    }
    if (!request->createFAR[0].presence) {
        UTLT_Error("No FAR");
        return NULL;
    }
    if (!request->pDNType.presence) {
        UTLT_Error("No PDN Type");
        return NULL;
    }
    if (!request->createPDR[0].pDI.presence) {
        UTLT_Error("PDR PDI error");
        return NULL;
    }
    if (!request->createPDR[0].pDI.uEIPAddress.presence) {
        UTLT_Error("UE IP Address error");
        return NULL;
    }
    if (!request->createPDR[0].pDI.networkInstance.presence) {
        UTLT_Error("Interface error");
        return NULL;
    }

    session = UpfSessionAdd((PfcpUeIpAddr *)
                &request->createPDR[0].pDI.uEIPAddress.value,
                request->createPDR[0].pDI.networkInstance.value,
                ((int8_t *)request->pDNType.value)[0]);
    UTLT_Assert(session, return NULL, "session add error");

    session->smfSeid = *(uint64_t*)request->cPFSEID.value;
    session->upfSeid = session->index+1;
    UTLT_Trace("UPF Establishment UPF SEID: %lu", session->upfSeid);

    return session;
}

