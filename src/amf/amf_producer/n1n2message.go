package amf_producer

import (
	"free5gc/lib/aper"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/amf/amf_producer/amf_producer_callback"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/gmm"
	"free5gc/src/amf/gmm/gmm_message"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"
	"net/http"
	"strconv"
	"strings"
)

// TS23502 4.2.3.3, 4.2.4.3, 4.3.2.2, 4.3.2.3, 4.3.3.2, 4.3.7
func HandleN1N2MessageTransferRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId, reqUri string, body models.N1N2MessageTransferRequest) {

	var response models.N1N2MessageTransferRspData
	var transferErr models.N1N2MessageTransferError
	var problem models.ProblemDetails
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()

	if strings.HasPrefix(ueContextId, "imsi") {
		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, transferErr)
			return
		}
	} else if strings.HasPrefix(ueContextId, "imei") {
		for _, ue1 := range amfSelf.UePool {
			if ue1.Pei == ueContextId {
				ue = ue1
				break
			}
		}
		if ue == nil {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, transferErr)
			return
		}
	}
	requestData := body.JsonData
	n2Info := body.BinaryDataN2Information
	n1Msg := body.BinaryDataN1Message
	anType := models.AccessType__3_GPP_ACCESS
	var smContext *amf_context.SmContext
	if requestData.N1MessageContainer != nil && requestData.N1MessageContainer.N1MessageClass == models.N1MessageClass_SM {
		smContext = ue.SmContextList[requestData.PduSessionId]
	}
	if smContext == nil && requestData.N2InfoContainer != nil && requestData.N2InfoContainer.N2InformationClass == models.N2InformationClass_SM {
		smContext = ue.SmContextList[requestData.PduSessionId]
	}
	if smContext != nil {
		anType = smContext.PduSessionContext.AccessType
	}
	onGoing := ue.OnGoing[anType]
	// TODO: Error Status 307, 403 in TS29.518 Table 6.1.3.5.3.1-3
	if onGoing != nil {
		switch onGoing.Procedure {
		case amf_context.OnGoingProcedurePaging:

			if requestData.Ppi == 0 || (onGoing.Ppi != 0 && onGoing.Ppi <= requestData.Ppi) {
				transferErr.Error = new(models.ProblemDetails)
				transferErr.Error.Status = 409
				transferErr.Error.Cause = "HIGHER_PRIORITY_REQUEST_ONGOING"
				amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusConflict, transferErr)
				return
			}
			amf_util.ClearT3513(ue)
			amf_producer_callback.SendN1N2TransferFailureNotification(ue, models.N1N2MessageTransferCause_UE_NOT_RESPONDING)
		case amf_context.OnGoingProcedureN2Handover:
			transferErr.Error = new(models.ProblemDetails)
			transferErr.Error.Status = 409
			transferErr.Error.Cause = "TEMPORARY_REJECT_HANDOVER_ONGOING"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusConflict, transferErr)
			return
		}
	}
	if !ue.Sm[anType].Check(gmm_state.REGISTERED) {
		transferErr.Error = new(models.ProblemDetails)
		transferErr.Error.Status = 409
		transferErr.Error.Cause = "TEMPORARY_REJECT_REGISTRATION_ONGOING"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusConflict, transferErr)
		return
	}

	if ue.CmConnect(anType) {
		response.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)

		if n2Info == nil {
			switch requestData.N1MessageContainer.N1MessageClass {
			case models.N1MessageClass_SM:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &requestData.PduSessionId, 0, nil, 0)
			case models.N1MessageClass_LPP:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeLPP, n1Msg, nil, 0, nil, 0)
			case models.N1MessageClass_SMS:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeSMS, n1Msg, nil, 0, nil, 0)
			case models.N1MessageClass_UPDP:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeUEPolicy, n1Msg, nil, 0, nil, 0)
			}
			return
		}
		if smContext != nil {
			smInfo := requestData.N2InfoContainer.SmInfo
			switch smInfo.N2InfoContent.NgapIeType {
			case models.NgapIeType_PDU_RES_SETUP_REQ:
				HttpLog.Debugln("AMF Transfer NGAP PDU Resource Setup Req from SMF")
				var nasPdu []byte
				var err error
				if n1Msg != nil {
					pduSessionId := uint8(smInfo.PduSessionId)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						logger.HttpLog.Errorln(err.Error())
					}
				}
				list := ngapType.PDUSessionResourceSetupListSUReq{}
				ngap_message.AppendPDUSessionResourceSetupListSUReq(&list, smInfo.PduSessionId, *smInfo.SNssai, nasPdu, n2Info)
				ngap_message.SendPDUSessionResourceSetupRequest(ue.RanUe[anType], nil, list)
			case models.NgapIeType_PDU_RES_MOD_REQ:
				HttpLog.Debugln("AMF Transfer NGAP PDU Resource Modify Req from SMF")
				var nasPdu []byte
				var err error
				if n1Msg != nil {
					pduSessionId := uint8(smInfo.PduSessionId)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						logger.HttpLog.Errorln(err.Error())
					}
				}
				list := ngapType.PDUSessionResourceModifyListModReq{}
				ngap_message.AppendPDUSessionResourceModifyListModReq(&list, smInfo.PduSessionId, nasPdu, n2Info)
				ngap_message.SendPDUSessionResourceModifyRequest(ue.RanUe[anType], list)

			case models.NgapIeType_PDU_RES_REL_CMD:
				HttpLog.Debugln("AMF Transfer NGAP PDU Resource Rel CMD from SMF")
				var nasPdu []byte
				var err error
				if n1Msg != nil {
					pduSessionId := uint8(smInfo.PduSessionId)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						logger.HttpLog.Errorln(err.Error())
					}
				}
				list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
				ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, smInfo.PduSessionId, n2Info)
				ngap_message.SendPDUSessionResourceReleaseCommand(ue.RanUe[anType], nasPdu, list)
			}
		} else {
			//ToDo: send n2 info for non pdu session case
		}
		return
	}

	// 409: transfer a N2 PDU Session Resource Release Command to a 5G-AN and if the UE is in CM-IDLE
	if smContext != nil && n2Info != nil && requestData.N2InfoContainer.SmInfo.N2InfoContent.NgapIeType == models.NgapIeType_PDU_RES_REL_CMD {
		transferErr.Error = new(models.ProblemDetails)
		transferErr.Error.Status = 409
		transferErr.Error.Cause = "UE_IN_CM_IDLE_STATE"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusConflict, transferErr)
		return
	}
	// 504: the UE in MICO mode or the UE is only registered over Non-3GPP access and its state is CM-IDLE
	if !ue.Sm[models.AccessType__3_GPP_ACCESS].Check(gmm_state.REGISTERED) {
		transferErr.Error = new(models.ProblemDetails)
		transferErr.Error.Status = 504
		transferErr.Error.Cause = "UE_NOT_REACHABLE"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusGatewayTimeout, transferErr)
		return
	}
	var pagingPriority *ngapType.PagingPriority
	locationHeader := amf_context.AMF_Self().GetIPv4Uri() + reqUri + "/" + strconv.Itoa(ue.N1N2MessageIDGenerator)
	// Case A (UE is CM-IDLE in 3GPP access and the associated access type is 3GPP access) in subclause 5.2.2.3.1.2 of TS29518
	if anType == models.AccessType__3_GPP_ACCESS {
		if requestData.SkipInd && n2Info == nil {
			response.Cause = models.N1N2MessageTransferCause_N1_MSG_NOT_TRANSFERRED
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
		} else {
			response.Cause = models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE
			message := amf_context.N1N2Message{
				Request:     body,
				Status:      response.Cause,
				ResourceUri: locationHeader,
			}
			ue.N1N2Message = &message
			ue.N1N2MessageIDGenerator++
			headers := http.Header{
				"Location": {locationHeader},
			}
			onGoing.Procedure = amf_context.OnGoingProcedurePaging
			onGoing.Ppi = requestData.Ppi
			amf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusAccepted, response)
			if onGoing.Ppi != 0 {
				pagingPriority = new(ngapType.PagingPriority)
				pagingPriority.Value = aper.Enumerated(onGoing.Ppi)
			}
			pkg, err := ngap_message.BuildPaging(ue, pagingPriority, false)
			if err != nil {
				logger.NgapLog.Errorf("Build Paging failed : %s", err.Error())
				return
			}
			ngap_message.SendPaging(ue, pkg)
		}
		// TODO: WAITING_FOR_ASYNCHRONOUS_TRANSFER
		return
	}
	// Case B (UE is CM-IDLE in Non-3GPP access but CM-CONNECTED in 3GPP access and the associated access type is Non-3GPP access)in subclause 5.2.2.3.1.2 of TS29518
	if ue.CmConnect(models.AccessType__3_GPP_ACCESS) {
		if n2Info == nil {
			response.Cause = models.N1N2MessageTransferCause_N1_N2_TRANSFER_INITIATED
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
			gmm_message.SendDLNASTransport(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &requestData.PduSessionId, 0, nil, 0)
		} else {
			response.Cause = models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE
			message := amf_context.N1N2Message{
				Request:     body,
				Status:      response.Cause,
				ResourceUri: locationHeader,
			}
			ue.N1N2Message = &message
			ue.N1N2MessageIDGenerator++
			headers := http.Header{
				"Location": {locationHeader},
			}
			amf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusAccepted, response)
			nasMsg, err := gmm_message.BuildNotification(ue, nasMessage.AccessTypeNon3GPP)
			if err != nil {
				logger.GmmLog.Errorf("Build Notification failed : %s", err.Error())
				return
			}
			gmm_message.SendNotification(ue.RanUe[models.AccessType__3_GPP_ACCESS], nasMsg)
		}
		return
	}
	// Case C ( UE is CM-IDLE in both Non-3GPP access and 3GPP access and the associated access ype is Non-3GPP access) in subclause 5.2.2.3.1.2 of TS29518
	response.Cause = models.N1N2MessageTransferCause_ATTEMPTING_TO_REACH_UE
	message := amf_context.N1N2Message{
		Request:     body,
		Status:      response.Cause,
		ResourceUri: locationHeader,
	}
	ue.N1N2Message = &message
	ue.N1N2MessageIDGenerator++

	headers := http.Header{
		"Location": {locationHeader},
	}
	onGoing.Procedure = amf_context.OnGoingProcedurePaging
	onGoing.Ppi = requestData.Ppi
	amf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusAccepted, response)
	if onGoing.Ppi != 0 {
		pagingPriority = new(ngapType.PagingPriority)
		pagingPriority.Value = aper.Enumerated(onGoing.Ppi)
	}
	pkg, err := ngap_message.BuildPaging(ue, pagingPriority, true)
	if err != nil {
		logger.NgapLog.Errorf("Build Paging failed : %s", err.Error())
	}
	ngap_message.SendPaging(ue, pkg)
}

func HandleN1N2MessageTransferStatusRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId, reqUri string) {
	var problem models.ProblemDetails
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()
	if strings.HasPrefix(ueContextId, "imsi") {

		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
			return
		}
	} else if strings.HasPrefix(ueContextId, "imei") {
		for _, ue1 := range amfSelf.UePool {
			if ue1.Pei == ueContextId {
				ue = ue1
				break
			}
		}
		if ue == nil {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
			return
		}
	}
	resourceUri := amf_context.AMF_Self().GetIPv4Uri() + reqUri
	n1n2Message := ue.N1N2Message
	if n1n2Message == nil || n1n2Message.ResourceUri != resourceUri {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
		return
	}
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, n1n2Message.Status)
}

func HandleN1N2MessageSubscirbeRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.UeN1N2InfoSubscriptionCreateData) {
	var response models.UeN1N2InfoSubscriptionCreatedData

	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()

	if strings.HasPrefix(ueContextId, "imsi") {
		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			ue = amfSelf.NewAmfUe(ueContextId)
			if err := gmm.InitAmfUeSm(ue); err != nil {
				HttpLog.Errorf("InitAmfUeSm error: %v", err.Error())
			}
		}
	}
	if ue != nil {
		newSubscriptionID := strconv.Itoa(ue.N1N2MessageSubscribeIDGenerator)
		ue.N1N2MessageSubscribeInfo[newSubscriptionID] = &body
		ue.N1N2SubscriptionID = newSubscriptionID
		response.N1n2NotifySubscriptionId = ue.N1N2SubscriptionID
		ue.N1N2MessageSubscribeIDGenerator++
	}
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusCreated, response)
}

func HandleN1N2MessageUnSubscribeRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, subscriptionId string) {
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()

	if strings.HasPrefix(ueContextId, "imsi") {
		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			_, ok := ue.N1N2MessageSubscribeInfo[subscriptionId]
			if ok {
				delete(ue.N1N2MessageSubscribeInfo, subscriptionId)
			}
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
		}
	}
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusBadRequest, nil)
}
