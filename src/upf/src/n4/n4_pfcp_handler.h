#ifndef __N4_PFCP_HANDLER_H__
#define __N4_PFCP_HANDLER_H__

#include "upf_context.h"
#include "pfcp_message.h"
#include "pfcp_xact.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

void UpfN4HandleCreatePdr(UpfSession *session, CreatePDR *createPdr);
void UpfN4HandleCreateFar(CreateFAR *createFar);
void UpfN4HandleUpdatePdr(UpdatePDR *updatePdr);
void UpfN4HandleUpdateFar(UpdateFAR *updateFar);
Status UpfN4HandleRemovePdr(uint16_t pdrId);
Status UpfN4HandleRemoveFar(uint32_t farId);
void UpfN4HandleSessionEstablishmentRequest(
        UpfSession *session, PfcpXact *pfcpXact, PFCPSessionEstablishmentRequest *request);
void UpfN4HandleSessionModificationRequest(
        UpfSession *session, PfcpXact *xact, PFCPSessionModificationRequest *request);
void UpfN4HandleSessionDeletionRequest(UpfSession *session, PfcpXact *xact, PFCPSessionDeletionRequest *request);
void UpfN4HandleSessionReportResponse(
        UpfSession *session, PfcpXact *xact, PFCPSessionReportResponse *response);
void UpfN4HandleAssociationSetupRequest(PfcpXact *xact, PFCPAssociationSetupRequest *request);
void UpfN4HandleAssociationUpdateRequest(PfcpXact *xact, PFCPAssociationUpdateRequest *request);
void UpfN4HandleAssociationReleaseRequest(PfcpXact *xact, PFCPAssociationReleaseRequest *request);
void UpfN4HandleHeartbeatRequest(PfcpXact *xact, HeartbeatRequest *request);
void UpfN4HandleHeartbeatResponse(PfcpXact *xact, HeartbeatResponse *response);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __N4_PFCP_HANDLER_H__ */
