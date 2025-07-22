package processor

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nudr_DataRepository "github.com/free5gc/openapi/udr/DataRepository"
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/internal/logger"
)

// ue_context_managemanet_service
func (p *Processor) GetAmf3gppAccessProcedure(c *gin.Context, ueID string, supportedFeatures string) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	var queryAmfContext3gppRequest Nudr_DataRepository.QueryAmfContext3gppRequest
	queryAmfContext3gppRequest.UeId = &ueID
	queryAmfContext3gppRequest.SupportedFeatures = &supportedFeatures

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	amf3GppAccessRegistration, err := clientAPI.AMF3GPPAccessRegistrationDocumentApi.
		QueryAmfContext3gpp(ctx, &queryAmfContext3gppRequest)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	c.JSON(http.StatusOK, amf3GppAccessRegistration.Amf3GppAccessRegistration)
}

func (p *Processor) GetAmfNon3gppAccessProcedure(c *gin.Context, queryAmfContextNon3gppParamOpts Nudr_DataRepository.
	QueryAmfContextNon3gppRequest, ueID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
	}
	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	amfNon3GppAccessRegistrationResponse, err := clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.
		QueryAmfContextNon3gpp(ctx, &queryAmfContextNon3gppParamOpts)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	c.JSON(http.StatusOK, amfNon3GppAccessRegistrationResponse.AmfNon3GppAccessRegistration)
}

func (p *Processor) RegistrationAmf3gppAccessProcedure(c *gin.Context,
	registerRequest models.Amf3GppAccessRegistration,
	ueID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	// TODO: EPS interworking with N26 is not supported yet in this stage
	var oldAmf3GppAccessRegContext *models.Amf3GppAccessRegistration
	var ue *udm_context.UdmUeContext

	if p.Context().UdmAmf3gppRegContextExists(ueID) {
		ue, _ = p.Context().UdmUeFindBySupi(ueID)
		oldAmf3GppAccessRegContext = ue.Amf3GppAccessRegistration
	}

	p.Context().CreateAmf3gppRegContext(ueID, registerRequest)

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	var createAmfContext3gppRequest Nudr_DataRepository.CreateAmfContext3gppRequest
	createAmfContext3gppRequest.UeId = &ueID
	createAmfContext3gppRequest.Amf3GppAccessRegistration = &registerRequest
	_, err = clientAPI.AMF3GPPAccessRegistrationDocumentApi.CreateAmfContext3gpp(ctx,
		&createAmfContext3gppRequest)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	// TS 23.502 4.2.2.2.2 14d: UDM initiate a Nudm_UECM_DeregistrationNotification to the old AMF
	// corresponding to the same (e.g. 3GPP) access, if one exists
	if oldAmf3GppAccessRegContext != nil {
		if !ue.SameAsStoredGUAMI3gpp(*oldAmf3GppAccessRegContext.Guami) {
			// Based on TS 23.502 4.2.2.2.2, If the serving NF removal reason indicated by the UDM is Initial Registration,
			// the old AMF invokes the Nsmf_PDUSession_ReleaseSMContext (SM Context ID). Thus we give different
			// dereg cause based on registration parameter from serving AMF
			deregReason := models.UdmUecmDeregistrationReason_UE_REGISTRATION_AREA_CHANGE
			if registerRequest.InitialRegistrationInd {
				deregReason = models.UdmUecmDeregistrationReason_UE_INITIAL_REGISTRATION
			}
			deregistData := models.UdmUecmDeregistrationData{
				DeregReason: deregReason,
				AccessType:  models.AccessType__3_GPP_ACCESS,
			}

			go func() {
				logger.UecmLog.Infof("Send DeregNotify to old AMF GUAMI=%v", oldAmf3GppAccessRegContext.Guami)
				pd := p.SendOnDeregistrationNotification(ueID,
					oldAmf3GppAccessRegContext.DeregCallbackUri,
					deregistData) // Deregistration Notify Triggered
				if pd != nil {
					logger.UecmLog.Errorf("RegistrationAmf3gppAccess: send DeregNotify fail %v", pd)
				}
			}()
		}

		c.JSON(http.StatusOK, registerRequest)
	} else {
		udmUe, _ := p.Context().UdmUeFindBySupi(ueID)
		c.Header("Location", udmUe.GetLocationURI(udm_context.LocationUriAmf3GppAccessRegistration))
		c.JSON(http.StatusCreated, registerRequest)
	}
}

func (p *Processor) RegisterAmfNon3gppAccessProcedure(c *gin.Context,
	registerRequest models.AmfNon3GppAccessRegistration,
	ueID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	var oldAmfNon3GppAccessRegContext *models.AmfNon3GppAccessRegistration
	if p.Context().UdmAmfNon3gppRegContextExists(ueID) {
		ue, _ := p.Context().UdmUeFindBySupi(ueID)
		oldAmfNon3GppAccessRegContext = ue.AmfNon3GppAccessRegistration
	}

	p.Context().CreateAmfNon3gppRegContext(ueID, registerRequest)

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	var createAmfContextNon3gppRequest Nudr_DataRepository.CreateAmfContextNon3gppRequest
	createAmfContextNon3gppRequest.UeId = &ueID
	createAmfContextNon3gppRequest.AmfNon3GppAccessRegistration = &registerRequest

	_, err = clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.CreateAmfContextNon3gpp(
		ctx, &createAmfContextNon3gppRequest)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	// TS 23.502 4.2.2.2.2 14d: UDM initiate a Nudm_UECM_DeregistrationNotification to the old AMF
	// corresponding to the same (e.g. 3GPP) access, if one exists
	if oldAmfNon3GppAccessRegContext != nil {
		deregistData := models.UdmUecmDeregistrationData{
			DeregReason: models.UdmUecmDeregistrationReason_UE_INITIAL_REGISTRATION,
			AccessType:  models.AccessType_NON_3_GPP_ACCESS,
		}
		p.SendOnDeregistrationNotification(ueID, oldAmfNon3GppAccessRegContext.DeregCallbackUri,
			deregistData) // Deregistration Notify Triggered

		return
	} else {
		udmUe, _ := p.Context().UdmUeFindBySupi(ueID)
		c.Header("Location", udmUe.GetLocationURI(udm_context.LocationUriAmfNon3GppAccessRegistration))
		c.JSON(http.StatusCreated, registerRequest)
	}
}

func (p *Processor) UpdateAmf3gppAccessProcedure(c *gin.Context,
	request models.Amf3GppAccessRegistrationModification,
	ueID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	var patchItemReqArray []models.PatchItem
	currentContext := p.Context().GetAmf3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateAmf3gppAccess] Empty Amf3gppRegContext")
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	if request.Guami != nil {
		udmUe, _ := p.Context().UdmUeFindBySupi(ueID)
		if udmUe.SameAsStoredGUAMI3gpp(*request.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateAmf3gppAccess - deregistration")
			request.PurgeFlag = true
		} else {
			logger.UecmLog.Errorln("INVALID_GUAMI")
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INVALID_GUAMI",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}

		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "guami"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = *request.Guami
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.PurgeFlag {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "purgeFlag"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.PurgeFlag
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.Pei != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "pei"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.Pei
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.ImsVoPs != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "imsVoPs"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.ImsVoPs
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.BackupAmfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "backupAmfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.BackupAmfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	var amfContext3gppRequest Nudr_DataRepository.AmfContext3gppRequest
	amfContext3gppRequest.UeId = &ueID
	amfContext3gppRequest.PatchItem = patchItemReqArray
	_, err = clientAPI.AMF3GPPAccessRegistrationDocumentApi.AmfContext3gpp(ctx,
		&amfContext3gppRequest)
	if err != nil {
		if apiErr, ok := err.(openapi.GenericOpenAPIError); ok {
			if amfContext3gppErr, ok2 := apiErr.Model().(Nudr_DataRepository.AmfContext3gppError); ok2 {
				problem := amfContext3gppErr.ProblemDetails
				c.JSON(int(problem.Status), problem)
				return
			}
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	if request.PurgeFlag {
		udmUe, _ := p.Context().UdmUeFindBySupi(ueID)
		udmUe.Amf3GppAccessRegistration = nil
	}

	c.Status(http.StatusNoContent)
}

func (p *Processor) UpdateAmfNon3gppAccessProcedure(c *gin.Context,
	request models.AmfNon3GppAccessRegistrationModification,
	ueID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	var patchItemReqArray []models.PatchItem
	currentContext := p.Context().GetAmfNon3gppRegContext(ueID)
	if currentContext == nil {
		logger.UecmLog.Errorln("[UpdateAmfNon3gppAccess] Empty AmfNon3gppRegContext")
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "CONTEXT_NOT_FOUND",
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	if request.Guami != nil {
		udmUe, _ := p.Context().UdmUeFindBySupi(ueID)
		if udmUe.SameAsStoredGUAMINon3gpp(*request.Guami) { // deregistration
			logger.UecmLog.Infoln("UpdateAmfNon3gppAccess - deregistration")
			request.PurgeFlag = true
		} else {
			logger.UecmLog.Errorln("INVALID_GUAMI")
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "INVALID_GUAMI",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}

		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "guami"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = *request.Guami
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.PurgeFlag {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "purgeFlag"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.PurgeFlag
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.Pei != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "pei"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.Pei
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.ImsVoPs != "" {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "imsVoPs"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.ImsVoPs
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	if request.BackupAmfInfo != nil {
		var patchItemTmp models.PatchItem
		patchItemTmp.Path = "/" + "backupAmfInfo"
		patchItemTmp.Op = models.PatchOperation_REPLACE
		patchItemTmp.Value = request.BackupAmfInfo
		patchItemReqArray = append(patchItemReqArray, patchItemTmp)
	}

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	var amfContextNon3gppRequest Nudr_DataRepository.AmfContextNon3gppRequest
	amfContextNon3gppRequest.UeId = &ueID
	amfContextNon3gppRequest.PatchItem = patchItemReqArray
	_, err = clientAPI.AMFNon3GPPAccessRegistrationDocumentApi.AmfContextNon3gpp(ctx,
		&amfContextNon3gppRequest)
	if err != nil {
		if apiErr, ok := err.(openapi.GenericOpenAPIError); ok {
			if amfContextNon3gppErr, ok2 := apiErr.Model().(Nudr_DataRepository.AmfContextNon3gppError); ok2 {
				problem := amfContextNon3gppErr.ProblemDetails
				c.JSON(int(problem.Status), problem)
				return
			}
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	c.Status(http.StatusNoContent)
}

func (p *Processor) DeregistrationSmfRegistrationsProcedure(c *gin.Context,
	ueID string,
	pduSessionID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	var pduSessionIDInt32 int32
	num, err := strconv.ParseInt(pduSessionID, 10, 32)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	pduSessionIDInt32 = int32(num)
	var deleteSmfRegistrationRequest Nudr_DataRepository.DeleteSmfRegistrationRequest
	deleteSmfRegistrationRequest.UeId = &ueID
	deleteSmfRegistrationRequest.PduSessionId = &pduSessionIDInt32
	_, err = clientAPI.SMFRegistrationDocumentApi.DeleteSmfRegistration(ctx, &deleteSmfRegistrationRequest)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	c.Status(http.StatusNoContent)
}

func (p *Processor) RegistrationSmfRegistrationsProcedure(
	c *gin.Context,
	smfRegistration *models.SmfRegistration,
	ueID string,
	pduSessionID string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	contextExisted := false
	p.Context().CreateSmfRegContext(ueID, pduSessionID)
	if !p.Context().UdmSmfRegContextNotExists(ueID) {
		contextExisted = true
	}

	pduID64, err := strconv.ParseInt(pduSessionID, 10, 32)
	if err != nil {
		logger.UecmLog.Errorln(err.Error())
	}
	pduID32 := int32(pduID64)

	var createSmfContext3gppRequest Nudr_DataRepository.CreateOrUpdateSmfRegistrationRequest
	createSmfContext3gppRequest.UeId = &ueID
	createSmfContext3gppRequest.SmfRegistration = smfRegistration
	createSmfContext3gppRequest.PduSessionId = &pduID32

	clientAPI, err := p.Consumer().CreateUDMClientToUDR(ueID)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	_, err = clientAPI.SMFRegistrationDocumentApi.CreateOrUpdateSmfRegistration(ctx, &createSmfContext3gppRequest)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	if contextExisted {
		c.Status(http.StatusNoContent)
	} else {
		udmUe, _ := p.Context().UdmUeFindBySupi(ueID)
		c.Header("Location", udmUe.GetLocationURI(udm_context.LocationUriSmfRegistration))
		c.JSON(http.StatusCreated, smfRegistration)
	}
}
