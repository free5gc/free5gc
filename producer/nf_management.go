package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"free5gc/lib/MongoDBLibrary"
	"free5gc/lib/TimeDecode"
	"free5gc/lib/http_wrapper"
	//"free5gc/lib/openapi"
	"free5gc/lib/openapi/Nnrf_NFManagement"
	"free5gc/lib/openapi/models"
	nrf_context "free5gc/src/nrf/context"
	//"free5gc/src/nrf/factory"
	"free5gc/src/nrf/logger"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
)

func HandleNFDeregisterRequest(request *http_wrapper.Request) *http_wrapper.Response {

	logger.ManagementLog.Infoln("Handle NFDeregisterRequest")
	nfInstanceId := request.Params["nfInstanceID"]

	problemDetails := NFDeregisterProcedure(nfInstanceId)

	if problemDetails != nil {
		logger.ManagementLog.Infoln("[NRF] Dergeister Success")
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}

}

func HandleGetNFInstanceRequest(request *http_wrapper.Request) *http_wrapper.Response {

	logger.ManagementLog.Infoln("Handle GetNFInstanceRequest")
	nfInstanceId := request.Params["nfInstanceID"]

	response := GetNFInstanceProcedure(nfInstanceId)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "UNSPECIFIED",
		}
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
}

func HandleNFRegisterRequest(request *http_wrapper.Request) *http_wrapper.Response {

	logger.ManagementLog.Infoln("Handle NFRegisterRequest")
	nfProfile := request.Body.(models.NfProfile)

	header, response, problemDetails := NFRegisterProcedure(nfProfile)

	if response != nil {

		logger.ManagementLog.Traceln("register success")
		return http_wrapper.NewResponse(http.StatusCreated, header, response)
	} else if problemDetails != nil {
		logger.ManagementLog.Traceln("register failed")
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.ManagementLog.Traceln("register failed")
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)

}

func HandleUpdateNFInstanceRequest(request *http_wrapper.Request) *http_wrapper.Response {

	logger.ManagementLog.Infoln("Handle UpdateNFInstanceRequest")
	nfInstanceID := request.Params["nfInstanceID"]
	patchJSON := request.Body.([]byte)

	response := UpdateNFInstanceProcedure(nfInstanceID, patchJSON)
	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {

		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func HandleGetNFInstancesRequest(request *http_wrapper.Request) *http_wrapper.Response {

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

		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}

	response, problemDetails := GetNFInstancesProcedure(nfType, limit)
	if response != nil {

		logger.ManagementLog.Traceln("GetNFInstances success")
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else if problemDetails != nil {
		logger.ManagementLog.Traceln("GetNFInstances failed")
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.ManagementLog.Traceln("GetNFInstances failed")
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func HandleRemoveSubscriptionRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.ManagementLog.Infoln("Handle RemoveSubscription")
	subscriptionID := request.Params["subscriptionID"]

	RemoveSubscriptionProcedure(subscriptionID)

	return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
}

func HandleUpdateSubscriptionRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.ManagementLog.Infoln("Handle UpdateSubscription")
	subscriptionID := request.Params["subscriptionID"]
	patchJSON := request.Body.([]byte)

	response := UpdateSubscriptionProcedure(subscriptionID, patchJSON)

	if response != nil {
		return http_wrapper.NewResponse(http.StatusOK, nil, response)
	} else {
		return http_wrapper.NewResponse(http.StatusNoContent, nil, nil)
	}
}

func HandleCreateSubscriptionRequest(request *http_wrapper.Request) *http_wrapper.Response {
	logger.ManagementLog.Infoln("Handle CreateSubscriptionRequest")
	subscription := request.Body.(models.NrfSubscriptionData)

	response, problemDetails := CreateSubscriptionProcedure(subscription)
	if response != nil {
		logger.ManagementLog.Traceln("CreateSubscription success")
		return http_wrapper.NewResponse(http.StatusCreated, nil, response)
	} else if problemDetails != nil {
		logger.ManagementLog.Traceln("CreateSubscription failed")
		return http_wrapper.NewResponse(int(problemDetails.Status), nil, problemDetails)
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	logger.ManagementLog.Traceln("CreateSubscription failed")
	return http_wrapper.NewResponse(http.StatusForbidden, nil, problemDetails)

}

func CreateSubscriptionProcedure(subscription models.NrfSubscriptionData) (response bson.M,
	problemDetails *models.ProblemDetails) {
	subscription.SubscriptionId = nrf_context.SetsubscriptionId()

	tmp, err := json.Marshal(subscription)
	if err != nil {
		logger.ManagementLog.Errorln("Marshal error in CreateSubscriptionProcedure: ", err)
	}
	var putData = bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.ManagementLog.Errorln("Unmarshal error in CreateSubscriptionProcedure: ", err)
	}

	// TODO: need to store Condition !
	if !MongoDBLibrary.RestfulAPIPost("Subscriptions", bson.M{"subscriptionId": subscription.SubscriptionId},
		putData) { // subscription id not exist before
		return putData, nil
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "CREATE_SUBSCIRPTION_ERROR",
		}

		return nil, problemDetails
	}
}

func UpdateSubscriptionProcedure(subscriptionID string, patchJSON []byte) (response map[string]interface{}) {
	collName := "Subscriptions"
	filter := bson.M{"subscriptionId": subscriptionID}

	if MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON) {
		response = MongoDBLibrary.RestfulAPIGetOne(collName, filter)
		return response
	} else {
		return nil
	}

}

func RemoveSubscriptionProcedure(subscriptionID string) {

	collName := "Subscriptions"
	filter := bson.M{"subscriptionId": subscriptionID}

	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)
}

func GetNFInstancesProcedure(nfType string, limit int) (response *nrf_context.UriList,
	problemDetail *models.ProblemDetails) {
	// nfType := c.Query("nf-type")
	// limit, err := strconv.Atoi(c.Query("limit"))
	collName := "urilist"
	filter := bson.M{"nfType": nfType}

	UL := MongoDBLibrary.RestfulAPIGetOne(collName, filter)
	logger.ManagementLog.Infoln("UL: ", UL)
	originalUL := &nrf_context.UriList{}
	err := mapstructure.Decode(UL, originalUL)
	if err != nil {
		logger.ManagementLog.Errorln("Decode error in GetNFInstancesProcedure: ", err)
		problemDetail := &models.ProblemDetails{
			Title:  "System failure",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
			Cause:  "SYSTEM_FAILURE",
		}
		return nil, problemDetail
	}
	nrf_context.NnrfUriListLimit(originalUL, limit)
	//c.JSON(http.StatusOK, originalUL)
	return originalUL, nil
}

func NFDeregisterProcedure(nfInstanceID string) (problemDetails *models.ProblemDetails) {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}

	nfProfilesRaw := MongoDBLibrary.RestfulAPIGetMany(collName, filter)
	time.Sleep(time.Duration(1) * time.Second)

	MongoDBLibrary.RestfulAPIDeleteMany(collName, filter)

	// nfProfile data for response
	nfProfiles, err := TimeDecode.Decode(nfProfilesRaw, time.RFC3339)
	if err != nil {
		logger.ManagementLog.Warnln("Time decode error: ", err)
		problemDetails = &models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "NOTIFICATION_ERROR",
			Detail: err.Error(),
		}
		return problemDetails
	}

	uriList := nrf_context.GetNofificationUri(nfProfiles[0])

	//make nfInstanceUri
	localIP := nrf_context.GetServiceIp()
	nfInstanceUri := fmt.Sprintf("%s%s%s%s%s", "https://", localIP, ":29510", "/nnrf-nfm/v1/nf-instances/", nfInstanceID)
	//set info for NotificationData
	Notification_event := models.NotificationEventType_DEREGISTERED

	for _, uri := range uriList {
		problemDetails = SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
		if problemDetails != nil {
			return problemDetails
		}
	}

	return nil
}

func UpdateNFInstanceProcedure(nfInstanceID string, patchJSON []byte) (response map[string]interface{}) {

	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}

	if MongoDBLibrary.RestfulAPIJSONPatch(collName, filter, patchJSON) {
		nf := MongoDBLibrary.RestfulAPIGetOne(collName, filter)

		nfProfilesRaw := []map[string]interface{}{
			nf,
		}

		nfProfiles, err := TimeDecode.Decode(nfProfilesRaw, time.RFC3339)
		if err != nil {
			logger.ManagementLog.Info(err.Error())
		}
		uriList := nrf_context.GetNofificationUri(nfProfiles[0])

		//make nfInstanceUri
		localIP := nrf_context.GetServiceIp()
		uri := fmt.Sprintf("%s%s%s%s%s", "https://", localIP, ":29510", "/nnrf-nfm/v1/nf-instances/", nfInstanceID)
		//set info for NotificationData
		Notification_event := models.NotificationEventType_PROFILE_CHANGED
		nfInstanceUri := uri

		for _, uri := range uriList {
			SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
		}

		return nf
	} else {
		return nil
	}
}

func GetNFInstanceProcedure(nfInstanceID string) (response map[string]interface{}) {
	collName := "NfProfile"
	filter := bson.M{"nfInstanceId": nfInstanceID}
	response = MongoDBLibrary.RestfulAPIGetOne(collName, filter)

	return response
}

func NFRegisterProcedure(nfProfile models.NfProfile) (header http.Header, response bson.M,
	problemDetails *models.ProblemDetails) {
	logger.ManagementLog.Traceln("[NRF] In NFRegisterProcedure")
	var nf models.NfProfile
	localIP := nrf_context.GetServiceIp()

	err := nrf_context.NnrfNFManagementDataModel(&nf, nfProfile)
	if err != nil {
		str1 := fmt.Sprint(nfProfile.HeartBeatTimer)
		problemDetails = &models.ProblemDetails{
			Title:  nfProfile.NfInstanceId,
			Status: http.StatusBadRequest,
			Detail: str1,
		}
		return nil, nil, problemDetails
	}

	//make location header
	locationHeaderValue := nrf_context.SetLocationHeader(nfProfile, localIP)
	//Marshal nf to bson
	tmp, err := json.Marshal(nf)
	if err != nil {
		logger.ManagementLog.Errorln("Marshal error in NFRegisterProcedure: ", err)
	}
	var putData = bson.M{}
	err = json.Unmarshal(tmp, &putData)
	if err != nil {
		logger.ManagementLog.Errorln("Unmarshal error in NFRegisterProcedure: ", err)
	}
	//set db info
	collName := "NfProfile"
	nfInstanceId := nf.NfInstanceId
	filter := bson.M{"nfInstanceId": nfInstanceId}

	// Update NF Profile case
	if MongoDBLibrary.RestfulAPIPutOne(collName, filter, putData) { //true insert
		logger.ManagementLog.Infoln("RestfulAPIPutOne True Insert")
		uriList := nrf_context.GetNofificationUri(nf)

		//set info for NotificationData
		Notification_event := models.NotificationEventType_PROFILE_CHANGED
		nfInstanceUri := locationHeaderValue

		//receive the rsp from handler
		for _, uri := range uriList {
			problemDetails = SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
			if problemDetails != nil {
				return nil, nil, problemDetails
			}
		}

		header = make(http.Header)
		header.Add("Location", locationHeaderValue)
		return header, putData, nil
	} else { // Create NF Profile case
		logger.ManagementLog.Infoln("Create NF Profile")
		uriList := nrf_context.GetNofificationUri(nf)
		//set info for NotificationData
		Notification_event := models.NotificationEventType_REGISTERED
		nfInstanceUri := locationHeaderValue

		for _, uri := range uriList {
			problemDetails = SendNFStatusNotify(Notification_event, nfInstanceUri, uri)
			if problemDetails != nil {
				return nil, nil, problemDetails
			}
		}

		header = make(http.Header)
		header.Add("Location", locationHeaderValue)
		logger.ManagementLog.Infoln("Location header: ", locationHeaderValue)
		return header, putData, nil
	}
}

func SendNFStatusNotify(Notification_event models.NotificationEventType, nfInstanceUri string,
	url string) *models.ProblemDetails {

	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	//url = fmt.Sprintf("%s%s", url, "/notification")

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
