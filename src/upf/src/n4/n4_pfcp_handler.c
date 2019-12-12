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
#include "up/up_gtp_path.h"


Status _AddGtpByPdrFar(UpfPdr *pdr, UpfFar *far) {
    UTLT_Assert(pdr, return STATUS_ERROR, "pdr error");
    UTLT_Assert(far, return STATUS_ERROR, "far error");

    Gtpv1TunDevNode *gtpv1Dev4 = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    char *ifname = gtpv1Dev4->ifname;
    char ueIpStr[INET_ADDRSTRLEN], ranIpStr[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &(pdr->ueIp.addr4), ueIpStr, INET_ADDRSTRLEN);
    inet_ntop(AF_INET, &(far->ranIp.addr4), ranIpStr, INET_ADDRSTRLEN);
    UTLT_Debug("GTP Add: interface: %s, sTEID(UPF): %d, oTEID(NR): %d, UE IP: %s, RAN IP: %s",
        ifname, pdr->upfGtpUTeid, far->upfN3Teid, ueIpStr, ranIpStr);
    return GtpTunnelAdd(ifname, pdr->upfGtpUTeid,
                        far->upfN3Teid, ueIpStr, ranIpStr);
}

Status _AddGtpByPdr(UpfPdr *pdr) {
    return _AddGtpByPdrFar(pdr, pdr->far);
}

Status _RemoveGtpByPdr(UpfPdr *pdr) {
    Status status;
    UTLT_Assert(pdr, return STATUS_ERROR, "pdr error");

    Gtpv1TunDevNode *gtpv1Dev4 = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    char *ifname = gtpv1Dev4->ifname;

    status = GtpTunnelDel(ifname, pdr->upfGtpUTeid);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "GTP Tunnel Remove UP error");

    return STATUS_OK;
}

void UpfN4HandleCreatePdr(UpfSession *session, CreatePDR *createPdr, UpfPdr **upfPdr) {
    UpfPdr *tmpPdr = NULL;
    uint32_t farId;
    UpfFar *far = NULL;

    UTLT_Assert(createPdr->pDRID.presence, return, "pdr id not presence");
    UTLT_Assert(createPdr->precedence.presence, return, "precedence not presence");
    UTLT_Assert(createPdr->pDI.presence, return, "Pdi not exist");
    UTLT_Assert(createPdr->pDI.sourceInterface.presence, return, "PDI SourceInterface not presence");

    tmpPdr = UpfPdrAdd(session);
    tmpPdr->precedence = ntohl(*((uint32_t *)createPdr->precedence.value));
    tmpPdr->sourceInterface = *((uint8_t *)(createPdr->pDI.sourceInterface.value));
    tmpPdr->pdrId = ntohs(*((uint16_t *)createPdr->pDRID.value));

    if ((tmpPdr->sourceInterface != PFCP_SRC_INTF_ACCESS) && (tmpPdr->sourceInterface != PFCP_SRC_INTF_SGILAN)) {
        UTLT_Assert(0, UpfPdrRemove(tmpPdr); return, "Source Interface not support");
    }

    // Get From RAN (dl)
    if (tmpPdr->sourceInterface == PFCP_SRC_INTF_SGILAN) {
        tmpPdr->ulDl = 1;
        PfcpFTeid *fTeid;
        UTLT_Assert(createPdr->pDI.localFTEID.presence, UpfPdrRemove(tmpPdr); return, "local F-TEID not found");
        fTeid = (PfcpFTeid *)createPdr->pDI.localFTEID.value;
        tmpPdr->upfGtpUTeid = ntohl(fTeid->teid);
        if (fTeid->v4 && fTeid->v6) {
            memcpy(&tmpPdr->dualStack, &fTeid->dualStack,
                sizeof(struct in_addr) + sizeof(struct in6_addr));
            // TODO: session ueIp dualStack
        } else if (fTeid->v4) {
            tmpPdr->addr4 = fTeid->addr4;
        } else if (fTeid->v6) {
            tmpPdr->addr6 = fTeid->addr6;
        }

        // UE IP
        if (createPdr->pDI.uEIPAddress.presence) {
            PfcpUeIpAddr *ueIp = (PfcpUeIpAddr*)createPdr->pDI.uEIPAddress.value;
            memcpy(&(tmpPdr->ueIp.addr4), &ueIp->addr4, sizeof(uint32_t));
        }
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

        UTLT_Info("[PFCP] - Create %s PDR[0x%04x] TEID[0x%08x] & FAR ID[0x%04x]",
                    tmpPdr->sourceInterface == PFCP_SRC_INTF_ACCESS ? "UL" : "DL",
                    tmpPdr->pdrId, tmpPdr->upfGtpUTeid, farId);

        // Add GTP Tunnel
        if (tmpPdr->sourceInterface == PFCP_SRC_INTF_SGILAN) {
            Status status = _AddGtpByPdrFar(tmpPdr, far);
            UTLT_Assert(status == STATUS_OK, , "Gtp Add failed");
        }
    }

    if (tmpPdr->sourceInterface == PFCP_SRC_INTF_ACCESS) {
        ListAppend(&session->ulPdrList, tmpPdr);
    } else {
        ListAppend(&session->dlPdrList, tmpPdr);
    }

    ListAppend(&Self()->pdrList, tmpPdr);
    tmpPdr->session = session;
    *upfPdr = tmpPdr;
}

void UpfN4HandleCreateFar(CreateFAR *createFar, UpfFar **upfFar) {
    UpfFar *tmpFar = NULL;
    UTLT_Assert(createFar->fARID.presence, return, "Far ID not presence");

    // Create FAR
    tmpFar = UpfFarAdd();
    tmpFar->farId = ntohl(*((uint32_t *)createFar->fARID.value));
    tmpFar->applyAction = *((uint8_t *)(createFar->applyAction.value));
    UTLT_Assert(createFar->forwardingParameters.presence, UpfFarRemove(tmpFar);
		return, "Miss of forwarding parameters");
    tmpFar->destinationInterface = *((uint8_t *)createFar->forwardingParameters.destinationInterface.value);

    if (tmpFar->applyAction & PFCP_FAR_APPLY_ACTION_FORW) {
        if (createFar->forwardingParameters.destinationInterface.presence) {
            // Down Link
            if (tmpFar->destinationInterface == PFCP_FAR_DEST_INTF_ACCESS) {
                PfcpOuterHdr *outerHdr;
                Ip ip;
                Status status;

                UTLT_Assert(createFar->forwardingParameters.outerHeaderCreation.presence,
			    return, "IE messing: outer header creation of create far");
                outerHdr = (PfcpOuterHdr *)(createFar->forwardingParameters.outerHeaderCreation.value);
                if (!(outerHdr->gtpuIpv4 || outerHdr->gtpuIpv6)) {
                    UTLT_Assert(0, UpfFarRemove(tmpFar); return, "Forwarding no GTPU");
                }

                tmpFar->upfN3Teid = ntohl(outerHdr->teid);
                status = PfcpOuterHdrToIp(outerHdr, &ip);
                tmpFar->ranIp = ip;
                UTLT_Assert(status == STATUS_OK, UpfFarRemove(tmpFar); return, "Outer hdr to IP error");

                // Maybe add tunnel also be here
                UpfPdr *tmpPdr = UpfPdrFindByFarId(tmpFar->farId);
                if (tmpPdr) {
                    Status status = _AddGtpByPdrFar(tmpPdr, tmpFar);
                    UTLT_Assert(status == STATUS_OK, , "Gtp Add failed");
                }
            }
        } else {
            UTLT_Trace("destinationInterface of createFAR not precent");
        }
    } else if (tmpFar->applyAction == PFCP_FAR_APPLY_ACTION_DROP) {
        // TODO
    } else if (tmpFar->applyAction == PFCP_FAR_APPLY_ACTION_BUFF) {
        // TODO
    } else {
        UTLT_Warning("FAR Forwarding rule not support");
    }
    UTLT_Info("[PFCP] - Create %s FAR ID[0x%08x]", tmpFar->destinationInterface == PFCP_FAR_DEST_INTF_ACCESS ? "DL" : "UL", tmpFar->farId);
    ListAppend(&Self()->farList, tmpFar);
    *upfFar = tmpFar;
}

void UpfN4HandleUpdatePdr(UpdatePDR *updatePdr, UpfPdr **upfPdr) {
    UpfPdr *tmpPdr = NULL;
    UTLT_Assert(updatePdr->pDRID.presence == 1, return, "updatePDR no pdrId");

    // Find PDR
    uint16_t pdrId = ntohs(*((uint16_t*)updatePdr->pDRID.value));
    tmpPdr = UpfPdrFindByPdrId(pdrId);

    if (tmpPdr) {
        // TODO: other IE of update PDR
        if (updatePdr->precedence.presence) {
            memcpy(&tmpPdr->precedence, (uint32_t *)updatePdr->precedence.value, sizeof(uint32_t));
        }
    } else {
        /* PDR missing in updatePDR */
        UTLT_Assert(0, , "[PFCP] Session Modification Request miss of PDR context");
    }

    *upfPdr = tmpPdr;
}

void UpfN4HandleUpdateFar(UpdateFAR *updateFar, UpfFar **upfFar) {
    Status status;
    UpfFar *tmpFar = NULL;
    UTLT_Assert(updateFar->fARID.presence, return, "Far ID not presence");

    // Find FAR
    uint32_t farId = ntohl(*((uint32_t *)updateFar->fARID.value));
    tmpFar = UpfFarFindByFarId(farId);

    if (tmpFar) {
        // update ApplyAction
        if (updateFar->applyAction.presence) {
            uint8_t applyAction = *((uint8_t *)(updateFar->applyAction.value));
            // The NOCP flag may only be set if the BUFF flag is set
            if (applyAction & PFCP_FAR_APPLY_ACTION_BUFF) {
                Gtpv1TunDevNode *gtpv1Dev4 = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
                char *ifname = gtpv1Dev4->ifname;
                UpfPdr *tmpPdr = UpfPdrFindByFarId(farId);
                UTLT_Assert(tmpPdr, return, "PDR with this FAR ID not found");

                status = GtpTunnelDel(ifname, tmpPdr->upfGtpUTeid);
                UTLT_Assert(status == STATUS_OK, return, "Buffering: GTP Tunnel Remove UP error");

                char ueIpStr[INET_ADDRSTRLEN];
                inet_ntop(AF_INET, &(tmpPdr->ueIp.addr4), ueIpStr, INET_ADDRSTRLEN);

                // point oTeid to iTeid & ranIp to upfIp to buffer the packet
                status = GtpTunnelAdd(ifname, tmpPdr->upfGtpUTeid, tmpPdr->upfGtpUTeid,
                    ueIpStr, gtpv1Dev4->ip);
                UTLT_Assert(status == STATUS_OK, return, "Buffering: GTP Tunnel Add to self error");

                if (applyAction & PFCP_FAR_APPLY_ACTION_NOCP) {
                    // NOCP
                }
            } else if (applyAction & PFCP_FAR_APPLY_ACTION_FORW) {
            } else if (applyAction & PFCP_FAR_APPLY_ACTION_DROP) {
                UpfPdr *tmpPdr = UpfPdrFindByFarId(farId);
                UTLT_Assert(tmpPdr, , "pdr by farid(%d) not found", farId);
                if (tmpPdr) {
                    _RemoveGtpByPdr(tmpPdr);
                }
	    }
            tmpFar->applyAction = applyAction;
        }
        // update Forwarding parameters
        if (updateFar->updateForwardingParameters.outerHeaderCreation.value) {
            Ip ip;
            PfcpOuterHdr *outerHeader;
            outerHeader = updateFar->updateForwardingParameters.outerHeaderCreation.value;

            if (tmpFar->applyAction & PFCP_FAR_APPLY_ACTION_FORW) {
                /* Remove Old GTP Link */
                UpfPdr *tmpPdr = UpfPdrFindByFarId(farId);
                UTLT_Assert(tmpPdr, return, "PDR with this FAR ID not found");
		if (tmpPdr) {
		    _RemoveGtpByPdr(tmpPdr);
		}

                tmpFar->upfN3Teid = ntohl(outerHeader->teid);
                status = PfcpOuterHdrToIp(outerHeader, &ip);
                UTLT_Assert(status == STATUS_OK, UpfFarRemove(tmpFar); return, "Outer header IP convert error");
                tmpFar->ranIp = ip;

                /* Add new GTP Link */
                status = _AddGtpByPdr(tmpPdr);
                UTLT_Assert(status == STATUS_OK, , "Add Gtp Tunnel error");

                // Send forward packet in buffer
                UpfSession *session = UpfPdrFindByFarId(farId)->session;
                Sock *sock = ((Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList))->sock1;
                status = UpfSessionPacketSend(session, sock);
                UTLT_Assert(status == STATUS_OK, , "forward the buffer error");
            } else {
                // if drop or buffer, just update info
                tmpFar->upfN3Teid = ntohl(outerHeader->teid);
                status = PfcpOuterHdrToIp(outerHeader, &ip);
                UTLT_Assert(status == STATUS_OK, UpfFarRemove(tmpFar); return, "Outer header IP convert error");
                tmpFar->ranIp = ip;
            }
        }
        // TODO: update Duplicating parameters
        // TODO: update BAR
    } else {
        /* FAR not found */
        UTLT_Assert(0, , "[PFCP] updateFar FAR not found");
    }

    *upfFar = tmpFar;
}

Status UpfN4HandleRemovePdr(uint16_t pdrId) {
    Status status;

    UpfPdr *pdr = UpfPdrFindByPdrId(pdrId);
    UTLT_Assert(pdr != NULL, return STATUS_ERROR, "Cannot find PDR by PdrId");

    UpfSession *session = pdr->session;
    if (session != NULL) {
        if (pdr->ulDl == 1) {
            // dl (1)
	    ListRemove(&session->dlPdrList, pdr);
        } else if (pdr->ulDl == 0) {
	    // ul (0)
	    ListRemove(&session->ulPdrList, pdr);
        }
    }

    if (pdr->far != NULL && pdr->far->applyAction == PFCP_FAR_APPLY_ACTION_FORW) {
        Gtpv1TunDevNode *gtpv1Dev4 = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
        char *ifname = gtpv1Dev4->ifname;
        status = GtpTunnelDel(ifname, pdr->upfGtpUTeid);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "GTP Remove error");
    }

    status = UpfPdrRemove(pdr);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "Remove PDR error");

    return STATUS_OK;
}

Status UpfN4HandleRemoveFar(uint32_t farId) {
    Status status;

    UpfFar *far = UpfFarFindByFarId(farId);
    UTLT_Assert(far != NULL, return STATUS_ERROR, "Cannot find FAR by FarId");

    UpfPdr *pdr = UpfPdrFindByFarId(farId);
    if (pdr) {
        if (far->applyAction == PFCP_FAR_APPLY_ACTION_FORW) {
            Gtpv1TunDevNode *gtpv1Dev4 = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
            char *ifname = gtpv1Dev4->ifname;
            status = GtpTunnelDel(ifname, pdr->upfGtpUTeid);
	    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "GTP Remove error");
        }
        pdr->far = NULL;
    }

    status = UpfFarRemove(far);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "Remove FAR error");

    return STATUS_OK;
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

    /* Create FAR */
    if (request->createFAR[0].presence) {
        UpfFar *far0;
        UpfN4HandleCreateFar(&request->createFAR[0], &far0);
    }
    if (request->createFAR[1].presence) {
        UpfFar *far1;
        UpfN4HandleCreateFar(&request->createFAR[1], &far1);
    }
    if (request->createPDR[0].presence) {
        UpfPdr *pdr0;
        UpfN4HandleCreatePdr(session, &request->createPDR[0], &pdr0);
    }
    if (request->createPDR[1].presence) {
        UpfPdr *pdr1;
        UpfN4HandleCreatePdr(session, &request->createPDR[1], &pdr1);
    }

    /* Update FAR */
    if (request->updateFAR.presence) {
        UpfFar *far;
        UpfN4HandleUpdateFar(&request->updateFAR, &far);
    }
    /* Update PDR */
    if (request->updatePDR.presence) {
        UpfPdr *pdr;
        UpfN4HandleUpdatePdr(&request->updatePDR, &pdr);
    }
    /* Remove PDR */
    if (request->removePDR.presence) {
        UTLT_Assert(request->removePDR.pDRID.presence == 1,
		    , "[PFCP] PdrId in removePDR not presence!");
        UpfN4HandleRemovePdr(*(uint16_t*)request->removePDR.pDRID.value);
    }
    /* Remove FAR */
    if (request->removeFAR.presence) {
      UTLT_Assert(request->removeFAR.fARID.presence == 1,
		  ,"[PFCP] FarId in removeFAR not presence");
        UpfN4HandleRemoveFar(*(uint32_t*)request->removeFAR.fARID.value);
    }
    //TODO

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

    // Remove all Gtp Tunnel of this session
    Gtpv1TunDevNode *gtpList = (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    char *ifname = gtpList->ifname;
    UpfPdr *tmpPdr;
    for (tmpPdr = (UpfPdr*)ListFirst(&session->dlPdrList); tmpPdr; tmpPdr = (UpfPdr*)ListNext(tmpPdr)) {
        status = GtpTunnelDel(ifname, tmpPdr->upfGtpUTeid);
        UTLT_Assert(status == STATUS_OK, , "Gtp Tunnel delete failed");
    }
    // TODO: Remove all PDR
    // TODO: Remove all FAR

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
