package consumer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/util"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	Nsmf_PDUSession "github.com/free5gc/openapi/smf/PDUSession"
)

var n2sminfocon = "N2SmInfo"

type nsmfService struct {
	consumer *Consumer

	PDUSessionMu sync.RWMutex

	PDUSessionClients map[string]*Nsmf_PDUSession.APIClient
}

func (s *nsmfService) getPDUSessionClient(uri string) *Nsmf_PDUSession.APIClient {
	if uri == "" {
		return nil
	}
	s.PDUSessionMu.RLock()
	client, ok := s.PDUSessionClients[uri]
	if ok {
		s.PDUSessionMu.RUnlock()
		return client
	}

	configuration := Nsmf_PDUSession.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nsmf_PDUSession.NewAPIClient(configuration)

	s.PDUSessionMu.RUnlock()
	s.PDUSessionMu.Lock()
	defer s.PDUSessionMu.Unlock()
	s.PDUSessionClients[uri] = client
	return client
}

func (s *nsmfService) SelectSmf(
	ue *amf_context.AmfUe,
	anType models.AccessType,
	pduSessionID int32,
	snssai models.Snssai,
	dnn string,
) (*amf_context.SmContext, uint8, error) {
	var (
		smfID  string
		smfUri string
	)

	ue.GmmLog.Infof("Select SMF [snssai: %+v, dnn: %+v]", snssai, dnn)

	nrfUri := ue.ServingAMF().NrfUri // default NRF URI is pre-configured by AMF

	nsiInformation := ue.GetNsiInformationFromSnssai(anType, snssai)
	if nsiInformation == nil {
		if ue.NssfUri == "" {
			// TODO: Set a timeout of NSSF Selection or will starvation here
			for {
				searchReq := Nnrf_NFDiscovery.SearchNFInstancesRequest{}
				if err := s.consumer.SearchNssfNSSelectionInstance(ue, nrfUri, models.NrfNfManagementNfType_NSSF,
					models.NrfNfManagementNfType_AMF, &searchReq); err != nil {
					ue.GmmLog.Errorf("AMF can not select an NSSF Instance by NRF[Error: %+v]", err)
					time.Sleep(2 * time.Second)
				} else {
					break
				}
			}
		}

		response, problemDetails, err := s.consumer.NSSelectionGetForPduSession(ue, snssai)
		if err != nil {
			err = fmt.Errorf("NSSelection Get Error[%+v]", err)
			return nil, nasMessage.Cause5GMMPayloadWasNotForwarded, err
		} else if problemDetails != nil {
			err = fmt.Errorf("NSSelection Get Failed Problem[%+v]", problemDetails)
			return nil, nasMessage.Cause5GMMPayloadWasNotForwarded, err
		}
		nsiInformation = response.NsiInformation
	}

	smContext := amf_context.NewSmContext(pduSessionID)
	smContext.SetSnssai(snssai)
	smContext.SetDnn(dnn)
	smContext.SetAccessType(anType)

	if nsiInformation == nil {
		ue.GmmLog.Warnf("nsiInformation is still nil, use default NRF[%s]", nrfUri)
	} else {
		smContext.SetNsInstance(nsiInformation.NsiId)
		nrfApiUri, err := url.Parse(nsiInformation.NrfId)
		if err != nil {
			ue.GmmLog.Errorf("Parse NRF URI error, use default NRF[%s]", nrfUri)
		} else {
			nrfUri = fmt.Sprintf("%s://%s", nrfApiUri.Scheme, nrfApiUri.Host)
		}
	}

	param := Nnrf_NFDiscovery.SearchNFInstancesRequest{
		ServiceNames: []models.ServiceName{models.ServiceName_NSMF_PDUSESSION},
		Dnn:          &dnn,
		Snssais:      []models.Snssai{snssai},
	}
	if ue.PlmnId.Mcc != "" {
		param.TargetPlmnList = append(param.TargetPlmnList, ue.PlmnId)
	}
	if amf_context.GetSelf().Locality != "" {
		param.PreferredLocality = &amf_context.GetSelf().Locality
	}

	ue.GmmLog.Debugf("Search SMF from NRF[%s]", nrfUri)

	result, err := s.consumer.SendSearchNFInstances(nrfUri, models.NrfNfManagementNfType_SMF,
		models.NrfNfManagementNfType_AMF, &param)
	if err != nil {
		return nil, nasMessage.Cause5GMMPayloadWasNotForwarded, err
	}

	if len(result.NfInstances) == 0 {
		err = fmt.Errorf("DNN[%s] is not supported or not subscribed in the slice[Snssai: %+v]", dnn, snssai)
		return nil, nasMessage.Cause5GMMDNNNotSupportedOrNotSubscribedInTheSlice, err
	}

	// select the first SMF, TODO: select base on other info
	for index := range result.NfInstances {
		smfUri = util.SearchNFServiceUri(&result.NfInstances[index], models.ServiceName_NSMF_PDUSESSION,
			models.NfServiceStatus_REGISTERED)
		if smfUri != "" {
			break
		}
	}
	smContext.SetSmfID(smfID)
	smContext.SetSmfUri(smfUri)
	return smContext, 0, nil
}

func (s *nsmfService) SendCreateSmContextRequest(ue *amf_context.AmfUe, smContext *amf_context.SmContext,
	requestType *models.RequestType, nasPdu []byte) (
	smContextRef string, errorResponse *models.PostSmContextsError,
	problemDetail *models.ProblemDetails, err1 error,
) {
	smContextCreateData := s.buildCreateSmContextRequest(ue, smContext, nil)

	postSmContextsRequest := Nsmf_PDUSession.PostSmContextsRequest{
		PostSmContextsRequest: &models.PostSmContextsRequest{
			JsonData:              &smContextCreateData,
			BinaryDataN1SmMessage: nasPdu,
		},
	}

	client := s.getPDUSessionClient(smContext.SmfUri())
	if client == nil {
		return "", nil, nil, openapi.ReportError("smf not found")
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NSMF_PDUSESSION, models.NrfNfManagementNfType_SMF)
	if err != nil {
		return "", nil, nil, err
	}
	postSmContextReponse, localErr := client.SMContextsCollectionApi.
		PostSmContexts(ctx, &postSmContextsRequest)
	if localErr == nil {
		location := postSmContextReponse.Location

		parts := strings.Split(location, "/")

		if len(parts) > 0 {
			location = parts[len(parts)-1]
		}

		smContextRef = "urn:uuid:" + location
	} else {
		err1 = localErr
		switch errType := localErr.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errModel := errType.Model().(type) {
			case Nsmf_PDUSession.PostSmContextsError:
				problemDetail = &errModel.ProblemDetails
				errorResponse = &errModel.PostSmContextsError
			case error:
				err1 = errModel
			default:
				err1 = openapi.ReportError("openapi error")
			}
		case error:
			problemDetail = openapi.ProblemDetailsSystemFailure(err1.Error())
		default:
			err1 = openapi.ReportError("server no response")
		}
	}
	return smContextRef, errorResponse, problemDetail, err1
}

func (s *nsmfService) buildCreateSmContextRequest(ue *amf_context.AmfUe, smContext *amf_context.SmContext,
	requestType *models.RequestType,
) (smContextCreateData models.SmfPduSessionSmContextCreateData) {
	context := amf_context.GetSelf()
	smContextCreateData.Supi = ue.Supi
	smContextCreateData.UnauthenticatedSupi = ue.UnauthenticatedSupi
	smContextCreateData.Pei = ue.Pei
	smContextCreateData.Gpsi = ue.Gpsi
	smContextCreateData.PduSessionId = smContext.PduSessionID()
	snssai := smContext.Snssai()
	smContextCreateData.SNssai = &snssai
	smContextCreateData.Dnn = smContext.Dnn()
	smContextCreateData.ServingNfId = context.NfId
	smContextCreateData.Guami = &context.ServedGuamiList[0]
	smContextCreateData.ServingNetwork = context.ServedGuamiList[0].PlmnId
	if requestType != nil {
		smContextCreateData.RequestType = *requestType
	}
	smContextCreateData.N1SmMsg = new(models.RefToBinaryData)
	smContextCreateData.N1SmMsg.ContentId = "n1SmMsg"
	smContextCreateData.AnType = smContext.AccessType()
	if ue.RatType != "" {
		smContextCreateData.RatType = ue.RatType
	}
	smContextCreateData.UeLocation = &ue.Location
	smContextCreateData.UeTimeZone = ue.TimeZone
	smContextCreateData.SmContextStatusUri = context.GetIPv4Uri() + factory.AmfCallbackResUriPrefix + "/smContextStatus/" +
		ue.Supi + "/" + strconv.Itoa(int(smContext.PduSessionID()))

	return smContextCreateData
}

// Upadate SmContext Request
// servingNfId, smContextStatusUri, guami, servingNetwork -> amf change
// anType -> anType change
// ratType -> ratType change
// presenceInLadn -> Service Request , Xn handover, N2 handover and dnn is a ladn
// ueLocation -> the user location has changed or the user plane of the PDU session is deactivated
// upCnxState -> request the activation or the deactivation of the user plane connection of the PDU session
// hoState -> the preparation, execution or cancellation of a handover of the PDU session
// toBeSwitch -> Xn Handover to request to switch the PDU session to a new downlink N3 tunnel endpoint
// failedToBeSwitch -> indicate that the PDU session failed to be setup in the target RAN
// targetId, targetServingNfId(preparation with AMF change) -> N2 handover
// release -> duplicated PDU Session Id in subclause 5.2.2.3.11, slice not available in subclause 5.2.2.3.12
// ngApCause -> e.g. the NGAP cause for requesting to deactivate the user plane connection of the PDU session.
// 5gMmCauseValue -> AMF received a 5GMM cause code from the UE e.g 5GMM Status message in response to
// a Downlink NAS Transport message carrying 5GSM payload
// anTypeCanBeChanged

func (s *nsmfService) SendUpdateSmContextActivateUpCnxState(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, accessType models.AccessType) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.UpCnxState = models.UpCnxState_ACTIVATING
	if !amf_context.CompareUserLocation(ue.Location, smContext.UserLocation()) {
		updateData.UeLocation = &ue.Location
	}
	if smContext.AccessType() != accessType {
		updateData.AnType = smContext.AccessType()
	}
	if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
		if amf_context.InTaiList(ue.Tai, ladn.TaiList) {
			updateData.PresenceInLadn = models.PresenceState_IN_AREA
		}
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
}

func (s *nsmfService) SendUpdateSmContextDeactivateUpCnxState(ue *amf_context.AmfUe,
	smContext *amf_context.SmContext, cause amf_context.CauseAll) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.UpCnxState = models.UpCnxState_DEACTIVATED
	updateData.UeLocation = &ue.Location
	if cause.Cause != nil {
		updateData.Cause = *cause.Cause
	}
	if cause.NgapCause != nil {
		updateData.NgApCause = cause.NgapCause
	}
	if cause.Var5GmmCause != nil {
		updateData.Var5gMmCauseValue = *cause.Var5GmmCause
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
}

func (s *nsmfService) SendUpdateSmContextChangeAccessType(ue *amf_context.AmfUe,
	smContext *amf_context.SmContext, anTypeCanBeChanged bool) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.AnTypeCanBeChanged = anTypeCanBeChanged
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
}

func (s *nsmfService) SendUpdateSmContextN2Info(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, n2SmInfo []byte) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.N2SmInfoType = n2SmType
	updateData.N2SmInfo = new(models.RefToBinaryData)
	updateData.N2SmInfo.ContentId = n2sminfocon
	updateData.UeLocation = &ue.Location
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, n2SmInfo)
}

func (s *nsmfService) SendUpdateSmContextXnHandover(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, n2SmInfo []byte) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = n2sminfocon
	}
	updateData.ToBeSwitched = true
	updateData.UeLocation = &ue.Location
	if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
		if amf_context.InTaiList(ue.Tai, ladn.TaiList) {
			updateData.PresenceInLadn = models.PresenceState_IN_AREA
		} else {
			updateData.PresenceInLadn = models.PresenceState_OUT_OF_AREA
		}
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, n2SmInfo)
}

func (s *nsmfService) SendUpdateSmContextXnHandoverFailed(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, n2SmInfo []byte) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = n2sminfocon
	}
	updateData.FailedToBeSwitched = true
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, n2SmInfo)
}

func (s *nsmfService) SendUpdateSmContextN2HandoverPreparing(
	ue *amf_context.AmfUe,
	smContext *amf_context.SmContext,
	n2SmType models.N2SmInfoType,
	n2SmInfo []byte, amfid string, targetId *models.NgRanTargetId) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = n2sminfocon
	}
	updateData.HoState = models.HoState_PREPARING
	updateData.TargetId = targetId
	// amf changed in same plmn
	if amfid != "" {
		updateData.TargetServingNfId = amfid
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, n2SmInfo)
}

func (s *nsmfService) SendUpdateSmContextN2HandoverPrepared(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, n2SmType models.N2SmInfoType, n2SmInfo []byte) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	if n2SmType != "" {
		updateData.N2SmInfoType = n2SmType
		updateData.N2SmInfo = new(models.RefToBinaryData)
		updateData.N2SmInfo.ContentId = n2sminfocon
	}
	updateData.HoState = models.HoState_PREPARED
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, n2SmInfo)
}

func (s *nsmfService) SendUpdateSmContextN2HandoverComplete(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, amfid string, guami *models.Guami) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.HoState = models.HoState_COMPLETED
	if amfid != "" {
		updateData.ServingNfId = amfid
		updateData.ServingNetwork = guami.PlmnId
		updateData.Guami = guami
	}
	if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
		if amf_context.InTaiList(ue.Tai, ladn.TaiList) {
			updateData.PresenceInLadn = models.PresenceState_IN_AREA
		} else {
			updateData.PresenceInLadn = models.PresenceState_OUT_OF_AREA
		}
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
}

func (s *nsmfService) SendUpdateSmContextN2HandoverCanceled(ue *amf_context.AmfUe,
	smContext *amf_context.SmContext, cause amf_context.CauseAll) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	// nolint openapi/model misspelling
	updateData.HoState = models.HoState_CANCELLED
	if cause.Cause != nil {
		updateData.Cause = *cause.Cause
	}
	if cause.NgapCause != nil {
		updateData.NgApCause = cause.NgapCause
	}
	if cause.Var5GmmCause != nil {
		updateData.Var5gMmCauseValue = *cause.Var5GmmCause
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
}

func (s *nsmfService) SendUpdateSmContextHandoverBetweenAccessType(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, targetAccessType models.AccessType, n1SmMsg []byte) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.AnType = targetAccessType
	if n1SmMsg != nil {
		updateData.N1SmMsg = new(models.RefToBinaryData)
		updateData.N1SmMsg.ContentId = "N1Msg"
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, n1SmMsg, nil)
}

func (s *nsmfService) SendUpdateSmContextHandoverBetweenAMF(
	ue *amf_context.AmfUe, smContext *amf_context.SmContext, amfid string, guami *models.Guami, activate bool) (
	*models.UpdateSmContextResponse200, *models.UpdateSmContextResponse400, *models.ProblemDetails, error,
) {
	updateData := models.SmfPduSessionSmContextUpdateData{}
	updateData.ServingNfId = amfid
	updateData.ServingNetwork = guami.PlmnId
	updateData.Guami = guami
	if activate {
		updateData.UpCnxState = models.UpCnxState_ACTIVATING
		if !amf_context.CompareUserLocation(ue.Location, smContext.UserLocation()) {
			updateData.UeLocation = &ue.Location
		}
		if ladn, ok := ue.ServingAMF().LadnPool[smContext.Dnn()]; ok {
			if amf_context.InTaiList(ue.Tai, ladn.TaiList) {
				updateData.PresenceInLadn = models.PresenceState_IN_AREA
			}
		}
	}
	return s.consumer.SendUpdateSmContextRequest(smContext, &updateData, nil, nil)
}

func (s *nsmfService) SendUpdateSmContextRequest(smContext *amf_context.SmContext,
	updateData *models.SmfPduSessionSmContextUpdateData, n1Msg []byte, n2Info []byte) (
	response *models.UpdateSmContextResponse200, errorResponse *models.UpdateSmContextResponse400,
	problemDetail *models.ProblemDetails, err1 error,
) {
	client := s.getPDUSessionClient(smContext.SmfUri())
	if client == nil {
		return nil, nil, nil, openapi.ReportError("smf not found")
	}

	smCtxRef := smContext.SmContextRef()
	updateSmContextRequest := Nsmf_PDUSession.UpdateSmContextRequest{
		SmContextRef: &smCtxRef,
		UpdateSmContextRequest: &models.UpdateSmContextRequest{
			JsonData:                  updateData,
			BinaryDataN1SmMessage:     n1Msg,
			BinaryDataN2SmInformation: n2Info,
		},
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NSMF_PDUSESSION, models.NrfNfManagementNfType_SMF)
	if err != nil {
		return nil, nil, nil, err
	}
	updateSmContextReponse, localErr := client.IndividualSMContextApi.
		UpdateSmContext(ctx, &updateSmContextRequest)
	if localErr == nil {
		response = &updateSmContextReponse.UpdateSmContextResponse200
	} else {
		err1 = localErr
		switch errType := localErr.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errModel := errType.Model().(type) {
			case Nsmf_PDUSession.UpdateSmContextError:
				problemDetail = &errModel.ProblemDetails
				errorResponse = &errModel.UpdateSmContextResponse400
			case error:
				err1 = errModel
			default:
				err1 = openapi.ReportError("openapi error")
			}
		case error:
			problemDetail = openapi.ProblemDetailsSystemFailure(err1.Error())
		default:
			err1 = openapi.ReportError("server no response")
		}
	}
	return response, errorResponse, problemDetail, err1
}

// Release SmContext Request

func (s *nsmfService) SendReleaseSmContextRequest(ue *amf_context.AmfUe, smContext *amf_context.SmContext,
	cause *amf_context.CauseAll, n2SmInfoType models.N2SmInfoType,
	n2Info []byte,
) (detail *models.ProblemDetails, err error) {
	client := s.getPDUSessionClient(smContext.SmfUri())
	if client == nil {
		return nil, openapi.ReportError("smf not found")
	}

	releaseData := s.buildReleaseSmContextRequest(ue, cause, n2SmInfoType, n2Info)

	smCtxRef := smContext.SmContextRef()
	releaseSmContextRequest := Nsmf_PDUSession.ReleaseSmContextRequest{
		SmContextRef: &smCtxRef,
		ReleaseSmContextRequest: &models.ReleaseSmContextRequest{
			JsonData: &releaseData,
		},
	}

	ctx, _, err := amf_context.GetSelf().GetTokenCtx(models.ServiceName_NSMF_PDUSESSION, models.NrfNfManagementNfType_SMF)
	if err != nil {
		return nil, err
	}
	_, localErr := client.IndividualSMContextApi.ReleaseSmContext(
		ctx, &releaseSmContextRequest)

	if localErr == nil {
		ue.SmContextList.Delete(smContext.PduSessionID())
	} else {
		err = localErr
		switch apiErr := localErr.(type) {
		// API error
		case openapi.GenericOpenAPIError:
			switch errorModel := apiErr.Model().(type) {
			case Nsmf_PDUSession.ReleaseSmContextError:
				detail = &errorModel.ProblemDetails
			case error:
				detail = openapi.ProblemDetailsSystemFailure(errorModel.Error())
			default:
				err = openapi.ReportError("openapi error")
			}
		case error:
			detail = openapi.ProblemDetailsSystemFailure(apiErr.Error())
		default:
			err = openapi.ReportError("openapi error")
		}
	}
	return detail, err
}

func (s *nsmfService) buildReleaseSmContextRequest(
	ue *amf_context.AmfUe, cause *amf_context.CauseAll, n2SmInfoType models.N2SmInfoType, n2Info []byte) (
	releaseData models.SmfPduSessionSmContextReleaseData,
) {
	if cause != nil {
		if cause.Cause != nil {
			releaseData.Cause = *cause.Cause
		}
		if cause.NgapCause != nil {
			releaseData.NgApCause = cause.NgapCause
		}
		if cause.Var5GmmCause != nil {
			releaseData.Var5gMmCauseValue = *cause.Var5GmmCause
		}
	}
	if ue.TimeZone != "" {
		releaseData.UeTimeZone = ue.TimeZone
	}
	if n2Info != nil {
		releaseData.N2SmInfoType = n2SmInfoType
		releaseData.N2SmInfo = &models.RefToBinaryData{
			ContentId: n2sminfocon,
		}
	}
	// TODO: other param(ueLocation...)
	return
}
