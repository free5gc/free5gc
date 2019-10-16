package udm_producer

import (
	"context"
	Nudr "free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
	"strconv"
)

func HandleGetAmData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	accessAndMobilitySubscriptionDataResp, res, err := clientAPI.AccessAndMobilitySubscriptionDataDocumentApi.QueryAmData(context.Background(), supi, plmnID, nil)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, accessAndMobilitySubscriptionDataResp)
}

func HandleGetIdTranslationResult(httpChannel chan udm_message.HandlerResponseMessage, gpsi string) {

	clientAPI := createUDMClientToUDR(gpsi, false)
	idTranslationResultResp, res, err := clientAPI.QueryIdentityDataBySUPIOrGPSIDocumentApi.GetIdentityData(context.Background(), gpsi, nil)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, idTranslationResultResp)
}

func HandleGetSupi(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	var subscriptionDataSetsReq models.SubscriptionDataSets
	var ueContextInSmfDataResp models.UeContextInSmfData
	pduSessionMap := make(map[string]models.PduSession)
	var pgwInfoArray []models.PgwInfo

	var queryAmDataParamOpts Nudr.QueryAmDataParamOpts
	var querySmfSelectDataParamOpts Nudr.QuerySmfSelectDataParamOpts

	amData, res, err1 := clientAPI.AccessAndMobilitySubscriptionDataDocumentApi.QueryAmData(context.Background(),
		supi, plmnID, &queryAmDataParamOpts)
	subscriptionDataSetsReq.AmData = &amData
	if err1 != nil {
		logger.SdmLog.Panic(err1.Error())
	}

	smfSelData, res, err2 := clientAPI.SMFSelectionSubscriptionDataDocumentApi.QuerySmfSelectData(context.Background(),
		supi, plmnID, &querySmfSelectDataParamOpts)
	subscriptionDataSetsReq.SmfSelData = &smfSelData
	if err2 != nil {
		logger.SdmLog.Panic(err2.Error())
	}

	traceData, res, err3 := clientAPI.TraceDataDocumentApi.QueryTraceData(context.Background(), supi, plmnID, nil)
	subscriptionDataSetsReq.TraceData = &traceData
	if err3 != nil {
		logger.SdmLog.Panic(err3.Error())
	}

	sessionManagementSubscriptionData, res, err4 := clientAPI.SessionManagementSubscriptionDataApi.QuerySmData(context.Background(), supi, plmnID, nil)
	subscriptionDataSetsReq.SmData = sessionManagementSubscriptionData
	if err4 != nil {
		logger.SdmLog.Panic(err4.Error())
	}

	pdusess, res, err := clientAPI.SMFRegistrationsCollectionApi.QuerySmfRegList(context.Background(), supi, nil)
	array := pdusess

	for _, element := range array {
		var pduSession models.PduSession
		pduSession.Dnn = element.Dnn
		pduSession.SmfInstanceId = element.SmfInstanceId
		pduSession.PlmnId = element.PlmnId
		pduSessionMap[strconv.Itoa(int(element.PduSessionId))] = pduSession
	}
	ueContextInSmfDataResp.PduSessions = pduSessionMap

	for _, element := range array {
		var pgwInfo models.PgwInfo
		pgwInfo.Dnn = element.Dnn
		pgwInfo.PgwFqdn = element.PgwFqdn
		pgwInfo.PlmnId = element.PlmnId
		pgwInfoArray = append(pgwInfoArray, pgwInfo)
	}
	ueContextInSmfDataResp.PgwInfo = pgwInfoArray

	subscriptionDataSetsReq.UecSmfData = &ueContextInSmfDataResp

	if res.StatusCode == http.StatusOK {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, subscriptionDataSetsReq)
	} else {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, err.(common.GenericOpenAPIError).Model().(models.ProblemDetails))

	}
}
func HandleGetSharedData(httpChannel chan udm_message.HandlerResponseMessage, sharedDataIds []string) {

	clientAPI := createUDMClientToUDR("", true)
	var getSharedDataParamOpts Nudr.GetSharedDataParamOpts

	sharedDataResp, res, err := clientAPI.RetrievalOfSharedDataApi.GetSharedData(context.Background(), sharedDataIds,
		&getSharedDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sharedDataResp)
}

func HandleGetSmData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	var querySmDataParamOpts Nudr.QuerySmDataParamOpts

	sessionManagementSubscriptionDataResp, res, err := clientAPI.SessionManagementSubscriptionDataApi.QuerySmData(context.Background(),
		supi, plmnID, &querySmDataParamOpts)

	if res.StatusCode == http.StatusNotFound {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, err.(common.GenericOpenAPIError).Model().(models.ProblemDetails))
	} else if res.StatusCode == http.StatusBadRequest {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusBadRequest, err.(common.GenericOpenAPIError).Model().(models.ProblemDetails))
	} else if res.StatusCode == http.StatusInternalServerError {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusInternalServerError, err.(common.GenericOpenAPIError).Model().(models.ProblemDetails))
	} else {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sessionManagementSubscriptionDataResp)
	}

}

func HandleGetNssai(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	var nssaiResp models.Nssai

	accessAndMobilitySubscriptionDataResp, res, err := clientAPI.AccessAndMobilitySubscriptionDataDocumentApi.QueryAmData(context.Background(),
		supi, plmnID, nil)
	nssaiResp = *accessAndMobilitySubscriptionDataResp.Nssai
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, nssaiResp)
}

func HandleGetSmfSelectData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	var querySmfSelectDataParamOpts Nudr.QuerySmfSelectDataParamOpts
	smfSelectionSubscriptionDataResp, res, err := clientAPI.SMFSelectionSubscriptionDataDocumentApi.QuerySmfSelectData(context.Background(),
		supi, plmnID, &querySmfSelectDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, smfSelectionSubscriptionDataResp)
}

func HandleSubscribeToSharedData(httpChannel chan udm_message.HandlerResponseMessage, sdmSubscription models.SdmSubscription) {

	// TODO
	clientAPI := createUDMClientToUDR("", true)
	sdmSubscriptionResp, res, err := clientAPI.SDMSubscriptionsCollectionApi.CreateSdmSubscriptions(context.Background(), "===", sdmSubscription)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusCreated, sdmSubscriptionResp)
}

func HandleSubscribe(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string, sdmSubscription models.SdmSubscription) {

	clientAPI := createUDMClientToUDR(supi, false)
	sdmSubscriptionResp, res, err := clientAPI.SDMSubscriptionsCollectionApi.CreateSdmSubscriptions(context.Background(), supi, sdmSubscription)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusCreated, sdmSubscriptionResp)
}

func HandleUnsubscribeForSharedData(httpChannel chan udm_message.HandlerResponseMessage, subscriptionID string) {

	// TODO
	clientAPI := createUDMClientToUDR("", true)
	res, err := clientAPI.SDMSubscriptionDocumentApi.RemovesdmSubscriptions(context.Background(), "====", subscriptionID)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}

func HandleUnsubscribe(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	res, err := clientAPI.SDMSubscriptionDocumentApi.RemovesdmSubscriptions(context.Background(), supi, subscriptionID)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}

func HandleModify(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string, sdmSubscription models.SdmSubscription) {

	clientAPI := createUDMClientToUDR(supi, false)
	res, err := clientAPI.SDMSubscriptionDocumentApi.Updatesdmsubscriptions(context.Background(), supi, subscriptionID, nil)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sdmSubscription)
}

func HandleModifyForSharedData(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string) {

	// TODO
	var sdmSubscription models.SdmSubscription
	clientAPI := createUDMClientToUDR(supi, false)
	res, err := clientAPI.SDMSubscriptionDocumentApi.Updatesdmsubscriptions(context.Background(), supi, subscriptionID, nil)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sdmSubscription)
}

func HandleGetTraceData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	traceDataResp, res, err := clientAPI.TraceDataDocumentApi.QueryTraceData(context.Background(), supi, plmnID, nil)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, traceDataResp)
}

func HandleGetUeContextInSmfData(httpChannel chan udm_message.HandlerResponseMessage, supi string) {

	clientAPI := createUDMClientToUDR(supi, false)
	var ueContextInSmfData models.UeContextInSmfData
	pduSessionMap := make(map[string]models.PduSession)
	var pgwInfoArray []models.PgwInfo

	pdusess, res, err := clientAPI.SMFRegistrationsCollectionApi.QuerySmfRegList(context.Background(), supi, nil)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	array := pdusess

	for _, element := range array {
		var pduSession models.PduSession
		pduSession.Dnn = element.Dnn
		pduSession.SmfInstanceId = element.SmfInstanceId
		pduSession.PlmnId = element.PlmnId
		pduSessionMap[strconv.Itoa(int(element.PduSessionId))] = pduSession
	}
	ueContextInSmfData.PduSessions = pduSessionMap

	for _, element := range array {
		var pgwInfo models.PgwInfo
		pgwInfo.Dnn = element.Dnn
		pgwInfo.PgwFqdn = element.PgwFqdn
		pgwInfo.PlmnId = element.PlmnId
		pgwInfoArray = append(pgwInfoArray, pgwInfo)
	}
	ueContextInSmfData.PgwInfo = pgwInfoArray

	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, ueContextInSmfData)
}
