package gmm

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mohae/deepcopy"
	"github.com/pkg/errors"

	"github.com/free5gc/amf/internal/context"
	gmm_common "github.com/free5gc/amf/internal/gmm/common"
	gmm_message "github.com/free5gc/amf/internal/gmm/message"
	"github.com/free5gc/amf/internal/logger"
	ngap_message "github.com/free5gc/amf/internal/ngap/message"
	"github.com/free5gc/amf/internal/sbi/consumer"
	callback "github.com/free5gc/amf/internal/sbi/processor/notifier"
	"github.com/free5gc/amf/internal/util"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/ngap/ngapConvert"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/util/fsm"
)

const psiArraySize = 16

func HandleULNASTransport(ue *context.AmfUe, anType models.AccessType,
	ulNasTransport *nasMessage.ULNASTransport,
) error {
	ue.GmmLog.Infoln("Handle UL NAS Transport")

	if ue.MacFailed {
		return fmt.Errorf("NAS message integrity check failed")
	}

	switch ulNasTransport.GetPayloadContainerType() {
	// TS 24.501 5.4.5.2.3 case a)
	case nasMessage.PayloadContainerTypeN1SMInfo:
		return transport5GSMMessage(ue, anType, ulNasTransport)
	case nasMessage.PayloadContainerTypeSMS:
		return fmt.Errorf("PayloadContainerTypeSMS has not been implemented yet in UL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeLPP:
		return fmt.Errorf("PayloadContainerTypeLPP has not been implemented yet in UL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeSOR:
		return fmt.Errorf("PayloadContainerTypeSOR has not been implemented yet in UL NAS TRANSPORT")
	case nasMessage.PayloadContainerTypeUEPolicy:
		ue.GmmLog.Infoln("AMF Transfer UEPolicy To PCF")
		callback.SendN1MessageNotify(ue, models.N1MessageClass_UPDP,
			ulNasTransport.PayloadContainer.GetPayloadContainerContents(), nil)
	case nasMessage.PayloadContainerTypeUEParameterUpdate:
		ue.GmmLog.Infoln("AMF Transfer UEParameterUpdate To UDM")
		upuMac, err := nasConvert.UpuAckToModels(ulNasTransport.PayloadContainer.GetPayloadContainerContents())
		if err != nil {
			return err
		}
		err = consumer.GetConsumer().PutUpuAck(ue, upuMac)
		if err != nil {
			return err
		}
		ue.GmmLog.Debugf("UpuMac[%s] in UPU ACK NAS Msg", upuMac)
	case nasMessage.PayloadContainerTypeMultiplePayload:
		return fmt.Errorf("PayloadContainerTypeMultiplePayload has not been implemented yet in UL NAS TRANSPORT")
	}
	return nil
}

func transport5GSMMessage(ue *context.AmfUe, anType models.AccessType,
	ulNasTransport *nasMessage.ULNASTransport,
) error {
	var pduSessionID int32

	ue.GmmLog.Info("Transport 5GSM Message to SMF")

	smMessage := ulNasTransport.PayloadContainer.GetPayloadContainerContents()

	if id := ulNasTransport.PduSessionID2Value; id != nil {
		pduSessionID = int32(id.GetPduSessionID2Value())
	} else {
		return errors.New("PDU Session ID is nil")
	}

	// case 1): looks up a PDU session routing context for the UE and the PDU session ID IE in case the Old PDU
	// session ID IE is not included
	if ulNasTransport.OldPDUSessionID == nil {
		smContext, smContextExist := ue.SmContextFindByPDUSessionID(pduSessionID)
		requestType := ulNasTransport.RequestType

		if requestType != nil {
			switch requestType.GetRequestTypeValue() {
			case nasMessage.ULNASTransportRequestTypeInitialEmergencyRequest:
				fallthrough
			case nasMessage.ULNASTransportRequestTypeExistingEmergencyPduSession:
				ue.GmmLog.Warnf("Emergency PDU Session is not supported")
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
					smMessage, pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
				return nil
			}
		}

		// AMF has a PDU session routing context for the PDU session ID and the UE
		if smContextExist {
			// case i) Request type IE is either not included
			if requestType == nil {
				return forward5GSMMessageToSMF(ue, anType, pduSessionID, smContext, smMessage)
			}

			switch requestType.GetRequestTypeValue() {
			case nasMessage.ULNASTransportRequestTypeInitialRequest:
				smContext.StoreULNASTransport(ulNasTransport)
				//  perform a local release of the PDU session identified by the PDU session ID and shall request
				// the SMF to perform a local release of the PDU session
				updateData := models.SmfPduSessionSmContextUpdateData{
					Release: true,
					Cause:   models.SmfPduSessionCause_REL_DUE_TO_DUPLICATE_SESSION_ID,
					SmContextStatusUri: fmt.Sprintf("%s"+factory.AmfCallbackResUriPrefix+"/smContextStatus/%s/%d",
						ue.ServingAMF().GetIPv4Uri(), ue.Guti, pduSessionID),
				}
				ue.GmmLog.Warningf("Duplicated PDU session ID[%d]", pduSessionID)
				smContext.SetDuplicatedPduSessionID(true)
				response, _, _, err := consumer.GetConsumer().SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
				if err != nil {
					ue.GmmLog.Errorf("Failed to update smContext, local release SmContext[%d]", pduSessionID)
					ue.SmContextList.Delete(pduSessionID)
					return err
				} else if response == nil {
					ue.GmmLog.Errorf("Response to update smContext is nil, local release SmContext[%d]", pduSessionID)
					ue.SmContextList.Delete(pduSessionID)
				} else if response != nil {
					smContext.SetUserLocation(ue.Location)
					responseData := response.JsonData
					n2Info := response.BinaryDataN2SmInformation
					if n2Info != nil {
						if responseData.N2SmInfoType == models.N2SmInfoType_PDU_RES_REL_CMD {
							ue.GmmLog.Debugln("AMF Transfer NGAP PDU Session Resource Release Command from SMF")
							list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
							ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2Info)
							ngap_message.SendPDUSessionResourceReleaseCommand(ue.RanUe[anType], nil, list)
						}
					}
				}

			// case ii) AMF has a PDU session routing context, and Request type is "existing PDU session"
			case nasMessage.ULNASTransportRequestTypeExistingPduSession:
				if ue.InAllowedNssai(smContext.Snssai(), anType) {
					return forward5GSMMessageToSMF(ue, anType, pduSessionID, smContext, smMessage)
				} else {
					ue.GmmLog.Errorf("S-NSSAI[%v] is not allowed for access type[%s] (PDU Session ID: %d)",
						smContext.Snssai(), anType, pduSessionID)
					gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
						smMessage, pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
				}
			// other requestType: AMF forward the 5GSM message, and the PDU session ID IE towards the SMF identified
			// by the SMF ID of the PDU session routing context
			default:
				return forward5GSMMessageToSMF(ue, anType, pduSessionID, smContext, smMessage)
			}
		} else { // AMF does not have a PDU session routing context for the PDU session ID and the UE
			if requestType == nil {
				ue.GmmLog.Warnf("Request type is nil")
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
					smMessage, pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
				return nil
			}
			switch requestType.GetRequestTypeValue() {
			// case iii) if the AMF does not have a PDU session routing context for the PDU session ID and the UE
			// and the Request type IE is included and is set to "initial request"
			case nasMessage.ULNASTransportRequestTypeInitialRequest:
				_, err := CreatePDUSession(ulNasTransport, ue, anType, pduSessionID, smMessage)
				return err
			case nasMessage.ULNASTransportRequestTypeModificationRequest:
				fallthrough
			case nasMessage.ULNASTransportRequestTypeExistingPduSession:
				if ue.UeContextInSmfData != nil {
					// TS 24.501 5.4.5.2.5 case a) 3)
					pduSessionIDStr := fmt.Sprintf("%d", pduSessionID)
					if ueContextInSmf, ok := ue.UeContextInSmfData.PduSessions[pduSessionIDStr]; !ok {
						gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
							smMessage, pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
					} else {
						// TS 24.501 5.4.5.2.3 case a) 1) iv)
						smContext = context.NewSmContext(pduSessionID)
						smContext.SetAccessType(anType)
						smContext.SetSmfID(ueContextInSmf.SmfInstanceId)
						smContext.SetDnn(ueContextInSmf.Dnn)
						smContext.SetPlmnID(*ueContextInSmf.PlmnId)
						ue.StoreSmContext(pduSessionID, smContext)
						return forward5GSMMessageToSMF(ue, anType, pduSessionID, smContext, smMessage)
					}
				} else {
					gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
						smMessage, pduSessionID, nasMessage.Cause5GMMPayloadWasNotForwarded, nil, 0)
				}
			default:
			}
		}
	} else {
		// TODO: implement SSC mode3 Op
		return fmt.Errorf("SSC mode3 operation has not been implemented yet")
	}
	return nil
}

func CreatePDUSession(ulNasTransport *nasMessage.ULNASTransport,
	ue *context.AmfUe,
	anType models.AccessType,
	pduSessionID int32,
	smMessage []uint8,
) (setNewSmContext bool, err error) {
	var (
		snssai models.Snssai
		dnn    string
	)
	// A) AMF shall select an SMF

	// If the S-NSSAI IE is not included and the user's subscription context obtained from UDM. AMF shall
	// select a default snssai
	if ulNasTransport.SNSSAI != nil {
		snssai = nasConvert.SnssaiToModels(ulNasTransport.SNSSAI)
	} else {
		if allowedNssai, ok := ue.AllowedNssai[anType]; ok {
			snssai = *allowedNssai[0].AllowedSnssai
		} else {
			return false, errors.New("Ue doesn't have allowedNssai")
		}
	}

	if ulNasTransport.DNN != nil {
		dnn = ulNasTransport.DNN.GetDNN()
	} else {
		// if user's subscription context obtained from UDM does not contain the default DNN for the,
		// S-NSSAI, the AMF shall use a locally configured DNN as the DNN
		dnn = ue.ServingAMF().SupportDnnLists[0]

		if ue.SmfSelectionData != nil {
			snssaiStr := util.SnssaiModelsToHex(snssai)
			if snssaiInfo, ok := ue.SmfSelectionData.SubscribedSnssaiInfos[snssaiStr]; ok {
				for _, dnnInfo := range snssaiInfo.DnnInfos {
					if dnnInfo.DefaultDnnIndicator {
						dnn = dnnInfo.Dnn.(string)
					}
				}
			}
		}
	}

	if newSmContext, cause, errSelectSmf := consumer.GetConsumer().SelectSmf(
		ue, anType, pduSessionID, snssai, dnn); errSelectSmf != nil {
		ue.GmmLog.Errorf("Select SMF failed: %+v", errSelectSmf)
		gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
			smMessage, pduSessionID, cause, nil, 0)
	} else {
		ue.Lock.Lock()
		defer ue.Lock.Unlock()

		smContextRef, errResponse, problemDetail, errSendReq := consumer.GetConsumer().SendCreateSmContextRequest(
			ue, newSmContext, nil, smMessage)
		if errSendReq != nil {
			ue.GmmLog.Errorf("CreateSmContextRequest Error: %+v", errSendReq)
			return false, nil
		} else if problemDetail != nil {
			// TODO: error handling
			return false, fmt.Errorf("failed to Create smContext[pduSessionID: %d], Error[%v]", pduSessionID, problemDetail)
		} else if errResponse != nil {
			ue.GmmLog.Warnf("PDU Session Establishment Request is rejected by SMF[pduSessionId:%d]",
				pduSessionID)
			gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
				errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
		} else {
			newSmContext.SetSmContextRef(smContextRef)
			newSmContext.SetUserLocation(deepcopy.Copy(ue.Location).(models.UserLocation))
			ue.StoreSmContext(pduSessionID, newSmContext)
			ue.GmmLog.Infof("create smContext[pduSessionID: %d] Success", pduSessionID)
			// TODO: handle response(response N2SmInfo to RAN if exists)
			return true, nil
		}
	}
	return false, nil
}

func forward5GSMMessageToSMF(
	ue *context.AmfUe,
	accessType models.AccessType,
	pduSessionID int32,
	smContext *context.SmContext,
	smMessage []byte,
) error {
	smContextUpdateData := models.SmfPduSessionSmContextUpdateData{
		N1SmMsg: &models.RefToBinaryData{
			ContentId: "N1SmMsg",
		},
	}
	smContextUpdateData.Pei = ue.Pei
	if !context.CompareUserLocation(ue.Location, smContext.UserLocation()) {
		smContextUpdateData.UeLocation = &ue.Location
	}

	if accessType != smContext.AccessType() {
		smContextUpdateData.AnType = accessType
	}

	response, errResponse, problemDetail, err := consumer.GetConsumer().SendUpdateSmContextRequest(smContext,
		&smContextUpdateData, smMessage, nil)

	if err != nil {
		// TODO: error handling
		ue.GmmLog.Errorf("Update SMContext error [pduSessionID: %d], Error[%v]", pduSessionID, err)
		return nil
	} else if problemDetail != nil {
		ue.GmmLog.Errorf("Update SMContext failed [pduSessionID: %d], problem[%v]", pduSessionID, problemDetail)
		return nil
	} else if errResponse != nil {
		errJSON := errResponse.JsonData
		n1Msg := errResponse.BinaryDataN1SmMessage
		ue.GmmLog.Warnf("PDU Session Modification Procedure is rejected by SMF[pduSessionId:%d], Error[%s]",
			pduSessionID, errJSON.Error.Cause)
		if n1Msg != nil {
			gmm_message.SendDLNASTransport(ue.RanUe[accessType], nasMessage.PayloadContainerTypeN1SMInfo,
				errResponse.BinaryDataN1SmMessage, pduSessionID, 0, nil, 0)
		}
		// TODO: handle n2 info transfer
	} else if response != nil {
		// update SmContext in AMF
		smContext.SetAccessType(accessType)
		smContext.SetUserLocation(ue.Location)

		responseData := response.JsonData
		var n1Msg []byte
		n2SmInfo := response.BinaryDataN2SmInformation
		if response.BinaryDataN1SmMessage != nil {
			ue.GmmLog.Debug("Receive N1 SM Message from SMF")
			n1Msg, err = gmm_message.BuildDLNASTransport(ue, accessType, nasMessage.PayloadContainerTypeN1SMInfo,
				response.BinaryDataN1SmMessage, uint8(pduSessionID), nil, nil, 0)
			if err != nil {
				return err
			}
		}

		if response.BinaryDataN2SmInformation != nil {
			ue.GmmLog.Debugf("Receive N2 SM Information[%s] from SMF", responseData.N2SmInfoType)
			switch responseData.N2SmInfoType {
			case models.N2SmInfoType_PDU_RES_MOD_REQ:
				list := ngapType.PDUSessionResourceModifyListModReq{}
				ngap_message.AppendPDUSessionResourceModifyListModReq(&list, pduSessionID, n1Msg, n2SmInfo)
				ngap_message.SendPDUSessionResourceModifyRequest(ue.RanUe[accessType], list)
			case models.N2SmInfoType_PDU_RES_REL_CMD:
				list := ngapType.PDUSessionResourceToReleaseListRelCmd{}
				ngap_message.AppendPDUSessionResourceToReleaseListRelCmd(&list, pduSessionID, n2SmInfo)
				ngap_message.SendPDUSessionResourceReleaseCommand(ue.RanUe[accessType], n1Msg, list)
			default:
				return fmt.Errorf("error N2 SM information type[%s]", responseData.N2SmInfoType)
			}
		} else if n1Msg != nil {
			ue.GmmLog.Debugf("AMF forward Only N1 SM Message to UE")
			ngap_message.SendDownlinkNasTransport(ue.RanUe[accessType], n1Msg, nil)
		}
	}
	return nil
}

// Handle cleartext IEs of Registration Request, which cleattext IEs defined in TS 24.501 4.4.6
func HandleRegistrationRequest(ue *context.AmfUe, anType models.AccessType, procedureCode int64,
	registrationRequest *nasMessage.RegistrationRequest,
) error {
	var guamiFromUeGuti models.Guami
	amfSelf := context.GetSelf()

	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	ue.GmmLog.Info("Handle Registration Request")

	if ue.RanUe[anType] == nil {
		return fmt.Errorf("RanUe is nil")
	}

	ue.SetOnGoing(anType, &context.OnGoing{
		Procedure: context.OnGoingProcedureRegistration,
	})

	ue.StopT3513()
	ue.StopT3565()

	// TS 24.501 8.2.6.21: if the UE is sending a REGISTRATION REQUEST message as an initial NAS message,
	// the UE has a valid 5G NAS security context and the UE needs to send non-cleartext IEs
	// TS 24.501 4.4.6: When the UE sends a REGISTRATION REQUEST or SERVICE REQUEST message that includes a NAS message
	// container IE, the UE shall set the security header type of the initial NAS message to "integrity protected"
	if registrationRequest.NASMessageContainer != nil && !ue.MacFailed {
		contents := registrationRequest.NASMessageContainer.GetNASMessageContainerContents()

		// TS 24.501 4.4.6: When the UE sends a REGISTRATION REQUEST or SERVICE REQUEST message that includes a NAS
		// message container IE, the UE shall set the security header type of the initial NAS message to
		// "integrity protected"; then the AMF shall decipher the value part of the NAS message container IE
		err := security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.ULCount.Get(), security.Bearer3GPP,
			security.DirectionUplink, contents)
		if err != nil {
			ue.SecurityContextAvailable = false
		} else {
			m := nas.NewMessage()
			if errGmmMessageDecode := m.GmmMessageDecode(&contents); errGmmMessageDecode != nil {
				return errGmmMessageDecode
			}

			messageType := m.GmmMessage.GmmHeader.GetMessageType()
			if messageType != nas.MsgTypeRegistrationRequest {
				return errors.New("The payload of NAS Message Container is not Registration Request")
			}
			// TS 24.501 4.4.6: The AMF shall consider the NAS message that is obtained from the NAS message container
			// IE as the initial NAS message that triggered the procedure
			registrationRequest = m.RegistrationRequest
		}
	}
	// TS 33.501 6.4.6 step 3: if the initial NAS message was protected but did not pass the integrity check
	ue.RetransmissionOfInitialNASMsg = ue.MacFailed

	ue.RegistrationRequest = registrationRequest
	ue.RegistrationType5GS = registrationRequest.NgksiAndRegistrationType5GS.GetRegistrationType5GS()
	switch ue.RegistrationType5GS {
	case nasMessage.RegistrationType5GSInitialRegistration:
		ue.GmmLog.Infof("RegistrationType: Initial Registration")
		ue.SecurityContextAvailable = false // need to start authentication procedure later
	case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
		ue.GmmLog.Infof("RegistrationType: Mobility Registration Updating")
		if ue.State[anType].Is(context.Deregistered) {
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMImplicitlyDeregistered, "")
			return fmt.Errorf("mobility registration updating was sent when the UE state was deregistered")
		}
	case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
		ue.GmmLog.Infof("RegistrationType: Periodic Registration Updating")
		if ue.State[anType].Is(context.Deregistered) {
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMImplicitlyDeregistered, "")
			return fmt.Errorf("periodic registration updating was sent when the UE state was deregistered")
		}
	case nasMessage.RegistrationType5GSEmergencyRegistration:
		return fmt.Errorf("not supported RegistrationType: Emergency Registration")
	case nasMessage.RegistrationType5GSReserved:
		ue.RegistrationType5GS = nasMessage.RegistrationType5GSInitialRegistration
		ue.GmmLog.Infof("RegistrationType: Reserved")
	default:
		ue.GmmLog.Infof("RegistrationType: %v, chage state to InitialRegistration", ue.RegistrationType5GS)
		ue.RegistrationType5GS = nasMessage.RegistrationType5GSInitialRegistration
	}

	mobileIdentity5GSContents := registrationRequest.MobileIdentity5GS.GetMobileIdentity5GSContents()
	if len(mobileIdentity5GSContents) < 1 {
		return errors.New("broken MobileIdentity5GS")
	}
	ue.IdentityTypeUsedForRegistration = nasConvert.GetTypeOfIdentity(mobileIdentity5GSContents[0])
	switch ue.IdentityTypeUsedForRegistration { // get type of identity
	case nasMessage.MobileIdentity5GSTypeNoIdentity:
		ue.GmmLog.Infof("MobileIdentity5GS: No Identity")
	case nasMessage.MobileIdentity5GSTypeSuci:
		if suci, plmnId, err := nasConvert.SuciToStringWithError(mobileIdentity5GSContents); err != nil {
			return fmt.Errorf("decode SUCI failed: %w", err)
		} else if plmnId == "" {
			return errors.New("empty plmnId")
		} else {
			ue.Suci = suci
			ue.PlmnId = util.PlmnIdStringToModels(plmnId)
		}
		ue.GmmLog.Infof("MobileIdentity5GS: SUCI[%s]", ue.Suci)
	case nasMessage.MobileIdentity5GSType5gGuti:
		guamiFromUeGutiTmp, guti, err := nasConvert.GutiToStringWithError(mobileIdentity5GSContents)
		if err != nil {
			return fmt.Errorf("decode GUTI failed: %w", err)
		}
		guamiFromUeGuti = guamiFromUeGutiTmp
		ue.PlmnId = util.PlmnIdNidToModelsPlmnId(*guamiFromUeGuti.PlmnId)
		ue.GmmLog.Infof("MobileIdentity5GS: GUTI[%s]", guti)

		// TODO: support multiple ServedGuami
		servedGuami := amfSelf.ServedGuamiList[0]
		if reflect.DeepEqual(guamiFromUeGuti, servedGuami) {
			ue.ServingAmfChanged = false
			// refresh 5G-GUTI according to 6.12.3 Subscription temporary identifier, TS33.501
			if ue.SecurityContextAvailable {
				context.GetSelf().FreeTmsi(int64(ue.Tmsi))
				context.GetSelf().AllocateGutiToUe(ue)
			}
		} else {
			ue.GmmLog.Infof("Serving AMF has changed: guamiFromUeGuti[%+v], servedGuami[%+v]",
				guamiFromUeGuti, servedGuami)
			ue.ServingAmfChanged = true
			context.GetSelf().FreeTmsi(int64(ue.Tmsi))
			ue.Guti = guti
		}
	case nasMessage.MobileIdentity5GSTypeImei:
		imei, err := nasConvert.PeiToStringWithError(mobileIdentity5GSContents)
		if err != nil {
			return fmt.Errorf("decode PEI failed: %w", err)
		}
		ue.Pei = imei
		ue.GmmLog.Infof("MobileIdentity5GS: PEI[%s]", imei)
	case nasMessage.MobileIdentity5GSTypeImeisv:
		imeisv, err := nasConvert.PeiToStringWithError(mobileIdentity5GSContents)
		if err != nil {
			return fmt.Errorf("decode PEI failed: %w", err)
		}
		ue.Pei = imeisv
		ue.GmmLog.Infof("MobileIdentity5GS: PEI[%s]", imeisv)
	}

	// NgKsi: TS 24.501 9.11.3.32
	switch registrationRequest.NgksiAndRegistrationType5GS.GetTSC() {
	case nasMessage.TypeOfSecurityContextFlagNative:
		ue.NgKsi.Tsc = models.ScType_NATIVE
	case nasMessage.TypeOfSecurityContextFlagMapped:
		ue.NgKsi.Tsc = models.ScType_MAPPED
	}
	ue.NgKsi.Ksi = int32(registrationRequest.NgksiAndRegistrationType5GS.GetNasKeySetIdentifiler())
	if ue.NgKsi.Tsc == models.ScType_NATIVE && ue.NgKsi.Ksi != 7 {
	} else {
		ue.NgKsi.Tsc = models.ScType_NATIVE
		ue.NgKsi.Ksi = 0
	}

	// Copy UserLocation from ranUe
	// TODO: This check due to RanUe may release during the process;it should be a better way to make this procedure
	// as an atomic operation
	if ue.RanUe[anType] != nil {
		ue.Location = ue.RanUe[anType].Location
		ue.Tai = ue.RanUe[anType].Tai
		if ue.RanUe[anType].Ran != nil {
			// ue.Ratype TS 23.502 4.2.2.1
			// The AMF determines Access Type and RAT Type as defined in clause 5.3.2.3 of TS 23.501 .
			ue.RatType = ue.RanUe[anType].Ran.UeRatType()
		}
	}

	// Check TAI
	if !context.InTaiList(ue.Tai, amfSelf.SupportTaiLists) {
		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMTrackingAreaNotAllowed, "")
		return fmt.Errorf("registration reject[tracking area not allowed]")
	}

	if registrationRequest.UESecurityCapability != nil {
		ue.UESecurityCapability = *registrationRequest.UESecurityCapability
	} else if registrationRequest.GetRegistrationType5GS() != nasMessage.RegistrationType5GSPeriodicRegistrationUpdating {
		// TS 23.501 8.2.6.4
		// The UE shall include this IE, unless the UE performs a periodic registration updating procedure.
		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
		return fmt.Errorf("UESecurityCapability is nil")
	}

	// TODO (TS 23.502 4.2.2.2 step 4): if UE's 5g-GUTI is included & serving AMF has changed
	// since last registration procedure, new AMF may invoke Namf_Communication_UEContextTransfer
	// to old AMF, including the complete registration request nas msg, to request UE's SUPI & UE Context
	if ue.ServingAmfChanged {
		if err := contextTransferFromOldAmf(ue, anType, guamiFromUeGuti); err != nil {
			ue.GmmLog.Warnf("[GMM] %+v", err)
			// if failed, give up to retrieve the old context and start a new authentication procedure.
			ue.ServingAmfChanged = false
			context.GetSelf().AllocateGutiToUe(ue) // refresh 5G-GUTI
		}
	}

	return nil
}

func contextTransferFromOldAmf(ue *context.AmfUe, anType models.AccessType, oldAmfGuami models.Guami) error {
	ue.GmmLog.Infof("ContextTransfer from old AMF[%s %s]", oldAmfGuami.PlmnId, oldAmfGuami.AmfId)

	amfSelf := context.GetSelf()
	searchOpt := Nnrf_NFDiscovery.SearchNFInstancesRequest{
		Guami: &oldAmfGuami,
	}
	if err := consumer.GetConsumer().SearchAmfCommunicationInstance(ue, amfSelf.NrfUri, models.NrfNfManagementNfType_AMF,
		models.NrfNfManagementNfType_AMF, &searchOpt); err != nil {
		return err
	}

	var transferReason models.TransferReason
	switch ue.RegistrationType5GS {
	case nasMessage.RegistrationType5GSInitialRegistration:
		transferReason = models.TransferReason_INIT_REG
	case nasMessage.RegistrationType5GSMobilityRegistrationUpdating:
		fallthrough
	case nasMessage.RegistrationType5GSPeriodicRegistrationUpdating:
		transferReason = models.TransferReason_MOBI_REG
	}

	ueContextTransferRspData, pd, err := consumer.GetConsumer().UEContextTransferRequest(ue, anType, transferReason)
	if pd != nil {
		if pd.Cause == "INTEGRITY_CHECK_FAIL" || pd.Cause == "CONTEXT_NOT_FOUND" {
			// TODO 9a. After successful authentication in new AMF, which is triggered by the integrity check failure
			// in old AMF at step 5, the new AMF invokes step 4 above again and indicates that the UE is validated
			// (i.e. through the reason parameter as specified in clause 5.2.2.2.2).
			return fmt.Errorf("can not retrieve UE context from old AMF[cause: %s]", pd.Cause)
		}
		return fmt.Errorf("UE Context Transfer Request Failed Problem[%+v]", pd)
	} else if err != nil {
		return fmt.Errorf("UE Context Transfer Request Error[%+v]", err)
	} else {
		ue.SecurityContextAvailable = true
		ue.MacFailed = false
	}

	ue.CopyDataFromUeContextModel(ueContextTransferRspData.UeContext)
	if ue.SecurityContextAvailable {
		ue.DerivateAlgKey()
	}
	return nil
}

func IdentityVerification(ue *context.AmfUe) bool {
	return ue.Supi != "" || len(ue.Suci) != 0
}

func HandleInitialRegistration(ue *context.AmfUe, anType models.AccessType) error {
	ue.GmmLog.Infoln("Handle InitialRegistration")

	amfSelf := context.GetSelf()

	// update Kgnb/Kn3iwf
	ue.UpdateSecurityContext(anType)

	// Registration with AMF re-allocation (TS 23.502 4.2.2.2.3)
	if len(ue.SubscribedNssai) == 0 {
		getSubscribedNssai(ue)
	}

	if err := handleRequestedNssai(ue, anType); err != nil {
		return err
	}

	if ue.RegistrationRequest.Capability5GMM != nil {
		ue.Capability5GMM = *ue.RegistrationRequest.Capability5GMM
	} else {
		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
		return fmt.Errorf("Capability5GMM is nil")
	}

	storeLastVisitedRegisteredTAI(ue, ue.RegistrationRequest.LastVisitedRegisteredTAI)

	if ue.RegistrationRequest.MICOIndication != nil {
		ue.GmmLog.Warnf("Receive MICO Indication[RAAI: %d], Not Supported",
			ue.RegistrationRequest.MICOIndication.GetRAAI())
	}

	// TODO: Negotiate DRX value if need (TS 23.501 5.4.5)
	negotiateDRXParameters(ue, ue.RegistrationRequest.RequestedDRXParameters)

	// TODO (step 10 optional): send Namf_Communication_RegistrationCompleteNotify to old AMF if need
	if ue.ServingAmfChanged {
		// If the AMF has changed the new AMF notifies the old AMF that the registration of the UE in the new AMF is completed
		req := models.UeRegStatusUpdateReqData{
			TransferStatus: models.UeContextTransferStatus_TRANSFERRED,
		}
		// TODO: based on locol policy, decide if need to change serving PCF for UE
		regStatusTransferComplete, problemDetails, err := consumer.GetConsumer().RegistrationStatusUpdate(ue, req)
		if problemDetails != nil {
			ue.GmmLog.Errorf("Registration Status Update Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			ue.GmmLog.Errorf("Registration Status Update Error[%+v]", err)
		} else if regStatusTransferComplete {
			ue.GmmLog.Infof("Registration Status Transfer complete")
		}
	}

	if len(ue.Pei) == 0 {
		gmm_message.SendIdentityRequest(ue.RanUe[anType], anType, nasMessage.MobileIdentity5GSTypeImei)
		return nil
	}

	// TODO (step 12 optional): the new AMF initiates ME identity check by invoking the
	// N5g-eir_EquipmentIdentityCheck_Get service operation

	if ue.ServingAmfChanged || ue.State[models.AccessType_NON_3_GPP_ACCESS].Is(context.Registered) ||
		!ue.ContextValid {
		if err := communicateWithUDM(ue, anType); err != nil {
			ue.GmmLog.Errorf("communicateWithUDM error: %v", err)
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMPLMNNotAllowed, "")
			return errors.Wrap(err, "communicateWithUDM failed")
		}
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesRequest{
		Supi: &ue.Supi,
	}
	if amfSelf.Locality != "" {
		param.PreferredLocality = &amfSelf.Locality
	}

	// TODO: (step 15) Should use PCF ID to select PCF
	// Retrieve PCF ID from old AMF
	// if ue.PcfId != "" {

	// }
	for {
		resp, err := consumer.GetConsumer().SendSearchNFInstances(
			amfSelf.NrfUri, models.NrfNfManagementNfType_PCF, models.NrfNfManagementNfType_AMF, &param)
		if err != nil {
			ue.GmmLog.Error("AMF can not select an PCF by NRF")
		} else {
			// select the first PCF, TODO: select base on other info
			var pcfUri string
			for index := range resp.NfInstances {
				pcfUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NPCF_AM_POLICY_CONTROL,
					models.NfServiceStatus_REGISTERED)
				if pcfUri != "" {
					ue.PcfId = resp.NfInstances[index].NfInstanceId
					break
				}
			}
			if ue.PcfUri = pcfUri; ue.PcfUri == "" {
				ue.GmmLog.Error("AMF can not select an PCF by NRF")
			} else {
				break
			}
		}
		time.Sleep(500 * time.Millisecond) // sleep a while when search NF Instance fail
	}

	problemDetails, err := consumer.GetConsumer().AMPolicyControlCreate(ue, anType)
	if problemDetails != nil {
		ue.GmmLog.Errorf("AM Policy Control Create Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		ue.GmmLog.Errorf("AM Policy Control Create Error[%+v]", err)
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
				// if numOfallowedTAs < int(servAreaRes.MaxNumOfTAs) {
				// 	TODO: based on AMF Policy, assign additional allowed area for UE,
				// 	and the upper limit is servAreaRes.MaxNumOfTAs (TS 29.507 4.2.2.3)
				// }
			}
		}
	}

	// TODO (step 18 optional):
	// If the AMF has changed and the old AMF has indicated an existing NGAP UE association towards a N3IWF, the new AMF
	// creates an NGAP UE association towards the N3IWF to which the UE is connectedsend N2 AMF mobility request to N3IWF
	// if anType == models.AccessType_NON_3_GPP_ACCESS && ue.ServingAmfChanged {
	// 	TODO: send N2 AMF Mobility Request
	// }

	amfSelf.AllocateRegistrationArea(ue, anType)
	ue.GmmLog.Debugf("Use original GUTI[%s]", ue.Guti)

	assignLadnInfo(ue, anType)

	amfSelf.AddAmfUeToUePool(ue, ue.Supi)
	ue.T3502Value = amfSelf.T3502Value
	if anType == models.AccessType__3_GPP_ACCESS {
		ue.T3512Value = amfSelf.T3512Value
	} else {
		ue.Non3gppDeregTimerValue = amfSelf.Non3gppDeregTimerValue
	}

	gmm_message.SendRegistrationAccept(ue, anType, nil, nil, nil, nil, nil)
	return nil
}

func HandleMobilityAndPeriodicRegistrationUpdating(ue *context.AmfUe, anType models.AccessType) error {
	ue.GmmLog.Infoln("Handle MobilityAndPeriodicRegistrationUpdating")

	amfSelf := context.GetSelf()

	if ue.RegistrationRequest.UpdateType5GS != nil {
		if ue.RegistrationRequest.UpdateType5GS.GetNGRanRcu() == nasMessage.NGRanRadioCapabilityUpdateNeeded {
			ue.UeRadioCapability = ""
			ue.UeRadioCapabilityForPaging = nil
		}
	}

	// Registration with AMF re-allocation (TS 23.502 4.2.2.2.3)
	if len(ue.SubscribedNssai) == 0 {
		getSubscribedNssai(ue)
	}

	if err := handleRequestedNssai(ue, anType); err != nil {
		return err
	}

	if ue.RegistrationRequest.Capability5GMM != nil {
		ue.Capability5GMM = *ue.RegistrationRequest.Capability5GMM
	} else if ue.RegistrationType5GS != nasMessage.RegistrationType5GSPeriodicRegistrationUpdating {
		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
		return fmt.Errorf("Capability5GMM is nil")
	}

	storeLastVisitedRegisteredTAI(ue, ue.RegistrationRequest.LastVisitedRegisteredTAI)

	if ue.RegistrationRequest.MICOIndication != nil {
		ue.GmmLog.Warnf("Receive MICO Indication[RAAI: %d], Not Supported",
			ue.RegistrationRequest.MICOIndication.GetRAAI())
	}

	// TODO: Negotiate DRX value if need (TS 23.501 5.4.5)
	negotiateDRXParameters(ue, ue.RegistrationRequest.RequestedDRXParameters)

	// TODO (step 10 optional): send Namf_Communication_RegistrationCompleteNotify to old AMF if need
	// if ue.ServingAmfChanged {
	// 	If the AMF has changed the new AMF notifies the old AMF that the registration of the UE in the new AMF is completed
	// }

	if len(ue.Pei) == 0 {
		gmm_message.SendIdentityRequest(ue.RanUe[anType], anType, nasMessage.MobileIdentity5GSTypeImei)
		return nil
	}

	// TODO (step 12 optional): the new AMF initiates ME identity check by invoking the
	// N5g-eir_EquipmentIdentityCheck_Get service operation

	if ue.ServingAmfChanged || ue.State[models.AccessType_NON_3_GPP_ACCESS].Is(context.Registered) ||
		!ue.ContextValid {
		if err := communicateWithUDM(ue, anType); err != nil {
			ue.GmmLog.Errorf("communicateWithUDM error: %v", err)
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMPLMNNotAllowed, "")
			return errors.Wrap(err, "communicateWithUDM failed")
		}
	}

	var reactivationResult *[psiArraySize]bool
	var errPduSessionId, errCause []uint8
	cxtList := ngapType.PDUSessionResourceSetupListCxtReq{}

	if ue.RegistrationRequest.UplinkDataStatus != nil {
		uplinkDataPsi := nasConvert.PSIToBooleanArray(ue.RegistrationRequest.UplinkDataStatus.Buffer)
		reactivationResult = new([psiArraySize]bool)
		allowReEstablishPduSession := true

		// determines that the UE is in non-allowed area or is not in allowed area
		if ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
			switch ue.AmPolicyAssociation.ServAreaRes.RestrictionType {
			case models.RestrictionType_ALLOWED_AREAS:
				allowReEstablishPduSession = context.TacInAreas(ue.Tai.Tac, ue.AmPolicyAssociation.ServAreaRes.Areas)
			case models.RestrictionType_NOT_ALLOWED_AREAS:
				allowReEstablishPduSession = !context.TacInAreas(ue.Tai.Tac, ue.AmPolicyAssociation.ServAreaRes.Areas)
			}
		}

		if !allowReEstablishPduSession {
			for pduSessionId, hasUplinkData := range uplinkDataPsi {
				if hasUplinkData {
					errPduSessionId = append(errPduSessionId, uint8(pduSessionId))
					errCause = append(errCause, nasMessage.Cause5GMMRestrictedServiceArea)
				}
			}
		} else {
			// There is no serviceType in MobilityAndPeriodicRegistrationUpdating
			errPduSessionId, errCause = reactivatePendingULDataPDUSession(ue, anType, 0, &uplinkDataPsi, 0, &cxtList,
				reactivationResult, errPduSessionId, errCause)
		}
	}

	var pduSessionStatus *[psiArraySize]bool
	if ue.RegistrationRequest.PDUSessionStatus != nil {
		pduSessionStatus = new([psiArraySize]bool)
		pduSessionPsi := nasConvert.PSIToBooleanArray(ue.RegistrationRequest.PDUSessionStatus.Buffer)
		releaseInactivePDUSession(ue, anType, &pduSessionPsi, pduSessionStatus)
	}

	// AllowedPDUSessionStatus indicate to the network PDU sessions associated with non-3GPP access that
	// are allowed to be re-established over 3GPP access
	if ue.RegistrationRequest.AllowedPDUSessionStatus != nil &&
		anType == models.AccessType__3_GPP_ACCESS && ue.N1N2Message != nil {
		allowedPsi := nasConvert.PSIToBooleanArray(ue.RegistrationRequest.AllowedPDUSessionStatus.Buffer)
		requestData := ue.N1N2Message.Request.JsonData
		n1Msg := ue.N1N2Message.Request.BinaryDataN1Message
		n2Info := ue.N1N2Message.Request.BinaryDataN2Information

		if n2Info == nil {
			// SMF has indicated pending downlink signaling only,
			// forward the received 5GSM message via 3GPP access to the UE
			// after the REGISTRATION ACCEPT message is sent
			gmm_message.SendRegistrationAccept(ue, anType, pduSessionStatus,
				reactivationResult, errPduSessionId, errCause, &cxtList)

			switch requestData.N1MessageContainer.N1MessageClass {
			case models.N1MessageClass_SM:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeN1SMInfo,
					n1Msg, requestData.PduSessionId, 0, nil, 0)
			case models.N1MessageClass_LPP:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeLPP,
					n1Msg, 0, 0, nil, 0)
			case models.N1MessageClass_SMS:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeSMS,
					n1Msg, 0, 0, nil, 0)
			case models.N1MessageClass_UPDP:
				gmm_message.SendDLNASTransport(ue.RanUe[anType], nasMessage.PayloadContainerTypeUEPolicy,
					n1Msg, 0, 0, nil, 0)
			}
			ue.N1N2Message = nil
			return nil
		}

		// SMF has indicated pending downlink data
		// notify the SMF that reactivation of the user-plane resources for the corresponding PDU session(s)
		// associated with non-3GPP access
		smContext, exist := ue.SmContextFindByPDUSessionID(requestData.PduSessionId)
		if !exist {
			ue.N1N2Message = nil
			return fmt.Errorf("SmContext[PDU Session ID:%d] not found", requestData.PduSessionId)
		}

		errPduSessionId, errCause = reestablishAllowedPDUSessionOver3GPP(ue, anType, smContext, &allowedPsi, &cxtList,
			reactivationResult, errPduSessionId, errCause)
	}

	if ue.LocationChanged && ue.RequestTriggerLocationChange {
		updateReq := models.PcfAmPolicyControlPolicyAssociationUpdateRequest{}
		updateReq.Triggers = append(updateReq.Triggers, models.PcfAmPolicyControlRequestTrigger_LOC_CH)
		updateReq.UserLoc = &ue.Location
		problemDetails, err := consumer.GetConsumer().AMPolicyControlUpdate(ue, updateReq)
		if problemDetails != nil {
			ue.GmmLog.Errorf("AM Policy Control Update Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			ue.GmmLog.Errorf("AM Policy Control Update Error[%v]", err)
		}
		ue.LocationChanged = false
	}

	// TODO (step 18 optional):
	// If the AMF has changed and the old AMF has indicated an existing NGAP UE association towards a N3IWF, the new AMF
	// creates an NGAP UE association towards the N3IWF to which the UE is connectedsend N2 AMF mobility request to N3IWF
	// if anType == models.AccessType_NON_3_GPP_ACCESS && ue.ServingAmfChanged {
	// 	TODO: send N2 AMF Mobility Request
	// }

	amfSelf.AllocateRegistrationArea(ue, anType)
	assignLadnInfo(ue, anType)

	// TODO: GUTI reassignment if need (based on operator poilcy)
	// TODO: T3512/Non3GPP de-registration timer reassignment if need (based on operator policy)

	gmm_message.SendRegistrationAccept(ue, anType, pduSessionStatus, reactivationResult,
		errPduSessionId, errCause, &cxtList)
	return nil
}

// TS 23.502 4.2.2.2.2 step 1
// If available, the last visited TAI shall be included in order to help the AMF produce Registration Area for the UE
func storeLastVisitedRegisteredTAI(ue *context.AmfUe, lastVisitedRegisteredTAI *nasType.LastVisitedRegisteredTAI) {
	if lastVisitedRegisteredTAI != nil {
		plmnID := nasConvert.PlmnIDToString(lastVisitedRegisteredTAI.Octet[1:4])
		nasTac := lastVisitedRegisteredTAI.GetTAC()
		tac := hex.EncodeToString(nasTac[:])

		tai := models.Tai{
			PlmnId: &models.PlmnId{
				Mcc: plmnID[:3],
				Mnc: plmnID[3:],
			},
			Tac: tac,
		}

		ue.LastVisitedRegisteredTai = tai
		ue.GmmLog.Debugf("Ue Last Visited Registered Tai; %v", ue.LastVisitedRegisteredTai)
	}
}

func negotiateDRXParameters(ue *context.AmfUe, requestedDRXParameters *nasType.RequestedDRXParameters) {
	if requestedDRXParameters != nil {
		switch requestedDRXParameters.GetDRXValue() {
		case nasMessage.DRXcycleParameterT32:
			ue.GmmLog.Tracef("Requested DRX: T = 32")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT32
		case nasMessage.DRXcycleParameterT64:
			ue.GmmLog.Tracef("Requested DRX: T = 64")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT64
		case nasMessage.DRXcycleParameterT128:
			ue.GmmLog.Tracef("Requested DRX: T = 128")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT128
		case nasMessage.DRXcycleParameterT256:
			ue.GmmLog.Tracef("Requested DRX: T = 256")
			ue.UESpecificDRX = nasMessage.DRXcycleParameterT256
		case nasMessage.DRXValueNotSpecified:
			fallthrough
		default:
			ue.UESpecificDRX = nasMessage.DRXValueNotSpecified
			ue.GmmLog.Tracef("Requested DRX: Value not specified")
		}
	}
}

func communicateWithUDM(ue *context.AmfUe, accessType models.AccessType) error {
	ue.GmmLog.Debugln("communicateWithUDM")
	amfSelf := context.GetSelf()

	// UDM selection described in TS 23.501 6.3.8
	// TODO: consider udm group id, Routing ID part of SUCI, GPSI or External Group ID (e.g., by the NEF)
	param := Nnrf_NFDiscovery.SearchNFInstancesRequest{
		Supi: &ue.Supi,
	}
	resp, err := consumer.GetConsumer().SendSearchNFInstances(
		amfSelf.NrfUri, models.NrfNfManagementNfType_UDM, models.NrfNfManagementNfType_AMF, &param)
	if err != nil {
		return errors.Errorf("AMF can not select an UDM by NRF: SendSearchNFInstances failed")
	}

	var uecmUri, sdmUri string
	for index := range resp.NfInstances {
		ue.UdmId = resp.NfInstances[index].NfInstanceId
		uecmUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NUDM_UECM,
			models.NfServiceStatus_REGISTERED)
		sdmUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NUDM_SDM,
			models.NfServiceStatus_REGISTERED)
		if uecmUri != "" && sdmUri != "" {
			break
		}
	}
	ue.NudmUECMUri = uecmUri
	ue.NudmSDMUri = sdmUri
	if ue.NudmUECMUri == "" || ue.NudmSDMUri == "" {
		return errors.Errorf("AMF can not select an UDM by NRF: SearchNFServiceUri failed")
	}

	problemDetails, err := consumer.GetConsumer().UeCmRegistration(ue, accessType, true)
	if problemDetails != nil {
		return errors.Errorf("%s", problemDetails.Cause)
	} else if err != nil {
		return errors.Wrap(err, "UECM_Registration Error")
	}

	// TS 23.502 4.2.2.2.1 14a-c.
	// "After a successful response is received, the AMF subscribes to be notified
	// 		using Nudm_SDM_Subscribe when the data requested is modified"
	problemDetails, err = consumer.GetConsumer().SDMGetAmData(ue)
	if problemDetails != nil {
		return errors.Errorf("%s", problemDetails.Cause)
	} else if err != nil {
		return errors.Wrap(err, "SDM_Get AmData Error")
	}

	problemDetails, err = consumer.GetConsumer().SDMGetSmfSelectData(ue)
	if problemDetails != nil {
		return errors.Errorf("%s", problemDetails.Cause)
	} else if err != nil {
		return errors.Wrap(err, "SDM_Get SmfSelectData Error")
	}

	problemDetails, err = consumer.GetConsumer().SDMGetUeContextInSmfData(ue)
	if problemDetails != nil {
		return errors.Errorf("%s", problemDetails.Cause)
	} else if err != nil {
		return errors.Wrap(err, "SDM_Get UeContextInSmfData Error")
	}

	problemDetails, err = consumer.GetConsumer().SDMSubscribe(ue)
	if problemDetails != nil {
		return errors.Errorf("%s", problemDetails.Cause)
	} else if err != nil {
		return errors.Wrap(err, "SDM Subscribe Error")
	}
	ue.ContextValid = true
	return nil
}

func getSubscribedNssai(ue *context.AmfUe) {
	amfSelf := context.GetSelf()
	if ue.NudmSDMUri == "" {
		param := Nnrf_NFDiscovery.SearchNFInstancesRequest{
			Supi: &ue.Supi,
		}
		for {
			err := consumer.GetConsumer().SearchUdmSdmInstance(
				ue, amfSelf.NrfUri, models.NrfNfManagementNfType_UDM, models.NrfNfManagementNfType_AMF, &param)
			if err != nil {
				ue.GmmLog.Errorf("AMF can not select an Nudm_SDM Instance by NRF[Error: %+v]", err)
				time.Sleep(2 * time.Second)
			} else {
				break
			}
		}
	}
	problemDetails, err := consumer.GetConsumer().SDMGetSliceSelectionSubscriptionData(ue)
	if problemDetails != nil {
		ue.GmmLog.Errorf("SDM_Get Slice Selection Subscription Data Failed Problem[%+v]", problemDetails)
	} else if err != nil {
		ue.GmmLog.Errorf("SDM_Get Slice Selection Subscription Data Error[%+v]", err)
	}
}

// TS 23.502 4.2.2.2.3 Registration with AMF Re-allocation
func handleRequestedNssai(ue *context.AmfUe, anType models.AccessType) error {
	amfSelf := context.GetSelf()

	if ue.RegistrationRequest.RequestedNSSAI != nil {
		logger.GmmLog.Infof("RequestedNssai: %+v", ue.RegistrationRequest.RequestedNSSAI)
		requestedNssai, err := nasConvert.RequestedNssaiToModels(ue.RegistrationRequest.RequestedNSSAI)
		if err != nil {
			return fmt.Errorf("decode failed at RequestedNSSAI[%s]", err)
		}

		needSliceSelection := false
		for _, requestedSnssai := range requestedNssai {
			ue.GmmLog.Infof("RequestedNssai - ServingSnssai: %+v, HomeSnssai: %+v",
				requestedSnssai.ServingSnssai, requestedSnssai.HomeSnssai)
			if ue.InSubscribedNssai(*requestedSnssai.ServingSnssai) {
				allowedSnssai := models.AllowedSnssai{
					AllowedSnssai: &models.Snssai{
						Sst: requestedSnssai.ServingSnssai.Sst,
						Sd:  requestedSnssai.ServingSnssai.Sd,
					},
					MappedHomeSnssai: requestedSnssai.HomeSnssai,
				}
				if !ue.InAllowedNssai(*allowedSnssai.AllowedSnssai, anType) {
					ue.AllowedNssai[anType] = append(ue.AllowedNssai[anType], allowedSnssai)
				}
			} else {
				needSliceSelection = true
				break
			}

			reqSnssai := models.Snssai{
				Sst: requestedSnssai.ServingSnssai.Sst,
				Sd:  requestedSnssai.ServingSnssai.Sd,
			}

			if !amfSelf.InPlmnSupportList(reqSnssai) {
				needSliceSelection = true
				logger.GmmLog.Warnf("RequestedNssai[%+v] is not supported by AMF", reqSnssai)
				break
			}
		}

		if needSliceSelection {
			if ue.NssfUri == "" {
				for {
					reqParam := Nnrf_NFDiscovery.SearchNFInstancesRequest{}
					errSearchNssf := consumer.GetConsumer().SearchNssfNSSelectionInstance(
						ue, amfSelf.NrfUri, models.NrfNfManagementNfType_NSSF, models.NrfNfManagementNfType_AMF, &reqParam)
					if errSearchNssf != nil {
						ue.GmmLog.Errorf("AMF can not select an NSSF Instance by NRF[Error: %+v]", errSearchNssf)
						time.Sleep(2 * time.Second)
					} else {
						break
					}
				}
			}

			// Step 4
			problemDetails, errNssfGetReg := consumer.GetConsumer().NSSelectionGetForRegistration(ue, requestedNssai)
			if problemDetails != nil {
				ue.GmmLog.Errorf("NSSelection Get Failed Problem[%+v]", problemDetails)
				gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
				return fmt.Errorf("handle Requested Nssai of UE failed")
			} else if errNssfGetReg != nil {
				ue.GmmLog.Errorf("NSSelection Get Error[%+v]", errNssfGetReg)
				gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
				return fmt.Errorf("handle Requested Nssai of UE failed")
			}

			// Step 5: Initial AMF send Namf_Communication_RegistrationCompleteNotify to old AMF
			req := models.UeRegStatusUpdateReqData{
				TransferStatus: models.UeContextTransferStatus_NOT_TRANSFERRED,
			}
			_, problemDetails, err = consumer.GetConsumer().RegistrationStatusUpdate(ue, req)
			if problemDetails != nil {
				ue.GmmLog.Errorf("Registration Status Update Failed Problem[%+v]", problemDetails)
			} else if err != nil {
				ue.GmmLog.Errorf("Registration Status Update Error[%+v]", err)
			}

			// Step 6
			searchTargetAmfQueryParam := Nnrf_NFDiscovery.SearchNFInstancesRequest{}
			if ue.NetworkSliceInfo != nil {
				networkSliceInfo := ue.NetworkSliceInfo
				if networkSliceInfo.TargetAmfSet != "" {
					// TS 29.531
					// TargetAmfSet format: ^[0-9]{3}-[0-9]{2-3}-[A-Fa-f0-9]{2}-[0-3][A-Fa-f0-9]{2}$
					// mcc-mnc-amfRegionId(8 bit)-AmfSetId(10 bit)
					targetAmfSetToken := strings.Split(networkSliceInfo.TargetAmfSet, "-")
					guami := amfSelf.ServedGuamiList[0]
					targetAmfPlmnId := models.PlmnId{
						Mcc: targetAmfSetToken[0],
						Mnc: targetAmfSetToken[1],
					}

					if !reflect.DeepEqual(util.PlmnIdNidToModelsPlmnId(*guami.PlmnId), targetAmfPlmnId) {
						searchTargetAmfQueryParam.TargetPlmnList = []models.PlmnId{targetAmfPlmnId}
						searchTargetAmfQueryParam.RequesterPlmnList = []models.PlmnId{util.PlmnIdNidToModelsPlmnId(*guami.PlmnId)}
					}

					searchTargetAmfQueryParam.AmfRegionId = &targetAmfSetToken[2]
					searchTargetAmfQueryParam.AmfSetId = &targetAmfSetToken[3]
				} else if len(networkSliceInfo.CandidateAmfList) > 0 {
					// TODO: select candidate Amf based on local poilcy
					searchTargetAmfQueryParam.TargetNfInstanceId = &networkSliceInfo.CandidateAmfList[0]
				}
			}

			sendReroute := true
			err = consumer.GetConsumer().SearchAmfCommunicationInstance(ue, amfSelf.NrfUri,
				models.NrfNfManagementNfType_AMF, models.NrfNfManagementNfType_AMF, &searchTargetAmfQueryParam)
			if err == nil {
				// Condition (A) Step 7: initial AMF find Target AMF via NRF ->
				// Send Namf_Communication_N1MessageNotify to Target AMF
				ueContext := consumer.GetConsumer().BuildUeContextModel(ue)
				registerContext := models.RegistrationContextContainer{
					UeContext:        &ueContext,
					AnType:           anType,
					AnN2ApId:         int32(ue.RanUe[anType].RanUeNgapId),
					RanNodeId:        ue.RanUe[anType].Ran.RanId,
					InitialAmfName:   amfSelf.Name,
					UserLocation:     &ue.Location,
					RrcEstCause:      ue.RanUe[anType].RRCEstablishmentCause,
					UeContextRequest: ue.RanUe[anType].UeContextRequest,
					AnN2IPv4Addr:     ue.RanUe[anType].Ran.Conn.RemoteAddr().String(),
					AllowedNssai: &models.AllowedNssai{
						AllowedSnssaiList: ue.AllowedNssai[anType],
						AccessType:        anType,
					},
				}
				if len(ue.NetworkSliceInfo.RejectedNssaiInPlmn) > 0 {
					registerContext.RejectedNssaiInPlmn = ue.NetworkSliceInfo.RejectedNssaiInPlmn
				}
				if len(ue.NetworkSliceInfo.RejectedNssaiInTa) > 0 {
					registerContext.RejectedNssaiInTa = ue.NetworkSliceInfo.RejectedNssaiInTa
				}

				var n1Message bytes.Buffer
				err = ue.RegistrationRequest.EncodeRegistrationRequest(&n1Message)
				if err != nil {
					logger.GmmLog.Errorf("re-encoding registration request message is failed: %+v", err)
				} else {
					err = callback.SendN1MessageNotifyAtAMFReAllocation(ue, n1Message.Bytes(), &registerContext)
					if err != nil {
						logger.GmmLog.Errorf("send N1MessageNotify failed: %+v", err)
					} else {
						sendReroute = false
					}
				}
			}
			if sendReroute {
				// Condition (B) Step 7: initial AMF can not find Target AMF via NRF -> Send Reroute NAS Request to RAN
				allowedNssaiNgap := ngapConvert.AllowedNssaiToNgap(ue.AllowedNssai[anType])
				ngap_message.SendRerouteNasRequest(ue, anType, nil, ue.RanUe[anType].InitialUEMessage, &allowedNssaiNgap)
				return err
			}
			return nil
		}
	}

	// if registration request has no requested nssai, or non of snssai in requested nssai is permitted by nssf
	// then use ue subscribed snssai which is marked as default as allowed nssai
	if len(ue.AllowedNssai[anType]) == 0 {
		for _, snssai := range ue.SubscribedNssai {
			if snssai.DefaultIndication {
				if amfSelf.InPlmnSupportList(*snssai.SubscribedSnssai) {
					allowedSnssai := models.AllowedSnssai{
						AllowedSnssai: snssai.SubscribedSnssai,
					}
					ue.AllowedNssai[anType] = append(ue.AllowedNssai[anType], allowedSnssai)
				}
			}
		}
	}
	return nil
}

func assignLadnInfo(ue *context.AmfUe, accessType models.AccessType) {
	amfSelf := context.GetSelf()

	ue.LadnInfo = nil
	if ue.RegistrationRequest.LADNIndication != nil {
		ue.LadnInfo = make([]factory.Ladn, 0)
		// request for LADN information
		if ue.RegistrationRequest.LADNIndication.GetLen() == 0 {
			if ue.HasWildCardSubscribedDNN() {
				for _, ladn := range amfSelf.LadnPool {
					if ue.TaiListInRegistrationArea(ladn.TaiList, accessType) {
						ue.LadnInfo = append(ue.LadnInfo, ladn)
					}
				}
			} else {
				for _, snssaiInfos := range ue.SmfSelectionData.SubscribedSnssaiInfos {
					for _, dnnInfo := range snssaiInfos.DnnInfos {
						if ladn, ok := amfSelf.LadnPool[dnnInfo.Dnn.(string)]; ok { // check if this dnn is a ladn
							if ue.TaiListInRegistrationArea(ladn.TaiList, accessType) {
								ue.LadnInfo = append(ue.LadnInfo, ladn)
							}
						}
					}
				}
			}
		} else {
			requestedLadnList := nasConvert.LadnToModels(ue.RegistrationRequest.LADNIndication.GetLADNDNNValue())
			for _, requestedLadn := range requestedLadnList {
				if ladn, ok := amfSelf.LadnPool[requestedLadn]; ok {
					if ue.TaiListInRegistrationArea(ladn.TaiList, accessType) {
						ue.LadnInfo = append(ue.LadnInfo, ladn)
					}
				}
			}
		}
	} else if ue.SmfSelectionData != nil {
		for _, snssaiInfos := range ue.SmfSelectionData.SubscribedSnssaiInfos {
			for _, dnnInfo := range snssaiInfos.DnnInfos {
				if dnnInfo.Dnn != "*" {
					if ladn, ok := amfSelf.LadnPool[dnnInfo.Dnn.(string)]; ok {
						if ue.TaiListInRegistrationArea(ladn.TaiList, accessType) {
							ue.LadnInfo = append(ue.LadnInfo, ladn)
						}
					}
				}
			}
		}
	}
}

func reactivatePendingULDataPDUSession(ue *context.AmfUe, anType models.AccessType, serviceType uint8,
	uplinkDataPsi *[psiArraySize]bool, dlPduSessionId int32, cxtList *ngapType.PDUSessionResourceSetupListCxtReq,
	reactivationResult *[psiArraySize]bool, errPduSessionId, errCause []uint8,
) ([]uint8, []uint8) {
	ue.SmContextList.Range(func(key, value interface{}) bool {
		pduSessionID := key.(int32)
		smContext := value.(*context.SmContext)

		// uplink data are pending for the corresponding PDU session identity
		if !uplinkDataPsi[pduSessionID] ||
			(pduSessionID == dlPduSessionId && serviceType == nasMessage.ServiceTypeMobileTerminatedServices) {
			// Skipping SendUpdateSmContextActivateUpCnxState for the following reason:
			//   In Step 4 of 4.2.3.2 UE Triggered Service Request in TS23.502
			//   > This procedure is triggered by the SMF but the PDU Session(s) identified by the UE
			//   > correlates to other PDU Session ID(s) than the one triggering the procedure
			// However, in the case of Mo-data etc., it cannot be skipped because AMF need to know
			// latest N2SmInformation even if the UE has known the N2Information received at
			// previous N1N2MessageTransfer.
			return true
		}

		// indicate the SMF to re-establish the user-plane resources for the corresponding PDU session
		// TODO: determine the UE presence in LADN service area and forward the UE presence
		// in LADN service area towards the SMF, if the corresponding PDU session is
		// a PDU session for LADN
		response, errRsp, problemDetail, err := consumer.GetConsumer().SendUpdateSmContextActivateUpCnxState(
			ue, smContext, anType)
		if err != nil {
			reactivationResult[pduSessionID] = true
			ue.GmmLog.Errorf("SendUpdateSmContextActivateUpCnxState[pduSessionID:%d] Error: %+v",
				pduSessionID, err)
		} else if response == nil {
			reactivationResult[pduSessionID] = true
			errPduSessionId = append(errPduSessionId, uint8(pduSessionID))
			cause := nasMessage.Cause5GMMProtocolErrorUnspecified
			if errRsp != nil {
				switch errRsp.JsonData.Error.Cause {
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
				ue.GmmLog.Errorf("Update SmContext Failed Problem[%+v]", problemDetail)
			} else if err != nil {
				ue.GmmLog.Errorf("Update SmContext Error[%v]", err.Error())
			}
		} else {
			ue.GmmLog.Infof("Re-active the pending uplink PDU Session[%d] over %q successfully",
				pduSessionID, smContext.AccessType())
			ngap_message.AppendPDUSessionResourceSetupListCxtReq(cxtList, pduSessionID,
				smContext.Snssai(), response.BinaryDataN1SmMessage, response.BinaryDataN2SmInformation)
		}
		return true
	})
	return errPduSessionId, errCause
}

func releaseInactivePDUSession(ue *context.AmfUe, anType models.AccessType, uePduStatus *[psiArraySize]bool,
	pduStatusResult *[psiArraySize]bool,
) {
	ue.SmContextList.Range(func(key, value interface{}) bool {
		pduSessionID := key.(int32)
		smContext := value.(*context.SmContext)

		if uePduStatus[pduSessionID] {
			pduStatusResult[pduSessionID] = true
			return true
		}

		// perform a local release of all those PDU session which are in 5GSM state PDU SESSION ACTIVE
		// on the AMF side associated with the access type the REGISTRATION REQUEST message is sent over,
		// but are indicated by the UE as being in 5GSM state PDU SESSION INACTIVE
		cause := models.SmfPduSessionCause_PDU_SESSION_STATUS_MISMATCH
		causeAll := &context.CauseAll{
			Cause: &cause,
		}
		ue.GmmLog.Infof("Release Inactive PDU Session[%d] over  %q", pduSessionID, smContext.AccessType())
		problemDetail, err := consumer.GetConsumer().SendReleaseSmContextRequest(ue, smContext, causeAll, "", nil)
		if problemDetail != nil {
			ue.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
		} else if err != nil {
			ue.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
		}
		return true
	})
}

func reestablishAllowedPDUSessionOver3GPP(ue *context.AmfUe, anType models.AccessType, smContext *context.SmContext,
	allowedPsi *[psiArraySize]bool, cxtList *ngapType.PDUSessionResourceSetupListCxtReq,
	reactivationResult *[psiArraySize]bool, errPduSessionId, errCause []uint8,
) ([]uint8, []uint8) {
	requestData := ue.N1N2Message.Request.JsonData

	if ue.N1N2Message == nil || ue.N1N2Message.Request.BinaryDataN2Information == nil {
		// no pending downlink data
		return errPduSessionId, errCause
	}

	if smContext == nil || smContext.AccessType() != models.AccessType_NON_3_GPP_ACCESS {
		return errPduSessionId, errCause
	}

	// SMF has indicated pending downlink data
	// notify the SMF that reactivation of the user-plane resources for the corresponding PDU session(s)
	// associated with non-3GPP access
	if reactivationResult == nil {
		reactivationResult = new([psiArraySize]bool)
	}
	if allowedPsi[requestData.PduSessionId] {
		// re-establish the PDU session associated with non-3GPP access over 3GPP access.
		// notify the SMF if the corresponding PDU session ID(s) associated with non-3GPP access
		// are indicated in the Allowed PDU session status IE
		// TODO: error handling
		response, errRes, _, err := consumer.GetConsumer().SendUpdateSmContextChangeAccessType(ue, smContext, true)
		if err != nil {
			reactivationResult[requestData.PduSessionId] = true
			ue.GmmLog.Errorf("SendUpdateSmContextActivateUpCnxState[pduSessionID:%d] Error: %+v",
				requestData.PduSessionId, err)
		} else if response == nil {
			ue.GmmLog.Warnf("failed to re-establish allowed PDU Session[%d] over 3GPP access",
				requestData.PduSessionId)
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
			ue.GmmLog.Infof("re-establish allowed PDU Session[%d] over 3GPP access successfully",
				requestData.PduSessionId)
			// the AMF and SMF update the associated access type of the corresponding PDU session
			smContext.SetUserLocation(deepcopy.Copy(ue.Location).(models.UserLocation))
			smContext.SetAccessType(models.AccessType__3_GPP_ACCESS)
			if response.BinaryDataN2SmInformation != nil &&
				response.JsonData.N2SmInfoType == models.N2SmInfoType_PDU_RES_SETUP_REQ {
				// discard the received 5GSM message for PDU session(s) associated with non-3GPP access
				ngap_message.AppendPDUSessionResourceSetupListCxtReq(cxtList, requestData.PduSessionId,
					smContext.Snssai(), nil, response.BinaryDataN2SmInformation)
			}
		}
	} else {
		// notify the SMF if the corresponding PDU session ID(s) associated with non-3GPP access
		// are not indicated in the Allowed PDU session status IE
		ue.GmmLog.Warnf("UE was reachable but did not accept to re-activate the PDU Session[%d]",
			requestData.PduSessionId)
		callback.SendN1N2TransferFailureNotification(ue,
			models.N1N2MessageTransferCause_UE_NOT_REACHABLE_FOR_SESSION)
	}
	return errPduSessionId, errCause
}

func getPDUSessionStatus(ue *context.AmfUe, anType models.AccessType) *[psiArraySize]bool {
	var pduStatusResult [psiArraySize]bool
	ue.SmContextList.Range(func(key, value interface{}) bool {
		pduSessionID := key.(int32)
		smContext := value.(*context.SmContext)

		if smContext.AccessType() != anType {
			return true
		}
		pduStatusResult[pduSessionID] = true
		return true
	})
	return &pduStatusResult
}

func HandleIdentityResponse(ue *context.AmfUe, identityResponse *nasMessage.IdentityResponse) error {
	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	ue.GmmLog.Info("Handle Identity Response")

	mobileIdentityContents := identityResponse.MobileIdentity.GetMobileIdentityContents()
	if len(mobileIdentityContents) < 1 {
		return errors.New("empty Mobile Identity")
	}
	if nasConvert.GetTypeOfIdentity(mobileIdentityContents[0]) != ue.RequestIdentityType {
		return fmt.Errorf("received identity type doesn't match request type")
	}

	if ue.T3570 != nil {
		ue.T3570.Stop()
		ue.T3570 = nil // clear the timer
	}

	switch nasConvert.GetTypeOfIdentity(mobileIdentityContents[0]) { // get type of identity
	case nasMessage.MobileIdentity5GSTypeSuci:
		if suci, plmnId, err := nasConvert.SuciToStringWithError(mobileIdentityContents); err != nil {
			return fmt.Errorf("decode SUCI failed: %w", err)
		} else if plmnId == "" {
			return errors.New("empty plmnId")
		} else {
			ue.Suci = suci
			ue.PlmnId = util.PlmnIdStringToModels(plmnId)
		}
		ue.GmmLog.Debugf("get SUCI: %s", ue.Suci)
	case nasMessage.MobileIdentity5GSType5gGuti:
		if ue.MacFailed {
			return fmt.Errorf("NAS message integrity check failed")
		}
		_, guti, err := nasConvert.GutiToStringWithError(mobileIdentityContents)
		if err != nil {
			return fmt.Errorf("decode GUTI failed: %w", err)
		}
		ue.Guti = guti
		ue.GmmLog.Debugf("get GUTI: %s", guti)
	case nasMessage.MobileIdentity5GSType5gSTmsi:
		if ue.MacFailed {
			return fmt.Errorf("NAS message integrity check failed")
		}
		sTmsi := hex.EncodeToString(mobileIdentityContents[1:])
		if tmp, err := strconv.ParseInt(sTmsi[4:], 10, 32); err != nil {
			return err
		} else {
			ue.Tmsi = int32(tmp)
		}
		ue.GmmLog.Debugf("get 5G-S-TMSI: %s", sTmsi)
	case nasMessage.MobileIdentity5GSTypeImei:
		if ue.MacFailed {
			return fmt.Errorf("NAS message integrity check failed")
		}
		imei, err := nasConvert.PeiToStringWithError(mobileIdentityContents)
		if err != nil {
			return fmt.Errorf("decode PEI failed: %w", err)
		}
		ue.Pei = imei
		ue.GmmLog.Debugf("get PEI: %s", imei)
	case nasMessage.MobileIdentity5GSTypeImeisv:
		if ue.MacFailed {
			return fmt.Errorf("NAS message integrity check failed")
		}
		imeisv, err := nasConvert.PeiToStringWithError(mobileIdentityContents)
		if err != nil {
			return fmt.Errorf("decode PEI failed: %w", err)
		}
		ue.Pei = imeisv
		ue.GmmLog.Debugf("get PEI: %s", imeisv)
	}
	return nil
}

// TS 24501 5.6.3.2
func HandleNotificationResponse(ue *context.AmfUe, notificationResponse *nasMessage.NotificationResponse) error {
	ue.GmmLog.Info("Handle Notification Response")

	if ue.MacFailed {
		return fmt.Errorf("NAS message integrity check failed")
	}

	ue.StopT3565()

	if notificationResponse != nil && notificationResponse.PDUSessionStatus != nil {
		psiArray := nasConvert.PSIToBooleanArray(notificationResponse.PDUSessionStatus.Buffer)
		for psi := 1; psi <= 15; psi++ {
			pduSessionId := int32(psi)
			if smContext, ok := ue.SmContextFindByPDUSessionID(pduSessionId); ok {
				if !psiArray[psi] {
					cause := models.SmfPduSessionCause_PDU_SESSION_STATUS_MISMATCH
					causeAll := &context.CauseAll{
						Cause: &cause,
					}
					problemDetail, err := consumer.GetConsumer().SendReleaseSmContextRequest(ue, smContext, causeAll, "", nil)
					if problemDetail != nil {
						ue.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
					} else if err != nil {
						ue.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
					}
				}
			}
		}
	}
	return nil
}

func HandleConfigurationUpdateComplete(ue *context.AmfUe,
	configurationUpdateComplete *nasMessage.ConfigurationUpdateComplete,
) error {
	ue.GmmLog.Info("Handle Configuration Update Complete")

	if ue.MacFailed {
		return fmt.Errorf("NAS message integrity check failed")
	}

	// Stop timer T3555 in TS 24.501 Figure 5.4.4.1.1 in handler
	ue.StopT3555()
	// TODO: Send acknowledgment by Nudm_SMD_Info_Service to UDM in handler
	//		import "github.com/free5gc/openapi/Nudm_SubscriberDataManagement" client.Info

	return nil
}

func AuthenticationProcedure(ue *context.AmfUe, accessType models.AccessType) (bool, error) {
	ue.GmmLog.Info("Authentication procedure")

	// Check whether UE has SUCI and SUPI
	if IdentityVerification(ue) {
		ue.GmmLog.Debugln("UE has SUCI / SUPI")
		if ue.SecurityContextIsValid() {
			ue.GmmLog.Debugln("UE has a valid security context - skip the authentication procedure")
			return true, nil
		}
	} else {
		// Request UE's SUCI by sending identity request
		ue.IdentityRequestSendTimes++
		gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
		return false, nil
	}

	amfSelf := context.GetSelf()

	// TODO: consider ausf group id, Routing ID part of SUCI
	param := Nnrf_NFDiscovery.SearchNFInstancesRequest{}
	resp, err := consumer.GetConsumer().SendSearchNFInstances(
		amfSelf.NrfUri, models.NrfNfManagementNfType_AUSF, models.NrfNfManagementNfType_AMF, &param)
	if err != nil {
		ue.GmmLog.Error("AMF can not select an AUSF by NRF")
		gmm_message.SendRegistrationReject(ue.RanUe[accessType], nasMessage.Cause5GMMCongestion, "")
		return false, err
	}

	// select the first AUSF, TODO: select base on other info
	var ausfUri string
	for index := range resp.NfInstances {
		ue.AusfId = resp.NfInstances[index].NfInstanceId
		ausfUri = util.SearchNFServiceUri(&resp.NfInstances[index], models.ServiceName_NAUSF_AUTH,
			models.NfServiceStatus_REGISTERED)
		if ausfUri != "" {
			break
		}
	}
	if ausfUri == "" {
		err = fmt.Errorf("AMF can not select an AUSF by NRF")
		ue.GmmLog.Error(err)
		gmm_message.SendRegistrationReject(ue.RanUe[accessType], nasMessage.Cause5GMMCongestion, "")
		return false, err
	}
	ue.AusfUri = ausfUri

	response, problemDetails, err := consumer.GetConsumer().SendUEAuthenticationAuthenticateRequest(ue, nil)
	if err != nil {
		ue.GmmLog.Errorf("Nausf_UEAU Authenticate Request Error: %+v", err)
		gmm_message.SendRegistrationReject(ue.RanUe[accessType], nasMessage.Cause5GMMCongestion, "")
		err = fmt.Errorf("Authentication procedure failed")
		ue.GmmLog.Error(err)
		return false, err
	} else if problemDetails != nil {
		ue.GmmLog.Warnf("Nausf_UEAU Authenticate Request Failed: %+v", problemDetails)
		var cause uint8
		switch problemDetails.Status {
		case http.StatusForbidden, http.StatusNotFound:
			cause = nasMessage.Cause5GMMIllegalUE
		default:
			cause = nasMessage.Cause5GMMCongestion
		}
		gmm_message.SendRegistrationReject(ue.RanUe[accessType], cause, "")
		err = fmt.Errorf("Authentication procedure failed")
		ue.GmmLog.Warn(err)
		return false, err
	}
	ue.AuthenticationCtx = response
	ue.ABBA = []uint8{0x00, 0x00} // set ABBA value as described at TS 33.501 Annex A.7.1

	gmm_message.SendAuthenticationRequest(ue.RanUe[accessType])
	return false, nil
}

// TS 24501 5.6.1
func HandleServiceRequest(ue *context.AmfUe, anType models.AccessType,
	serviceRequest *nasMessage.ServiceRequest,
) error {
	if ue == nil {
		return fmt.Errorf("AmfUe is nil")
	}

	ue.GmmLog.Info("Handle Service Request")

	ue.StopT3513()
	ue.StopT3565()

	// Set No ongoing
	if procedure := ue.OnGoing(anType).Procedure; procedure == context.OnGoingProcedurePaging {
		ue.SetOnGoing(anType, &context.OnGoing{
			Procedure: context.OnGoingProcedureNothing,
		})
	} else if procedure != context.OnGoingProcedureNothing {
		ue.GmmLog.Warnf("UE should not in OnGoing[%s]", procedure)
	}

	var pduStatusResult *[psiArraySize]bool
	if serviceRequest.PDUSessionStatus != nil {
		pduStatusResult = getPDUSessionStatus(ue, anType)
	}

	// Send Authtication / Security Procedure not support
	if !ue.SecurityContextIsValid() {
		ue.GmmLog.Warnf("No Security Context : SUPI[%s]", ue.Supi)
		gmm_message.SendServiceReject(ue.RanUe[anType], pduStatusResult,
			nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork)
		ngap_message.SendUEContextReleaseCommand(ue.RanUe[anType],
			context.UeContextN2NormalRelease, ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
		return nil
	}

	// TS 24.501 8.2.6.21: if the UE is sending a REGISTRATION REQUEST message as an initial NAS message,
	// the UE has a valid 5G NAS security context and the UE needs to send non-cleartext IEs
	// TS 24.501 4.4.6: When the UE sends a REGISTRATION REQUEST or SERVICE REQUEST message that includes a NAS message
	// container IE, the UE shall set the security header type of the initial NAS message to "integrity protected"
	if serviceRequest.NASMessageContainer != nil {
		contents := serviceRequest.NASMessageContainer.GetNASMessageContainerContents()

		// TS 24.501 4.4.6: When the UE sends a REGISTRATION REQUEST or SERVICE REQUEST message that includes a NAS
		// message container IE, the UE shall set the security header type of the initial NAS message to
		// "integrity protected"; then the AMF shall decipher the value part of the NAS message container IE
		err := security.NASEncrypt(ue.CipheringAlg, ue.KnasEnc, ue.ULCount.Get(), security.Bearer3GPP,
			security.DirectionUplink, contents)

		if err != nil {
			ue.SecurityContextAvailable = false
		} else {
			m := nas.NewMessage()
			if errGmmMessageDecode := m.GmmMessageDecode(&contents); errGmmMessageDecode != nil {
				return errGmmMessageDecode
			}

			messageType := m.GmmMessage.GmmHeader.GetMessageType()
			if messageType != nas.MsgTypeServiceRequest {
				return errors.New("The payload of NAS message Container is not service request")
			}
			// TS 24.501 4.4.6: The AMF shall consider the NAS message that is obtained from the NAS message container
			// IE as the initial NAS message that triggered the procedure
			serviceRequest = m.ServiceRequest
		}
		// TS 33.501 6.4.6 step 3: if the initial NAS message was protected but did not pass the integrity check
		ue.RetransmissionOfInitialNASMsg = ue.MacFailed
	}

	serviceType := serviceRequest.GetServiceTypeValue()
	var reactivationResult *[psiArraySize]bool
	var errPduSessionId, errCause []uint8
	var dlPduSessionId int32
	cxtList := ngapType.PDUSessionResourceSetupListCxtReq{}

	if serviceType == nasMessage.ServiceTypeEmergencyServices ||
		serviceType == nasMessage.ServiceTypeEmergencyServicesFallback {
		ue.GmmLog.Warnf("emergency service is not supported")
		gmm_message.SendServiceReject(ue.RanUe[anType], pduStatusResult, nasMessage.Cause5GMM5GSServicesNotAllowed)
		ngap_message.SendUEContextReleaseCommand(ue.RanUe[anType],
			context.UeContextN2NormalRelease, ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
		return nil
	}

	if serviceType == nasMessage.ServiceTypeSignalling {
		err := gmm_message.SendServiceAccept(ue, anType, cxtList, pduStatusResult, nil, nil, nil)
		return err
	}

	var N1N2ReqData *models.N1N2MessageTransferReqData
	var n1Msg, n2Info []byte
	if ue.N1N2Message != nil {
		N1N2ReqData = ue.N1N2Message.Request.JsonData
		n1Msg = ue.N1N2Message.Request.BinaryDataN1Message
		n2Info = ue.N1N2Message.Request.BinaryDataN2Information
		if n2Info != nil {
			if N1N2ReqData.N2InfoContainer.N2InformationClass == models.N2InformationClass_SM {
				dlPduSessionId = N1N2ReqData.N2InfoContainer.SmInfo.PduSessionId
			} else {
				ue.N1N2Message = nil
				return fmt.Errorf("service request triggered by network has not implemented about non SM N2Info")
			}
		}
	}

	if serviceRequest.UplinkDataStatus != nil {
		uplinkDataPsi := nasConvert.PSIToBooleanArray(serviceRequest.UplinkDataStatus.Buffer)
		if reactivationResult == nil {
			reactivationResult = new([psiArraySize]bool)
		}
		errPduSessionId, errCause = reactivatePendingULDataPDUSession(ue, anType, serviceType, &uplinkDataPsi,
			dlPduSessionId, &cxtList, reactivationResult, errPduSessionId, errCause)
	}

	if serviceRequest.PDUSessionStatus != nil {
		uePduStatus := nasConvert.PSIToBooleanArray(serviceRequest.PDUSessionStatus.Buffer)
		if pduStatusResult == nil {
			pduStatusResult = new([psiArraySize]bool)
		}
		releaseInactivePDUSession(ue, anType, &uePduStatus, pduStatusResult)
	}

	switch serviceType {
	case nasMessage.ServiceTypeMobileTerminatedServices: // Trigger by Network
		if ue.N1N2Message != nil {
			// downlink signaling only
			if n2Info == nil {
				err := gmm_message.SendServiceAccept(ue, anType, cxtList, pduStatusResult,
					reactivationResult, errPduSessionId, errCause)
				if err != nil {
					return err
				}
				switch N1N2ReqData.N1MessageContainer.N1MessageClass {
				case models.N1MessageClass_SM:
					gmm_message.SendDLNASTransport(ue.RanUe[anType],
						nasMessage.PayloadContainerTypeN1SMInfo, n1Msg, N1N2ReqData.PduSessionId, 0, nil, 0)
				case models.N1MessageClass_LPP:
					gmm_message.SendDLNASTransport(ue.RanUe[anType],
						nasMessage.PayloadContainerTypeLPP, n1Msg, 0, 0, nil, 0)
				case models.N1MessageClass_SMS:
					gmm_message.SendDLNASTransport(ue.RanUe[anType],
						nasMessage.PayloadContainerTypeSMS, n1Msg, 0, 0, nil, 0)
				case models.N1MessageClass_UPDP:
					gmm_message.SendDLNASTransport(ue.RanUe[anType],
						nasMessage.PayloadContainerTypeUEPolicy, n1Msg, 0, 0, nil, 0)
				}
				ue.N1N2Message = nil
				return nil
			}

			// TODO: Area of validity for the N2 SM information
			smInfo := N1N2ReqData.N2InfoContainer.SmInfo
			smContext, ok := ue.SmContextFindByPDUSessionID(N1N2ReqData.PduSessionId)
			if !ok {
				return fmt.Errorf("service request triggered by network error for pduSession[%d] does not exist",
					N1N2ReqData.PduSessionId)
			}

			if smContext.AccessType() == models.AccessType_NON_3_GPP_ACCESS {
				if serviceRequest.AllowedPDUSessionStatus != nil {
					allowPduSessionPsi := nasConvert.PSIToBooleanArray(serviceRequest.AllowedPDUSessionStatus.Buffer)
					errPduSessionId, errCause = reestablishAllowedPDUSessionOver3GPP(ue, anType, smContext,
						&allowPduSessionPsi, &cxtList, reactivationResult, errPduSessionId, errCause)
				}
			} else if smInfo.N2InfoContent.NgapIeType == models.AmfCommunicationNgapIeType_PDU_RES_SETUP_REQ {
				var nasPdu []byte
				var err error
				if n1Msg != nil {
					pduSessionId := uint8(smInfo.PduSessionId)
					nasPdu, err = gmm_message.BuildDLNASTransport(ue, anType, nasMessage.PayloadContainerTypeN1SMInfo,
						n1Msg, pduSessionId, nil, nil, 0)
					if err != nil {
						return err
					}
				}
				ngap_message.AppendPDUSessionResourceSetupListCxtReq(&cxtList, smInfo.PduSessionId, *smInfo.SNssai,
					nasPdu, n2Info)
			}
			err := gmm_message.SendServiceAccept(ue, anType, cxtList, pduStatusResult,
				reactivationResult, errPduSessionId, errCause)
			if err != nil {
				return err
			}
		}

		// downlink signaling
		if ue.ConfigurationUpdateCommandFlags != nil {
			err := gmm_message.SendServiceAccept(ue, anType, cxtList,
				pduStatusResult, reactivationResult, errPduSessionId, errCause)
			if err != nil {
				return err
			}
			gmm_message.SendConfigurationUpdateCommand(ue, anType, ue.ConfigurationUpdateCommandFlags)
			ue.ConfigurationUpdateCommandFlags = nil
		}
	case nasMessage.ServiceTypeData:
		if anType == models.AccessType__3_GPP_ACCESS {
			if ue.AmPolicyAssociation != nil && ue.AmPolicyAssociation.ServAreaRes != nil {
				var accept bool
				switch ue.AmPolicyAssociation.ServAreaRes.RestrictionType {
				case models.RestrictionType_ALLOWED_AREAS:
					accept = context.TacInAreas(ue.Tai.Tac, ue.AmPolicyAssociation.ServAreaRes.Areas)
				case models.RestrictionType_NOT_ALLOWED_AREAS:
					accept = !context.TacInAreas(ue.Tai.Tac, ue.AmPolicyAssociation.ServAreaRes.Areas)
				}

				if !accept {
					gmm_message.SendServiceReject(ue.RanUe[anType], nil, nasMessage.Cause5GMMRestrictedServiceArea)
					return nil
				}
			}
		}
		err := gmm_message.SendServiceAccept(ue, anType, cxtList, pduStatusResult,
			reactivationResult, errPduSessionId, errCause)
		if err != nil {
			return err
		}
	case nasMessage.ServiceTypeHighPriorityAccess:
		// TODO: support HighPriorityAccess
		err := gmm_message.SendServiceAccept(ue, anType, cxtList, pduStatusResult,
			reactivationResult, errPduSessionId, errCause)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("service type[%d] is not supported", serviceType)
	}
	if len(errPduSessionId) != 0 {
		ue.GmmLog.Info(errPduSessionId, errCause)
	}
	ue.N1N2Message = nil
	return nil
}

// TS 24.501 5.4.1
func HandleAuthenticationResponse(ue *context.AmfUe, accessType models.AccessType,
	authenticationResponse *nasMessage.AuthenticationResponse,
) error {
	ue.GmmLog.Info("Handle Authentication Response")

	ue.StopT3560()

	if ue.AuthenticationCtx == nil {
		return fmt.Errorf("ue authentication context is nil")
	}

	switch ue.AuthenticationCtx.AuthType {
	case models.AusfUeAuthenticationAuthType__5_G_AKA:
		var av5gAka models.Av5gAka
		if err := mapstructure.Decode(ue.AuthenticationCtx.Var5gAuthData, &av5gAka); err != nil {
			return fmt.Errorf("Var5gAuthData Convert Type Error")
		}
		if authenticationResponse.AuthenticationResponseParameter == nil {
			return fmt.Errorf("AuthenticationResponseParamete is nil")
		}
		resStar := authenticationResponse.AuthenticationResponseParameter.GetRES()

		// Calculate HRES* (TS 33.501 Annex A.5)
		p0, err := hex.DecodeString(av5gAka.Rand)
		if err != nil {
			return err
		}
		p1 := resStar[:]
		concat := []byte{}
		concat = append(concat, p0...)
		concat = append(concat, p1...)
		hResStarBytes := sha256.Sum256(concat)
		hResStar := hex.EncodeToString(hResStarBytes[16:])

		if hResStar != av5gAka.HxresStar {
			ue.GmmLog.Errorf("HRES* Validation Failure (received: %s, expected: %s)", hResStar, av5gAka.HxresStar)

			if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti && ue.IdentityRequestSendTimes == 0 {
				ue.IdentityRequestSendTimes++
				gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
				return nil
			} else {
				gmm_message.SendAuthenticationReject(ue.RanUe[accessType], "")
				return GmmFSM.SendEvent(ue.State[accessType], AuthFailEvent, fsm.ArgsType{
					ArgAmfUe:      ue,
					ArgAccessType: accessType,
				}, logger.GmmLog)
			}
		}

		response, problemDetails, err := consumer.GetConsumer().SendAuth5gAkaConfirmRequest(
			ue, hex.EncodeToString(resStar[:]))
		if err != nil {
			return err
		} else if problemDetails != nil {
			ue.GmmLog.Debugf("Auth5gAkaConfirm Error[Problem Detail: %+v]", problemDetails)
			return nil
		}
		switch response.AuthResult {
		case models.AusfUeAuthenticationAuthResult_SUCCESS:
			ue.UnauthenticatedSupi = false
			ue.Kseaf = response.Kseaf
			ue.Supi = response.Supi
			ue.DerivateKamf()
			ue.GmmLog.Debugln("ue.DerivateKamf()", ue.Kamf)
			return GmmFSM.SendEvent(ue.State[accessType], AuthSuccessEvent, fsm.ArgsType{
				ArgAmfUe:      ue,
				ArgAccessType: accessType,
				ArgEAPSuccess: false,
				ArgEAPMessage: "",
			}, logger.GmmLog)
		case models.AusfUeAuthenticationAuthResult_FAILURE:
			if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti && ue.IdentityRequestSendTimes == 0 {
				ue.IdentityRequestSendTimes++
				gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
				return nil
			} else {
				gmm_message.SendAuthenticationReject(ue.RanUe[accessType], "")
				return GmmFSM.SendEvent(ue.State[accessType], AuthFailEvent, fsm.ArgsType{
					ArgAmfUe:      ue,
					ArgAccessType: accessType,
				}, logger.GmmLog)
			}
		}
	case models.AusfUeAuthenticationAuthType_EAP_AKA_PRIME:
		response, pd, err := consumer.GetConsumer().SendEapAuthConfirmRequest(ue, *authenticationResponse.EAPMessage)
		if err != nil {
			return err
		} else if pd != nil {
			ue.GmmLog.Debugf("EapAuthConfirm Error[Problem Detail: %+v]", pd)
			return nil
		}

		switch response.AuthResult {
		case models.AusfUeAuthenticationAuthResult_SUCCESS:
			ue.UnauthenticatedSupi = false
			ue.Kseaf = response.KSeaf
			ue.Supi = response.Supi
			ue.DerivateKamf()
			// TODO: select enc/int algorithm based on ue security capability & amf's policy,
			// then generate KnasEnc, KnasInt
			return GmmFSM.SendEvent(ue.State[accessType], AuthSuccessEvent, fsm.ArgsType{
				ArgAmfUe:      ue,
				ArgAccessType: accessType,
				ArgEAPSuccess: true,
				ArgEAPMessage: response.EapPayload,
			}, logger.GmmLog)
		case models.AusfUeAuthenticationAuthResult_FAILURE:
			if ue.IdentityTypeUsedForRegistration == nasMessage.MobileIdentity5GSType5gGuti && ue.IdentityRequestSendTimes == 0 {
				ue.IdentityRequestSendTimes++
				gmm_message.SendAuthenticationResult(ue.RanUe[accessType], false, response.EapPayload)
				gmm_message.SendIdentityRequest(ue.RanUe[accessType], accessType, nasMessage.MobileIdentity5GSTypeSuci)
				return nil
			} else {
				gmm_message.SendAuthenticationReject(ue.RanUe[accessType], response.EapPayload)
				return GmmFSM.SendEvent(ue.State[accessType], AuthFailEvent, fsm.ArgsType{
					ArgAmfUe:      ue,
					ArgAccessType: accessType,
				}, logger.GmmLog)
			}
		case models.AusfUeAuthenticationAuthResult_ONGOING:
			ue.AuthenticationCtx.Var5gAuthData = response.EapPayload
			if _, exists := response.Links["eap-session"]; exists {
				ue.AuthenticationCtx.Links = response.Links
			}
			gmm_message.SendAuthenticationRequest(ue.RanUe[accessType])
		}
	}

	return nil
}

func HandleAuthenticationError(ue *context.AmfUe, anType models.AccessType) error {
	ue.GmmLog.Info("Handle Authentication Error")
	if ue.RegistrationRequest != nil {
		gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMTrackingAreaNotAllowed, "")
	}

	return nil
}

func HandleAuthenticationFailure(ue *context.AmfUe, anType models.AccessType,
	authenticationFailure *nasMessage.AuthenticationFailure,
) error {
	ue.GmmLog.Info("Handle Authentication Failure")

	ue.StopT3560()

	cause5GMM := authenticationFailure.Cause5GMM.GetCauseValue()

	switch ue.AuthenticationCtx.AuthType {
	case models.AusfUeAuthenticationAuthType__5_G_AKA:
		switch cause5GMM {
		case nasMessage.Cause5GMMMACFailure:
			ue.GmmLog.Warnln("Authentication Failure Cause: Mac Failure")
			gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
			return GmmFSM.SendEvent(
				ue.State[anType],
				AuthFailEvent,
				fsm.ArgsType{
					ArgAmfUe:      ue,
					ArgAccessType: anType,
				},
				logger.GmmLog,
			)
		case nasMessage.Cause5GMMNon5GAuthenticationUnacceptable:
			ue.GmmLog.Warnln("Authentication Failure Cause: Non-5G Authentication Unacceptable")
			gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
			return GmmFSM.SendEvent(
				ue.State[anType],
				AuthFailEvent,
				fsm.ArgsType{
					ArgAmfUe:      ue,
					ArgAccessType: anType,
				},
				logger.GmmLog,
			)
		case nasMessage.Cause5GMMngKSIAlreadyInUse:
			ue.GmmLog.Warnln("Authentication Failure Cause: NgKSI Already In Use")
			ue.AuthFailureCauseSynchFailureTimes = 0
			ue.GmmLog.Warnln("Select new NgKsi")
			// select new ngksi
			if ue.NgKsi.Ksi < 6 { // ksi is range from 0 to 6
				ue.NgKsi.Ksi += 1
			} else {
				ue.NgKsi.Ksi = 0
			}
			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		case nasMessage.Cause5GMMSynchFailure: // TS 24.501 5.4.1.3.7 case f
			ue.GmmLog.Warn("Authentication Failure 5GMM Cause: Synch Failure")

			ue.AuthFailureCauseSynchFailureTimes++
			if ue.AuthFailureCauseSynchFailureTimes >= 2 {
				ue.GmmLog.Warnf("2 consecutive Synch Failure, terminate authentication procedure")
				gmm_message.SendAuthenticationReject(ue.RanUe[anType], "")
				return GmmFSM.SendEvent(
					ue.State[anType],
					AuthFailEvent,
					fsm.ArgsType{
						ArgAmfUe:      ue,
						ArgAccessType: anType,
					},
					logger.GmmLog,
				)
			}

			var av5gAka models.Av5gAka
			if err := mapstructure.Decode(ue.AuthenticationCtx.Var5gAuthData, &av5gAka); err != nil {
				ue.GmmLog.Error("Var5gAuthData Convert Type Error")
				return err
			}

			if authenticationFailure.AuthenticationFailureParameter == nil {
				return errors.New("AuthenticationFailureParameter is nil")
			}
			auts := authenticationFailure.AuthenticationFailureParameter.GetAuthenticationFailureParameter()
			resynchronizationInfo := &models.ResynchronizationInfo{
				Auts: hex.EncodeToString(auts[:]),
				Rand: av5gAka.Rand,
			}

			response, pd, err := consumer.GetConsumer().SendUEAuthenticationAuthenticateRequest(ue, resynchronizationInfo)
			if err != nil {
				return err
			} else if pd != nil {
				ue.GmmLog.Errorf("Nausf_UEAU Authenticate Request Error[Problem Detail: %+v]", pd)
				return nil
			}
			ue.AuthenticationCtx = response
			ue.ABBA = []uint8{0x00, 0x00}

			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		}
	case models.AusfUeAuthenticationAuthType_EAP_AKA_PRIME:
		switch cause5GMM {
		case nasMessage.Cause5GMMngKSIAlreadyInUse:
			ue.GmmLog.Warn("Authentication Failure 5GMM Cause: NgKSI Already In Use")
			if ue.NgKsi.Ksi < 6 { // ksi is range from 0 to 6
				ue.NgKsi.Ksi += 1
			} else {
				ue.NgKsi.Ksi = 0
			}
			gmm_message.SendAuthenticationRequest(ue.RanUe[anType])
		default:
		}
	}

	return nil
}

func HandleRegistrationComplete(ue *context.AmfUe, accessType models.AccessType,
	registrationComplete *nasMessage.RegistrationComplete,
) error {
	ue.GmmLog.Info("Handle Registration Complete")

	ue.StopT3550()

	// Release existed old SmContext when Initial Registration completed
	if ue.RegistrationType5GS == nasMessage.RegistrationType5GSInitialRegistration {
		ue.SmContextList.Range(func(key, value interface{}) bool {
			smContext := value.(*context.SmContext)

			if smContext.AccessType() == accessType {
				problemDetail, err := consumer.GetConsumer().SendReleaseSmContextRequest(ue, smContext, nil, "", nil)
				if problemDetail != nil {
					ue.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
				} else if err != nil {
					ue.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
				}
			}
			return true
		})
	}

	// Send NITZ information to UE
	configurationUpdateCommandFlags := &context.ConfigurationUpdateCommandFlags{
		NeedNITZ: true,
	}
	gmm_message.SendConfigurationUpdateCommand(ue, accessType, configurationUpdateCommandFlags)

	// if registrationComplete.SORTransparentContainer != nil {
	// 	TODO: if at regsitration procedure 14b, udm provide amf Steering of Roaming info & request an ack,
	// 	AMF provides the UE's ack with Nudm_SDM_Info (SOR not supportted in this stage)
	// }

	// TODO: if
	//	1. AMF has evaluated the support of IMS Voice over PS Sessions (TS 23.501 5.16.3.2)
	//	2. AMF determines that it needs to update the Homogeneous Support of IMS Voice over PS Sessions (TS 23.501 5.16.3.3)
	// Then invoke Nudm_UECM_Update to send "Homogeneous Support of IMS Voice over PS Sessions" indication to udm

	if ue.RegistrationRequest.UplinkDataStatus == nil &&
		ue.RegistrationRequest.GetFOR() == nasMessage.FollowOnRequestNoPending {
		ngap_message.SendUEContextReleaseCommand(ue.RanUe[accessType], context.UeContextN2NormalRelease,
			ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
	}
	return GmmFSM.SendEvent(ue.State[accessType], ContextSetupSuccessEvent, fsm.ArgsType{
		ArgAmfUe:      ue,
		ArgAccessType: accessType,
	}, logger.GmmLog)
}

// TS 33.501 6.7.2
func HandleSecurityModeComplete(ue *context.AmfUe, anType models.AccessType, procedureCode int64,
	securityModeComplete *nasMessage.SecurityModeComplete,
) error {
	ue.GmmLog.Info("Handle Security Mode Complete")

	if ue.MacFailed {
		return fmt.Errorf("NAS message integrity check failed")
	}

	ue.StopT3560()

	if ue.SecurityContextIsValid() {
		// update Kgnb/Kn3iwf
		ue.UpdateSecurityContext(anType)
	}

	if securityModeComplete.IMEISV != nil {
		ue.GmmLog.Debugln("receieve IMEISV")
		if pei, err := nasConvert.PeiToStringWithError(securityModeComplete.IMEISV.Octet[:]); err != nil {
			gmm_message.SendRegistrationReject(ue.RanUe[anType], nasMessage.Cause5GMMProtocolErrorUnspecified, "")
			return fmt.Errorf("decode PEI failed: %w", err)
		} else {
			ue.Pei = pei
		}
	}

	// TODO: AMF shall set the NAS COUNTs to zero if horizontal derivation of KAMF is performed
	if securityModeComplete.NASMessageContainer != nil {
		contents := securityModeComplete.NASMessageContainer.GetNASMessageContainerContents()
		m := nas.NewMessage()
		if err := m.GmmMessageDecode(&contents); err != nil {
			return err
		}

		argsType := fsm.ArgsType{ArgAmfUe: ue, ArgAccessType: anType, ArgProcedureCode: procedureCode}
		event := SecurityModeSuccessEvent
		switch m.GmmMessage.GmmHeader.GetMessageType() {
		case nas.MsgTypeRegistrationRequest:
			argsType[ArgNASMessage] = m.GmmMessage.RegistrationRequest
		case nas.MsgTypeServiceRequest:
			argsType[ArgNASMessage] = m.GmmMessage.ServiceRequest
			if !ue.State[anType].Is(context.Registered) {
				gmm_message.SendServiceReject(ue.RanUe[anType], nil, nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork)
				ue.GmmLog.Warnf("Service Request was sent when UE state was not Registered")
				ngap_message.SendUEContextReleaseCommand(ue.RanUe[anType],
					context.UeContextN2NormalRelease, ngapType.CausePresentNas, ngapType.CauseNasPresentNormalRelease)
				event = SecurityModeFailEvent
			}
		default:
			ue.GmmLog.Errorln("nas message container Iei type error")
			return errors.New("nas message container Iei type error")
		}
		return GmmFSM.SendEvent(ue.State[anType], event, argsType, logger.GmmLog)
	}
	return GmmFSM.SendEvent(ue.State[anType], SecurityModeSuccessEvent, fsm.ArgsType{
		ArgAmfUe:         ue,
		ArgAccessType:    anType,
		ArgProcedureCode: procedureCode,
		ArgNASMessage:    ue.RegistrationRequest,
	}, logger.GmmLog)
}

func HandleSecurityModeReject(ue *context.AmfUe, anType models.AccessType,
	securityModeReject *nasMessage.SecurityModeReject,
) error {
	ue.GmmLog.Info("Handle Security Mode Reject")

	ue.StopT3560()

	cause := securityModeReject.Cause5GMM.GetCauseValue()
	ue.GmmLog.Warnf("Reject Cause: %s", nasMessage.Cause5GMMToString(cause))
	ue.GmmLog.Error("UE reject the security mode command, abort the ongoing procedure")
	return nil
}

// TS 23.502 4.2.2.3
func HandleDeregistrationRequest(ue *context.AmfUe, anType models.AccessType,
	deregistrationRequest *nasMessage.DeregistrationRequestUEOriginatingDeregistration,
) error {
	ue.GmmLog.Info("Handle Deregistration Request(UE Originating)")

	targetDeregistrationAccessType := deregistrationRequest.GetAccessType()
	ue.SmContextList.Range(func(key, value interface{}) bool {
		smContext := value.(*context.SmContext)

		if smContext.AccessType() == anType ||
			targetDeregistrationAccessType == nasMessage.AccessTypeBoth {
			problemDetail, err := consumer.GetConsumer().SendReleaseSmContextRequest(ue, smContext, nil, "", nil)
			if problemDetail != nil {
				ue.GmmLog.Errorf("Release SmContext Failed Problem[%+v]", problemDetail)
			} else if err != nil {
				ue.GmmLog.Errorf("Release SmContext Error[%v]", err.Error())
			}
		}
		return true
	})

	if ue.AmPolicyAssociation != nil {
		terminateAmPolicyAssocaition := true
		switch anType {
		case models.AccessType__3_GPP_ACCESS:
			terminateAmPolicyAssocaition = ue.State[models.AccessType_NON_3_GPP_ACCESS].Is(context.Deregistered)
		case models.AccessType_NON_3_GPP_ACCESS:
			terminateAmPolicyAssocaition = ue.State[models.AccessType__3_GPP_ACCESS].Is(context.Deregistered)
		}

		if terminateAmPolicyAssocaition {
			problemDetails, err := consumer.GetConsumer().AMPolicyControlDelete(ue)
			if problemDetails != nil {
				ue.GmmLog.Errorf("AM Policy Control Delete Failed Problem[%+v]", problemDetails)
			} else if err != nil {
				ue.GmmLog.Errorf("AM Policy Control Delete Error[%v]", err.Error())
			}
		}
	}

	gmm_common.PurgeAmfUeSubscriberData(ue)

	// if Deregistration type is not switch-off, send Deregistration Accept
	if deregistrationRequest.GetSwitchOff() == 0 {
		gmm_message.SendDeregistrationAccept(ue.RanUe[anType])
	}

	// TS 23.502 4.2.6, 4.12.3
	switch targetDeregistrationAccessType {
	case nasMessage.AccessType3GPP:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		return GmmFSM.SendEvent(ue.State[models.AccessType__3_GPP_ACCESS], DeregistrationAcceptEvent, fsm.ArgsType{
			ArgAmfUe:      ue,
			ArgAccessType: anType,
		}, logger.GmmLog)
	case nasMessage.AccessTypeNon3GPP:
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		return GmmFSM.SendEvent(ue.State[models.AccessType_NON_3_GPP_ACCESS], DeregistrationAcceptEvent, fsm.ArgsType{
			ArgAmfUe:      ue,
			ArgAccessType: anType,
		}, logger.GmmLog)
	case nasMessage.AccessTypeBoth:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}

		err := GmmFSM.SendEvent(ue.State[models.AccessType__3_GPP_ACCESS], DeregistrationAcceptEvent, fsm.ArgsType{
			ArgAmfUe:      ue,
			ArgAccessType: anType,
		}, logger.GmmLog)
		if err != nil {
			ue.GmmLog.Errorln(err)
		}
		return GmmFSM.SendEvent(ue.State[models.AccessType_NON_3_GPP_ACCESS], DeregistrationAcceptEvent, fsm.ArgsType{
			ArgAmfUe:      ue,
			ArgAccessType: anType,
		}, logger.GmmLog)
	}

	return nil
}

// TS 23.502 4.2.2.3
func HandleDeregistrationAccept(ue *context.AmfUe, anType models.AccessType,
	deregistrationAccept *nasMessage.DeregistrationAcceptUETerminatedDeregistration,
) error {
	ue.GmmLog.Info("Handle Deregistration Accept(UE Terminated)")

	ue.StopT3522()

	switch ue.DeregistrationTargetAccessType {
	case nasMessage.AccessType3GPP:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
	case nasMessage.AccessTypeNon3GPP:
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
	case nasMessage.AccessTypeBoth:
		if ue.RanUe[models.AccessType__3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType__3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
		if ue.RanUe[models.AccessType_NON_3_GPP_ACCESS] != nil {
			ngap_message.SendUEContextReleaseCommand(ue.RanUe[models.AccessType_NON_3_GPP_ACCESS],
				context.UeContextReleaseUeContext, ngapType.CausePresentNas, ngapType.CauseNasPresentDeregister)
		}
	}

	ue.DeregistrationTargetAccessType = 0
	return nil
}

func HandleStatus5GMM(ue *context.AmfUe, anType models.AccessType, status5GMM *nasMessage.Status5GMM) error {
	ue.GmmLog.Info("Handle Staus 5GMM")
	if ue.MacFailed {
		return fmt.Errorf("NAS message integrity check failed")
	}

	cause := status5GMM.Cause5GMM.GetCauseValue()
	ue.GmmLog.Errorf("Error condition [Cause Value: %s]", nasMessage.Cause5GMMToString(cause))
	return nil
}
