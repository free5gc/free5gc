package udm_producer

import (
	"context"
	"fmt"
	"free5gc/lib/Nudm_SubscriberDataManagement"
	Nudr "free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
	"strconv"

	"github.com/antihax/optional"
)

func HandleGetAmData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string, supportedFeatures string) {
	var queryAmDataParamOpts Nudr.QueryAmDataParamOpts
	queryAmDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)

	clientAPI := createUDMClientToUDR(supi, false)
	accessAndMobilitySubscriptionDataResp, res, err := clientAPI.AccessAndMobilitySubscriptionDataDocumentApi.QueryAmData(context.Background(),
		supi, plmnID, &queryAmDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.AccessAndMobilitySubscriptionData = &accessAndMobilitySubscriptionDataResp
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *udmUe.AccessAndMobilitySubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetIdTranslationResult(httpChannel chan udm_message.HandlerResponseMessage, gpsi string) {

	var idTranslationResult models.IdTranslationResult
	var getIdentityDataParamOpts Nudr.GetIdentityDataParamOpts
	clientAPI := createUDMClientToUDR(gpsi, false)
	idTranslationResultResp, res, err := clientAPI.QueryIdentityDataBySUPIOrGPSIDocumentApi.GetIdentityData(context.Background(), gpsi, &getIdentityDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		idList := udm_context.UDM_Self().GpsiSupiList
		idList = idTranslationResultResp
		if idList.SupiList != nil {
			idTranslationResult.Supi = udm_context.GetCorrespondingSupi(idList) // GetCorrespondingSupi get corresponding Supi(here IMSI) matching the given Gpsi from the queried SUPI list from UDR
			idTranslationResult.Gpsi = gpsi
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, idTranslationResult)
		} else {
			var problemDetail models.ProblemDetails
			problemDetail.Cause = "USER_NOT_FOUND" // SupiList are empty
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetail)
		}
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}

}

func HandleGetSupi(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string, dataSetNames []string, supportedFeatures string) {

	clientAPI := createUDMClientToUDR(supi, false)
	var subscriptionDataSets, subsDataSetBody models.SubscriptionDataSets
	var ueContextInSmfDataResp models.UeContextInSmfData
	pduSessionMap := make(map[string]models.PduSession)
	var pgwInfoArray []models.PgwInfo

	var queryAmDataParamOpts Nudr.QueryAmDataParamOpts
	queryAmDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	var querySmfSelectDataParamOpts Nudr.QuerySmfSelectDataParamOpts
	var queryTraceDataParamOpts Nudr.QueryTraceDataParamOpts
	var querySmDataParamOpts Nudr.QuerySmDataParamOpts

	queryAmDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	querySmfSelectDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	udm_context.CreateSubsDataSetsForUe(supi, subsDataSetBody)

	var body models.AccessAndMobilitySubscriptionData
	udm_context.CreateAccessMobilitySubsDataForUe(supi, body)
	amData, res1, err1 := clientAPI.AccessAndMobilitySubscriptionDataDocumentApi.QueryAmData(context.Background(), supi, plmnID, &queryAmDataParamOpts)
	if err1 != nil {
		var problemDetails models.ProblemDetails
		if res1 == nil {
			fmt.Println(err1.Error())
		} else if err1.Error() != res1.Status {
			fmt.Println(err1.Error())
		} else {
			problemDetails.Cause = err1.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res1.StatusCode, problemDetails)
		}
		return
	}
	if res1.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.AccessAndMobilitySubscriptionData = &amData
		subscriptionDataSets.AmData = &amData
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		fmt.Printf(problemDetails.Cause)
	}

	var smfSelSubsbody models.SmfSelectionSubscriptionData
	udm_context.CreateSmfSelectionSubsDatadforUe(supi, smfSelSubsbody)
	smfSelData, res2, err2 := clientAPI.SMFSelectionSubscriptionDataDocumentApi.QuerySmfSelectData(context.Background(),
		supi, plmnID, &querySmfSelectDataParamOpts)
	if err2 != nil {
		var problemDetails models.ProblemDetails
		if res2 == nil {
			fmt.Println(err2.Error())
		} else if err2.Error() != res2.Status {
			fmt.Println(err2.Error())
		} else {
			problemDetails.Cause = err2.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res2.StatusCode, problemDetails)
		}
		return
	}
	if res2.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.SmfSelSubsData = &smfSelData
		subscriptionDataSets.SmfSelData = &smfSelData
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		fmt.Printf(problemDetails.Cause)
	}

	var TraceDatabody models.TraceData
	udm_context.CreateTraceDataforUe(supi, TraceDatabody)
	traceData, res3, err3 := clientAPI.TraceDataDocumentApi.QueryTraceData(context.Background(), supi, plmnID, &queryTraceDataParamOpts)
	if err3 != nil {
		var problemDetails models.ProblemDetails
		if res3 == nil {
			fmt.Println(err3.Error())
		} else if err3.Error() != res3.Status {
			fmt.Println(err3.Error())
		} else {
			problemDetails.Cause = err3.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res3.StatusCode, problemDetails)
		}
		return
	}
	if res3.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.TraceData = &traceData
		udmUe.TraceDataResponse.TraceData = &traceData
		subscriptionDataSets.TraceData = &traceData
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		fmt.Printf(problemDetails.Cause)
	}

	sessionManagementSubscriptionData, res4, err4 := clientAPI.SessionManagementSubscriptionDataApi.QuerySmData(context.Background(), supi, plmnID, &querySmDataParamOpts)
	if err4 != nil {
		var problemDetails models.ProblemDetails
		if res4 == nil {
			fmt.Println(err4.Error())
		} else if err4.Error() != res4.Status {
			fmt.Println(err4.Error())
		} else {
			problemDetails.Cause = err4.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res4.StatusCode, problemDetails)
		}
		return
	}
	if res4.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		smData, _, _, _ := udm_context.ManageSmData(sessionManagementSubscriptionData, "", "")
		udmUe.SessionManagementSubsData = smData
		subscriptionDataSets.SmData = sessionManagementSubscriptionData
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		fmt.Printf(problemDetails.Cause)
	}

	var UeContextInSmfbody models.UeContextInSmfData
	var querySmfRegListParamOpts Nudr.QuerySmfRegListParamOpts
	querySmfRegListParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	udm_context.CreateUeContextInSmfDataforUe(supi, UeContextInSmfbody)
	pdusess, res, err := clientAPI.SMFRegistrationsCollectionApi.QuerySmfRegList(context.Background(), supi, &querySmfRegListParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	for _, element := range pdusess {
		var pduSession models.PduSession
		pduSession.Dnn = element.Dnn
		pduSession.SmfInstanceId = element.SmfInstanceId
		pduSession.PlmnId = element.PlmnId
		pduSessionMap[strconv.Itoa(int(element.PduSessionId))] = pduSession
	}
	ueContextInSmfDataResp.PduSessions = pduSessionMap

	for _, element := range pdusess {
		var pgwInfo models.PgwInfo
		pgwInfo.Dnn = element.Dnn
		pgwInfo.PgwFqdn = element.PgwFqdn
		pgwInfo.PlmnId = element.PlmnId
		pgwInfoArray = append(pgwInfoArray, pgwInfo)
	}
	ueContextInSmfDataResp.PgwInfo = pgwInfoArray

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.UeCtxtInSmfData = &ueContextInSmfDataResp
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		fmt.Printf(problemDetails.Cause)
	}

	if (res.StatusCode == http.StatusOK) && (res1.StatusCode == http.StatusOK) && (res2.StatusCode == http.StatusOK) && (res3.StatusCode == http.StatusOK) && (res4.StatusCode == http.StatusOK) {
		subscriptionDataSets.UecSmfData = &ueContextInSmfDataResp
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, subscriptionDataSets)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetSharedData(httpChannel chan udm_message.HandlerResponseMessage, sharedDataIds []string, supportedFeatures string) {

	clientAPI := createUDMClientToUDR("", true)
	var getSharedDataParamOpts Nudr.GetSharedDataParamOpts
	getSharedDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)

	sharedDataResp, res, err := clientAPI.RetrievalOfSharedDataApi.GetSharedData(context.Background(), sharedDataIds,
		&getSharedDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udm_context.UDM_Self().SharedSubsDataMap = udm_context.MappingSharedData(sharedDataResp)
		sharedData := udm_context.ObtainRequiredSharedData(sharedDataIds, sharedDataResp)
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sharedData)
	} else {
		var problemDetail models.ProblemDetails
		problemDetail.Cause = "DATA_NOT_FOUND"
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetail)
	}
}

func HandleGetSmData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string, Dnn string, Snssai string, supportedFeatures string) {
	logger.Handlelog.Infof("HandleGetSmData SUPI[%s] PLMNID[%s] DNN[%s] SNssai[%s]", supi, plmnID, Dnn, Snssai)

	clientAPI := createUDMClientToUDR(supi, false)
	var querySmDataParamOpts Nudr.QuerySmDataParamOpts
	querySmDataParamOpts.SingleNssai = optional.NewInterface(Snssai)

	sessionManagementSubscriptionDataResp, res, err := clientAPI.SessionManagementSubscriptionDataApi.QuerySmData(context.Background(),
		supi, plmnID, &querySmDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		var snssaikey string
		var AllDnnConfigsbyDnn []models.DnnConfiguration
		var AllDnns []map[string]models.DnnConfiguration
		udmUe.SessionManagementSubsData, snssaikey, AllDnnConfigsbyDnn, AllDnns = udm_context.ManageSmData(sessionManagementSubscriptionDataResp, Snssai, Dnn)

		var rspSMSubDataList = make([]models.SessionManagementSubscriptionData, 0, 4)
		for _, eachSMSubData := range udmUe.SessionManagementSubsData {
			rspSMSubDataList = append(rspSMSubDataList, eachSMSubData)
		}

		switch {
		case Snssai == "" && Dnn == "":
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, AllDnns)
		case Snssai != "" && Dnn == "":
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, udmUe.SessionManagementSubsData[snssaikey].DnnConfigurations)
		case Snssai == "" && Dnn != "":
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, AllDnnConfigsbyDnn)
		case Snssai != "" && Dnn != "":
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, rspSMSubDataList)
		default:
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, udmUe.SessionManagementSubsData)
		}
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetNssai(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string, supportedFeatures string) {

	var queryAmDataParamOpts Nudr.QueryAmDataParamOpts
	queryAmDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	var nssaiResp models.Nssai
	clientAPI := createUDMClientToUDR(supi, false)

	accessAndMobilitySubscriptionDataResp, res, err := clientAPI.AccessAndMobilitySubscriptionDataDocumentApi.QueryAmData(context.Background(),
		supi, plmnID, &queryAmDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}
	nssaiResp = *accessAndMobilitySubscriptionDataResp.Nssai

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.Nssai = &nssaiResp
		if plmnID != "" {
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, udmUe.Nssai)
		} else {
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, udmUe.Nssai)
		}
	} else {
		var problemDetail models.ProblemDetails
		problemDetail.Cause = "DATA_NOT_FOUND"
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetail)
	}
}

func HandleGetSmfSelectData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string, supportedFeatures string) {

	var querySmfSelectDataParamOpts Nudr.QuerySmfSelectDataParamOpts
	querySmfSelectDataParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	var body models.SmfSelectionSubscriptionData
	clientAPI := createUDMClientToUDR(supi, false)
	udm_context.CreateSmfSelectionSubsDatadforUe(supi, body)

	smfSelectionSubscriptionDataResp, res, err := clientAPI.SMFSelectionSubscriptionDataDocumentApi.QuerySmfSelectData(context.Background(),
		supi, plmnID, &querySmfSelectDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.SmfSelSubsData = &smfSelectionSubscriptionDataResp
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *udmUe.SmfSelSubsData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleSubscribeToSharedData(httpChannel chan udm_message.HandlerResponseMessage, sdmSubscription models.SdmSubscription) {
	var body *models.SdmSubscription
	cfg := Nudm_SubscriberDataManagement.NewConfiguration()
	udmClientAPI := Nudm_SubscriberDataManagement.NewAPIClient(cfg)
	udm_context.CreateSubstoNotifSharedData(sdmSubscription.SubscriptionId, body)

	sdmSubscriptionResp, res, err := udmClientAPI.SubscriptionCreationForSharedDataApi.SubscribeToSharedData(context.Background(), sdmSubscription)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
			fmt.Println(res.Status)
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusCreated {
		h := make(http.Header)
		udmUe := udm_context.CreateUdmUe(sdmSubscriptionResp.SubscriptionId)
		udmUe.SubscribeToNotifSharedDataChange = &sdmSubscriptionResp
		h.Set("Location", udmUe.GetLocationURI2(udm_context.LocationUriSharedDataSubscription, "supi"))
		udm_message.SendHttpResponseMessage(httpChannel, h, http.StatusCreated, *udmUe.SubscribeToNotifSharedDataChange)
	} else if res.StatusCode == http.StatusNotFound {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "UNSUPPORTED_RESOURCE_URI"
		problemDetails.Status = 501
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotImplemented, problemDetails)
	}
}

func HandleSubscribe(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string, sdmSubscription models.SdmSubscription) {

	var body *models.SdmSubscription
	clientAPI := createUDMClientToUDR(supi, false)
	udm_context.CreateSubscriptiontoNotifChange(subscriptionID, body)

	sdmSubscriptionResp, res, err := clientAPI.SDMSubscriptionsCollectionApi.CreateSdmSubscriptions(context.Background(), supi, sdmSubscription)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusCreated {
		h := make(http.Header)
		udmUe := udm_context.CreateUdmUe(subscriptionID)
		udmUe.SubscribeToNotifChange = &sdmSubscriptionResp
		h.Set("Location", udmUe.GetLocationURI2(udm_context.LocationUriSdmSubscription, supi))
		udm_message.SendHttpResponseMessage(httpChannel, h, http.StatusCreated, *udmUe.SubscribeToNotifChange)

	} else if res.StatusCode == http.StatusNotFound {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "UNSUPPORTED_RESOURCE_URI"
		problemDetails.Status = 501
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotImplemented, problemDetails)
	}
}

func HandleUnsubscribeForSharedData(httpChannel chan udm_message.HandlerResponseMessage, subscriptionID string) {

	cfg := Nudm_SubscriberDataManagement.NewConfiguration()
	udmClientAPI := Nudm_SubscriberDataManagement.NewAPIClient(cfg)

	res, err := udmClientAPI.SubscriptionDeletionForSharedDataApi.UnsubscribeForSharedData(context.Background(), subscriptionID)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}
	if res.StatusCode == http.StatusNoContent {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleUnsubscribe(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string) {

	clientAPI := createUDMClientToUDR(supi, false)
	res, err := clientAPI.SDMSubscriptionDocumentApi.RemovesdmSubscriptions(context.Background(), "====", subscriptionID)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}
	if res.StatusCode == http.StatusNoContent {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleModify(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string, sdmSubsModification models.SdmSubsModification) {

	clientAPI := createUDMClientToUDR(supi, false)
	sdmSubscription := models.SdmSubscription{}
	body := Nudr.UpdatesdmsubscriptionsParamOpts{
		SdmSubscription: optional.NewInterface(sdmSubscription),
	}
	res, err := clientAPI.SDMSubscriptionDocumentApi.Updatesdmsubscriptions(context.Background(), supi, subscriptionID, &body)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sdmSubscription)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleModifyForSharedData(httpChannel chan udm_message.HandlerResponseMessage, supi string, subscriptionID string, sdmSubsModification models.SdmSubsModification) {

	clientAPI := createUDMClientToUDR(supi, false)
	var sdmSubscription models.SdmSubscription
	sdmSubs := models.SdmSubscription{}
	body := Nudr.UpdatesdmsubscriptionsParamOpts{
		SdmSubscription: optional.NewInterface(sdmSubs),
	}

	res, err := clientAPI.SDMSubscriptionDocumentApi.Updatesdmsubscriptions(context.Background(), supi, subscriptionID, &body)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, sdmSubscription)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetTraceData(httpChannel chan udm_message.HandlerResponseMessage, supi string, plmnID string) {

	var body models.TraceData
	var queryTraceDataParamOpts Nudr.QueryTraceDataParamOpts
	clientAPI := createUDMClientToUDR(supi, false)
	udm_context.CreateTraceDataforUe(supi, body)

	traceDataRes, res, err := clientAPI.TraceDataDocumentApi.QueryTraceData(context.Background(), supi, plmnID, &queryTraceDataParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.TraceData = &traceDataRes
		udmUe.TraceDataResponse.TraceData = &traceDataRes

		if plmnID != "" {
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *udmUe.TraceDataResponse.TraceData)
		} else {
			udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *udmUe.TraceDataResponse.TraceData) // If "plmn-id" is not included, UDM shall return the Trace Data for the SUPI associated to the HPLMN
		}
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}

}

func HandleGetUeContextInSmfData(httpChannel chan udm_message.HandlerResponseMessage, supi string, supportedFeatures string) {

	var body models.UeContextInSmfData
	var ueContextInSmfData models.UeContextInSmfData
	var pgwInfoArray []models.PgwInfo
	var querySmfRegListParamOpts Nudr.QuerySmfRegListParamOpts
	querySmfRegListParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)
	clientAPI := createUDMClientToUDR(supi, false)
	pduSessionMap := make(map[string]models.PduSession)
	udm_context.CreateUeContextInSmfDataforUe(supi, body)

	pdusess, res, err := clientAPI.SMFRegistrationsCollectionApi.QuerySmfRegList(context.Background(), supi, &querySmfRegListParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		if res == nil {
			fmt.Println(err.Error())
		} else if err.Error() != res.Status {
			fmt.Println(err.Error())
		} else {
			problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
			udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		}
		return
	}

	for _, element := range pdusess {
		var pduSession models.PduSession
		pduSession.Dnn = element.Dnn
		pduSession.SmfInstanceId = element.SmfInstanceId
		pduSession.PlmnId = element.PlmnId
		pduSessionMap[strconv.Itoa(int(element.PduSessionId))] = pduSession
	}
	ueContextInSmfData.PduSessions = pduSessionMap

	for _, element := range pdusess {
		var pgwInfo models.PgwInfo
		pgwInfo.Dnn = element.Dnn
		pgwInfo.PgwFqdn = element.PgwFqdn
		pgwInfo.PlmnId = element.PlmnId
		pgwInfoArray = append(pgwInfoArray, pgwInfo)
	}
	ueContextInSmfData.PgwInfo = pgwInfoArray

	if res.StatusCode == http.StatusOK {
		udmUe := udm_context.CreateUdmUe(supi)
		udmUe.UeCtxtInSmfData = &ueContextInSmfData
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *udmUe.UeCtxtInSmfData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}
