package pcf_producer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	// "context"

	"context"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/cydev/zero"
	"github.com/jinzhu/copier"
	// "github.com/antihax/optional"
)

// PostAppSessions - Creates a new Individual Application Session Context resource
func PostAppSessionsContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.AppSessionContext) {
	var appSessionContext models.AppSessionContext = body
	var problemDetails models.ProblemDetails

	logger.PolicyAuthorizationlog.Traceln("AppSessionContext to store: ", appSessionContext)
	pcfSelf := pcf_context.PCF_Self()
	ue := pcfSelf.UePool[appSessionContext.AscReqData.Supi]
	if appSessionContext.AscReqData.UeIpv4 == "" && appSessionContext.AscReqData.UeIpv6 == "" && appSessionContext.AscReqData.UeMac == "" {
		goto END
	}
	if ue == nil {
		if appSessionContext.AscReqData.UeIpv4 != "" {
			ue = pcfSelf.PCFUeFindByIPv4(appSessionContext.AscReqData.UeIpv4)
		} else if appSessionContext.AscReqData.UeIpv6 != "" {
			ue = pcfSelf.PCFUeFindByIPv6(appSessionContext.AscReqData.UeIpv6)

		} else {
			//TODO: Support Mac Address
			logger.PolicyAuthorizationlog.Warnf("MAC Address[%s] have not implemented yet", appSessionContext.AscReqData.UeMac)
			goto END
		}
	}

	if appSessionContext.AscReqData.NotifUri != "" && appSessionContext.AscReqData.SuppFeat != "" {
		if appSessionContext.AscReqData.AspId != "" {
			ue.AspId = appSessionContext.AscReqData.AspId
		}
		// else {
		// 	logger.PolicyAuthorizationlog.Warnln("AppSessionContext Id is nil")
		// }
		if ue.AppSessionIdStore != nil {
			respHeader := make(http.Header)
			respHeader.Set("Location", "app-sessions/"+ue.AppSessionIdStore.AppSessionId)
			pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 303, nil)
			return
		} else {

			appSessionContext.AscRespData = &models.AppSessionContextRespData{
				ServAuthInfo: appSessionContext.AscRespData.ServAuthInfo,
				SuppFeat:     appSessionContext.AscReqData.SuppFeat,
			}
			// Gate Control
			if appSessionContext.AscReqData.MedComponents != nil {
				GateControl(ue, appSessionContext)
			}
			// Initial BDT policy indication
			if appSessionContext.AscReqData.BdtRefId != "" {
				InitialBackgroundDataTransferPolicyIndication(ue, appSessionContext)
			}

			// Initial provisioning of sponsored connectivity information
			if appSessionContext.AscReqData.AspId != "" && appSessionContext.AscReqData.SponId != "" {
				if appSessionContext.AscReqData.AspId != "" && appSessionContext.AscReqData.SponId != "" {
					smfSuppFeat := strings.Split(ue.SmPolicyControlStore.Policy.SuppFeat, "")
					for index := range smfSuppFeat {
						if smfSuppFeat[index] != "13" && index == len(smfSuppFeat) {
							problemDetails.Status = 403
							problemDetails.Cause = "REQUESTED_SERVICE_NOT_AUTHORIZED"
							pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
							return
						}
						if smfSuppFeat[index] == "13" {
							for pccRuleID := range ue.SmPolicyControlStore.Policy.PccRules {
								pccRuleTemp := ue.SmPolicyControlStore.Policy.PccRules[pccRuleID]
								for refUmDataIndex := range pccRuleTemp.RefUmData {
									umId := pccRuleTemp.RefUmData[refUmDataIndex]
									if umId == appSessionContext.AscReqData.SponId {
										if !zero.IsZero(ue.SmPolicyControlStore.Policy.UmDecs[umId]) {
											umData := ue.SmPolicyControlStore.Policy.UmDecs[umId]
											usgThres := appSessionContext.AscReqData.EvSubsc.UsgThres
											if usgThres != nil {
												if umData.VolumeThreshold != usgThres.TotalVolume || umData.VolumeThresholdUplink != usgThres.UplinkVolume || umData.VolumeThresholdDownlink != usgThres.DownlinkVolume || umData.TimeThreshold != usgThres.Duration {
													umData.NextVolThreshold = usgThres.TotalVolume
													umData.NextVolThresholdUplink = usgThres.UplinkVolume
													umData.NextVolThresholdDownlink = usgThres.DownlinkVolume
													umData.NextTimeThreshold = usgThres.Duration
												}
											} else {
												problemDetails.Status = 403
												problemDetails.Cause = "UNAUTHORIZED_SPONSORED_DATA_CONNECTIVITY"
												pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
												return
											}
										}
									}
								}
							}

						}
					}
				}
			}
			// Subscriptions to Service Data Flow QoS notification
			// Subscription to Service Data Flow Deactivation
			// Initial provisioning of traffic routing information
			if appSessionContext.AscReqData.SuppFeat != "" {
				suppFeat := strings.Split(appSessionContext.AscReqData.SuppFeat, "")
				for index := range suppFeat {
					if suppFeat[index] == "1" {
						InitialProvisioningOfTrafficRoutingInformation(ue, appSessionContext)
						break
					}
				}
			}

			for index := range appSessionContext.AscReqData.EvSubsc.Events {
				switch appSessionContext.AscReqData.EvSubsc.Events[index].Event {
				case "PLMN_CHG":
					appSessionContext.EvsNotif.PlmnId.Mnc = ue.SmPolicyControlStore.Context.ServingNetwork.Mnc
					appSessionContext.EvsNotif.PlmnId.Mcc = ue.SmPolicyControlStore.Context.ServingNetwork.Mcc
				case "USAGE_REPORT":
					appSessionContext.EvsNotif.UsgRep.Duration = ue.SmPolicyControlStore.Policy.UmDecs[appSessionContext.AscReqData.SponId].TimeThreshold
					appSessionContext.EvsNotif.UsgRep.TotalVolume = ue.SmPolicyControlStore.Policy.UmDecs[appSessionContext.AscReqData.SponId].VolumeThreshold
					appSessionContext.EvsNotif.UsgRep.DownlinkVolume = ue.SmPolicyControlStore.Policy.UmDecs[appSessionContext.AscReqData.SponId].VolumeThresholdDownlink
					appSessionContext.EvsNotif.UsgRep.UplinkVolume = ue.SmPolicyControlStore.Policy.UmDecs[appSessionContext.AscReqData.SponId].VolumeThresholdUplink
				case "ACCESS_TYPE_CHANGE":
					appSessionContext.EvsNotif.AccessType = ue.SmPolicyControlStore.Context.AccessType
					appSessionContext.EvsNotif.RatType = ue.SmPolicyControlStore.Context.RatType
				// appSessionContext.EvsNotif.AnGwAddr = ue.SmPolicyControlStore.Context.AccNetChId
				case "FAILED_RESOURECES_ALLOCATION":
					afEventNotification := models.AfEventNotification{
						Event: "FAILED_RESOURCES_ALLOCATION",
					}
					appSessionContext.EvsNotif.EvNotifs = append(appSessionContext.EvsNotif.EvNotifs, afEventNotification)
					// appSessionContext.EvsNotif.FailedResourcAllocReports
				case "QOS_NOTIF_CONTROL":

				case "SUCCESSFUL_RESOURCES_ALLOCATION":
					afEventNotification := models.AfEventNotification{
						Event: "SUCCESSFUL_RESOURCES_ALLOCATION",
					}
					appSessionContext.EvsNotif.EvNotifs = append(appSessionContext.EvsNotif.EvNotifs, afEventNotification)
				}
			}
			// Subscription to resources allocation outcome
			// Invocation of Multimedia Priority Services
			// Support of content versioning

			appSession := pcf_context.AppSessionIdStore{
				AppSessionId:      uuid.New().String(),
				AppSessionContext: appSessionContext,
			}
			ue.AppSessionIdStore = &appSession
			// pcf_context.AppSessionContextStore = append(pcf_context.AppSessionContextStore, appSession)
			respHeader := make(http.Header)
			respHeader.Set("Location", ue.AppSessionIdStore.AppSessionId)
			pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 201, appSessionContext)
			return

		}
	}
END:
	problemDetails.Status = 403
	problemDetails.Cause = "REQUESTED_SERVICE_NOT_AUTHORIZED"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
}

// DeleteAppSession - Deletes an existing Individual Application Session Context
func DeleteAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string) {
	AppSessionId := appSessionId
	logger.PolicyAuthorizationlog.Traceln("AppSessionId: ", AppSessionId)
	var eventsNotification *models.EventsNotification
	pcfUeContext := pcf_context.PCF_Self().UePool

	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		if AppSessionId == pcfUeContext[key].AppSessionIdStore.AppSessionId {
			if pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvSubsUri != "" && pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvNotifs != nil {
				eventsNotification = pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif
				pcfUeContext[key].AppSessionIdStore = &pcf_context.AppSessionIdStore{}
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, eventsNotification)
				// CreateSmPolicyNotifyContext(fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId), "terminate")
				return
			}
			if pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvSubsUri == "" && pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvNotifs == nil {
				pcfUeContext[key].AppSessionIdStore = &pcf_context.AppSessionIdStore{}
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, nil)
				// CreateSmPolicyNotifyContext(fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId), "terminate")
				return
			}
		}
	}
}

// GetAppSession - Reads an existing Individual Application Session Context
func GetAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string) {
	var problemDetails models.ProblemDetails
	AppSessionId := appSessionId
	logger.PolicyAuthorizationlog.Traceln("AppSessionId: ", AppSessionId)
	pcfUeContext := pcf_context.PCF_Self().UePool

	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		if AppSessionId == pcfUeContext[key].AppSessionIdStore.AppSessionId {
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, pcfUeContext[key].AppSessionIdStore.AppSessionContext)
			return
		}
	}
	// can not found
	problemDetails.Status = 404
	problemDetails.Cause = "CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)
}

// ModAppSession - Modifies an existing Individual Application Session Context
func ModAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string, body models.AppSessionContextUpdateData) {
	var appSessionContextUpdateData models.AppSessionContextUpdateData = body
	var appSessionContext models.AppSessionContext
	var problemDetails models.ProblemDetails
	AppSessionId := appSessionId
	logger.PolicyAuthorizationlog.Traceln("AppSessionId: ", AppSessionId)
	pcfUeContext := pcf_context.PCF_Self().UePool
	if zero.IsZero(appSessionContextUpdateData) {
		problemDetails.Status = 403
		problemDetails.Cause = "REQUESTED_SERVICE_NOT_AUTHORIZED"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
		return
	}

	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		if AppSessionId == pcfUeContext[key].AppSessionIdStore.AppSessionId {
			if appSessionContextUpdateData.AfAppId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfAppId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfAppId = appSessionContextUpdateData.AfAppId
			}
			if !reflect.DeepEqual(appSessionContextUpdateData.AfRoutReq, pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfRoutReq) {
				if err := copier.Copy(&pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfRoutReq, &appSessionContextUpdateData.AfRoutReq); err != nil {
					logger.PolicyAuthorizationlog.Warnln("AfRoutReq copy error: ", err)
				}
			}
			if appSessionContextUpdateData.AspId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AspId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AspId = appSessionContextUpdateData.AspId
			}
			if appSessionContextUpdateData.BdtRefId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.BdtRefId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.BdtRefId = appSessionContextUpdateData.BdtRefId
			}
			if !reflect.DeepEqual(appSessionContextUpdateData.EvSubsc, pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc) {
				if err := copier.Copy(&pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc, &appSessionContextUpdateData.EvSubsc); err != nil {
					logger.PolicyAuthorizationlog.Warnln("EvSubsc copy error: ", err)
				}
			}
			if !reflect.DeepEqual(pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.MedComponents, appSessionContextUpdateData.MedComponents) {
				if err := copier.Copy(&pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.MedComponents, &appSessionContextUpdateData.MedComponents); err != nil {
					logger.PolicyAuthorizationlog.Warnln("MedComponents copy error: ", err)
				}
			}
			if appSessionContextUpdateData.SponId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.SponId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.SponId = appSessionContextUpdateData.SponId
			}
			if appSessionContextUpdateData.SponStatus != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.SponStatus {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.SponStatus = appSessionContextUpdateData.SponStatus
			}
			if appSessionContextUpdateData.MpsId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.MpsId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.MpsId = appSessionContextUpdateData.MpsId
			}
			appSessionContext = pcfUeContext[key].AppSessionIdStore.AppSessionContext
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, appSessionContext)
			return
		}
	}
	problemDetails.Status = 404
	problemDetails.Cause = "APPLICATION_SESSION_CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)
}

// DeleteEventsSubsc - deletes the Events Subscription subresource
func DeleteEventsSubscContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string) {
	var problemDetails models.ProblemDetails
	AppSessionId := appSessionId
	logger.PolicyAuthorizationlog.Traceln("AppSessionId: ", AppSessionId)
	pcfUeContext := pcf_context.PCF_Self().UePool

	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		if AppSessionId == pcfUeContext[key].AppSessionIdStore.AppSessionId {
			pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc = &models.EventsSubscReqData{}
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, nil)
			return
		}
	}
	problemDetails.Status = 404
	problemDetails.Cause = "APPLICATION_SESSION_CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)
}

// UpdateEventsSubsc - creates or modifies an Events Subscription subresource
func UpdateEventsSubscContext(httpChannel chan pcf_message.HttpResponseMessage, appSessionId string, body models.EventsSubscReqData) {
	var eventsSubscReqData models.EventsSubscReqData = body
	var problemDetails models.ProblemDetails
	AppSessionId := appSessionId
	logger.PolicyAuthorizationlog.Traceln("AppSessionId: ", AppSessionId)
	pcfUeContext := pcf_context.PCF_Self().UePool

	if eventsSubscReqData.Events == nil {
		problemDetails.Status = 403
		problemDetails.Cause = "REQUESTED_SERVICE_NOT_AUTHORIZED"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
		return
	}
	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		if AppSessionId == pcfUeContext[key].AppSessionIdStore.AppSessionId {
			if zero.IsZero(pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc) && !zero.IsZero(eventsSubscReqData.Events) {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc = &eventsSubscReqData
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, eventsSubscReqData)
				return
			}
			if !zero.IsZero(pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc) && !zero.IsZero(eventsSubscReqData.Events) {
				if !reflect.DeepEqual(pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.Events, eventsSubscReqData.Events) {
					pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.Events = eventsSubscReqData.Events
				}
				if pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.NotifUri != eventsSubscReqData.NotifUri {
					pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.NotifUri = eventsSubscReqData.NotifUri
				}
				if pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.UsgThres != eventsSubscReqData.UsgThres {
					pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.UsgThres = eventsSubscReqData.UsgThres
				}
				// CreateSmPolicyNotifyContext(fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId), "update")
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, eventsSubscReqData)
				return
			}
		}
	}
	problemDetails.Status = 404
	problemDetails.Cause = "APPLICATION_SESSION_CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)
}

func Npcf_PolicyAuthorization_Notify(id string, send_type string) {
	var eventsNotification models.EventsNotification
	var terminationInfo models.TerminationInfo
	pcfUeContext := pcf_context.PCF_Self().UePool
	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		idTemp := fmt.Sprint(pcfUeContext[key].SmPolicyControlStore.Context.PduSessionId)
		if id == idTemp {
			if send_type == "update" {
				url := pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.NotifUri
				if err := copier.Copy(eventsNotification, pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif); err != nil {
					logger.PolicyAuthorizationlog.Warnln("Binding error: ", err)
				}
				data, err := json.Marshal(&eventsNotification)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnln("JSON Marshal error: ", err)
				}
				req2 := bytes.NewBuffer([]byte(data))
				Uri := "https://localhost:29514/" + url + "/notify"
				req, err := http.NewRequest("POST", Uri, req2)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnln("Naf update fail error message is: ", err)
				}

				client := &http.Client{}
				_, err = client.Do(req)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnln("Naf update fail error message is : ", err)
				}
				return

			}
			if send_type == "terminate" {
				url := pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.NotifUri
				terminationInfo.ResUri = "https://localhost:29514/" + url + "/terminate"
				terminationInfo.TermCause = "PDU_SESSION_TERMINATION"
				data, err := json.Marshal(&terminationInfo)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnln("JSON Marshal error: ", err)
				}
				req2 := bytes.NewBuffer([]byte(data))
				req, err := http.NewRequest("POST", terminationInfo.ResUri, req2)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnln("Naf delete fail error message is: ", err)
				}
				req.Header.Set("X-Custom-Header", "myvalue")
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				_, err = client.Do(req)
				if err != nil {
					logger.PolicyAuthorizationlog.Warnln("Naf delete fail error message is : ", err)
				}
				return
			}
		}
	}
}

// GateControl - Add FlowStatus into PccRule
func GateControl(pcfUeContext *pcf_context.UeContext, appSessionContext models.AppSessionContext) {
	type Status string
	const (
		Enabled         Status = "ENABLED"
		Disabled        Status = "DISABLED"
		EnabledUplink   Status = "ENABLED-UPLINK"
		EnabledDownlink Status = "ENABLED-DOWNLINK"
	)
	var status = Disabled
	mediaComponents := appSessionContext.AscReqData.MedComponents
	for medCompN := range mediaComponents {
		mediaComponentTemp := mediaComponents[medCompN]
		for fNum := range mediaComponentTemp.MedSubComps {
			medSubCompsTemp := mediaComponentTemp.MedSubComps[fNum]
			if len(medSubCompsTemp.EthfDescs) != 0 {
				for pccRuleID := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
					pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID]
					for flowInfo := range pccRuleTemp.FlowInfos {
						pccRuleTemp.FlowInfos[flowInfo].EthFlowDescription = &medSubCompsTemp.EthfDescs[flowInfo]
					}
				}
			}
			if medSubCompsTemp.FlowUsage == "RTCP" {
				for pccRuleID := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
					pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID]
					if pccRuleTemp.RefTcData == nil {
						pccRuleTemp.RefTcData = []string{uuid.New().String()}
					}
					for refTcDataindex := range pccRuleTemp.RefTcData {
						tcID := pccRuleTemp.RefTcData[refTcDataindex]
						if pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs == nil {
							pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs = make(map[string]models.TrafficControlData)
						}
						traffContTemp := pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs[tcID]
						traffContTemp.FlowStatus = models.FlowStatus_ENABLED
						break
					}
					pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID] = pccRuleTemp
					break
				}
				break
			}
		}
		if mediaComponentTemp.FStatus == models.FlowStatus_ENABLED_DOWNLINK {
			status = EnabledDownlink
		}
		if mediaComponentTemp.FStatus == models.FlowStatus_ENABLED_UPLINK {
			status = EnabledUplink
		}
		if mediaComponentTemp.FStatus == models.FlowStatus_ENABLED {
			status = Enabled
		}
		switch status {
		case "ENABLED":
			for pccRuleID := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
				pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID]
				if pccRuleTemp.RefTcData == nil {
					pccRuleTemp.RefTcData = []string{uuid.New().String()}
				}
				for refTcDataindex := range pccRuleTemp.RefTcData {
					tcID := pccRuleTemp.RefTcData[refTcDataindex]
					if pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs == nil {
						pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs = make(map[string]models.TrafficControlData)
					}
					traffContTemp := pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs[tcID]
					traffContTemp.FlowStatus = models.FlowStatus_ENABLED
				}
				pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID] = pccRuleTemp
			}
		case "ENABLED-DOWNLINK":
			for pccRuleID := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
				pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID]
				if pccRuleTemp.RefTcData == nil {
					pccRuleTemp.RefTcData = []string{uuid.New().String()}
				}
				for refTcDataindex := range pccRuleTemp.RefTcData {
					tcID := pccRuleTemp.RefTcData[refTcDataindex]
					if pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs == nil {
						pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs = make(map[string]models.TrafficControlData)
					}
					traffContTemp := pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs[tcID]
					traffContTemp.FlowStatus = models.FlowStatus_ENABLED_DOWNLINK
				}
				pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID] = pccRuleTemp
			}
		case "ENABLED-UPLINK":
			for pccRuleID := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
				pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID]
				if pccRuleTemp.RefTcData == nil {
					pccRuleTemp.RefTcData = []string{uuid.New().String()}
				}
				for refTcDataindex := range pccRuleTemp.RefTcData {
					tcID := pccRuleTemp.RefTcData[refTcDataindex]
					if pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs == nil {
						pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs = make(map[string]models.TrafficControlData)
					}
					traffContTemp := pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs[tcID]
					traffContTemp.FlowStatus = models.FlowStatus_ENABLED_UPLINK
				}
				pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID] = pccRuleTemp
			}
		case "DISABLED":
			for pccRuleID := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
				pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID]
				if pccRuleTemp.RefTcData == nil {
					pccRuleTemp.RefTcData = []string{uuid.New().String()}
				}
				for refTcDataindex := range pccRuleTemp.RefTcData {
					tcID := pccRuleTemp.RefTcData[refTcDataindex]
					if pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs == nil {
						pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs = make(map[string]models.TrafficControlData)
					}
					traffContTemp := pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs[tcID]
					traffContTemp.FlowStatus = models.FlowStatus_DISABLED
				}
				pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRuleID] = pccRuleTemp
			}
		}
	}
}

func InitialBackgroundDataTransferPolicyIndication(pcfUeContext *pcf_context.UeContext, appSessionContext models.AppSessionContext) {
	if !zero.IsZero(pcfUeContext.BdtPolicyStore) {
		if pcfUeContext.BdtPolicyStore.TransfPolicies[pcfUeContext.BdtPolicyStore.SelTransPolicyId].RecTimeInt.StopTime.Before(time.Now()) {
			appSessionContext.AscRespData.ServAuthInfo = "TP_EXPIRED"
		}
		if pcfUeContext.BdtPolicyStore.TransfPolicies[pcfUeContext.BdtPolicyStore.SelTransPolicyId].RecTimeInt.StartTime.Before(time.Now()) {
			appSessionContext.AscRespData.ServAuthInfo = "TP_NOT_YET_OCCURRED"
		}
	} else {
		client := pcf_util.GetNudrClient("https://localhost:29504")
		bdtdata, _, err := client.DefaultApi.PolicyDataBdtDataGet(context.Background())
		if bdtdata != nil && err == nil {
			for index := range bdtdata {
				if bdtdata[index].BdtRefId == appSessionContext.AscReqData.BdtRefId {
					if bdtdata[index].TransPolicy.RecTimeInt.StartTime.After(time.Now()) {
						appSessionContext.AscRespData.ServAuthInfo = "TP_NOT_YET_OCCURRED"
					}
					if bdtdata[index].TransPolicy.RecTimeInt.StopTime.Before(time.Now()) {
						appSessionContext.AscRespData.ServAuthInfo = "TP_EXPIRED"
					}
					break
				}
			}
		} else {
			appSessionContext.AscRespData.ServAuthInfo = "TP_NOT_KNOWN"
		}
	}
}

func InitialProvisioningOfTrafficRoutingInformation(pcfUeContext *pcf_context.UeContext, appSessionContext models.AppSessionContext) {
	for pccRulesindex := range pcfUeContext.SmPolicyControlStore.Policy.PccRules {
		pccRuleTemp := pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRulesindex]
		if pccRuleTemp.RefTcData == nil {
			pccRuleTemp.RefTcData = []string{uuid.New().String()}
		} else {
			for refTcDataindex := range pccRuleTemp.RefTcData {
				tcId := pccRuleTemp.RefTcData[refTcDataindex]
				if pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs == nil {
					pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs = make(map[string]models.TrafficControlData)
				}
				traffContTemp := pcfUeContext.SmPolicyControlStore.Policy.TraffContDecs[tcId]
				for medCompN := range appSessionContext.AscReqData.MedComponents {
					if appSessionContext.AscReqData.MedComponents[medCompN].AfRoutReq.RouteToLocs != nil {
						medComponent := appSessionContext.AscReqData.MedComponents[medCompN]
						for index := range medComponent.AfRoutReq.RouteToLocs {
							traffContTemp.RouteToLocs = append(traffContTemp.RouteToLocs, medComponent.AfRoutReq.RouteToLocs[index])
						}
					}
				}
				if traffContTemp.RouteToLocs == nil {
					for index := range appSessionContext.AscReqData.AfRoutReq.RouteToLocs {
						traffContTemp.RouteToLocs = append(traffContTemp.RouteToLocs, appSessionContext.AscReqData.AfRoutReq.RouteToLocs[index])
					}
				}
				traffContTemp.RouteToLocs = appSessionContext.AscReqData.AfRoutReq.RouteToLocs
				traffContTemp.UpPathChgEvent = appSessionContext.AscReqData.AfRoutReq.UpPathChgSub
			}
		}
		pccRuleTemp.AppReloc = appSessionContext.AscReqData.AfRoutReq.AppReloc
		pcfUeContext.SmPolicyControlStore.Policy.PccRules[pccRulesindex] = pccRuleTemp
	}
}
