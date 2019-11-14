#define TRACE_MODULE _n4_dispatcher

#include "utlt_debug.h"
#include "utlt_event.h"
#include "n4_pfcp_handler.h"
#include "pfcp_xact.h"
#include "pfcp_path.h"
#include "n4_pfcp_build.h"

void UpfDispatcher(const Event *event) {
    switch ((UpfEvent)event->type) {
    case UPF_EVENT_SESSION_REPORT: {
        Status status;
        PfcpHeader header;
        Bufblk *bufBlk = NULL;
        PfcpXact *xact = NULL;

        uint64_t *seid = (uint64_t*)event->arg0;
        uint16_t *pdrId = (uint16_t*)event->arg1;

        UpfSession *session = UpfSessionFindBySeid(*seid);
        UTLT_Assert(session != NULL, return, "Session not find by seid: %d", *seid);

        memset(&header, 0, sizeof(PfcpHeader));
        header.type = PFCP_SESSION_REPORT_REQUEST;
        header.seid = *seid;

        status = UpfN4BuildSessionReportRequestDownlinkDataReport(&bufBlk, header.type, session, *pdrId);
        UTLT_Assert(status == STATUS_OK, return, "Build Session Report Request error");

        xact = PfcpXactLocalCreate(session->pfcpNode, &header, bufBlk);
        UTLT_Assert(xact, return, "pfcpXactLocalCreate error");

        status = PfcpXactCommit(xact);
        UTLT_Assert(status == STATUS_OK, return, "xact commit error");

	BufblkFree(bufBlk);

        break;
    }
    case UPF_EVENT_N4_MESSAGE: {
        Status status;
        Bufblk *bufBlk = NULL;
        Bufblk *recvBufBlk = (Bufblk *)event->arg0;
        PfcpNode *upf = (PfcpNode *)event->arg1;
        PfcpMessage *pfcpMessage = NULL;
        PfcpXact *xact = NULL;
        UpfSession *session = NULL;

        UTLT_Assert(recvBufBlk, return, "recv buffer no data");
        bufBlk = BufblkAlloc(1, sizeof(PfcpMessage));
        UTLT_Assert(bufBlk, return, "create buffer error");
        pfcpMessage = bufBlk->buf;
        UTLT_Assert(pfcpMessage, return, "pfcpMessage assigned error");

        status = PfcpParseMessage(pfcpMessage, recvBufBlk);
        UTLT_Assert(status == STATUS_OK, return, "PfcpParseMessage error");

        if (pfcpMessage->header.seidP) {

            // if SEID presence
            if (!pfcpMessage->header.seid) {
                // without SEID
                if (pfcpMessage->header.type == PFCP_SESSION_ESTABLISHMENT_REQUEST) {
                    session = UpfSessionAddByMessage(pfcpMessage);
                } else {
                    UTLT_Assert(0, BufblkFree(recvBufBlk); BufblkFree(bufBlk); return,
                            "no SEID but not SESSION ESTABLISHMENT");
                }
            } else {
                // with SEID
                session = UpfSessionFindBySeid(pfcpMessage->header.seid);
            }

            UTLT_Assert(session, BufblkFree(recvBufBlk); BufblkFree(bufBlk); return,
                    "do not find / establish session");

            if (pfcpMessage->header.type != PFCP_SESSION_REPORT_RESPONSE) {
                session->pfcpNode = upf;
            }

            status = PfcpXactReceive(session->pfcpNode, &pfcpMessage->header, &xact);

            if (status != STATUS_OK) {
                BufblkFree(bufBlk);
                BufblkFree(recvBufBlk);
                break;
            }
        } else {
            status = PfcpXactReceive(upf, &pfcpMessage->header, &xact);

            if (status != STATUS_OK) {
                BufblkFree(bufBlk);
                BufblkFree(recvBufBlk);
                break;
            }
        }

        switch (pfcpMessage->header.type) {
        case PFCP_HEARTBEAT_REQUEST:
            UpfN4HandleHeartbeatRequest(
                xact, &pfcpMessage->heartbeatRequest);
            break;
        case PFCP_HEARTBEAT_RESPONSE:
            UpfN4HandleHeartbeatResponse(
                xact, &pfcpMessage->heartbeatResponse);
            break;
        case PFCP_ASSOCIATION_SETUP_REQUEST:
            UpfN4HandleAssociationSetupRequest(
                xact, &pfcpMessage->pFCPAssociationSetupRequest);
            break;
        case PFCP_ASSOCIATION_UPDATE_REQUEST:
            UpfN4HandleAssociationUpdateRequest(
                xact, &pfcpMessage->pFCPAssociationUpdateRequest);
            break;
        case PFCP_ASSOCIATION_RELEASE_RESPONSE:
            UpfN4HandleAssociationReleaseRequest(
                xact, &pfcpMessage->pFCPAssociationReleaseRequest);
            break;
        case PFCP_SESSION_ESTABLISHMENT_REQUEST:
            UpfN4HandleSessionEstablishmentRequest(session,
                xact, &pfcpMessage->pFCPSessionEstablishmentRequest);
            break;
        case PFCP_SESSION_MODIFICATION_REQUEST:
            UpfN4HandleSessionModificationRequest(session,
                xact, &pfcpMessage->pFCPSessionModificationRequest);
            break;
        case PFCP_SESSION_DELETION_REQUEST:
            UpfN4HandleSessionDeletionRequest(session,
                xact, &pfcpMessage->pFCPSessionDeletionRequest);
            break;
        case PFCP_SESSION_REPORT_RESPONSE:
            UpfN4HandleSessionReportResponse(session,
                xact, &pfcpMessage->pFCPSessionReportResponse);
            break;
        default:
            UTLT_Error("No implement pfcp type: %d", pfcpMessage->header.type);
        }
        BufblkFree(bufBlk);
        BufblkFree(recvBufBlk);

        break;
    }
    case UPF_EVENT_N4_T3_RESPONSE:
    case UPF_EVENT_N4_T3_HOLDING: {
        uint8_t type;
        PfcpXactTimeout(*((uint32_t*)event->arg0), (UpfEvent)event->type, &type);
        break;
    }
    default: {
        UTLT_Error("No handler for event type: %d", event->type);
        break;
    }
    }
}
