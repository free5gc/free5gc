package handler

import (
	"fmt"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/pfcp/pfcpUdp"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
	pfcp_message "github.com/free5gc/smf/internal/pfcp/message"
	"github.com/free5gc/smf/pkg/service"
)

func HandlePfcpHeartbeatRequest(msg *pfcpUdp.Message) {
	h := msg.PfcpMessage.Header
	pfcp_message.SendHeartbeatResponse(msg.RemoteAddr, h.SequenceNumber)
}

func HandlePfcpPfdManagementRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP PFD Management Request handling is not implemented")
}

func HandlePfcpAssociationSetupRequest(msg *pfcpUdp.Message) {
	req := msg.PfcpMessage.Body.(pfcp.PFCPAssociationSetupRequest)

	nodeID := req.NodeID
	if nodeID == nil {
		logger.PfcpLog.Errorln("pfcp association needs NodeID")
		return
	}
	logger.PfcpLog.Infof("Handle PFCP Association Setup Request with NodeID[%s]",
		nodeID.ResolveNodeIdToIp().String())

	upf := smf_context.RetrieveUPFNodeByNodeID(*nodeID)
	if upf == nil {
		logger.PfcpLog.Errorf("can't find UPF[%s]", nodeID.ResolveNodeIdToIp().String())
		return
	}

	// Response with PFCP Association Setup Response
	cause := pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}
	pfcp_message.SendPfcpAssociationSetupResponse(msg.RemoteAddr, cause)
}

func HandlePfcpAssociationUpdateRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Association Update Request handling is not implemented")
}

func HandlePfcpAssociationReleaseRequest(msg *pfcpUdp.Message) {
	pfcpMsg := msg.PfcpMessage.Body.(pfcp.PFCPAssociationReleaseRequest)

	var cause pfcpType.Cause
	upf := smf_context.RetrieveUPFNodeByNodeID(*pfcpMsg.NodeID)

	if upf != nil {
		smf_context.RemoveUPFNodeByNodeID(*pfcpMsg.NodeID)
		cause.CauseValue = pfcpType.CauseRequestAccepted
	} else {
		cause.CauseValue = pfcpType.CauseNoEstablishedPfcpAssociation
	}

	pfcp_message.SendPfcpAssociationReleaseResponse(msg.RemoteAddr, cause)
}

func HandlePfcpNodeReportRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Node Report Request handling is not implemented")
}

func HandlePfcpSessionSetDeletionRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Set Deletion Request handling is not implemented")
}

func HandlePfcpSessionSetDeletionResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Set Deletion Response handling is not implemented")
}

func HandlePfcpSessionReportRequest(msg *pfcpUdp.Message) {
	var cause pfcpType.Cause

	req := msg.PfcpMessage.Body.(pfcp.PFCPSessionReportRequest)
	SEID := msg.PfcpMessage.Header.SEID
	smContext := smf_context.GetSMContextBySEID(SEID)
	seqFromUPF := msg.PfcpMessage.Header.SequenceNumber

	if smContext == nil {
		logger.PfcpLog.Errorf("PFCP Session SEID[%d] not found", SEID)
		cause.CauseValue = pfcpType.CauseSessionContextNotFound
		pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, 0)
		return
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	upfNodeID := smContext.GetNodeIDByLocalSEID(SEID)
	upfNodeIDtoIP := upfNodeID.ResolveNodeIdToIp()
	if upfNodeIDtoIP.IsUnspecified() {
		logger.PduSessLog.Errorf("Invalid PFCP Session Report Request : no PFCP session found with SEID %d", SEID)
		cause.CauseValue = pfcpType.CauseNoEstablishedPfcpAssociation
		pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, 0)
		return
	}
	upfNodeIDtoIPStr := upfNodeIDtoIP.String()

	pfcpCtx := smContext.PFCPContext[upfNodeIDtoIPStr]
	if pfcpCtx == nil {
		logger.PfcpLog.Errorf("pfcpCtx [nodeId: %v, seid:%d] not found", upfNodeID, SEID)
		cause.CauseValue = pfcpType.CauseNoEstablishedPfcpAssociation
		pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, 0)
		return
	}
	remoteSEID := pfcpCtx.RemoteSEID

	upf := smf_context.RetrieveUPFNodeByNodeID(upfNodeID)
	if upf == nil {
		logger.PfcpLog.Errorf("can't find UPF[%s]", upfNodeIDtoIPStr)
		cause.CauseValue = pfcpType.CauseNoEstablishedPfcpAssociation
		pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, 0)
		return
	}
	if err := upf.IsAssociated(); err != nil {
		logger.PfcpLog.Warnf("PFCP Session Report Request rejected: %+v", err)
		cause.CauseValue = pfcpType.CauseNoEstablishedPfcpAssociation
		pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, 0)
	}

	if smContext.UpCnxState == models.UpCnxState_DEACTIVATED {
		if req.ReportType.Dldr {
			downlinkDataReport := req.DownlinkDataReport

			if downlinkDataReport.DownlinkDataServiceInformation != nil {
				logger.PfcpLog.Warnf(
					"PFCP Session Report Request DownlinkDataServiceInformation handling is not implemented")
			}

			n1n2Request := models.N1N2MessageTransferRequest{}

			// TS 23.502 4.2.3.3 3a. Send Namf_Communication_N1N2MessageTransfer Request, SMF->AMF
			if n2SmBuf, err := smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext); err != nil {
				logger.PduSessLog.Errorln("Build PDUSessionResourceSetupRequestTransfer failed:", err)
			} else {
				n1n2Request.BinaryDataN2Information = n2SmBuf
			}

			n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
				PduSessionId: smContext.PDUSessionID,
				// Temporarily assign SMF itself,
				// TODO: TS 23.502 4.2.3.3 5. Namf_Communication_N1N2TransferFailureNotification
				N1n2FailureTxfNotifURI: fmt.Sprintf("%s://%s:%d",
					smf_context.GetSelf().URIScheme,
					smf_context.GetSelf().RegisterIPv4,
					smf_context.GetSelf().SBIPort),
				N2InfoContainer: &models.N2InfoContainer{
					N2InformationClass: models.N2InformationClass_SM,
					SmInfo: &models.N2SmInformation{
						PduSessionId: smContext.PDUSessionID,
						N2InfoContent: &models.N2InfoContent{
							NgapIeType: models.AmfCommunicationNgapIeType_PDU_RES_SETUP_REQ,
							NgapData: &models.RefToBinaryData{
								ContentId: "N2SmInformation",
							},
						},
						SNssai: smContext.SNssai,
					},
				},
			}

			ctx, _, errToken := smf_context.GetSelf().GetTokenCtx(models.ServiceName_NAMF_COMM, models.NrfNfManagementNfType_AMF)
			if errToken != nil {
				logger.PfcpLog.Warnf("Get NAMF_COMM context failed: %s", errToken)
				return
			}
			rspData, err := service.GetApp().Consumer().
				N1N2MessageTransfer(ctx, smContext.Supi, n1n2Request, smContext.CommunicationClientApiPrefix)
			if err != nil {
				logger.ConsumerLog.Warnf("Send N1N2Transfer failed: %s", err)
				return
			}

			if rspData.Cause == models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE {
				logger.PfcpLog.Infof("Receive %v, AMF is able to page the UE", rspData.Cause)
			}
			if rspData.Cause == models.N1N2MessageTransferCause_UE_NOT_RESPONDING {
				logger.PfcpLog.Warnf("%v", rspData.Cause)
				// TODO: TS 23.502 4.2.3.3 3c. Failure indication
			}
		}
	}

	if req.ReportType.Usar && req.UsageReport != nil {
		smContext.HandleReports(req.UsageReport, nil, nil, upfNodeID, "")
		// After receiving the Usage Report, it should send charging request to the CHF
		// and update the URR with the quota or other charging information according to
		// the charging response
		service.GetApp().Processor().ReportUsageAndUpdateQuota(smContext)
	}

	// TS 23.502 4.2.3.3 2b. Send Data Notification Ack, SMF->UPF
	cause.CauseValue = pfcpType.CauseRequestAccepted
	pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, remoteSEID)
}
