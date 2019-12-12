package udr_producer

import (
	// "context"
	"encoding/json"
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/udr_handler/udr_message"
	"net/http"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	// "strconv"
	// "strings"
)

func HandleCreateAccessAndMobilityData(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.AccessAndMobilityData) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleDeleteAccessAndMobilityData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQueryAccessAndMobilityData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQueryAmData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.amData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	accessAndMobilitySubscriptionData := RestfulAPIGetOne(collName, filter)

	if accessAndMobilitySubscriptionData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, accessAndMobilitySubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleAmfContext3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.contextData.amf3gppAccess"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleCreateAmfContext3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.Amf3GppAccessRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.amf3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAmfContext3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.amf3gppAccess"
	filter := bson.M{"ueId": ueId}

	amf3GppAccessRegistration := RestfulAPIGetOne(collName, filter)

	if amf3GppAccessRegistration != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amf3GppAccessRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleAmfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleCreateAmfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.AmfNon3GppAccessRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.amfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAmfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.amfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	amfNon3GppAccessRegistration := RestfulAPIGetOne(collName, filter)

	if amfNon3GppAccessRegistration != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amfNon3GppAccessRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleModifyAmfSubscriptionInfo(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.contextData.eeSubscriptions.amfSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleModifyAuthentication(respChan chan udr_message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleQueryAuthSubsData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	filter := bson.M{"ueId": ueId}

	authenticationSubscription := RestfulAPIGetOne(collName, filter)

	if authenticationSubscription != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, authenticationSubscription)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateAuthenticationSoR(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.SorData) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.ueUpdateConfirmationData.sorData"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAuthSoR(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.ueUpdateConfirmationData.sorData"
	filter := bson.M{"ueId": ueId}

	sorData := RestfulAPIGetOne(collName, filter)

	if sorData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sorData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateAuthenticationStatus(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.AuthEvent) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.authenticationData.authenticationStatus"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryAuthenticationStatus(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.authenticationData.authenticationStatus"
	filter := bson.M{"ueId": ueId}

	authEvent := RestfulAPIGetOne(collName, filter)

	if authEvent != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, authEvent)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleApplicationDataInfluenceDataGet(respChan chan udr_message.HandlerResponseMessage) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdDelete(respChan chan udr_message.HandlerResponseMessage, influenceId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdPatch(respChan chan udr_message.HandlerResponseMessage, influenceId string, body models.TrafficInfluDataPatch) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataInfluenceIdPut(respChan chan udr_message.HandlerResponseMessage, influenceId string, body models.TrafficInfluData) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifyGet(respChan chan udr_message.HandlerResponseMessage) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifyPost(respChan chan udr_message.HandlerResponseMessage, body models.TrafficInfluSub) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete(respChan chan udr_message.HandlerResponseMessage, subscriptionId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet(respChan chan udr_message.HandlerResponseMessage, subscriptionId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut(respChan chan udr_message.HandlerResponseMessage, subscriptionId string, body models.TrafficInfluSub) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataPfdsAppIdDelete(respChan chan udr_message.HandlerResponseMessage, appId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataPfdsAppIdGet(respChan chan udr_message.HandlerResponseMessage, appId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataPfdsAppIdPut(respChan chan udr_message.HandlerResponseMessage, appId string, body models.PfdDataForApp) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleApplicationDataPfdsGet(respChan chan udr_message.HandlerResponseMessage) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifyPost(respChan chan udr_message.HandlerResponseMessage, body models.ExposureDataSubscription) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifySubIdDelete(respChan chan udr_message.HandlerResponseMessage, subId string) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleExposureDataSubsToNotifySubIdPut(respChan chan udr_message.HandlerResponseMessage, subId string, body models.ExposureDataSubscription) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandlePolicyDataBdtDataBdtReferenceIdDelete(respChan chan udr_message.HandlerResponseMessage, bdtReferenceId string) {
	collName := "policyData.bdtData"
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataBdtDataBdtReferenceIdGet(respChan chan udr_message.HandlerResponseMessage, bdtReferenceId string) {
	collName := "policyData.bdtData"
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	bdtData := RestfulAPIGetOne(collName, filter)

	if bdtData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, bdtData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataBdtDataBdtReferenceIdPut(respChan chan udr_message.HandlerResponseMessage, bdtReferenceId string, body models.BdtData) {
	putData := toBsonM(body)
	putData["bdtReferenceId"] = bdtReferenceId

	collName := "policyData.bdtData"
	filter := bson.M{"bdtReferenceId": bdtReferenceId}

	isExisted := RestfulAPIPutOneNotUpdate(collName, filter, putData)

	if !isExisted {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	} else {
		problemDetails := models.ProblemDetails{
			Cause: "UPDATE_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandlePolicyDataBdtDataGet(respChan chan udr_message.HandlerResponseMessage) {
	collName := "policyData.bdtData"
	filter := bson.M{}

	bdtDataArray := RestfulAPIGetMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, bdtDataArray)
}

func HandlePolicyDataPlmnsPlmnIdUePolicySetGet(respChan chan udr_message.HandlerResponseMessage, plmnId string) {
	collName := "policyData.plmns.uePolicySet"
	filter := bson.M{"plmnId": plmnId}

	uePolicySet := RestfulAPIGetOne(collName, filter)

	if uePolicySet != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, uePolicySet)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataSponsorConnectivityDataSponsorIdGet(respChan chan udr_message.HandlerResponseMessage, sponsorId string) {
	collName := "policyData.sponsorConnectivityData"
	filter := bson.M{"sponsorId": sponsorId}

	sponsorConnectivityData := RestfulAPIGetOne(collName, filter)

	if sponsorConnectivityData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sponsorConnectivityData)
	} else {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	}
}

func HandlePolicyDataSubsToNotifyPost(respChan chan udr_message.HandlerResponseMessage, body models.PolicyDataSubscription) {
	postData := toBsonM(body)

	collName := "policyData.subsToNotify"
	filter := bson.M{}

	RestfulAPIPost(collName, filter, postData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, postData)
}

func HandlePolicyDataSubsToNotifySubsIdDelete(respChan chan udr_message.HandlerResponseMessage, subsId string) {
	collName := "policyData.subsToNotify"
	filter := bson.M{"subsId": subsId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataSubsToNotifySubsIdPut(respChan chan udr_message.HandlerResponseMessage, subsId string, body models.PolicyDataSubscription) {
	putData := toBsonM(body)
	putData["subsId"] = subsId

	collName := "policyData.subsToNotify"
	filter := bson.M{"subsId": subsId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, putData)
}

func HandlePolicyDataUesUeIdAmDataGet(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.amData"
	filter := bson.M{"ueId": ueId}

	amPolicyData := RestfulAPIGetOne(collName, filter)

	if amPolicyData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amPolicyData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataGet(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	operatorSpecificDataContainerMapCover := RestfulAPIGetOne(collName, filter)

	if operatorSpecificDataContainerMapCover != nil {
		operatorSpecificDataContainerMap := operatorSpecificDataContainerMapCover["operatorSpecificDataContainerMap"]
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, operatorSpecificDataContainerMap)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataPatch(respChan chan udr_message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "policyData.ues.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatchExtend(collName, filter, patchJSON, "operatorSpecificDataContainerMap")

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandlePolicyDataUesUeIdOperatorSpecificDataPut(respChan chan udr_message.HandlerResponseMessage, ueId string, body map[string]models.OperatorSpecificDataContainer) {
	// json.NewDecoder(c.Request.Body).Decode(&operatorSpecificDataContainerMap)

	collName := "policyData.ues.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	putData := map[string]interface{}{"operatorSpecificDataContainerMap": body}
	putData["ueId"] = ueId

	RestfulAPIPutOne(collName, filter, putData)

	// for k, v := range operatorSpecificDataContainerMap {
	// 	specificDataElementName := k

	// 	putData := toBsonM(body)
	// 	putData["ueId"] = ueId
	// 	putData["SpecificDataElementName"] = specificDataElementName

	// 	filterTmp := bson.M{"ueId": ueId, "SpecificDataElementName": specificDataElementName}
	// 	RestfulAPIPutOne(collName, filterTmp, putData)
	// }

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandlePolicyDataUesUeIdSmDataGet(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.smData"
	filter := bson.M{"ueId": ueId}

	smPolicyData := RestfulAPIGetOne(collName, filter)
	if smPolicyData != nil {
		var smPolicyDataResp models.SmPolicyData

		err := mapstructure.Decode(smPolicyData, &smPolicyDataResp)
		if err != nil {
			panic(err)
		}

		{
			collName := "policyData.ues.smData.usageMonData"
			filter := bson.M{"ueId": ueId}
			usageMonDataArray := RestfulAPIGetMany(collName, filter)

			if usageMonDataArray != nil {
				var tmp []models.UsageMonData
				err := mapstructure.Decode(usageMonDataArray, &tmp)
				if err != nil {
					panic(err)
				}
				tmp2 := make(map[string]models.UsageMonData)
				for _, element := range tmp {
					tmp2[element.LimitId] = element
				}
				smPolicyDataResp.UmData = tmp2
			}

		}

		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smPolicyDataResp)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdSmDataPatch(respChan chan udr_message.HandlerResponseMessage, ueId string, body map[string]models.UsageMonData) {
	collName := "policyData.ues.smData.usageMonData"
	// filter := bson.M{"ueId": ueId}

	for k, v := range body {
		limitId := k
		filterTmp := bson.M{"ueId": ueId, "LimitId": limitId}
		RestfulAPIMergePatch(collName, filterTmp, toBsonM(v))
	}

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdDelete(respChan chan udr_message.HandlerResponseMessage, ueId string, usageMonId string) {
	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdGet(respChan chan udr_message.HandlerResponseMessage, ueId string, usageMonId string) {
	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	usageMonData := RestfulAPIGetOne(collName, filter)

	if usageMonData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, usageMonData)
	} else {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	}
}

func HandlePolicyDataUesUeIdSmDataUsageMonIdPut(respChan chan udr_message.HandlerResponseMessage, ueId string, usageMonId string, body models.UsageMonData) {
	putData := toBsonM(body)
	putData["ueId"] = ueId
	putData["usageMonId"] = usageMonId

	collName := "policyData.ues.smData.usageMonData"
	filter := bson.M{"ueId": ueId, "usageMonId": usageMonId}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	} else {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
	}
}

func HandlePolicyDataUesUeIdUePolicySetGet(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "policyData.ues.uePolicySet"
	filter := bson.M{"ueId": ueId}

	uePolicySet := RestfulAPIGetOne(collName, filter)

	if uePolicySet != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, uePolicySet)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePolicyDataUesUeIdUePolicySetPatch(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.UePolicySet) {
	patchData := toBsonM(body)
	patchData["ueId"] = ueId

	collName := "policyData.ues.uePolicySet"
	filter := bson.M{"ueId": ueId}

	RestfulAPIMergePatch(collName, filter, patchData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandlePolicyDataUesUeIdUePolicySetPut(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.UePolicySet) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "policyData.ues.uePolicySet"
	filter := bson.M{"ueId": ueId}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	} else {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	}
}

func HandleCreateAMFSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string, body []models.AmfSubscriptionInfo) {
	var putDataArray []map[string]interface{}
	var filterArray []bson.M
	for _, amfSubscriptionInfo := range body {
		putData := toBsonM(amfSubscriptionInfo)
		putData["ueId"] = ueId
		putData["subsId"] = subsId
		putData["AmfInstanceId"] = amfSubscriptionInfo.AmfInstanceId
		putDataArray = append(putDataArray, putData)
		filterArray = append(filterArray, bson.M{"ueId": ueId, "subsId": subsId, "AmfInstanceId": amfSubscriptionInfo.AmfInstanceId})
	}

	collName := "subscriptionData.contextData.eeSubscriptions.amfSubscriptions"

	RestfulAPIPutMany(collName, filterArray, putDataArray)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleRemoveAmfSubscriptionsInfo(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string) {
	collName := "subscriptionData.contextData.eeSubscriptions.amfSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	RestfulAPIDeleteMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryEEData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.eeProfileData"
	filter := bson.M{"ueId": ueId}

	eeProfileData := RestfulAPIGetOne(collName, filter)

	if eeProfileData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, eeProfileData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleRemoveEeGroupSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueGroupId string, subsId string) {
	collName := "subscriptionData.groupData.eeSubscriptions"
	filter := bson.M{"ueGroupId": ueGroupId, "subsId": subsId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleUpdateEeGroupSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueGroupId string, subsId string, body models.EeSubscription) {
	putData := toBsonM(body)
	putData["ueGroupId"] = ueGroupId
	putData["subsId"] = subsId

	collName := "subscriptionData.groupData.eeSubscriptions"
	filter := bson.M{"ueGroupId": ueGroupId, "subsId": subsId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleCreateEeGroupSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueGroupId string, body models.EeSubscription) {
	postData := toBsonM(body)
	postData["ueGroupId"] = ueGroupId

	collName := "subscriptionData.groupData.eeSubscriptions"
	filter := bson.M{"ueGroupId": ueGroupId}

	RestfulAPIPost(collName, filter, postData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, postData)
}

func HandleQueryEeGroupSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueGroupId string) {
	collName := "subscriptionData.groupData.eeSubscriptions"
	filter := bson.M{"ueGroupId": ueGroupId}

	eeSubscriptions := RestfulAPIGetMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, eeSubscriptions)
}

func HandleRemoveeeSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string) {
	collName := "subscriptionData.contextData.eeSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleUpdateEesubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string, body models.EeSubscription) {
	putData := toBsonM(body)
	putData["ueId"] = ueId
	putData["subsId"] = subsId

	collName := "subscriptionData.contextData.eeSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleCreateEeSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.EeSubscription) {
	postData := toBsonM(body)
	postData["ueId"] = ueId

	collName := "subscriptionData.contextData.eeSubscriptions"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPost(collName, filter, postData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, postData)
}

func HandleQueryeesubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.eeSubscriptions"
	filter := bson.M{"ueId": ueId}

	eeSubscriptions := RestfulAPIGetMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, eeSubscriptions)
}

func HandlePatchOperSpecData(respChan chan udr_message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleQueryOperSpecData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.operatorSpecificData"
	filter := bson.M{"ueId": ueId}

	operatorSpecificDataContainer := RestfulAPIGetOne(collName, filter)

	// The key of the map is operator specific data element name and the value is the operator specific data of the UE.

	if operatorSpecificDataContainer != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, operatorSpecificDataContainer)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetppData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.ppData"
	filter := bson.M{"ueId": ueId}

	ppData := RestfulAPIGetOne(collName, filter)

	if ppData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, ppData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateSessionManagementData(respChan chan udr_message.HandlerResponseMessage, ueId string, pduSessionId int32, body models.PduSessionManagementData) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleDeleteSessionManagementData(respChan chan udr_message.HandlerResponseMessage, ueId string, pduSessionId int32) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQuerySessionManagementData(respChan chan udr_message.HandlerResponseMessage, ueId string, pduSessionId int32) {
	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, map[string]interface{}{})
}

func HandleQueryProvisionedData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	var provisionedDataSets models.ProvisionedDataSets

	{
		collName := "subscriptionData.provisionedData.amData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		accessAndMobilitySubscriptionData := RestfulAPIGetOne(collName, filter)
		if accessAndMobilitySubscriptionData != nil {
			var tmp models.AccessAndMobilitySubscriptionData
			err := mapstructure.Decode(accessAndMobilitySubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.AmData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smfSelectionSubscriptionData := RestfulAPIGetOne(collName, filter)
		if smfSelectionSubscriptionData != nil {
			var tmp models.SmfSelectionSubscriptionData
			err := mapstructure.Decode(smfSelectionSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmfSelData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsSubscriptionData := RestfulAPIGetOne(collName, filter)
		if smsSubscriptionData != nil {
			var tmp models.SmsSubscriptionData
			err := mapstructure.Decode(smsSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsSubsData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		sessionManagementSubscriptionDatas := RestfulAPIGetMany(collName, filter)
		if sessionManagementSubscriptionDatas != nil {
			var tmp []models.SessionManagementSubscriptionData
			err := mapstructure.Decode(sessionManagementSubscriptionDatas, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmData = tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.traceData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		traceData := RestfulAPIGetOne(collName, filter)
		if traceData != nil {
			var tmp models.TraceData
			err := mapstructure.Decode(traceData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.TraceData = &tmp
		}
	}

	{
		collName := "subscriptionData.provisionedData.smsMngData"
		filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}
		smsManagementSubscriptionData := RestfulAPIGetOne(collName, filter)
		if smsManagementSubscriptionData != nil {
			var tmp models.SmsManagementSubscriptionData
			err := mapstructure.Decode(smsManagementSubscriptionData, &tmp)
			if err != nil {
				panic(err)
			}
			provisionedDataSets.SmsMngData = &tmp
		}
	}

	if !reflect.DeepEqual(provisionedDataSets, models.ProvisionedDataSets{}) {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, provisionedDataSets)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleModifyPpData(respChan chan udr_message.HandlerResponseMessage, ueId string, patchItem []models.PatchItem) {
	collName := "subscriptionData.ppData"
	filter := bson.M{"ueId": ueId}

	patchJSON, _ := json.Marshal(patchItem)
	success := RestfulAPIJSONPatch(collName, filter, patchJSON)

	if success {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
	} else {
		var problemDetails = models.ProblemDetails{
			Cause: "MODIFY_NOT_ALLOWED",
		}
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
	}
}

func HandleGetAmfSubscriptionInfo(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string) {
	collName := "subscriptionData.contextData.eeSubscriptions.amfSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	amfSubscriptionInfos := RestfulAPIGetMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, amfSubscriptionInfos)
}

func HandleGetIdentityData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.identityData"
	filter := bson.M{"ueId": ueId}

	identityData := RestfulAPIGetOne(collName, filter)

	if identityData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, identityData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetOdbData(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.operatorDeterminedBarringData"
	filter := bson.M{"ueId": ueId}

	operatorDeterminedBarringData := RestfulAPIGetOne(collName, filter)

	if operatorDeterminedBarringData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, operatorDeterminedBarringData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleGetSharedData(respChan chan udr_message.HandlerResponseMessage, sharedDataIds []string) {
	sharedDataIdArray := sharedDataIds

	collName := "subscriptionData.sharedData"
	var sharedDataArray []map[string]interface{}
	for _, sharedDataId := range sharedDataIdArray {
		filter := bson.M{"sharedDataId": sharedDataId}
		sharedData := RestfulAPIGetOne(collName, filter)
		sharedDataArray = append(sharedDataArray, sharedData)
	}

	if sharedDataArray != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sharedDataArray)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleRemovesdmSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string) {
	collName := "subscriptionData.contextData.sdmSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleUpdatesdmsubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, subsId string, body models.SdmSubscription) {
	putData := toBsonM(body)
	putData["ueId"] = ueId
	putData["subsId"] = subsId

	collName := "subscriptionData.contextData.sdmSubscriptions"
	filter := bson.M{"ueId": ueId, "subsId": subsId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleCreateSdmSubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.SdmSubscription) {
	postData := toBsonM(body)
	postData["ueId"] = ueId

	collName := "subscriptionData.contextData.sdmSubscriptions"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPost(collName, filter, postData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, postData)
}

func HandleQuerysdmsubscriptions(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.sdmSubscriptions"
	filter := bson.M{"ueId": ueId}

	sdmSubscriptions := RestfulAPIGetMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sdmSubscriptions)
}

func HandleQuerySmData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string, singleNssai models.Snssai, dnn string) {
	collName := "subscriptionData.provisionedData.smData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	if !reflect.DeepEqual(singleNssai, models.Snssai{}) {
		if singleNssai.Sd == "" {
			filter["singleNssai.sst"] = singleNssai.Sst
		} else {
			filter["singleNssai.sst"] = singleNssai.Sst
			filter["singleNssai.sd"] = singleNssai.Sd
		}
	}

	if dnn != "" {
		filter["dnnConfigurations."+dnn] = bson.M{"$exists": true}
	}

	sessionManagementSubscriptionDatas := RestfulAPIGetMany(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, sessionManagementSubscriptionDatas)
}

func HandleCreateSmfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, pduSessionId int32, body models.SmfRegistration) {
	pduSessionIdInt := pduSessionId

	putData := toBsonM(body)
	putData["ueId"] = ueId
	putData["pduSessionId"] = int32(pduSessionIdInt)

	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	isExisted := RestfulAPIPutOne(collName, filter, putData)

	if !isExisted {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, putData)
	} else {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, putData)
	}
}

func HandleDeleteSmfContext(respChan chan udr_message.HandlerResponseMessage, ueId string, pduSessionId string) {
	pduSessionIdInt, _ := strconv.ParseInt(pduSessionId, 10, 32)

	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQuerySmfRegistration(respChan chan udr_message.HandlerResponseMessage, ueId string, pduSessionId string) {
	pduSessionIdInt, _ := strconv.ParseInt(pduSessionId, 10, 32)

	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId, "pduSessionId": pduSessionIdInt}

	smfRegistration := RestfulAPIGetOne(collName, filter)

	if smfRegistration != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smfRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleQuerySmfRegList(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smfRegistrations"
	filter := bson.M{"ueId": ueId}

	smfRegList := RestfulAPIGetMany(collName, filter)

	if smfRegList == nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smfRegList)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}

}

func HandleQuerySmfSelectData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smfSelectionSubscriptionData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smfSelectionSubscriptionData := RestfulAPIGetOne(collName, filter)

	if smfSelectionSubscriptionData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smfSelectionSubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateSmsfContext3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.SmsfRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleDeleteSmsfContext3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQuerySmsfContext3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"ueId": ueId}

	smsfRegistration := RestfulAPIGetOne(collName, filter)

	if smsfRegistration != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsfRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleCreateSmsfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string, body models.SmsfRegistration) {
	putData := toBsonM(body)
	putData["ueId"] = ueId

	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIPutOne(collName, filter, putData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleDeleteSmsfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQuerySmsfContextNon3gpp(respChan chan udr_message.HandlerResponseMessage, ueId string) {
	collName := "subscriptionData.contextData.smsfNon3gppAccess"
	filter := bson.M{"ueId": ueId}

	smsfRegistration := RestfulAPIGetOne(collName, filter)

	if smsfRegistration != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsfRegistration)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleQuerySmsMngData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smsMngData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smsManagementSubscriptionData := RestfulAPIGetOne(collName, filter)

	if smsManagementSubscriptionData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsManagementSubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleQuerySmsData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.smsData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	smsSubscriptionData := RestfulAPIGetOne(collName, filter)

	if smsSubscriptionData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, smsSubscriptionData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}

func HandlePostSubscriptionDataSubscriptions(respChan chan udr_message.HandlerResponseMessage, body models.SubscriptionDataSubscriptions) {
	postData := toBsonM(body)

	collName := "subscriptionData.contextData.sdmSubscriptions"
	filter := bson.M{"ueId": body.UeId}

	RestfulAPIPost(collName, filter, postData)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, postData)
}

func HandleRemovesubscriptionDataSubscriptions(respChan chan udr_message.HandlerResponseMessage, subsId string) {
	collName := "subscriptionData.contextData.smsf3gppAccess"
	filter := bson.M{"subsId": subsId}

	RestfulAPIDeleteOne(collName, filter)

	udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNoContent, map[string]interface{}{})
}

func HandleQueryTraceData(respChan chan udr_message.HandlerResponseMessage, ueId string, servingPlmnId string) {
	collName := "subscriptionData.provisionedData.traceData"
	filter := bson.M{"ueId": ueId, "servingPlmnId": servingPlmnId}

	traceData := RestfulAPIGetOne(collName, filter)

	if traceData != nil {
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, traceData)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		udr_message.SendHttpResponseMessage(respChan, nil, http.StatusNotFound, problemDetails)
	}
}
