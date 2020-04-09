package udm_producer_callback

import (
	"context"
	"free5gc/lib/Nudm_SubscriberDataManagement"
	"free5gc/lib/Nudm_UEContextManagement"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"
)

func SendOnDataChangeNotification(ueId string, notifyItems []models.NotifyItem) {
	udmSelf := udm_context.UDM_Self().UdmUePool[ueId]
	configuration := Nudm_SubscriberDataManagement.NewConfiguration()
	clientAPI := Nudm_SubscriberDataManagement.NewAPIClient(configuration)

	for _, subscriptionDataSubscription := range udmSelf.UdmSubsToNotify {
		onDataChangeNotificationurl := subscriptionDataSubscription.OriginalCallbackReference
		dataChangeNotification := models.ModificationNotification{}
		dataChangeNotification.NotifyItems = notifyItems
		httpResponse, err := clientAPI.DataChangeNotificationCallbackDocumentApi.OnDataChangeNotification(context.TODO(), onDataChangeNotificationurl, dataChangeNotification)
		if err != nil {
			if httpResponse == nil {
				logger.HttpLog.Error(err.Error())
			} else if err.Error() != httpResponse.Status {
				logger.HttpLog.Errorln(err.Error())
			}
		}
	}
}

func SendOnDeregistrationNotification(ueId string, onDeregistrationNotificationUrl string, deregistData models.DeregistrationData) {
	configuration := Nudm_UEContextManagement.NewConfiguration()
	clientAPI := Nudm_UEContextManagement.NewAPIClient(configuration)

	httpResponse, err := clientAPI.DeregistrationNotificationCallbackApi.DeregistrationNotify(context.TODO(), onDeregistrationNotificationUrl, deregistData)
	if err != nil {
		if httpResponse == nil {
			logger.HttpLog.Error(err.Error())
		} else if err.Error() != httpResponse.Status {
			logger.HttpLog.Errorln(err.Error())
		}
	}
}
