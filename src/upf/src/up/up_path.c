/*
  The file is for handling UP issue.
  Including Echo Request, Echo response, Error Indication, End marker, buffer
  This will write a little route rule for UE route
  This will create udp/ipv4 socket to handle first four.
  This will create name pipe for kernel sending unmatch packet up to user space.
 */

#define TRACE_MODULE _up_path

#include "up_path.h"

#include "utlt_debug.h"
#include "utlt_network.h"
#include "utlt_buff.h"
#include "pfcp_types.h"
#include "upf_context.h"
#include "knet_route.h"
#include "gtp_path.h"
#include "gtp_header.h"
#include "gtp_buffer.h"
#include "gtp_tunnel.h"
#include "libgtp5gnl/gtp5g.h"

Status UpRouteInit() {
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

    UTLT_Free(routeEntries);

    return STATUS_OK;
}

Status UpRouteTerminate() {
    Status status;

    UTLT_Info("Removing APN routes");
    for (ApnNode *it = ListFirst(&Self()->apnList); it != NULL;
         it = ListNext(it)) {
        // TODO: only get the first dev here, there will be only one gtp tun
        // dev in the future
        Gtpv1TunDevNode *gtpTunDev = ListFirst(&Self()->gtpv1DevList);

        status = KnetDelRoute(gtpTunDev->ifname, it->subnetIP,
                              it->subnetPrefix, NULL, 0);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "");
    }

    return STATUS_OK;
}

Status GTPv1ServerInit() {
    Status status;

    // TODO : One IP mapping to one dev, need to discuss
    status = GtpDevListCreate(Self()->epfd, AF_INET, &Self()->gtpv1DevList,
                              GtpHandler, NULL);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "GtpLinkListCreate Fail");

    return STATUS_OK;
}

Status GTPv1ServerTerminate() {
    Status status = STATUS_OK;

    status = GtpDevListFree(Self()->epfd, &Self()->gtpv1DevList);
    UTLT_Assert(status == STATUS_OK, status |= STATUS_ERROR,
                "GTPv1 tunnel list free fail");

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
    UTLT_Assert((gtpHdr->flags >> 5) == 1, goto FREEBUFBLK,
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
        .flags = 0x30 + (gtpHdr->flags & 0x03),
        .type = GTPV1_ECHO_RESPONSE,
    };

    Bufblk *optPkt = BufblkAlloc(1, 0x40);
    if (gtpRespHrd.flags & 0x03) {
        Gtpv1OptHeader *opthrd = (void *)((uint8_t *) data + GTPV1_HEADER_LEN);
        Gtpv1OptHeader gtpOptHrd = {
            ._seqNum = (gtpRespHrd.flags & 0x02) ? htons(ntohs(opthrd->_seqNum) + 1) : 0,
            .nPdnNum = (gtpRespHrd.flags & 0x01) ? opthrd->nPdnNum : 0,
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
    int gtpPayloadLen = GTPV1_HEADER_LEN + ((gtpHdr->flags & 0x02) ? GTPV1_OPT_HEADER_LEN : 0);

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
// Create name pipe
Status BufferServerInit() {
    Status status = STATUS_OK;

    // HACK: read path from config
    Self()->buffSock = BufferServerCreate(SOCK_DGRAM, Self()->buffSockPath,
                                          BufferHandler, NULL);
    if (&Self()->buffSock == NULL) {
        UTLT_Error("Buffering PIPE cannot create");
        return STATUS_ERROR;
    }

    status = BufferEpollRegister(Self()->epfd, Self()->buffSock);
    UTLT_Assert(status == STATUS_OK, return status, "epoll register error");

    return status;
}

// Clear name pipe
Status BufferServerTerminate() {
    Status status = STATUS_OK;

    status = BufferEpollDeregister(Self()->epfd, Self()->buffSock);
    UTLT_Assert(status == STATUS_OK, return status, "epoll deregister error");

    status = BufferServerFree(Self()->buffSock);
    UTLT_Assert(status == STATUS_OK, return status, "PIPE cannot be free");

    return status;
}

// Handle when recv a packet from GTP Dev but not matching any rules
Status BufferHandler(Sock *sock, void *data) {
    UTLT_Assert(sock, return STATUS_ERROR, "Unix socket not found");

    UTLT_Debug("BufferHandler get event");

    uint8_t farAction;
    uint16_t pdrId;
    Bufblk *pktbuf = BufblkAlloc(1, MAX_OF_BUFFER_PACKET_SIZE);

    // BufferRecv return -1 if error
    int readNum = BufferRecv(sock, pktbuf, &pdrId, &farAction);
    UTLT_Assert(readNum >= 0, goto ERROR_AND_FREE, "Buffer receive fail");

    if (farAction & PFCP_FAR_APPLY_ACTION_BUFF) {
        // Store packet in context &Self->bufPacketHash
        UpfBufPacket *packetStorage = UpfBufPacketFindByPdrId(pdrId);
        UTLT_Assert(packetStorage, goto ERROR_AND_FREE,
                    "Cannot find matching PDR ID buffer slot");

        Status status;
        // protect data write with spinlock
        // instead of protect code block with mutex
        UTLT_Assert(!pthread_spin_lock(&Self()->buffLock), goto ERROR_AND_FREE,
                    "spin lock buffLock error");
        if (packetStorage->packetBuffer) {
            // if packetBuffer not null, just add packet followed
            status = BufblkBuf(packetStorage->packetBuffer, pktbuf);
            UTLT_Assert(status == STATUS_OK, goto ERROR_AND_FREE,
                        "block add behand old buffer error");
            // free the pktbuf
            UTLT_Assert(BufblkFree(pktbuf) == STATUS_OK, return STATUS_ERROR,
                        "Bufblk free fail");
        } else {
            // if packetBuffer null, allocate space
            // reuse the pktbuf, so don't free it
            packetStorage->packetBuffer = pktbuf;
        }
        while (pthread_spin_unlock(&Self()->buffLock)) {
            // if unlock failed, keep trying
            UTLT_Error("spin unlock error");
        }

        // If NOCP, Send event to notify SMF
        uint64_t seid = ((UpfSession*)packetStorage->sessionPtr)->upfSeid;
        UTLT_Debug("buffer NOCP to SMF: SEID: %u, PDRID: %u", seid, pdrId);
        status = EventSend(Self()->eventQ, UPF_EVENT_SESSION_REPORT, 2,
                           seid, pdrId);
        UTLT_Assert(status == STATUS_OK, ,
                    "DL data message event send to N4 failed");
    } else {
        UTLT_Warning("apply action no BUF tag, drop it");
        UTLT_Assert(BufblkFree(pktbuf) == STATUS_OK, return STATUS_ERROR,
                    "Bufblk free fail");
    }

    // When no jump to ERROR_AND_FREE, it means ok
    return STATUS_OK;

ERROR_AND_FREE:
    UTLT_Assert(BufblkFree(pktbuf) == STATUS_OK, , "Bufblk free fail");
    return STATUS_ERROR;
}

Status UpSendPacketByPdrFar(UpfPdr *pdr, UpfFar *far, Sock *sock) {
    UTLT_Assert(pdr, return STATUS_ERROR, "PDR error");
    UTLT_Assert(far, return STATUS_ERROR, "FAR error");
    UTLT_Assert(sock, return STATUS_ERROR, "Send packet sock error");
    Status status = STATUS_OK;

    // if GTP IPV4
    uint16_t *desPtr = gtp5g_far_get_outer_header_creation_description(far);
    if (desPtr && (*desPtr & 1)) {
        uint32_t *teidPtr = gtp5g_far_get_outer_header_creation_teid(far);
        UTLT_Assert(teidPtr, return STATUS_ERROR, "TEID not found");
        Gtpv1Header gtpHdr = {
            .flags = 0x30,
            .type = GTPV1_T_PDU,
            ._teid = *teidPtr,
        };

        uint16_t pdrId = *gtp5g_pdr_get_id(pdr);
        UpfBufPacket *bufStorage = UpfBufPacketFindByPdrId(pdrId);
        if (bufStorage->packetBuffer) {
            UTLT_Assert(!pthread_spin_lock(&Self()->buffLock),
                        return STATUS_ERROR, "spin lock buffLock error");

            Bufblk *sendBuf = BufblkAlloc(1, 0x40);
            gtpHdr._length = htons(bufStorage->packetBuffer->len);
            BufblkBytes(sendBuf, (void*)&gtpHdr, GTPV1_HEADER_LEN);
            BufblkBuf(sendBuf, bufStorage->packetBuffer);

            status = GtpSend(sock, sendBuf);
            UTLT_Assert(status == STATUS_OK, return status, "GtpSend failed");
            BufblkClear(sendBuf);

            while (pthread_spin_unlock(&Self()->buffLock)) {
                // if unlock failed, keep trying
                UTLT_Error("spin unlock error");
            }
            status = BufblkFree(bufStorage->packetBuffer);
        } else {
            UTLT_Debug("bufStorage is NULL");
        }
        UTLT_Assert(status == STATUS_OK, return status,
                    "Free packet buffer failed");
        bufStorage->packetBuffer = NULL;
    } else {
        UTLT_Warning("outer header creatation not implement: "
                     "GTP-IPV6, IPV4, IPV6");
    }

    return status;
}
