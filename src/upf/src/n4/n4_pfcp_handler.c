#define TRACE_MODULE _n4_pfcp_handler

#include <endian.h>
#include <arpa/inet.h>
#include <net/if.h>

#include "utlt_list.h"
#include "utlt_network.h"

#include "upf_context.h"
#include "pfcp_message.h"
#include "pfcp_xact.h"
#include "pfcp_convert.h"
#include "gtp_path.h"
#include "n4_pfcp_build.h"
#include "up/up_path.h"
#include "libgtp5gnl/gtp5g.h"
#include "libgtp5gnl/gtp5gnl.h"
#include "gtp_tunnel.h"


#define _PDR_ADD 0
#define _PDR_MOD 1
#define _PDR_DEL 2
#define _FAR_ADD 0
#define _FAR_MOD 1
#define _FAR_DEL 2

#define _addr4ToStr(addrPtr, ipStr)                         \
    char ipStr[INET_ADDRSTRLEN];                            \
    inet_ntop(AF_INET, addrPtr, ipStr, INET_ADDRSTRLEN);

Status _pushPdrToKernel(struct gtp5g_pdr *pdr, int action) {
    UTLT_Assert(pdr, return STATUS_ERROR, "push PDR not found");
    Status status;

    Gtpv1TunDevNode *gtpv1Dev4 =
        (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No gtp device");
    char *ifname = gtpv1Dev4->ifname;

    uint16_t pdrId = *(uint16_t*)gtp5g_pdr_get_id(pdr);

    if (gtp5g_pdr_get_precedence(pdr)) {
        UTLT_Debug("precendence: %u", ntohl(*gtp5g_pdr_get_precedence(pdr)));
    }
    if (gtp5g_pdr_get_far_id(pdr)) {
        UTLT_Debug("farId: %u", ntohl(*gtp5g_pdr_get_far_id(pdr)));
    }
    if (gtp5g_pdr_get_outer_header_removal(pdr)) {
        UTLT_Debug("outer header removal: %u",
                   *gtp5g_pdr_get_outer_header_removal(pdr));
    }
    if (gtp5g_pdr_get_ue_addr_ipv4(pdr)) {
        _addr4ToStr(gtp5g_pdr_get_ue_addr_ipv4(pdr), ipStr);
        UTLT_Debug("ue ip: %s", ipStr);
    }
    if (gtp5g_pdr_get_local_f_teid_teid(pdr)) {
        UTLT_Debug("teid: %u", ntohl(*gtp5g_pdr_get_local_f_teid_teid(pdr)));
    }
    if (gtp5g_pdr_get_local_f_teid_gtpu_addr_ipv4(pdr)) {
        _addr4ToStr(gtp5g_pdr_get_local_f_teid_gtpu_addr_ipv4(pdr), ipStr);
        UTLT_Debug("gtpu ip: %s", ipStr);
    }

    switch (action) {
    case _PDR_ADD:
        UTLT_Debug("PDR add to kernel, dev: %s, pdr id: %u",
                   ifname, ntohs(*gtp5g_pdr_get_id(pdr)));
        status = GtpTunnelAddPdr(ifname, pdr);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Add PDR failed");
        break;
    case _PDR_MOD:
        UTLT_Debug("PDR modify to kernel, dev: %s, pdr id: %u",
                   ifname, ntohs(*gtp5g_pdr_get_id(pdr)));
        status = GtpTunnelModPdr(ifname, pdr);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modify PDR failed");
        break;
    case _PDR_DEL:
        UTLT_Debug("PDR delete to kernel, dev: %s, pdr id: %u",
                   ifname, ntohs(*gtp5g_pdr_get_id(pdr)));
        status = GtpTunnelDelPdr(ifname, pdrId);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Delete PDR failed");
        break;
    default:
        UTLT_Assert(0, return STATUS_ERROR,
                    "PDR Action %d not defined", action);
    }

    return STATUS_OK;
}

Status _pushFarToKernel(struct gtp5g_far *far, int action) {
    UTLT_Assert(far, return STATUS_ERROR, "push FAR not found");
    Status status;

    Gtpv1TunDevNode *gtpv1Dev4 =
        (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
    char *ifname = gtpv1Dev4->ifname;

    uint32_t farId = *(uint32_t*)gtp5g_far_get_id(far);

    if (gtp5g_far_get_apply_action(far)) {
        UTLT_Debug("apply action: %u", *gtp5g_far_get_apply_action(far));
    }
    if (gtp5g_far_get_outer_header_creation_description(far)) {
        UTLT_Debug("description: %u",
                   *gtp5g_far_get_outer_header_creation_description(far));
        if (gtp5g_far_get_outer_header_creation_peer_addr_ipv4(far)) {
            _addr4ToStr(gtp5g_far_get_outer_header_creation_peer_addr_ipv4(far),
                        ipStr);
            UTLT_Debug("peer ipv4: %s", ipStr);
        }
        if (gtp5g_far_get_outer_header_creation_port(far)) {
            UTLT_Debug("port: %u",
                       ntohs(*gtp5g_far_get_outer_header_creation_port(far)));
        }
    }

    switch (action) {
    case _FAR_ADD:
        UTLT_Debug("FAR add to kernel, dev: %s, far id: %u",
                   ifname, ntohl(*gtp5g_far_get_id(far)));
        status = GtpTunnelAddFar(ifname, far);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Add FAR failed");
        break;
    case _FAR_MOD:
        UTLT_Debug("FAR modify to kernel, dev: %s, far id: %u",
                   ifname, ntohl(*gtp5g_far_get_id(far)));
        status = GtpTunnelModFar(ifname, far);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modify FAR failed");
        break;
    case _FAR_DEL:
        UTLT_Debug("FAR delete to kernel, dev: %s, far id: %u",
                   ifname, ntohl(*gtp5g_far_get_id(far)));
        status = GtpTunnelDelFar(ifname, farId);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Delete FAR failed");
        break;
    default:
        UTLT_Assert(0, return STATUS_ERROR,
                    "FAR Action %d not defined", action);
    }

    return STATUS_OK;
}

Status UpfN4HandleCreatePdr(UpfSession *session, CreatePDR *createPdr) {
    UTLT_Debug("Handle Create PDR");
    UpfPdr *tmpPdr = NULL;

    UTLT_Assert(createPdr->pDRID.presence, return STATUS_ERROR,
                "pdr id not presence");
    UTLT_Assert(createPdr->precedence.presence, return STATUS_ERROR,
                "precedence not presence");
    UTLT_Assert(createPdr->pDI.presence, return STATUS_ERROR,
                "Pdi not exist");
    UTLT_Assert(createPdr->pDI.sourceInterface.presence,
                return STATUS_ERROR, "PDI SourceInterface not presence");

    tmpPdr = gtp5g_pdr_alloc();
    UTLT_Assert(tmpPdr, return STATUS_ERROR, "pdr allocate error");

    // PdrId
    uint16_t pdrId = *((uint16_t *)createPdr->pDRID.value);
    UTLT_Debug("PDR ID: %u", ntohs(pdrId));
    gtp5g_pdr_set_id(tmpPdr, pdrId);
    gtp5g_pdr_set_unix_sock_path(tmpPdr, Self()->buffSockPath);

    // precedence
    uint32_t precedence = *((uint32_t *)createPdr->precedence.value);
    gtp5g_pdr_set_precedence(tmpPdr, precedence);

    // source interface
    //uint8_t sourceInterface =
    //    *((uint8_t *)(createPdr->pDI.sourceInterface.value));

    // F-TEID
    if (createPdr->pDI.localFTEID.presence) {
        PfcpFTeid *fTeid = (PfcpFTeid*)createPdr->pDI.localFTEID.value;
        uint32_t teid = fTeid->teid;

        if (fTeid->v4 && fTeid->v6) {
            // TODO: Dual Stack
        } else if (fTeid->v4) {
            gtp5g_pdr_set_local_f_teid(tmpPdr, teid, &(fTeid->addr4));
        } else if (fTeid->v6) {
            // TODO: ipv6
            //gtp5g_pdr_set_local_f_teid(tmpPdr, teid, &(fTeid->addr6));
        }
    }

    // UE IP
    if (createPdr->pDI.uEIPAddress.presence) {
        PfcpUeIpAddr *ueIp =
            (PfcpUeIpAddr*)createPdr->pDI.uEIPAddress.value;
        if (ueIp->v4 && ueIp->v6) {
            // TODO: Dual Stack
        } else if (ueIp->v4) {
            gtp5g_pdr_set_ue_addr_ipv4(tmpPdr, &(ueIp->addr4));
        } else if (ueIp->v6) {
            // TODO: IPv6
        }
    }

    // PDI SDF filter
    uint16_t flowDescriptionLen = 0;
    char *flowDescription = NULL;
    if (createPdr->pDI.sDFFilter.presence) {
        // Decode SDF
        PfcpSDFFilterDescription des =
            *(PfcpSDFFilterDescription*)createPdr->pDI.sDFFilter.value;

        for (size_t idx = 2; idx < createPdr->pDI.sDFFilter.len;
        /* Do nothing here */) {
            if (des.fd) {
                flowDescriptionLen =
                    *(uint16_t*)((uint8_t*)createPdr->pDI.sDFFilter.value+idx);

                flowDescription =
                    UTLT_Calloc(flowDescriptionLen + 1, sizeof(uint8_t));
                UTLT_Assert(flowDescription, idx += 2 + flowDescriptionLen; continue,
                            "flow description allocate error");

                memcpy(flowDescription,
                       (uint8_t*)createPdr->pDI.sDFFilter.value+idx+2,
                       flowDescriptionLen);

                gtp5g_pdr_set_sdf_filter_description(tmpPdr, flowDescription);

                idx += 2 + flowDescriptionLen;
            } else if (des.ttc) {
                /*
                uint16_t tosTrafficClass =
                    *(uint16_t*)((uint8_t*)createPdr->pDI.sDFFilter.value+idx);
                */
                UTLT_Warning("SDF ToS traffic class not implemented");
                idx += 2;
            } else if (des.spi) {
                /*
                uint32_t securityParameterIndex =
                    *(uint32_t*)((uint8_t*)createPdr->pDI.sDFFilter.value+idx);
                */
                UTLT_Warning("SDF security paramenter index not implemented");
                idx += 4;
            } else if (des.fl) {
                // TODO: Flow Label
                UTLT_Warning("SDF flow label not implemented");
                idx += 3;
            } else if (des.bid) {
                /*
                uint32_t sDFFilterId =
                    *(uint32_t*)((uint8_t*)createPdr->pDI.sDFFilter.value+idx);
                */
                UTLT_Warning("SDF filter id not implemented now");
                idx += 4;
            } else {
                UTLT_Warning("Other tag not implements");
                idx++;
            }
        }
    }

    // Outer Header Removal
    if (createPdr->outerHeaderRemoval.presence) {
        uint8_t outerHeader =
            *(uint8_t*)createPdr->outerHeaderRemoval.value;
        gtp5g_pdr_set_outer_header_removal(tmpPdr, outerHeader);
    }

    // FAR ID
    if (createPdr->fARID.presence) {
        uint32_t farId = *((uint32_t *)createPdr->fARID.value);
        gtp5g_pdr_set_far_id(tmpPdr, farId);
    }

    // Send PDR to kernel
    Status status = _pushPdrToKernel(tmpPdr, _PDR_ADD);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PDR not pushed to kernel");
    gtp5g_pdr_free(tmpPdr);
    UTLT_Assert(tmpPdr != NULL, return STATUS_ERROR,
                "Free PDR struct error");
    if (flowDescription) {
        UTLT_Assert(UTLT_Free(flowDescription) == STATUS_OK, ,
                    "Free flow description error");
    }

    // Set session point to pdr
    UpfPdrId *pdrIdPtr = UpfPdrIdAdd(pdrId);
    UTLT_Assert(pdrIdPtr, return STATUS_ERROR, "PdrId Add error");
    ListAppend(&session->pdrIdList, pdrIdPtr);
    // Set buff relate pdr to session
    UpfBufPacketAdd(session, pdrId);

    return STATUS_OK;
}

Status UpfN4HandleCreateFar(CreateFAR *createFar) {
    UTLT_Debug("Handle Create FAR");
    UpfFar *tmpFar = NULL;
    UTLT_Assert(createFar->fARID.presence, return STATUS_ERROR,
                "Far ID not presence");
    UTLT_Assert(createFar->applyAction.presence,
                return STATUS_ERROR, "Apply Action not presence");

    // Create FAR
    tmpFar = gtp5g_far_alloc();
    UTLT_Assert(tmpFar, return STATUS_ERROR, "FAR allocate error");

    // FarId
    uint32_t farId = *((uint32_t *)createFar->fARID.value);
    UTLT_Debug("FAR ID: %u", ntohl(farId));
    gtp5g_far_set_id(tmpFar, farId);

    // Apply Action
    uint8_t applyAction = *((uint8_t *)(createFar->applyAction.value));
    gtp5g_far_set_apply_action(tmpFar, applyAction);

    // Forwarding Parameters
    if (createFar->forwardingParameters.presence) {
        // Destination Interface
        /*
          if (createFar->forwardingParameters.destinationInterface.presence) {
          uint8_t destinationInterface =
          *((uint8_t *)(createFar->forwardingParameters.destinationInterface.value));
          }
          // Network Instance
          if (createFar->forwardingParameters.networkInstance.presence) {
          }
        */
        // Outer Header Creation
        if (createFar->forwardingParameters.outerHeaderCreation.presence) {
            PfcpOuterHdr *outerHdr = (PfcpOuterHdr *)
                (createFar->forwardingParameters.outerHeaderCreation.value);
            uint16_t description = *((uint16_t *)outerHdr);

            if (outerHdr->gtpuIpv4 || outerHdr->udpIpv4) {
                gtp5g_far_set_outer_header_creation(tmpFar, description,
                                                    outerHdr->teid,
                                                    &(outerHdr->addr4),
                                                    htons(2152));
            } else if (outerHdr->udpIpv4) {
                // only with UDP enabled has port number
                gtp5g_far_set_outer_header_creation(tmpFar, description,
                                                    0, &(outerHdr->addr4),
                                                    outerHdr->port);
            }
        }
    }

    // Send FAR to kernel
    Status status = _pushFarToKernel(tmpFar, _FAR_ADD);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "FAR not pushed to kernel");
    gtp5g_far_free(tmpFar);
    UTLT_Assert(tmpFar != NULL, return STATUS_ERROR,
                "Free FAR struct error");

    return STATUS_OK;
}

Status UpfN4HandleUpdatePdr(UpdatePDR *updatePdr) {
    UTLT_Debug("Handle Update PDR");
    UpfPdr *tmpPdr = NULL;
    UTLT_Assert(updatePdr->pDRID.presence == 1,
                return STATUS_ERROR, "updatePDR no pdrId");

    // Find PDR
    uint16_t pdrId = *((uint16_t*)updatePdr->pDRID.value);

    /*
    Gtpv1TunDevNode *gtpv1Dev4 =
        (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
    tmpPdr = GtpTunnelFindPdrById(gtpv1Dev4->ifname, pdrId);
    UTLT_Assert(tmpPdr, return STATUS_ERROR,
                "[PFCP] UpdatePDR PDR[%u] not found", ntohs(pdrId));
    */
    tmpPdr = gtp5g_pdr_alloc();
    UTLT_Assert(tmpPdr, return STATUS_ERROR, "pdr alloc error");
    gtp5g_pdr_set_id(tmpPdr, pdrId);

    // TODO: other IE of update PDR
    if (updatePdr->outerHeaderRemoval.presence) {
        gtp5g_pdr_set_outer_header_removal(tmpPdr,
            *((uint8_t*)(updatePdr->outerHeaderRemoval.value)));
    }
    if (updatePdr->precedence.presence) {
        gtp5g_pdr_set_precedence(tmpPdr,
                                 *((uint32_t *)(updatePdr->precedence.value)));
    }
    if (updatePdr->pDI.presence) {
        if (updatePdr->pDI.localFTEID.presence) {
            PfcpFTeid *fTeid = (PfcpFTeid*)updatePdr->pDI.localFTEID.value;
            uint32_t teid = fTeid->teid;

            if (fTeid->v4 && fTeid->v6) {
                // TODO: Dual Stack
            } else if (fTeid->v4) {
                gtp5g_pdr_set_local_f_teid(tmpPdr, teid, &(fTeid->addr4));
            } else if (fTeid->v6) {
                // TODO: ipv6
                //gtp5g_pdr_set_local_f_teid(tmpPdr, teid, &(fTeid->addr6));
            }
        }
        if (updatePdr->pDI.uEIPAddress.presence) {
            PfcpUeIpAddr *ueIp =
                (PfcpUeIpAddr*)updatePdr->pDI.uEIPAddress.value;
            if (ueIp->v4 && ueIp->v6) {
                // TODO: Dual Stack
            } else if (ueIp->v4) {
                gtp5g_pdr_set_ue_addr_ipv4(tmpPdr, &(ueIp->addr4));
            } else if (ueIp->v6) {
                // TODO: IPv6
            }
        }
    }
    if (updatePdr->fARID.presence) {
        gtp5g_pdr_set_far_id(tmpPdr, *(uint32_t *)updatePdr->fARID.value);
    }

    // update PDR to kernel
    Status status = _pushPdrToKernel(tmpPdr, _PDR_MOD);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PDR not pushed to kernel");
    gtp5g_pdr_free(tmpPdr);
    UTLT_Assert(tmpPdr != NULL, return STATUS_ERROR,
                "Free PDR struct error");

    return STATUS_OK;
}

Status UpfN4HandleUpdateFar(UpdateFAR *updateFar) {
    UTLT_Debug("Handle Update FAR");
    UpfFar *tmpFar = NULL;
    UTLT_Assert(updateFar->fARID.presence,
                return STATUS_ERROR, "Far ID not presence");

    // Record if need to send buffer
    uint8_t sendBufMarker = 0;
    // Find FAR
    uint32_t farId = *((uint32_t *)updateFar->fARID.value);

    Gtpv1TunDevNode *gtpv1Dev4 =
        (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
    char *gtpIfname = gtpv1Dev4->ifname;

    UpfFar *oldFar = GtpTunnelFindFarById(gtpIfname, farId);
    UTLT_Assert(oldFar, return STATUS_ERROR,
                "[PFCP] UpdateFAR FAR[%u] not found", ntohl(farId));

    // Create been updated far
    tmpFar = gtp5g_far_alloc();
    UTLT_Assert(tmpFar, return STATUS_ERROR, "FAR allocate error");
    gtp5g_far_set_id(tmpFar, farId);

    // update Apply Action
    if (updateFar->applyAction.presence) {
        uint8_t applyAction = *(uint8_t*)updateFar->applyAction.value;
        gtp5g_far_set_apply_action(tmpFar, applyAction);
        if ((applyAction & PFCP_FAR_APPLY_ACTION_FORW)) {
            //(*gtp5g_far_get_apply_action(oldFar) & PFCP_FAR_APPLY_ACTION_BUFF)) {
            sendBufMarker = 1;
        }
    }
    // update Forwarding parameters
    if (updateFar->updateForwardingParameters.outerHeaderCreation.value) {
        PfcpOuterHdr *outerHdr = (PfcpOuterHdr *)
            (updateFar->updateForwardingParameters.outerHeaderCreation.value);
        uint16_t description = *((uint16_t *)outerHdr);

        if (outerHdr->gtpuIpv4) {
            gtp5g_far_set_outer_header_creation(tmpFar, description,
                                                outerHdr->teid,
                                                &(outerHdr->addr4),
                                                htons(2152));
        } else if (outerHdr->udpIpv4) {
            // TODO: Check if no teid, which should put
            gtp5g_far_set_outer_header_creation(tmpFar, description,
                                                0, &(outerHdr->addr4),
                                                outerHdr->port);
        }
    }
    // TODO: update Duplicating parameters
    // TODO: update BAR

    // TODO: update FAR to kernel
    Status status = _pushFarToKernel(tmpFar, _FAR_MOD);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "FAR not pushed to kernel");
    gtp5g_far_free(tmpFar);
    UTLT_Assert(tmpFar != NULL, return STATUS_ERROR,
                "Free FAR struct error");
    gtp5g_far_free(oldFar);
    UTLT_Assert(oldFar != NULL, return STATUS_ERROR,
                "Free FAR struct error");

    // Send buffer if need
    // Send here because I'm not sure updateFar will update forwarding
    // parameter or not
    if (sendBufMarker) {
        UpfFar *newFar = GtpTunnelFindFarById(gtpIfname, farId);
        UTLT_Assert(newFar, return STATUS_ERROR,
                    "[PFCP] UpdateFAR FAR[%u] not found", ntohl(farId));
        // Check if pdr Packet full of packet
        int pdrNum = *gtp5g_far_get_related_pdr_num(newFar);
        uint16_t *pdrIdList = gtp5g_far_get_related_pdr_list(newFar);
        // Send buffer back to Gtp Dev
        Sock *sock = ((Gtpv1TunDevNode*)
                        ListFirst(&Self()->gtpv1DevList))->sock;
        sock->remoteAddr._family = sock->localAddr._family;
        sock->remoteAddr._port = sock->localAddr._port;
        if (sock->localAddr._family == AF_INET) {
            sock->remoteAddr.s4.sin_addr =
                *gtp5g_far_get_outer_header_creation_peer_addr_ipv4(newFar);
        } else {
            // TODO: IPv6
        }

        for (size_t idx = 0; idx < pdrNum; ++idx) {
            UpfPdr *tmpPdr = GtpTunnelFindPdrById(gtpIfname, pdrIdList[idx]);
            UpSendPacketByPdrFar(tmpPdr, newFar, sock);
            gtp5g_pdr_free(tmpPdr);
            UTLT_Assert(tmpPdr != NULL, return STATUS_ERROR,
                        "Free PDR struct error");
        }
        // FIXME: free pdrIdList

        gtp5g_far_free(newFar);
        UTLT_Assert(newFar != NULL, return STATUS_ERROR,
                    "Free FAR struct error");
    }

    return STATUS_OK;
}

Status UpfN4HandleRemovePdr(UpfSession *session, uint16_t pdrId) {
    UTLT_Debug("Handle Remove PDR[%u]", ntohs(pdrId));
    UTLT_Assert(ntohs(pdrId), return STATUS_ERROR, "pdrId cannot be 0");
    UTLT_Assert(session, return STATUS_ERROR,
                "session not found");

    UpfPdrId *sessionPdrIdPtr = ListFirst(&session->pdrIdList);
    while (sessionPdrIdPtr) {
        if (sessionPdrIdPtr->pdrId == pdrId) {
            Gtpv1TunDevNode *gtpv1Dev4 =
                (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
            UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
            Status status = GtpTunnelDelPdr(gtpv1Dev4->ifname, pdrId);
            UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                        "PDR[%u] delete failed", ntohs(pdrId));

            // Remove PDR ID from session
            ListRemove(&session->pdrIdList, sessionPdrIdPtr);
            UpfPdrIdRemove(sessionPdrIdPtr);
            // Remove buff
            UpfBufPacket *tmpBufPacket = UpfBufPacketFindByPdrId(pdrId);
            UpfBufPacketRemove(tmpBufPacket);
            return STATUS_OK;
        }

        sessionPdrIdPtr = (UpfPdrId *)ListNext(sessionPdrIdPtr);
    }

    UTLT_Warning("PDR[%u] not in this session, PDR not removed", ntohs(pdrId));
    return STATUS_ERROR;
}

Status UpfN4HandleRemoveFar(uint32_t farId) {
    UTLT_Debug("Handle Remove FAR[%u]", ntohl(farId));
    UTLT_Assert(ntohl(farId), return STATUS_ERROR,
                "farId should not be 0");

    // TODO: here can be speedup like
    //UpfPdr *pdr = GtpTunnelFindPdrByFarId(gtp5g_int_name, farId);
    //if (pdr) {
    //    gtp5g_pdr_set_far_id(pdr, 0);
    //}
    //Status status = GtpTunnelDelFar(gtp5g_int_name, farId);
    //UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
    //            "FAR delete error");

    Gtpv1TunDevNode *gtpv1Dev4 =
        (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
    char *gtpIfName = gtpv1Dev4->ifname;
    UpfFar *far = GtpTunnelFindFarById(gtpIfName, farId);
    UTLT_Assert(far != NULL, return STATUS_ERROR,
                "Cannot find FAR[%u] by FarId", ntohl(farId));

    // Set FarId to 0 if the PDR has this far
    int pdrNum = *(int*)gtp5g_far_get_related_pdr_num(far);
    uint16_t *pdrList = gtp5g_far_get_related_pdr_list(far);
    for (size_t idx = 0; idx < pdrNum; ++idx) {
        UpfPdr * tmpPdr = GtpTunnelFindPdrById(gtpIfName, pdrList[idx]);
        gtp5g_pdr_set_far_id(tmpPdr, 0);
        _pushPdrToKernel(tmpPdr, _PDR_MOD);

        gtp5g_pdr_free(tmpPdr);
        UTLT_Assert(tmpPdr != NULL, continue,
                    "Free pdr error");
    }
    // FIXME: free pdrList

    Status status = _pushFarToKernel(far, _FAR_DEL);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "FAR not pushed to kernel");
    gtp5g_far_free(far);
    UTLT_Assert(far != NULL, return STATUS_ERROR,
                "Free FAR error");

    return STATUS_OK;
}

Status UpfN4HandleSessionEstablishmentRequest(UpfSession *session, PfcpXact *pfcpXact,
                                              PFCPSessionEstablishmentRequest *request) {
    Status status;
    uint8_t cause = PFCP_CAUSE_REQUEST_ACCEPTED;

    UTLT_Assert(session, return STATUS_ERROR, "Upf Session error");
    UTLT_Assert(pfcpXact, return STATUS_ERROR, "pfcpXact error");
    //UTLT_Assert(pfcpXact->gtpBuf, return,
    //  "GTP buffer of pfcpXact error");
    //UTLT_Assert(pfcpXact->gtpXact, return,
    // "GTP Xact of pfcpXact error");

    /* First FAR */
    if (request->createFAR[0].presence) {
        status = UpfN4HandleCreateFar(&request->createFAR[0]);
        // TODO: if error, which cause, and pull out the rule from kernel that
        // has been set, maybe need to pull out session as well
        UTLT_Assert(status == STATUS_OK, cause = PFCP_CAUSE_REQUEST_REJECTED,
                    "Create FAR error");
    }
    if (request->createPDR[1].presence) {
        status = UpfN4HandleCreateFar(&request->createFAR[1]);
        UTLT_Assert(status == STATUS_OK, cause = PFCP_CAUSE_REQUEST_REJECTED,
                    "Create FAR error");
    }
    /* Then PDR (order is important) */
    if (request->createPDR[0].presence) {
        status = UpfN4HandleCreatePdr(session, &request->createPDR[0]);
        UTLT_Assert(status == STATUS_OK, cause = PFCP_CAUSE_REQUEST_REJECTED,
                    "Create PDR Error");
    }
    if (request->createPDR[1].presence) {
        status = UpfN4HandleCreatePdr(session, &request->createPDR[1]);
        UTLT_Assert(status == STATUS_OK, cause = PFCP_CAUSE_REQUEST_REJECTED,
                    "Create PDR 2 Error");
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

    PfcpHeader header;
    Bufblk *bufBlk = NULL;
    PfcpFSeid *smfFSeid = NULL;

    if (!request->cPFSEID.presence) {
        UTLT_Error("Session Establishment Response: No CP F-SEID");
        cause = PFCP_CAUSE_MANDATORY_IE_MISSING;
    }

    smfFSeid = request->cPFSEID.value;
    session->smfSeid = be64toh(smfFSeid->seid);

    /* Send Response */
    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_SESSION_ESTABLISHMENT_RESPONSE;
    header.seid = session->smfSeid;

    status = UpfN4BuildSessionEstablishmentResponse(&bufBlk, header.type,
                                                    session, cause, request);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "N4 build error");

    status = PfcpXactUpdateTx(pfcpXact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "pfcpXact update TX error");

    status = PfcpXactCommit(pfcpXact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "xact commit error");

    UTLT_Info("[PFCP] Session Establishment Response");
    return STATUS_OK;
}

Status UpfN4HandleSessionModificationRequest(UpfSession *session, PfcpXact *xact,
                                             PFCPSessionModificationRequest *request) {
    UTLT_Assert(session, return STATUS_ERROR, "Session error");
    UTLT_Assert(xact, return STATUS_ERROR, "xact error");

    Status status;
    PfcpHeader header;
    Bufblk *bufBlk;

    /* Create FAR */
    if (request->createFAR[0].presence) {
        status = UpfN4HandleCreateFar(&request->createFAR[0]);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Create FAR error");
    }
    if (request->createFAR[1].presence) {
        status = UpfN4HandleCreateFar(&request->createFAR[1]);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Create FAR2 error");
    }
    /* Create PDR */
    if (request->createPDR[0].presence) {
        status = UpfN4HandleCreatePdr(session, &request->createPDR[0]);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Create PDR error");
    }
    if (request->createPDR[1].presence) {
        status = UpfN4HandleCreatePdr(session, &request->createPDR[1]);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Create PDR2 error");
    }

    /* Update FAR */
    if (request->updateFAR.presence) {
        UTLT_Assert(request->updateFAR.fARID.presence == 1, ,
                    "[PFCP] FarId in updateFAR not presence");
        status = UpfN4HandleUpdateFar(&request->updateFAR);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Update FAR error");
    }
    /* Update PDR */
    if (request->updatePDR.presence) {
        UTLT_Assert(request->updatePDR.pDRID.presence == 1, ,
                    "[PFCP] PdrId in updatePDR not presence!");
        status = UpfN4HandleUpdatePdr(&request->updatePDR);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Update PDR error");
    }

    /* Remove PDR */
    if (request->removePDR.presence) {
        UTLT_Assert(request->removePDR.pDRID.presence == 1, ,
                    "[PFCP] PdrId in removePDR not presence!");
        status = UpfN4HandleRemovePdr(session, *(uint16_t*)
                                      request->removePDR.pDRID.value);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Remove PDR error");
    }
    /* Remove FAR */
    if (request->removeFAR.presence) {
        UTLT_Assert(request->removeFAR.fARID.presence == 1, ,
                    "[PFCP] FarId in removeFAR not presence");
        status = UpfN4HandleRemoveFar(*(uint32_t*)
                                      request->removeFAR.fARID.value);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Modification: Remove FAR error");
    }

    /* Send Session Modification Response */
    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_SESSION_MODIFICATION_RESPONSE;
    header.seid = session->smfSeid;

    status = UpfN4BuildSessionModificationResponse(&bufBlk, header.type,
                                                   session, request);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PFCP Commit error");

    UTLT_Info("[PFCP] Session Modification Response");
    return STATUS_OK;
}

Status UpfN4HandleSessionDeletionRequest(UpfSession *session, PfcpXact *xact,
                                         PFCPSessionDeletionRequest *request) {
    UTLT_Assert(session, return STATUS_ERROR, "session error");
    UTLT_Assert(xact, return STATUS_ERROR, "xact error");

    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;

    // Remove all PDR and FAR
    // PDR will not overlap between session so just remove it
    Gtpv1TunDevNode *gtpv1Dev4 =
        (Gtpv1TunDevNode*)ListFirst(&Self()->gtpv1DevList);
    char *ifname = gtpv1Dev4->ifname;

    UTLT_Assert(gtpv1Dev4, return STATUS_ERROR, "No GTP Device");
    uint16_t pdrId;
    UpfPdrId *pdrIdPtr;
    // Always get first one because the first one before have been deleted
    while((pdrIdPtr = ListFirst(&session->pdrIdList))) {
        pdrId = pdrIdPtr->pdrId;
        // Remove PDR before far, but save farId first
        UpfPdr *tmpPdr = GtpTunnelFindPdrById(ifname, pdrId);
        uint32_t farId = *gtp5g_pdr_get_far_id(tmpPdr);

        status = GtpTunnelDelPdr(gtpv1Dev4->ifname, pdrId);
        UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                    "Remove PDR[%u] error", pdrId);
        gtp5g_pdr_free(tmpPdr);

        // Remove FAR
        status = GtpTunnelDelFar(ifname, farId);
        if (status != STATUS_OK) {
            // status not important
            UTLT_Debug("Remove FAR[%u] error, "
                       "but it may be n PDR point to same FAR", farId);
        }
        ListRemove(&session->pdrIdList, pdrIdPtr);
        UpfPdrIdRemove(pdrIdPtr);
    }

    /* delete session */
    UpfSessionRemove(session);

    /* Send Session Deletion Response */
    memset(&header, 0, sizeof(PfcpHeader));

    header.type = PFCP_SESSION_DELETION_RESPONSE;
    header.seid = session->smfSeid;

    status = UpfN4BuildSessionDeletionResponse(&bufBlk, header.type,
                                               session, request);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "xact commit error");

    UTLT_Info("[PFCP] Session Deletion Response");
    return STATUS_OK;
}

Status UpfN4HandleSessionReportResponse(UpfSession *session, PfcpXact *xact,
                                        PFCPSessionReportResponse *response) {
    Status status;

    UTLT_Assert(session, return STATUS_ERROR, "Session error");
    UTLT_Assert(xact, return STATUS_ERROR, "xact error");
    UTLT_Assert(response->cause.presence, return STATUS_ERROR,
                "SessionReportResponse error: no Cause");

    // TODO: check if need update TX

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "xact commit error");

    UTLT_Info("[PFCP] Session Report Response");
    return STATUS_OK;
}

Status UpfN4HandleAssociationSetupRequest(PfcpXact *xact,
                                          PFCPAssociationSetupRequest *request) {
    PfcpNodeId *nodeId;

    UTLT_Assert(xact, return STATUS_ERROR, "xact error");
    UTLT_Assert(xact->gnode, return STATUS_ERROR,
                "gNode of xact error");
    UTLT_Assert(request->nodeID.presence, return STATUS_ERROR,
                "Request missing nodeId");

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
        UTLT_Assert(0, return STATUS_ERROR,
                    "Request no node id type");
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
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "xact commit error");

    UTLT_Info("[PFCP] Association Setup Response");
    return STATUS_OK;
}

Status UpfN4HandleAssociationUpdateRequest(PfcpXact *xact,
                                           PFCPAssociationUpdateRequest *request) {
    // TODO
    UTLT_Info("[PFCP] TODO Association Update Request");
    return STATUS_OK;
}

Status UpfN4HandleAssociationReleaseRequest(PfcpXact *xact,
                                            PFCPAssociationReleaseRequest *request) {
    UTLT_Assert(xact, return STATUS_ERROR, "xact error");
    UTLT_Assert(xact->gnode, return STATUS_ERROR,
                "gNode of xact error");
    UTLT_Assert(request->nodeID.presence, return STATUS_ERROR,
                "Request missing nodeId");

    // Clear all session releated to this node
    HashIndex *sessionHashIdx = NULL;
    UpfSession *session = NULL;

    for (sessionHashIdx = UpfSessionFirst(); sessionHashIdx;
         sessionHashIdx = UpfSessionNext(sessionHashIdx)) {
        session = UpfSessionThis(sessionHashIdx);
        // Clear transaction node
        if (session->pfcpNode == xact->gnode) {
            UpfSessionRemove(session);
        }
    }
    // TODO: Check if I need to remove gnode in transaction

    // Build Response
    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;

    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_ASSOCIATION_RELEASE_RESPONSE;
    header.seid = 0;

    status = UpfN4BuildAssociationReleaseResponse(&bufBlk, header.type);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "xact commit error");

    UTLT_Info("[PFCP] Association Release Request");
    return STATUS_OK;
}

Status UpfN4HandleHeartbeatRequest(PfcpXact *xact, HeartbeatRequest *request) {
    Status status;
    PfcpHeader header;
    Bufblk *bufBlk = NULL;

    UTLT_Info("[PFCP] Heartbeat Request");

    /* Send */
    memset(&header, 0, sizeof(PfcpHeader));
    header.type = PFCP_HEARTBEAT_RESPONSE;
    header.seid = 0;

    status = UpfN4BuildHeartbeatResponse(&bufBlk, header.type);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "N4 build error");

    status = PfcpXactUpdateTx(xact, &header, bufBlk);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "PfcpXactUpdateTx error");

    status = PfcpXactCommit(xact);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "xact commit error");

    UTLT_Info("[PFCP] Heartbeat Response");
    return STATUS_OK;
}

Status UpfN4HandleHeartbeatResponse(PfcpXact *xact,
                                    HeartbeatResponse *response) {
    // if rsv response, nothing to do, else peer may be not alive
    UTLT_Info("[PFCP] Heartbeat Response");
    return STATUS_OK;
}
