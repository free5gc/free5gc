#define TRACE_MODULE _upf_context

#include "upf_context.h"

#include <string.h>
#include <stdlib.h>
#include <arpa/inet.h>
#include <pthread.h>
#include <netinet/in.h>

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

#define MAX_POOL_OF_PDR (MAX_POOL_OF_BEARER * 2)
#define MAX_POOL_OF_FAR (MAX_POOL_OF_SESS * 2)
#define MAX_POOL_OF_QER (MAX_POOL_OF_SESS * 2)
#define MAX_POOL_OF_URR (MAX_POOL_OF_UE)
#define MAX_POOL_OF_BAR (MAX_POOL_OF_UE)

#define MAX_NUM_OF_SUBNET       16

IndexDeclare(upfSessionPool, UpfSession, MAX_POOL_OF_SESS);
IndexDeclare(upfPdrPool, UpfPdr, MAX_POOL_OF_PDR);
IndexDeclare(upfFarPool, UpfFar, MAX_POOL_OF_FAR);
IndexDeclare(upfQerPool, UpfQer, MAX_POOL_OF_QER);
IndexDeclare(upfUrrPool, UpfUrr, MAX_POOL_OF_URR);
IndexDeclare(upfBarPool, UpfBar, MAX_POOL_OF_BAR);

static UpfContext self;
static _Bool upfContextInitialized = 0;

UpfContext *Self() {
    return &self;
}

Status UpfContextInit() {
    UTLT_Assert(upfContextInitialized == 0, return STATUS_ERROR, "UPF context has been initialized!");

    memset(&self, 0, sizeof(UpfContext));

    // TODO : Add GTPv1 init here
    ListInit(&self.gtpv1DevList);
    ListInit(&self.gtpv1v6DevList);

    // TODO : Add PFCP init here
    ListInit(&self.pfcpIPList);
    ListInit(&self.pfcpIPv6List);

    // TODO : Add by self if context has been updated
    // TODO: check if gtp node need to init?
    ListInit(&self.gtpv1DevList);
    ListInit(&self.gtpv1v6DevList);
    ListInit(&self.pfcpIPList);
    ListInit(&self.pfcpIPv6List);
    ListInit(&self.ranS1uList);
    ListInit(&self.upfN4List);
    ListInit(&self.apnList);
    ListInit(&self.pdrList);
    ListInit(&self.farList);
    ListInit(&self.qerList);
    ListInit(&self.urrList);

    self.recoveryTime = htonl(time((time_t *)NULL));

    // Set Default Value
    self.gtpDevNamePrefix = "upfgtp";
    self.gtpv1Port = GTPV1_U_UDP_PORT; // defined in utlt_3gpptypes instead of GTP_V1_PORT defined in GTP_PATH;
    self.pfcpPort = PFCP_UDP_PORT;
    self.gtpv1DevSN = 0;

    // Init Resource
    IndexInit(&upfSessionPool, MAX_POOL_OF_SESS);
    IndexInit(&upfPdrPool, MAX_POOL_OF_PDR);
    IndexInit(&upfFarPool, MAX_POOL_OF_FAR);
    IndexInit(&upfQerPool, MAX_POOL_OF_QER);
    IndexInit(&upfUrrPool, MAX_POOL_OF_URR);
    IndexInit(&upfBarPool, MAX_POOL_OF_BAR);

    PfcpNodeInit(); // init pfcp node for upfN4List (it will used pfcp node)
    TimerListInit(&self.timerServiceList);

    self.sessionHash = HashMake();

    upfContextInitialized = 1;

    return STATUS_OK;
}

// TODO : Need to Remove List Members iterativelyatively
Status UpfContextTerminate() {
    UTLT_Assert(upfContextInitialized == 1, return STATUS_ERROR, "UPF context has been terminated!");

    Status status = STATUS_OK;

    UpfSessionRemoveAll();

    UTLT_Assert(self.sessionHash, , "Hash Table missing?!");
    HashDestroy(self.sessionHash);

    // Terminate resource
    IndexTerminate(&upfBarPool);
    IndexTerminate(&upfUrrPool);
    IndexTerminate(&upfQerPool);
    IndexTerminate(&upfFarPool);
    IndexTerminate(&upfPdrPool);
    IndexTerminate(&upfSessionPool);

    PfcpRemoveAllNodes(&self.upfN4List);
    PfcpNodeTerminate();

    // // TODO: remove gtpv1TunnelList, ranS1uList, upfN4LIst, apnList, pdrList, farList, qerList, urrLIist
    Gtpv1DevListFree(&self.gtpv1DevList);
    Gtpv1DevListFree(&self.gtpv1v6DevList);
    SockNodeListFree(&self.pfcpIPList);
    SockNodeListFree(&self.pfcpIPv6List);

    UpfApnRemoveAll();

    upfContextInitialized = 0;

    return status;
}

/**
 * @param  *natifname: nullable
 */
ApnNode *UpfApnAdd(const char *apnName, const char *ip, const char *prefix, const char *natifname) {
    UTLT_Assert(strlen(apnName) <= MAX_APN_LEN, return NULL, "apn name should not longer than %d", MAX_APN_LEN);

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
        node = nextNode;
    }

    return STATUS_OK;
}

// TODO: check this function
UpfPdr *UpfPdrAdd(UpfSession *session) {
    UpfPdr *pdr = NULL;

    UTLT_Assert(session, return NULL, "no session");

    IndexAlloc(&upfPdrPool, pdr);
    UTLT_Assert(pdr, return NULL, "PDR context allocation failed");

    pdr->pdrId = 0;
    //pdr->outerHdrRemove = 0;

    pdr->pfcpNode = NULL;
    pdr->far = NULL;
    pdr->qer = NULL;
    pdr->urr = NULL;
    pdr->session = NULL;

    return pdr;
}

Status UpfPdrRemove(UpfPdr *pdr) {
    UTLT_Assert(pdr, return STATUS_ERROR, "PDR error");
    UTLT_Assert(pdr->session, return STATUS_ERROR, "PDR no session");

    if (pdr->pdrId) {
        // Check if UL or  DL
        if (pdr->sourceInterface == PFCP_SRC_INTF_ACCESS) {
            ListRemove(&pdr->session->ulPdrList, pdr);
        } else {
            ListRemove(&pdr->session->dlPdrList, pdr);
        }
    }

    if (pdr->far) {
        UpfFarRemove(pdr->far);
    }

    IndexFree(&upfPdrPool, pdr);

    return STATUS_OK;
}

UpfPdr *UpfPdrFindByPdrId(uint16_t pdrId) {
    int idx;
    for (idx = 0; idx < IndexSize(&upfPdrPool); ++idx) {
        UpfPdr *pdr = IndexFind(&upfPdrPool, idx);
        if (pdr->pdrId == pdrId) {
            return pdr;
        }
    }

    return NULL;
}

UpfPdr *UpfPdrFindByFarId(uint32_t farId) {
    int idx;
    for (idx = 0; idx < IndexSize(&upfPdrPool); ++idx) {
        UpfPdr *pdr = IndexFind(&upfPdrPool, idx);
        if (pdr->far && pdr->far->farId == farId) {
            return pdr;
        }
    }

    return NULL;
}

UpfPdr *UpfPdrFidByUpfGtpUTeid(uint32_t teid) {

    HashIndex *hashIdx = NULL;

    for (hashIdx = UpfSessionFirst(); hashIdx; hashIdx = UpfSessionNext(hashIdx)) {
        UpfSession *session = UpfSessionThis(hashIdx);
        UTLT_Assert(session, return NULL, "session from hash error");

        UpfPdr *pdr = NULL;
        UpfPdr *defaultPdr = NULL;

        /* Save default PDR */
        defaultPdr = ListFirst(&session->ulPdrList);
        UTLT_Assert(defaultPdr, return NULL, "No default PDR");

        /* Find */
        pdr = defaultPdr;
        for (; pdr; pdr = ListNext(pdr)) {
            if (pdr->sourceInterface != PFCP_SRC_INTF_ACCESS) {
                continue;
            }

            if (pdr->upfGtpUTeid == teid) {
                return pdr;
            }
        }
    }

    return NULL;
}

UpfFar *UpfFarAdd() {
    UpfFar *far = NULL;

    IndexAlloc(&upfFarPool, far);
    UTLT_Assert(far, return NULL, "FAR context allocation failed");

    far->farId = far->index;

    far->pfcpNode = NULL;
    far->bar = NULL;
    //far->gtpNode = NULL;

    ListAppend(&self.farList, far);

    return far;
}

Status UpfFarRemove(UpfFar *far) {
    UTLT_Assert(far, return STATUS_ERROR, "far error");

    ListRemove(&Self()->farList, far);

    IndexFree(&upfFarPool, far);

    return STATUS_OK;
}

UpfFar *UpfFarFindByFarId(uint32_t farId) {
    UpfFar *far = NULL;

    far = ListFirst(&Self()->farList);
    while (far) {
        if (far->farId == farId) {
            break;
        }
        far = ListNext(far);
    }

    // if return NULL, no FAR has the farId
    return far;
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

void SessionHashKeygen(uint8_t *out, int *outLen, uint8_t *imsi, int imsiLen, uint8_t *apn) {
    memcpy(out, imsi, imsiLen);
    strncpy((char *)(out + imsiLen), (char*)apn, MAX_APN_LEN + 1);
    *outLen = imsiLen + strlen((char *)(out + imsiLen));

    return;
}

UpfSession *UpfSessionAdd(PfcpUeIpAddr *ueIp, uint8_t *apn, uint8_t pdnType) {
    UpfSession *session = NULL;

    IndexAlloc(&upfSessionPool, session);
    UTLT_Assert(session, return NULL, "session alloc error");

    //session->gtpNode = NULL;

    if (self.pfcpAddr) {
        session->upfSeid = ((uint64_t)self.pfcpAddr->s4.sin_addr.s_addr << 32) | session->index;
    } else if (self.pfcpAddr6) {
        uint32_t *ptr = (uint32_t *)self.pfcpAddr6->s6.sin6_addr.s6_addr;
        session->upfSeid = (((uint64_t)(*ptr)) << 32) | session->index; // TODO: check if correct
    }
    session->upfSeid = htobe64(session->upfSeid);
    //UTLT_Info()
    session->upfSeid = 0; // TODO: check why

    /* IMSI APN Hash */
    /* APN */
    strncpy((char*)session->pdn.apn, (char*)apn, MAX_APN_LEN + 1);

    ListInit(&session->dlPdrList);
    ListInit(&session->ulPdrList);

    session->pdn.paa.pdnType = pdnType;
    if (pdnType == PFCP_PDN_TYPE_IPV4) {
        session->pdn.paa.addr4 = ueIp->addr4;
    } else if (pdnType == PFCP_PDN_TYPE_IPV6) {
        session->pdn.paa.addr6 = ueIp->addr6;
    } else if (pdnType == PFCP_PDN_TYPE_IPV4V6) {
        // TODO
        // session->ueIpv4 = UpfUeIPAlloc(AF_INET, apn);
        // UTLT_Assert(session->ueIpv4, UpfSessionRemove(session); return NULL, "Cannot allocate IPv4");

        // session->ueIpv6 = UpfUeIPAlloc(AF_INET6, apn);
        // UTLT_Assert(session->ueIpv6, UpfSessionRemove(session); return NULL, "Cannot allocate IPv6");

        // session->pdn.paa.dualStack.addr4 = session->ueIpv4->addr4;
        // session->pdn.paa.dualStack.addr6 = session->ueIpv6->addr6;
    } else {
        UTLT_Assert(0, return NULL, "UnSupported PDN Type(%d)", pdnType);
    }

    /* Generate Hash Key: IP + APN */
    if (pdnType == PFCP_PDN_TYPE_IPV4) {
        SessionHashKeygen(session->hashKey, &session->hashKeylen, (uint8_t *)&session->ueIpv4.addr4, 4, apn);
    } else {
        SessionHashKeygen(session->hashKey, &session->hashKeylen, (uint8_t *)&session->ueIpv6.addr6, IPV6_LEN, apn);
    }

    HashSet(self.sessionHash, session->hashKey, session->hashKeylen, session);

    /* initial the session's packIdx to 0 */
    session->pktBufIdx = 0;
    pthread_mutex_init(&session->bufLock, 0);

    return session;
}

Status UpfSessionRemove(UpfSession *session) {
    UpfPdr *pdr;
    UTLT_Assert(self.sessionHash, return STATUS_ERROR, "sessionHash error");
    UTLT_Assert(session, return STATUS_ERROR, "session error");

    pthread_mutex_destroy(&session->bufLock);
    HashSet(self.sessionHash, session->hashKey, session->hashKeylen, NULL);

    // if (session->ueIpv4) {
    //     UpfUeIPFree(session->ueIpv4);
    // }
    // if (session->ueIpv6) {
    //     UpfUeIPFree(session->ueIpv6);
    // }

    pdr = ListFirst(&session->dlPdrList);
    while (pdr) {
        UpfPdrRemove(pdr);
        pdr = ListFirst(&session->dlPdrList);
    }

    pdr = ListFirst(&session->ulPdrList);
    while (pdr) {
        UpfPdrRemove(pdr);
        pdr = ListFirst(&session->ulPdrList);
    }

    IndexFree(&upfSessionPool, session);

    return STATUS_OK;
}

Status UpfSessionRemoveAll() {
    HashIndex *hashIdx = NULL;
    UpfSession *session = NULL;

    for (hashIdx = UpfSessionFirst(); hashIdx; hashIdx = UpfSessionNext(hashIdx)) {
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

    PFCPSessionEstablishmentRequest *request = &message->pFCPSessionEstablishmentRequest;

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

    session = UpfSessionAdd((PfcpUeIpAddr *)&request->createPDR[0].pDI.uEIPAddress.value,
            request->createPDR[0].pDI.networkInstance.value,
            ((int8_t *)request->pDNType.value)[0]);
    UTLT_Assert(session, return NULL, "session add error");

    session->smfSeid = *(uint64_t*)request->cPFSEID.value;
    session->upfSeid = session->index+1;
    UTLT_Trace("UPF Establishment UPF SEID: %lu", session->upfSeid);

    return session;
}

UpfSession *UpfSessionFindByPdrTeid(uint32_t teid) {
    UpfPdr *pdr = NULL;
    for (pdr = ListFirst(&Self()->pdrList); pdr; pdr = ListNext(pdr)) {
        if (pdr->upfGtpUTeid == teid) {
            return pdr->session;
        }
    }
    return NULL;
}
