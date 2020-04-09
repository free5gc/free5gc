package pcf_producer

import (
	"context"
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"
	"strings"
	"time"

	"github.com/cydev/zero"
)

// Initial provisioning of service information (DONE)
// Gate control (DONE)
// Initial provisioning of sponsored connectivity information (DONE)
// Subscriptions to Service Data Flow QoS notification control (DONE)
// Subscription to Service Data Flow Deactivation (DONE)
// Initial provisioning of traffic routing information (DONE)
// Subscription to resources allocation outcome (DONE)
// Invocation of Multimedia Priority Services (TODO)
// Support of content versioning (TODO)
// PostAppSessions - Creates a new Individual Application Session Context resource
func PostAppSessionsContext(httpChannel chan pcf_message.HttpResponseMessage, request models.AppSessionContext) {
	logger.PolicyAuthorizationlog.Traceln("Handle Create AppSessions")
	reqData := request.AscReqData
	pcfSelf := pcf_context.PCF_Self()
	// Initial BDT policy indication(the only one which is not related to session)
	if reqData.BdtRefId != "" {
		err := handleBackgroundDataTransferPolicyIndication(pcfSelf, &request)
		if err != nil {
			sendProblemDetail(httpChannel, err.Error(), pcf_util.ERROR_REQUEST_PARAMETERS)
			return
		}
		appSessionId := fmt.Sprintf("BdtRefId-%s", reqData.BdtRefId)
		data := pcf_context.AppSessionData{
			AppSessionId:      appSessionId,
			AppSessionContext: &request,
		}
		pcfSelf.AppSessionPool[appSessionId] = &data
		locationHeader := pcf_util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSessionId)
		headers := http.Header{
			"Location": {locationHeader},
		}
		logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Create", appSessionId)
		pcf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusCreated, request)
		return
	}
	if request.AscReqData.UeIpv4 == "" && request.AscReqData.UeIpv6 == "" && request.AscReqData.UeMac == "" {
		sendProblemDetail(httpChannel, "Ue UeIpv4 and UeIpv6 and UeMac are all empty", pcf_util.ERROR_REQUEST_PARAMETERS)
		return
	}
	smPolicy, err := pcfSelf.SessionBinding(request.AscReqData)
	if err != nil {
		sendProblemDetail(httpChannel, fmt.Sprintf("Session Binding failed[%s]", err.Error()), pcf_util.PDU_SESSION_NOT_AVAILABLE)
		return
	}
	ue := smPolicy.PcfUe
	updateSMpolicy := false
	nSuppFeat := pcf_util.GetNegotiateSuppFeat(reqData.SuppFeat, pcfSelf.PcfSuppFeats[models.ServiceName_NPCF_POLICYAUTHORIZATION])
	traffRoutSupp := pcf_util.CheckSuppFeat(nSuppFeat, 1) && pcf_util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 1) // InfluenceOnTrafficRouting = 1 in 29514 &  Traffic Steering Control support = 1 in 29512
	relatedPccRuleIds := make(map[string]string)

	if reqData.MedComponents != nil {
		// Handle Pcc rules
		maxPrecedence := getMaxPrecedence(smPolicy.PolicyDecision.PccRules)
		for _, mediaComponent := range reqData.MedComponents {
			var pccRule *models.PccRule
			// TODO: use specific algorithm instead of default, details in subsclause 7.3.3 of TS 29513
			var var5qi int32 = 9
			if mediaComponent.MedType != "" {
				var5qi = pcf_util.MediaTypeTo5qiMap[mediaComponent.MedType]
			}

			if mediaComponent.MedSubComps != nil {
				for _, mediaSubComponent := range mediaComponent.MedSubComps {
					flowInfos, err := getFlowInfos(mediaSubComponent)
					if err != nil {
						sendProblemDetail(httpChannel, err.Error(), pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
						return
					}
					pccRule = pcf_util.GetPccRuleByFlowInfos(smPolicy.PolicyDecision.PccRules, flowInfos)
					if pccRule == nil {
						pccRule = pcf_util.CreatePccRule(smPolicy.PccRuleIdGenarator, maxPrecedence+1, nil, false)
						// Set QoS Data
						// TODO: use real arp
						qosData := pcf_util.CreateQosData(smPolicy.PccRuleIdGenarator, var5qi, 8)
						if var5qi <= 4 {
							// update Qos Data accroding to request BitRate
							var ul, dl bool
							qosData, ul, dl = updateQos_subComp(qosData, mediaComponent, mediaSubComponent)
							err = modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
							if err != nil {
								return
							}
						}
						// Set PackfiltId
						for i := range flowInfos {
							flowInfos[i].PackFiltId = pcf_util.GetPackFiltId(smPolicy.PackFiltIdGenarator)
							smPolicy.PackFiltMapToPccRuleId[flowInfos[i].PackFiltId] = pccRule.PccRuleId
							smPolicy.PackFiltIdGenarator++
						}
						// Set flowsInfo in Pcc Rule
						pccRule.FlowInfos = flowInfos
						// Set Traffic Control Data
						tcData := pcf_util.CreateTcData(smPolicy.PccRuleIdGenarator, mediaSubComponent.FStatus)
						pcf_util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, &tcData, &qosData, nil, nil)
						smPolicy.PccRuleIdGenarator++
						maxPrecedence++
					} else {
						// update qos
						var qosData models.QosData
						for _, qosId := range pccRule.RefQosData {
							qosData = smPolicy.PolicyDecision.QosDecs[qosId]
							if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
								var ul, dl bool
								qosData, ul, dl = updateQos_subComp(smPolicy.PolicyDecision.QosDecs[qosId], mediaComponent, mediaSubComponent)
								err = modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
								if err != nil {
									fmt.Println(err.Error())
									return
								}
								smPolicy.PolicyDecision.QosDecs[qosData.QosId] = qosData
							}
						}
					}
					// Initial provisioning of traffic routing information
					if traffRoutSupp {
						InitialProvisioningOfTrafficRoutingInformation(smPolicy, pccRule, mediaComponent.AfRoutReq, reqData.AfRoutReq)
					}
					smPolicy.PolicyDecision.PccRules[pccRule.PccRuleId] = *pccRule
					key := fmt.Sprintf("%d-%d", mediaComponent.MedCompN, mediaSubComponent.FNum)
					relatedPccRuleIds[key] = pccRule.PccRuleId
					updateSMpolicy = true
				}
				continue
			} else if mediaComponent.AfAppId != "" {
				// if mediaComponent.AfAppId has value -> find pccRule by reqData.AfAppId, otherwise create a new pcc rule
				pccRule = pcf_util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, mediaComponent.AfAppId)
				if pccRule != nil {
					pccRule.AppId = mediaComponent.AfAppId
				}
			} else if reqData.AfAppId != "" {
				pccRule = pcf_util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, reqData.AfAppId)
				if pccRule != nil {
					pccRule.AppId = reqData.AfAppId
				}
			} else {
				sendProblemDetail(httpChannel, "Media Component needs flows of subComp or afAppId", pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return
			}

			if pccRule == nil { // create new pcc rule
				pccRule = pcf_util.CreatePccRule(smPolicy.PccRuleIdGenarator, maxPrecedence+1, nil, false)
				if mediaComponent.AfAppId != "" {
					pccRule.AppId = mediaComponent.AfAppId
				} else {
					pccRule.AppId = reqData.AfAppId
				}
				// Set QoS Data
				// TODO: use real arp
				qosData := pcf_util.CreateQosData(smPolicy.PccRuleIdGenarator, var5qi, 8)
				if var5qi <= 4 {
					// update Qos Data accroding to request BitRate
					var ul, dl bool
					qosData, ul, dl = updateQos_Comp(qosData, mediaComponent)
					err = modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
					if err != nil {
						return
					}
				}

				// Set Traffic Control Data
				tcData := pcf_util.CreateTcData(smPolicy.PccRuleIdGenarator, mediaComponent.FStatus)
				pcf_util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, &tcData, &qosData, nil, nil)
				smPolicy.PccRuleIdGenarator++
				maxPrecedence++
			} else {
				// update qos
				var qosData models.QosData
				for _, qosId := range pccRule.RefQosData {
					qosData = smPolicy.PolicyDecision.QosDecs[qosId]
					if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
						var ul, dl bool
						qosData, ul, dl = updateQos_Comp(smPolicy.PolicyDecision.QosDecs[qosId], mediaComponent)
						err = modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
						if err != nil {
							return
						}
						smPolicy.PolicyDecision.QosDecs[qosData.QosId] = qosData
					}
				}
			}
			key := fmt.Sprintf("%d", mediaComponent.MedCompN)
			relatedPccRuleIds[key] = pccRule.PccRuleId
			// TODO : handle temporal or spatial validity
			// Initial provisioning of traffic routing information
			if traffRoutSupp {
				InitialProvisioningOfTrafficRoutingInformation(smPolicy, pccRule, mediaComponent.AfRoutReq, reqData.AfRoutReq)
			}
			updateSMpolicy = true
		}
	} else if reqData.AfAppId != "" {
		// Initial provisioning of traffic routing information
		if reqData.AfRoutReq != nil && traffRoutSupp {
			decision := smPolicy.PolicyDecision
			cnt := 0
			for _, rule := range smPolicy.PolicyDecision.PccRules {
				if rule.AppId == reqData.AfAppId {
					tcData := models.TrafficControlData{
						TcId:       strings.ReplaceAll(rule.PccRuleId, "PccRule", "Tc"),
						FlowStatus: models.FlowStatus_ENABLED,
					}
					tcData.RouteToLocs = append(tcData.RouteToLocs, reqData.AfRoutReq.RouteToLocs...)
					tcData.UpPathChgEvent = reqData.AfRoutReq.UpPathChgSub
					rule.RefTcData = []string{tcData.TcId}
					rule.AppReloc = reqData.AfRoutReq.AppReloc
					pcf_util.SetPccRuleRelatedData(decision, &rule, &tcData, nil, nil, nil)
					updateSMpolicy = true
					key := fmt.Sprintf("appId-%s-%d", reqData.AfAppId, cnt)
					relatedPccRuleIds[key] = rule.PccRuleId
					cnt++
				}
			}
			// Create a Pcc Rule if afappId dose not match any pcc rule
			if !updateSMpolicy {
				maxPrecedence := getMaxPrecedence(smPolicy.PolicyDecision.PccRules)
				pccRule := pcf_util.CreatePccRule(smPolicy.PccRuleIdGenarator, maxPrecedence+1, nil, false)
				pccRule.AppId = reqData.AfAppId
				qosData := models.QosData{
					QosId:                pcf_util.GetQosId(smPolicy.PccRuleIdGenarator),
					DefQosFlowIndication: true,
				}
				tcData := pcf_util.CreateTcData(smPolicy.PccRuleIdGenarator, "")
				pccRule.RefTcData = []string{tcData.TcId}
				pccRule.RefQosData = []string{qosData.QosId}
				pcf_util.SetPccRuleRelatedData(decision, pccRule, &tcData, &qosData, nil, nil)
				smPolicy.PccRuleIdGenarator++
				updateSMpolicy = true
				key := fmt.Sprintf("appId-%s", reqData.AfAppId)
				relatedPccRuleIds[key] = pccRule.PccRuleId
			}
		} else {
			sendProblemDetail(httpChannel, "Traffic routing not supported", pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return
		}
	} else {
		sendProblemDetail(httpChannel, "AF Request need AfAppId or Media Component to match Service Data Flow", pcf_util.ERROR_REQUEST_PARAMETERS)
		return
	}

	// Event Subscription
	eventSubs := make(map[models.AfEvent]models.AfNotifMethod)
	if reqData.EvSubsc != nil {
		for _, subs := range reqData.EvSubsc.Events {
			if subs.NotifMethod == "" {
				// default value "EVENT_DETECTION"
				subs.NotifMethod = models.AfNotifMethod_EVENT_DETECTION
			}
			eventSubs[subs.Event] = subs.NotifMethod
			var trig models.PolicyControlRequestTrigger
			switch subs.Event {
			case models.AfEvent_ACCESS_TYPE_CHANGE:
				trig = models.PolicyControlRequestTrigger_AC_TY_CH
			// case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
			// 	// Subscription to Service Data Flow Deactivation
			// 	trig = models.PolicyControlRequestTrigger_RES_RELEASE
			case models.AfEvent_PLMN_CHG:
				trig = models.PolicyControlRequestTrigger_PLMN_CH
			case models.AfEvent_QOS_NOTIF:
				// Subscriptions to Service Data Flow QoS notification control
				for _, pccRuleId := range relatedPccRuleIds {
					pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
					for _, qosId := range pccRule.RefQosData {
						qosData := smPolicy.PolicyDecision.QosDecs[qosId]
						qosData.Qnc = true
						smPolicy.PolicyDecision.QosDecs[qosId] = qosData
					}
				}
				trig = models.PolicyControlRequestTrigger_QOS_NOTIF
			case models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
				// Subscription to resources allocation outcome
				trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
			case models.AfEvent_USAGE_REPORT:
				trig = models.PolicyControlRequestTrigger_US_RE
			default:
				logger.PolicyAuthorizationlog.Warn("AF Event is unknown")
				continue
			}
			if !pcf_util.CheckPolicyControlReqTrig(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig) {
				smPolicy.PolicyDecision.PolicyCtrlReqTriggers = append(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig)
				updateSMpolicy = true
			}

		}
	}

	// Initial provisioning of sponsored connectivity information
	if reqData.AspId != "" && reqData.SponId != "" {
		supp := pcf_util.CheckSuppFeat(nSuppFeat, 2) && pcf_util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 12) // SponsoredConnectivity = 2 in 29514 &  SponsoredConnectivity support = 12 in 29512
		if !supp {
			sendProblemDetail(httpChannel, "Sponsored Connectivity not supported", pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return
		}
		umId := pcf_util.GetUmId(reqData.AspId, reqData.SponId)
		umData, err := extractUmData(umId, eventSubs, reqData.EvSubsc.UsgThres)
		if err != nil {
			sendProblemDetail(httpChannel, err.Error(), pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return
		}
		err = handleSponsoredConnectivityInformation(smPolicy, relatedPccRuleIds, reqData.AspId, reqData.SponId, reqData.SponStatus, umData, &updateSMpolicy)
		if err != nil {
			return
		}
	}

	// Allocate App Session Id
	appSessionId := ue.AllocUeAppSessionId(pcfSelf)
	request.AscRespData = &models.AppSessionContextRespData{
		SuppFeat: nSuppFeat,
	}
	// Associate App Session to SMPolicy
	smPolicy.AppSessions[appSessionId] = true
	data := pcf_context.AppSessionData{
		AppSessionId:      appSessionId,
		AppSessionContext: &request,
		SmPolicyData:      smPolicy,
	}
	if len(relatedPccRuleIds) > 0 {
		data.RelatedPccRuleIds = relatedPccRuleIds
		data.PccRuleIdMapToCompId = reverseStringMap(relatedPccRuleIds)
	}
	request.EvsNotif = &models.EventsNotification{}
	// Set Event Subsciption related Data
	if len(eventSubs) > 0 {
		data.Events = eventSubs
		data.EventUri = reqData.EvSubsc.NotifUri
		if _, exist := eventSubs[models.AfEvent_PLMN_CHG]; exist {
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_PLMN_CHG,
			}
			request.EvsNotif.EvNotifs = append(request.EvsNotif.EvNotifs, afNotif)
			plmnId := smPolicy.PolicyContext.ServingNetwork
			if plmnId != nil {
				request.EvsNotif.PlmnId = &models.PlmnId{
					Mcc: plmnId.Mcc,
					Mnc: plmnId.Mnc,
				}
			}
		}
		if _, exist := eventSubs[models.AfEvent_ACCESS_TYPE_CHANGE]; exist {
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_ACCESS_TYPE_CHANGE,
			}
			request.EvsNotif.EvNotifs = append(request.EvsNotif.EvNotifs, afNotif)
			request.EvsNotif.AccessType = smPolicy.PolicyContext.AccessType
			request.EvsNotif.RatType = smPolicy.PolicyContext.RatType
		}
	}
	if request.EvsNotif.EvNotifs == nil {
		request.EvsNotif = nil
	}
	pcfSelf.AppSessionPool[appSessionId] = &data
	locationHeader := pcf_util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSessionId)
	headers := http.Header{
		"Location": {locationHeader},
	}
	logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Create", appSessionId)
	pcf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusCreated, request)
	// Send Notification to SMF
	if updateSMpolicy {
		smPolicyId := fmt.Sprintf("%s-%d", ue.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      pcf_util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyId),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		SendSMPolicyUpdateNotification(ue, smPolicyId, notification)
		logger.PolicyAuthorizationlog.Tracef("Send SM Policy[%s] Update Notification", smPolicyId)
	}
}

// DeleteAppSession - Deletes an existing Individual Application Session Context
func DeleteAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string, requset *models.EventsSubscReqData) {

	logger.PolicyAuthorizationlog.Tracef("Handle Del AppSessions, AppSessionId[%s]", appSessionId)
	pcfSelf := pcf_context.PCF_Self()

	appSession := pcfSelf.AppSessionPool[appSessionId]
	if appSession == nil {
		sendProblemDetail(httpChannel, "can't find app session", pcf_util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		return
	}
	if requset != nil {
		logger.PolicyAuthorizationlog.Warnf("Delete AppSessions does not support with Event Subscription")
	}
	// Remove related pcc rule resourse
	smPolicy := appSession.SmPolicyData
	for _, pccRuleId := range appSession.RelatedPccRuleIds {
		err := smPolicy.RemovePccRule(pccRuleId)
		if err != nil {
			logger.PolicyAuthorizationlog.Warnf(err.Error())
		}
	}

	delete(smPolicy.AppSessions, appSessionId)

	logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Del", appSessionId)

	// TODO: AccUsageReport
	// if appSession.AccUsage != nil {

	// 	resp := models.AppSessionContext{
	// 		EvsNotif: &models.EventsNotification{
	// 			UsgRep: appSession.AccUsage,
	// 		},
	// 	}
	// 	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, resp)
	// } else {
	// }
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)

	delete(pcfSelf.AppSessionPool, appSessionId)

	smPolicy.ArrangeExistEventSubscription()

	// Notify SMF About Pcc Rule moval
	smPolicyId := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
	notification := models.SmPolicyNotification{
		ResourceUri:      pcf_util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyId),
		SmPolicyDecision: smPolicy.PolicyDecision,
	}
	SendSMPolicyUpdateNotification(smPolicy.PcfUe, smPolicyId, notification)
	logger.PolicyAuthorizationlog.Tracef("Send SM Policy[%s] Update Notification", smPolicyId)

}

// GetAppSession - Reads an existing Individual Application Session Context
func GetAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string) {
	logger.PolicyAuthorizationlog.Tracef("Handle Get AppSessions, AppSessionId[%s]", appSessionId)
	pcfSelf := pcf_context.PCF_Self()

	appSession := pcfSelf.AppSessionPool[appSessionId]
	if appSession == nil {
		sendProblemDetail(httpChannel, "can't find app session", pcf_util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		return
	}
	logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Get", appSessionId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *appSession.AppSessionContext)
}

// ModAppSession - Modifies an existing Individual Application Session Context
func ModAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string, request models.AppSessionContextUpdateData) {

	logger.PolicyAuthorizationlog.Tracef("Handle Modi AppSessions, AppSessionId[%s]", appSessionId)
	pcfSelf := pcf_context.PCF_Self()
	appSession := pcfSelf.AppSessionPool[appSessionId]
	if appSession == nil {
		sendProblemDetail(httpChannel, "can't find app session", pcf_util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		return
	}
	appContext := appSession.AppSessionContext
	if request.BdtRefId != "" {
		appContext.AscReqData.BdtRefId = request.BdtRefId
		err := handleBackgroundDataTransferPolicyIndication(pcfSelf, appContext)
		if err != nil {
			sendProblemDetail(httpChannel, err.Error(), pcf_util.ERROR_REQUEST_PARAMETERS)
			return
		}
		logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Updated", appSessionId)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *appContext)
		return

	}
	smPolicy := appSession.SmPolicyData
	if smPolicy == nil {
		sendProblemDetail(httpChannel, "Can't find related PDU Session", pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
		return
	}
	traffRoutSupp := pcf_util.CheckSuppFeat(appContext.AscRespData.SuppFeat, 1) && pcf_util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 1) // InfluenceOnTrafficRouting = 1 in 29514 &  Traffic Steering Control support = 1 in 29512
	relatedPccRuleIds := make(map[string]string)
	// Event Subscription
	eventSubs := make(map[models.AfEvent]models.AfNotifMethod)
	updateSMpolicy := false

	if request.MedComponents != nil {
		maxPrecedence := getMaxPrecedence(smPolicy.PolicyDecision.PccRules)
		for compN, mediaComponent := range request.MedComponents {
			removeMediaComp(appSession, compN)
			if zero.IsZero(mediaComponent) {
				// remove MediaComp(media Comp is null)
				continue
			}
			// modify MediaComp(remove and reinstall again)
			var pccRule *models.PccRule
			// TODO: use specific algorithm instead of default, details in subsclause 7.3.3 of TS 29513
			var var5qi int32 = 9
			if mediaComponent.MedType != "" {
				var5qi = pcf_util.MediaTypeTo5qiMap[mediaComponent.MedType]
			}
			qosMediaComp := models.MediaComponent{
				MarBwDl: mediaComponent.MarBwDl,
				MarBwUl: mediaComponent.MarBwUl,
				MirBwDl: mediaComponent.MirBwDl,
				MirBwUl: mediaComponent.MirBwUl,
			}
			if mediaComponent.MedSubComps != nil {
				for _, mediaSubComponentRm := range mediaComponent.MedSubComps {
					mediaSubComponent := models.MediaSubComponent(mediaSubComponentRm)
					flowInfos, err := getFlowInfos(mediaSubComponent)
					if err != nil {
						sendProblemDetail(httpChannel, err.Error(), pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
						return
					}
					pccRule = pcf_util.GetPccRuleByFlowInfos(smPolicy.PolicyDecision.PccRules, flowInfos)
					if pccRule == nil {
						pccRule = pcf_util.CreatePccRule(smPolicy.PccRuleIdGenarator, maxPrecedence+1, nil, false)
						// Set QoS Data
						// TODO: use real arp
						qosData := pcf_util.CreateQosData(smPolicy.PccRuleIdGenarator, var5qi, 8)
						if var5qi <= 4 {
							// update Qos Data accroding to request BitRate
							var ul, dl bool

							qosData, ul, dl = updateQos_subComp(qosData, qosMediaComp, mediaSubComponent)
							err = modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
							if err != nil {
								return
							}
						}
						// Set PackfiltId
						for i := range flowInfos {
							flowInfos[i].PackFiltId = pcf_util.GetPackFiltId(smPolicy.PackFiltIdGenarator)
							smPolicy.PackFiltMapToPccRuleId[flowInfos[i].PackFiltId] = pccRule.PccRuleId
							smPolicy.PackFiltIdGenarator++
						}
						// Set flowsInfo in Pcc Rule
						pccRule.FlowInfos = flowInfos
						// Set Traffic Control Data
						tcData := pcf_util.CreateTcData(smPolicy.PccRuleIdGenarator, mediaSubComponent.FStatus)
						pcf_util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, &tcData, &qosData, nil, nil)
						smPolicy.PccRuleIdGenarator++
						maxPrecedence++
					} else {
						// update qos
						var qosData models.QosData
						for _, qosId := range pccRule.RefQosData {
							qosData = smPolicy.PolicyDecision.QosDecs[qosId]
							if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
								var ul, dl bool
								qosData, ul, dl = updateQos_subComp(smPolicy.PolicyDecision.QosDecs[qosId], qosMediaComp, mediaSubComponent)
								err = modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
								if err != nil {
									fmt.Println(err.Error())
									return
								}
								smPolicy.PolicyDecision.QosDecs[qosData.QosId] = qosData
							}
						}
					}
					// Modify provisioning of traffic routing information
					if traffRoutSupp {
						ModifyProvisioningOfTrafficRoutingInformation(smPolicy, pccRule, mediaComponent.AfRoutReq, request.AfRoutReq)
					}
					smPolicy.PolicyDecision.PccRules[pccRule.PccRuleId] = *pccRule
					key := fmt.Sprintf("%d-%d", mediaComponent.MedCompN, mediaSubComponent.FNum)
					relatedPccRuleIds[key] = pccRule.PccRuleId
					updateSMpolicy = true
				}
				continue
			} else if mediaComponent.AfAppId != "" {
				// if mediaComponent.AfAppId has value -> find pccRule by reqData.AfAppId, otherwise create a new pcc rule
				pccRule = pcf_util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, mediaComponent.AfAppId)
				if pccRule != nil {
					pccRule.AppId = mediaComponent.AfAppId
				}
			} else if request.AfAppId != "" {
				pccRule = pcf_util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, request.AfAppId)
				if pccRule != nil {
					pccRule.AppId = request.AfAppId
				}
			} else {
				sendProblemDetail(httpChannel, "Media Component needs flows of subComp or afAppId", pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return
			}

			if pccRule == nil { // create new pcc rule
				pccRule = pcf_util.CreatePccRule(smPolicy.PccRuleIdGenarator, maxPrecedence+1, nil, false)
				if mediaComponent.AfAppId != "" {
					pccRule.AppId = mediaComponent.AfAppId
				} else {
					pccRule.AppId = request.AfAppId
				}
				// Set QoS Data
				// TODO: use real arp
				qosData := pcf_util.CreateQosData(smPolicy.PccRuleIdGenarator, var5qi, 8)
				if var5qi <= 4 {
					// update Qos Data accroding to request BitRate
					var ul, dl bool
					qosData, ul, dl = updateQos_Comp(qosData, qosMediaComp)
					err := modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
					if err != nil {
						return
					}
				}

				// Set Traffic Control Data
				tcData := pcf_util.CreateTcData(smPolicy.PccRuleIdGenarator, mediaComponent.FStatus)
				pcf_util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, &tcData, &qosData, nil, nil)
				smPolicy.PccRuleIdGenarator++
				maxPrecedence++
			} else {
				// update qos
				var qosData models.QosData
				for _, qosId := range pccRule.RefQosData {
					qosData = smPolicy.PolicyDecision.QosDecs[qosId]
					if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
						var ul, dl bool
						qosData, ul, dl = updateQos_Comp(smPolicy.PolicyDecision.QosDecs[qosId], qosMediaComp)
						err := modifyRemainBitRate(httpChannel, smPolicy, &qosData, ul, dl)
						if err != nil {
							return
						}
						smPolicy.PolicyDecision.QosDecs[qosData.QosId] = qosData
					}
				}
			}
			key := fmt.Sprintf("%d", mediaComponent.MedCompN)
			relatedPccRuleIds[key] = pccRule.PccRuleId
			// TODO : handle temporal or spatial validity
			// Modify provisioning of traffic routing information
			if traffRoutSupp {
				ModifyProvisioningOfTrafficRoutingInformation(smPolicy, pccRule, mediaComponent.AfRoutReq, request.AfRoutReq)
			}
			updateSMpolicy = true

		}
	}

	// Merge Original PccRuleId and new
	for key, pccRuleId := range appSession.RelatedPccRuleIds {
		relatedPccRuleIds[key] = pccRuleId
	}

	if request.EvSubsc != nil {
		for _, subs := range request.EvSubsc.Events {
			if subs.NotifMethod == "" {
				// default value "EVENT_DETECTION"
				subs.NotifMethod = models.AfNotifMethod_EVENT_DETECTION
			}
			eventSubs[subs.Event] = subs.NotifMethod
			var trig models.PolicyControlRequestTrigger
			switch subs.Event {
			case models.AfEvent_ACCESS_TYPE_CHANGE:
				trig = models.PolicyControlRequestTrigger_AC_TY_CH
			// case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
			// 	// Subscription to Service Data Flow Deactivation
			// 	trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
			case models.AfEvent_PLMN_CHG:
				trig = models.PolicyControlRequestTrigger_PLMN_CH
			case models.AfEvent_QOS_NOTIF:
				// Subscriptions to Service Data Flow QoS notification control
				for _, pccRuleId := range relatedPccRuleIds {
					pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
					for _, qosId := range pccRule.RefQosData {
						qosData := smPolicy.PolicyDecision.QosDecs[qosId]
						qosData.Qnc = true
						smPolicy.PolicyDecision.QosDecs[qosId] = qosData
					}
				}
				trig = models.PolicyControlRequestTrigger_QOS_NOTIF
			case models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
				// Subscription to resources allocation outcome
				trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
			case models.AfEvent_USAGE_REPORT:
				trig = models.PolicyControlRequestTrigger_US_RE
			default:
				logger.PolicyAuthorizationlog.Warn("AF Event is unknown")
				continue
			}
			if !pcf_util.CheckPolicyControlReqTrig(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig) {
				smPolicy.PolicyDecision.PolicyCtrlReqTriggers = append(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig)
				updateSMpolicy = true
			}

		}
		// update Context
		if appContext.AscReqData.EvSubsc == nil {
			appContext.AscReqData.EvSubsc = new(models.EventsSubscReqData)
		}
		appContext.AscReqData.EvSubsc.Events = request.EvSubsc.Events
		if request.EvSubsc.NotifUri != "" {
			appContext.AscReqData.EvSubsc.NotifUri = request.EvSubsc.NotifUri
			appSession.EventUri = request.EvSubsc.NotifUri
		}
		if request.EvSubsc.UsgThres != nil {
			appContext.AscReqData.EvSubsc.UsgThres = threshRmToThresh(request.EvSubsc.UsgThres)
		}

	} else {
		// remove eventSubs
		appSession.Events = nil
		appSession.EventUri = ""
		appContext.AscReqData.EvSubsc = nil
	}

	// Moification provisioning of sponsored connectivity information
	if request.AspId != "" && request.SponId != "" {
		umId := pcf_util.GetUmId(request.AspId, request.SponId)
		umData, err := extractUmData(umId, eventSubs, threshRmToThresh(request.EvSubsc.UsgThres))
		if err != nil {
			sendProblemDetail(httpChannel, err.Error(), pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return
		}
		err = handleSponsoredConnectivityInformation(smPolicy, relatedPccRuleIds, request.AspId, request.SponId, request.SponStatus, umData, &updateSMpolicy)
		if err != nil {
			return
		}
	}

	if len(relatedPccRuleIds) > 0 {
		appSession.RelatedPccRuleIds = relatedPccRuleIds
		appSession.PccRuleIdMapToCompId = reverseStringMap(relatedPccRuleIds)

	}
	appContext.EvsNotif = &models.EventsNotification{}
	// Set Event Subsciption related Data
	if len(eventSubs) > 0 {
		appSession.Events = eventSubs
		if _, exist := eventSubs[models.AfEvent_PLMN_CHG]; exist {
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_PLMN_CHG,
			}
			appContext.EvsNotif.EvNotifs = append(appContext.EvsNotif.EvNotifs, afNotif)
			plmnId := smPolicy.PolicyContext.ServingNetwork
			if plmnId != nil {
				appContext.EvsNotif.PlmnId = &models.PlmnId{
					Mcc: plmnId.Mcc,
					Mnc: plmnId.Mnc,
				}
			}
		}
		if _, exist := eventSubs[models.AfEvent_ACCESS_TYPE_CHANGE]; exist {
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_ACCESS_TYPE_CHANGE,
			}
			appContext.EvsNotif.EvNotifs = append(appContext.EvsNotif.EvNotifs, afNotif)
			appContext.EvsNotif.AccessType = smPolicy.PolicyContext.AccessType
			appContext.EvsNotif.RatType = smPolicy.PolicyContext.RatType
		}
	}
	if appContext.EvsNotif.EvNotifs == nil {
		appContext.EvsNotif = nil
	}

	// TODO: MPS Sevice
	logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Updated", appSessionId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *appContext)

	smPolicy.ArrangeExistEventSubscription()

	// Send Notification to SMF
	if updateSMpolicy {
		smPolicyId := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      pcf_util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyId),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		SendSMPolicyUpdateNotification(smPolicy.PcfUe, smPolicyId, notification)
		logger.PolicyAuthorizationlog.Tracef("Send SM Policy[%s] Update Notification", smPolicyId)
	}

}

// DeleteEventsSubsc - deletes the Events Subscription subresource
func DeleteEventsSubscContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string) {

	logger.PolicyAuthorizationlog.Tracef("Handle Del AppSessions Events Subsc, AppSessionId[%s]", appSessionId)
	pcfSelf := pcf_context.PCF_Self()

	appSession := pcfSelf.AppSessionPool[appSessionId]
	if appSession == nil {
		sendProblemDetail(httpChannel, "can't find app session", pcf_util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		return
	}
	appSession.Events = nil
	appSession.EventUri = ""
	appSession.AppSessionContext.EvsNotif = nil
	appSession.AppSessionContext.AscReqData.EvSubsc = nil

	changed := appSession.SmPolicyData.ArrangeExistEventSubscription()

	logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Del Events Subsc success", appSessionId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)

	smPolicy := appSession.SmPolicyData
	// Send Notification to SMF
	if changed {
		smPolicyId := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      pcf_util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyId),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		SendSMPolicyUpdateNotification(smPolicy.PcfUe, smPolicyId, notification)
		logger.PolicyAuthorizationlog.Tracef("Send SM Policy[%s] Update Notification", smPolicyId)
	}
}

// UpdateEventsSubsc - creates or modifies an Events Subscription subresource
func UpdateEventsSubscContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string, request models.EventsSubscReqData) {

	logger.PolicyAuthorizationlog.Tracef("Handle Put AppSessions Events Subsc, AppSessionId[%s]", appSessionId)
	pcfSelf := pcf_context.PCF_Self()

	appSession := pcfSelf.AppSessionPool[appSessionId]
	if appSession == nil {
		sendProblemDetail(httpChannel, "can't find app session", pcf_util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		return
	}
	smPolicy := appSession.SmPolicyData
	eventSubs := make(map[models.AfEvent]models.AfNotifMethod)

	updataSmPolicy := false
	created := false
	if appSession.Events == nil {
		created = true
	}

	for _, subs := range request.Events {
		if subs.NotifMethod == "" {
			// default value "EVENT_DETECTION"
			subs.NotifMethod = models.AfNotifMethod_EVENT_DETECTION
		}
		eventSubs[subs.Event] = subs.NotifMethod
		var trig models.PolicyControlRequestTrigger
		switch subs.Event {
		case models.AfEvent_ACCESS_TYPE_CHANGE:
			trig = models.PolicyControlRequestTrigger_AC_TY_CH
		// case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
		// 	// Subscription to Service Data Flow Deactivation
		// 	trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
		case models.AfEvent_PLMN_CHG:
			trig = models.PolicyControlRequestTrigger_PLMN_CH
		case models.AfEvent_QOS_NOTIF:
			// Subscriptions to Service Data Flow QoS notification control
			for _, pccRuleId := range appSession.RelatedPccRuleIds {
				pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
				for _, qosId := range pccRule.RefQosData {
					qosData := smPolicy.PolicyDecision.QosDecs[qosId]
					qosData.Qnc = true
					smPolicy.PolicyDecision.QosDecs[qosId] = qosData
				}
			}
			trig = models.PolicyControlRequestTrigger_QOS_NOTIF
		case models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
			// Subscription to resources allocation outcome
			trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
		case models.AfEvent_USAGE_REPORT:
			trig = models.PolicyControlRequestTrigger_US_RE
		default:
			logger.PolicyAuthorizationlog.Warn("AF Event is unknown")
			continue
		}
		if !pcf_util.CheckPolicyControlReqTrig(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig) {
			smPolicy.PolicyDecision.PolicyCtrlReqTriggers = append(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig)
			updataSmPolicy = true
		}

	}
	appContext := appSession.AppSessionContext
	// update Context
	if appContext.AscReqData.EvSubsc == nil {
		appContext.AscReqData.EvSubsc = new(models.EventsSubscReqData)
	}
	appContext.AscReqData.EvSubsc.Events = request.Events
	appContext.AscReqData.EvSubsc.UsgThres = request.UsgThres
	appContext.AscReqData.EvSubsc.NotifUri = request.NotifUri
	appContext.EvsNotif = nil
	// update app Session
	appSession.EventUri = request.NotifUri
	appSession.Events = eventSubs

	resp := models.UpdateEventsSubscResponse{
		EvSubsc: request,
	}
	appContext.EvsNotif = &models.EventsNotification{
		EvSubsUri: request.NotifUri,
	}
	// Set Event Subsciption related Data
	if len(eventSubs) > 0 {
		if _, exist := eventSubs[models.AfEvent_PLMN_CHG]; exist {
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_PLMN_CHG,
			}
			appContext.EvsNotif.EvNotifs = append(appContext.EvsNotif.EvNotifs, afNotif)
			plmnId := smPolicy.PolicyContext.ServingNetwork
			if plmnId != nil {
				appContext.EvsNotif.PlmnId = &models.PlmnId{
					Mcc: plmnId.Mcc,
					Mnc: plmnId.Mnc,
				}
			}
		}
		if _, exist := eventSubs[models.AfEvent_ACCESS_TYPE_CHANGE]; exist {
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_ACCESS_TYPE_CHANGE,
			}
			appContext.EvsNotif.EvNotifs = append(appContext.EvsNotif.EvNotifs, afNotif)
			appContext.EvsNotif.AccessType = smPolicy.PolicyContext.AccessType
			appContext.EvsNotif.RatType = smPolicy.PolicyContext.RatType
		}
	}
	if appContext.EvsNotif.EvNotifs == nil {
		appContext.EvsNotif = nil
	}

	resp.EvsNotif = appContext.EvsNotif

	if created {
		locationHeader := fmt.Sprintf("%s/events-subscription", pcf_util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSessionId))
		headers := http.Header{
			"Location": {locationHeader},
		}
		logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Create Subscription", appSessionId)
		pcf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusCreated, resp)
	} else if resp.EvsNotif != nil {
		logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Modify Subscription", appSessionId)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, resp)
	} else {
		logger.PolicyAuthorizationlog.Tracef("App Session Id[%s] Modify Subscription", appSessionId)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, resp)
	}

	changed := appSession.SmPolicyData.ArrangeExistEventSubscription()

	// Send Notification to SMF
	if updataSmPolicy || changed {
		smPolicyId := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      pcf_util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyId),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		SendSMPolicyUpdateNotification(smPolicy.PcfUe, smPolicyId, notification)
		logger.PolicyAuthorizationlog.Tracef("Send SM Policy[%s] Update Notification", smPolicyId)
	}
}

func SendAppSessionEventNotification(appSession *pcf_context.AppSessionData, request models.EventsNotification) {
	logger.PolicyAuthorizationlog.Tracef("Send App Session Event Notification")
	if appSession == nil {
		logger.PolicyAuthorizationlog.Warnln("Send App Session Event Notification Error[appSession is nil]")
		return
	}
	uri := appSession.EventUri
	if uri != "" {
		request.EvSubsUri = fmt.Sprintf("%s/events-subscription", pcf_util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSession.AppSessionId))
		client := pcf_util.GetNpcfPolicyAuthorizationCallbackClient()
		httpResponse, err := client.PolicyAuthorizationEventNotificationApi.PolicyAuthorizationEventNotification(context.Background(), uri, request)
		if err != nil {
			if httpResponse != nil {
				logger.PolicyAuthorizationlog.Warnf("Send App Session Event Notification Error[%s]", httpResponse.Status)
			} else {
				logger.PolicyAuthorizationlog.Warnf("Send App Session Event Notification Failed[%s]", err.Error())
			}
			return
		} else if httpResponse == nil {
			logger.PolicyAuthorizationlog.Warnln("Send App Session Event Notification Failed[HTTP Response is nil]")
			return
		}
		if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
			logger.PolicyAuthorizationlog.Warnf("Send App Session Event Notification Failed")
		} else {
			logger.PolicyAuthorizationlog.Tracef("Send App Session Event Notification Success")
		}
	}
}

func SendAppSessionTermination(appSession *pcf_context.AppSessionData, request models.TerminationInfo) {
	logger.PolicyAuthorizationlog.Tracef("Send App Session Termination")
	if appSession == nil {
		logger.PolicyAuthorizationlog.Warnln("Send App Session Termination Error[appSession is nil]")
		return
	}
	uri := appSession.AppSessionContext.AscReqData.NotifUri
	if uri != "" {
		request.ResUri = pcf_util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSession.AppSessionId)
		client := pcf_util.GetNpcfPolicyAuthorizationCallbackClient()
		httpResponse, err := client.PolicyAuthorizationTerminateRequestApi.PolicyAuthorizationTerminateRequest(context.Background(), uri, request)
		if err != nil {
			if httpResponse != nil {
				logger.PolicyAuthorizationlog.Warnf("Send App Session Termination Error[%s]", httpResponse.Status)
			} else {
				logger.PolicyAuthorizationlog.Warnf("Send App Session Termination Failed[%s]", err.Error())
			}
			return
		} else if httpResponse == nil {
			logger.PolicyAuthorizationlog.Warnln("Send App Session Termination Failed[HTTP Response is nil]")
			return
		}
		if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
			logger.PolicyAuthorizationlog.Warnf("Send App Session Termination Failed")
		} else {
			logger.PolicyAuthorizationlog.Tracef("Send App Session Termination Success")
		}
	}
}

// Handle Create/ Modify  Background Data Transfer Policy Indication
func handleBackgroundDataTransferPolicyIndication(pcfSelf *pcf_context.PCFContext, appContext *models.AppSessionContext) (err error) {
	req := appContext.AscReqData
	respData := models.AppSessionContextRespData{
		ServAuthInfo: models.ServAuthInfo_NOT_KNOWN,
		SuppFeat:     pcf_util.GetNegotiateSuppFeat(req.SuppFeat, pcfSelf.PcfSuppFeats[models.ServiceName_NPCF_POLICYAUTHORIZATION]),
	}
	client := pcf_util.GetNudrClient(getDefaultUdrUri(pcfSelf))
	bdtData, resp, err1 := client.DefaultApi.PolicyDataBdtDataBdtReferenceIdGet(context.Background(), req.BdtRefId)
	if err1 != nil {
		return fmt.Errorf("UDR Get BdtDate error[%s]", err1.Error())
	} else if resp == nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("UDR Get BdtDate error")
	} else {
		startTime, err1 := time.Parse(pcf_util.TimeFormat, bdtData.TransPolicy.RecTimeInt.StartTime)
		if err1 != nil {
			return err1
		}
		stopTime, err1 := time.Parse(pcf_util.TimeFormat, bdtData.TransPolicy.RecTimeInt.StopTime)
		if err1 != nil {
			return err1
		}
		if startTime.After(time.Now()) {
			respData.ServAuthInfo = models.ServAuthInfo_NOT_YET_OCURRED
		} else if stopTime.Before(time.Now()) {
			respData.ServAuthInfo = models.ServAuthInfo_EXPIRED
		}
	}
	appContext.AscRespData = &respData
	return nil
}

// provisioning of sponsored connectivity information
func handleSponsoredConnectivityInformation(smPolicy *pcf_context.UeSmPolicyData, relatedPccRuleIds map[string]string, aspId, sponId string, sponStatus models.SponsoringStatus, umData *models.UsageMonitoringData, updateSMpolicy *bool) (err error) {
	if sponStatus == models.SponsoringStatus_DISABLED {
		logger.PolicyAuthorizationlog.Debugf("Sponsored Connectivity is disabled by AF")
		umId := pcf_util.GetUmId(aspId, sponId)
		for _, pccRuleId := range relatedPccRuleIds {
			pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
			for _, chgId := range pccRule.RefChgData {
				// disables sponsoring a service
				chgData := smPolicy.PolicyDecision.ChgDecs[chgId]
				if chgData.AppSvcProvId == aspId && chgData.SponsorId == sponId {
					chgData.SponsorId = ""
					chgData.AppSvcProvId = ""
					chgData.ReportingLevel = models.ReportingLevel_SER_ID_LEVEL
					smPolicy.PolicyDecision.ChgDecs[chgId] = chgData
					*updateSMpolicy = true
				}
			}
			if pccRule.RefUmData != nil {
				pccRule.RefUmData = nil
				smPolicy.PolicyDecision.PccRules[pccRuleId] = pccRule
			}
			// disable the usage monitoring
			// TODO: As a result, PCF gets the accumulated usage of the sponsored data connectivity
			delete(smPolicy.PolicyDecision.UmDecs, umId)
		}
	} else {

		if umData != nil {
			supp := pcf_util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 5) // UMC support = 5 in 29512
			if !supp {
				err = fmt.Errorf("Usage Monitor Control is not supported in SMF")
				return
			}
		}
		chgIdUsed := false
		chgId := pcf_util.GetChgId(smPolicy.ChargingIdGenarator)
		for _, pccRuleId := range relatedPccRuleIds {
			pccRule := smPolicy.PolicyDecision.PccRules[pccRuleId]
			chgData := models.ChargingData{
				ChgId: chgId,
			}
			if pccRule.RefChgData != nil {
				chgId := pccRule.RefChgData[0]
				chgData = smPolicy.PolicyDecision.ChgDecs[chgId]
			} else {
				chgIdUsed = true
			}
			// TODO: PCF, based on operator policies, shall check whether it is required to validate the sponsored connectivity data.
			// If it is required, it shall perform the authorizations based on sponsored data connectivity profiles.
			// If the authorization fails, the PCF shall send HTTP "403 Forbidden" with the "cause" attribute set to "UNAUTHORIZED_SPONSORED_DATA_CONNECTIVITY"
			pccRule.RefChgData = []string{chgData.ChgId}
			chgData.ReportingLevel = models.ReportingLevel_SPON_CON_LEVEL
			chgData.SponsorId = sponId
			chgData.AppSvcProvId = aspId
			if umData != nil {
				pccRule.RefUmData = []string{umData.UmId}
			}
			pcf_util.SetPccRuleRelatedData(smPolicy.PolicyDecision, &pccRule, nil, nil, &chgData, umData)
			*updateSMpolicy = true
		}
		if chgIdUsed {
			smPolicy.ChargingIdGenarator++
		}
		// TODO: handling UE is roaming in VPLMN case
	}
	return
}

func getMaxPrecedence(pccRules map[string]models.PccRule) (maxVaule int32) {
	maxVaule = 0
	for _, rule := range pccRules {
		if rule.Precedence > maxVaule {
			maxVaule = rule.Precedence
		}
	}
	return
}

/*
func getFlowInfos(comp models.MediaComponent) (flows []models.FlowInformation, err error) {
	for _, subComp := range comp.MedSubComps {
		if subComp.EthfDescs != nil {
			return nil, fmt.Errorf("Flow Description with Mac Address does not support")
		}
		fStatus := subComp.FStatus
		if subComp.FlowUsage == models.FlowUsage_RTCP {
			fStatus = models.FlowStatus_ENABLED
		} else if fStatus == "" {
			fStatus = comp.FStatus
		}
		if fStatus == models.FlowStatus_REMOVED {
			continue
		}
		// gate control
		statusUsage := map[models.FlowDirection]bool{
			models.FlowDirection_UPLINK:   true,
			models.FlowDirection_DOWNLINK: true,
		}
		switch fStatus {
		case models.FlowStatus_ENABLED_UPLINK:
			statusUsage[models.FlowDirection_DOWNLINK] = false
		case models.FlowStatus_ENABLED_DOWNLINK:
			statusUsage[models.FlowDirection_UPLINK] = false
		case models.FlowStatus_DISABLED:
			statusUsage[models.FlowDirection_DOWNLINK] = false
			statusUsage[models.FlowDirection_UPLINK] = false
		}
		for _, desc := range subComp.FDescs {
			flowDesc, flowDir, err := flowDescriptionFromN5toN7(desc)
			if err != nil {
				return nil, err
			}
			flowInfo := models.FlowInformation{
				FlowDescription:   flowDesc,
				FlowDirection:     models.FlowDirectionRm(flowDir),
				PacketFilterUsage: statusUsage[flowDir],
				TosTrafficClass:   subComp.TosTrCl,
			}
			flows = append(flows, flowInfo)
		}
	}
	return
}
*/

func getFlowInfos(subComp models.MediaSubComponent) (flows []models.FlowInformation, err error) {
	if subComp.EthfDescs != nil {
		return nil, fmt.Errorf("Flow Description with Mac Address does not support")
	}
	fStatus := subComp.FStatus
	if subComp.FlowUsage == models.FlowUsage_RTCP {
		fStatus = models.FlowStatus_ENABLED
	}
	if fStatus == models.FlowStatus_REMOVED {
		return nil, nil
	}
	// gate control
	statusUsage := map[models.FlowDirection]bool{
		models.FlowDirection_UPLINK:   true,
		models.FlowDirection_DOWNLINK: true,
	}
	switch fStatus {
	case models.FlowStatus_ENABLED_UPLINK:
		statusUsage[models.FlowDirection_DOWNLINK] = false
	case models.FlowStatus_ENABLED_DOWNLINK:
		statusUsage[models.FlowDirection_UPLINK] = false
	case models.FlowStatus_DISABLED:
		statusUsage[models.FlowDirection_DOWNLINK] = false
		statusUsage[models.FlowDirection_UPLINK] = false
	}
	for _, desc := range subComp.FDescs {
		flowDesc, flowDir, err := flowDescriptionFromN5toN7(desc)
		if err != nil {
			return nil, err
		}
		flowInfo := models.FlowInformation{
			FlowDescription:   flowDesc,
			FlowDirection:     models.FlowDirectionRm(flowDir),
			PacketFilterUsage: statusUsage[flowDir],
			TosTrafficClass:   subComp.TosTrCl,
		}
		flows = append(flows, flowInfo)
	}
	return
}

func flowDescriptionFromN5toN7(n5Flow string) (n7Flow string, direction models.FlowDirection, err error) {
	if strings.HasPrefix(n5Flow, "permit out") {
		n7Flow = n5Flow
		direction = models.FlowDirection_DOWNLINK
	} else if strings.HasPrefix(n5Flow, "permit in") {
		n7Flow = strings.Replace(n5Flow, "permit in", "permit out", -1)
		direction = models.FlowDirection_UPLINK
	} else if strings.HasPrefix(n5Flow, "permit inout") {
		n7Flow = strings.Replace(n5Flow, "permit inout", "permit out", -1)
		direction = models.FlowDirection_BIDIRECTIONAL
	} else {
		err = fmt.Errorf("Invaild flow Description[%s]", n5Flow)
	}
	return
}
func updateQos_Comp(qosData models.QosData, comp models.MediaComponent) (updatedQosData models.QosData, ulExist, dlExist bool) {
	updatedQosData = qosData
	if comp.FStatus == models.FlowStatus_REMOVED {
		updatedQosData.MaxbrDl = ""
		updatedQosData.MaxbrUl = ""
		return
	}
	maxBwUl := 0.0
	maxBwDl := 0.0
	minBwUl := 0.0
	minBwDl := 0.0
	for _, subsComp := range comp.MedSubComps {
		for _, flow := range subsComp.FDescs {
			_, dir, _ := flowDescriptionFromN5toN7(flow)
			both := false
			if dir == models.FlowDirection_BIDIRECTIONAL {
				both = true
			}
			if subsComp.FlowUsage != models.FlowUsage_RTCP {
				// not RTCP
				if both || dir == models.FlowDirection_UPLINK {
					ulExist = true
					if comp.MarBwUl != "" {
						bwUl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
						maxBwUl += bwUl
					}
					if comp.MirBwUl != "" {
						bwUl, _ := pcf_context.ConvertBitRateToKbps(comp.MirBwUl)
						minBwUl += bwUl
					}
				}
				if both || dir == models.FlowDirection_DOWNLINK {
					dlExist = true
					if comp.MarBwDl != "" {
						bwDl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
						maxBwDl += bwDl
					}
					if comp.MirBwDl != "" {
						bwDl, _ := pcf_context.ConvertBitRateToKbps(comp.MirBwDl)
						minBwDl += bwDl
					}
				}
			} else {
				if both || dir == models.FlowDirection_UPLINK {
					ulExist = true
					if subsComp.MarBwUl != "" {
						bwUl, _ := pcf_context.ConvertBitRateToKbps(subsComp.MarBwUl)
						maxBwUl += bwUl
					} else if comp.MarBwUl != "" {
						bwUl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
						maxBwUl += (0.05 * bwUl)
					}
				}
				if both || dir == models.FlowDirection_DOWNLINK {
					dlExist = true
					if subsComp.MarBwDl != "" {
						bwDl, _ := pcf_context.ConvertBitRateToKbps(subsComp.MarBwDl)
						maxBwDl += bwDl
					} else if comp.MarBwDl != "" {
						bwDl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
						maxBwDl += (0.05 * bwDl)
					}
				}
			}
		}
	}
	// update Downlink MBR
	if maxBwDl == 0.0 {
		updatedQosData.MaxbrDl = comp.MarBwDl
	} else {
		updatedQosData.MaxbrDl = pcf_context.ConvertBitRateToString(maxBwDl)
	}
	// update Uplink MBR
	if maxBwUl == 0.0 {
		updatedQosData.MaxbrUl = comp.MarBwUl
	} else {
		updatedQosData.MaxbrUl = pcf_context.ConvertBitRateToString(maxBwUl)
	}
	// if gbr == 0 then assign gbr = mbr

	// update Downlink GBR
	if minBwDl != 0.0 {
		updatedQosData.GbrDl = pcf_context.ConvertBitRateToString(minBwDl)
	}
	// update Uplink GBR
	if minBwUl != 0.0 {
		updatedQosData.GbrUl = pcf_context.ConvertBitRateToString(minBwUl)
	}
	return
}

func updateQos_subComp(qosData models.QosData, comp models.MediaComponent, subsComp models.MediaSubComponent) (updatedQosData models.QosData, ulExist, dlExist bool) {
	updatedQosData = qosData
	if comp.FStatus == models.FlowStatus_REMOVED {
		updatedQosData.MaxbrDl = ""
		updatedQosData.MaxbrUl = ""
		return
	}
	maxBwUl := 0.0
	maxBwDl := 0.0
	minBwUl := 0.0
	minBwDl := 0.0
	for _, flow := range subsComp.FDescs {
		_, dir, _ := flowDescriptionFromN5toN7(flow)
		both := false
		if dir == models.FlowDirection_BIDIRECTIONAL {
			both = true
		}
		if subsComp.FlowUsage != models.FlowUsage_RTCP {
			// not RTCP
			if both || dir == models.FlowDirection_UPLINK {
				ulExist = true
				if comp.MarBwUl != "" {
					bwUl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
					maxBwUl += bwUl
				}
				if comp.MirBwUl != "" {
					bwUl, _ := pcf_context.ConvertBitRateToKbps(comp.MirBwUl)
					minBwUl += bwUl
				}
			}
			if both || dir == models.FlowDirection_DOWNLINK {
				dlExist = true
				if comp.MarBwDl != "" {
					bwDl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
					maxBwDl += bwDl
				}
				if comp.MirBwDl != "" {
					bwDl, _ := pcf_context.ConvertBitRateToKbps(comp.MirBwDl)
					minBwDl += bwDl
				}
			}
		} else {
			if both || dir == models.FlowDirection_UPLINK {
				ulExist = true
				if subsComp.MarBwUl != "" {
					bwUl, _ := pcf_context.ConvertBitRateToKbps(subsComp.MarBwUl)
					maxBwUl += bwUl
				} else if comp.MarBwUl != "" {
					bwUl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
					maxBwUl += (0.05 * bwUl)
				}
			}
			if both || dir == models.FlowDirection_DOWNLINK {
				dlExist = true
				if subsComp.MarBwDl != "" {
					bwDl, _ := pcf_context.ConvertBitRateToKbps(subsComp.MarBwDl)
					maxBwDl += bwDl
				} else if comp.MarBwDl != "" {
					bwDl, _ := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
					maxBwDl += (0.05 * bwDl)
				}
			}
		}
	}

	// update Downlink MBR
	if maxBwDl == 0.0 {
		updatedQosData.MaxbrDl = comp.MarBwDl
	} else {
		updatedQosData.MaxbrDl = pcf_context.ConvertBitRateToString(maxBwDl)
	}
	// update Uplink MBR
	if maxBwUl == 0.0 {
		updatedQosData.MaxbrUl = comp.MarBwUl
	} else {
		updatedQosData.MaxbrUl = pcf_context.ConvertBitRateToString(maxBwUl)
	}
	// if gbr == 0 then assign gbr = mbr
	// update Downlink GBR
	if minBwDl != 0.0 {
		updatedQosData.GbrDl = pcf_context.ConvertBitRateToString(minBwDl)
	}
	// update Uplink GBR
	if minBwUl != 0.0 {
		updatedQosData.GbrUl = pcf_context.ConvertBitRateToString(minBwUl)
	}
	return
}

func removeMediaComp(appSession *pcf_context.AppSessionData, compN string) {
	idMaps := appSession.RelatedPccRuleIds
	smPolicy := appSession.SmPolicyData
	if idMaps != nil {
		if appSession.AppSessionContext.AscReqData.MedComponents == nil {
			return
		}
		comp, exist := appSession.AppSessionContext.AscReqData.MedComponents[compN]
		if !exist {
			return
		}
		if comp.MedSubComps != nil {
			for fNum := range comp.MedSubComps {
				key := fmt.Sprintf("%s-%s", compN, fNum)
				pccRuleId := idMaps[key]
				err := smPolicy.RemovePccRule(pccRuleId)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnf(err.Error())
				}
				delete(appSession.RelatedPccRuleIds, key)
				delete(appSession.PccRuleIdMapToCompId, pccRuleId)
			}
		} else {
			pccRuleId := idMaps[compN]
			err := smPolicy.RemovePccRule(pccRuleId)
			if err != nil {
				logger.PolicyAuthorizationlog.Warnf(err.Error())
			}
			delete(appSession.RelatedPccRuleIds, compN)
			delete(appSession.PccRuleIdMapToCompId, pccRuleId)
		}
		delete(appSession.AppSessionContext.AscReqData.MedComponents, compN)
	}
}

// func removeMediaSubComp(appSession *pcf_context.AppSessionData, compN, fNum string) {
// 	key := fmt.Sprintf("%s-%s", compN, fNum)
// 	idMaps := appSession.RelatedPccRuleIds
// 	smPolicy := appSession.SmPolicyData
// 	if idMaps != nil {
// 		if appSession.AppSessionContext.AscReqData.MedComponents == nil {
// 			return
// 		}
// 		if comp, exist := appSession.AppSessionContext.AscReqData.MedComponents[compN]; exist {
// 			pccRuleId := idMaps[key]
// 			smPolicy.RemovePccRule(pccRuleId)
// 			delete(appSession.RelatedPccRuleIds, key)
// 			delete(comp.MedSubComps, fNum)
// 			appSession.AppSessionContext.AscReqData.MedComponents[compN] = comp
// 		}
// 	}
// 	return
// }

func threshRmToThresh(threshrm *models.UsageThresholdRm) *models.UsageThreshold {
	if threshrm == nil {
		return nil
	}
	return &models.UsageThreshold{
		Duration:       threshrm.Duration,
		TotalVolume:    threshrm.TotalVolume,
		DownlinkVolume: threshrm.DownlinkVolume,
		UplinkVolume:   threshrm.UplinkVolume,
	}
}

func extractUmData(umId string, eventSubs map[models.AfEvent]models.AfNotifMethod, threshold *models.UsageThreshold) (umData *models.UsageMonitoringData, err error) {
	if _, umExist := eventSubs[models.AfEvent_USAGE_REPORT]; umExist {
		if threshold == nil {
			return nil, fmt.Errorf("UsageThreshold is nil in USAGE REPORT Subscription")

		} else {
			tmp := pcf_util.CreateUmData(umId, *threshold)
			umData = &tmp
		}
	}
	return
}

func modifyRemainBitRate(httpChannel chan pcf_message.HttpResponseMessage, smPolicy *pcf_context.UeSmPolicyData, qosData *models.QosData, ulExist, dlExist bool) (err error) {
	// if request GBR == 0, qos GBR = MBR
	// if request GBR > remain GBR, qos GBR = remain GBR
	if ulExist {
		if qosData.GbrUl == "" {
			err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.MaxbrUl)
			if err != nil {
				qosData.GbrUl = pcf_context.DecreaseRamainBitRateToZero(smPolicy.RemainGbrUL)
			} else {
				qosData.GbrUl = qosData.MaxbrUl
			}
		} else {
			err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.GbrUl)
			if err != nil {
				sendProblemDetail(httpChannel, err.Error(), pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return
			}
		}
	}
	if dlExist {
		if qosData.GbrDl == "" {
			err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrDL, qosData.MaxbrDl)
			if err != nil {
				qosData.GbrDl = pcf_context.DecreaseRamainBitRateToZero(smPolicy.RemainGbrDL)
			} else {
				qosData.GbrDl = qosData.MaxbrDl
			}
		} else {
			err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrDL, qosData.GbrDl)
			if err != nil {
				// if Policy failed, revert remain GBR to original GBR
				pcf_context.IncreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.GbrUl)
				sendProblemDetail(httpChannel, err.Error(), pcf_util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return
			}
		}
	}
	return nil
}

func InitialProvisioningOfTrafficRoutingInformation(smPolicy *pcf_context.UeSmPolicyData, pccRule *models.PccRule, compAfRoutReq, reqAfRoutReq *models.AfRoutingRequirement) {
	for _, tcId := range pccRule.RefTcData {
		tcData := smPolicy.PolicyDecision.TraffContDecs[tcId]
		if compAfRoutReq != nil {
			tcData.RouteToLocs = append(tcData.RouteToLocs, compAfRoutReq.RouteToLocs...)
			tcData.UpPathChgEvent = compAfRoutReq.UpPathChgSub
			pccRule.AppReloc = compAfRoutReq.AppReloc
		} else if reqAfRoutReq != nil {
			tcData.RouteToLocs = append(tcData.RouteToLocs, reqAfRoutReq.RouteToLocs...)
			tcData.UpPathChgEvent = reqAfRoutReq.UpPathChgSub
			pccRule.AppReloc = reqAfRoutReq.AppReloc
		}
		smPolicy.PolicyDecision.TraffContDecs[tcData.TcId] = tcData
	}
}

func ModifyProvisioningOfTrafficRoutingInformation(smPolicy *pcf_context.UeSmPolicyData, pccRule *models.PccRule, compAfRoutReq, reqAfRoutReq *models.AfRoutingRequirementRm) {
	for _, tcId := range pccRule.RefTcData {
		tcData := smPolicy.PolicyDecision.TraffContDecs[tcId]
		if compAfRoutReq != nil {
			tcData.RouteToLocs = append(tcData.RouteToLocs, compAfRoutReq.RouteToLocs...)
			tcData.UpPathChgEvent = compAfRoutReq.UpPathChgSub
			pccRule.AppReloc = compAfRoutReq.AppReloc
		} else if reqAfRoutReq != nil {
			tcData.RouteToLocs = append(tcData.RouteToLocs, reqAfRoutReq.RouteToLocs...)
			tcData.UpPathChgEvent = reqAfRoutReq.UpPathChgSub
			pccRule.AppReloc = reqAfRoutReq.AppReloc
		}
		smPolicy.PolicyDecision.TraffContDecs[tcData.TcId] = tcData
	}
}

func reverseStringMap(srcMap map[string]string) map[string]string {
	if srcMap == nil {
		return nil
	}
	reverseMap := make(map[string]string)
	for key, value := range srcMap {
		reverseMap[value] = key
	}
	return reverseMap
}

func sendProblemDetail(httpChannel chan pcf_message.HttpResponseMessage, errDetail, errCause string) {
	rsp := pcf_util.GetProblemDetail(errDetail, errCause)
	logger.PolicyAuthorizationlog.Error(rsp.Detail)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
}
