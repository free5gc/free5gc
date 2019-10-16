#include "n4_pfcp_path.h"

#include <errno.h>
#include "utlt_event.h"
#include "utlt_buff.h"
#include "utlt_debug.h"
#include "n4_pfcp_handler.h"
#include "upf_context.h"
#include "pfcp_path.h"

static int _pfcpReceiveCB(Sock *sock, void *data) {
    //Event event;
    Status status;
    Bufblk *bufBlk = NULL;
    SockAddr from;
    PfcpNode *upf;
    PfcpHeader *pfcpHeader = NULL;

    UTLT_Assert(sock, return -1, "");

    status = PfcpReceiveFrom(sock, &bufBlk, &from);
    if (status != STATUS_OK) {
        if (errno == EAGAIN) {
            return 0;
        }
        return -1;
    }

    pfcpHeader = (PfcpHeader *)bufBlk->buf;

    if (pfcpHeader->version > PFCP_VERSION) {
        unsigned char vFail[8];
        PfcpHeader *pfcpOut = (PfcpHeader *)vFail;
        
        UTLT_Info("Unsupported PFCP version: %d", pfcpHeader->version);
        pfcpOut->flags = (PFCP_VERSION << 5);
        pfcpOut->type = PFCP_VERSION_NOT_SUPPORTED_RESPONSE;
        pfcpOut->length = htons(4);
        pfcpOut->sqn_only = pfcpHeader->sqn_only;
        SockSendTo(sock, vFail, 8); // TODO: must check localAddress / remoteAddress / fd is correct?
        BufblkFree(bufBlk);
        return 0;
    }

    upf = PfcpFindNodeSockAddr(&Self()->upfN4List, &from);
    if (!upf) {
        PfcpFSeid fSeid;
        memset(&fSeid, 0, sizeof(fSeid));
        // IPv4
        if (sock->remoteAddr._family == AF_INET) {
            fSeid.v4 = 1;
            //fSeid.seid = 0; // TOOD: check SEID value
            fSeid.addr4 = from.s4.sin_addr;
            upf = PfcpAddNodeWithSeid(&Self()->upfN4List, &fSeid,
                    Self()->pfcpPort, 0, 1, 0); // TODO: check noIpv4, noIpv6, preferIpv4, originally from context.no_ipv4
            UTLT_Assert(upf, return STATUS_ERROR, "");

            upf->sock = Self()->pfcpSock;
        }
        if (sock->remoteAddr._family == AF_INET6) {
            fSeid.v6 = 1;
            //fSeid.seid = 0;
            fSeid.addr6 = from.s6.sin6_addr;
            upf = PfcpAddNodeWithSeid(&Self()->upfN4List, &fSeid,
                    Self()->pfcpPort, 1, 0, 0);
            UTLT_Assert(upf, return STATUS_ERROR, "");

            upf->sock = Self()->pfcpSock6;
        }
    }

    UTLT_Assert(upf, return -1, "PFCP node not found");

    //event.type = UPF_EVENT_N4_MESSAGE;
    //event.arg0 = (uintptr_t)bufBlk;
    //event.arg1 = (uintptr_t)upf;
    status = EventSend(Self()->eventQ, UPF_EVENT_N4_MESSAGE, 2, bufBlk, upf);
    if (status != STATUS_OK) {
        UTLT_Error("UPF EventSend error");
        BufblkFree(bufBlk);
        return 0;
    }

    return 0;
}

Status PfcpServerInit() {
    Status status;

    status = PfcpServerList(&Self()->pfcpIPList, _pfcpReceiveCB, Self()->epfd);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "Create PFCP Server for IPv4 error");
    status = PfcpServerList(&Self()->pfcpIPv6List, _pfcpReceiveCB, Self()->epfd);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "Create PFCP Server for IPv6 error");

    Self()->pfcpSock = PfcpLocalSockFirst(&Self()->pfcpIPList);
    Self()->pfcpSock6 = PfcpLocalSockFirst(&Self()->pfcpIPv6List);
    Self()->pfcpAddr = PfcpLocalAddrFirst(&Self()->pfcpIPList);
    Self()->pfcpAddr6 = PfcpLocalAddrFirst(&Self()->pfcpIPv6List);

    UTLT_Assert(Self()->pfcpAddr || Self()->pfcpAddr6, return STATUS_ERROR, "No PFCP Server");

    return STATUS_OK;
}

Status PfcpServerTerminate() {
    SockListFree(&Self()->pfcpIPList);
    SockListFree(&Self()->pfcpIPv6List);

    return STATUS_OK;
}
