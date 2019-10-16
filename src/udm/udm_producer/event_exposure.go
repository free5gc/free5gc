package udm_producer

import (
	"context"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	m "free5gc/lib/openapi/models"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
)

func HandleCreateEeSubscription(httpChannel chan udm_message.HandlerResponseMessage, ueIdentity string, subscriptionID string, eesubscription m.EeSubscription) {

	clientAPI := createUDMClientToUDR(ueIdentity, false)
	eeSubscriptionResp, res, err := clientAPI.EventExposureSubscriptionsCollectionApi.CreateEeSubscriptions(context.Background(),
		ueIdentity, eesubscription)
	if err != nil {
		var problemDetails m.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusCreated, eeSubscriptionResp)
}
func HandleDeleteEeSubscription(httpChannel chan udm_message.HandlerResponseMessage, ueIdentity string, subscriptionID string) {

	clientAPI := createUDMClientToUDR(ueIdentity, false)
	res, err := clientAPI.EventExposureSubscriptionDocumentApi.RemoveeeSubscriptions(context.Background(), ueIdentity, subscriptionID)
	if err != nil {
		var problemDetails m.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}

func HandleUpdateEeSubscription(httpChannel chan udm_message.HandlerResponseMessage, ueIdentity string, subscriptionID string) {

	clientAPI := createUDMClientToUDR(ueIdentity, false)
	res, err := clientAPI.EventExposureSubscriptionDocumentApi.UpdateEesubscriptions(context.Background(), ueIdentity, subscriptionID, nil)
	if err != nil {
		var problemDetails m.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(httpChannel, nil, res.StatusCode, problemDetails)
		return
	}
	udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}
