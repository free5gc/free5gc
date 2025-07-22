package processor

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cydev/zero"
	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/pcf/PolicyAuthorization"
	pcf_context "github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
)

const (
	Precedence_Initial      int32 = 22
	Precedence_Maximum      int32 = 254
	Precedence_Default_path int32 = 255
)

func transferAfRoutReqRmToAfRoutReq(AfRoutReqRm *models.AfRoutingRequirementRm) *models.AfRoutingRequirement {
	spVal := models.SpatialValidity{
		PresenceInfoList: AfRoutReqRm.SpVal.PresenceInfoList,
	}
	afRoutReq := models.AfRoutingRequirement{
		AppReloc:     AfRoutReqRm.AppReloc,
		RouteToLocs:  AfRoutReqRm.RouteToLocs,
		SpVal:        &spVal,
		TempVals:     AfRoutReqRm.TempVals,
		UpPathChgSub: AfRoutReqRm.UpPathChgSub,
	}
	return &afRoutReq
}

func transferMedCompRmToMedComp(medCompRm *models.MediaComponentRm) *models.MediaComponent {
	medSubComps := make(map[string]models.MediaSubComponent)
	for id, medSubCompRm := range medCompRm.MedSubComps {
		medSubComps[id] = models.MediaSubComponent(*medSubCompRm)
	}
	medComp := models.MediaComponent{
		AfAppId:     medCompRm.AfAppId,
		AfRoutReq:   transferAfRoutReqRmToAfRoutReq(medCompRm.AfRoutReq),
		ContVer:     medCompRm.ContVer,
		Codecs:      medCompRm.Codecs,
		FStatus:     medCompRm.FStatus,
		MarBwDl:     medCompRm.MarBwDl,
		MarBwUl:     medCompRm.MarBwUl,
		MedCompN:    medCompRm.MedCompN,
		MedSubComps: medSubComps,
		MedType:     medCompRm.MedType,
		MirBwDl:     medCompRm.MirBwDl,
		MirBwUl:     medCompRm.MirBwUl,
		ResPrio:     medCompRm.ResPrio,
	}
	return &medComp
}

// Handle Create/ Modify  Media SubComponent
func handleMediaSubComponent(smPolicy *pcf_context.UeSmPolicyData, medComp *models.MediaComponent,
	medSubComp *models.MediaSubComponent, var5qi int32,
) (*models.PccRule, *models.ProblemDetails) {
	var flowInfos []models.FlowInformation
	if tempFlowInfos, err := getFlowInfos(medSubComp); err != nil {
		problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
		return nil, &problemDetail
	} else {
		flowInfos = tempFlowInfos
	}
	pccRule := util.GetPccRuleByFlowInfos(smPolicy.PolicyDecision.PccRules, flowInfos)
	if pccRule == nil {
		precedence := getAvailablePrecedence(smPolicy.PolicyDecision.PccRules)
		pccRule = util.CreatePccRule(smPolicy.PccRuleIdGenerator, precedence, nil, "")
		// Set QoS Data
		// TODO: use real arp
		qosData := util.CreateQosData(smPolicy.PccRuleIdGenerator, var5qi, 8)
		if var5qi <= 4 {
			// update Qos Data according to request BitRate
			var ul, dl bool

			qosData, ul, dl = updateQosInMedSubComp(&qosData, medComp, medSubComp)
			if problemDetails := modifyRemainBitRate(smPolicy, &qosData, ul, dl); problemDetails != nil {
				return nil, problemDetails
			}
		}
		// Set PackfiltId
		for i := range flowInfos {
			flowInfos[i].PackFiltId = util.GetPackFiltId(smPolicy.PackFiltIdGenerator)
			smPolicy.PackFiltMapToPccRuleId[flowInfos[i].PackFiltId] = pccRule.PccRuleId
			smPolicy.PackFiltIdGenerator++
		}
		// Set flowsInfo in Pcc Rule
		pccRule.FlowInfos = flowInfos
		// Set Traffic Control Data
		tcData := util.CreateTcData(smPolicy.PccRuleIdGenerator, "", medSubComp.FStatus)
		util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, tcData, &qosData, nil, nil)
		smPolicy.PccRuleIdGenerator++
	} else {
		// update qos
		var qosData models.QosData
		for _, qosID := range pccRule.RefQosData {
			qosData = *smPolicy.PolicyDecision.QosDecs[qosID]
			if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
				var ul, dl bool
				qosData, ul, dl = updateQosInMedSubComp(smPolicy.PolicyDecision.QosDecs[qosID], medComp, medSubComp)
				if problemDetails := modifyRemainBitRate(smPolicy, &qosData, ul, dl); problemDetails != nil {
					logger.PolicyAuthLog.Error(problemDetails.Detail)
					return nil, problemDetails
				}
				smPolicy.PolicyDecision.QosDecs[qosData.QosId] = &qosData
			}
		}
	}
	smPolicy.PolicyDecision.PccRules[pccRule.PccRuleId] = pccRule
	return pccRule, nil
}

// HandlePostAppSessionsContext - Creates a new Individual Application Session Context resource
// Initial provisioning of service information (DONE)
// Gate control (DONE)
// Initial provisioning of sponsored connectivity information (DONE)
// Subscriptions to Service Data Flow QoS notification control (DONE)
// Subscription to Service Data Flow Deactivation (DONE)
// Initial provisioning of traffic routing information (DONE)
// Subscription to resources allocation outcome (DONE)
// Invocation of Multimedia Priority Services (TODO)
// Support of content versioning (TODO)
func (p *Processor) HandlePostAppSessionsContext(
	c *gin.Context,
	appSessionContext models.AppSessionContext,
) {
	logger.PolicyAuthLog.Traceln("Handle Create AppSessions")

	response, locationHeader, problemDetails := p.postAppSessCtxProcedure(&appSessionContext)

	if response != nil {
		c.Header("Location", locationHeader)
		c.JSON(http.StatusCreated, response)
		return
	} else if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	c.JSON(http.StatusForbidden, problemDetails)
}

func (p *Processor) postAppSessCtxProcedure(appSessCtx *models.AppSessionContext) (*models.AppSessionContext,
	string, *models.ProblemDetails,
) {
	ascReqData := appSessCtx.AscReqData
	pcfSelf := p.Context()

	// Initial BDT policy indication(the only one which is not related to session)
	if ascReqData.BdtRefId != "" {
		if err := p.handleBDTPolicyInd(pcfSelf, appSessCtx); err != nil {
			problemDetail := util.GetProblemDetail(err.Error(), util.ERROR_REQUEST_PARAMETERS)
			return nil, "", &problemDetail
		}
		appSessID := fmt.Sprintf("BdtRefId-%s", ascReqData.BdtRefId)
		data := pcf_context.AppSessionData{
			AppSessionId:      appSessID,
			AppSessionContext: appSessCtx,
		}
		pcfSelf.AppSessionPool.Store(appSessID, &data)
		locationHeader := util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSessID)
		logger.PolicyAuthLog.Tracef("App Session Id[%s] Create", appSessID)
		return appSessCtx, locationHeader, nil
	}
	if ascReqData.UeIpv4 == "" && ascReqData.UeIpv6 == "" && ascReqData.UeMac == "" {
		problemDetail := util.GetProblemDetail("Ue UeIpv4 and UeIpv6 and UeMac are all empty", util.ERROR_REQUEST_PARAMETERS)
		return nil, "", &problemDetail
	}
	if ascReqData.AfRoutReq != nil && ascReqData.Dnn == "" {
		problemDetail := util.GetProblemDetail("DNN shall be present", util.ERROR_REQUEST_PARAMETERS)
		return nil, "", &problemDetail
	}
	var smPolicy *pcf_context.UeSmPolicyData
	if tempSmPolicy, err := pcfSelf.SessionBinding(ascReqData); err != nil {
		problemDetail := util.GetProblemDetail(fmt.Sprintf("Session Binding failed[%s]",
			err.Error()), util.PDU_SESSION_NOT_AVAILABLE)
		return nil, "", &problemDetail
	} else {
		smPolicy = tempSmPolicy
	}
	logger.PolicyAuthLog.Infof("Session Binding Success - UeIpv4[%s], UeIpv6[%s], UeMac[%s]",
		ascReqData.UeIpv4, ascReqData.UeIpv6, ascReqData.UeMac)
	ue := smPolicy.PcfUe
	updateSMpolicy := false

	var requestSuppFeat openapi.SupportedFeature
	if tempRequestSuppFeat, err := openapi.NewSupportedFeature(ascReqData.SuppFeat); err != nil {
		logger.PolicyAuthLog.Error(err.Error())
	} else {
		requestSuppFeat = tempRequestSuppFeat
	}

	nSuppFeat := pcfSelf.PcfSuppFeats[models.ServiceName_NPCF_POLICYAUTHORIZATION].NegotiateWith(requestSuppFeat).String()
	// InfluenceOnTrafficRouting = 1 in 29514 &  Traffic Steering Control support = 1 in 29512
	traffRoutSupp := util.CheckSuppFeat(nSuppFeat, 1) && util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 1)
	relatedPccRuleIds := make(map[string]string)

	if ascReqData.MedComponents != nil {
		// Handle Pcc rules
		precedence := getAvailablePrecedence(smPolicy.PolicyDecision.PccRules)
		for _, medComp := range ascReqData.MedComponents {
			var pccRule *models.PccRule
			var appID string
			var routeReq *models.AfRoutingRequirement
			// TODO: use specific algorithm instead of default, details in subsclause 7.3.3 of TS 29513
			var var5qi int32 = 9
			if medComp.MedType != "" {
				var5qi = util.MediaTypeTo5qiMap[medComp.MedType]
			}

			if medComp.MedSubComps != nil {
				for _, medSubComp := range medComp.MedSubComps {
					if tempPccRule, problemDetail := handleMediaSubComponent(smPolicy,
						&medComp, &medSubComp, var5qi); problemDetail != nil {
						return nil, "", problemDetail
					} else {
						pccRule = tempPccRule
					}
					key := fmt.Sprintf("%d-%d", medComp.MedCompN, medSubComp.FNum)
					relatedPccRuleIds[key] = pccRule.PccRuleId
					updateSMpolicy = true
				}
				continue
			} else if medComp.AfAppId != "" {
				appID = medComp.AfAppId
				routeReq = medComp.AfRoutReq
			} else if ascReqData.AfAppId != "" {
				appID = ascReqData.AfAppId
				routeReq = ascReqData.AfRoutReq
			} else {
				problemDetail := util.GetProblemDetail("Media Component needs flows of subComp or afAppId",
					util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return nil, "", &problemDetail
			}

			// Find pccRule by AfAppId, otherwise create a new pcc rule
			pccRule = util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, appID)
			if pccRule == nil {
				pccRule = util.CreatePccRule(smPolicy.PccRuleIdGenerator, precedence, nil, appID)
				// Set QoS Data
				// TODO: use real arp
				qosData := util.CreateQosData(smPolicy.PccRuleIdGenerator, var5qi, 8)
				if var5qi <= 4 {
					// update Qos Data according to request BitRate
					var ul, dl bool
					qosData, ul, dl = updateQosInMedComp(qosData, &medComp)
					if problemDetails := modifyRemainBitRate(smPolicy, &qosData, ul, dl); problemDetails != nil {
						return nil, "", problemDetails
					}
				}
				util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, nil, &qosData, nil, nil)
				smPolicy.PccRuleIdGenerator++
				if precedence < Precedence_Maximum {
					precedence++
				}
			} else {
				// update pccRule's qos
				var qosData models.QosData
				for _, qosID := range pccRule.RefQosData {
					qosData = *smPolicy.PolicyDecision.QosDecs[qosID]
					if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
						var ul, dl bool
						qosData, ul, dl = updateQosInMedComp(*smPolicy.PolicyDecision.QosDecs[qosID], &medComp)
						if problemDetails := modifyRemainBitRate(smPolicy, &qosData, ul, dl); problemDetails != nil {
							return nil, "", problemDetails
						}
						smPolicy.PolicyDecision.QosDecs[qosData.QosId] = &qosData
					}
				}
			}
			// Initial provisioning of traffic routing information
			if traffRoutSupp {
				pccRule = provisioningOfTrafficRoutingInfo(smPolicy, appID, routeReq, medComp.FStatus)
			}
			key := fmt.Sprintf("%d", medComp.MedCompN)
			relatedPccRuleIds[key] = pccRule.PccRuleId
			updateSMpolicy = true
		}
	} else if ascReqData.AfAppId != "" {
		// Initial provisioning of traffic routing information
		if ascReqData.AfRoutReq != nil && traffRoutSupp {
			logger.PolicyAuthLog.Infof("AF influence on Traffic Routing - AppId[%s]", ascReqData.AfAppId)
			pccRule := provisioningOfTrafficRoutingInfo(smPolicy, ascReqData.AfAppId, ascReqData.AfRoutReq, "")
			key := fmt.Sprintf("appID-%s", ascReqData.AfAppId)
			relatedPccRuleIds[key] = pccRule.PccRuleId
			updateSMpolicy = true
		} else {
			problemDetail := util.GetProblemDetail("Traffic routing not supported", util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return nil, "", &problemDetail
		}
	} else {
		problemDetail := util.GetProblemDetail("AF Request need AfAppId or Media Component to match Service Data Flow",
			util.ERROR_REQUEST_PARAMETERS)
		return nil, "", &problemDetail
	}

	// Event Subscription
	eventSubs := make(map[models.PcfPolicyAuthorizationAfEvent]models.AfNotifMethod)
	if ascReqData.EvSubsc != nil {
		for _, subs := range ascReqData.EvSubsc.Events {
			if subs.NotifMethod == "" {
				// default value "EVENT_DETECTION"
				subs.NotifMethod = models.AfNotifMethod_EVENT_DETECTION
			}
			eventSubs[subs.Event] = subs.NotifMethod
			var trig models.PolicyControlRequestTrigger
			switch subs.Event {
			case models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE:
				trig = models.PolicyControlRequestTrigger_AC_TY_CH
			// case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
			// 	// Subscription to Service Data Flow Deactivation
			// 	trig = models.PolicyControlRequestTrigger_RES_RELEASE
			case models.PcfPolicyAuthorizationAfEvent_PLMN_CHG:
				trig = models.PolicyControlRequestTrigger_PLMN_CH
			case models.PcfPolicyAuthorizationAfEvent_QOS_NOTIF:
				// Subscriptions to Service Data Flow QoS notification control
				for _, pccRuleID := range relatedPccRuleIds {
					pccRule := smPolicy.PolicyDecision.PccRules[pccRuleID]
					for _, qosID := range pccRule.RefQosData {
						qosData := smPolicy.PolicyDecision.QosDecs[qosID]
						qosData.Qnc = true
						smPolicy.PolicyDecision.QosDecs[qosID] = qosData
					}
				}
				trig = models.PolicyControlRequestTrigger_QOS_NOTIF
			case models.PcfPolicyAuthorizationAfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
				// Subscription to resources allocation outcome
				trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
			case models.PcfPolicyAuthorizationAfEvent_USAGE_REPORT:
				trig = models.PolicyControlRequestTrigger_US_RE
			default:
				logger.PolicyAuthLog.Warn("AF Event is unknown")
				continue
			}
			if !util.CheckPolicyControlReqTrig(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig) {
				smPolicy.PolicyDecision.PolicyCtrlReqTriggers = append(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig)
				updateSMpolicy = true
			}
		}
	}

	// Initial provisioning of sponsored connectivity information
	if ascReqData.AspId != "" && ascReqData.SponId != "" {
		// SponsoredConnectivity = 2 in 29514 &  SponsoredConnectivity support = 12 in 29512
		supp := util.CheckSuppFeat(nSuppFeat, 2) && util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 12)
		if !supp {
			problemDetail := util.GetProblemDetail("Sponsored Connectivity not supported", util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return nil, "", &problemDetail
		}
		umID := util.GetUmId(ascReqData.AspId, ascReqData.SponId)
		var umData *models.UsageMonitoringData
		if tempUmData, err := extractUmData(umID, eventSubs, ascReqData.EvSubsc.UsgThres); err != nil {
			problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return nil, "", &problemDetail
		} else {
			umData = tempUmData
		}
		if err := handleSponsoredConnectivityInformation(smPolicy, relatedPccRuleIds, ascReqData.AspId,
			ascReqData.SponId, ascReqData.SponStatus, umData, &updateSMpolicy); err != nil {
			problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			return nil, "", &problemDetail
		}
	}

	// Allocate App Session Id
	appSessID := ue.AllocUeAppSessionId(pcfSelf)
	appSessCtx.AscRespData = &models.AppSessionContextRespData{
		SuppFeat: nSuppFeat,
	}
	// Associate App Session to SMPolicy
	smPolicy.AppSessions[appSessID] = true
	data := pcf_context.AppSessionData{
		AppSessionId:      appSessID,
		AppSessionContext: appSessCtx,
		SmPolicyData:      smPolicy,
	}
	if len(relatedPccRuleIds) > 0 {
		data.RelatedPccRuleIds = relatedPccRuleIds
		data.PccRuleIdMapToCompId = reverseStringMap(relatedPccRuleIds)
	}
	appSessCtx.EvsNotif = &models.PcfPolicyAuthorizationEventsNotification{}
	// Set Event Subsciption related Data
	if len(eventSubs) > 0 {
		data.Events = eventSubs
		data.EventUri = ascReqData.EvSubsc.NotifUri
		if _, exist := eventSubs[models.PcfPolicyAuthorizationAfEvent_PLMN_CHG]; exist {
			afNotif := models.PcfPolicyAuthorizationAfEventNotification{
				Event: models.PcfPolicyAuthorizationAfEvent_PLMN_CHG,
			}
			appSessCtx.EvsNotif.EvNotifs = append(appSessCtx.EvsNotif.EvNotifs, afNotif)
			plmnID := smPolicy.PolicyContext.ServingNetwork
			if plmnID != nil {
				appSessCtx.EvsNotif.PlmnId = &models.PlmnIdNid{
					Mcc: plmnID.Mcc,
					Mnc: plmnID.Mnc,
				}
			}
		}
		if _, exist := eventSubs[models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE]; exist {
			afNotif := models.PcfPolicyAuthorizationAfEventNotification{
				Event: models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE,
			}
			appSessCtx.EvsNotif.EvNotifs = append(appSessCtx.EvsNotif.EvNotifs, afNotif)
			appSessCtx.EvsNotif.AccessType = smPolicy.PolicyContext.AccessType
			appSessCtx.EvsNotif.RatType = smPolicy.PolicyContext.RatType
		}
	}
	if appSessCtx.EvsNotif.EvNotifs == nil {
		appSessCtx.EvsNotif = nil
	}
	pcfSelf.AppSessionPool.Store(appSessID, &data)
	locationHeader := util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSessID)
	logger.PolicyAuthLog.Infof("App Session Id[%s] Create", appSessID)
	// Send Notification to SMF
	if updateSMpolicy {
		smPolicyID := fmt.Sprintf("%s-%d", ue.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyID),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &notification)
	}
	return appSessCtx, locationHeader, nil
}

// HandleDeleteAppSession - Deletes an existing Individual Application Session Context
func (p *Processor) HandleDeleteAppSessionContext(
	c *gin.Context,
	appSessionId string,
	eventsSubscReqData *models.PcfPolicyAuthorizationEventsSubscReqData,
) {
	logger.PolicyAuthLog.Infof("Handle Del AppSessions, AppSessionId[%s]", appSessionId)

	pcfSelf := p.Context()
	var appSession *pcf_context.AppSessionData
	if val, ok := pcfSelf.AppSessionPool.Load(appSessionId); ok {
		appSession = val.(*pcf_context.AppSessionData)
	}
	if appSession == nil {
		problemDetail := util.GetProblemDetail("can't find app session", util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}
	if eventsSubscReqData != nil {
		logger.PolicyAuthLog.Warn("Delete AppSessions does not support with Event Subscription")
	}
	// Remove related pcc rule resource
	smPolicy := appSession.SmPolicyData
	deletedSmPolicyDec := models.SmPolicyDecision{}
	for _, pccRuleID := range appSession.RelatedPccRuleIds {
		if err := smPolicy.RemovePccRule(pccRuleID, &deletedSmPolicyDec); err != nil {
			logger.PolicyAuthLog.Warn(err.Error())
		}
	}

	delete(smPolicy.AppSessions, appSessionId)

	logger.PolicyAuthLog.Infof("App Session Id[%s] Del", appSessionId)

	// TODO: AccUsageReport
	// if appSession.AccUsage != nil {

	// 	resp := models.AppSessionContext{
	// 		EvsNotif: &models.PcfPolicyAuthorizationEventsNotification{
	// 			UsgRep: appSession.AccUsage,
	// 		},
	// 	}
	// 	message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, resp)
	// } else {
	// }

	pcfSelf.AppSessionPool.Delete(appSessionId)

	smPolicy.ArrangeExistEventSubscription()

	// Notify SMF About Pcc Rule moval
	smPolicyID := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
	notification := models.SmPolicyNotification{
		ResourceUri:      util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyID),
		SmPolicyDecision: &deletedSmPolicyDec,
	}
	go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &notification)
	logger.PolicyAuthLog.Tracef("Send SM Policy[%s] Update Notification", smPolicyID)
	c.JSON(http.StatusNoContent, nil)
}

// HandleGetAppSession - Reads an existing Individual Application Session Context
func (p *Processor) HandleGetAppSessionContext(
	c *gin.Context,
	appSessionId string,
) {
	logger.PolicyAuthLog.Infof("Handle Get AppSessions, AppSessionId[%s]", appSessionId)

	pcfSelf := p.Context()

	var appSession *pcf_context.AppSessionData
	if val, ok := pcfSelf.AppSessionPool.Load(appSessionId); ok {
		appSession = val.(*pcf_context.AppSessionData)
	}
	if appSession == nil {
		problemDetail := util.GetProblemDetail("can't find app session", util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}
	logger.PolicyAuthLog.Tracef("App Session Id[%s] Get", appSessionId)
	c.JSON(http.StatusOK, appSession.AppSessionContext)
}

// HandleModAppSession - Modifies an existing Individual Application Session Context
func (p *Processor) HandleModAppSessionContext(
	c *gin.Context,
	appSessionId string,
	appSessionContextUpdateData models.AppSessionContextUpdateData,
) {
	logger.PolicyAuthLog.Infof("Handle Modify AppSessions, AppSessionId[%s]", appSessionId)

	pcfSelf := p.Context()
	var appSession *pcf_context.AppSessionData
	if val, ok := pcfSelf.AppSessionPool.Load(appSessionId); ok {
		appSession = val.(*pcf_context.AppSessionData)
	}
	if appSession == nil {
		problemDetail := util.GetProblemDetail("can't find app session", util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		c.JSON((int)(problemDetail.Status), problemDetail)
		return
	}
	appSessCtx := appSession.AppSessionContext
	if appSessionContextUpdateData.BdtRefId != "" {
		appSessCtx.AscReqData.BdtRefId = appSessionContextUpdateData.BdtRefId
		if err := p.handleBDTPolicyInd(pcfSelf, appSessCtx); err != nil {
			problemDetail := util.GetProblemDetail(err.Error(), util.ERROR_REQUEST_PARAMETERS)
			c.JSON(int(problemDetail.Status), problemDetail)
			return
		}
		logger.PolicyAuthLog.Tracef("App Session Id[%s] Updated", appSessionId)
		c.JSON(http.StatusOK, appSessCtx)
		return
	}
	smPolicy := appSession.SmPolicyData
	if smPolicy == nil {
		problemDetail := util.GetProblemDetail("Can't find related PDU Session", util.REQUESTED_SERVICE_NOT_AUTHORIZED)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}
	// InfluenceOnTrafficRouting = 1 in 29514 &  Traffic Steering Control support = 1 in 29512
	traffRoutSupp := util.CheckSuppFeat(appSessCtx.AscRespData.SuppFeat,
		1) && util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 1)
	relatedPccRuleIds := make(map[string]string)
	// Event Subscription
	eventSubs := make(map[models.PcfPolicyAuthorizationAfEvent]models.AfNotifMethod)
	updateSMpolicy := false

	if appSessionContextUpdateData.MedComponents != nil {
		precedence := getAvailablePrecedence(smPolicy.PolicyDecision.PccRules)
		for compN, medCompRm := range appSessionContextUpdateData.MedComponents {
			medComp := transferMedCompRmToMedComp(medCompRm)
			removeMediaComp(appSession, compN)
			if zero.IsZero(medComp) {
				// remove MediaComp(media Comp is null)
				continue
			}
			// modify MediaComp(remove and reinstall again)
			var pccRule *models.PccRule
			var appID string
			var routeReq *models.AfRoutingRequirement
			// TODO: use specific algorithm instead of default, details in subsclause 7.3.3 of TS 29513
			var var5qi int32 = 9
			if medComp.MedType != "" {
				var5qi = util.MediaTypeTo5qiMap[medComp.MedType]
			}
			if medComp.MedSubComps != nil {
				for _, medSubComp := range medComp.MedSubComps {
					if tempPccRule, problemDetail := handleMediaSubComponent(smPolicy, medComp,
						&medSubComp, var5qi); problemDetail != nil {
						c.JSON(int(problemDetail.Status), problemDetail)
						return
					} else {
						pccRule = tempPccRule
					}
					key := fmt.Sprintf("%d-%d", medComp.MedCompN, medSubComp.FNum)
					relatedPccRuleIds[key] = pccRule.PccRuleId
					updateSMpolicy = true
				}
				continue
			} else if medComp.AfAppId != "" {
				// if medComp.AfAppId has value -> find pccRule by reqData.AfAppId, otherwise create a new pcc rule
				appID = medComp.AfAppId
				routeReq = medComp.AfRoutReq
			} else if appSessionContextUpdateData.AfAppId != "" {
				appID = appSessionContextUpdateData.AfAppId
				routeReq = medComp.AfRoutReq
			} else {
				problemDetail := util.GetProblemDetail("Media Component needs flows of subComp or afAppId",
					util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				c.JSON(int(problemDetail.Status), problemDetail)
				return
			}

			pccRule = util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, appID)
			if pccRule == nil { // create new pcc rule
				pccRule = util.CreatePccRule(smPolicy.PccRuleIdGenerator, precedence, nil, appID)
				// Set QoS Data
				// TODO: use real arp
				qosData := util.CreateQosData(smPolicy.PccRuleIdGenerator, var5qi, 8)
				if var5qi <= 4 {
					// update Qos Data according to request BitRate
					var ul, dl bool
					qosData, ul, dl = updateQosInMedComp(qosData, medComp)
					if problemDetail := modifyRemainBitRate(smPolicy, &qosData, ul, dl); problemDetail != nil {
						c.JSON(int(problemDetail.Status), problemDetail)
						return
					}
				}
				util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, nil, &qosData, nil, nil)
				smPolicy.PccRuleIdGenerator++
				if precedence < Precedence_Maximum {
					precedence++
				}
			} else {
				// update qos
				var qosData models.QosData
				for _, qosID := range pccRule.RefQosData {
					qosData = *smPolicy.PolicyDecision.QosDecs[qosID]
					if qosData.Var5qi == var5qi && qosData.Var5qi <= 4 {
						var ul, dl bool
						qosData, ul, dl = updateQosInMedComp(*smPolicy.PolicyDecision.QosDecs[qosID], medComp)
						if problemDetail := modifyRemainBitRate(smPolicy, &qosData, ul, dl); problemDetail != nil {
							c.JSON(int(problemDetail.Status), problemDetail)
							return
						}
						smPolicy.PolicyDecision.QosDecs[qosData.QosId] = &qosData
					}
				}
			}
			// Modify provisioning of traffic routing information
			if traffRoutSupp {
				pccRule = provisioningOfTrafficRoutingInfo(smPolicy, appID, routeReq, medComp.FStatus)
			}
			key := fmt.Sprintf("%d", medComp.MedCompN)
			relatedPccRuleIds[key] = pccRule.PccRuleId
			updateSMpolicy = true
		}
	}

	// Update of traffic routing information
	// TODO: check ascUpdateData.AfAppId with appSessCtx.AscReqData.AfAppId (now ascUpdateData.AfAppId is empty)
	if appSessionContextUpdateData.AfRoutReq != nil && traffRoutSupp {
		logger.PolicyAuthLog.Infof("Update Traffic Routing info - [%+v]", appSessionContextUpdateData.AfRoutReq)
		appSessCtx.AscReqData.AfRoutReq = transferAfRoutReqRmToAfRoutReq(appSessionContextUpdateData.AfRoutReq)
		// Update SmPolicyDecision
		pccRule := provisioningOfTrafficRoutingInfo(smPolicy,
			appSessCtx.AscReqData.AfAppId, appSessCtx.AscReqData.AfRoutReq, "")
		key := fmt.Sprintf("appID-%s", appSessCtx.AscReqData.AfAppId)
		relatedPccRuleIds[key] = pccRule.PccRuleId
		updateSMpolicy = true
	}

	// Merge Original PccRuleId and new
	for key, pccRuleID := range appSession.RelatedPccRuleIds {
		relatedPccRuleIds[key] = pccRuleID
	}

	if appSessionContextUpdateData.EvSubsc != nil {
		for _, subs := range appSessionContextUpdateData.EvSubsc.Events {
			if subs.NotifMethod == "" {
				// default value "EVENT_DETECTION"
				subs.NotifMethod = models.AfNotifMethod_EVENT_DETECTION
			}
			eventSubs[subs.Event] = subs.NotifMethod
			var trig models.PolicyControlRequestTrigger
			switch subs.Event {
			case models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE:
				trig = models.PolicyControlRequestTrigger_AC_TY_CH
			// case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
			// 	// Subscription to Service Data Flow Deactivation
			// 	trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
			case models.PcfPolicyAuthorizationAfEvent_PLMN_CHG:
				trig = models.PolicyControlRequestTrigger_PLMN_CH
			case models.PcfPolicyAuthorizationAfEvent_QOS_NOTIF:
				// Subscriptions to Service Data Flow QoS notification control
				for _, pccRuleID := range relatedPccRuleIds {
					pccRule := smPolicy.PolicyDecision.PccRules[pccRuleID]
					for _, qosID := range pccRule.RefQosData {
						qosData := smPolicy.PolicyDecision.QosDecs[qosID]
						qosData.Qnc = true
						smPolicy.PolicyDecision.QosDecs[qosID] = qosData
					}
				}
				trig = models.PolicyControlRequestTrigger_QOS_NOTIF
			case models.PcfPolicyAuthorizationAfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
				// Subscription to resources allocation outcome
				trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
			case models.PcfPolicyAuthorizationAfEvent_USAGE_REPORT:
				trig = models.PolicyControlRequestTrigger_US_RE
			default:
				logger.PolicyAuthLog.Warn("AF Event is unknown")
				continue
			}
			if !util.CheckPolicyControlReqTrig(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig) {
				smPolicy.PolicyDecision.PolicyCtrlReqTriggers = append(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig)
				updateSMpolicy = true
			}
		}
		// update Context
		if appSessCtx.AscReqData.EvSubsc == nil {
			appSessCtx.AscReqData.EvSubsc = new(models.PcfPolicyAuthorizationEventsSubscReqData)
		}
		appSessCtx.AscReqData.EvSubsc.Events = appSessionContextUpdateData.EvSubsc.Events
		if appSessionContextUpdateData.EvSubsc.NotifUri != "" {
			appSessCtx.AscReqData.EvSubsc.NotifUri = appSessionContextUpdateData.EvSubsc.NotifUri
			appSession.EventUri = appSessionContextUpdateData.EvSubsc.NotifUri
		}
		if appSessionContextUpdateData.EvSubsc.UsgThres != nil {
			appSessCtx.AscReqData.EvSubsc.UsgThres = threshRmToThresh(appSessionContextUpdateData.EvSubsc.UsgThres)
		}
	} else {
		// remove eventSubs
		appSession.Events = nil
		appSession.EventUri = ""
		appSessCtx.AscReqData.EvSubsc = nil
	}

	// Moification provisioning of sponsored connectivity information
	if appSessionContextUpdateData.AspId != "" && appSessionContextUpdateData.SponId != "" {
		umID := util.GetUmId(appSessionContextUpdateData.AspId, appSessionContextUpdateData.SponId)
		var umData *models.UsageMonitoringData
		if tempUmData, err := extractUmData(umID, eventSubs,
			threshRmToThresh(appSessionContextUpdateData.EvSubsc.UsgThres)); err != nil {
			problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			c.JSON(int(problemDetail.Status), problemDetail)
			return
		} else {
			umData = tempUmData
		}
		if err := handleSponsoredConnectivityInformation(smPolicy, relatedPccRuleIds, appSessionContextUpdateData.AspId,
			appSessionContextUpdateData.SponId, appSessionContextUpdateData.SponStatus, umData, &updateSMpolicy); err != nil {
			problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
			c.JSON(int(problemDetail.Status), problemDetail)
			return
		}
	}

	if len(relatedPccRuleIds) > 0 {
		appSession.RelatedPccRuleIds = relatedPccRuleIds
		appSession.PccRuleIdMapToCompId = reverseStringMap(relatedPccRuleIds)
	}
	appSessCtx.EvsNotif = &models.PcfPolicyAuthorizationEventsNotification{}
	// Set Event Subsciption related Data
	if len(eventSubs) > 0 {
		appSession.Events = eventSubs
		if _, exist := eventSubs[models.PcfPolicyAuthorizationAfEvent_PLMN_CHG]; exist {
			afNotif := models.PcfPolicyAuthorizationAfEventNotification{
				Event: models.PcfPolicyAuthorizationAfEvent_PLMN_CHG,
			}
			appSessCtx.EvsNotif.EvNotifs = append(appSessCtx.EvsNotif.EvNotifs, afNotif)
			plmnID := smPolicy.PolicyContext.ServingNetwork
			if plmnID != nil {
				appSessCtx.EvsNotif.PlmnId = &models.PlmnIdNid{
					Mcc: plmnID.Mcc,
					Mnc: plmnID.Mnc,
				}
			}
		}
		if _, exist := eventSubs[models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE]; exist {
			afNotif := models.PcfPolicyAuthorizationAfEventNotification{
				Event: models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE,
			}
			appSessCtx.EvsNotif.EvNotifs = append(appSessCtx.EvsNotif.EvNotifs, afNotif)
			appSessCtx.EvsNotif.AccessType = smPolicy.PolicyContext.AccessType
			appSessCtx.EvsNotif.RatType = smPolicy.PolicyContext.RatType
		}
	}
	if appSessCtx.EvsNotif.EvNotifs == nil {
		appSessCtx.EvsNotif = nil
	}

	// TODO: MPS Service
	logger.PolicyAuthLog.Tracef("App Session Id[%s] Updated", appSessionId)

	smPolicy.ArrangeExistEventSubscription()

	// Send Notification to SMF
	if updateSMpolicy {
		smPolicyID := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyID),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &notification)
		logger.PolicyAuthLog.Tracef("Send SM Policy[%s] Update Notification", smPolicyID)
	}
	c.JSON(http.StatusOK, appSessCtx)
}

// HandleDeleteEventsSubsc - deletes the Events Subscription subresource
func (p *Processor) HandleDeleteEventsSubscContext(
	c *gin.Context,
	appSessionId string,
) {
	logger.PolicyAuthLog.Tracef("Handle Del AppSessions Events Subsc, AppSessionId[%s]", appSessionId)

	pcfSelf := p.Context()
	var appSession *pcf_context.AppSessionData
	if val, ok := pcfSelf.AppSessionPool.Load(appSessionId); ok {
		appSession = val.(*pcf_context.AppSessionData)
	}
	if appSession == nil {
		problemDetail := util.GetProblemDetail("can't find app session", util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}

	appSession.Events = nil
	appSession.EventUri = ""
	appSession.AppSessionContext.EvsNotif = nil
	appSession.AppSessionContext.AscReqData.EvSubsc = nil

	logger.PolicyAuthLog.Tracef("App Session Id[%s] Del Events Subsc success", appSessionId)

	smPolicy := appSession.SmPolicyData
	// Send Notification to SMF
	if changed := appSession.SmPolicyData.ArrangeExistEventSubscription(); changed {
		smPolicyID := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyID),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &notification)
		logger.PolicyAuthLog.Tracef("Send SM Policy[%s] Update Notification", smPolicyID)
	}
	c.JSON(http.StatusNoContent, nil)
}

// HandleUpdateEventsSubsc - creates or modifies an Events Subscription subresource
func (p *Processor) HandleUpdateEventsSubscContext(
	c *gin.Context,
	appSessionId string,
	eventsSubscReqData models.PcfPolicyAuthorizationEventsSubscReqData,
) {
	logger.PolicyAuthLog.Tracef("Handle Put AppSessions Events Subsc, AppSessionId[%s]", appSessionId)

	pcfSelf := p.Context()

	var appSession *pcf_context.AppSessionData
	if val, ok := pcfSelf.AppSessionPool.Load(appSessionId); ok {
		appSession = val.(*pcf_context.AppSessionData)
	}
	if appSession == nil {
		problemDetail := util.GetProblemDetail("can't find app session", util.APPLICATION_SESSION_CONTEXT_NOT_FOUND)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}
	smPolicy := appSession.SmPolicyData
	eventSubs := make(map[models.PcfPolicyAuthorizationAfEvent]models.AfNotifMethod)

	updataSmPolicy := false
	created := appSession.Events == nil

	for _, subs := range eventsSubscReqData.Events {
		if subs.NotifMethod == "" {
			// default value "EVENT_DETECTION"
			subs.NotifMethod = models.AfNotifMethod_EVENT_DETECTION
		}
		eventSubs[subs.Event] = subs.NotifMethod
		var trig models.PolicyControlRequestTrigger
		switch subs.Event {
		case models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE:
			trig = models.PolicyControlRequestTrigger_AC_TY_CH
		// case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
		// 	// Subscription to Service Data Flow Deactivation
		// 	trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
		case models.PcfPolicyAuthorizationAfEvent_PLMN_CHG:
			trig = models.PolicyControlRequestTrigger_PLMN_CH
		case models.PcfPolicyAuthorizationAfEvent_QOS_NOTIF:
			// Subscriptions to Service Data Flow QoS notification control
			for _, pccRuleID := range appSession.RelatedPccRuleIds {
				pccRule := smPolicy.PolicyDecision.PccRules[pccRuleID]
				for _, qosID := range pccRule.RefQosData {
					qosData := smPolicy.PolicyDecision.QosDecs[qosID]
					qosData.Qnc = true
					smPolicy.PolicyDecision.QosDecs[qosID] = qosData
				}
			}
			trig = models.PolicyControlRequestTrigger_QOS_NOTIF
		case models.PcfPolicyAuthorizationAfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
			// Subscription to resources allocation outcome
			trig = models.PolicyControlRequestTrigger_SUCC_RES_ALLO
		case models.PcfPolicyAuthorizationAfEvent_USAGE_REPORT:
			trig = models.PolicyControlRequestTrigger_US_RE
		default:
			logger.PolicyAuthLog.Warn("AF Event is unknown")
			continue
		}
		if !util.CheckPolicyControlReqTrig(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig) {
			smPolicy.PolicyDecision.PolicyCtrlReqTriggers = append(smPolicy.PolicyDecision.PolicyCtrlReqTriggers, trig)
			updataSmPolicy = true
		}
	}
	appSessCtx := appSession.AppSessionContext
	// update Context
	if appSessCtx.AscReqData.EvSubsc == nil {
		appSessCtx.AscReqData.EvSubsc = new(models.PcfPolicyAuthorizationEventsSubscReqData)
	}
	appSessCtx.AscReqData.EvSubsc.Events = eventsSubscReqData.Events
	appSessCtx.AscReqData.EvSubsc.UsgThres = eventsSubscReqData.UsgThres
	appSessCtx.AscReqData.EvSubsc.NotifUri = eventsSubscReqData.NotifUri
	appSessCtx.EvsNotif = nil
	// update app Session
	appSession.EventUri = eventsSubscReqData.NotifUri
	appSession.Events = eventSubs

	appSessCtx.EvsNotif = &models.PcfPolicyAuthorizationEventsNotification{
		EvSubsUri: eventsSubscReqData.NotifUri,
	}
	// Set Event Subsciption related Data
	if len(eventSubs) > 0 {
		if _, exist := eventSubs[models.PcfPolicyAuthorizationAfEvent_PLMN_CHG]; exist {
			afNotif := models.PcfPolicyAuthorizationAfEventNotification{
				Event: models.PcfPolicyAuthorizationAfEvent_PLMN_CHG,
			}
			appSessCtx.EvsNotif.EvNotifs = append(appSessCtx.EvsNotif.EvNotifs, afNotif)
			plmnID := smPolicy.PolicyContext.ServingNetwork
			if plmnID != nil {
				appSessCtx.EvsNotif.PlmnId = &models.PlmnIdNid{
					Mcc: plmnID.Mcc,
					Mnc: plmnID.Mnc,
				}
			}
		}
		if _, exist := eventSubs[models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE]; exist {
			afNotif := models.PcfPolicyAuthorizationAfEventNotification{
				Event: models.PcfPolicyAuthorizationAfEvent_ACCESS_TYPE_CHANGE,
			}
			appSessCtx.EvsNotif.EvNotifs = append(appSessCtx.EvsNotif.EvNotifs, afNotif)
			appSessCtx.EvsNotif.AccessType = smPolicy.PolicyContext.AccessType
			appSessCtx.EvsNotif.RatType = smPolicy.PolicyContext.RatType
		}
	}
	if appSessCtx.EvsNotif.EvNotifs == nil {
		appSessCtx.EvsNotif = nil
	}

	resp := models.EventsSubscPutData{
		Events:                    eventsSubscReqData.Events,
		NotifUri:                  eventsSubscReqData.NotifUri,
		ReqQosMonParams:           eventsSubscReqData.ReqQosMonParams,
		QosMon:                    eventsSubscReqData.QosMon,
		ReqAnis:                   eventsSubscReqData.ReqAnis,
		UsgThres:                  eventsSubscReqData.UsgThres,
		NotifCorreId:              eventsSubscReqData.NotifCorreId,
		AfAppIds:                  eventsSubscReqData.AfAppIds,
		DirectNotifInd:            eventsSubscReqData.DirectNotifInd,
		AccessType:                appSessCtx.EvsNotif.AccessType,
		AnGwAddr:                  appSessCtx.EvsNotif.AnGwAddr,
		EvSubsUri:                 appSessCtx.EvsNotif.EvSubsUri,
		EvNotifs:                  appSessCtx.EvsNotif.EvNotifs,
		FailedResourcAllocReports: appSessCtx.EvsNotif.FailedResourcAllocReports,
		PlmnId:                    appSessCtx.EvsNotif.PlmnId,
		QncReports:                appSessCtx.EvsNotif.QncReports,
		RatType:                   appSessCtx.EvsNotif.RatType,
		UsgRep:                    appSessCtx.EvsNotif.UsgRep,
	}

	changed := appSession.SmPolicyData.ArrangeExistEventSubscription()

	// Send Notification to SMF
	if updataSmPolicy || changed {
		smPolicyID := fmt.Sprintf("%s-%d", smPolicy.PcfUe.Supi, smPolicy.PolicyContext.PduSessionId)
		notification := models.SmPolicyNotification{
			ResourceUri:      util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyID),
			SmPolicyDecision: smPolicy.PolicyDecision,
		}
		go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &notification)
		logger.PolicyAuthLog.Tracef("Send SM Policy[%s] Update Notification", smPolicyID)
	}
	if created {
		locationHeader := fmt.Sprintf("%s/events-subscription",
			util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSessionId))
		logger.PolicyAuthLog.Tracef("App Session Id[%s] Create Subscription", appSessionId)
		c.Header("Location", locationHeader)
		c.JSON(http.StatusCreated, resp)
	} else if appSessCtx.EvsNotif != nil {
		logger.PolicyAuthLog.Tracef("App Session Id[%s] Modify Subscription", appSessionId)
		c.JSON(http.StatusOK, resp)
	} else {
		logger.PolicyAuthLog.Tracef("App Session Id[%s] Modify Subscription", appSessionId)
		c.JSON(http.StatusNoContent, nil)
	}
}

func (p *Processor) SendAppSessionEventNotification(appSession *pcf_context.AppSessionData,
	request models.PcfPolicyAuthorizationEventsNotification,
) {
	logger.PolicyAuthLog.Tracef("Send App Session Event Notification")
	if appSession == nil {
		logger.PolicyAuthLog.Warnln("Send App Session Event Notification Error[appSession is nil]")
		return
	}
	uri := appSession.EventUri

	if uri != "" {
		ctx, _, err := p.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NrfNfManagementNfType_PCF)
		if err != nil {
			return
		}

		request.EvSubsUri = fmt.Sprintf("%s/events-subscription",
			util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSession.AppSessionId))
		client := util.GetNpcfPolicyAuthorizationCallbackClient()
		req := PolicyAuthorization.UpdateEventsSubscEventNotificationPostRequest{
			PcfPolicyAuthorizationEventsNotification: &request,
		}
		rsp, err := client.EventsSubscriptionDocumentApi.UpdateEventsSubscEventNotificationPost(
			ctx, uri, &req)
		if err != nil {
			if rsp == nil {
				logger.PolicyAuthLog.Warnf("Send App Session Event Notification Error")
			} else {
				logger.PolicyAuthLog.Warnf("Send App Session Event Notification Failed[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.PolicyAuthLog.Warnln("Send App Session Event Notification Failed[HTTP Response is nil]")
			return
		}
	}
}

func (p *Processor) SendAppSessionTermination(appSession *pcf_context.AppSessionData, request models.TerminationInfo) {
	logger.PolicyAuthLog.Tracef("Send App Session Termination")
	if appSession == nil {
		logger.PolicyAuthLog.Warnln("Send App Session Termination Error[appSession is nil]")
		return
	}
	uri := appSession.AppSessionContext.AscReqData.NotifUri

	if uri != "" {
		ctx, _, err := p.Context().GetTokenCtx(models.ServiceName_NPCF_POLICYAUTHORIZATION, models.NrfNfManagementNfType_PCF)
		if err != nil {
			return
		}

		request.ResUri = util.GetResourceUri(models.ServiceName_NPCF_POLICYAUTHORIZATION, appSession.AppSessionId)
		client := util.GetNpcfPolicyAuthorizationCallbackClient()
		req := PolicyAuthorization.PostAppSessionsTerminationRequestPostRequest{
			TerminationInfo: &request,
		}
		rsp, err := client.ApplicationSessionsCollectionApi.PostAppSessionsTerminationRequestPost(
			ctx, uri, &req)
		if err != nil {
			if rsp != nil {
				logger.PolicyAuthLog.Warnf("Send App Session Termination Error")
			} else {
				logger.PolicyAuthLog.Warnf("Send App Session Termination Failed[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.PolicyAuthLog.Warnln("Send App Session Termination Failed[HTTP Response is nil]")
			return
		}
	}
}

// Handle Create/ Modify Background Data Transfer Policy Indication
func (p *Processor) handleBDTPolicyInd(pcfSelf *pcf_context.PCFContext,
	appSessCtx *models.AppSessionContext,
) (err error) {
	req := appSessCtx.AscReqData

	var requestSuppFeat openapi.SupportedFeature
	if tempRequestSuppFeat, err1 := openapi.NewSupportedFeature(req.SuppFeat); err1 != nil {
		logger.PolicyAuthLog.Errorf("Sponsored Connectivity is disabled by AF")
	} else {
		requestSuppFeat = tempRequestSuppFeat
	}
	respData := models.AppSessionContextRespData{
		ServAuthInfo: models.ServAuthInfo_TP_NOT_KNOWN,
		SuppFeat: pcfSelf.PcfSuppFeats[models.ServiceName_NPCF_POLICYAUTHORIZATION].NegotiateWith(
			requestSuppFeat).String(),
	}

	udrUri := p.getDefaultUdrUri(pcfSelf)
	if udrUri == "" {
		err = fmt.Errorf("can't find any UDR which supported to this PCF")
		return err
	}
	resp, pd, err := p.Consumer().GetBdtData(udrUri, req.BdtRefId)
	bdtData := resp.BdtData
	if err != nil {
		return fmt.Errorf("UDR Get BdtData error[%s]", err.Error())
	} else if pd != nil {
		return fmt.Errorf("UDR Get BdtData fault[%s]", pd.Detail)
	} else if resp == nil {
		return fmt.Errorf("UDR Get BdtData error")
	} else {
		if bdtData.TransPolicy.RecTimeInt.StartTime.After(time.Now()) {
			respData.ServAuthInfo = models.ServAuthInfo_TP_NOT_YET_OCURRED
		} else if bdtData.TransPolicy.RecTimeInt.StopTime.Before(time.Now()) {
			respData.ServAuthInfo = models.ServAuthInfo_TP_EXPIRED
		}
	}
	appSessCtx.AscRespData = &respData
	return nil
}

// provisioning of sponsored connectivity information
func handleSponsoredConnectivityInformation(smPolicy *pcf_context.UeSmPolicyData, relatedPccRuleIds map[string]string,
	aspID, sponID string, sponStatus models.SponsoringStatus, umData *models.UsageMonitoringData,
	updateSMpolicy *bool,
) error {
	if sponStatus == models.SponsoringStatus_DISABLED {
		logger.PolicyAuthLog.Debugf("Sponsored Connectivity is disabled by AF")
		umID := util.GetUmId(aspID, sponID)
		for _, pccRuleID := range relatedPccRuleIds {
			pccRule := smPolicy.PolicyDecision.PccRules[pccRuleID]
			for _, chgID := range pccRule.RefChgData {
				// disables sponsoring a service
				chgData := smPolicy.PolicyDecision.ChgDecs[chgID]
				if chgData.AppSvcProvId == aspID && chgData.SponsorId == sponID {
					chgData.SponsorId = ""
					chgData.AppSvcProvId = ""
					chgData.ReportingLevel = models.ReportingLevel_SER_ID_LEVEL
					smPolicy.PolicyDecision.ChgDecs[chgID] = chgData
					*updateSMpolicy = true
				}
			}
			if pccRule.RefUmData != nil {
				pccRule.RefUmData = nil
				smPolicy.PolicyDecision.PccRules[pccRuleID] = pccRule
			}
			// disable the usage monitoring
			// TODO: As a result, PCF gets the accumulated usage of the sponsored data connectivity
			delete(smPolicy.PolicyDecision.UmDecs, umID)
		}
	} else {
		if umData != nil {
			supp := util.CheckSuppFeat(smPolicy.PolicyDecision.SuppFeat, 5) // UMC support = 5 in 29512
			if !supp {
				err := fmt.Errorf("usage Monitor Control is not supported in SMF")
				return err
			}
		}
		chgIDUsed := false
		chgID := util.GetChgId(smPolicy.ChargingIdGenerator)
		for _, pccRuleID := range relatedPccRuleIds {
			pccRule := smPolicy.PolicyDecision.PccRules[pccRuleID]
			chgData := models.ChargingData{
				ChgId: chgID,
			}
			if pccRule.RefChgData != nil {
				chgID := pccRule.RefChgData[0]
				chgData = *smPolicy.PolicyDecision.ChgDecs[chgID]
			} else {
				chgIDUsed = true
			}
			// TODO: PCF, based on operator policies, shall check whether it is required to
			// validate the sponsored connectivity data.
			// If it is required, it shall perform the authorizations based on sponsored data connectivity profiles.
			// If the authorization fails, the PCF shall send HTTP "403 Forbidden" with the "cause" attribute set to
			// "UNAUTHORIZED_SPONSORED_DATA_CONNECTIVITY"
			pccRule.RefChgData = []string{chgData.ChgId}
			chgData.ReportingLevel = models.ReportingLevel_SPON_CON_LEVEL
			chgData.SponsorId = sponID
			chgData.AppSvcProvId = aspID
			if umData != nil {
				pccRule.RefUmData = []string{umData.UmId}
			}
			util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, nil, nil, &chgData, umData)
			*updateSMpolicy = true
		}
		if chgIDUsed {
			smPolicy.ChargingIdGenerator++
		}
		// TODO: handling UE is roaming in VPLMN case
	}
	return nil
}

// The highest precedence for pcc rule is 23, the precedence lower than 254 is
// all set to 254 and leave 255 for default path
func getAvailablePrecedence(pccRules map[string]*models.PccRule) (maxVaule int32) {
	maxVaule = Precedence_Initial
	for _, rule := range pccRules {
		if rule.Precedence != Precedence_Default_path && rule.Precedence > maxVaule {
			maxVaule = rule.Precedence
		}
	}
	if maxVaule < Precedence_Maximum {
		maxVaule++
	}
	return
}

/*
func getFlowInfos(comp models.MediaComponent) (flows []models.FlowInformation, err error) {
	for _, subComp := range comp.MedSubComps {
		if subComp.EthfDescs != nil {
			return nil, fmt.Errorf("flow Description with Mac Address does not support")
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
			flowDesc, flowDir, err := flowDescFromN5toN7(desc)
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

func getFlowInfos(subComp *models.MediaSubComponent) ([]models.FlowInformation, error) {
	var flows []models.FlowInformation
	if subComp.EthfDescs != nil {
		return nil, fmt.Errorf("flow Description with Mac Address does not support")
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
		flowDesc, flowDir, err := flowDescFromN5toN7(desc)
		if err != nil {
			return nil, err
		}
		flowInfo := models.FlowInformation{
			FlowDescription:   flowDesc,
			FlowDirection:     flowDir,
			PacketFilterUsage: statusUsage[flowDir],
			TosTrafficClass:   subComp.TosTrCl,
		}
		flows = append(flows, flowInfo)
	}
	return flows, nil
}

func flowDescFromN5toN7(n5Flow string) (n7Flow string, direction models.FlowDirection, err error) {
	if strings.HasPrefix(n5Flow, "permit out") {
		n7Flow = n5Flow
		direction = models.FlowDirection_DOWNLINK
	} else if strings.HasPrefix(n5Flow, "permit in") {
		n7Flow = strings.ReplaceAll(n5Flow, "permit in", "permit out")
		direction = models.FlowDirection_UPLINK
	} else if strings.HasPrefix(n5Flow, "permit inout") {
		n7Flow = strings.ReplaceAll(n5Flow, "permit inout", "permit out")
		direction = models.FlowDirection_BIDIRECTIONAL
	} else {
		err = fmt.Errorf("invaild flow Description[%s]", n5Flow)
	}
	return
}

func updateQosInMedComp(qosData models.QosData, comp *models.MediaComponent) (models.QosData,
	bool, bool,
) {
	var dlExist bool
	var ulExist bool
	updatedQosData := qosData
	if comp.FStatus == models.FlowStatus_REMOVED {
		updatedQosData.MaxbrDl = ""
		updatedQosData.MaxbrUl = ""
		return updatedQosData, ulExist, dlExist
	}
	maxBwUl := 0.0
	maxBwDl := 0.0
	minBwUl := 0.0
	minBwDl := 0.0
	for _, subsComp := range comp.MedSubComps {
		for _, flow := range subsComp.FDescs {
			_, dir, err := flowDescFromN5toN7(flow)
			if err != nil {
				logger.PolicyAuthLog.Errorf(
					"flowDescFromN5toN7 error in updateQosInMedComp: %+v", err)
			}
			both := false
			if dir == models.FlowDirection_BIDIRECTIONAL {
				both = true
			}
			if subsComp.FlowUsage != models.FlowUsage_RTCP {
				// not RTCP
				if both || dir == models.FlowDirection_UPLINK {
					ulExist = true
					if comp.MarBwUl != "" {
						bwUl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						maxBwUl += bwUl
					}
					if comp.MirBwUl != "" {
						bwUl, err := pcf_context.ConvertBitRateToKbps(comp.MirBwUl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						minBwUl += bwUl
					}
				}
				if both || dir == models.FlowDirection_DOWNLINK {
					dlExist = true
					if comp.MarBwDl != "" {
						bwDl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						maxBwDl += bwDl
					}
					if comp.MirBwDl != "" {
						bwDl, err := pcf_context.ConvertBitRateToKbps(comp.MirBwDl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						minBwDl += bwDl
					}
				}
			} else {
				if both || dir == models.FlowDirection_UPLINK {
					ulExist = true
					if subsComp.MarBwUl != "" {
						bwUl, err := pcf_context.ConvertBitRateToKbps(subsComp.MarBwUl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						maxBwUl += bwUl
					} else if comp.MarBwUl != "" {
						bwUl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						maxBwUl += (0.05 * bwUl)
					}
				}
				if both || dir == models.FlowDirection_DOWNLINK {
					dlExist = true
					if subsComp.MarBwDl != "" {
						bwDl, err := pcf_context.ConvertBitRateToKbps(subsComp.MarBwDl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
						maxBwDl += bwDl
					} else if comp.MarBwDl != "" {
						bwDl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
						if err != nil {
							logger.PolicyAuthLog.Errorf(
								"pcf_context ConvertBitRateToKbps error in updateQosInMedComp: %+v", err)
						}
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
	return updatedQosData, ulExist, dlExist
}

func updateQosInMedSubComp(qosData *models.QosData, comp *models.MediaComponent,
	subsComp *models.MediaSubComponent,
) (updatedQosData models.QosData, ulExist, dlExist bool) {
	updatedQosData = *qosData
	if comp.FStatus == models.FlowStatus_REMOVED {
		updatedQosData.MaxbrDl = ""
		updatedQosData.MaxbrUl = ""
		return updatedQosData, ulExist, dlExist
	}
	maxBwUl := 0.0
	maxBwDl := 0.0
	minBwUl := 0.0
	minBwDl := 0.0
	for _, flow := range subsComp.FDescs {
		_, dir, err := flowDescFromN5toN7(flow)
		if err != nil {
			logger.PolicyAuthLog.Errorf(
				"flowDescFromN5toN7 error in updateQosInMedSubComp: %+v", err)
		}
		both := false
		if dir == models.FlowDirection_BIDIRECTIONAL {
			both = true
		}
		if subsComp.FlowUsage != models.FlowUsage_RTCP {
			// not RTCP
			if both || dir == models.FlowDirection_UPLINK {
				ulExist = true
				if comp.MarBwUl != "" {
					bwUl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					maxBwUl += bwUl
				}
				if comp.MirBwUl != "" {
					bwUl, err := pcf_context.ConvertBitRateToKbps(comp.MirBwUl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					minBwUl += bwUl
				}
			}
			if both || dir == models.FlowDirection_DOWNLINK {
				dlExist = true
				if comp.MarBwDl != "" {
					bwDl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					maxBwDl += bwDl
				}
				if comp.MirBwDl != "" {
					bwDl, err := pcf_context.ConvertBitRateToKbps(comp.MirBwDl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					minBwDl += bwDl
				}
			}
		} else {
			if both || dir == models.FlowDirection_UPLINK {
				ulExist = true
				if subsComp.MarBwUl != "" {
					bwUl, err := pcf_context.ConvertBitRateToKbps(subsComp.MarBwUl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					maxBwUl += bwUl
				} else if comp.MarBwUl != "" {
					bwUl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwUl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					maxBwUl += (0.05 * bwUl)
				}
			}
			if both || dir == models.FlowDirection_DOWNLINK {
				dlExist = true
				if subsComp.MarBwDl != "" {
					bwDl, err := pcf_context.ConvertBitRateToKbps(subsComp.MarBwDl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
					maxBwDl += bwDl
				} else if comp.MarBwDl != "" {
					bwDl, err := pcf_context.ConvertBitRateToKbps(comp.MarBwDl)
					if err != nil {
						logger.PolicyAuthLog.Errorf(
							"pcf_context ConvertBitRateToKbps error in updateQosInMedSubComp: %+v", err)
					}
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
	return updatedQosData, ulExist, dlExist
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
				pccRuleID := idMaps[key]
				err := smPolicy.RemovePccRule(pccRuleID, nil)
				if err != nil {
					logger.PolicyAuthLog.Warn(err.Error())
				}
				delete(appSession.RelatedPccRuleIds, key)
				delete(appSession.PccRuleIdMapToCompId, pccRuleID)
			}
		} else {
			pccRuleID := idMaps[compN]
			err := smPolicy.RemovePccRule(pccRuleID, nil)
			if err != nil {
				logger.PolicyAuthLog.Warn(err.Error())
			}
			delete(appSession.RelatedPccRuleIds, compN)
			delete(appSession.PccRuleIdMapToCompId, pccRuleID)
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
// 			pccRuleID := idMaps[key]
// 			smPolicy.RemovePccRule(pccRuleID, nil)
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

func extractUmData(umID string, eventSubs map[models.PcfPolicyAuthorizationAfEvent]models.AfNotifMethod,
	threshold *models.UsageThreshold,
) (umData *models.UsageMonitoringData, err error) {
	if _, umExist := eventSubs[models.PcfPolicyAuthorizationAfEvent_USAGE_REPORT]; umExist {
		if threshold == nil {
			return nil, fmt.Errorf("usageThreshold is nil in USAGE REPORT Subscription")
		} else {
			tmp := util.CreateUmData(umID, *threshold)
			umData = &tmp
		}
	}
	return
}

func modifyRemainBitRate(smPolicy *pcf_context.UeSmPolicyData, qosData *models.QosData,
	ulExist, dlExist bool,
) *models.ProblemDetails {
	// if request GBR == 0, qos GBR = MBR
	// if request GBR > remain GBR, qos GBR = remain GBR
	if ulExist {
		if qosData.GbrUl == "" {
			// err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.MaxbrUl)
			if err := pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.MaxbrUl); err != nil {
				qosData.GbrUl = pcf_context.DecreaseRamainBitRateToZero(smPolicy.RemainGbrUL)
			} else {
				qosData.GbrUl = qosData.MaxbrUl
			}
		} else {
			// err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.GbrUl)
			if err := pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.GbrUl); err != nil {
				problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				// sendProblemDetail(httpChannel, err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return &problemDetail
			}
		}
	}
	if dlExist {
		if qosData.GbrDl == "" {
			// err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrDL, qosData.MaxbrDl)
			if err := pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrDL, qosData.MaxbrDl); err != nil {
				qosData.GbrDl = pcf_context.DecreaseRamainBitRateToZero(smPolicy.RemainGbrDL)
			} else {
				qosData.GbrDl = qosData.MaxbrDl
			}
		} else {
			// err = pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrDL, qosData.GbrDl)
			if err := pcf_context.DecreaseRamainBitRate(smPolicy.RemainGbrDL, qosData.GbrDl); err != nil {
				// if Policy failed, revert remain GBR to original GBR
				pcf_context.IncreaseRamainBitRate(smPolicy.RemainGbrUL, qosData.GbrUl)
				problemDetail := util.GetProblemDetail(err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				// sendProblemDetail(httpChannel, err.Error(), util.REQUESTED_SERVICE_NOT_AUTHORIZED)
				return &problemDetail
			}
		}
	}
	return nil
}

func provisioningOfTrafficRoutingInfo(smPolicy *pcf_context.UeSmPolicyData, appID string,
	routeReq *models.AfRoutingRequirement, fStatus models.FlowStatus,
) *models.PccRule {
	var tcData *models.TrafficControlData

	// TODO : handle temporal or spatial validity
	pccRule := util.GetPccRuleByAfAppId(smPolicy.PolicyDecision.PccRules, appID)
	if pccRule != nil {
		// Update TcData
		var tcID string
		if len(pccRule.RefTcData) > 0 {
			// 1 PCC rule only supports 1 TrafficControlData
			// TODO: 1 PCC rule supports multiple TrafficControlData
			// Re-use the original tcID
			tcID = pccRule.RefTcData[0]
			if smPolicy.PolicyDecision.TraffContDecs == nil {
				logger.PolicyAuthLog.Errorf("TraffContDecs is nil, tcID[%s]", tcID)
				tcData = util.CreateTcData(0, tcID, fStatus)
			} else {
				tcData = smPolicy.PolicyDecision.TraffContDecs[tcID]
				if tcData == nil {
					logger.PolicyAuthLog.Errorf("TraffContDecs[%s] not found", tcID)
					tcData = util.CreateTcData(0, tcID, fStatus)
				}
			}
		} else {
			// tcID's number equals to pccRuleID's number
			tcID = strings.ReplaceAll(pccRule.PccRuleId, "PccRule", "Tc")
			tcData = util.CreateTcData(0, tcID, fStatus)
			pccRule.RefTcData = []string{tcID}
		}
		tcData.RouteToLocs = routeReq.RouteToLocs
		tcData.UpPathChgEvent = routeReq.UpPathChgSub
		pccRule.AppReloc = routeReq.AppReloc
		util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, tcData, nil, nil, nil)
		logger.PolicyAuthLog.Infof("Update Traffic Control Data[%s] in PCC rule[%s]",
			tcID, pccRule.PccRuleId)
	} else {
		// Create a Pcc Rule if afappID dose not match any pcc rule
		precedence := getAvailablePrecedence(smPolicy.PolicyDecision.PccRules)
		pccRule = util.CreatePccRule(smPolicy.PccRuleIdGenerator, precedence, nil, appID)
		tcData = util.CreateTcData(smPolicy.PccRuleIdGenerator, "", fStatus)
		tcData.RouteToLocs = routeReq.RouteToLocs
		tcData.UpPathChgEvent = routeReq.UpPathChgSub
		// pccRule.RefTcData = []string{tcData.TcId}
		util.SetPccRuleRelatedData(smPolicy.PolicyDecision, pccRule, tcData, nil, nil, nil)
		smPolicy.PccRuleIdGenerator++
		logger.PolicyAuthLog.Infof("Create PCC rule[%s] with new Traffic Control Data[%s]",
			pccRule.PccRuleId, tcData.TcId)
	}
	return pccRule
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
