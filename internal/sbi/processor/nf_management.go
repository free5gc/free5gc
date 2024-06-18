package processor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
	"github.com/free5gc/util/httpwrapper"
	timedecode "github.com/free5gc/util/mapstruct"
	"github.com/free5gc/util/mongoapi"
)

func (p *Processor) HandleNFDeregisterRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle NFDeregisterRequest")
	nfInstanceId := request.Params["nfInstanceID"]

	problemDetails := p.NFDeregisterProcedure(nfInstanceId)

	if problemDetails != nil {
		logger.NfmLog.Infoln("[NRF] Dergeister Success")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func (p *Processor) HandleGetNFInstanceRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle GetNFInstanceRequest")
	nfInstanceId := request.Params["nfInstanceID"]

	response := p.GetNFInstanceProcedure(nfInstanceId)

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "UNSPECIFIED",
		}
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func (p *Processor) HandleNFRegisterRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle NFRegisterRequest")
	nfProfile := request.Body.(models.NfProfile)

	header, response, isUpdate, problemDetails := p.NFRegisterProcedure(nfProfile)

	if response != nil {
		if isUpdate {
			logger.NfmLog.Traceln("update success")
			return httpwrapper.NewResponse(http.StatusOK, header, response)
		}
		logger.NfmLog.Traceln("register success")
		return httpwrapper.NewResponse(http.StatusCreated, header, response)
	} else if problemDetails != nil {
		logger.NfmLog.Traceln("register failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.NfmLog.Traceln("register failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func (p *Processor) HandleUpdateNFInstanceRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle UpdateNFInstanceRequest")
	nfInstanceID := request.Params["nfInstanceID"]
	patchJSON := request.Body.([]byte)

	response := p.UpdateNFInstanceProcedure(nfInstanceID, patchJSON)
	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func (p *Processor) HandleGetNFInstancesRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle GetNFInstancesRequest")
	nfType := request.Query.Get("nf-type")
	limit_param := request.Query.Get("limit")
	limit := 0
	if limit_param != "" {
		var err error
		limit, err = strconv.Atoi(request.Query.Get("limit"))
		if err != nil {
			logger.NfmLog.Errorln("Error in string conversion: ", limit)
			problemDetails := models.ProblemDetails{
				Title:  "Invalid Parameter",
				Status: http.StatusBadRequest,
				Detail: err.Error(),
			}

			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}
		if limit < 1 {
			problemDetails := models.ProblemDetails{
				Title:  "Invalid Parameter",
				Status: http.StatusBadRequest,
				Detail: "limit must be greater than 0",
			}
			return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
		}
	}

	response, problemDetails := p.GetNFInstancesProcedure(nfType, limit)
	if response != nil {
		logger.NfmLog.Traceln("GetNFInstances success")
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		logger.NfmLog.Traceln("GetNFInstances failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.NfmLog.Traceln("GetNFInstances failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func (p *Processor) HandleRemoveSubscriptionRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle RemoveSubscription")
	subscriptionID := request.Params["subscriptionID"]

	p.RemoveSubscriptionProcedure(subscriptionID)

	return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
}

func (p *Processor) HandleUpdateSubscriptionRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle UpdateSubscription")
	subscriptionID := request.Params["subscriptionID"]
	patchJSON := request.Body.([]byte)

	response := p.UpdateSubscriptionProcedure(subscriptionID, patchJSON)

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func (p *Processor) HandleCreateSubscriptionRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.NfmLog.Infoln("Handle CreateSubscriptionRequest")
	subscription := request.Body.(models.NrfSubscriptionData)

	response, problemDetails := p.CreateSubscriptionProcedure(subscription)
	if response != nil {
		logger.NfmLog.Traceln("CreateSubscription success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, response)
	} else if problemDetails != nil {
		logger.NfmLog.Traceln("CreateSubscription failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.NfmLog.Traceln("CreateSubscription failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func (p *Processor) CreateSubscriptionProcedure(
	subscription models.NrfSubscriptionData,
) (bson.M, *models.ProblemDetails) {
	subscriptionID, err := nrf_context.SetsubscriptionId()
	if err != nil {
		logger.NfmLog.Errorf("Unable to create subscription ID in CreateSubscriptionProcedure: %+v", err)
		return nil, &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "CREATE_SUBSCRIPTION_ERROR",
		}
	}
	subscription.SubscriptionId = subscriptionID

	tmp, err := json.Marshal(subscription)
	if err != nil {
		logger.NfmLog.Errorln("Marshal error in CreateSubscriptionProcedure: ", err)
		return nil, &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "CREATE_SUBSCRIPTION_ERROR",
		}
	}
	putData := bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.NfmLog.Errorln("Unmarshal error in CreateSubscriptionProcedure: ", err)
		return nil, &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "CREATE_SUBSCRIPTION_ERROR",
		}
	}

	// TODO: need to store Condition !
	existed, err := mongoapi.RestfulAPIPost("Subscriptions", bson.M{"subscriptionId": subscription.SubscriptionId},
		putData) // subscription id not exist before
	if err != nil || existed {
		if err != nil {
			logger.NfmLog.Errorf("CreateSubscriptionProcedure err: %+v", err)
		}
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "CREATE_SUBSCRIPTION_ERROR",
		}
		return nil, problemDetails
	}
	return putData, nil
}

func (p *Processor) UpdateSubscriptionProcedure(subscriptionID string, patchJSON []byte) map[string]interface{} {
	collName := "Subscriptions"
	filter := bson.M{"subscriptionId": subscriptionID}

	if err := mongoapi.RestfulAPIJSONPatch(collName, filter, patchJSON); err != nil {
		return nil
	} else {
		if response, err := mongoapi.RestfulAPIGetOne(collName, filter); err == nil {
			return response
		}
		return nil
	}
}

func (p *Processor) RemoveSubscriptionProcedure(subscriptionID string) {
	collName := "Subscriptions"
	filter := bson.M{"subscriptionId": subscriptionID}

	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		logger.NfmLog.Errorf("RemoveSubscriptionProcedure err: %+v", err)
	}
}

func (p *Processor) GetNFInstancesProcedure(nfType string, limit int) (*nrf_context.UriList, *models.ProblemDetails) {
	collName := "urilist"
	filter := bson.M{"nfType": nfType}
	if nfType == "" {
		// if the query parameter is not present, do not filter by nfType
		filter = bson.M{}
	}

	ULs, err := mongoapi.RestfulAPIGetMany(collName, filter)
	if err != nil {
		logger.NfmLog.Errorf("GetNFInstancesProcedure err: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetail
	}
	logger.NfmLog.Infoln("ULs: ", ULs)
	rspUriList := &nrf_context.UriList{}
	for _, UL := range ULs {
		originalUL := &nrf_context.UriList{}
		if err := mapstructure.Decode(UL, originalUL); err != nil {
			logger.NfmLog.Errorf("Decode error in GetNFInstancesProcedure: %+v", err)
			problemDetail := &models.ProblemDetails{
				Title:  "System failure",
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
				Cause:  "SYSTEM_FAILURE",
			}
			return nil, problemDetail
		}
		rspUriList.Link.Item = append(rspUriList.Link.Item, originalUL.Link.Item...)
		if nfType != "" && rspUriList.NfType == "" {
			rspUriList.NfType = originalUL.NfType
		}
	}

	nrf_context.NnrfUriListLimit(rspUriList, limit)
	return rspUriList, nil
}

func (p *Processor) NFDeregisterProcedure(nfInstanceID string) *models.ProblemDetails {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}

	nfProfilesRaw, err := mongoapi.RestfulAPIGetMany(collName, filter)
	if err != nil {
		logger.NfmLog.Errorf("NFDeregisterProcedure err: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return problemDetail
	}
	time.Sleep(time.Duration(1) * time.Second)

	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		logger.NfmLog.Errorf("NFDeregisterProcedure err: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return problemDetail
	}

	// nfProfile data for response
	var nfProfiles []models.NfProfile
	if err := timedecode.Decode(nfProfilesRaw, &nfProfiles); err != nil {
		logger.NfmLog.Warnln("Time decode error: ", err)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NOTIFICATION_ERROR",
			Detail: err.Error(),
		}
		return problemDetails
	}

	if len(nfProfiles) == 0 {
		logger.NfmLog.Warnf("NFProfile[%s] not found", nfInstanceID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_URI_STRUCTURE_NOT_FOUND",
			Detail: fmt.Sprintf("NFProfile[%s] not found", nfInstanceID),
		}
		return problemDetails
	}

	uriList := nrf_context.GetNofificationUri(nfProfiles[0])
	nfInstanceType := nfProfiles[0].NfType
	nfInstanceUri := nrf_context.GetNfInstanceURI(nfInstanceID)
	// set info for NotificationData
	Notification_event := models.NotificationEventType_DEREGISTERED

	for _, uri := range uriList {
		problemDetails := p.Consumer().SendNFStatusNotify(Notification_event, nfInstanceUri, uri, nil)
		if problemDetails != nil {
			return problemDetails
		}
	}

	collNameURI := "urilist"
	filterURI := bson.M{"nfType": nfProfiles[0].NfType}
	putData := bson.M{"_link.item": bson.M{"href": nfInstanceUri}, "multi": true}
	if err := mongoapi.RestfulAPIPullOne(collNameURI, filterURI, putData); err != nil {
		logger.NfmLog.Errorf("NFDeregisterProcedure err: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return problemDetail
	}
	if factory.NrfConfig.GetOAuth() {
		nfCertPath := oauth.GetNFCertPath(factory.NrfConfig.GetCertBasePath(), string(nfInstanceType), nfInstanceID)
		err := os.Remove(nfCertPath)
		if err != nil {
			logger.NfmLog.Warningf("Can not delete NFCertPem file: %v: %v", nfCertPath, err)
		}
	}
	return nil
}

func (p *Processor) UpdateNFInstanceProcedure(nfInstanceID string, patchJSON []byte) map[string]interface{} {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}

	if err := mongoapi.RestfulAPIJSONPatch(collName, filter, patchJSON); err != nil {
		logger.NfmLog.Errorf("UpdateNFInstanceProcedure err: %+v", err)
		return nil
	}

	nf, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.NfmLog.Errorf("UpdateNFInstanceProcedure err: %+v", err)
		return nil
	}

	nfProfilesRaw := []map[string]interface{}{
		nf,
	}

	var nfProfiles []models.NfProfile
	if err := timedecode.Decode(nfProfilesRaw, &nfProfiles); err != nil {
		logger.NfmLog.Errorf("UpdateNFInstanceProcedure err: %+v", err)
	}

	if len(nfProfiles) == 0 {
		logger.NfmLog.Warnf("NFProfile[%s] not found", nfInstanceID)
		return nil
	}

	uriList := nrf_context.GetNofificationUri(nfProfiles[0])

	// set info for NotificationData
	Notification_event := models.NotificationEventType_PROFILE_CHANGED
	nfInstanceUri := nrf_context.GetNfInstanceURI(nfInstanceID)

	for _, uri := range uriList {
		p.Consumer().SendNFStatusNotify(Notification_event, nfInstanceUri, uri, &nfProfiles[0])
	}

	return nf
}

func (p *Processor) GetNFInstanceProcedure(nfInstanceID string) map[string]interface{} {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}
	response, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.NfmLog.Errorf("GetNFInstanceProcedure err: %+v", err)
		return nil
	}

	return response
}

func (p *Processor) NFRegisterProcedure(
	nfProfile models.NfProfile,
) (
	header http.Header, response bson.M,
	update bool, problemDetails *models.ProblemDetails,
) {
	logger.NfmLog.Traceln("[NRF] In NFRegisterProcedure")
	var nf models.NfProfile

	err := nrf_context.NnrfNFManagementDataModel(&nf, nfProfile)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
		return nil, nil, false, problemDetails
	}

	// make location header
	locationHeaderValue := nrf_context.SetLocationHeader(nfProfile)
	// Marshal nf to bson
	tmp, err := json.Marshal(nf)
	if err != nil {
		logger.NfmLog.Errorln("Marshal error in NFRegisterProcedure: ", err)
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, nil, false, problemDetails
	}
	putData := bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.NfmLog.Errorln("Unmarshal error in NFRegisterProcedure: ", err)
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, nil, false, problemDetails
	}
	// set db info
	collName := "NfProfile"
	nfInstanceId := nf.NfInstanceId
	filter := bson.M{"nfInstanceId": nfInstanceId}

	// Update NF Profile case
	existed, err := mongoapi.RestfulAPIPutOne(collName, filter, putData)
	if err != nil {
		logger.NfmLog.Errorf("NFRegisterProcedure err: %+v", err)
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, nil, false, problemDetails
	}

	if existed {
		logger.NfmLog.Infoln("RestfulAPIPutOne Update")
		uriList := nrf_context.GetNofificationUri(nf)

		// set info for NotificationData
		Notification_event := models.NotificationEventType_PROFILE_CHANGED
		nfInstanceUri := locationHeaderValue

		// receive the rsp from handler
		for _, uri := range uriList {
			problemDetails := p.Consumer().SendNFStatusNotify(Notification_event, nfInstanceUri, uri, &nfProfile)
			if problemDetails != nil {
				return nil, nil, true, problemDetails
			}
		}

		header := make(http.Header)
		header.Add("Location", locationHeaderValue)
		return header, putData, true, nil
	} else { // Create NF Profile case
		logger.NfmLog.Infoln("Create NF Profile")
		uriList := nrf_context.GetNofificationUri(nf)
		// set info for NotificationData
		Notification_event := models.NotificationEventType_REGISTERED
		nfInstanceUri := locationHeaderValue

		for _, uri := range uriList {
			problemDetails := p.Consumer().SendNFStatusNotify(Notification_event, nfInstanceUri, uri, &nfProfile)
			if problemDetails != nil {
				return nil, nil, false, problemDetails
			}
		}

		header := make(http.Header)
		header.Add("Location", locationHeaderValue)
		logger.NfmLog.Infoln("Location header: ", locationHeaderValue)

		if factory.NrfConfig.GetOAuth() {
			// Generate NF's pubkey certificate with root certificate
			err := nrf_context.SignNFCert(string(nf.NfType), nfInstanceId)
			if err != nil {
				logger.NfmLog.Warnln(err)
			}
		}
		return header, putData, false, nil
	}
}
