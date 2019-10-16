package pcf_producer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"

	"github.com/jinzhu/copier"

	"github.com/antihax/optional"
)

func DeletePoliciesPolAssoId(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	var problem models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()
	for key := range pcfUeContext {
		var URITemp = pcf_context.AmpolicyUri + pcfUeContext[key].PolAssociationIDStore.PolAssoId
		if URI == URITemp {
			var subsId = pcfUeContext[key].PolAssociationIDStore.PolAssoId
			delete(pcfUeContext, key)
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, nil)
			client := pcf_util.GetNudrClient()
			_, err := client.DefaultApi.PolicyDataSubsToNotifySubsIdDelete(context.Background(), subsId)

			if err == nil {
				logger.AMpolicylog.Println("unsubscribition ok")
			}
			return
		}

		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
		return
	}
}

// PoliciesPolAssoIdGet -
func GetPoliciesPolAssoId(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	var problem models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()
	logger.AMpolicylog.Traceln("ampolicy request uri", URI)
	for key := range pcfUeContext {
		if pcfUeContext[key].PolAssociationIDStore == nil {
			continue
		}
		var URITemp = pcf_context.AmpolicyUri + pcfUeContext[key].PolAssociationIDStore.PolAssoId
		if URI == URITemp {
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, pcfUeContext[key].PolAssociationIDStore.PolAssoidTemp)
			return
		}
		problem.Status = 404
		problem.Cause = "CONTEXT_NOT_FOUND"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
		return
	}

}
func UpdatePostPoliciesPolAssoId(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.PolicyAssociationUpdateRequest) {
	var policyAssociationUpdateRequest models.PolicyAssociationUpdateRequest = body
	var policyUpdate models.PolicyUpdate
	var problem models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()
	for key := range pcfUeContext {
		var URITemp = pcf_context.AmpolicyUri + pcfUeContext[key].PolAssociationIDStore.PolAssoId + "/update"
		if URI == URITemp {
			for triggersindex := range policyAssociationUpdateRequest.Triggers {
				if policyAssociationUpdateRequest.Triggers[triggersindex] == "LOC_CH" {
					if err := copier.Copy(&pcfUeContext[key].PolAssociationIDStore.PolAssoidTemp.Request.UserLoc, &policyAssociationUpdateRequest.UserLoc); err != nil {
						logger.AMpolicylog.Warnln("Copy LOC_CH fail: ", err)
					}
				}
				if policyAssociationUpdateRequest.Triggers[triggersindex] == "PRA_CH" {
					if err := copier.Copy(&pcfUeContext[key].PolAssociationIDStore.PolAssoidTemp.Pras, &policyAssociationUpdateRequest.PraStatuses); err != nil {
						logger.AMpolicylog.Warnln("Copy PRA_CH fail: ", err)
					}
				}
				if policyAssociationUpdateRequest.Triggers[triggersindex] == "SERV_AREA_CH" {
					if err := copier.Copy(&pcfUeContext[key].PolAssociationIDStore.PolAssoidTemp.ServAreaRes, policyAssociationUpdateRequest.ServAreaRes); err != nil {
						logger.AMpolicylog.Warnln("Copy SERV_AREA_CH fail: ", err)
					}
				}
				if policyAssociationUpdateRequest.Triggers[triggersindex] == "RFSP_CH" {
					if err := copier.Copy(&pcfUeContext[key].PolAssociationIDStore.PolAssoidTemp.Rfsp, &policyAssociationUpdateRequest.Rfsp); err != nil {
						logger.AMpolicylog.Warnln("Copy RFSP_CH fail: ", err)

					}
				}
			}
			policyUpdate.Triggers = &policyAssociationUpdateRequest.Triggers
			policyUpdate.ResourceUri = URI
			pcfUeContext[key].PolAssociationIDStore.PolAssoidUpdateTemp = policyUpdate
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, policyUpdate)
			return
		}
		if URI != URITemp {
			problem.Status = 404
			problem.Cause = "CONTEXT_NOT_FOUND"
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problem)
			return
		}
	}
}

// PoliciesPost -
func PostPolicies(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.PolicyAssociationRequest) {
	var policyAssociationRequest models.PolicyAssociationRequest = body
	var policyDataSubsToNotifyPostParamOpts Nudr_DataRepository.PolicyDataSubsToNotifyPostParamOpts
	var policyDataSubscription models.PolicyDataSubscription
	var subscriptiondata models.SubscriptionData
	var guami models.Guami
	var problem models.ProblemDetails
	pcfUeContext := pcf_context.GetPCFUeContext()
	if policyAssociationRequest.NotificationUri != "" && policyAssociationRequest.Supi != "" && policyAssociationRequest.SuppFeat != "" {
		for key := range pcfUeContext {
			if pcfUeContext[key].PolAssociationIDStore.PolAssoidDataSubscriptionTemp.NotificationUri == "" {
				ueId := policyAssociationRequest.Supi
				client := pcf_util.GetNudrClient()
				amPolicyData, _, err := client.DefaultApi.PolicyDataUesUeIdAmDataGet(context.Background(), ueId)

				if err == nil {
					pcfUeContext[key].PolAssociationIDStore.PolAssoidSubcCatsTemp = amPolicyData
				}
				if pcfUeContext[key].PolAssociationIDStore.PolAssoidSubcCatsTemp.SubscCats == nil {
					logger.AMpolicylog.Warnln("Nudr_DataRepository Amdata SubscCats is empty")
				} else {
					return
				}
			}
		}

		policyDataSubscription.NotificationUri = pcf_context.NotifiUri
		policyDataSubscription.MonitoredResourceUris = []string{policyAssociationRequest.NotificationUri}
		policyDataSubscription.SupportedFeatures = policyAssociationRequest.SuppFeat

		policyDataSubsToNotifyPostParamOpts.PolicyDataSubscription = optional.NewInterface(policyDataSubscription)
		client := pcf_util.GetNudrClient()
		policydatasubscription, _, err := client.DefaultApi.PolicyDataSubsToNotifyPost(context.Background(), &policyDataSubsToNotifyPostParamOpts)

		if err == nil {
			pcfUeContext[policyAssociationRequest.Supi].PolAssociationIDStore.PolAssoidDataSubscriptionTemp = policydatasubscription
		}
		supi := policyAssociationRequest.Supi
		if err := pcf_context.NewPCFUe(supi); err != nil {
			logger.AMpolicylog.Warnln("Create pcf ue context fail")
		}

		polAssociationContext := pcf_context.PolAssociationIDStore{
			PolAssoId: policyAssociationRequest.Supi,
			PolAssoidTemp: models.PolicyAssociation{
				SuppFeat: "suppfeat",
				Request:  policyAssociationRequest,
			},
		}

		pcfUeContext[supi].PolAssociationIDStore = &polAssociationContext

		if policyAssociationRequest.Guami != guami {
			client := pcf_util.GetNamfClient()
			subscriptiondata, _, err := client.SubscriptionsCollectionDocumentApi.AMFStatusChangeSubscribe(context.Background(), subscriptiondata)
			if err == nil {
				pcfUeContext[supi].PolAssociationIDStore.PolAssoidSubscriptiondataTemp = subscriptiondata
			}
		}

		logger.AMpolicylog.Traceln("ampolicy association id", pcfUeContext[supi].PolAssociationIDStore.PolAssoId)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, pcfUeContext[supi].PolAssociationIDStore)
		return

	} else {
		problem.Status = 400
		problem.Cause = "USER_UNKNOWN"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 400, problem)
		return
	}
}

func AMPolicyUpdateNotification(id string, send_type string) {
	var policyUpdate models.PolicyUpdate
	var terminationNotification models.TerminationNotification
	pcfUeContext := pcf_context.GetPCFUeContext()
	for key := range pcfUeContext {
		idTemp := fmt.Sprint(pcfUeContext[key].PolAssociationIDStore.PolAssoId)
		if id == idTemp {

			updateurl := pcf_util.PCF_BASIC_PATH + pcfUeContext[key].PolAssociationIDStore.PolAssoidUpdateTemp.ResourceUri + "/update"
			terminateurl := pcf_util.PCF_BASIC_PATH + pcfUeContext[key].PolAssociationIDStore.PolAssoidUpdateTemp.ResourceUri + "/terminate"
			if send_type == "update" {
				policyUpdate.ResourceUri = updateurl
				policyUpdate = pcfUeContext[key].PolAssociationIDStore.PolAssoidUpdateTemp

				bs, err := json.Marshal(&policyUpdate)
				if err != nil {
					logger.AMpolicylog.Warnln("Marshal error message is : ", err)
				}
				req2 := bytes.NewBuffer([]byte(bs))
				req, err := http.NewRequest("POST", updateurl, req2)
				if err != nil {
					logger.AMpolicylog.Warnln("Namf Update fail error message is : ", err)
				}
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				_, err = client.Do(req)
				if err != nil {
					logger.AMpolicylog.Warnln("Namf Update fail error message is : ", err)
				}
				return

			} else if send_type == "terminate" {
				terminationNotification.ResourceUri = terminateurl
				terminationNotification.Cause = "UNSPECIFIED"

				bs, err := json.Marshal(&terminationNotification)
				if err != nil {
					fmt.Printf("JSON Marshal error: %v", err)
				}
				req2 := bytes.NewBuffer([]byte(bs))
				req, err := http.NewRequest("POST", terminateurl, req2)
				if err != nil {
					fmt.Printf("POST error: %v", err)
				}
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				_, err = client.Do(req)
				if err != nil {
					fmt.Println("Nsmf Delete fail error message is : ", err)
				}
				return
			} else {
				return
			}
		}
	}
}
