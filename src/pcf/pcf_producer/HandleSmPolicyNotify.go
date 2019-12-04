package pcf_producer

import (
	"context"
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"

	"github.com/gin-gonic/gin"
)

func HandleSmPolicyNotify(httpChannel chan pcf_message.HttpResponseMessage, supi string, body models.PolicyDataChangeNotification) {
	policyDataChangeNotification := body
	var problem models.ProblemDetails
	var smPolicyDecision models.SmPolicyDecision
	pcfUeContext := pcf_context.PCF_Self().UePool
	counter := false
	client := pcf_util.GetNudrClient("https://localhost:29504")
	UeContext := pcf_context.PCF_Self().UePool
	if UeContext[supi] == nil {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
	}
	_, resp, err := client.DefaultApi.PolicyDataUesUeIdSmDataGet(context.Background(), supi, nil)
	if err != nil {
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
	}
	if resp.StatusCode == 204 {
		UeContext[supi].SmPolicyControlStore = nil
		// pcf_producer_callback.CreateSmPolicyNotifyContext(fmt.Sprint(UeContext[supi].SmPolicyControlStore.Context.PduSessionId), "terminate", nil)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, nil)
	}
	if resp.StatusCode == 200 {
		for key := range pcfUeContext {

			SupiTemp := fmt.Sprint(pcfUeContext[key].Supi)
			if supi == SupiTemp {
				snssai := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.SliceInfo.Sst) + pcfUeContext[key].SmPolicyControlStore.Context.SliceInfo.Sd
				dnn := pcfUeContext[key].SmPolicyControlStore.Context.Dnn
				if policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Ipv4Index != 0 {
					smPolicyDecision.Ipv4Index = policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Ipv4Index
					pcfUeContext[key].SmPolicyControlStore.Policy.Ipv4Index = smPolicyDecision.Ipv4Index
					counter = true
				}
				if policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Ipv6Index != 0 {
					smPolicyDecision.Ipv6Index = policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Ipv6Index
					pcfUeContext[key].SmPolicyControlStore.Policy.Ipv6Index = smPolicyDecision.Ipv6Index
					counter = true
				}
				if pcfUeContext[key].SmPolicyControlStore.Policy.Online != policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Online {
					smPolicyDecision.Online = policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Online
					pcfUeContext[key].SmPolicyControlStore.Policy.Online = smPolicyDecision.Online
					counter = true
				}
				if pcfUeContext[key].SmPolicyControlStore.Policy.Offline != policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Offline {
					smPolicyDecision.Offline = policyDataChangeNotification.SmPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[dnn].Offline
					pcfUeContext[key].SmPolicyControlStore.Policy.Offline = smPolicyDecision.Offline
					counter = true
				}
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, gin.H{})
				if counter {
					// pcf_producer_callback.CreateSmPolicyNotifyContext(fmt.Sprint(UeContext[supi].SmPolicyControlStore.Context.PduSessionId), "update", &smPolicyDecision)
					pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, nil)
				}

			}

		}
	}
}
