package udm_producer

import (
	"context"
	"fmt"
	"free5gc/lib/Nudr_DataRepository"
	Nudr "free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/factory"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_consumer"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
	"strconv"
	"strings"

	"github.com/antihax/optional"
)

func createUDMClientToUDR(id string, nonUe bool) *Nudr_DataRepository.APIClient {
	var addr string
	if !nonUe {
		addr = getUdrUri(id)
	}
	if addr == "" {
		// dafault
		if !nonUe {
			logger.Handlelog.Warnf("Use default UDR Uri bacause ID[%s] does not match any UDR", id)
		}
		config := factory.UdmConfig
		udrclient := config.Configuration.Udrclient
		addr = fmt.Sprintf("%s://%s:%d", udrclient.Scheme, udrclient.Ipv4Adrr, udrclient.Port)
	}
	cfg := Nudr.NewConfiguration()
	cfg.SetBasePath(addr)
	clientAPI := Nudr.NewAPIClient(cfg)
	return clientAPI
}

func getUdrUri(id string) string {
	// supi
	if strings.Contains(id, "imsi") || strings.Contains(id, "nai") {
		udmUe := udm_context.UDM_Self().UdmUePool[id]
		if udmUe != nil {
			if udmUe.UdrUri == "" {
				udmUe.UdrUri = udm_consumer.SendNFIntancesUDR(id, udm_consumer.NFDiscoveryToUDRParamSupi)
			}
			return udmUe.UdrUri
		} else {
			udmUe = udm_context.CreateUdmUe(id)
			udmUe.UdrUri = udm_consumer.SendNFIntancesUDR(id, udm_consumer.NFDiscoveryToUDRParamSupi)
			return udmUe.UdrUri
		}
	} else if strings.Contains(id, "pei") {
		for _, udmUe := range udm_context.UDM_Self().UdmUePool {
			if udmUe.Amf3GppAccessRegistration != nil && udmUe.Amf3GppAccessRegistration.Pei == id {
				if udmUe.UdrUri != "" {
					udmUe.UdrUri = udm_consumer.SendNFIntancesUDR(udmUe.Supi, udm_consumer.NFDiscoveryToUDRParamSupi)
				}
				return udmUe.UdrUri
			} else if udmUe.AmfNon3GppAccessRegistration != nil && udmUe.AmfNon3GppAccessRegistration.Pei == id {
				if udmUe.UdrUri != "" {
					udmUe.UdrUri = udm_consumer.SendNFIntancesUDR(udmUe.Supi, udm_consumer.NFDiscoveryToUDRParamSupi)
				}
				return udmUe.UdrUri
			}
		}
	} else if strings.Contains(id, "extgroupid") {
		// extra group id
		return udm_consumer.SendNFIntancesUDR(id, udm_consumer.NFDiscoveryToUDRParamExtGroupId)
	} else if strings.Contains(id, "msisdn") || strings.Contains(id, "extid") {
		// gpsi
		return udm_consumer.SendNFIntancesUDR(id, udm_consumer.NFDiscoveryToUDRParamGpsi)
	}
	return ""
}

func HandleGetAmf3gppAccess(respChan chan udm_message.HandlerResponseMessage, ueID string, supportedFeatures string) {
	var queryAmfContext3gppParamOpts Nudr_DataRepository.QueryAmfContext3gppParamOpts
	queryAmfContext3gppParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)

	clientAPI := createUDMClientToUDR(ueID, false)
	amf3GppAccessRegistration, resp, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.QueryAmfContext3gpp(context.Background(), ueID, &queryAmfContext3gppParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amf3GppAccessRegistration)
}

func HandleGetAmfNon3gppAccess(respChan chan udm_message.HandlerResponseMessage, ueID string, supportedFeatures string) {
	var queryAmfContextNon3gppParamOpts Nudr_DataRepository.QueryAmfContextNon3gppParamOpts
	queryAmfContextNon3gppParamOpts.SupportedFeatures = optional.NewString(supportedFeatures)

	clientAPI := createUDMClientToUDR(ueID, false)
	amfNon3GppAccessRegistration, resp, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.QueryAmfContextNon3gpp(context.Background(), ueID, &queryAmfContextNon3gppParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amfNon3GppAccessRegistration)
}

func HandleRegistrationAmf3gppAccess(respChan chan udm_message.HandlerResponseMessage, ueID string, body models.Amf3GppAccessRegistration) {
	contextExisted := false
	udm_context.CreateAmf3gppRegContext(ueID, body)
	if !udm_context.UdmAmf3gppRegContextNotExists(ueID) {
		contextExisted = true
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	var createAmfContext3gppParamOpts Nudr_DataRepository.CreateAmfContext3gppParamOpts
	optInterface := optional.NewInterface(body)
	createAmfContext3gppParamOpts.Amf3GppAccessRegistration = optInterface
	resp, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.CreateAmfContext3gpp(context.Background(), ueID, &createAmfContext3gppParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}

	if contextExisted {
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, nil)
	} else {
		h := make(http.Header)
		udmUe := udm_context.UDM_Self().UdmUePool[ueID]
		h.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriAmf3GppAccessRegistration))
		udm_message.SendHttpResponseMessage(respChan, h, http.StatusCreated, body)
	}
}

func HandleRegisterAmfNon3gppAccess(respChan chan udm_message.HandlerResponseMessage, ueID string, body models.AmfNon3GppAccessRegistration) {
	contextExisted := false
	udm_context.CreateAmfNon3gppRegContext(ueID, body)
	if !udm_context.UdmAmfNon3gppRegContextNotExists(ueID) {
		contextExisted = true
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	var createAmfContextNon3gppParamOpts Nudr_DataRepository.CreateAmfContextNon3gppParamOpts
	optInterface := optional.NewInterface(body)
	createAmfContextNon3gppParamOpts.AmfNon3GppAccessRegistration = optInterface
	resp, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.CreateAmfContextNon3gpp(context.Background(), ueID, &createAmfContextNon3gppParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}

	if contextExisted {
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, nil)
	} else {
		h := make(http.Header)
		udmUe := udm_context.UDM_Self().UdmUePool[ueID]
		h.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriAmfNon3GppAccessRegistration))
		udm_message.SendHttpResponseMessage(respChan, h, http.StatusCreated, body)
	}
}

func HandleUpdateAmf3gppAccess(respChan chan udm_message.HandlerResponseMessage, ueID string, body models.Amf3GppAccessRegistrationModification) {
	var patchItemReqArray []models.PatchItem
	currentContext := udm_context.GetAmf3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateAmf3gppAccess] Empty Amf3gppRegContext")
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "CONTEXT_NOT_FOUND"
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	if body.Guami != nil {
		udmUe := udm_context.UDM_Self().UdmUePool[ueID]
		if udmUe.SameAsStoredGUAMI3gpp(*body.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateAmf3gppAccess - deregistration")
			body.PurgeFlag = true
		} else {
			var problemDetails models.ProblemDetails
			problemDetails.Cause = "INVALID_GUAMI"
			logger.UecmLog.Errorln("INVALID_GUAMI")
			udm_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
		}

		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Guami"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = *body.Guami
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.PurgeFlag {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "PurgeFlag"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.PurgeFlag
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.Pei != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Pei"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.Pei
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.ImsVoPs != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "ImsVoPs"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.ImsVoPs
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.BackupAmfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "BackupAmfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.BackupAmfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.AmfContext3gpp(context.Background(), ueID, patchItemReqArray)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, nil)
}

func HandleUpdateAmfNon3gppAccess(respChan chan udm_message.HandlerResponseMessage, ueID string, body models.AmfNon3GppAccessRegistrationModification) {
	var patchItemReqArray []models.PatchItem
	currentContext := udm_context.GetAmfNon3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateAmfNon3gppAccess] Empty AmfNon3gppRegContext")
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "CONTEXT_NOT_FOUND"
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
		return
	}

	if body.Guami != nil {
		udmUe := udm_context.UDM_Self().UdmUePool[ueID]
		if udmUe.SameAsStoredGUAMINon3gpp(*body.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateAmfNon3gppAccess - deregistration")
			body.PurgeFlag = true
		} else {
			var problemDetails models.ProblemDetails
			problemDetails.Cause = "INVALID_GUAMI"
			logger.UecmLog.Errorln("INVALID_GUAMI")
			udm_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
		}

		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Guami"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = *body.Guami
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.PurgeFlag {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "PurgeFlag"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.PurgeFlag
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.Pei != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "Pei"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.Pei
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.ImsVoPs != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "ImsVoPs"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.ImsVoPs
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if body.BackupAmfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "BackupAmfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = body.BackupAmfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.AmfContextNon3gpp(context.Background(), ueID, patchItemReqArray)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, nil)
}

func HandleDeregistrationSmfRegistrations(respChan chan udm_message.HandlerResponseMessage, ueID string, pduSessionID string) {
	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.SMFRegistrationDocumentApi.DeleteSmfContext(context.Background(), ueID, pduSessionID)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, nil)
}

func HandleRegistrationSmfRegistrations(respChan chan udm_message.HandlerResponseMessage, ueID string, pduSessionID string, body models.SmfRegistration) {
	contextExisted := false
	udm_context.CreateSmfRegContext(ueID, pduSessionID)
	if !udm_context.UdmSmfRegContextNotExists(ueID) {
		contextExisted = true
	}

	pduID64, err := strconv.ParseInt(pduSessionID, 10, 32)
	if err != nil {
		logger.UecmLog.Errorln(err.Error())
	}
	pduID32 := int32(pduID64)

	var createSmfContextNon3gppParamOpts Nudr_DataRepository.CreateSmfContextNon3gppParamOpts
	optInterface := optional.NewInterface(body)
	createSmfContextNon3gppParamOpts.SmfRegistration = optInterface

	clientAPI := createUDMClientToUDR(ueID, false)
	resp, err := clientAPI.SMFRegistrationDocumentApi.CreateSmfContextNon3gpp(context.Background(), ueID, pduID32, &createSmfContextNon3gppParamOpts)
	if err != nil {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
		return
	}

	if contextExisted {
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, nil)
	} else {
		h := make(http.Header)
		udmUe := udm_context.UDM_Self().UdmUePool[ueID]
		h.Set("Location", udmUe.GetLocationURI(udm_context.LocationUriSmfRegistration))
		udm_message.SendHttpResponseMessage(respChan, h, http.StatusCreated, body)
	}
}
