package amf_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/amf_handler/amf_message"
	"net/http"
	"reflect"
	"strconv"
)

func HandleAMFStatusChangeSubscribeRequest(httpChannel chan amf_message.HandlerResponseMessage, body models.SubscriptionData) {
	var response models.SubscriptionData
	var problem models.ProblemDetails
	var guami models.Guami
	amfSelf := amf_context.AMF_Self()

	for _, guami = range body.GuamiList {
		for _, servedGumi := range amfSelf.ServedGuamiList {
			if reflect.DeepEqual(guami, servedGumi) {
				//AMF status is available
				response.GuamiList = append(response.GuamiList, guami)
			}
		}
	}

	if response.GuamiList != nil {
		newSubscriptionID := strconv.Itoa(amfSelf.AMFStatusSubscriptionIDGenerator)
		amfSelf.AMFStatusSubscriptions[newSubscriptionID] = new(models.SubscriptionData)
		locationHeader := body.AmfStatusUri + "/" + newSubscriptionID
		headers := http.Header{
			"Location": {locationHeader},
		}
		amfSelf.AMFStatusSubscriptions[newSubscriptionID].AmfStatusUri = locationHeader
		amfSelf.AMFStatusSubscriptions[newSubscriptionID].GuamiList = response.GuamiList
		amfSelf.AMFStatusSubscriptionIDGenerator++
		amf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusCreated, response)
	} else {
		problem.Status = 403
		problem.Cause = "UNSPECIFIED"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
	}
}

func HandleAMFStatusChangeUnSubscribeRequest(httpChannel chan amf_message.HandlerResponseMessage, subscriptionId string) {
	var problem models.ProblemDetails
	amfSelf := amf_context.AMF_Self()
	_, ok := amfSelf.AMFStatusSubscriptions[subscriptionId]

	if !ok {
		problem.Status = 403
		problem.Cause = "SUBSCRIPTION_NOT_FOUND "
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problem)
	} else {
		delete(amfSelf.AMFStatusSubscriptions, subscriptionId)
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
	}
}

func HandleAMFStatusChangeSubscribeModfy(httpChannel chan amf_message.HandlerResponseMessage, subscriptionId string, body models.SubscriptionData) {
	var problem models.ProblemDetails
	var response models.SubscriptionData
	amfSelf := amf_context.AMF_Self()
	_, ok := amfSelf.AMFStatusSubscriptions[subscriptionId]
	if !ok {
		problem.Status = 403
		problem.Cause = "Forbidden"
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problem)
	} else {
		amfGuamiList := amfSelf.AMFStatusSubscriptions[subscriptionId].GuamiList
		// clear GuamiList
		amfGuamiList = amfGuamiList[0:0]
		for _, guamiList := range body.GuamiList {
			amfGuamiList = append(amfGuamiList, guamiList)
			response.GuamiList = append(response.GuamiList, guamiList)
		}

		amfSelf.AMFStatusSubscriptions[subscriptionId].AmfStatusUri = body.AmfStatusUri
		response.AmfStatusUri = body.AmfStatusUri
		amf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusAccepted, response)
	}
}
