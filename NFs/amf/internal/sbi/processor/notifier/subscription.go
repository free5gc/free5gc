package callback

import (
	"context"
	"reflect"

	amf_context "github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	Namf_Communication "github.com/free5gc/openapi/amf/Communication"
	"github.com/free5gc/openapi/models"
)

func SendAmfStatusChangeNotify(amfStatus string, guamiList []models.Guami) {
	amfSelf := amf_context.GetSelf()

	amfSelf.AMFStatusSubscriptions.Range(func(key, value interface{}) bool {
		subscriptionData := value.(models.AmfCommunicationSubscriptionData)

		configuration := Namf_Communication.NewConfiguration()
		client := Namf_Communication.NewAPIClient(configuration)
		amfStatusNotification := models.AmfStatusChangeNotification{}
		amfStatusInfo := models.AmfStatusInfo{}

		for _, guami := range guamiList {
			for _, subGumi := range subscriptionData.GuamiList {
				if reflect.DeepEqual(guami, subGumi) {
					// AMF status is available
					amfStatusInfo.GuamiList = append(amfStatusInfo.GuamiList, guami)
				}
			}
		}

		amfStatusInfo = models.AmfStatusInfo{
			StatusChange:     (models.StatusChange)(amfStatus),
			TargetAmfRemoval: "",
			TargetAmfFailure: "",
		}

		amfStatusNotification.AmfStatusInfoList = append(amfStatusNotification.AmfStatusInfoList, amfStatusInfo)
		uri := subscriptionData.AmfStatusUri

		amfStatusNotificationReq := Namf_Communication.AmfStatusChangeNotifyRequest{
			AmfStatusChangeNotification: &amfStatusNotification,
		}
		logger.ProducerLog.Infof("[AMF] Send Amf Status Change Notify to %s", uri)
		_, err := client.IndividualSubscriptionDocumentApi.
			AmfStatusChangeNotify(context.Background(), uri, &amfStatusNotificationReq)
		if err != nil {
			HttpLog.Errorln(err.Error())
		}
		return true
	})
}
