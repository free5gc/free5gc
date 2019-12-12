#include "up_gtp_path.h"

#include "utlt_debug.h"
#include "utlt_buff.h"
#include "utlt_network.h"
#include "utlt_list.h"
#include "knet_route.h"
#include "gtp_header.h"
#include "gtp_path.h"
#include "upf_context.h"

Status GtpRouteInit() {
    Status status;
    
    for (ApnNode *it = ListFirst(&Self()->apnList); it != NULL; it = ListNext(it)) {
        // TODO: only get the first dev here, there will be only one gtp tun dev in the future
        Gtpv1TunDevNode *gtpTunDev = ListFirst(&Self()->gtpv1DevList);

        status = KnetAddRoute(gtpTunDev->ifname, it->subnetIP, it->subnetPrefix, NULL, 0);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "");
    }
    
    // Get routes from main IPv4 routing table and print
    ListNode *routeEntries = KnetGetRoutes(AF_INET, RT_TABLE_MAIN);
    UTLT_Assert(routeEntries, return STATUS_ERROR, "");

    UTLT_Info("APN routes added, main routing table:");
    KnetPrintRoutes(routeEntries);
    KnetRtListFree(routeEntries);

    return STATUS_OK;
}

Status GtpRouteTerminate() {
    Status status;

    UTLT_Info("Removing APN routes");
    for (ApnNode *it = ListFirst(&Self()->apnList); it != NULL; it = ListNext(it)) {
        // TODO: only get the first dev here, there will be only one gtp tun dev in the future
        Gtpv1TunDevNode *gtpTunDev = ListFirst(&Self()->gtpv1DevList);

        status = KnetDelRoute(gtpTunDev->ifname, it->subnetIP, it->subnetPrefix, NULL, 0);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "");
    }

    return STATUS_OK;
}

Status GTPv1ServerInit() {
    Status status;

    // TODO : One IP mapping to one dev, need to discuss
    status = GtpDevListCreate(Self()->epfd, AF_INET, &Self()->gtpv1DevList, GtpHandler, NULL);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "GtpLinkListCreate fail");

    return STATUS_OK;
}

// TODO : GTPv1 Server Terminate
Status GTPv1ServerTerminate() {
    Status status = STATUS_OK;

    UTLT_Assert(GtpTunListFree(Self()->epfd, &Self()->gtpv1DevList) == STATUS_OK,
                status |= STATUS_ERROR, "GTPv1 tunnel list free fail");

    return status;
}

// TODO : Need to handle buffer and drop (Rule will be set at far)
Status GtpHandler(Sock *sock, void *data) {
    UTLT_Assert(sock, return STATUS_ERROR, "GTP socket not found");
    Status status = STATUS_ERROR;

    Bufblk *pktbuf = BufblkAlloc(1, MAX_OF_GTPV1_PACKET_SIZE);
    int readNum = GtpRecv(sock, pktbuf);
    UTLT_Assert(readNum >= 0, goto FREEBUFBLK, "GTP receive fail");

    // TODO : Need to handle buffering and reject, including GTP and general packet
    // Not only GTP packet
    Gtpv1Header *gtpHdr = pktbuf->buf;
    UTLT_Assert(gtpHdr->version == 1, goto FREEBUFBLK,
                "Only handle the GTP version 1 in user plane");

    switch (gtpHdr->type) {
        case GTPV1_ECHO_REQUEST :
            status = GtpHandleEchoRequest(sock, gtpHdr);
            break;
        case GTPV1_ECHO_RESPONSE :
            status = GtpHandleEchoResponse(gtpHdr);
            break;
        case GTPV1_ERROR_INDICATION :

            break;
        case GTPV1_END_MARK :
            // TODO : Need to deal with the UE packet that does not have tunnel yet
            status = GtpHandleEndMark(sock, gtpHdr);
            break;
        case GTPV1_T_PDU :
            // TODO : Check PDR, FAR to forward packet, or maybe do paging and buffer UE packet
            status = GtpHandleTPDU(sock, pktbuf);
            break;
        default : 
            UTLT_Warning("This type[%d] of GTPv1 header does not implement yet", gtpHdr->type);
    }

FREEBUFBLK:
    UTLT_Assert(BufblkFree(pktbuf) == STATUS_OK, , "Bufblk free fail");

    return status;
}

Status GtpHandleEchoRequest(Sock *sock, void *data) {
    UTLT_Assert(data, return STATUS_ERROR, "GTP data is NULL");

    Gtpv1Header *gtpHdr = data;
    UTLT_Assert(gtpHdr->type == GTPV1_ECHO_REQUEST, return STATUS_ERROR, 
                "The type of GTP data is not 'Echo Request'");

    Status status = STATUS_OK;

    // Build the Echo Response packet
    Gtpv1Header gtpRespHrd = {
        .version = 1,
        .PT = 1,
        .type = GTPV1_ECHO_RESPONSE,
        .seqFlag = gtpHdr->seqFlag,
        .pn = gtpHdr->pn,
    };

    Bufblk *optPkt = BufblkAlloc(1, 0x40);
    if (gtpRespHrd.seqFlag | gtpRespHrd.pn) {
        Gtp1OptHeader *opthrd = (void *)((uint8_t *) data + GTPV1_HEADER_LEN);
        Gtp1OptHeader gtpOptHrd = {
            ._seqNum = gtpRespHrd.seqFlag ? htons(ntohs(opthrd->_seqNum) + 1) : 0,
            .nPdnNum = gtpRespHrd.pn ? opthrd->nPdnNum : 0,
        };
        BufblkBytes(optPkt, (void *) &gtpOptHrd, sizeof(gtpOptHrd));
    }

    /* Recover IE */
    uint8_t recoverType = 14, recoverCnt = 0;
    BufblkBytes(optPkt, (void *) &recoverType, 1);
    BufblkBytes(optPkt, (void *) &recoverCnt, 1);

    gtpRespHrd._length = htons(optPkt->len);

    Bufblk *pkt = BufblkAlloc(1, 0x40);
    BufblkBytes(pkt, (void *) &gtpRespHrd, GTPV1_HEADER_LEN);
    BufblkBuf(pkt, optPkt);

    BufblkFree(optPkt);

    UTLT_Assert(GtpSend(sock, pkt) == STATUS_OK, status = STATUS_ERROR,
                "GTP Send fail");

    BufblkFree(pkt);

    return status;
}

Status GtpHandleEchoResponse(void *data) {
    UTLT_Assert(data, return STATUS_ERROR, "GTP data is NULL");

    Gtpv1Header *gtpHdr = data;
    UTLT_Assert(gtpHdr->type == GTPV1_ECHO_RESPONSE, return STATUS_ERROR, 
                "The type of GTP data is not 'Echo Response'");

    // TODO : Check the peer device exists, and ....
    // 29.281 says the restart conter shall be ignore by the receiver

    return STATUS_OK;
}

// TODO : Need to check fepc code
Status GtpHandleEndMark(Sock *sock, void *data) {
    UTLT_Assert(sock && data, return STATUS_ERROR, "GTP data is NULL");
    Status status = STATUS_ERROR;
    /*
    Gtpv1Header *gtpHdr = data;
    int teid = ntohl(gtpHdr->_teid);
    int gtpPayloadLen = GTPV1_HEADER_LEN + (gtpHdr->seqFlag ? GTPV1_OPT_HEADER_LEN : 0);
    
    */

    // TODO : Check PDR, FAR to forward packet, or maybe do paging and buffer UE packet
    /*
    UTLT_Assert(GtpSend(sock, pktbuf) == STATUS_OK, goto FREEBUFBLK, "GTP Send fail");

    status = STATUS_OK;

FREEBUFBLK:
    UTLT_Assert(BufblkFree(pktbuf) == STATUS_OK, status = STATUS_ERROR, "Bufblk free fail");
*/
    status = STATUS_OK;
    return status;
}

Status GtpHandleTPDU(Sock *sock, Bufblk *data) {
    UTLT_Assert(sock && data, return STATUS_ERROR, "GTP data is NULL");
    Status status = STATUS_ERROR;

    Gtpv1Header *gtpHdr = data->buf;
    int teid = ntohl(gtpHdr->_teid);

    UpfSession *sess = UpfSessionFindByPdrTeid(teid);
    if (sess == NULL) {
        UTLT_Warning("The TEID[0x%x] does not match any session", teid);
        return STATUS_OK;
    }
    
    // Store the packet and trigger Data Notification
    UTLT_Assert(UpfSessionPacketRecv(sess, data) == STATUS_OK, return STATUS_ERROR,
                "UPF Store non-tunnel UE packet fail : TEID[0x%x]", teid);

    status = STATUS_OK;
    return status;
}

Status UpfSessionPacketSend(UpfSession *session, Sock *sock) {
    UTLT_Assert(session, return STATUS_ERROR, "UPF session is NULL");
    UTLT_Assert(sock, return STATUS_ERROR, "Socket is NULL");

    UpfPdr *pdr = ListFirst(&session->dlPdrList);
    UTLT_Assert(pdr, return STATUS_ERROR, "There is no PDR in this session");

    uint32_t teid = pdr->upfGtpUTeid;
    Gtpv1Header gtpv1Hdr = {
        .version = 1,
        .type = GTPV1_T_PDU,
        ._teid = htons(teid),
    };

    pthread_mutex_lock(&session->bufLock);
    // Set RAN IP and Port
    sock->remoteAddr._family = sock->localAddr._family;
    sock->remoteAddr._port = sock->localAddr._port;    // Default : 2152
    if (sock->localAddr._family == AF_INET)
        sock->remoteAddr.s4.sin_addr = pdr->far->ranIp.addr4;
    else
        sock->remoteAddr.s6.sin6_addr = pdr->far->ranIp.addr6;

    Bufblk *sendBuf = BufblkAlloc(1, 0x40);
    for (int i = 0; i < session->pktBufIdx; i++) {
        gtpv1Hdr._length = htons(session->packetBuffer[i]->len);
        BufblkBytes(sendBuf, (void *) &gtpv1Hdr, GTPV1_HEADER_LEN);
        BufblkBuf(sendBuf, session->packetBuffer[i]);

        Status status = GtpSend(sock, sendBuf);
        UTLT_Assert(status == STATUS_OK, , "GTP Send Error");
        BufblkClear(sendBuf);
        BufblkFree(session->packetBuffer[i]);
    }
    BufblkFree(sendBuf);
    pthread_mutex_unlock(&session->bufLock);

    return STATUS_OK;
}

Status UpfSessionPacketRecv(UpfSession *session, Bufblk *pktBuf) {
    Status status = STATUS_OK;

    UTLT_Assert(session, return STATUS_ERROR, "UPF session is NULL");
    UTLT_Assert(pktBuf, return STATUS_ERROR, "Packet buffer is NULL");

    UpfPdr *pdr = ListFirst(&session->dlPdrList);
    // TODO : Find rule for UE packet
    /*
    for (; pdr; pdr = ListNext(pdr)) {

    }
    */
    UTLT_Assert(pdr, return STATUS_ERROR, "There is no PDR in this session");

    pthread_mutex_lock(&session->bufLock);
    if (session->pktBufIdx <= 0 && (pdr->far->applyAction & PFCP_FAR_APPLY_ACTION_NOCP)) {
        // Trigger DL data notification
        status = EventSend(Self()->eventQ, UPF_EVENT_SESSION_REPORT, 2, &session->upfSeid, &pdr->pdrId);
        UTLT_Assert(status == STATUS_OK, , "DL data message event send to N4 fail");
    }

    uint32_t pktBufIdx = (status != STATUS_OK ? session->pktBufIdx * -1 : session->pktBufIdx);
    uint32_t teid = pdr->upfGtpUTeid;
    UTLT_Assert(pktBufIdx < MAX_NUM_OF_PACKET_BUFFER_SIZE, return STATUS_ERROR,
                "The buffer in this session is full : DL TEID[0x%x]", teid);

    Bufblk *recvBuf = BufblkAlloc(1, 0x40);
    BufblkBuf(recvBuf, pktBuf);
    session->packetBuffer[pktBufIdx] = recvBuf;

    session->pktBufIdx = (status != STATUS_OK ? session->pktBufIdx - 1 : session->pktBufIdx + 1);
    pthread_mutex_unlock(&session->bufLock);

    return STATUS_OK;
}

Status UpfSessionPacketClear(UpfSession *session) {
    UTLT_Assert(session, return STATUS_ERROR, "UPF session is NULL");

    pthread_mutex_lock(&session->bufLock);
    uint32_t pktBufIdx = (session->pktBufIdx < 0 ? session->pktBufIdx * -1 : session->pktBufIdx);
    for (int i = 0; i < pktBufIdx; i++) {
        BufblkFree(session->packetBuffer[i]);
    }
    pthread_mutex_unlock(&session->bufLock);

    return STATUS_OK;
}

