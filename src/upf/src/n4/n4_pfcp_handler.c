#define TRACE_MODULE _n4_pfcp_handler

#include <endian.h>
#include <arpa/inet.h>

#include "utlt_list.h"
#include "utlt_network.h"

#include "upf_context.h"
#include "pfcp_message.h"
#include "pfcp_xact.h"
#include "pfcp_convert.h"
#include "gtp_path.h"
#include "n4_pfcp_build.h"


void UpfN4HandleCreatePdr(UpfSession *session, CreatePDR *createPdr, UpfPdr **upfPdr) {
    UpfPdr *tmpPdr = NULL;
    uint32_t farId;
    UpfFar *far = NULL;

    UTLT_Assert(createPdr->pDRID.presence, return, "pdr id not presence");
    UTLT_Assert(createPdr->precedence.presence, return, "precedence not presence");
    UTLT_Assert(createPdr->pDI.presence, return, "Pdi not exist");
    UTLT_Assert(createPdr->pDI.sourceInterface.presence, return, "PDI SourceInterface not presence");

    tmpPdr = UpfPdrAdd(session);
    tmpPdr->presence = ntohl(*((uint32_t *)createPdr->precedence.value));
    tmpPdr->sourceInterface = *((uint8_t *)(createPdr->pDI.sourceInterface.value));

    if ((tmpPdr->sourceInterface != PFCP_SRC_INTF_ACCESS) && (tmpPdr->sourceInterface != PFCP_SRC_INTF_SGILAN)) {
        UTLT_Assert(far, UpfPdrRemove(tmpPdr); return, "Source Interface not support");
    }

    // Get From RAN
    if (tmpPdr->sourceInterface == PFCP_SRC_INTF_ACCESS) {
        PfcpFTeid *fTeid;
        UTLT_Assert(createPdr->pDI.localFTEID.presence, UpfPdrRemove(tmpPdr); return, "local F-TEID not found");
        fTeid = (PfcpFTeid *)createPdr->pDI.localFTEID.value;
        tmpPdr->upfGtpUTeid = ntohl(fTeid->teid);
        // TODO: check IP address, fTeid->ipv4, ipv6
    }

    if (createPdr->outerHeaderRemoval.presence) {
        tmpPdr->outerHeaderRemove = *((uint8_t *)(createPdr->outerHeaderRemoval.value));
        if ((tmpPdr->outerHeaderRemove != PFCP_OUTER_HDR_RMV_DESC_GTPU_IP4) &&
            (tmpPdr->outerHeaderRemove != PFCP_OUTER_HDR_RMV_DESC_GTPU_IP6)) {
            UTLT_Assert(0, UpfPdrRemove(tmpPdr); return, "Outer header removal type not support");
        }
    }

    if (createPdr->fARID.presence) {
        farId = ntohl(*((uint32_t *)createPdr->fARID.value));
        far = UpfFarFindByFarId(farId);
        UTLT_Assert(far, UpfPdrRemove(tmpPdr); return, "FAR ID not found");
        tmpPdr->far = far;

        UTLT_Info("[PFCP] - Create %s PDR[0x%04x] TEID[0x%08x] & FAR ID[0x%08x]",
                    tmpPdr->sourceInterface == PFCP_SRC_INTF_ACCESS ? "UL" : "DL",
                    tmpPdr->pdrId, tmpPdr->upfGtpUTeid, farId);
    }

    if (tmpPdr->sourceInterface == PFCP_SRC_INTF_ACCESS) {
        ListAppend(&session->ulPdrList, tmpPdr);
    } else {
        ListAppend(&session->dlPdrList, tmpPdr);
    }

    tmpPdr->pdrId = ntohs(*((uint16_t *)createPdr->pDRID.value));

    *upfPdr = tmpPdr;
}

void UpfN4HandleCreateFar(CreateFAR *createFar, UpfFar **upfFar) {
    UpfFar *tmpFar = NULL;
    UTLT_Assert(createFar->fARID.presence, return, "Far ID not presence");

    // Create FAR
    tmpFar = UpfFarAdd();
    tmpFar->farId = ntohl(*((uint32_t *)createFar->fARID.value));
    tmpFar->applyAction = *((uint8_t *)(createFar->applyAction.value));
    UTLT_Assert(createFar->forwardingParameters.presence, UpfFarRemove(tmpFar); return, "Miss of forwarding parameters");
    tmpFar->destinationInterface = *((uint8_t *)createFar->forwardingParameters.destinationInterface.value); // TODO: fix bug of forwardingParameters type

    if (tmpFar->applyAction == PFCP_FAR_APPLY_ACTION_FORW) {
        // Down Link
        if (tmpFar->destinationInterface == PFCP_FAR_DEST_INTF_ACCESS) {
            /*
            PfcpOuterHdr *outerHdr;
            Ip ip;
            SockNode *gtpNode = NULL;
            Status status;

            UTLT_Assert(createFar->forwardingParameters.outerHeaderCreation.presence, return, "IE messing: outer header creation of create far"); // TODO: fix forwardingParameters IE type
            outerHdr = (PfcpOuterHdr *)(createFar->forwardingParameters.outerHeaderCreation.value);
            if (!(outerHdr->gtpuIpv4 || outerHdr->gtpuIpv6)) {
                UTLT_Assert(0, UpfFarRemove(tmpFar); return, "Forwarding no GTPU");
            }

            tmpFar->upfN3Teid = ntohl(outerHdr->teid);
            status = PfcpOuterHdrToIp(outerHdr, &ip);
            UTLT_Assert(status == STATUS_OK, UpfFarRemove(tmpFar); return, "Outer hdr to IP error");

            //gtpNode = GtpFindNodeByIp(&Self()->ranS1uList, &ip); // TODO: this function
            if (!gtpNode) {
                //gtpNode = GtpAddNodeWithIp(); // TODO: this function
                UTLT_Assert(gtpNode, UpfFarRemove(tmpFar); return, "node create failed");
                Sock *tmpSock = GtpClientCreate(AF_INET, gtpNode->ip, GetPort(&gtpNode->sock->remoteAddr)); // TODO: choose ipv4 or ipv6
                UTLT_Assert(tmpSock, return, "Gtp Client Create error");
            }
            tmpFar->gtpNode = gtpNode;
            */
        }
    }
    UTLT_Info("[PFCP] - Create %s FAR ID[0x%08x]", tmpFar->destinationInterface == PFCP_FAR_DEST_INTF_ACCESS ? "DL" : "UL", tmpFar->farId);
    ListAppend(&Self()->farList, tmpFar);
    *upfFar = tmpFar;
}

void UpfN4HandleSessionEstablishmentRequest(
        UpfSession *session, PfcpXact *pfcpXact, PFCPSessionEstablishmentRequest *request) {
    UpfPdr *pdr0;
    UpfPdr *pdr1;
    UpfFar *far0;
    UpfFar *far1;

    UTLT_Assert(session, return, "Upf Session error");
    UTLT_Assert(pfcpXact, return, "pfcpXact error");
    //UTLT_Assert(pfcpXact->gtpBuf, return, "GTP buffer of pfcpXact error");
    //UTLT_Assert(pfcpXact->gtpXact, return, "GTP Xact of pfcpXact error");

    if (request->createFAR[0].presence) {
        UpfN4HandleCreateFar(&request->createFAR[0], &far0);
    }
    if (request->createFAR[1].presence) {
        UpfN4HandleCreateFar(&request->createFAR[1], &far1);
    }
    if (request->createURR.presence) {
        // TODO
    }
    if (request->createBAR.presence) {
        // TODO
    }
    if (request->createQER.presence) {
        // TODO
    }
    if (request->createPDR[0].presence) {
        UpfN4HandleCreatePdr(session, &request->createPDR[0], &pdr0);
    }
    if (request->createPDR[1].presence) {
        UpfN4HandleCreatePdr(session, &request->createPDR[1], &pdr1);
    }

    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;
    PfcpFSeid *smfFSeid = NULL;

    UTLT_Info("[PFCP] Session Establishment Response");

    if (!request->cPFSEID.presence) {
        UTLT_Error("Session Establishment Response: No CP F-SEID");
        return;
    }

    smfFSeid = request->cPFSEID.value;
    session->smfSeid = be64toh(smfFSeid->seid);

    /* Send Response */
    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_SESSION_ESTABLISHMENT_RESPONSE;
    header.seid = session->smfSeid;

    status = UpfN4BuildSessionEstablishmentResponse(
        &bufBlk, header.type, session, request);
    UTLT_Assert(status == STATUS_OK, return, "N4 build error");

    status = PfcpXactUpdateTx(pfcpXact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return, "pfcpXact update TX error");

    status = PfcpXactCommit(pfcpXact);
    UTLT_Assert(status == STATUS_OK, return, "xact commit error");

    return;
}

void UpfN4HandleSessionModificationRequest(
        UpfSession *session, PfcpXact *xact, PFCPSessionModificationRequest *request) {
    UTLT_Assert(session, return, "Session error");
    UTLT_Assert(xact, return, "xact error");

    Status status;
    PfcpHeader header;
    Bufblk *bufBlk;

    /* Update PDR */
    if (request->updatePDR.presence) {
        UpfPdr *tmpPdr = UpfPdrFindByPdrId(&request->updatePDR.pDRID);
        if (tmpPdr) {
            memcpy(&tmpPdr->presence, (uint32_t *)&request->updatePDR.precedence.value, sizeof(uint32_t)); // TODO: check why set presence to precedence
        } else {
            /* PDR missing in updatePDR */
            UTLT_Assert(0, , "[PFCP] Session Modification Request miss of PDR context");
        }
    }

    /* Update FAR */
    if (request->updateFAR.presence) {
        UpfFar *tmpFar = UpfFarFindByFarId(ntohl(*(uint32_t *)request->updateFAR.fARID.value));

        if (tmpFar) {
            /* Update RAN GTPU-TEID */
            if (request->updateFAR.updateForwardingParameters.outerHeaderCreation.value) { // TODO: fix pyDocxPfcp bug to parse updateForwardingParameters IE bugs
                Ip ip;
                SockNode *gtpNode = NULL;
                PfcpOuterHdr *outerHeader;
                outerHeader = request->updateFAR.updateForwardingParameters.outerHeaderCreation.value;
                tmpFar->upfN3Teid = ntohl(outerHeader->teid);
                status = PfcpOuterHdrToIp(outerHeader, &ip);
                UTLT_Assert(status == STATUS_OK, UpfFarRemove(tmpFar); return, "Outer header IP convert error");

                /* Update GTP Link */
                if (!gtpNode) {
                    // gtpNode = GtpAddNodeWithIp(); // TODO: this function (Self()->ranS1uList, &ip, Self()->gtpv1Port, Self()->parameter.noIpv4, Self()->parameter.noIpv6, Self()->parameter.preferIpv4)
                    UTLT_Assert(gtpNode, UpfFarRemove(tmpFar); return, "GTP node create error");
                    Sock *tmpSock = GtpClientCreate(AF_INET, gtpNode->ip, GetPort(&gtpNode->sock->remoteAddr)); // TODO: choose ipv4 or ipv6
                    UTLT_Assert(tmpSock, return, "GTP client create error");
                }
                tmpFar->gtpNode = gtpNode;
            }

            if (request->updateFAR.applyAction.presence) {
                tmpFar->applyAction = *(uint8_t *)request->updateFAR.applyAction.value;
            }
        } else {
            /* FAR not found */
            UTLT_Assert(0, , "[PFCP] Session Modification: updateFar no FAR");
        }
    }

    if (request->createPDR[0].presence) {
        PfcpFTeid *fTeid;
        PfcpUeIpAddr *ueIp;
        if (request->createPDR[0].pDI.presence) {
            fTeid = (PfcpFTeid *)request->createPDR[0].pDI.localFTEID.value;
            ueIp = (PfcpUeIpAddr*)request->createPDR[0].pDI.uEIPAddress.value;
        } else {
            UTLT_Warning("create PDR PDI missing");
        }
        if (request->createFAR[0].presence) {
            if (request->createFAR[0].forwardingParameters.presence) {
                PfcpOuterHdr *hdr = request->createFAR[0].forwardingParameters.outerHeaderCreation.value;
                Gtpv1TunDevNode *gtpList = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
                char *ifname = gtpList->ifname;
                char ueIpStr[INET_ADDRSTRLEN], ranIpStr[INET_ADDRSTRLEN];
                inet_ntop(AF_INET, &(ueIp->addr4), ueIpStr, INET_ADDRSTRLEN);
                inet_ntop(AF_INET, &(hdr->addr4), ranIpStr, INET_ADDRSTRLEN);
                GtpTunnelAdd(ifname, ntohl(fTeid->teid), ntohl(hdr->teid), ueIpStr, ranIpStr);
            } else {
                UTLT_Info("forwarding parameter missing");
            }
        } else {
            UTLT_Info("create FAR missing");
        }
    }

    /* Send Session Modification Response */
    UTLT_Info("[PFCP] Session Modification Response");

    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_SESSION_MODIFICATION_RESPONSE;
    header.seid = session->smfSeid;

    status = UpfN4BuildSessionModificationResponse(
        &bufBlk, header.type, session, request);
    UTLT_Assert(status == STATUS_OK, return, "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return, "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return, "PFCP Commit error");

    return;
}

void UpfN4HandleSessionDeletionRequest(UpfSession *session, PfcpXact *xact, PFCPSessionDeletionRequest *request) {
    UTLT_Assert(session, return, "session error");
    UTLT_Assert(xact, return, "xact error");

    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;

    /* Send Session Deletion Response */
    UTLT_Info("[PFCP] Session Deletion Response");
    memset(&header, 0, sizeof(PfcpHeader));

    header.type = PFCP_SESSION_DELETION_RESPONSE;
    header.seid = session->smfSeid;

    status = UpfN4BuildSessionDeletionResponse(&bufBlk, header.type, session, request);
    UTLT_Assert(status == STATUS_OK, return, "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return, "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return, "xact commit error");

    /* delete session */
    UpfSessionRemove(session);
    // TODO
    // GtpTunnelDel(const char *ifname, int iteid, int oteid);

    return;
}

void UpfN4HandleSessionReportResponse(
        UpfSession *session, PfcpXact *xact, PFCPSessionReportResponse *response) {
    Status status;

    UTLT_Assert(session, return, "Session error");
    UTLT_Assert(xact, return, "xact error");
    UTLT_Assert(response->cause.presence, return, "SessionReportResponse error: no Cause");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return, "xact commit error");

    return;
}

void UpfN4HandleAssociationSetupRequest(PfcpXact *xact, PFCPAssociationSetupRequest *request) {
    PfcpNodeId *nodeId;

    UTLT_Assert(xact, return, "xact error");
    UTLT_Assert(xact->gnode, return, "gNode of xact error");
    UTLT_Assert(request->nodeID.presence, return, "Request missing nodeId");

    nodeId = (PfcpNodeId *)request->nodeID.value;

    xact->gnode->nodeId.type = nodeId->type;
    switch (nodeId->type) {
        case PFCP_NODE_ID_IPV4:
            xact->gnode->nodeId.addr4 = nodeId->addr4;
            break;
        case PFCP_NODE_ID_IPV6:
            xact->gnode->nodeId.addr6 = nodeId->addr6;
            break;
        default:
            UTLT_Assert(0, return, "Request no node id type");
            break;
    }

    /* Accept */
    xact->gnode->state = PFCP_NODE_ST_ASSOCIATED;

    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;

    /* Send */
    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_ASSOCIATION_SETUP_RESPONSE;
    header.seid = 0;

    status = UpfN4BuildAssociationSetupResponse(&bufBlk, header.type);
    UTLT_Assert(status == STATUS_OK, return, "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return, "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return, "xact commit error");

    return;
}

void UpfN4HandleAssociationUpdateRequest(PfcpXact *xact, PFCPAssociationUpdateRequest *request) {
    // TODO
    UTLT_Info("[PFCP] Association Update Request");
}

void UpfN4HandleAssociationReleaseRequest(PfcpXact *xact, PFCPAssociationReleaseRequest *request) {
    // TODO
    UTLT_Info("[PFCP] Association Release Request");
}

void UpfN4HandleHeartbeatRequest(PfcpXact *xact, HeartbeatRequest *request) {
    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;

    UTLT_Info("[PFCP] Heartbeat Request");

    /* Send */
    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_HEARTBEAT_RESPONSE;
    header.seid = 0;

    status = UpfN4BuildHeartbeatResponse(&bufBlk, header.type);
    UTLT_Assert(status == STATUS_OK, return, "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return, "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return, "xact commit error");

    return;
}

void UpfN4HandleHeartbeatResponse(PfcpXact *xact, HeartbeatResponse *response) {
    // if rsv response, nothing to do, else peer may be not alive
    UTLT_Info("[PFCP] Heartbeat Response");
    return;
}
