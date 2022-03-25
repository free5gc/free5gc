package TestComm

import (
	"github.com/free5gc/openapi/models"
)

const (
	UeN1N2InfoSubsriptionCreateData = "UeN1N2InfoSubsriptionCreateData"
)

var ConsumerAMFN1N2MessageSubscribeRequsetTable = make(map[string]*models.UeN1N2InfoSubscriptionCreateData)

func init() {
	ConsumerAMFN1N2MessageSubscribeRequsetTable[UeN1N2InfoSubsriptionCreateData] = &models.UeN1N2InfoSubscriptionCreateData{
		N2NotifyCallbackUri: "https://127.0.0.1:29503/n2InfoNotify",
		N2InformationClass:  models.N2InformationClass_NRP_PA,
		N1NotifyCallbackUri: "https://127.0.0.1:29507/n1MessageNotify",
		N1MessageClass:      models.N1MessageClass__5_GMM,
		NfId:                "AmfId",
	}
}
