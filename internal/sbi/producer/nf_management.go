package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi/Nnrf_NFManagement"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	timedecode "github.com/free5gc/util/mapstruct"
	"github.com/free5gc/util/mongoapi"
)

func HandleNFDeregisterRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle NFDeregisterRequest")
	nfInstanceId := request.Params["nfInstanceID"]

	problemDetails := NFDeregisterProcedure(nfInstanceId)

	if problemDetails != nil {
		logger.ManagementLog.Infoln("[NRF] Dergeister Success")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func HandleGetNFInstanceRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle GetNFInstanceRequest")
	nfInstanceId := request.Params["nfInstanceID"]

	response := GetNFInstanceProcedure(nfInstanceId)

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

func HandleNFRegisterRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle NFRegisterRequest")
	nfProfile := request.Body.(models.NfProfile)

	header, response, isUpdate, problemDetails := NFRegisterProcedure(nfProfile)

	if response != nil {
		if isUpdate {
			logger.ManagementLog.Traceln("update success")
			return httpwrapper.NewResponse(http.StatusOK, header, response)
		}
		logger.ManagementLog.Traceln("register success")
		return httpwrapper.NewResponse(http.StatusCreated, header, response)
	} else if problemDetails != nil {
		logger.ManagementLog.Traceln("register failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.ManagementLog.Traceln("register failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func HandleUpdateNFInstanceRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle UpdateNFInstanceRequest")
	nfInstanceID := request.Params["nfInstanceID"]
	patchJSON := request.Body.([]byte)

	response := UpdateNFInstanceProcedure(nfInstanceID, patchJSON)
	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func HandleGetNFInstancesRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle GetNFInstancesRequest")
	nfType := request.Query.Get("nf-type")
	limit, err := strconv.Atoi(request.Query.Get("limit"))
	if err != nil {
		logger.ManagementLog.Errorln("Error in string conversion: ", limit)
		problemDetails := models.ProblemDetails{
			Title:  "Invalid Parameter",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}

		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	response, problemDetails := GetNFInstancesProcedure(nfType, limit)
	if response != nil {
		logger.ManagementLog.Traceln("GetNFInstances success")
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		logger.ManagementLog.Traceln("GetNFInstances failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.ManagementLog.Traceln("GetNFInstances failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func HandleRemoveSubscriptionRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle RemoveSubscription")
	subscriptionID := request.Params["subscriptionID"]

	RemoveSubscriptionProcedure(subscriptionID)

	return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
}

func HandleUpdateSubscriptionRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle UpdateSubscription")
	subscriptionID := request.Params["subscriptionID"]
	patchJSON := request.Body.([]byte)

	response := UpdateSubscriptionProcedure(subscriptionID, patchJSON)

	if response != nil {
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return httpwrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func HandleCreateSubscriptionRequest(request *httpwrapper.Request) *httpwrapper.Response {
	logger.ManagementLog.Infoln("Handle CreateSubscriptionRequest")
	subscription := request.Body.(models.NrfSubscriptionData)

	response, problemDetails := CreateSubscriptionProcedure(subscription)
	if response != nil {
		logger.ManagementLog.Traceln("CreateSubscription success")
		return httpwrapper.NewResponse(http.StatusCreated, nil, response)
	} else if problemDetails != nil {
		logger.ManagementLog.Traceln("CreateSubscription failed")
		return httpwrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.ManagementLog.Traceln("CreateSubscription failed")
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func CreateSubscriptionProcedure(subscription models.NrfSubscriptionData) (bson.M, *models.ProblemDetails) {
	subscription.SubscriptionId = nrf_context.SetsubscriptionId()

	tmp, err := json.Marshal(subscription)
	if err != nil {
		logger.ManagementLog.Errorln("Marshal error in CreateSubscriptionProcedure: ", err)
	}
	putData := bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.ManagementLog.Errorln("Unmarshal error in CreateSubscriptionProcedure: ", err)
	}

	// TODO: need to store Condition !
	existed, err := mongoapi.RestfulAPIPost("Subscriptions", bson.M{"subscriptionId": subscription.SubscriptionId},
		putData) // subscription id not exist before
	if err != nil || existed {
		if err != nil {
			logger.ManagementLog.Errorf("CreateSubscriptionProcedure err: %+v", err)
		}
		problemDetails := &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "CREATE_SUBSCRIPTION_ERROR",
		}
		return nil, problemDetails
	}
	return putData, nil
}

func UpdateSubscriptionProcedure(subscriptionID string, patchJSON []byte) map[string]interface{} {
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

func RemoveSubscriptionProcedure(subscriptionID string) {
	collName := "Subscriptions"
	filter := bson.M{"subscriptionId": subscriptionID}

	if err := mongoapi.RestfulAPIDeleteMany(collName, filter); err != nil {
		logger.ManagementLog.Errorf("RemoveSubscriptionProcedure err: %+v", err)
	}
}

func GetNFInstancesProcedure(nfType string, limit int) (*nrf_context.UriList, *models.ProblemDetails) {
	// nfType := c.Query("nf-type")
	// limit, err := strconv.Atoi(c.Query("limit"))
	collName := "urilist"
	filter := bson.M{"nfType": nfType}

	UL, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.ManagementLog.Errorf("GetNFInstancesProcedure err: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetail
	}
	logger.ManagementLog.Infoln("UL: ", UL)
	originalUL := &nrf_context.UriList{}
	if err := mapstructure.Decode(UL, originalUL); err != nil {
		logger.ManagementLog.Errorf("Decode error in GetNFInstancesProcedure: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetail
	}
	nrf_context.NnrfUriListLimit(originalUL, limit)
	// c.JSON(http.StatusOK, originalUL)
	return originalUL, nil
}

func NFDeregisterProcedure(nfInstanceID string) *models.ProblemDetails {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}

	nfProfilesRaw, err := mongoapi.RestfulAPIGetMany(collName, filter)
	if err != nil {
		logger.ManagementLog.Errorf("NFDeregisterProcedure err: %+v", err)
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
		logger.ManagementLog.Errorf("NFDeregisterProcedure err: %+v", err)
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
		logger.ManagementLog.Warnln("Time decode error: ", err)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NOTIFICATION_ERROR",
			Detail: err.Error(),
		}
		return problemDetails
	}

	if len(nfProfiles) == 0 {
		logger.ManagementLog.Warnf("NFProfile[%s] not found", nfInstanceID)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_URI_STRUCTURE_NOT_FOUND",
			Detail: fmt.Sprintf("NFProfile[%s] not found", nfInstanceID),
		}
		return problemDetails
	}

	uriList := nrf_context.GetNofificationUri(nfProfiles[0])

	nfInstanceUri := nrf_context.GetNfInstanceURI(nfInstanceID)
	// set info for NotificationData
	Notification_event := models.NotificationEventType_DEREGISTERED

	for _, uri := range uriList {
		problemDetails := SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
		if problemDetails != nil {
			return problemDetails
		}
	}

	collNameURI := "urilist"
	filterURI := bson.M{"nfType": nfProfiles[0].NfType}
	putData := bson.M{"_link.item": bson.M{"href": nfInstanceUri}, "multi": true}
	if err := mongoapi.RestfulAPIPullOne(collNameURI, filterURI, putData); err != nil {
		logger.ManagementLog.Errorf("NFDeregisterProcedure err: %+v", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return problemDetail
	}

	return nil
}

func UpdateNFInstanceProcedure(nfInstanceID string, patchJSON []byte) map[string]interface{} {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}

	if err := mongoapi.RestfulAPIJSONPatch(collName, filter, patchJSON); err != nil {
		logger.ManagementLog.Errorf("UpdateNFInstanceProcedure err: %+v", err)
		return nil
	}

	nf, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.ManagementLog.Errorf("UpdateNFInstanceProcedure err: %+v", err)
		return nil
	}

	nfProfilesRaw := []map[string]interface{}{
		nf,
	}

	var nfProfiles []models.NfProfile
	if err := timedecode.Decode(nfProfilesRaw, &nfProfiles); err != nil {
		logger.ManagementLog.Errorf("UpdateNFInstanceProcedure err: %+v", err)
	}

	if len(nfProfiles) == 0 {
		logger.ManagementLog.Warnf("NFProfile[%s] not found", nfInstanceID)
		return nil
	}

	uriList := nrf_context.GetNofificationUri(nfProfiles[0])

	// set info for NotificationData
	Notification_event := models.NotificationEventType_PROFILE_CHANGED
	nfInstanceUri := nrf_context.GetNfInstanceURI(nfInstanceID)

	for _, uri := range uriList {
		SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
	}

	return nf
}

func GetNFInstanceProcedure(nfInstanceID string) map[string]interface{} {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}
	response, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.ManagementLog.Errorf("GetNFInstanceProcedure err: %+v", err)
		return nil
	}

	return response
}

func NFRegisterProcedure(
	nfProfile models.NfProfile,
) (
	header http.Header, response bson.M,
	update bool, problemDetails *models.ProblemDetails,
) {
	logger.ManagementLog.Traceln("[NRF] In NFRegisterProcedure")
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
		logger.ManagementLog.Errorln("Marshal error in NFRegisterProcedure: ", err)
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
		logger.ManagementLog.Errorln("Unmarshal error in NFRegisterProcedure: ", err)
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
		logger.ManagementLog.Errorf("NFRegisterProcedure err: %+v", err)
		problemDetails := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, nil, false, problemDetails
	}

	if existed {
		logger.ManagementLog.Infoln("RestfulAPIPutOne Update")
		uriList := nrf_context.GetNofificationUri(nf)

		// set info for NotificationData
		Notification_event := models.NotificationEventType_PROFILE_CHANGED
		nfInstanceUri := locationHeaderValue

		// receive the rsp from handler
		for _, uri := range uriList {
			problemDetails := SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
			if problemDetails != nil {
				return nil, nil, true, problemDetails
			}
		}

		header := make(http.Header)
		header.Add("Location", locationHeaderValue)
		return header, putData, true, nil
	} else { // Create NF Profile case
		logger.ManagementLog.Infoln("Create NF Profile")
		uriList := nrf_context.GetNofificationUri(nf)
		// set info for NotificationData
		Notification_event := models.NotificationEventType_REGISTERED
		nfInstanceUri := locationHeaderValue

		for _, uri := range uriList {
			problemDetails := SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
			if problemDetails != nil {
				return nil, nil, false, problemDetails
			}
		}

		header := make(http.Header)
		header.Add("Location", locationHeaderValue)
		logger.ManagementLog.Infoln("Location header: ", locationHeaderValue)

		if factory.NrfConfig.GetOAuth() {
			// Generate NF's pubkey certificate with root certificate
			err := nrf_context.SignNFCert(string(nf.NfType), nfInstanceId)
			if err != nil {
				logger.ManagementLog.Warnln(err)
			}
		}
		return header, putData, false, nil
	}
}

func SendNFStatusNotify(Notification_event models.NotificationEventType, nfInstanceUri string,
	url string,
) *models.ProblemDetails {
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	// url = fmt.Sprintf("%s%s", url, "/notification")

	configuration.SetBasePathNoGroup(url)
	notifcationData := models.NotificationData{
		Event:         Notification_event,
		NfInstanceUri: nfInstanceUri,
	}
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	res, err := client.NotificationApi.NotificationPost(context.TODO(), notifcationData)
	if err != nil {
		logger.ManagementLog.Infof("Notify fail: %v", err)
		problemDetails := &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NOTIFICATION_ERROR",
			Detail: err.Error(),
		}
		return problemDetails
	}
	if res != nil {
		defer func() {
			if resCloseErr := res.Body.Close(); resCloseErr != nil {
				logger.ManagementLog.Errorf("NotificationApi response body cannot close: %+v", resCloseErr)
			}
		}()
		if status := res.StatusCode; status != http.StatusNoContent {
			logger.ManagementLog.Warnln("Error status in NotificationPost: ", status)
			problemDetails := &models.ProblemDetails{
				Status: int32(status),
				Cause:  "NOTIFICATION_ERROR",
			}
			return problemDetails
		}
	}
	return nil
}
