package gmm_handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/fsm"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasConvert"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_consumer"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_ngap/ngap_message"
	"free5gc/src/amf/amf_producer/amf_producer_callback"
	"free5gc/src/amf/amf_util"
	"free5gc/src/amf/gmm/gmm_event"
	"free5gc/src/amf/gmm/gmm_message"
	"free5gc/src/amf/gmm/gmm_state"
	"free5gc/src/amf/logger"
	"reflect"
	"strconv"
	"time"

	"github.com/antihax/optional"
	"github.com/mitchellh/mapstructure"
	"github.com/mohae/deepcopy"
)

func HandleULNASTransport(ue *amf_context.AmfUe, anType models.AccessType, ulNasTransport *nasMessage.ULNASTransport) error {
	switch ulNasTransport.GetPayloadContainerType() {
	case nasMessage.PayloadContainerTypeN1SMInfo:

		var pduSessionId int32
		var requestType models.RequestType
		var sNssai *models.Snssai
		var dnn string
		if nasPduSessionId := ulNasTransport.PduSessionID2Value; nasPduSessionId != nil {
			pduSessionId = int32(nasPduSessionId.GetPduSessionID2Value())
		} else {
			return fmt.Errorf("NAS N1 SM info need PDU Session ID")
		}
		// TODO: implement SSC mode3 Op
		if oldPduSessionId := ulNasTransport.OldPDUSessionID; oldPduSessionId != nil {
			return fmt.Errorf("SSC mode3 operation has not been implemented yet")
		}
		if NasSnssai := ulNasTransport.SNSSAI; NasSnssai != nil {
			snssai := nasConvert.SnssaiToModels(NasSnssai)
			sNssai = &snssai
		}
		if NasDnn := ulNasTransport.DNN; NasDnn != nil {
			dnn = string(NasDnn.GetDNN())
		}
		if ulNasTransport.RequestType != nil {
			NasRequestType := ulNasTransport.RequestType.GetRequestTypeValue()

			switch NasRequestType {
			case nasMessage.ULNASTransportRequestTypeInitialRequest:
				requestType = models.RequestType_INITIAL_REQUEST
			case nasMessage.ULNASTransportRequestTypeExistingPduSession:
				requestType = models.RequestType_EXISTING_PDU_SESSION
			case nasMessage.ULNASTransportRequestTypeInitialEmergencyRequest:
				requestType = models.RequestType_INITIAL_EMERGENCY_REQUEST
			case nasMessage.ULNASTransportRequestTypeExistingEmergencyPduSession:
				requestType = models.RequestType_EXISTING_EMERGENCY_PDU_SESSION
			}
		}

		m := nas.NewMessage()
		err := m.GsmMessageDecode(&ulNasTransport.PayloadContainer.Buffer)
		if err != nil {
			return err
		}
		switch m.GsmMessage.GetMessageType() {
		case nas.MsgTypePDUSessionEstablishmentRequest:
			return HandlePDUSessionEstablishmentRequest(ue, anType, ulNasTransport.PayloadContainer.Buffer, pduSessionId, requestType, sNssai, dnn)
		case nas.MsgTypePDUSessionModificationRequest, nas.MsgTypePDUSessionModificationComplete, nas.MsgTypePDUSessionModificationCommandReject:
			return HandlePDUSessionModificationForward(ue, anType, ulNasTransport.PayloadContainer.Buffer, pduSessionId)
		case nas.MsgTypePDUSessionReleaseRequest, nas.MsgTypePDUSessionReleaseComplete, nas.MsgTypePDUSessionReleaseReject:
			return HandlePDUSessionReleaseForward(ue, anType, ulNasTransport.PayloadContainer.Buffer, pduSessionId)
		case nas.MsgTypePDUSessionAuthenticationComplete:
			return HandlePDUSessionAuthenticationForward(ue, anType, ulNasTransport.PayloadContainer.GetPayloadContainerContents(), pduSessionId)
		case nas.MsgTypeStatus5GSM:
			return HandleStatus5GSM(ue, anType, ulNasTransport.PayloadContainer.Buffer, pduSessionId)
		}
	case nasMessage.PayloadContainerTypeSMS:
		return fmt.Errorf("PayloadContainerTypeSMS has not been implemented yet in UL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeLPP:
		return fmt.Errorf("PayloadContainerTypeLPP has not been implemented yet in UL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeSOR:
		return fmt.Errorf("PayloadContainerTypeSOR has not been implemented yet in UL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeUEPolicy:
		logger.GmmLog.Infoln("AMF Transfer UEPolicy To PCF")
		amf_producer_callback.SendN1MessageNotify(ue, models.N1MessageClass_UPDP, ulNasTransport.PayloadContainer.GetPayloadContainerContents(), nil)
	case nasMessage.PayloadContainerTypeUEParameterUpdate:
		logger.GmmLog.Infoln("AMF Transfer UEParameterUpdate To UDM")
		upuMac, err := nasConvert.UpuAckToModels(ulNasTransport.PayloadContainer.GetPayloadContainerContents())
		if err != nil {
			return err
		}
		err = amf_consumer.PutUpuAck(ue, upuMac)
		if err != nil {
			return err
		}
		logger.GmmLog.Debugf("UpuMac[%s] in UPU ACK NAS Msg", upuMac)
	case nasMessage.PayloadContainerTypeMultiplePayload:
		return fmt.Errorf("PayloadContainerTypeMultiplePayload has not been implemented yet in UL NAS TRANSPORT")
	}
	return nil
}
func HandlePDUSessionEstablishmentRequest(ue *amf_context.AmfUe, anType models.AccessType, payload []byte, pduSessionID int32, requestType models.RequestType, sNssai *models.Snssai, dnn string) error {
	// TODO Request Type Emergency requset
	var pduSession models.PduSessionContext
	pduSession.PduSessionId = pduSessionID
	pduSession.AccessType = anType
	if requestType == models.RequestType_INITIAL_REQUEST {
		amfSelf := amf_context.AMF_Self()

		if sNssai == nil {
			if ue.SmfSelectionData != nil {
				for snssai, sNssaiInfo := range ue.SmfSelectionData.SubscribedSnssaiInfos {
					var err error
					sNssai, err = amf_util.SnssaiHexToModels(snssai)
					if err != nil {
						return err
					}
					if dnn == "" {
						for _, dnnInfo := range sNssaiInfo.DnnInfos {
							if dnnInfo.DefaultDnnIndicator {
								dnn = dnnInfo.Dnn
								break
							}
						}
					}

				}
			}
			if sNssai == nil {
				snssai := ue.AllowedNssai[anType]
				if len(snssai) > 0 {
					sNssai = &snssai[0]
				} else {
					err := fmt.Errorf("Ue[%s] doesn't have allowedNssai\n", ue.Supi)
					logger.GmmLog.Errorf(err.Error())
					return err
				}
			}

		}
		pduSession.SNssai = sNssai
		if dnn == "" {
			// default DNN decided by AMF
			dnn = amfSelf.SupportDnnLists[0]
		}
		pduSession.Dnn = dnn

		// NRF use local configure
		// TODO: SMF selection in TS 23502 clause 4.3.2.2.3.2
		//		step 1. AMF invokes the Nnssf_NSSelection_Get with the S-NSSAI from the Allowed NSSAI requested by the UE, PLMN ID of the SUPI, TAI of the UE
		// 		step 2. NSSF returns the appropriate NRF and may return a NSI ID corresponding to the Network Slice instance.
		tmp := []models.Snssai{
			{
				Sst: 1,
				Sd:  "010203",
			},
		}
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NSMF_PDUSESSION}),
			Dnn:          optional.NewString(dnn),
			Snssais:      optional.NewInterface(amf_util.MarshToJsonString(tmp)),
		}

		if ue.PlmnId.Mcc != "" {
			param.TargetPlmnList = optional.NewInterface(amf_util.MarshToJsonString(ue.PlmnId))
		}
		result, err := amf_consumer.SendSearchNFInstances(amfSelf.NrfUri, models.NfType_SMF, models.NfType_AMF, param)
		if err != nil || result.NfInstances == nil {
			err = fmt.Errorf("DNN[%s] is not support by network and AMF can not select an SMF by NRF\n", dnn)
			logger.GmmLog.Errorf(err.Error())
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMDNNNotSupportedOrNotSubscribedInTheSlice, nil, 0)
			return err
		}
		logger.GmmLog.Info("result, err := amf_consumer.SendSearchNFInstances(amfSelf.NrfUri, models.NfType_SMF, models.NfType_AMF, param)")
		// select the first SMF, TODO: select base on other info
		var smfID, smfUri string
		for _, nfProfile := range result.NfInstances {
			for _, service := range *nfProfile.NfServices {
				if service.ServiceName == models.ServiceName_NSMF_PDUSESSION && service.NfServiceStatus == models.NfServiceStatus_REGISTERED {
					smfID = nfProfile.NfInstanceId
					if nfProfile.Fqdn != "" {
						smfUri = nfProfile.Fqdn
					} else if service.Fqdn != "" {
						smfUri = service.Fqdn
					} else if service.ApiPrefix != "" {
						smfUri = service.ApiPrefix
					} else if service.IpEndPoints != nil {
						point := (*service.IpEndPoints)[0]
						if point.Port != 0 {
							smfUri = fmt.Sprintf("%s://%s:%d", service.Scheme, point.Ipv4Address, point.Port)
						}
					}
				}
				if smfUri != "" {
					break
				}
			}
			if smfUri != "" {
				break
			}
		}

		if smfUri == "" {
			err := fmt.Errorf("AMF can not select an SMF by NRF")
			logger.GmmLog.Errorf(err.Error())
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
			return err
		}

		// Store PduSessionContext For duplicated PDU Session Id
		if smContext, ok := ue.SmContextList[pduSessionID]; ok {
			ue.StoredSmContext[pduSessionID] = &amf_context.StoredSmContext{
				SmfId:             smfID,
				SmfUri:            smfUri,
				PduSessionContext: &pduSession,
				AnType:            anType,
				Payload:           payload,
			}
			updateData := models.SmContextUpdateData{
				Release: true,
				Cause:   models.Cause_REL_DUE_TO_DUPLICATE_SESSION_ID,
			}
			response, _, _, err := amf_consumer.SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, updateData, nil, nil)
			if err != nil {
				return err
			} else if response == nil {
				err := fmt.Errorf("PDU Session ID[%d] can't be released in DUPLICATE_SESSION_ID case for request Type \"Initial Request\" in PDU Session Establishment Request\n", pduSessionID)
				logger.GmmLog.Errorf(err.Error())
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
			}
			return nil
		}

		smContextCreateData := amf_consumer.BuildCreateSmContextRequest(ue, pduSession, requestType)

		response, smContextRef, errResponse, problemDetail, err := amf_consumer.SendCreateSmContextRequest(ue, smfUri, payload, smContextCreateData)
		if response != nil {
			var smContext amf_context.SmContext
			pduSession.SmContextRef = smContextRef
			logger.GmmLog.Infof("smconetxt ef : %s", smContextRef)

			smContext.PduSessionContext = &pduSession
			smContext.UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
			smContext.SmfUri = smfUri
			smContext.SmfId = smfID
			ue.SmContextList[pduSession.PduSessionId] = &smContext
			logger.GmmLog.Infof("Http create smContext[pduSessionID: %d] Success", pduSession.PduSessionId)
			// TODO: handle response(response N2SmInfo to RAN if exists)
		} else if errResponse != nil {
			logger.GmmLog.Warnf("PDU Session Establishment Request is rejected by SMF[pduSessionId:%d]\n", pduSession.PduSessionId)
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmInfoToUe, &pduSessionID, 0, nil, 0)
		} else if err != nil {
			return err
		} else {
			// TODO: error handling
			err1 := fmt.Errorf("Failed to Create smContext[pduSessionID: %d], Error[%v]\n", pduSessionID, problemDetail)
			logger.GmmLog.Errorf(err1.Error())
			return err1
		}
	} else if requestType == models.RequestType_EXISTING_PDU_SESSION {
		smContext, ok := ue.SmContextList[pduSessionID]
		if !ok {
			err := fmt.Errorf("PDU Session ID[%d] dose not exist for request Type \"Existing Pdu Session\" in PDU Session Establishment Request\n", pduSessionID)
			logger.GmmLog.Errorf(err.Error())
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
			return err
		}
		// TODO: Compare Smfid PLMN with AMF
		if smContext.SmfId == "" {
			err := fmt.Errorf("PDU Session ID[%d] correspond SMF-ID does not exist\n", pduSessionID)
			logger.GmmLog.Errorf(err.Error())
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
			return err
		}
		targetAccessType := models.AccessType__3_GPP_ACCESS
		if anType == models.AccessType__3_GPP_ACCESS {
			targetAccessType = models.AccessType_NON_3_GPP_ACCESS
		}
		if !ue.InAllowedNssai(*smContext.PduSessionContext.SNssai, targetAccessType) {
			err := fmt.Errorf("S-NSSAI[%v] is not present in the Allowed NSSAI of target Access\n", *smContext.PduSessionContext.SNssai)
			logger.GmmLog.Errorf(err.Error())
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
			return err
		}

		response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextHandoverBetweenAccessType(ue, pduSessionID, targetAccessType, payload)
		if response != nil {
			smContext.PduSessionContext.AccessType = targetAccessType
			logger.GmmLog.Infof("Success update smContext Handover Between AccessType[pduSessionID: %d]", pduSessionID)
			// TODO: handle response(response N2SmInfo to RAN if exists)
		} else if errResponse != nil {
			logger.GmmLog.Warnf("PDU Session Establishment Request is rejected by SMF[pduSessionId:%d]\n", pduSession.PduSessionId)
			if errResponse.BinaryDataN1SmMessage != nil {
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
			}
		} else if err != nil {
			return err
		} else {
			// TODO: error handling
			logger.GmmLog.Errorf("Failed to Create smContext[pduSessionID: %d], Error[%v]\n", pduSessionID, problemDetail)
			return nil
		}
	} else {
		return fmt.Errorf("RequestType INITIAL_EMERGENCY_REQUEST and EXISTING_EMERGENCY_PDU_SESSION haven't been supported")
	}
	return nil
}

func HandlePDUSessionModificationForward(ue *amf_context.AmfUe, anType models.AccessType, payload []byte, pduSessionID int32) (err error) {

	smContext, ok := ue.SmContextList[pduSessionID]
	if !ok {
		err := fmt.Errorf("PDU Session ID[%d] in PDU Session Modification Procedure is not exist\n", pduSessionID)
		logger.GmmLog.Errorf(err.Error())
		gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
		return err
	}
	smContextUpdateData := models.SmContextUpdateData{
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "N1SmMsg",
		},
	}
	smContextUpdateData.Pei = ue.Pei
	smContextUpdateData.Gpsi = ue.Gpsi
	if !amf_context.CompareUserLocation(ue.Location, smContext.UserLocation) {
		smContextUpdateData.UeLocation = &ue.Location
	}
	response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, smContextUpdateData, payload, nil)
	if response != nil {
		smContext.UserLocation = ue.Location
		responseData := response.JsonData
		n1Msg := response.BinaryDataN1SmMessage
		n2Info := response.BinaryDataN2SmInformation
		if n2Info != nil {
			switch responseData.N2SmInfoType {
			case models.N2SmInfoType_PDU_RES_MOD_REQ:
				logger.GmmLog.Infoln("AMF Transfer NGAP PDU Resource Modify Req from SMF")
				var nasPdu []byte
				if n1Msg != nil {
					pduSessionId := uint8(pduSessionID)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						return
					}
				}
				list := ngapType.PDUSessionResourceModifyListModReq{}
				ngap_message.AppendPDUSessionResourceModifyListModReq(&list, pduSessionID, nasPdu, n2Info)
				ngap_message.SendPDUSessionResourceModifyRequest(ue.RanUe[anType], list)
			}
		} else if n1Msg != nil {
			logger.GmmLog.Infoln("AMF tansfer Only N1 PDU Session Modification Msg to UE")
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
		}
	} else if errResponse != nil {
		errJSON := errResponse.JsonData
		n1Msg := errResponse.BinaryDataN1SmMessage
		logger.GmmLog.Warnf("PDU Session Modification Procedure is rejected by SMF[pduSessionId:%d], Error[%s]\n", pduSessionID, errJSON.Error.Cause)
		if n1Msg != nil {
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
		}
		// TODO: handle n2 info transfer
	} else if err != nil {
		return err
	} else {
		// TODO: error handling
		err1 := fmt.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
		logger.GmmLog.Errorf(err1.Error())
		return err1
	}

	return nil
}

func HandlePDUSessionReleaseForward(ue *amf_context.AmfUe, anType models.AccessType, payload []byte, pduSessionID int32) (err error) {

	smContext, ok := ue.SmContextList[pduSessionID]
	if !ok {
		err := fmt.Errorf("PDU Session ID[%d] in PDU Session Release Request is not exist\n", pduSessionID)
		logger.GmmLog.Errorf(err.Error())
		gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
		return err
	}
	smContextUpdateData := models.SmContextUpdateData{
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "N1SmMsg",
		},
	}
	smContextUpdateData.UeLocation = &ue.Location
	response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, smContextUpdateData, payload, nil)
	if response != nil {
		smContext.UserLocation = ue.Location
		responseData := response.JsonData
		n1Msg := response.BinaryDataN1SmMessage
		n2Info := response.BinaryDataN2SmInformation
		if n2Info != nil {
			switch responseData.N2SmInfoType {
			case models.N2SmInfoType_PDU_RES_REL_CMD:
				logger.GmmLog.Infoln("AMF Transfer NGAP PDU Session Resource Rel Co from SMF")
				var nasPdu []byte
				if n1Msg != nil {
					pduSessionId := uint8(pduSessionID)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						return
					}
				}
				list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
				ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2Info)
				ngap_message.SendPDUSessionResourceReleaseCommand(ue.RanUe[anType], nasPdu, list)
			}
		} else if n1Msg != nil {
			logger.GmmLog.Infoln("AMF tansfer Only N1 PDU Session Release Msg to UE")
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionID, 0, nil, 0)
		}
	} else if errResponse != nil {
		errJSON := errResponse.JsonData
		n1Msg := errResponse.BinaryDataN1SmMessage
		logger.GmmLog.Warnf("PDU Session Release Procedure is rejected by SMF[pduSessionId:%d], Error[%s]\n", pduSessionID, errJSON.Error.Cause)
		if n1Msg != nil {
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
		}
	} else if err != nil {
		return err
	} else {
		// TODO: error handling
		err1 := fmt.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
		logger.GmmLog.Errorf(err1.Error())
		return err1
	}
	return nil
}
func HandlePDUSessionAuthenticationForward(ue *amf_context.AmfUe, anType models.AccessType, payload []byte, pduSessionID int32) error {

	smContext, ok := ue.SmContextList[pduSessionID]
	if !ok {
		err := fmt.Errorf("PDU Session ID[%d] in PDU Session Authentication Complete is not exist\n", pduSessionID)
		logger.GmmLog.Errorf(err.Error())
		gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
		return err
	}
	smContextUpdateData := models.SmContextUpdateData{
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "N1SmMsg",
		},
	}
	response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, smContextUpdateData, payload, nil)
	if response != nil {
		logger.GmmLog.Infoln("AMF tansfer PDU Session Authentication Complete to SMF Success")
	} else if errResponse != nil {
		errJSON := errResponse.JsonData
		n1Msg := errResponse.BinaryDataN1SmMessage
		logger.GmmLog.Warnf("PDU Session Authentication Complete is rejected by SMF[pduSessionId:%d], Error[%s]\n", pduSessionID, errJSON.Error.Cause)
		if n1Msg != nil {
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, errResponse.BinaryDataN1SmMessage, &pduSessionID, 0, nil, 0)
		}
	} else if err != nil {
		return err
	} else {
		// TODO: error handling
		err1 := fmt.Errorf("Failed to Update smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
		logger.GmmLog.Errorf(err1.Error())
		return err1
	}
	return nil
}

func HandleRegistrationRequest(ue *amf_context.AmfUe, anType models.AccessType, registrationRequest *nasMessage.RegistrationRequest) error {

	logger.GmmLog.Info("[AMF] Handle Registration Request")

	amf_util.ClearT3513(ue)
	amf_util.ClearT3565(ue)

	amfSelf := amf_context.AMF_Self()

	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	ranUe := ue.RanUe[ue.GetAnType()]
	if ranUe == nil {
		return fmt.Errorf("RanUe is nil")
	}

	ue.RegistrationRequest = registrationRequest
	ue.RegistrationType5GS = registrationRequest.NgksiAndRegistrationType5GS.GetRegistrationType5GS()
	switch ue.RegistrationType5GS {
	case nasMessage.RegistrationType5GSInitialRegistration:
		logger.GmmLog.Debugf("RegistrationType: Initial Registration")
	case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
		logger.GmmLog.Debugf("RegistrationType: Mobility Registration Updating")
	case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
		logger.GmmLog.Debugf("RegistrationType: Periodic Registration Updating")
	case nasMessage.RegistrationType5GSEmergencyRegistration:
		return fmt.Errorf("Not Supportted RegistrationType: Emergency Registration")
	}

	mobileIdentity5GSContents := registrationRequest.MobileIdentity5GS.GetMobileIdentity5GSContents()
	ue.IdentityTypeUsedForRegistration = nasConvert.GetTypeOfIdentity(mobileIdentity5GSContents[0])
	switch ue.IdentityTypeUsedForRegistration { // get type of identity
	case nasMessage.MobileIdentity5GSTypeNoIdentity:
		logger.GmmLog.Debugf("No Identity")
	case nasMessage.MobileIdentity5GSTypeSuci:
		ue.Suci = mobileIdentity5GSContents
		var plmnId string
		// TODO: handle the situation if suci is encrypted
		plmnId, ue.Supi = nasConvert.GetPlmnIdAndSupiFromSuci(mobileIdentity5GSContents)
		ue.PlmnId = amf_util.PlmnIdStringToModels(plmnId)
		logger.GmmLog.Debugf("SUCI[SUPI: %s]", ue.Supi)
	case nasMessage.MobileIdentity5GSType5gGuti:
		guami, guti := nasConvert.GutiToString(mobileIdentity5GSContents)
		ue.Guti = guti
		logger.GmmLog.Debugf("GUTI: %s", guti)

		servedGuami := amfSelf.ServedGuamiList[0]
		if reflect.DeepEqual(guami, servedGuami) {
			logger.GmmLog.Debugf("Serving AMF has changed")
			ue.ServingAmfChanged = true
		} else {
			ue.ServingAmfChanged = false
		}
	case nasMessage.MobileIdentity5GSTypeImei:
		imei := nasConvert.PeiToString(mobileIdentity5GSContents)
		ue.Pei = imei
		logger.GmmLog.Debugf("PEI: %s", imei)
	case nasMessage.MobileIdentity5GSTypeImeisv:
		imeisv := nasConvert.PeiToString(mobileIdentity5GSContents)
		ue.Pei = imeisv
		logger.GmmLog.Debugf("PEI: %s", imeisv)
	}

	// NgKsi: TS 24.501 9.11.3.32
	switch registrationRequest.NgksiAndRegistrationType5GS.GetTSC() {
	case nasMessage.TypeOfSecurityContextFlagNative:
		ue.NgKsi.Tsc = models.ScType_NATIVE
	case nasMessage.TypeOfSecurityContextFlagMapped:
		ue.NgKsi.Tsc = models.ScType_MAPPED
	}
	ue.NgKsi.Ksi = int32(registrationRequest.NgksiAndRegistrationType5GS.GetNasKeySetIdentifiler())

	// Copy UserLocation from ranUe
	ue.Location = ranUe.Location
	ue.Tai = ranUe.Tai

	// Check TAI
	if !amf_context.InTaiList(ue.Tai, amfSelf.SupportTaiLists) {
		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMTrackingAreaNotAllowed, "")
		return fmt.Errorf("Registration Reject[Tracking area not allowed]")
	}

	if ue.RegistrationType5GS != nasMessage.RegistrationType5GSPeriodicRegistrationUpdating {
		if registrationRequest.Capability5GMM != nil {
			ue.Capability5GMM = *registrationRequest.Capability5GMM
		} else {
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
			return fmt.Errorf("Capability5GMM is nil")
		}
		if registrationRequest.UESecurityCapability != nil {
			ue.NasUESecurityCapability = *registrationRequest.UESecurityCapability
			ue.SecurityCapabilities.NREncryptionAlgorithms, ue.SecurityCapabilities.NRIntegrityProtectionAlgorithms,
				ue.SecurityCapabilities.EUTRAEncryptionAlgorithms, ue.SecurityCapabilities.EUTRAIntegrityProtectionAlgorithms =
				nasConvert.UESecurityCapabilityToByteArray(registrationRequest.UESecurityCapability.Buffer)
		} else {
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
			return fmt.Errorf("UESecurityCapability is nil")
		}
	}

	// TS 23.502 4.2.2.2.2 step 1
	// If available, the last visited TAI shall be included in order to help the AMF produce Registration Area for the UE
	if registrationRequest.LastVisitedRegisteredTAI != nil {
		plmnID := nasConvert.PlmnIDToString(registrationRequest.LastVisitedRegisteredTAI.Octet[1:4])
		nasTac := registrationRequest.LastVisitedRegisteredTAI.GetTAC()
		tac := hex.EncodeToString(nasTac[:])

		lastVisitedRegisteredTai := models.Tai{
			PlmnId: &models.PlmnId{
				Mcc: plmnID[:3],
				Mnc: plmnID[3:],
			},
			Tac: tac,
		}

		ue.LastVisitedRegisteredTai = lastVisitedRegisteredTai
		logger.GmmLog.Debugf("Ue Last Visited Registered Tai; %v", ue.LastVisitedRegisteredTai)
	}

	if registrationRequest.MICOIndication != nil {
		logger.GmmLog.Warnf("Receive MICO Indication[RAAI: %d], Not Supported", registrationRequest.MICOIndication.GetRAAI())
	}

	// TODO: Negotiate DRX value if need (TS 23.501 5.4.5)
	if registrationRequest.RequestedDRXParameters != nil {
		switch registrationRequest.RequestedDRXParameters.GetDRXValue() {
		case nasMessage.DRXcycleParameterT32:
			logger.GmmLog.Tracef("Requested DRX: T = 32")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT32
		case nasMessage.DRXcycleParameterT64:
			logger.GmmLog.Tracef("Requested DRX: T = 64")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT64
		case nasMessage.DRXcycleParameterT128:
			logger.GmmLog.Tracef("Requested DRX: T = 128")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT128
		case nasMessage.DRXcycleParameterT256:
			logger.GmmLog.Tracef("Requested DRX: T = 256")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT256
		case nasMessage.DRXValueNotSpecified:
			fallthrough
		default:
			ue.UESpecificDRX = nasMessage.DRXValueNotSpecified
			logger.GmmLog.Tracef("Requested DRX: Value not specified")
		}
	}

	// TODO (TS 23.502 4.2.2.2 step 4): if UE's 5g-GUTI is included & serving AMF has changed since last registration procedure,
	// new AMF may invoke Namf_Communication_UEContextTransfer to old AMF, including the complete registration request nas
	// msg, to request UE's SUPI & UE Context
	if ue.ServingAmfChanged {
	}
	return nil
}

func HandleStatus5GSM(ue *amf_context.AmfUe, anType models.AccessType, payload []byte, pduSessionID int32) (err error) {

	logger.GmmLog.Info("[AMF] Handle Status5GSM")
	smContext, ok := ue.SmContextList[pduSessionID]
	if !ok {
		err := fmt.Errorf("PDU Session ID[%d] in Status 5GSM Procedure is not exist\n", pduSessionID)
		logger.GmmLog.Errorf(err.Error())
		gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo, payload, &pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
		return err
	}
	smContextUpdateData := models.SmContextUpdateData{
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "N1SmMsg",
		},
	}

	_, _, _, err = amf_consumer.SendUpdateSmContextRequest(ue, smContext.SmfUri, smContext.PduSessionContext.SmContextRef, smContextUpdateData, payload, nil)
	if err != nil {
		return err
	}

	return nil
}

func HandleInitialRegistration(ue *amf_context.AmfUe, anType models.AccessType) error {

	amfSelf := amf_context.AMF_Self()
	// Common Registration procedure
	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	if len(ue.Suci) == 0 {
		gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeSuci)
		return nil
	}

	if !ue.SecurityContextIsValid() {
		return startAuthenticationProcedure(ue, anType)
	}
	// update Kgnb/Kn3iwf
	ue.UpdateSecurityContext(anType)

	// Registration with AMF re-allocation (TS 23.502 4.2.2.2.3)
	if len(ue.SubscribedNssai) == 0 {
		// TODO(optional): Nudm_SDM_Get(SUPI, slice selection subscription data) to get UE's slice selection subscription info
	}

	var requestedNssai []models.Snssai
	if ue.RegistrationRequest.RequestedNSSAI != nil {
		requestedNssai = nasConvert.RequestedNssaiToModels(ue.RegistrationRequest.RequestedNSSAI)
		logger.GmmLog.Debugf("RequestedNssai: %+v", requestedNssai)
		// TODO(optional): Nnssf_NSSelection_Get(requested nssai, subscribed s-nssai)
		// TODO(optional): Namf_Communication_RegistrationCompleteNotify(failure cause) to old AMF if amf decides to reroute nas msg to another AMF
		// TODO(optional): Nnrf_NfDiscovery_Request(NF Type, Amf Set): if initial AMF does not have target AMF address
		// TODO(if decide to reroute): Namf_Communication_N1MessageNotify to target AMF or Send Reroute Nas Request to RAN
	} else {
		// assign amf support snssai to ue
		ue.AllowedNssai[anType] = append(ue.AllowedNssai[anType], amfSelf.PlmnSupportList[0].SNssaiList...)
	}

	// TODO (step 10 optional): send Namf_Communication_RegistrationCompleteNotify to old AMF if need
	if ue.ServingAmfChanged {
		// If the AMF has changed the new AMF notifies the old AMF that the registration of the UE in the new AMF is completed
	}

	if len(ue.Pei) == 0 {
		gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeImei)
		return nil
	}

	// TODO (step 12 optional): the new AMF initiates ME identity check by invoking the N5g-eir_EquipmentIdentityCheck_Get service operation

	if ue.ServingAmfChanged || ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Check(gmm_state.REGISTERED) || !ue.ContextValid {
		// UDM selection described in TS 23.501 6.3.8
		// TODO: consider udm group id, Routing ID part of SUCI, GPSI or External Group ID (e.g., by the NEF)
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			Supi: optional.NewString(ue.Supi),
		}
		resp, err := amf_consumer.SendSearchNFInstances(amfSelf.NrfUri, models.NfType_UDM, models.NfType_AMF, param)
		if err != nil {
			logger.GmmLog.Error("AMF can not select an UDM by NRF")
			return err
		}

		// select the first UDM, TODO: select base on other info
		var uecmUri, sdmUri string
		for _, nfProfile := range resp.NfInstances {
			ue.UdmId = nfProfile.NfInstanceId
			uecmUri = amf_util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDM_UECM, models.NfServiceStatus_REGISTERED)
			sdmUri = amf_util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDM_SDM, models.NfServiceStatus_REGISTERED)
			if uecmUri != "" && sdmUri != "" {
				break
			}
		}
		ue.NudmUECMUri = uecmUri
		ue.NudmSDMUri = sdmUri
		if ue.NudmUECMUri == "" || ue.NudmSDMUri == "" {
			err := fmt.Errorf("AMF can not select an UDM by NRF")
			logger.GmmLog.Errorf(err.Error())
			return err
		}

		problemDetails, err := amf_consumer.UeCmRegistration(ue, anType, true)
		if problemDetails != nil {
			logger.GmmLog.Errorf("UECM_Registration Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("UECM_Registration Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMGetAmData(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM_Get AmData Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM_Get AmData Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMGetSmfSelectData(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM_Get SmfSelectData Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM_Get SmfSelectData Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMGetUeContextInSmfData(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM_Get UeContextInSmfData Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM_Get UeContextInSmfData Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMSubscribe(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM Subscribe Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM Subscribe Error[%+v]", err)
		}

		ue.ContextValid = true
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		Supi: optional.NewString(ue.Supi),
	}
	for {
		resp, err := amf_consumer.SendSearchNFInstances(amfSelf.NrfUri, models.NfType_PCF, models.NfType_AMF, param)
		if err != nil {
			logger.GmmLog.Error("AMF can not select an PCF by NRF")
		} else {
			// select the first PCF, TODO: select base on other info
			var pcfUri string
			for _, nfProfile := range resp.NfInstances {
				pcfUri = amf_util.SearchNFServiceUri(nfProfile, models.ServiceName_NPCF_AM_POLICY_CONTROL, models.NfServiceStatus_REGISTERED)
				if pcfUri != "" {
					break
				}
			}
			if ue.PcfUri = pcfUri; ue.PcfUri == "" {
				logger.GmmLog.Error("AMF can not select an PCF by NRF")
			} else {
				break
			}
		}
		time.Sleep(500 * time.Millisecond) // sleep a while when search NF Instance fail
	}

	problemDetails, err := amf_consumer.AMPolicyControlCreate(ue)
	if problemDetails != nil {
		logger.GmmLog.Errorf("AM Policy Control Create Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		logger.GmmLog.Errorf("AM Policy Control Create Error[%+v]", err)
	}

	// Service Area Restriction are applicable only to 3GPP access
	if anType == models.AccessType__3_GPP_ACCESS {
		if ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
			servAreaRes := ue.AmPolicyAssociation.ServAreaRes
			if servAreaRes.RestrictionType == models.RestrictionType_ALLOWED_AREAS {
				numOfallowedTAs := 0
				for _, area := range servAreaRes.Areas {
					numOfallowedTAs += len(area.Tacs)
				}
				if numOfallowedTAs < int(servAreaRes.MaxNumOfTAs) {
					// TODO: based on AMF Policy, assign additional allowed area for UE,
					// and the upper limit is servAreaRes.MaxNumOfTAs (TS 29.507 4.2.2.3)
				}
			}
		}
	}

	// TODO (step 18 optional): If the AMF has changed and the old AMF has indicated an existing NGAP UE association towards a N3IWF, the new AMF
	// creates an NGAP UE association towards the N3IWF to which the UE is connectedsend N2 AMF mobility request to N3IWF
	if anType == models.AccessType_NON_3_GPP_ACCESS && ue.ServingAmfChanged {
		// TODO: send N2 AMF Mobility Request
	}

	amfSelf.AllocateRegistrationArea(ue, anType)
	logger.GmmLog.Debugf("Allocate GUTI[%s]", ue.Guti)

	assignLadnInfo(ue)

	amfSelf.AddAmfUeToUePool(ue, ue.Supi)
	ue.T3502Value = amfSelf.T3502Value
	if anType == models.AccessType__3_GPP_ACCESS {
		ue.T3512Value = amfSelf.T3512Value
	} else {
		ue.Non3gppDeregistrationTimerValue = amfSelf.Non3gppDeregistrationTimerValue
	}
	gmm_message.SendRegistrationAccept(ue, anType, nil, nil, nil, nil, nil)
	return ue.Sm[anType].Transfer(gmm_state.INITIAL_CONTEXT_SETUP, nil)
}

func HandleMobilityAndPeriodicRegistrationUpdating(ue *amf_context.AmfUe, anType models.AccessType, procedureCode int64) error {

	logger.GmmLog.Infoln("[AMF] Handle MobilityAndPeriodicRegistrationUpdating")

	amfSelf := amf_context.AMF_Self()
	initialContextSetup := false
	if procedureCode == ngapType.ProcedureCodeInitialUEMessage {
		initialContextSetup = true
	}
	// Common Registration procedure
	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	if len(ue.Suci) == 0 {
		gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeSuci)
		return nil
	}

	if !ue.SecurityContextIsValid() {
		return startAuthenticationProcedure(ue, anType)
	}

	if ue.RegistrationRequest.UpdateType5GS != nil {
		if ue.RegistrationRequest.UpdateType5GS.GetNGRanRcu() == nasMessage.NGRanRadioCapabilityUpdateNeeded {
			ue.UeRadioCapability = ""
			ue.UeRadioCapabilityForPaging = nil
		}
	}

	// Registration with AMF re-allocation (TS 23.502 4.2.2.2.3)
	if len(ue.SubscribedNssai) == 0 {
		// TODO(optional): Nudm_SDM_Get(SUPI, slice selection subscription data) to get UE's slice selection subscription info
	}

	var requestedNssai []models.Snssai
	if ue.RegistrationRequest.RequestedNSSAI != nil {
		requestedNssai = nasConvert.RequestedNssaiToModels(ue.RegistrationRequest.RequestedNSSAI)
		logger.GmmLog.Debugf("RequestedNssai: %+v", requestedNssai)
		// TODO(optional): Nnssf_NSSelection_Get(requested nssai, subscribed s-nssai)
		// TODO(optional): Namf_Communication_RegistrationCompleteNotify(failure cause) to old AMF if amf decides to reroute nas msg to another AMF
		// TODO(optional): Nnrf_NfDiscovery_Request(NF Type, Amf Set): if initial AMF does not have target AMF address
		// TODO(if decide to reroute): Namf_Communication_N1MessageNotify to target AMF or Send Reroute Nas Request to RAN
	} else {
		ue.AllowedNssai[anType] = nil
		ue.AllowedNssai[anType] = append(ue.AllowedNssai[anType], amfSelf.PlmnSupportList[0].SNssaiList...)
	}

	// TODO (step 10 optional): send Namf_Communication_RegistrationCompleteNotify to old AMF if need
	if ue.ServingAmfChanged {
		// If the AMF has changed the new AMF notifies the old AMF that the registration of the UE in the new AMF is completed
	}

	if len(ue.Pei) == 0 {
		gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeImei)
		return nil
	}

	// TODO (step 12 optional): the new AMF initiates ME identity check by invoking the N5g-eir_EquipmentIdentityCheck_Get service operation

	if ue.ServingAmfChanged || ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Check(gmm_state.REGISTERED) || !ue.ContextValid {
		// UDM selection described in TS 23.501 6.3.8
		// TODO: consider udm group id, Routing ID part of SUCI, GPSI or External Group ID (e.g., by the NEF)
		param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
			Supi: optional.NewString(ue.Supi),
		}
		resp, err := amf_consumer.SendSearchNFInstances(amfSelf.NrfUri, models.NfType_UDM, models.NfType_AMF, param)
		if err != nil {
			logger.GmmLog.Error("AMF can not select an UDM by NRF")
			return err
		}

		var uecmUri, sdmUri string
		for _, nfProfile := range resp.NfInstances {
			ue.UdmId = nfProfile.NfInstanceId
			uecmUri = amf_util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDM_UECM, models.NfServiceStatus_REGISTERED)
			sdmUri = amf_util.SearchNFServiceUri(nfProfile, models.ServiceName_NUDM_SDM, models.NfServiceStatus_REGISTERED)
			if uecmUri != "" && sdmUri != "" {
				break
			}
		}
		ue.NudmUECMUri = uecmUri
		ue.NudmSDMUri = sdmUri

		problemDetails, err := amf_consumer.UeCmRegistration(ue, anType, true)
		if problemDetails != nil {
			logger.GmmLog.Errorf("UECM_Registration Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("UECM_Registration Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMGetAmData(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM_Get AmData Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM_Get AmData Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMGetSmfSelectData(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM_Get SmfSelectData Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM_Get SmfSelectData Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMGetUeContextInSmfData(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM_Get UeContextInSmfData Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM_Get UeContextInSmfData Error[%+v]", err)
		}

		problemDetails, err = amf_consumer.SDMSubscribe(ue)
		if problemDetails != nil {
			logger.GmmLog.Errorf("SDM Subscribe Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.GmmLog.Errorf("SDM Subscribe Error[%+v]", err)
		}

		ue.ContextValid = true
	}

	var reactivationResult *[16]bool
	var errPduSessionId, errCause []uint8
	ctxList := ngapType.PDUSessionResourceSetupListCxtReq{}
	suList := ngapType.PDUSessionResourceSetupListSUReq{}

	// TODO: determine if UE is in non-allowed area
	if ue.RegistrationRequest.UplinkDataStatus != nil {
		uplinkDataPsi := nasConvert.PSIToBooleanArray(ue.RegistrationRequest.UplinkDataStatus.Buffer)
		reactivationResult = new([16]bool)

		for psi := 1; psi <= 15; psi++ {
			pduSessionId := int32(psi)
			if _, exists := ue.SmContextList[pduSessionId]; !exists {
				continue
			}

			smContext := ue.SmContextList[pduSessionId]
			if uplinkDataPsi[psi] && smContext.PduSessionContext.AccessType == models.AccessType__3_GPP_ACCESS { // uplink data are pending for the corresponding PDU session identity
				response, errResponse, problemDetail, err := amf_consumer.SendUpdateSmContextActivateUpCnxState(ue, pduSessionId, anType)
				if response == nil {
					reactivationResult[pduSessionId] = true
					errPduSessionId = append(errPduSessionId, uint8(pduSessionId))
					cause := nasMessage.Cause5GMMProtocolErrorUnspecified
					if errResponse != nil {
						switch errResponse.JsonData.Error.Cause {
						case "OUT_OF_LADN_SERVICE_AREA":
							cause = nasMessage.Cause5GMMLADNNotAvailable
						case "PRIORITIZED_SERVICES_ONLY":
							cause = nasMessage.Cause5GMMRestrictedServiceArea
						case "DNN_CONGESTION", "S-NSSAI_CONGESTION":
							cause = nasMessage.Cause5GMMInsufficientUserPlaneResourcesForThePDUSession
						}
					}
					errCause = append(errCause, cause)

					if problemDetail != nil {
						logger.GmmLog.Errorf("Update SmContext Failed Problem[%+v]", problemDetail)
					} else if err != nil {
						logger.GmmLog.Errorf("Update SmContext Error[%v]", err.Error())
					}
				} else {
					if initialContextSetup {
						ngap_message.AppendPDUSessionResourceSetupListCxtReq(&ctxList, pduSessionId, *smContext.PduSessionContext.SNssai, response.BinaryDataN1SmMessage, response.BinaryDataN2SmInformation)
					} else {
						ngap_message.AppendPDUSessionResourceSetupListSUReq(&suList, pduSessionId, *smContext.PduSessionContext.SNssai, response.BinaryDataN1SmMessage, response.BinaryDataN2SmInformation)
					}
				}
			}
		}
	}

	var pduSessionStatus *[16]bool
	if ue.RegistrationRequest.PDUSessionStatus != nil {
		pduSessionStatus = new([16]bool)
		psiArray := nasConvert.PSIToBooleanArray(ue.RegistrationRequest.PDUSessionStatus.Buffer)
		for psi := 1; psi <= 15; psi++ {
			pduSessionId := int32(psi)
			if _, exists := ue.SmContextList[pduSessionId]; !exists {
				continue
			}
			if !psiArray[psi] && ue.SmContextList[pduSessionId].PduSessionContext.AccessType == anType {
				cause := models.Cause_PDU_SESSION_STATUS_MISMATCH
				causeAll := &amf_context.CauseAll{
					Cause: &cause,
				}
				releaseData := amf_consumer.BuildReleaseSmContextRequest(ue, causeAll, "", nil)
				problemDetail, err := amf_consumer.SendReleaseSmContextRequest(ue, pduSessionId, releaseData)
				if problemDetail != nil {
					pduSessionStatus[psi] = true
					logger.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
				} else if err != nil {
					pduSessionStatus[psi] = true
					logger.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
				} else {
					pduSessionStatus[psi] = false
				}
			} else {
				pduSessionStatus[psi] = false
			}
		}
	}

	if ue.RegistrationRequest.AllowedPDUSessionStatus != nil {
		allowedPsis := nasConvert.PSIToBooleanArray(ue.RegistrationRequest.AllowedPDUSessionStatus.Buffer)
		if ue.N1N2Message != nil {
			requestData := ue.N1N2Message.Request.JsonData
			n1Msg := ue.N1N2Message.Request.BinaryDataN1Message
			n2Info := ue.N1N2Message.Request.BinaryDataN2Information

			// downlink signalling
			if n2Info == nil {
				if len(suList.List) != 0 {
					nasPdu, err := gmm_message.BuildRegistrationAccept(ue, anType, pduSessionStatus, reactivationResult, errPduSessionId, errCause)
					if err != nil {
						return err
					}
					ngap_message.SendPDUSessionResourceSetupRequest(ue.RanUe[anType], nasPdu, suList)
				} else {
					gmm_message.SendRegistrationAccept(ue, anType, pduSessionStatus, reactivationResult, errPduSessionId, errCause, &ctxList)
				}
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
				ue.N1N2Message = nil
				return nil
			}

			smInfo := requestData.N2InfoContainer.SmInfo
			smContext, exist := ue.SmContextList[requestData.PduSessionId]
			if !exist {
				ue.N1N2Message = nil
				return fmt.Errorf("Pdu Session Id not Exists")
			}

			if smContext.PduSessionContext.AccessType == models.AccessType_NON_3_GPP_ACCESS {
				if reactivationResult == nil {
					reactivationResult = new([16]bool)
				}
				if allowedPsis[requestData.PduSessionId] {
					// TODO: error handling
					response, errRes, _, _ := amf_consumer.SendUpdateSmContextChangeAccessType(ue, requestData.PduSessionId, true)
					if response == nil {
						reactivationResult[requestData.PduSessionId] = true
						errPduSessionId = append(errPduSessionId, uint8(requestData.PduSessionId))
						cause := nasMessage.Cause5GMMProtocolErrorUnspecified
						if errRes != nil {
							switch errRes.JsonData.Error.Cause {
							case "OUT_OF_LADN_SERVICE_AREA":
								cause = nasMessage.Cause5GMMLADNNotAvailable
							case "PRIORITIZED_SERVICES_ONLY":
								cause = nasMessage.Cause5GMMRestrictedServiceArea
							case "DNN_CONGESTION", "S-NSSAI_CONGESTION":
								cause = nasMessage.Cause5GMMInsufficientUserPlaneResourcesForThePDUSession
							}
						}
						errCause = append(errCause, cause)
					} else {
						smContext.UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
						smContext.PduSessionContext.AccessType = models.AccessType__3_GPP_ACCESS
						if response.BinaryDataN2SmInformation != nil && response.JsonData.N2SmInfoType == models.N2SmInfoType_PDU_RES_SETUP_REQ {
							ngap_message.AppendPDUSessionResourceSetupListSUReq(&suList, requestData.PduSessionId, *smContext.PduSessionContext.SNssai, nil, response.BinaryDataN2SmInformation)
						}
					}
				} else {
					logger.GmmLog.Warnf("UE was reachable but did not accept to re-activate the PDU Session[%d]", requestData.PduSessionId)
					amf_producer_callback.SendN1N2TransferFailureNotification(ue, models.N1N2MessageTransferCause_UE_NOT_REACHABLE_FOR_SESSION)
				}
			} else if smInfo.N2InfoContent.NgapIeType == models.NgapIeType_PDU_RES_SETUP_REQ {
				var nasPdu []byte
				var err error
				if n1Msg != nil {
					pduSessionId := uint8(smInfo.PduSessionId)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						return err
					}
				}
				ngap_message.AppendPDUSessionResourceSetupListSUReq(&suList, smInfo.PduSessionId, *smInfo.SNssai, nasPdu, n2Info)
			}
		}
	}

	// TODO (step 18 optional): If the AMF has changed and the old AMF has indicated an existing NGAP UE association towards a N3IWF, the new AMF
	// creates an NGAP UE association towards the N3IWF to which the UE is connectedsend N2 AMF mobility request to N3IWF
	if anType == models.AccessType_NON_3_GPP_ACCESS && ue.ServingAmfChanged {
		// TODO: send N2 AMF Mobility Request
	}

	amfSelf.AllocateRegistrationArea(ue, anType)
	assignLadnInfo(ue)

	// TODO: GUTI reassignment if need (based on operator poilcy)
	// TODO: T3512/Non3GPP de-registration timer reassignment if need (based on operator policy)

	if procedureCode == ngapType.ProcedureCodeInitialUEMessage {
		gmm_message.SendRegistrationAccept(ue, anType, pduSessionStatus, reactivationResult, errPduSessionId, errCause, &ctxList)
		return ue.Sm[anType].Transfer(gmm_state.INITIAL_CONTEXT_SETUP, nil)
	} else if procedureCode == ngapType.ProcedureCodeUplinkNASTransport {
		nasPdu, err := gmm_message.BuildRegistrationAccept(ue, anType, pduSessionStatus, reactivationResult, errPduSessionId, errCause)
		if err != nil {
			logger.GmmLog.Error(err.Error())
		}
		if len(suList.List) != 0 {
			ngap_message.SendPDUSessionResourceSetupRequest(ue.RanUe[anType], nasPdu, suList)
		} else {
			ngap_message.SendDownlinkNasTransport(ue.RanUe[anType], nasPdu)
		}
		ue.ClearRegistrationRequestData()
		return ue.Sm[anType].Transfer(gmm_state.REGISTERED, nil)
	}

	return nil
}

func assignLadnInfo(ue *amf_context.AmfUe) {

	amfSelf := amf_context.AMF_Self()

	ue.LadnInfo = nil
	if ue.RegistrationRequest.LADNIndication != nil {
		// request for LADN information
		if ue.RegistrationRequest.LADNIndication.GetLen() == 0 {
			if ue.HasWildCardSubscribedDNN() {
				for _, ladn := range amfSelf.LadnPool {
					if ue.TaiListInRegistrationArea(ladn.TaiLists) {
						ue.LadnInfo = append(ue.LadnInfo, *ladn)
					}
				}
			} else {
				for _, snssaiInfos := range ue.SmfSelectionData.SubscribedSnssaiInfos {
					for _, dnnInfo := range snssaiInfos.DnnInfos {
						if ladn, ok := amfSelf.LadnPool[dnnInfo.Dnn]; ok { // check if this dnn is a ladn
							if ue.TaiListInRegistrationArea(ladn.TaiLists) {
								ue.LadnInfo = append(ue.LadnInfo, *ladn)
							}
						}
					}
				}
			}
		} else {
			requestedLadnList := nasConvert.LadnToModels(ue.RegistrationRequest.LADNIndication.GetLADNDNNValue())
			for _, requestedLadn := range requestedLadnList {
				if ladn, ok := amfSelf.LadnPool[requestedLadn]; ok {
					if ue.TaiListInRegistrationArea(ladn.TaiLists) {
						ue.LadnInfo = append(ue.LadnInfo, *ladn)
					}
				}
			}
		}
	} else if ue.SmfSelectionData != nil {
		for _, snssaiInfos := range ue.SmfSelectionData.SubscribedSnssaiInfos {
			for _, dnnInfo := range snssaiInfos.DnnInfos {
				if dnnInfo.Dnn != "*" {
					if ladn, ok := amfSelf.LadnPool[dnnInfo.Dnn]; ok {
						if ue.TaiListInRegistrationArea(ladn.TaiLists) {
							ue.LadnInfo = append(ue.LadnInfo, *ladn)
						}
					}
				}
			}
		}
	}
}

func HandleIdentityResponse(ue *amf_context.AmfUe, identityResponse *nasMessage.IdentityResponse) error {

	logger.GmmLog.Info("[AMF] Handle Identity Response")

	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	mobileIdentityContents := identityResponse.MobileIdentity.GetMobileIdentityContents()
	switch nasConvert.GetTypeOfIdentity(mobileIdentityContents[0]) { // get type of identity
	case nasMessage.MobileIdentity5GSTypeSuci:
		ue.Suci = mobileIdentityContents
		var plmnId string
		// TODO: handle the situation if suci is encrypted
		plmnId, ue.Supi = nasConvert.GetPlmnIdAndSupiFromSuci(mobileIdentityContents)
		ue.PlmnId = amf_util.PlmnIdStringToModels(plmnId)
		logger.GmmLog.Debugf("get SUCI[SUPI: %s]", ue.Supi)
	case nasMessage.MobileIdentity5GSType5gGuti:
		_, guti := nasConvert.GutiToString(mobileIdentityContents)
		ue.Guti = guti
		logger.GmmLog.Debugf("get GUTI: %s", guti)
	case nasMessage.MobileIdentity5GSType5gSTmsi:
		sTmsi := hex.EncodeToString(mobileIdentityContents[1:])
		tmp, _ := strconv.ParseInt(sTmsi[4:], 10, 32)
		ue.Tmsi = int32(tmp)
		logger.GmmLog.Debugf("get 5G-S-TMSI: %s", sTmsi)
	case nasMessage.MobileIdentity5GSTypeImei:
		imei := nasConvert.PeiToString(mobileIdentityContents)
		ue.Pei = imei
		logger.GmmLog.Debugf("get PEI: %s", imei)
	case nasMessage.MobileIdentity5GSTypeImeisv:
		imeisv := nasConvert.PeiToString(mobileIdentityContents)
		ue.Pei = imeisv
		logger.GmmLog.Debugf("get PEI: %s", imeisv)
	}
	return nil
}

// TS 24501 5.6.3.2
func HandleNotificationResponse(ue *amf_context.AmfUe, notificationResponse *nasMessage.NotificationResponse) error {

	logger.GmmLog.Info("[AMF] Handle Notification Response")
	amf_util.ClearT3565(ue)
	if notificationResponse != nil && notificationResponse.PDUSessionStatus != nil {
		psiArray := nasConvert.PSIToBooleanArray(notificationResponse.PDUSessionStatus.Buffer)
		for psi := 1; psi <= 15; psi++ {
			pduSessionId := int32(psi)
			if _, exists := ue.SmContextList[pduSessionId]; !exists {
				continue
			}
			if !psiArray[psi] {
				cause := models.Cause_PDU_SESSION_STATUS_MISMATCH
				causeAll := &amf_context.CauseAll{
					Cause: &cause,
				}
				releaseData := amf_consumer.BuildReleaseSmContextRequest(ue, causeAll, "", nil)
				problemDetail, err := amf_consumer.SendReleaseSmContextRequest(ue, pduSessionId, releaseData)
				if problemDetail != nil {
					logger.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
				} else if err != nil {
					logger.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())

				}
			}
		}
	}
	return nil
}

func HandleConfigurationUpdateComplete(ue *amf_context.AmfUe, configurationUpdateComplete *nasMessage.ConfigurationUpdateComplete) error {

	logger.GmmLog.Info("[AMF] Handle Configuration Update Complete")

	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	// TODO: Stop timer T3555 in TS 24.501 Figure 5.4.4.1.1 in handler
	// TODO: Send acknowledgment by Nudm_SMD_Info_Service to UDM in handler
	//		import "free5gc/lib/Nudm_SubscriberDataManagement" client.Info

	return nil
}

func startAuthenticationProcedure(ue *amf_context.AmfUe, anType models.AccessType) error {

	logger.GmmLog.Info("Start authentication procedure")

	amfSelf := amf_context.AMF_Self()
	// TODO: consider ausf group id, Routing ID part of SUCI
	param := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		Supi: optional.NewString(ue.Supi),
	}
	resp, err := amf_consumer.SendSearchNFInstances(amfSelf.NrfUri, models.NfType_AUSF, models.NfType_AMF, param)
	if err != nil {
		logger.GmmLog.Error("AMF can not select an AUSF by NRF")
		return err
	}

	// select the first AUSF, TODO: select base on other info
	var ausfUri string
	for _, nfProfile := range resp.NfInstances {
		ue.AusfId = nfProfile.NfInstanceId
		ausfUri = amf_util.SearchNFServiceUri(nfProfile, models.ServiceName_NAUSF_AUTH, models.NfServiceStatus_REGISTERED)
		if ausfUri != "" {
			break
		}
	}
	if ausfUri == "" {
		err := fmt.Errorf("AMF can not select an AUSF by NRF")
		logger.GmmLog.Errorf(err.Error())
		return err
	}
	ue.AusfUri = ausfUri

	response, problemDetails, err := amf_consumer.SendUEAuthenticationAuthenticateRequest(ue, nil)
	if err != nil {
		return err
	} else if problemDetails != nil {
		logger.GmmLog.Errorf("Nausf_UEAU Authenticate Request Error[Problem Detail: %+v]", problemDetails)
		return nil
	}
	ue.AuthenticationCtx = response
	ue.ABBA = []uint8{0x00, 0x00} // set ABBA value as described at TS 33.501 Annex A.7.1

	ue.T3560RetryTimes = 0
	gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
	return ue.Sm[anType].Transfer(gmm_state.AUTHENTICATION, nil)
}

// TS 24501 5.6.1
func HandleServiceRequest(ue *amf_context.AmfUe, anType models.AccessType, procedureCode int64, serviceRequest *nasMessage.ServiceRequest) (err error) {
	logger.GmmLog.Info("Handle Service Reqeust")

	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	// Clear Timer
	amf_util.ClearT3513(ue)
	amf_util.ClearT3565(ue)

	// Send Authtication / Security Procedure not support
	if !ue.SecurityContextIsValid() {
		logger.GmmLog.Warnf("No Security Context : SUPI[%s]", ue.Supi)
		gmm_message.SendServiceReject(ue.RanUe[anType], nil, nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork)
		ngap_message.SendUEContextReleaseCommand(ue.RanUe[anType], amf_context.UeContextN2NormalRelease, ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
		return nil
	}

	serviceType := serviceRequest.GetServiceTypeValue()
	var reactivationResult, acceptPduSessionPsi *[16]bool
	var errPduSessionId, errCause []uint8
	var targetPduSessionId int32
	suList := ngapType.PDUSessionResourceSetupListSUReq{}
	ctxList := ngapType.PDUSessionResourceSetupListCxtReq{}
	initCxt := procedureCode == ngapType.ProcedureCodeInitialUEMessage
	if serviceType == nasMessage.ServiceTypeSignalling {
		err = sendServiceAccept(initCxt, ue, anType, ctxList, suList, nil, nil, nil, nil)
		return
	}
	if ue.N1N2Message != nil {
		requestData := ue.N1N2Message.Request.JsonData
		if ue.N1N2Message.Request.BinaryDataN2Information != nil {
			if requestData.N2InfoContainer.N2InformationClass == models.N2InformationClass_SM {
				targetPduSessionId = requestData.N2InfoContainer.SmInfo.PduSessionId
			} else {
				ue.N1N2Message = nil
				return fmt.Errorf("Service Request triggered by Network has not implemented about non SM N2Info")
			}
		}
	}

	if serviceRequest.UplinkDataStatus != nil {
		uplinkDataPsi := nasConvert.PSIToBooleanArray(serviceRequest.UplinkDataStatus.Buffer)
		reactivationResult = new([16]bool)
		for pduSessionId, smContext := range ue.SmContextList {
			if pduSessionId == targetPduSessionId {
				continue
			}
			if uplinkDataPsi[pduSessionId] && smContext.PduSessionContext.AccessType == models.AccessType__3_GPP_ACCESS {
				response, errRes, _, _ := amf_consumer.SendUpdateSmContextActivateUpCnxState(ue, pduSessionId, models.AccessType__3_GPP_ACCESS)
				if response == nil {
					reactivationResult[pduSessionId] = true
					errPduSessionId = append(errPduSessionId, uint8(pduSessionId))
					cause := nasMessage.Cause5GMMProtocolErrorUnspecified
					if errRes != nil {
						switch errRes.JsonData.Error.Cause {
						case "OUT_OF_LADN_SERVICE_AREA":
							cause = nasMessage.Cause5GMMLADNNotAvailable
						case "PRIORITIZED_SERVICES_ONLY":
							cause = nasMessage.Cause5GMMRestrictedServiceArea
						case "DNN_CONGESTION", "S-NSSAI_CONGESTION":
							cause = nasMessage.Cause5GMMInsufficientUserPlaneResourcesForThePDUSession
						}
					}
					errCause = append(errCause, cause)
				} else if initCxt {
					ngap_message.AppendPDUSessionResourceSetupListCxtReq(&ctxList, pduSessionId, *smContext.PduSessionContext.SNssai, nil, response.BinaryDataN2SmInformation)
				} else {
					ngap_message.AppendPDUSessionResourceSetupListSUReq(&suList, pduSessionId, *smContext.PduSessionContext.SNssai, nil, response.BinaryDataN2SmInformation)

				}
			}
		}
	}
	if serviceRequest.PDUSessionStatus != nil {
		acceptPduSessionPsi = new([16]bool)
		psiArray := nasConvert.PSIToBooleanArray(serviceRequest.PDUSessionStatus.Buffer)
		for pduSessionId, smContext := range ue.SmContextList {
			if smContext.PduSessionContext.AccessType == anType {
				if !psiArray[pduSessionId] {
					cause := models.Cause_PDU_SESSION_STATUS_MISMATCH
					causeAll := &amf_context.CauseAll{
						Cause: &cause,
					}
					releaseData := amf_consumer.BuildReleaseSmContextRequest(ue, causeAll, "", nil)
					problemDetail, err := amf_consumer.SendReleaseSmContextRequest(ue, pduSessionId, releaseData)
					if problemDetail != nil {
						logger.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
					} else if err != nil {
						logger.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
					}
				} else {
					acceptPduSessionPsi[pduSessionId] = true
				}
			}
		}
	}
	switch serviceType {
	case nasMessage.ServiceTypeMobileTerminatedServices: // Trigger by Network
		if ue.N1N2Message != nil {
			requestData := ue.N1N2Message.Request.JsonData
			n1Msg := ue.N1N2Message.Request.BinaryDataN1Message
			n2Info := ue.N1N2Message.Request.BinaryDataN2Information

			// downlink signalling
			if n2Info == nil {
				err = sendServiceAccept(initCxt, ue, anType, ctxList, suList, acceptPduSessionPsi, reactivationResult, errPduSessionId, errCause)
				if err != nil {
					return
				}
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
				ue.N1N2Message = nil
				return nil
			}
			// TODO: Area of validity for the N2 SM information
			smInfo := requestData.N2InfoContainer.SmInfo
			smContext, exist := ue.SmContextList[requestData.PduSessionId]
			if !exist {
				ue.N1N2Message = nil
				return fmt.Errorf("Service Request triggered by Network error for pduSessionId does not exist")
			}

			if smContext.PduSessionContext.AccessType == models.AccessType_NON_3_GPP_ACCESS {
				if serviceRequest.AllowedPDUSessionStatus != nil {
					allowPduSessionPsi := nasConvert.PSIToBooleanArray(serviceRequest.AllowedPDUSessionStatus.Buffer)
					if reactivationResult == nil {
						reactivationResult = new([16]bool)
					}
					if allowPduSessionPsi[requestData.PduSessionId] {
						response, errRes, _, _ := amf_consumer.SendUpdateSmContextChangeAccessType(ue, requestData.PduSessionId, true)
						if response == nil {
							reactivationResult[requestData.PduSessionId] = true
							errPduSessionId = append(errPduSessionId, uint8(requestData.PduSessionId))
							cause := nasMessage.Cause5GMMProtocolErrorUnspecified
							if errRes != nil {
								switch errRes.JsonData.Error.Cause {
								case "OUT_OF_LADN_SERVICE_AREA":
									cause = nasMessage.Cause5GMMLADNNotAvailable
								case "PRIORITIZED_SERVICES_ONLY":
									cause = nasMessage.Cause5GMMRestrictedServiceArea
								case "DNN_CONGESTION", "S-NSSAI_CONGESTION":
									cause = nasMessage.Cause5GMMInsufficientUserPlaneResourcesForThePDUSession
								}
							}
							errCause = append(errCause, cause)
						} else {
							smContext.UserLocation = deepcopy.Copy(ue.Location).(models.UserLocation)
							smContext.PduSessionContext.AccessType = models.AccessType__3_GPP_ACCESS
							if response.BinaryDataN2SmInformation != nil && response.JsonData.N2SmInfoType == models.N2SmInfoType_PDU_RES_SETUP_REQ {
								if initCxt {
									ngap_message.AppendPDUSessionResourceSetupListCxtReq(&ctxList, requestData.PduSessionId, *smContext.PduSessionContext.SNssai, nil, response.BinaryDataN2SmInformation)
								} else {
									ngap_message.AppendPDUSessionResourceSetupListSUReq(&suList, requestData.PduSessionId, *smContext.PduSessionContext.SNssai, nil, response.BinaryDataN2SmInformation)
								}
							}
						}
					} else {
						logger.GmmLog.Warnf("UE was reachable but did not accept to re-activate the PDU Session[%d]", requestData.PduSessionId)
						amf_producer_callback.SendN1N2TransferFailureNotification(ue, models.N1N2MessageTransferCause_UE_NOT_REACHABLE_FOR_SESSION)
					}
				}
			} else if smInfo.N2InfoContent.NgapIeType == models.NgapIeType_PDU_RES_SETUP_REQ {
				var nasPdu []byte
				if n1Msg != nil {
					pduSessionId := uint8(smInfo.PduSessionId)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, &pduSessionId, nil, nil, 0)
					if err != nil {
						return err
					}
				}
				if initCxt {
					ngap_message.AppendPDUSessionResourceSetupListCxtReq(&ctxList, smInfo.PduSessionId, *smInfo.SNssai, nasPdu, n2Info)
				} else {
					ngap_message.AppendPDUSessionResourceSetupListSUReq(&suList, smInfo.PduSessionId, *smInfo.SNssai, nasPdu, n2Info)
				}
			}
			err = sendServiceAccept(initCxt, ue, anType, ctxList, suList, acceptPduSessionPsi, reactivationResult, errPduSessionId, errCause)
			if err != nil {
				return
			}

		}
	case nasMessage.ServiceTypeData:
		if anType == models.AccessType__3_GPP_ACCESS {
			if ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
				var accept bool
				switch ue.AmPolicyAssociation.ServAreaRes.RestrictionType {
				case models.RestrictionType_ALLOWED_AREAS:
					accept = amf_context.TacInAreas(ue.Tai.Tac, ue.AmPolicyAssociation.ServAreaRes.Areas)
				case models.RestrictionType_NOT_ALLOWED_AREAS:
					accept = !amf_context.TacInAreas(ue.Tai.Tac, ue.AmPolicyAssociation.ServAreaRes.Areas)
				}

				if !accept {
					gmm_message.SendServiceReject(ue.RanUe[anType], nil, nasMessage.Cause5GMMRestrictedServiceArea)
					return nil
				}
			}
			err = sendServiceAccept(initCxt, ue, anType, ctxList, suList, acceptPduSessionPsi, reactivationResult, errPduSessionId, errCause)
			if err != nil {
				return
			}
		} else {
			err = sendServiceAccept(initCxt, ue, anType, ctxList, suList, acceptPduSessionPsi, reactivationResult, errPduSessionId, errCause)
			if err != nil {
				return
			}
		}
	default:
		return fmt.Errorf("Service Type[%d] is not supported", serviceType)
	}
	if len(errPduSessionId) != 0 {
		logger.GmmLog.Info(errPduSessionId, errCause)
	}
	ue.N1N2Message = nil
	return nil
}

func sendServiceAccept(initCxt bool, ue *amf_context.AmfUe, anType models.AccessType, ctxList ngapType.PDUSessionResourceSetupListCxtReq, suList ngapType.PDUSessionResourceSetupListSUReq, pDUSessionStatus *[16]bool, reactivationResult *[16]bool, errPduSessionId, errCause []uint8) error {
	if initCxt {
		// update Kgnb/Kn3iwf
		ue.UpdateSecurityContext(anType)

		nasPdu, err := gmm_message.BuildServiceAccept(ue, pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
		if err != nil {
			return err
		}
		if len(ctxList.List) != 0 {
			ngap_message.SendInitialContextSetupRequest(ue, anType, nasPdu, nil, &ctxList, nil, nil, nil)
		} else {
			ngap_message.SendInitialContextSetupRequest(ue, anType, nasPdu, nil, nil, nil, nil, nil)
		}
	} else if len(suList.List) != 0 {
		nasPdu, err := gmm_message.BuildServiceAccept(ue, pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
		if err != nil {
			return err
		}
		ngap_message.SendPDUSessionResourceSetupRequest(ue.RanUe[anType], nasPdu, suList)
	} else {
		gmm_message.SendServiceAccept(ue.RanUe[anType], pDUSessionStatus, reactivationResult, errPduSessionId, errCause)
	}
	return nil
}

func HandleAuthenticationResponse(ue *amf_context.AmfUe, anType models.AccessType, authenticationResponse *nasMessage.AuthenticationResponse) error {

	logger.GmmLog.Info("[AMF] Handle Authentication Response")

	amf_util.ClearT3560(ue)

	if ue.AuthenticationCtx == nil {
		return fmt.Errorf("Ue Authentication Context is nil")
	}

	switch ue.AuthenticationCtx.AuthType {
	case models.AuthType__5_G_AKA:
		var av5gAka models.Av5gAka
		if err := mapstructure.Decode(ue.AuthenticationCtx.Var5gAuthData, &av5gAka); err != nil {
			return fmt.Errorf("Var5gAuthData Convert Type Error")
		}
		resStar := authenticationResponse.AuthenticationResponseParameter.GetRES()

		// Calculate HRES* (TS 33.501 Annex A.5)
		p0, _ := hex.DecodeString(av5gAka.Rand)
		p1 := resStar[:]
		concat := append(p0, p1...)
		hResStarBytes := sha256.Sum256(concat)
		hResStar := hex.EncodeToString(hResStarBytes[16:])

		if hResStar != av5gAka.HxresStar {
			logger.GmmLog.Errorf("HRES* Validation Failure")

			if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti {
				gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeSuci)
				return ue.Sm[anType].Transfer(gmm_state.DE_REGISTERED, nil)
			} else {
				ue.ClearRegistrationRequestData()
				gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
				return ue.Sm[anType].Transfer(gmm_state.DE_REGISTERED, nil)
			}
		}

		response, problemDetails, err := amf_consumer.SendAuth5gAkaConfirmRequest(ue, hex.EncodeToString(resStar[:]))
		if err != nil {
			return err
		} else if problemDetails != nil {
			logger.GmmLog.Debugf("Auth5gAkaConfirm Error[Problem Detail: %+v]", problemDetails)
			return nil
		}
		switch response.AuthResult {
		case models.AuthResult_SUCCESS:
			ue.UnauthenticatedSupi = false
			ue.Kseaf = response.Kseaf
			ue.DerivateKamf()
			gmm_message.SendSecurityModeCommand(ue.RanUe[anType], false, "")
			return ue.Sm[anType].Transfer(gmm_state.SECURITY_MODE, nil)
		case models.AuthResult_FAILURE:
			if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti {
				gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeSuci)
				return ue.Sm[anType].Transfer(gmm_state.DE_REGISTERED, nil)
			} else {
				ue.ClearRegistrationRequestData()
				gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
				return ue.Sm[anType].Transfer(gmm_state.DE_REGISTERED, nil)
			}
		}
	case models.AuthType_EAP_AKA_PRIME:
		response, problemDetails, err := amf_consumer.SendEapAuthConfirmRequest(ue, *authenticationResponse.EAPMessage)
		if err != nil {
			return err
		} else if problemDetails != nil {
			logger.GmmLog.Debugf("EapAuthConfirm Error[Problem Detail: %+v]", problemDetails)
			return nil
		}

		switch response.AuthResult {
		case models.AuthResult_SUCCESS:
			ue.UnauthenticatedSupi = false
			ue.Kseaf = response.KSeaf
			ue.DerivateKamf()
			// TODO: select enc/int algorithm based on ue security capability & amf's policy,
			// then generate KnasEnc, KnasInt
			gmm_message.SendSecurityModeCommand(ue.RanUe[anType], true, response.EapPayload)
			return ue.Sm[anType].Transfer(gmm_state.SECURITY_MODE, nil)
		case models.AuthResult_FAILURE:
			if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti {
				gmm_message.SendAuthenticationResult(ue.RanUe[anType], false, response.EapPayload)
				gmm_message.SendIdentityRequest(ue.RanUe[anType], nasMessage.MobileIdentity5GSTypeSuci)
				return ue.Sm[anType].Transfer(gmm_state.DE_REGISTERED, nil)
			} else {
				gmm_message.SendAuthenticationReject(ue.RanUe[anType], response.EapPayload)
			}
			return ue.Sm[anType].Transfer(gmm_state.DE_REGISTERED, nil)
		case models.AuthResult_ONGOING:
			ue.AuthenticationCtx.Var5gAuthData = response.EapPayload
			if _, exists := response.Links["link"]; exists {
				ue.AuthenticationCtx.Links = response.Links
			}
			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		}
	}

	return nil
}

func HandleAuthenticationFailure(ue *amf_context.AmfUe, anType models.AccessType, authenticationFailure *nasMessage.AuthenticationFailure) error {

	logger.GmmLog.Info("[AMF] Handle Authentication Failure")

	amf_util.ClearT3560(ue)

	cause5GMM := authenticationFailure.Cause5GMM.GetCauseValue()

	if ue.AuthenticationCtx.AuthType == models.AuthType__5_G_AKA {
		switch cause5GMM {
		case nasMessage.Cause5GMMMACFailure:
			logger.GmmLog.Warn("5GMM Cause[Mac Failure]")
			ue.AuthFailureCauseSynchFailureTimes = 0
			gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
			return ue.Sm[anType].Transfer(gmm_state.EXCEPTION, nil)
		case nasMessage.Cause5GMMNon5GAuthenticationUnacceptable:
			logger.GmmLog.Warn("5GMM Cause[Non-5G Authentication Unacceptable]")
			ue.AuthFailureCauseSynchFailureTimes = 0
			gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
			return ue.Sm[anType].Transfer(gmm_state.EXCEPTION, nil)
		case nasMessage.Cause5GMMngKSIAlreadyInUse:
			logger.GmmLog.Warn("5GMM Cause[NgKSI Already In Use]")
			ue.AuthFailureCauseSynchFailureTimes = 0
			logger.GmmLog.Debug("Select new NgKsi")
			// select new ngksi
			if ue.NgKsi.Ksi < 6 { // ksi is range from 0 to 6
				ue.NgKsi.Ksi += 1
			} else {
				ue.NgKsi.Ksi = 0
			}
			ue.T3560RetryTimes = 0
			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		case nasMessage.Cause5GMMSynchFailure:
			logger.GmmLog.Warn("5GMM Cause[Synch Failure]")
			ue.AuthFailureCauseSynchFailureTimes++
			if ue.AuthFailureCauseSynchFailureTimes >= 2 {
				logger.GmmLog.Warnf("2 consecutive Synch Failure, terminate authentication procedure")
				ue.AuthFailureCauseSynchFailureTimes = 0
				gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
				return ue.Sm[anType].Transfer(gmm_state.EXCEPTION, nil)
			}
			auts := authenticationFailure.AuthenticationFailureParameter.GetAuthenticationFailureParameter()
			resynchronizationInfo := &models.ResynchronizationInfo{
				Auts: hex.EncodeToString(auts[:]),
			}

			response, problemDetails, err := amf_consumer.SendUEAuthenticationAuthenticateRequest(ue, resynchronizationInfo)
			if err != nil {
				return err
			} else if problemDetails != nil {
				logger.GmmLog.Errorf("Nausf_UEAU Authenticate Request Error[Problem Detail: %+v]", problemDetails)
				return nil
			}
			ue.AuthenticationCtx = response
			ue.ABBA = []uint8{0x00, 0x00}

			ue.T3560RetryTimes = 0
			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		}
	} else if ue.AuthenticationCtx.AuthType == models.AuthType_EAP_AKA_PRIME {
		switch cause5GMM {
		case nasMessage.Cause5GMMngKSIAlreadyInUse:
			logger.GmmLog.Warn("5GMM Cause[NgKSI Already In Use]")
			if ue.NgKsi.Ksi < 6 { // ksi is range from 0 to 6
				ue.NgKsi.Ksi += 1
			} else {
				ue.NgKsi.Ksi = 0
			}
			ue.T3560RetryTimes = 0
			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		}
	}

	return nil
}

func HandleRegistrationComplete(ue *amf_context.AmfUe, anType models.AccessType, registrationComplete *nasMessage.RegistrationComplete) error {

	logger.GmmLog.Info("[AMF] Handle Registration Complete")

	amf_util.ClearT3550(ue)

	if registrationComplete.SORTransparentContainer != nil {
		// TODO: if at regsitration procedure 14b, udm provide amf Steering of Roaming info & request an ack,
		// AMF provides the UE's ack with Nudm_SDM_Info (SOR not supportted in this stage)
	}

	// TODO: if
	//	1. AMF has evaluated the support of IMS Voice over PS Sessions (TS 23.501 5.16.3.2)
	//	2. AMF determines that it needs to update the Homogeneous Support of IMS Voice over PS Sessions (TS 23.501 5.16.3.3)
	// Then invoke Nudm_UECM_Update to send "Homogeneous Support of IMS Voice over PS Sessions" indication to udm

	if ue.RegistrationRequest.UplinkDataStatus == nil && ue.RegistrationRequest.GetFOR() == nasMessage.FollowOnRequestNoPending {
		ngap_message.SendUEContextReleaseCommand(ue.RanUe[anType], amf_context.UeContextN2NormalRelease, ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
	}

	ue.ClearRegistrationRequestData()
	return ue.Sm[anType].Transfer(gmm_state.REGISTERED, nil)
}

// TODO: finish it; TS 33.501 6.7.2
func HandleSecurityModeComplete(ue *amf_context.AmfUe, anType models.AccessType, procedureCode int64, securityModeComplete *nasMessage.SecurityModeComplete) error {

	logger.GmmLog.Info("[AMF] Handle Security Mode Complete")

	// stop T3560
	amf_util.ClearT3560(ue)
	if ue.SecurityContextIsValid() {
		// update Kgnb/Kn3iwf
		ue.UpdateSecurityContext(anType)
	}

	if securityModeComplete.IMEISV != nil {
		ue.Pei = nasConvert.PeiToString(securityModeComplete.IMEISV.Octet[:])
	}

	// TODO: AMF shall set the NAS COUNTs to zero if horizontal derivation of KAMF is performed
	if securityModeComplete.NASMessageContainer != nil {
		nasMessageContainer := securityModeComplete.NASMessageContainer
		m := nas.NewMessage()
		_ = m.GmmMessageDecode(&nasMessageContainer.Buffer)

		switch m.GmmMessage.GmmHeader.GetMessageType() {
		case nas.MsgTypeRegistrationRequest:
			logger.GmmLog.Traceln("[AMF] Handle MsgTypeRegistrationRequest")
			args := make(fsm.Args)
			args[gmm_event.AMF_UE] = ue
			args[gmm_event.PROCEDURE_CODE] = procedureCode
			args[gmm_event.GMM_MESSAGE] = m.GmmMessage
			_ = ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.REGISTERED, nil)
			return ue.Sm[models.AccessType__3_GPP_ACCESS].SendEvent(gmm_event.EVENT_GMM_MESSAGE, args)
		case nas.MsgTypeServiceRequest:
			logger.GmmLog.Traceln("[AMF] Handle MsgTypeServiceRequest")
			args := make(fsm.Args)
			args[gmm_event.AMF_UE] = ue
			args[gmm_event.PROCEDURE_CODE] = procedureCode
			args[gmm_event.GMM_MESSAGE] = m.GmmMessage
			_ = ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.REGISTERED, nil)
			return ue.Sm[models.AccessType__3_GPP_ACCESS].SendEvent(gmm_event.EVENT_GMM_MESSAGE, args)
		default:
			logger.GmmLog.Errorln("nas message container Iei type error")
		}
	} else {
		logger.GmmLog.Info("[AMF] Handle Security Mode Complete NULL")
	}

	return nil
}

func HandleSecurityModeReject(ue *amf_context.AmfUe, anType models.AccessType, securityModeReject *nasMessage.SecurityModeReject) error {

	logger.GmmLog.Info("[AMF] Handle Security Mode Reject")

	// stop T3560
	amf_util.ClearT3560(ue)

	logger.GmmLog.Warnf("Reject Cause: %d", securityModeReject.Cause5GMM.GetCauseValue())

	return ue.Sm[anType].Transfer(gmm_state.EXCEPTION, nil)
}

// TODO: finish it TS 23.502 4.2.2.3
func HandleDeregistrationRequest(ue *amf_context.AmfUe, anType models.AccessType, deregistrationRequest *nasMessage.DeregistrationRequestUEOriginatingDeregistration) error {

	logger.GmmLog.Info("[AMF] Handle Deregistration Request(UE Originating)")

	targetDeregistrationAccessType := deregistrationRequest.GetAccessType()
	for pduSessionId, smContext := range ue.SmContextList {
		if smContext.PduSessionContext.AccessType != anType && targetDeregistrationAccessType != nasMessage.AccessTypeBoth {
			continue
		}

		releaseData := amf_consumer.BuildReleaseSmContextRequest(ue, nil, "", nil)
		problemDetail, err := amf_consumer.SendReleaseSmContextRequest(ue, pduSessionId, releaseData)
		if problemDetail != nil {
			logger.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
		} else if err != nil {
			logger.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
		}
	}

	if ue.AmPolicyAssociation != nil {
		terminateAmPolicyAssocaition := true
		switch anType {
		case models.AccessType__3_GPP_ACCESS:
			terminateAmPolicyAssocaition = ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Check(gmm_state.DE_REGISTERED)
		case models.AccessType_NON_3_GPP_ACCESS:
			terminateAmPolicyAssocaition = ue.Sm[models.AccessType__3_GPP_ACCESS].Check(gmm_state.DE_REGISTERED)
		}

		if terminateAmPolicyAssocaition {
			problemDetails, err := amf_consumer.AMPolicyControlDelete(ue)
			if problemDetails != nil {
				logger.GmmLog.Errorf("AM Policy Control Delete Failed Problem[%+v]", problemDetails)
			} else if err != nil {
				logger.GmmLog.Errorf("AM Policy Control Delete Error[%v]", err.Error())
			}
		}
	}

	// if Deregistration type is not switch-off, send Deregistration Accept
	if deregistrationRequest.GetSwitchOff() == 0 {
		gmm_message.SendDeregistrationAccept(ue.RanUe[anType])
	}

	// TS 23.502 4.2.6, 4.12.3
	switch targetDeregistrationAccessType {
	case nasMessage.AccessType3GPP:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		return ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.DE_REGISTERED, nil)
	case nasMessage.AccessTypeNon3GPP:
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		return ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Transfer(gmm_state.DE_REGISTERED, nil)
	case nasMessage.AccessTypeBoth:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}

		if err := ue.Sm[models.AccessType__3_GPP_ACCESS].Transfer(gmm_state.DE_REGISTERED, nil); err != nil {
			return err
		}
		if err := ue.Sm[models.AccessType_NON_3_GPP_ACCESS].Transfer(gmm_state.DE_REGISTERED, nil); err != nil {
			return err
		}
	}

	return nil
}

// TS 23.502 4.2.2.3
func HandleDeregistrationAccept(ue *amf_context.AmfUe, anType models.AccessType, deregistrationAccept *nasMessage.DeregistrationAcceptUETerminatedDeregistration) error {

	logger.GmmLog.Info("[AMF] Handle Deregistration Accept(UE Terminated)")

	amf_util.ClearT3522(ue)

	switch ue.DeregistrationTargetAccessType {
	case nasMessage.AccessType3GPP:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
	case nasMessage.AccessTypeNon3GPP:
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
	case nasMessage.AccessTypeBoth:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS], amf_context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
	}

	ue.DeregistrationTargetAccessType = 0
	return nil
}

func HandleStatus5GMM(ue *amf_context.AmfUe, anType models.AccessType, status5GMM *nasMessage.Status5GMM) error {

	logger.GmmLog.Info("Handle Staus 5GMM")

	logger.GmmLog.Errorf("Error condition [Cause Value: %d]", status5GMM.Cause5GMM.Octet)

	return nil
}
