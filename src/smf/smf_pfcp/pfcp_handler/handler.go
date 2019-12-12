package pfcp_handler

import (
	"context"
	"fmt"
	"free5gc/lib/http_wrapper"
	"free5gc/lib/openapi/models"
	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_context"
	"free5gc/src/smf/smf_handler/smf_message"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
	"net/http"
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

	if rsp.UPFSEID != nil {
		UPFSEID := rsp.UPFSEID
		smContext.RemoteSEID = UPFSEID.Seid
	}

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

func HandlePfcpSessionModificationResponse(msg *pfcpUdp.Message, HttpResponseQueue *smf_message.ResponseQueue) {
	pfcpRsp := msg.PfcpMessage.Body.(pfcp.PFCPSessionModificationResponse)

	SEID := msg.PfcpMessage.Header.SEID
	seqNum := msg.PfcpMessage.Header.SequenceNumber

	if HttpResponseQueue.CheckItemExist(seqNum) {
		if pfcpRsp.Cause.CauseValue == pfcpType.CauseRequestAccepted {
			resQueueItem := HttpResponseQueue.GetItem(seqNum)

			resQueueItem.RspChan <- smf_message.HandlerResponseMessage{HTTPResponse: &http_wrapper.Response{
				Header: nil,
				Status: http.StatusOK,
				Body:   resQueueItem.ResponseBody,
			}}

			HttpResponseQueue.DeleteItem(seqNum)
			//if smContext.SMState == smf_context.PDUSessionInactive {
			//	smNasBuf, _ := smf_context.BuildGSMPDUSessionEstablishmentAccept(smContext)
			//		n1n2Request := models.N1N2MessageTransferRequest{}
			//			n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
			//					N1MessageContainer: &models.N1MessageContainer{
			//							N1MessageClass:   "SM",
			//								N1MessageContent: &models.RefToBinaryData{ContentId: "GSM_NAS"},
			//								},
			//								}
			//									n1n2Request.BinaryDataN1Message = smNasBuf

			// 	logger.PfcpLog.Warnf("N1N2 Transfer")

			//rspData, _, err := smContext.CommunicationClient.N1N2MessageCollectionDocumentApi.N1N2MessageTransfer(context.Background(), smContext.Supi, n1n2Request)
			//if err != nil {
			//		logger.PfcpLog.Warnf("Send N1N2Transfer failed")
			//		}
			//			if rspData.Cause == models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED {
			//					logger.PfcpLog.Warnf("%v", rspData.Cause)
			//					}
			// 		smContext.SMState = smf_context.PDUSessionActive
			// }
			logger.PfcpLog.Infof("PFCP Session Modification Success[%d]\n", SEID)
		} else {
			logger.PfcpLog.Infof("PFCP Session Modification Failed[%d]\n", SEID)
		}
	} else {
		logger.PfcpLog.Infof("[PFCP Modification RSP] Can't find corresponding seq num[%d]\n", seqNum)
	}

}

func HandlePfcpSessionDeletionResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Deletion Response handling is not implemented")
}

func HandlePfcpSessionReportRequest(msg *pfcpUdp.Message) {
	req := msg.PfcpMessage.Body.(pfcp.PFCPSessionReportRequest)

	SEID := msg.PfcpMessage.Header.SEID
	smContext := smf_context.GetSMContextBySEID(SEID)
	seqFromUPF := msg.PfcpMessage.Header.SequenceNumber

	var cause pfcpType.Cause

	if smContext == nil {
		logger.PfcpLog.Warnf("PFCP Session Report Request Found SM Context NULL, Request Rejected")
		cause.CauseValue = pfcpType.CauseRequestRejected
		// TODO fix: SEID should be the value sent by UPF but now the SEID value is from sm context
		pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, SEID)
		return
	}

	if req.ReportType.Dldr {
		downlinkDataReport := req.DownlinkDataReport
		pdrID := downlinkDataReport.PDRID.RuleId

		if downlinkDataReport.DownlinkDataServiceInformation != nil {
			logger.PfcpLog.Warnf("PFCP Session Report Request DownlinkDataServiceInformation handling is not implemented")
		}

		if smContext.Tunnel.DLPDR.PDRID == pdrID {
			// TS 23.502 4.2.3.3 2b. Send Data Notification Ack, SMF->UPF
			cause.CauseValue = pfcpType.CauseRequestAccepted
			// TODO fix: SEID should be the value sent by UPF but now the SEID value is from sm context
			pfcp_message.SendPfcpSessionReportResponse(msg.RemoteAddr, cause, seqFromUPF, SEID)

			// TS 23.502 4.2.3.3 3a. Send Namf_Communication_N1N2MessageTransfer Request, SMF->AMF
			n2SmBuf, _ := smf_context.BuildPDUSessionResourceSetupRequestTransfer(smContext)
			n1n2Request := models.N1N2MessageTransferRequest{}
			n1n2Request.JsonData = &models.N1N2MessageTransferReqData{
				PduSessionId: smContext.PDUSessionID,
				// Temporarily assign SMF itself, TODO: TS 23.502 4.2.3.3 5. Namf_Communication_N1N2TransferFailureNotification
				N1n2FailureTxfNotifURI: fmt.Sprintf("%s://%s:%d", smf_context.SMF_Self().URIScheme, smf_context.SMF_Self().HTTPAddress, smf_context.SMF_Self().HTTPPort),
				N2InfoContainer: &models.N2InfoContainer{
					N2InformationClass: models.N2InformationClass_SM,
					SmInfo: &models.N2SmInformation{
						PduSessionId: smContext.PDUSessionID,
						N2InfoContent: &models.N2InfoContent{
							NgapIeType: models.NgapIeType_PDU_RES_SETUP_REQ,
							NgapData: &models.RefToBinaryData{
								ContentId: "N2SmInformation",
							},
						},
						SNssai: smContext.Snssai,
					},
				},
			}
			n1n2Request.BinaryDataN2Information = n2SmBuf

			rspData, _, err := smContext.CommunicationClient.N1N2MessageCollectionDocumentApi.N1N2MessageTransfer(context.Background(), smContext.Supi, n1n2Request)
			if err != nil {
				logger.PfcpLog.Warnf("Send N1N2Transfer failed")
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
}

func HandlePfcpSessionReportResponse(msg *pfcpUdp.Message) {
	logger.PfcpLog.Warnf("PFCP Session Report Response handling is not implemented")
}
