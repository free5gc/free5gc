package udm_producer

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler/udm_message"
	"net/http"
)

func HandleCreateEeSubscription(httpChannel chan udm_message.HandlerResponseMessage, ueIdentity string, subscriptionID string, eesubscription models.EeSubscription) {

	var body models.CreatedEeSubscription
	udm_context.CreateEeSusbContext(ueIdentity, body)
	clientAPI := createUDMClientToUDR(ueIdentity, false)
	eeSubscriptionResp, res, err := clientAPI.EventExposureSubscriptionsCollectionApi.CreateEeSubscriptions(context.Background(),
		ueIdentity, eesubscription)
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
		udmue := udm_context.CreateUdmUe(ueIdentity)
		udmue.CreatedEeSubscription.EeSubscription = &eeSubscriptionResp
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusCreated, udmue.CreatedEeSubscription.EeSubscription)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	}
}

func HandleDeleteEeSubscription(httpChannel chan udm_message.HandlerResponseMessage, ueIdentity string, subscriptionID string) {

	clientAPI := createUDMClientToUDR(ueIdentity, false)
	res, err := clientAPI.EventExposureSubscriptionDocumentApi.RemoveeeSubscriptions(context.Background(), ueIdentity, subscriptionID)
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

func HandleUpdateEeSubscription(httpChannel chan udm_message.HandlerResponseMessage, ueIdentity string, subscriptionID string) {

	clientAPI := createUDMClientToUDR(ueIdentity, false)
	patchItem := models.PatchItem{}
	patchItem.Value = models.EeSubscription{}
	body := Nudr_DataRepository.UpdateEesubscriptionsParamOpts{
		EeSubscription: optional.NewInterface(patchItem.Value),
	}
	res, err := clientAPI.EventExposureSubscriptionDocumentApi.UpdateEesubscriptions(context.Background(), ueIdentity, subscriptionID, &body)
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
	} else if res.StatusCode == http.StatusNotFound {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "DATA_NOT_FOUND"
		problemDetails.Status = 404
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNotFound, problemDetails)
	} else {
		var problemDetails models.ProblemDetails
		problemDetails.Status = 403
		udm_message.SendHttpResponseMessage(httpChannel, nil, http.StatusForbidden, problemDetails)
	}
}
