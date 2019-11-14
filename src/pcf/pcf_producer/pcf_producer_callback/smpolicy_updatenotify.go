package pcf_producer_callback

import (
	"context"
	"fmt"
	"free5gc/lib/Npcf_SMPolicy"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_util"
)

// SmPoliciesSmPolicyUpdateNotify -
func CreateSmPolicyNotifyContext(id string, send_type string, policydecision *models.SmPolicyDecision) {
	resourceURI := pcf_util.PCF_BASIC_PATH + pcf_context.SmpolicyUri + id
	var smPolicyNotification models.SmPolicyNotification
	var terminationNotification models.TerminationNotification
	var url string
	pcfUeContext := pcf_context.PCF_Self().UePool
	configuration := Npcf_SMPolicy.NewConfiguration()
	client := Npcf_SMPolicy.NewAPIClient(configuration)
	for key := range pcfUeContext {
		if pcfUeContext[key].SmPolicyControlStore == nil {
			continue
		}
		idTemp := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId)
		if id == idTemp {

			url = pcfUeContext[key].SmPolicyControlStore.Context.NotificationUri
			if send_type == "update" {
				smPolicyNotification.ResourceUri = resourceURI + "/update"
				smPolicyNotification.SmPolicyDecision = policydecision
				_, err := client.NotifyApi.SMNotificationUri(context.Background(), url, smPolicyNotification)

				if err != nil {
					logger.SMpolicylog.Warnln("SMPolicy UpdateNotify POST error: ", err)
				}
				return

			} else if send_type == "terminate" {
				terminationNotification.ResourceUri = resourceURI + "/delete"
				terminationNotification.Cause = "UNSPECIFIED"
				_, err := client.NotifyApi.SMTerminationUri(context.Background(), url, terminationNotification)

				if err != nil {
					logger.SMpolicylog.Warnln("SMPolicy UpdateNotify POST error: ", err)
				}
				return
			} else {
				return
			}
		}
	}
}
