package pfcp_handler

import (
	"context"
	"free5gc/lib/openapi/models"

	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
)

func HandlePfcpHeartbeatRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Heartbeat Request handling is not implemented")
}

func HandlePfcpHeartbeatResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Heartbeat Response handling is not implemented")
}

func HandlePfcpPfdManagementRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP PFD Management Request handling is not implemented")
}

func HandlePfcpPfdManagementResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP PFD Management Response handling is not implemented")
}

func HandlePfcpAssociationSetupRequest(msg *pfcpUdp.Message) {
	//pfcpMsg := msg.PfcpMessage.Body.(pfcp.PFCPAssociationSetupRequest)

	// TODO: check if request is valid

	//upfId, err := generateUpfIdFromNodeId(*pfcpMsg.NodeID)
	//if err != nil {
	//	logger.PfcpLog.Errorf(err.Error())
	//	return
	//}

	//upfNode := smf_context.RetrieveUPFNodeByUpfId(upfId)
	//upfNode.NodeID = pfcpMsg.NodeID
	//upfNode.RecoveryTimeStamp = pfcpMsg.RecoveryTimeStamp
	//upfNode.UPFunctionFeatures = pfcpMsg.UPFunctionFeatures
	//upfNode.UserPlaneIPResourceInformation = pfcpMsg.UserPlaneIPResourceInformation
	//upfNode.NodeID = *pfcpMsg.NodeID
	//upfNode.UPIPInfo = *pfcpMsg.UserPlaneIPResourceInformation

	// Response with PFCP Association Setup Response
	cause := pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}
	pfcp_message.SendPfcpAssociationSetupResponse(msg.RemoteAddr, cause)
}

func HandlePfcpAssociationSetupResponse(msg *pfcpUdp.Message) {
	req := msg.PfcpMessage.Body.(pfcp.PFCPAssociationSetupResponse)

	if req.Cause.CauseValue == pfcpType.CauseRequestAccepted {
		if req.NodeID == nil {
			logger.PfcpLog.Errorln("Association Setup Response Node ID not found")
		}

		if req.UserPlaneIPResourceInformation != nil {
			upf := smf_context.AddUPF(req.NodeID)
			upf.UPIPInfo = *req.UserPlaneIPResourceInformation
			logger.PfcpLog.Infof("UPF[%s]", upf.UPIPInfo.NetworkInstance)
		}
	}
}

func HandlePfcpAssociationUpdateRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Association Update Request handling is not implemented")
}

func HandlePfcpAssociationUpdateResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Association Update Response handling is not implemented")
}

// Deprecated: PFCP Association Release Request should be initiated by the CP function
func HandlePfcpAssociationReleaseRequest(msg *pfcpUdp.Message) {
	pfcpMsg := msg.PfcpMessage.Body.(pfcp.PFCPAssociationReleaseRequest)

	var cause pfcpType.Cause
	upfNode := smf_context.RetrieveUPFNodeByNodeId(*pfcpMsg.NodeID)
	if upfNode != nil {
		smf_context.RemoveUPFNodeByNodeId(*pfcpMsg.NodeID)
		cause.CauseValue = pfcpType.CauseRequestAccepted
	} else {
		cause.CauseValue = pfcpType.CauseNoEstablishedPfcpAssociation
	}
	pfcp_message.SendPfcpAssociationReleaseResponse(msg.RemoteAddr, cause)
}

func HandlePfcpAssociationReleaseResponse(msg *pfcpUdp.Message) {
	pfcpMsg := msg.PfcpMessage.Body.(pfcp.PFCPAssociationReleaseResponse)

	if pfcpMsg.Cause.CauseValue == pfcpType.CauseRequestAccepted {
		smf_context.RemoveUPFNodeByNodeId(*pfcpMsg.NodeID)
	}
}

func HandlePfcpVersionNotSupportedResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Version Not Support Response handling is not implemented")
}

func HandlePfcpNodeReportRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Node Report Request handling is not implemented")
}

func HandlePfcpNodeReportResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Node Report Response handling is not implemented")
}

func HandlePfcpSessionSetDeletionRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Set Deletion Request handling is not implemented")
}

func HandlePfcpSessionSetDeletionResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Set Deletion Response handling is not implemented")
}

func HandlePfcpSessionEstablishmentResponse(msg *pfcpUdp.Message) {
	rsp := msg.PfcpMessage.Body.(pfcp.PFCPSessionEstablishmentResponse)

	SEID := msg.PfcpMessage.Header.SEID
	smContext := smf_context.GetSMContextBySEID(SEID)
	if rsp.Cause.CauseValue == pfcpType.CauseRequestAccepted {
		smNasBuf, _ := smf_context.BuildGSMPDUSessionEstablishmentAccept(smContext)
		n1n2Request := models.N1N2MessageTransferRequest{}
		n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
			N1MessageContainer: &models.N1MessageContainer{
				N1MessageClass:   "SM",
				N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
			},
		}
		n1n2Request.BinaryDataN1Message = smNasBuf

		rspData, _, err := smContext.CommunicationClient.N1N2MessageCollectionDocumentApi.N1N2MessageTransfer(context.Background(), smContext.Supi, n1n2Request)
		if err != nil {
			logger.PfcpLog.Warnf("Send N1N2Transfer failed")
		}
		if rspData.Cause == models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED {
			logger.PfcpLog.Warnf("%v", rspData.Cause)
		}
	}
}

func HandlePfcpSessionModificationResponse(msg *pfcpUdp.Message) {
	rsp := msg.PfcpMessage.Body.(pfcp.PFCPSessionModificationResponse)

	SEID := msg.PfcpMessage.Header.SEID
	if rsp.Cause.CauseValue == pfcpType.CauseRequestAccepted {
		// if smContext.SMState == smf_context.PDUSessionInactive {
		// 	smNasBuf, _ := smf_context.BuildGSMPDUSessionEstablishmentAccept(smContext)
		// 	n1n2Request := models.N1N2MessageTransferRequest{}
		// 	n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
		// 		N1MessageContainer: &models.N1MessageContainer{
		// 			N1MessageClass:   "SM",
		// 			N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
		// 		},
		// 	}
		// 	n1n2Request.BinaryDataN1Message = smNasBuf

		// 	logger.PfcpLog.Warnf("N1N2 Transfer")

		// 	rspData, _, err := smContext.CommunicationClient.N1N2MessageCollectionDocumentApi.N1N2MessageTransfer(context.Background(), smContext.Supi, n1n2Request)
		// 	if err != nil {
		// 		logger.PfcpLog.Warnf("Send N1N2Transfer failed")
		// 	}
		// 	if rspData.Cause == models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED {
		// 		logger.PfcpLog.Warnf("%v", rspData.Cause)
		// 	}
		// 	smContext.SMState = smf_context.PDUSessionActive
		// }
		logger.PfcpLog.Infof("PFCP Session Modification Success[%d]\n", SEID)
	} else {
		logger.PfcpLog.Infof("PFCP Session Modification Failed[%d]\n", SEID)
	}

}

func HandlePfcpSessionDeletionResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Deletion Response handling is not implemented")
}

func HandlePfcpSessionReportRequest(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Report Request handling is not implemented")
}

func HandlePfcpSessionReportResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Report Response handling is not implemented")
}
