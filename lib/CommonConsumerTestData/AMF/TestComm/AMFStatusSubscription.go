package TestComm

import (
	"free5gc/lib/openapi/models"
)

const (
	AMFStatusSubscription403      = "AMFStatusSubscription403"
	AMFStatusSubscription200      = "AMFStatusSubscription200"
	AMFStatusUnSubscription403    = "AMFStatusUnSubscription403"
	AMFStatusUnSubscription204    = "AMFStatusUnSubscription204"
	AMFStatusSubscriptionModfy403 = "AMFStatusSubscriptionModfy403"
	AMFStatusSubscriptionModfy200 = "AMFStatusSubscriptionModfy200"
)

var ConsumerAMFStatusSubscriptionTable = make(map[string]models.SubscriptionData)

func init() {
	ConsumerAMFStatusSubscriptionTable[AMFStatusSubscription403] = models.SubscriptionData{
		AmfStatusUri: "",
		GuamiList:    nil,
	}

	ConsumerAMFStatusSubscriptionTable[AMFStatusSubscription200] = models.SubscriptionData{
		AmfStatusUri: "https://127.0.0.1:29333/AMFStatusNotify",
		GuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnId{
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

var ConsumerAMFStatusChangeSubscribeModfyTable = make(map[string]models.SubscriptionData)

func init() {
	ConsumerAMFStatusChangeSubscribeModfyTable[AMFStatusSubscriptionModfy403] = models.SubscriptionData{
		AmfStatusUri: "",
		GuamiList:    nil,
	}

	ConsumerAMFStatusChangeSubscribeModfyTable[AMFStatusSubscriptionModfy200] = models.SubscriptionData{
		AmfStatusUri: "https://127.0.0.1:29333/AMFStatusNotify/1",
		GuamiList: []models.Guami{
			{
				PlmnId: &models.PlmnId{
					Mcc: "208",
					Mnc: "93",
				},
				AmfId: "cafe00",
			},
		},
	}
}
