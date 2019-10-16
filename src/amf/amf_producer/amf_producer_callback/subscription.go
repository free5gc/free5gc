package amf_producer_callback

import (
	"context"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"reflect"
)

func SendAmfStatusChangeNotify(amfStatus string, guamiList []models.Guami) {
	amfSelf := amf_context.AMF_Self()

	for _, amfStatusSubscriptions := range amfSelf.AMFStatusSubscriptions {

		configuration := Namf_Communication.NewConfiguration()
		client := Namf_Communication.NewAPIClient(configuration)
		amfStatusNotification := models.AmfStatusChangeNotification{}
		var amfStatusInfo = models.AmfStatusInfo{}

		for _, guami := range guamiList {
			for _, subGumi := range amfStatusSubscriptions.GuamiList {
				if reflect.DeepEqual(guami, subGumi) {
					//AMF status is available
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
		uri := amfStatusSubscriptions.AmfStatusUri

		httpResponse, err := client.AmfStatusChangeCallbackDocumentApiServiceCallbackDocumentApi.AmfStatusChangeNotify(context.Background(), uri, amfStatusNotification)
		if err != nil {
			if httpResponse == nil {
				HttpLog.Errorln(err.Error())
			} else if err.Error() != httpResponse.Status {
				HttpLog.Errorln(err.Error())
			}
		}
	}
}
