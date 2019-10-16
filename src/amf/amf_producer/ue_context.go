package amf_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"free5gc/src/amf/gmm"
	"free5gc/src/amf/logger"
	"net/http"
	"strings"
)

func HandleCreateUeContextRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.CreateUeContextRequest) {
	var response models.CreateUeContextResponse
	var rspErr models.UeContextCreateError
	var problem models.ProblemDetails
	amfSelf := amf_context.AMF_Self()

	ueContextCreateData := body.JsonData

	if ueContextCreateData.UeContext == nil || ueContextCreateData.TargetId == nil || ueContextCreateData.PduSessionList == nil || ueContextCreateData.SourceToTargetData == nil || ueContextCreateData.N2NotifyUri == "" {
		{
			rspErr.Error = &problem
			problem.Status = 403
			problem.Cause = "HANDOVER_FAILURE"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, rspErr)
			return
		}
	}
	// create the UE context in target amf
	ue := amfSelf.NewAmfUe(ueContextId)
	if err := gmm.InitAmfUeSm(ue); err != nil {
		HttpLog.Errorf("InitAmfUeSm error: %v", err.Error())
	}
	//amfSelf.AmfRanSetByRanId(*ueContextCreateData.TargetId.RanNodeId)
	// ue.N1N2Message[ueContextId] = &amf_context.N1N2Message{}
	// ue.N1N2Message[ueContextId].Request.JsonData = &models.N1N2MessageTransferReqData{
	// 	N2InfoContainer: &models.N2InfoContainer{
	// 		SmInfo: &models.N2SmInformation{
	// 			N2InfoContent: ueContextCreateData.SourceToTargetData,
	// 		},
	// 	},
	// }
	ue.HandoverNotifyUri = ueContextCreateData.N2NotifyUri

	amfSelf.AmfRanFindByRanId(*ueContextCreateData.TargetId.RanNodeId)
	supportedTAI := amf_context.NewSupportedTAI()
	supportedTAI.Tai.Tac = ueContextCreateData.TargetId.Tai.Tac
	supportedTAI.Tai.PlmnId = ueContextCreateData.TargetId.Tai.PlmnId
	ue.N1N2MessageSubscribeInfo[ueContextId] = &models.UeN1N2InfoSubscriptionCreateData{
		N2NotifyCallbackUri: ueContextCreateData.N2NotifyUri,
	}
	ue.UnauthenticatedSupi = ueContextCreateData.UeContext.SupiUnauthInd
	//should be smInfo list

	for _, smInfo := range ueContextCreateData.PduSessionList {
		if smInfo.N2InfoContent.NgapIeType == "NgapIeType_HANDOVER_REQUIRED" {
			// ue.N1N2Message[amfSelf.Uri].Request.JsonData.N2InfoContainer.SmInfo = &smInfo
		}
	}

	ue.RoutingIndicator = ueContextCreateData.UeContext.RoutingIndicator

	// optional
	ue.UdmGroupId = ueContextCreateData.UeContext.UdmGroupId
	ue.AusfGroupId = ueContextCreateData.UeContext.AusfGroupId
	//ueContextCreateData.UeContext.HpcfId
	ue.RatType = ueContextCreateData.UeContext.RestrictedRatList[0] //minItem = -1
	//ueContextCreateData.UeContext.ForbiddenAreaList
	//ueContextCreateData.UeContext.ServiceAreaRestriction
	//ueContextCreateData.UeContext.RestrictedCoreNwTypeList

	//it's not in 5.2.2.1.1 step 2a, so don't support
	//ue.Gpsi = ueContextCreateData.UeContext.GpsiList
	//ue.Pei = ueContextCreateData.UeContext.Pei
	//ueContextCreateData.UeContext.GroupList
	//ueContextCreateData.UeContext.DrxParameter
	//ueContextCreateData.UeContext.SubRfsp
	//ueContextCreateData.UeContext.UsedRfsp
	//ue.UEAMBR = ueContextCreateData.UeContext.SubUeAmbr
	//ueContextCreateData.UeContext.SmsSupport
	//ueContextCreateData.UeContext.SmsfId
	//ueContextCreateData.UeContext.SeafData
	//ueContextCreateData.UeContext.Var5gMmCapability
	//ueContextCreateData.UeContext.PcfId
	//ueContextCreateData.UeContext.PcfAmPolicyUri
	//ueContextCreateData.UeContext.AmPolicyReqTriggerList
	//ueContextCreateData.UeContext.EventSubscriptionList
	//ueContextCreateData.UeContext.MmContextList
	//ue.CurPduSession.PduSessionId = ueContextCreateData.UeContext.SessionContextList.
	//ue.TraceData = ueContextCreateData.UeContext.TraceData
	response.JsonData = &models.UeContextCreatedData{
		UeContext: &models.UeContext{
			Supi: ueContextCreateData.UeContext.Supi,
		},
	}

	// response.JsonData.TargetToSourceData = ue.N1N2Message[ueContextId].Request.JsonData.N2InfoContainer.SmInfo.N2InfoContent
	response.JsonData.PduSessionList = ueContextCreateData.PduSessionList
	response.JsonData.PcfReselectedInd = false // TODO:When  Target AMF selects a nw PCF for AM policy, set the flag to true.

	//response.UeContext = ueContextCreateData.UeContext
	//response.TargetToSourceData = ue.N1N2Message[amfSelf.Uri].Request.JsonData.N2InfoContainer.SmInfo.N2InfoContent
	//response.PduSessionList = ueContextCreateData.PduSessionList
	//response.PcfReselectedInd = false // TODO:When  Target AMF selects a nw PCF for AM policy, set the flag to true.
	//

	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusCreated, response)
}

func HandleUEContextReleaseRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.UeContextRelease) {
	var problem models.ProblemDetails
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()
	ueContextRelease := body

	// emergency handle
	if ueContextRelease.Supi != "" {
		if ueContextRelease.UnauthenticatedSupi {

		}
	}

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
			problem.Status = 403
			problem.Cause = "SUPI_OR_PEI_UNKNOWN"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
			return
		}
	}
	if ue != nil {
		ue.Remove()
	}
	// TODO : Ngap handle
	//if ueContextRelease.NgapCause.Group == ngapType.CauseRadioNetwork{
	//	if ueContextRelease.NgapCause.Value == 	ngapType.CauseRadioNetworkPresentHandoverCancelled {
	//
	//	}
	//}
	//ueContextRelease.NgapCause.Value
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}

func HandleUEContextTransferRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.UeContextTransferRequest) {
	var response models.UeContextTransferResponse

	var problem models.ProblemDetails
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()

	if body.JsonData == nil {
		problem.Status = 403
		problem.Cause = "CONTEXT_NOT_FOUND"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
		return
	}
	UeContextTransferReqData := body.JsonData

	if UeContextTransferReqData.AccessType == "" || UeContextTransferReqData.Reason == "" {
		problem.Status = 403
		problem.Cause = "CONTEXT_NOT_FOUND"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
		return
	}

	if strings.HasPrefix(ueContextId, "imsi") {
		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			problem.Status = 403
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
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
			problem.Status = 403
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
			return
		}
	}
	response.JsonData = new(models.UeContextTransferRspData)
	ueContextTransferRspData := response.JsonData

	if ue != nil {
		if ue.GetAnType() != UeContextTransferReqData.AccessType {
			for _, tai := range ue.RegistrationArea[ue.GetAnType()] {
				if UeContextTransferReqData.PlmnId == tai.PlmnId {
					// TODO : generate N2 signalling
				}
			}
		}
		if UeContextTransferReqData.Reason == models.TransferReason_INIT_REG {
			// TODO optional
			//m := nas.NewMessage()
			//m.GmmMessage = nas.NewGmmMessage()
			//m.GmmHeader.SetMessageType(nas.MsgTypeRegistrationRequest)
			//m.GmmMessageDecode(&body.BinaryDataN1Message)
			//
			//registrationType5GS := m.RegistrationRequest.NgksiAndRegistrationType5GS.GetRegistrationType5GS()
			//switch registrationType5GS {
			//default:
			//	logger.ProducerLog.Debugln(registrationType5GS)
			//}
			//mobileIdentity5GSContents := m.RegistrationRequest.MobileIdentity5GS.GetMobileIdentity5GSContents()
			//switch mobileIdentity5GSContents[0] & 0x07 {
			//// cover guti and compare
			//}
			ueContextTransferRspData.UeContext = &models.UeContext{
				Supi:                     ue.Supi,
				SupiUnauthInd:            ue.UnauthenticatedSupi,
				GpsiList:                 nil,
				Pei:                      "",
				UdmGroupId:               ue.UdmGroupId,
				AusfGroupId:              ue.AusfGroupId,
				RoutingIndicator:         ue.RoutingIndicator,
				GroupList:                nil,
				DrxParameter:             "",
				SubRfsp:                  0,
				UsedRfsp:                 0,
				SubUeAmbr:                nil,
				SmsSupport:               "",
				SmsfId:                   "",
				SeafData:                 nil,
				Var5gMmCapability:        "",
				PcfId:                    "",
				PcfAmPolicyUri:           "",
				AmPolicyReqTriggerList:   nil,
				HpcfId:                   "",
				RestrictedRatList:        []models.RatType{ue.RatType},
				ForbiddenAreaList:        nil,
				ServiceAreaRestriction:   nil,
				RestrictedCoreNwTypeList: nil,
				EventSubscriptionList:    nil,
				MmContextList:            nil,
				SessionContextList:       nil,
				TraceData:                nil,
			}
		} else if UeContextTransferReqData.Reason == models.TransferReason_MOBI_REG {
			ueContextTransferRspData.UeContext = &models.UeContext{
				Supi:                     ue.Supi,
				SupiUnauthInd:            ue.UnauthenticatedSupi,
				GpsiList:                 nil,
				Pei:                      "",
				UdmGroupId:               ue.UdmGroupId,
				AusfGroupId:              ue.AusfGroupId,
				RoutingIndicator:         ue.RoutingIndicator,
				GroupList:                nil,
				DrxParameter:             "",
				SubRfsp:                  0,
				UsedRfsp:                 0,
				SubUeAmbr:                nil,
				SmsSupport:               "",
				SmsfId:                   "",
				SeafData:                 nil,
				Var5gMmCapability:        "",
				PcfId:                    "",
				PcfAmPolicyUri:           "",
				AmPolicyReqTriggerList:   nil,
				HpcfId:                   "",
				RestrictedRatList:        []models.RatType{ue.RatType},
				ForbiddenAreaList:        nil,
				ServiceAreaRestriction:   nil,
				RestrictedCoreNwTypeList: nil,
				EventSubscriptionList:    nil,
				MmContextList:            nil,
				SessionContextList:       nil,
				TraceData:                nil,
			}
			ueContextTransferRspData.UeRadioCapability = &models.N2InfoContent{
				NgapMessageType: 0,
				NgapIeType:      models.NgapIeType_UE_RADIO_CAPABILITY,
				NgapData: &models.RefToBinaryData{
					ContentId: "1",
				},
			}
			b := []byte(ue.UeRadioCapability)
			copy(response.BinaryDataN2Information, b)
		} else {
			logger.ProducerLog.Errorln("error Reason")
			problem.Status = 403
			problem.Cause = "CONTEXT_NOT_FOUND"
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
		}
	}
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
}

func HandleAssignEbiDataRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.AssignEbiData) {
	var response models.AssignedEbiData
	var assignEbiError models.AssignEbiError
	var assignEbiFailed models.AssignEbiFailed
	var problem models.ProblemDetails
	var ue *amf_context.AmfUe
	var ok bool
	amfSelf := amf_context.AMF_Self()

	if strings.HasPrefix(ueContextId, "imsi") {
		if ue, ok = amfSelf.UePool[ueContextId]; !ok {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			assignEbiError.Error = &problem
			assignEbiFailed.PduSessionId = body.PduSessionId
			assignEbiFailed.FailedArpList = nil
			assignEbiError.FailureDetails = &assignEbiFailed
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, assignEbiError)
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
			assignEbiError.Error = &problem
			assignEbiFailed.PduSessionId = body.PduSessionId
			assignEbiFailed.FailedArpList = nil
			assignEbiError.FailureDetails = &assignEbiFailed
			amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, assignEbiError)
			return
		}
	}

	if ue != nil {
		if ue.SmContextList[body.PduSessionId] != nil {
			response.PduSessionId = body.PduSessionId
			response.AssignedEbiList = ue.SmContextList[body.PduSessionId].PduSessionContext.AllocatedEbiList
		} else {
			logger.ProducerLog.Errorln("ue.SmContextList is nil")
		}

	}
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
}

func HandleRegistrationStatusUpdateRequest(httpChannel chan amf_message.HandlerResponseMessage, ueContextId string, body models.UeRegStatusUpdateReqData) {
	var response models.UeRegStatusUpdateRspData
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
	} else {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
		return
	}

	if ue != nil {
		if body.TransferStatus == models.UeContextTransferStatus_TRANSFERRED {
			// remove the individual ueContext resource and release any PDU session(s)
			ue.Remove()
			if body.PcfReselectedInd {
				// TODO: send terminate AM Policy Association to the old PCF
			} else {

			}
		} else {
			// NOT_TRANSFERRED

		}
	}
	response.RegStatusTransferComplete = true
	amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, response)
}
