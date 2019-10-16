#define TRACE_MODULE _n4_dispatcher

#include "utlt_debug.h"
#include "utlt_event.h"
#include "n4_pfcp_handler.h"
#include "pfcp_xact.h"
#include "pfcp_path.h"

void UpfDispatcher(const Event *event) {
    switch ((UpfEvent)event->type) {
    case UPF_EVENT_N4_MESSAGE:;
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

        if (!pfcpMessage->header.seidP) {
            status = PfcpXactReceive(upf, &pfcpMessage->header, &xact);

            if (status != STATUS_OK) {
                BufblkFree(bufBlk);
                BufblkFree(recvBufBlk);
                break;
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
            default:
                UTLT_Error("No implement pfcp type: %d", pfcpMessage->header.type);
            }
            BufblkFree(bufBlk);
            BufblkFree(recvBufBlk);
            break;
        }

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
        session->pfcpNode = upf;

        status = PfcpXactReceive(session->pfcpNode, &pfcpMessage->header, &xact);
        if (status != STATUS_OK) {
            BufblkFree(recvBufBlk);
            BufblkFree(bufBlk);
            return;
        }

        switch (pfcpMessage->header.type) {
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
    default:
        UTLT_Error("No handler for event type: %d", event->type);
        break;
    }
}
