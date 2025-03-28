package TestComm

import (
	"github.com/free5gc/openapi/models"
)

const (
	AMFStatusSubscription403      = "AMFStatusSubscription403"
	AMFStatusSubscription201      = "AMFStatusSubscription201"
	AMFStatusUnSubscription403    = "AMFStatusUnSubscription403"
	AMFStatusUnSubscription204    = "AMFStatusUnSubscription204"
	AMFStatusSubscriptionModfy403 = "AMFStatusSubscriptionModfy403"
	AMFStatusSubscriptionModfy200 = "AMFStatusSubscriptionModfy200"
)

var ConsumerAMFStatusSubscriptionTable = make(map[string]models.AmfCommunicationSubscriptionData)

func init() {
	ConsumerAMFStatusSubscriptionTable[AMFStatusSubscription403] = models.AmfCommunicationSubscriptionData{
		AmfStatusUri: "",
		GuamiList:    nil,
	}

	ConsumerAMFStatusSubscriptionTable[AMFStatusSubscription201] = models.AmfCommunicationSubscriptionData{
		AmfStatusUri: "https://127.0.0.1:29333/AMFStatusNotify",
		GuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnIdNid{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "cafe00",
			},
		},
	}
}

var ConsumerAMFStatusUnSubscriptionTable = make(map[string]string)

func init() {
	ConsumerAMFStatusUnSubscriptionTable[AMFStatusUnSubscription403] = "0"
	ConsumerAMFStatusUnSubscriptionTable[AMFStatusUnSubscription204] = "1"
}

var ConsumerAMFStatusChangeSubscribeModfyTable = make(map[string]models.AmfCommunicationSubscriptionData)

func init() {
	ConsumerAMFStatusChangeSubscribeModfyTable[AMFStatusSubscriptionModfy403] = models.AmfCommunicationSubscriptionData{
		AmfStatusUri: "",
		GuamiList:    nil,
	}

	ConsumerAMFStatusChangeSubscribeModfyTable[AMFStatusSubscriptionModfy200] = models.AmfCommunicationSubscriptionData{
		AmfStatusUri: "https://127.0.0.1:29333/AMFStatusNotify/1",
		GuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnIdNid{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "cafe00",
			},
		},
	}
}
