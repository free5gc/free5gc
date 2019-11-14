package pcf_producer

import (
	"context"
	"fmt"
	"free5gc/lib/Npcf_SMPolicy"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"
	"time"

	"github.com/antihax/optional"
	"github.com/gin-gonic/gin"
)

var policyDataSubscriptionStore []models.PolicyDataSubscription
var smPolicyDataStore []models.SmPolicyData

// SmPoliciesPost -
func CreateSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.SmPolicyContextData) {
	var smPolicyContextData models.SmPolicyContextData = body
	var smPolicyDecision models.SmPolicyDecision
	var smPolicyControl models.SmPolicyControl
	var problemDetails models.ProblemDetails
	var policyDataUesUeIdSmDataGetParamOpts Nudr_DataRepository.PolicyDataUesUeIdSmDataGetParamOpts
	var policyDataSubsToNotifyPostParamOpts Nudr_DataRepository.PolicyDataSubsToNotifyPostParamOpts
	var applicationDataInfluenceDataGetParamOpts Nudr_DataRepository.ApplicationDataInfluenceDataGetParamOpts
	var policyDataSubscription models.PolicyDataSubscription
	var sid int32
	pcfUeContext := pcf_context.PCF_Self().UePool

	if (smPolicyContextData.PduSessionId != 0) && (smPolicyContextData.Dnn != "") && (smPolicyContextData.NotificationUri != "") && (smPolicyContextData.PduSessionType != "") && (smPolicyContextData.SliceInfo.Sd != "") && (smPolicyContextData.SliceInfo.Sst != 0) {
		supi := smPolicyContextData.Supi
		if pcfUeContext[supi] == nil {
			problemDetails.Title = "Not found Supi"
			problemDetails.Status = 400
			problemDetails.Cause = "USER_UNKNOWN"
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 400, problemDetails)
			return
		}
		sid = smPolicyContextData.PduSessionId
		if pcfUeContext[supi].SmPolicyControlStore != nil {
			if sid == pcfUeContext[supi].SmPolicyControlStore.Context.PduSessionId {
				sid_1 := pcfUeContext[supi].SmPolicyControlStore.Context.PduSessionId
				uri := pcf_util.PCF_BASIC_PATH + pcf_context.SmpolicyUri + fmt.Sprintf("%d", sid_1)
				respHeader := make(http.Header)
				respHeader.Set("Location", uri)
				pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 303, nil)
				return
			}
		}
		policyDataSubscription.NotificationUri = pcf_context.NotifiUri + smPolicyContextData.Supi
		policyDataSubscription.MonitoredResourceUris = []string{"A", "B"}
		policyDataSubscription.SupportedFeatures = smPolicyContextData.SuppFeat
		policyDataUesUeIdSmDataGetParamOpts.Snssai = optional.NewInterface(smPolicyContextData.SliceInfo)
		applicationDataInfluenceDataGetParamOpts.Supis = optional.NewInterface(smPolicyContextData.Supi)
		policyDataSubsToNotifyPostParamOpts.PolicyDataSubscription = optional.NewInterface(policyDataSubscription)
		ueid := smPolicyContextData.Supi

		//Query
		client := pcf_util.GetNudrClient()
		smPolicyData, _, err := client.DefaultApi.PolicyDataUesUeIdSmDataGet(context.Background(), ueid, &policyDataUesUeIdSmDataGetParamOpts)
		//trafficlnfIuData, _, err := client.DefaultApi.ApplicationDataInfluenceDataGet(context.Background(), &applicationDataInfluenceDataGetParamOpts)
		if err == nil {
			smPolicyDataStore = append(smPolicyDataStore, smPolicyData)
			//trafficlnfIuData = append(trafficlnfIuData, trafficlnfIuData)
		} else {
			var smData = models.SmPolicyData{
				SmPolicySnssaiData: map[string]models.SmPolicySnssaiData{
					"Snssai": {Snssai: smPolicyContextData.SliceInfo},
				},
			}
			smPolicyDataStore = append(smPolicyDataStore, smData)
			//trafficlnfIuData = append(trafficlnfIuData, trafficlnfIuData)
		}

		//Subscription
		policyDataSubscription, _, err := client.DefaultApi.PolicyDataSubsToNotifyPost(context.Background(), &policyDataSubsToNotifyPostParamOpts)
		if err == nil {
			policyDataSubscriptionStore = append(policyDataSubscriptionStore, policyDataSubscription)
		} else {
			logger.SMpolicylog.Warnln("Npcf_SMPolicy Subscribe fail error message is : ", err)
		}
		snssai := fmt.Sprint(smPolicyContextData.SliceInfo.Sst) + smPolicyContextData.SliceInfo.Sd
		switch smPolicyContextData.PduSessionType {
		case "IPV4":
			if smPolicyContextData.Ipv4Address == "" {
				if smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv4Index != 0 {
					ipv4index := smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv4Index
					smPolicyDecision.Ipv4Index = ipv4index
				} else {
					smPolicyDecision.Ipv4Index = pcf_context.Ipv4Index()
				}
			}
		case "IPV6":
			if smPolicyContextData.Ipv6AddressPrefix == "" {
				if smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv6Index != 0 {
					ipv6index := smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv6Index
					smPolicyDecision.Ipv6Index = ipv6index
				} else {
					smPolicyDecision.Ipv6Index = pcf_context.Ipv6Index()
				}
			}
		case "IPV4V6":
			if smPolicyContextData.Ipv4Address == "" {
				if smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv4Index != 0 {
					ipv4index := smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv4Index
					smPolicyDecision.Ipv4Index = ipv4index
				} else {
					smPolicyDecision.Ipv4Index = pcf_context.Ipv4Index()
				}
			}
			if smPolicyContextData.Ipv6AddressPrefix == "" {
				if smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv6Index != 0 {
					ipv6index := smPolicyData.SmPolicySnssaiData[snssai].SmPolicyDnnData[smPolicyContextData.Dnn].Ipv6Index
					smPolicyDecision.Ipv6Index = ipv6index
				} else {
					smPolicyDecision.Ipv6Index = pcf_context.Ipv6Index()
				}
			}
		default:
		}
		if smPolicyContextData.SuppFeat != "" {
			supp := smPolicyContextData.SuppFeat
			smPolicyDecision.SuppFeat = supp
		}
		formatTime, err := time.Parse(pcf_context.GetTimeformat(), time.Now().Format(pcf_context.GetTimeformat()))
		if err == nil {
			smPolicyDecision.RevalidationTime = &formatTime
		}
		smPolicyDecision.PccRules = make(map[string]models.PccRule)
		smPolicyDecision.PccRules["default"] = models.PccRule{
			PccRuleId: "default",
		}
		smPolicyDecision.ChargingInfo = &models.ChargingInformation{
			PrimaryChfAddress:   "string",
			SecondaryChfAddress: "string",
		}
		if smPolicyContextData.Offline || smPolicyContextData.Online {
			smPolicyDecision.Offline = smPolicyContextData.Offline
			smPolicyDecision.Online = smPolicyContextData.Online
		} else {
			smPolicyDecision.Offline = true
			smPolicyDecision.Online = true
		}

		smPolicyControl.Context = &smPolicyContextData
		smPolicyControl.Policy = &smPolicyDecision
		pcfUeContext[supi].SmPolicyControlStore = &smPolicyControl
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, smPolicyDecision)

	} else {

		problemDetails.Title = "Not found Mandatory Request Data"
		problemDetails.Status = 400
		problemDetails.Cause = "ERROR_INITIAL_PARAMETERS"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 400, problemDetails)
	}
}

// SmPoliciesSmPolicyIdDeletePost -
func DeleteSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	var problemDetails models.ProblemDetails
	var policyDataUesUeIdSmDataPatchParamOpts Nudr_DataRepository.PolicyDataUesUeIdSmDataPatchParamOpts
	var usageMonData models.UsageMonData
	logger.SMpolicylog.Traceln("URL: ", ReqURI)
	pcfUeContext := pcf_context.PCF_Self().UePool

	for key := range pcfUeContext {
		if pcfUeContext[key].SmPolicyControlStore == nil {
			continue
		}
		PduSessionIDTemp := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId)
		if ReqURI == PduSessionIDTemp {
			ueid := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.Supi)
			snssai := pcfUeContext[key].SmPolicyControlStore.Context.SliceInfo
			ipv4index := pcfUeContext[key].SmPolicyControlStore.Policy.Ipv4Index
			if ipv4index != 0 {
				pcf_context.DeleteIpv4index(ipv4index)
			}
			pcfUeContext[key].SmPolicyControlStore = nil
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, gin.H{})
			//notify policyAuthorization
			Npcf_PolicyAuthorization_Notify(fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId), "terminate")
			for i := 0; i < len(smPolicyDataStore); i++ {
				if snssai == smPolicyDataStore[i].SmPolicySnssaiData["Snssai"].Snssai {
					usageMonData.LimitId = fmt.Sprint(smPolicyDataStore[i].UmDataLimits["LimitId"].LimitId)
				}
			}
			policyDataUesUeIdSmDataPatchParamOpts.RequestBody = optional.NewInterface(usageMonData)
			//patchquery
			client := pcf_util.GetNudrClient()
			_, err := client.DefaultApi.PolicyDataUesUeIdSmDataPatch(context.Background(), ueid, &policyDataUesUeIdSmDataPatchParamOpts)
			if err != nil {
				logger.SMpolicylog.Warnln("Npcf Delete Query fail error message is : ", err)
			}
			//unsubscribe
			_, err = client.DefaultApi.PolicyDataSubsToNotifySubsIdDelete(context.Background(), "Subsid")
			if err == nil {
				return
			}
			return

		}
	}
	problemDetails.Status = 404
	problemDetails.Cause = "APPLICATION_SESSION_CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)

}

// SmPoliciesSmPolicyIdGet -
func GetSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	var problemDetails models.ProblemDetails
	logger.SMpolicylog.Traceln("URL: ", ReqURI)
	pcfUeContext := pcf_context.PCF_Self().UePool
	for key := range pcfUeContext {
		if pcfUeContext[key].SmPolicyControlStore == nil {
			continue
		}
		PduSessionIDTemp := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId)
		if ReqURI == PduSessionIDTemp {
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, pcfUeContext[key].SmPolicyControlStore)
			return
		}
	}
	problemDetails.Status = 404
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)
}

// SmPoliciesSmPolicyIdUpdatePost -
func UpdateSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.SmPolicyUpdateContextData) {
	var smPolicyUpdateContextData models.SmPolicyUpdateContextData = body
	var smPolicyDecision models.SmPolicyDecision
	var problemDetails models.ProblemDetails
	var policyDataUesUeIdSmDataGetParamOpts Nudr_DataRepository.PolicyDataUesUeIdSmDataGetParamOpts
	var formatTimeStrAdd string
	pcfUeContext := pcf_context.PCF_Self().UePool
	if len(smPolicyUpdateContextData.RepPolicyCtrlReqTriggers) > 0 {
		for key := range pcfUeContext {
			if pcfUeContext[key].SmPolicyControlStore == nil {
				continue
			}
			PduSessionIDTemp := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId)
			if ReqURI == PduSessionIDTemp {
				policyDataUesUeIdSmDataGetParamOpts.Snssai = optional.NewInterface(pcfUeContext[key].SmPolicyControlStore.Context.SliceInfo)
				ueid := fmt.Sprint(pcfUeContext[key].Supi)
				//Query
				client := pcf_util.GetNudrClient()
				smPolicyData, _, err := client.DefaultApi.PolicyDataUesUeIdSmDataGet(context.Background(), ueid, &policyDataUesUeIdSmDataGetParamOpts)
				if err == nil {
					smPolicyDataStore = append(smPolicyDataStore, smPolicyData)
				} else {
					//PolicyAuthorization Terminate Notify
					Npcf_PolicyAuthorization_Notify(fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId), "terminate")
					logger.SMpolicylog.Warnln("Nudr Query fail error message is : ", err)
				}
				//PolicyAuthorization Update Notify
				Npcf_PolicyAuthorization_Notify(fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId), "update")
				suppfeat := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.SuppFeat)

				smPolicyDecision.ChargingInfo = &models.ChargingInformation{
					PrimaryChfAddress:   "string",
					SecondaryChfAddress: "string",
				}
				if smPolicyUpdateContextData.TraceReq != nil {
					pcfUeContext[key].SmPolicyControlStore.Context.TraceReq = smPolicyUpdateContextData.TraceReq
				}
				smPolicyDecision.SuppFeat = suppfeat
				for x := 0; x < len(smPolicyUpdateContextData.RepPolicyCtrlReqTriggers); x++ {
					switch smPolicyUpdateContextData.RepPolicyCtrlReqTriggers[x] {
					case "PLMN_CH": //PLMN Change
					case "RES_MO_RE": // UE notice to SMF
					case "AC_TY_CH": //notice SMF that UE status changed
					case "UE_IP_CH": //SMF notice PCF "ipv4Address" & ipv6AddressPrefix
					case "UE_MAC_CH": //SMF notice PCF when SMF detect new ue mac
					case "AN_CH_COR":
					case "US_RE":
					case "APP_STA": //no response
					case "APP_STO": //no response
					case "AN_INFO":
					case "CM_SES_FAIL":
					case "PS_DA_OFF": //3GPP PS Data Off status changed
					case "DEF_QOS_CH":
						smPolicyDecision.SessRules = make(map[string]models.SessionRule)
						smPolicyDecision.SessRules["default"] = models.SessionRule{
							SessRuleId: "default",
						}
					case "SE_AMBR_CH":
						smPolicyDecision.PccRules = make(map[string]models.PccRule)
						smPolicyDecision.PccRules["default"] = models.PccRule{
							PccRuleId: "default",
						}
					case "QOS_NOTIF":
					case "NO_CREDIT":
					case "PRA_CH":
					case "SAREA_CH":
					case "SCNN_CH":
					case "RE_TIMEOUT":
						formatTimeStr := time.Now()
						formatTimeStr = formatTimeStr.Add(time.Second * 60)
						formatTimeStrAdd = formatTimeStr.Format(pcf_context.GetTimeformat())
						formatTime, err := time.Parse(pcf_context.GetTimeformat(), formatTimeStrAdd)
						if err == nil {
							smPolicyDecision.RevalidationTime = &formatTime
						}
					case "RES_RELEASE":
					case "SUCC_RES_ALLO": //SMF response PCF
					case "RAT_TY_CH": //SMF notice PCF
					case "REF_QOS_IND_CH": //PCF response PolicyDecision "reflectiveQosTimer"(Option)
					}
				}

				pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, smPolicyDecision)
				return

			}

		}

		problemDetails.Status = 403
		problemDetails.Cause = "ERROR_CONFLICTING_REQUEST"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
		return
	} else {
		problemDetails.Status = 400
		problemDetails.Cause = "ERROR_TRIGGER_EVENT"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 400, problemDetails)
		return
	}
}

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
