package pcf_producer

import (
	// "context"

	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"net/http"
	"reflect"

	"github.com/cydev/zero"
	uuid "github.com/satori/go.uuid"
	// "github.com/jinzhu/copier"
	// "github.com/antihax/optional"
)

// PostAppSessions - Creates a new Individual Application Session Context resource
func PostAppSessionsContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.AppSessionContext) {
	var appSessionContext models.AppSessionContext = body
	var problemDetails models.ProblemDetails
	logger.PolicyAuthorizationlog.Traceln("AppSessionContext to store: ", appSessionContext)
	pcfUeContext := pcf_context.GetPCFUeContext()

	if appSessionContext.AscReqData.NotifUri != "" && appSessionContext.AscReqData.SuppFeat != "" {
		if appSessionContext.AscReqData.Supi == "" {
			problemDetails.Title = "Not found Supi"
			problemDetails.Status = 400
			problemDetails.Cause = "USER_UNKNOWN"
			pcf_message.SendHttpResponseMessage(httpChannel, nil, 400, problemDetails)
			return
		}
		if appSessionContext.AscReqData.AspId != "" {
			if err := pcf_context.AddAspIdToUe(appSessionContext.AscReqData.Supi, appSessionContext.AscReqData.AspId); err != nil {
				logger.PolicyAuthorizationlog.Warnln("Add AspId to UE fail")
			}
		}
		supi_temp := appSessionContext.AscReqData.Supi
		for key := range pcfUeContext {
			if pcfUeContext[key].Supi == supi_temp && !zero.IsZero(pcfUeContext[key].AppSessionIdStore) {
				respHeader := make(http.Header)
				respHeader.Set("Location", "app-sessions/"+pcfUeContext[key].AppSessionIdStore.AppSessionId)
				pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 303, nil)
				return
			}
			if pcfUeContext[key].Supi == supi_temp && zero.IsZero(pcfUeContext[key].AppSessionIdStore) {
				rid := uuid.NewV4()
				var appSessionContext = models.AppSessionContext{
					AscReqData: &models.AppSessionContextReqData{
						AfAppId:       appSessionContext.AscReqData.AfAppId,
						AfRoutReq:     appSessionContext.AscReqData.AfRoutReq,
						AspId:         appSessionContext.AscReqData.AspId,
						BdtRefId:      appSessionContext.AscReqData.BdtRefId,
						Dnn:           appSessionContext.AscReqData.Dnn,
						EvSubsc:       appSessionContext.AscReqData.EvSubsc,
						MedComponents: appSessionContext.AscReqData.MedComponents,
						IpDomain:      appSessionContext.AscReqData.IpDomain,
						MpsId:         appSessionContext.AscReqData.MpsId,
						NotifUri:      appSessionContext.AscReqData.NotifUri,
						SliceInfo:     appSessionContext.AscReqData.SliceInfo,
						SponId:        appSessionContext.AscReqData.SponId,
						SponStatus:    appSessionContext.AscReqData.SponStatus,
						Supi:          appSessionContext.AscReqData.Supi,
						SuppFeat:      appSessionContext.AscReqData.SuppFeat,
						UeIpv4:        appSessionContext.AscReqData.UeIpv4,
						UeIpv6:        appSessionContext.AscReqData.UeIpv6,
						UeMac:         appSessionContext.AscReqData.UeMac,
					},
					AscRespData: &models.AppSessionContextRespData{
						ServAuthInfo: appSessionContext.AscRespData.ServAuthInfo,
						SuppFeat:     appSessionContext.AscReqData.SuppFeat,
					},
					EvsNotif: &models.EventsNotification{
						AccessType:                appSessionContext.EvsNotif.AccessType,
						AnGwAddr:                  appSessionContext.EvsNotif.AnGwAddr,
						EvSubsUri:                 appSessionContext.EvsNotif.EvSubsUri,
						EvNotifs:                  appSessionContext.EvsNotif.EvNotifs,
						FailedResourcAllocReports: appSessionContext.EvsNotif.FailedResourcAllocReports,
						PlmnId:                    appSessionContext.EvsNotif.PlmnId,
						QncReports:                appSessionContext.EvsNotif.QncReports,
						RatType:                   appSessionContext.EvsNotif.RatType,
						UsgRep:                    appSessionContext.EvsNotif.UsgRep,
					},
				}
				appSession := pcf_context.AppSessionIdStore{
					AppSessionId:      rid.String(),
					AppSessionContext: appSessionContext,
				}
				pcfUeContext[key].AppSessionIdStore = &appSession
				// pcf_context.AppSessionContextStore = append(pcf_context.AppSessionContextStore, appSession)
				respHeader := make(http.Header)
				respHeader.Set("Location", pcfUeContext[key].AppSessionIdStore.AppSessionId)
				pcf_message.SendHttpResponseMessage(httpChannel, respHeader, 201, appSessionContext)
				return
			}
		}
	}
	if appSessionContext.AscReqData.NotifUri == "" || appSessionContext.AscReqData.SuppFeat == "" {
		problemDetails.Status = 403
		problemDetails.Cause = "REQUESTED_SERVICE_NOT_AUTHORIZED"
		pcf_message.SendHttpResponseMessage(httpChannel, nil, 403, problemDetails)
		return
	}
}

// DeleteAppSession - Deletes an existing Individual Application Session Context
func DeleteAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	URI := ReqURI
	logger.PolicyAuthorizationlog.Traceln("URL: ", URI)
	var eventsNotification *models.EventsNotification
	pcfUeContext := pcf_context.GetPCFUeContext()

	for key := range pcfUeContext {
		// if pcfUeContext[key].AppSessionIdStore != nil {
		// 	continue
		// }
		AppSessionId_temp := "/npcf-policyauthorization/v1/app-sessions/" + pcfUeContext[key].AppSessionIdStore.AppSessionId + "/delete"
		if URI == AppSessionId_temp {
			if pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvSubsUri != "" && pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvNotifs != nil {
				eventsNotification = pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif
				pcfUeContext[key].AppSessionIdStore = &pcf_context.AppSessionIdStore{}
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, eventsNotification)
				return
			}
			if pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvSubsUri == "" && pcfUeContext[key].AppSessionIdStore.AppSessionContext.EvsNotif.EvNotifs == nil {
				pcfUeContext[key].AppSessionIdStore = &pcf_context.AppSessionIdStore{}
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 204, nil)
				return
			}
		}
	}
}

// GetAppSession - Reads an existing Individual Application Session Context
func GetAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	var problemDetails models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()

	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		AppSessionId_temp := "/npcf-policyauthorization/v1/app-sessions/" + pcfUeContext[key].AppSessionIdStore.AppSessionId
		if URI == AppSessionId_temp {
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
func ModAppSessionContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.AppSessionContextUpdateData) {
	var appSessionContextUpdateData models.AppSessionContextUpdateData = body
	var appSessionContext models.AppSessionContext
	var problemDetails models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()
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
		AppSessionId_temp := "/npcf-policyauthorization/v1/app-sessions/" + pcfUeContext[key].AppSessionIdStore.AppSessionId
		if URI == AppSessionId_temp {
			if appSessionContextUpdateData.AfAppId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfAppId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfAppId = appSessionContextUpdateData.AfAppId
			}
			if !reflect.DeepEqual(appSessionContextUpdateData.AfRoutReq, pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AfRoutReq) {
				//	if err := deepcopy.Copy(&AppSessionContextStore[index].AppSessionContext.AscReqData.AfRoutReq, &appSessionContextUpdateData.AfRoutReq); err != nil {
				//		fmt.Printf("Binding error: %v", err)
				//	}
			}
			if appSessionContextUpdateData.AspId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AspId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.AspId = appSessionContextUpdateData.AspId
			}
			if appSessionContextUpdateData.BdtRefId != pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.BdtRefId {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.BdtRefId = appSessionContextUpdateData.BdtRefId
			}
			if !reflect.DeepEqual(appSessionContextUpdateData.EvSubsc, pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc) {
				//if err := deepcopy.Copy(&AppSessionContextStore[index].AppSessionContext.AscReqData.EvSubsc, &appSessionContextUpdateData.EvSubsc); err != nil {
				//	fmt.Printf("Binding error: %v", err)
				//}
			}
			if !reflect.DeepEqual(pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.MedComponents, appSessionContextUpdateData.MedComponents) {
				//if err := deepcopy.Copy(&AppSessionContextStore[index].AppSessionContext.AscReqData.MedComponents, &appSessionContextUpdateData.MedComponents); err != nil {
				//	fmt.Printf("Binding error: %v", err)
				//}
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
func DeleteEventsSubscContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string) {
	var problemDetails models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()

	for key := range pcfUeContext {
		if pcfUeContext[key].AppSessionIdStore == nil {
			continue
		}
		AppSessionId_temp := "/npcf-policyauthorization/v1/app-sessions/" + pcfUeContext[key].AppSessionIdStore.AppSessionId + "/events-subscription"
		if URI == AppSessionId_temp {
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
func UpdateEventsSubscContext(httpChannel chan pcf_message.HttpResponseMessage, ReqURI string, body models.EventsSubscReqData) {
	var eventsSubscReqData models.EventsSubscReqData = body
	var problemDetails models.ProblemDetails
	URI := ReqURI
	pcfUeContext := pcf_context.GetPCFUeContext()

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
		AppSessionId_temp := "/npcf-policyauthorization/v1/app-sessions/" + pcfUeContext[key].AppSessionIdStore.AppSessionId + "/events-subscription"
		if URI == AppSessionId_temp {
			if pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc == nil && eventsSubscReqData.Events != nil {
				pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc = &eventsSubscReqData
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 201, eventsSubscReqData)
				return
			}
			if pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc != nil && eventsSubscReqData.Events != nil {
				if !reflect.DeepEqual(pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.Events, eventsSubscReqData.Events) {
					pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.Events = eventsSubscReqData.Events
				}
				if pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.NotifUri != eventsSubscReqData.NotifUri {
					pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.NotifUri = eventsSubscReqData.NotifUri
				}
				if pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.UsgThres != eventsSubscReqData.UsgThres {
					pcfUeContext[key].AppSessionIdStore.AppSessionContext.AscReqData.EvSubsc.UsgThres = eventsSubscReqData.UsgThres
				}
				pcf_message.SendHttpResponseMessage(httpChannel, nil, 200, eventsSubscReqData)
				return
			}
		}
	}
	problemDetails.Status = 404
	problemDetails.Cause = "APPLICATION_SESSION_CONTEXT_NOT_FOUND"
	pcf_message.SendHttpResponseMessage(httpChannel, nil, 404, problemDetails)
}

// func PCFEventNotification(notifUri string, send_type string) {
// 	var eventsNotification models.EventsNotification
// 	for index := range AppSessionContextStore {
// 		var UriTemp = AppSessionContextStore[index].AppSessionContext.AscReqData.EvSubsc.NotifUri
// 		if notifUri == UriTemp {
// 			if send_type == "update" {
// 				//if err := deepcopy.Copy(eventsNotification, AppSessionContextStore[index].AppSessionContext.EvsNotif); err != nil {
// 				//	fmt.Printf("Binding error: %v", err)
// 				//}
// 				data, err := json.Marshal(&eventsNotification)
// 				if err != nil {
// 					logger.PolicyAuthorizationlog.Warnln("JSON Marshal error: ", err)
// 				}
// 				req2 := bytes.NewBuffer([]byte(data))
// 				var Uri = "https://localhost:8081/" + notifUri + "/notify"
// 				req, err := http.NewRequest("POST", Uri, req2)
// 				if err != nil {
// 					logger.PolicyAuthorizationlog.Warnln("Naf update fail error message is: ", err)
// 				}
// 				req.Header.Set("X-Custom-Header", "myvalue")
// 				req.Header.Set("Content-Type", "application/json")

// 				client := &http.Client{}
// 				_, err = client.Do(req)
// 				if err != nil {
// 					logger.PolicyAuthorizationlog.Warnln("Naf update fail error message is : ", err)
// 				}
// 				return

// 			}
// 		}
// 	}
// }

// func PCFEventTermination(notifUri string, send_type string) {
// 	var terminationInfo models.TerminationInfo
// 	for index := range AppSessionContextStore {
// 		var UriTemp = AppSessionContextStore[index].AppSessionContext.AscReqData.NotifUri
// 		if notifUri == UriTemp {
// 			if send_type == "terminate" {
// 				terminationInfo.ResUri = "https://localhost:8081/" + notifUri + "/terminate"
// 				terminationInfo.TermCause = "PDU_SESSION_TERMINATION"
// 				data, err := json.Marshal(&terminationInfo)
// 				if err != nil {
// 					logger.PolicyAuthorizationlog.Warnln("JSON Marshal error: ", err)
// 				}
// 				req2 := bytes.NewBuffer([]byte(data))
// 				req, err := http.NewRequest("POST", terminationInfo.ResUri, req2)
// 				if err != nil {
// 					logger.PolicyAuthorizationlog.Warnln("Naf delete fail error message is: ", err)
// 				}
// 				req.Header.Set("X-Custom-Header", "myvalue")
// 				req.Header.Set("Content-Type", "application/json")

// 				client := &http.Client{}
// 				_, err = client.Do(req)
// 				if err != nil {
// 					logger.PolicyAuthorizationlog.Warnln("Naf delete fail error message is : ", err)
// 				}
// 				return
// 			}
// 		}
// 	}
// }
