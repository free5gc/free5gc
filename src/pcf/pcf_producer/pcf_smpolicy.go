package pcf_producer

import (
	"context"
	"fmt"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"
	"strconv"
	"strings"

	"github.com/antihax/optional"
)

// SmPoliciesPost -
func CreateSmPolicy(httpChannel chan pcf_message.HttpResponseMessage, request models.SmPolicyContextData) {
	var err error

	logger.SMpolicylog.Tracef("Handle Create SM Policy Request")
	pcfSelf := pcf_context.PCF_Self()
	ue := pcfSelf.UePool[request.Supi]
	if ue == nil {
		rsp := pcf_util.GetProblemDetail("Supi is not supported in PCF", pcf_util.USER_UNKNOWN)
		logger.SMpolicylog.Warnf("Supi[%s] is not supported in PCF", request.Supi)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
	}
	udrUri := getUdrUri(ue)
	if udrUri == "" {
		rsp := pcf_util.GetProblemDetail("Can't find corresponding UDR with UE", pcf_util.USER_UNKNOWN)
		logger.SMpolicylog.Warnf("Can't find corresponding UDR with UE[%s]", ue.Supi)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	var smData models.SmPolicyData
	smPolicyId := fmt.Sprintf("%s-%d", ue.Supi, request.PduSessionId)
	smPolicyData := ue.SmPolicyData[smPolicyId]
	if smPolicyData == nil || smPolicyData.SmPolicyData == nil {
		client := pcf_util.GetNudrClient(udrUri)
		param := Nudr_DataRepository.PolicyDataUesUeIdSmDataGetParamOpts{
			Snssai: optional.NewInterface(pcf_util.MarshToJsonString(*request.SliceInfo)),
			Dnn:    optional.NewString(request.Dnn),
		}
		var response *http.Response
		smData, response, err = client.DefaultApi.PolicyDataUesUeIdSmDataGet(context.Background(), ue.Supi, &param)
		if err != nil || response == nil || response.StatusCode != http.StatusOK {
			rsp := pcf_util.GetProblemDetail("Can't find UE SM Policy Data in UDR", pcf_util.USER_UNKNOWN)
			logger.SMpolicylog.Warnf("Can't find UE[%s] SM Policy Data in UDR", ue.Supi)
			pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
			return
			// logger.SMpolicylog.Warnln("Nudr Query failed [%s]", err.Error())
		}
		//TODO: subscribe to UDR
	} else {
		smData = *smPolicyData.SmPolicyData
	}
	amPolicy := ue.FindAMPolicy(request.AccessType, request.ServingNetwork)
	if amPolicy == nil {
		rsp := pcf_util.GetProblemDetail("Can't find corresponding AM Policy", pcf_util.POLICY_CONTEXT_DENIED)
		logger.SMpolicylog.Warnf("Can't find corresponding AM Policy")
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	// TODO: check service restrict
	if ue.Gpsi == "" {
		ue.Gpsi = request.Gpsi
	}
	if ue.Pei == "" {
		ue.Pei = request.Pei
	}
	if smPolicyData != nil {
		delete(ue.SmPolicyData, smPolicyId)
	}
	smPolicyData = ue.NewUeSmPolicyData(smPolicyId, request, &smData)
	// Policy Decision
	decision := models.SmPolicyDecision{
		SessRules: make(map[string]models.SessionRule),
	}
	SessRuleId := fmt.Sprintf("SessRuleId-%d", request.PduSessionId)
	sessRule := models.SessionRule{
		AuthSessAmbr: request.SubsSessAmbr,
		SessRuleId:   SessRuleId,
		// RefUmData
		// RefCondData
	}
	defQos := request.SubsDefQos
	if defQos != nil {
		sessRule.AuthDefQos = &models.AuthorizedDefaultQos{
			Var5qi:        defQos.Var5qi,
			Arp:           defQos.Arp,
			PriorityLevel: defQos.PriorityLevel,
			// AverWindow
			// MaxDataBurstVol
		}
	}
	decision.SessRules[SessRuleId] = sessRule
	// TODO: See how UDR used
	dnnData := pcf_util.GetSMPolicyDnnData(smData, request.SliceInfo, request.Dnn)
	if dnnData != nil {
		decision.Online = dnnData.Online
		decision.Offline = dnnData.Offline
		decision.Ipv4Index = dnnData.Ipv4Index
		decision.Ipv6Index = dnnData.Ipv6Index
		// Set Aggregate GBR if exist
		if dnnData.GbrDl != "" {
			gbrDL, err := pcf_context.ConvertBitRateToKbps(dnnData.GbrDl)
			if err != nil {
				logger.SMpolicylog.Warnf(err.Error())
			} else {
				smPolicyData.RemainGbrDL = &gbrDL
				logger.SMpolicylog.Tracef("SM Policy Dnn[%s] Data Aggregate DL GBR[%.2f Kbps]", request.Dnn, gbrDL)
			}
		}
		if dnnData.GbrUl != "" {
			gbrUL, err := pcf_context.ConvertBitRateToKbps(dnnData.GbrUl)
			if err != nil {
				logger.SMpolicylog.Warnf(err.Error())
			} else {
				smPolicyData.RemainGbrUL = &gbrUL
				logger.SMpolicylog.Tracef("SM Policy Dnn[%s] Data Aggregate UL GBR[%.2f Kbps]", request.Dnn, gbrUL)

			}
		}
	} else {
		logger.SMpolicylog.Warnf("Policy Subscription Info: SMPolicyDnnData is null for dnn[%s] in UE[%s]", request.Dnn, ue.Supi)
		decision.Online = request.Online
		decision.Offline = request.Offline
	}
	decision.SuppFeat = pcf_util.GetNegotiateSuppFeat(request.SuppFeat, pcfSelf.PcfSuppFeats[models.ServiceName_NPCF_SMPOLICYCONTROL])
	decision.QosFlowUsage = request.QosFlowUsage
	// TODO: Trigger about UMC, ADC, NetLoc,...
	decision.PolicyCtrlReqTriggers = pcf_util.PolicyControlReqTrigToArray(0x40780f)
	smPolicyData.PolicyDecision = &decision
	// TODO: PCC rule, PraInfo ...
	locationHeader := pcf_util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyId)
	headers := http.Header{
		"Location": {locationHeader},
	}
	logger.SMpolicylog.Tracef("SMPolicy PduSessionId[%d] Create", request.PduSessionId)
	pcf_message.SendHttpResponseMessage(httpChannel, headers, 201, decision)
}

// SmPoliciesSmPolicyIdDeletePost -
func DeleteSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, smPolicyId string) {
	logger.AMpolicylog.Traceln("Handle SM Policy Delete")

	ue := pcf_context.PCF_Self().PCFUeFindByPolicyId(smPolicyId)
	if ue == nil || ue.SmPolicyData[smPolicyId] == nil {
		rsp := pcf_util.GetProblemDetail("smPolicyId not found in PCF", pcf_util.CONTEXT_NOT_FOUND)
		logger.SMpolicylog.Warnf(rsp.Detail)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}

	pcfSelf := pcf_context.PCF_Self()
	smPolicy := ue.SmPolicyData[smPolicyId]

	// Unsubscrice UDR
	delete(ue.SmPolicyData, smPolicyId)
	logger.SMpolicylog.Tracef("SMPolicy SmPolicyId[%s] DELETE", smPolicyId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)

	// Release related App Session
	terminationInfo := models.TerminationInfo{
		TermCause: models.TerminationCause_PDU_SESSION_TERMINATION,
	}
	for appSessionId := range smPolicy.AppSessions {
		if appSession, exist := pcfSelf.AppSessionPool[appSessionId]; exist {
			SendAppSessionTermination(appSession, terminationInfo)
			delete(pcfSelf.AppSessionPool, appSessionId)
			logger.SMpolicylog.Tracef("SMPolicy[%s] DELETE Related AppSession[%s]", smPolicyId, appSessionId)

		}
	}

}

// SmPoliciesSmPolicyIdGet -
func GetSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, smPolicyId string) {
	logger.SMpolicylog.Traceln("Handle GET SM Policy Request")

	ue := pcf_context.PCF_Self().PCFUeFindByPolicyId(smPolicyId)
	if ue == nil || ue.SmPolicyData[smPolicyId] == nil {
		rsp := pcf_util.GetProblemDetail("smPolicyId not found in PCF", pcf_util.CONTEXT_NOT_FOUND)
		logger.SMpolicylog.Warnf(rsp.Detail)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	smPolicyData := ue.SmPolicyData[smPolicyId]
	rsp := models.SmPolicyControl{
		Policy:  smPolicyData.PolicyDecision,
		Context: smPolicyData.PolicyContext,
	}
	logger.SMpolicylog.Tracef("SMPolicy SmPolicyId[%s] GET", smPolicyId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, rsp)
}

// SmPoliciesSmPolicyIdUpdatePost -
func UpdateSmPolicyContext(httpChannel chan pcf_message.HttpResponseMessage, smPolicyId string, request models.SmPolicyUpdateContextData) {

	logger.SMpolicylog.Traceln("Handle SM Policy Update")

	ue := pcf_context.PCF_Self().PCFUeFindByPolicyId(smPolicyId)
	if ue == nil || ue.SmPolicyData[smPolicyId] == nil {
		rsp := pcf_util.GetProblemDetail("smPolicyId not found in PCF", pcf_util.CONTEXT_NOT_FOUND)
		logger.SMpolicylog.Warnf(rsp.Detail)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	smPolicy := ue.SmPolicyData[smPolicyId]
	smPolicyDecision := smPolicy.PolicyDecision
	smPolicyContext := smPolicy.PolicyContext
	errCause := ""

	// For App Session Notification
	afEventsNotification := models.EventsNotification{}
	for _, trigger := range request.RepPolicyCtrlReqTriggers {
		switch trigger {
		case models.PolicyControlRequestTrigger_PLMN_CH: // PLMN Change
			if request.ServingNetwork == nil {
				errCause = "Serving Network is nil in Trigger PLMN_CH"
				break
			}
			smPolicyContext.ServingNetwork = request.ServingNetwork
			afEventsNotification.PlmnId = &models.PlmnId{
				Mcc: request.ServingNetwork.Mcc,
				Mnc: request.ServingNetwork.Mnc,
			}
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_PLMN_CHG,
			}
			afEventsNotification.EvNotifs = append(afEventsNotification.EvNotifs, afNotif)

			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_RES_MO_RE: // UE intiate resource modification to SMF (subsclause 4.2.4.17 in TS29512)
			req := request.UeInitResReq
			if req == nil {
				errCause = "UeInitResReq is nil in Trigger RES_MO_RE"
				break
			}
			switch req.RuleOp {
			case models.RuleOperation_CREATE_PCC_RULE:
				if req.ReqQos == nil || len(req.PackFiltInfo) < 1 {
					errCause = "Parameter Erroneous/Missing in Create Pcc Rule"
					break
				}
				// TODO: Packet Filters are covered by outstanding pcc rule
				id := smPolicy.PccRuleIdGenarator
				infos := pcf_util.ConvertPacketInfoToFlowInformation(req.PackFiltInfo)
				// Set PackFiltId
				for i := range infos {
					infos[i].PackFiltId = pcf_util.GetPackFiltId(smPolicy.PackFiltIdGenarator)
					smPolicy.PackFiltIdGenarator++
				}
				pccRule := pcf_util.CreatePccRule(id, req.Precedence, infos, false)
				// Add Traffic control Data
				tcData := pcf_util.CreateTcData(id, "")
				// TODO: ARP use real Data
				qosData := pcf_util.CreateQosData(id, int32(req.ReqQos.Var5qi), 15)
				// TODO: Set MBR
				var err error
				// Set GBR
				qosData.GbrDl, qosData.GbrUl, err = smPolicy.DecreaseRemainGBR(req.ReqQos)
				if err != nil {
					rsp := pcf_util.GetProblemDetail(err.Error(), pcf_util.ERROR_TRAFFIC_MAPPING_INFO_REJECTED)
					logger.SMpolicylog.Warnf(rsp.Detail)
					pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
					return
				}
				if qosData.GbrDl != "" {
					logger.SMpolicylog.Tracef("SM Policy Dnn[%s] Data Aggregate decrease %s and then DL GBR remain[%.2f Kbps]", smPolicyContext.Dnn, qosData.GbrDl, *smPolicy.RemainGbrDL)
				}
				if qosData.GbrUl != "" {
					logger.SMpolicylog.Tracef("SM Policy Dnn[%s] Data Aggregate decrease %s and then UL GBR remain[%.2f Kbps]", smPolicyContext.Dnn, qosData.GbrUl, *smPolicy.RemainGbrUL)
				}
				pcf_util.SetPccRuleRelatedData(smPolicyDecision, pccRule, &tcData, &qosData, nil, nil)
				// link Packet filters to PccRule
				for _, info := range infos {
					smPolicy.PackFiltMapToPccRuleId[info.PackFiltId] = pccRule.PccRuleId
				}
				smPolicy.PccRuleIdGenarator++
			case models.RuleOperation_DELETE_PCC_RULE:
				if req.PccRuleId == "" {
					errCause = "Parameter Erroneous/Missing in Create Pcc Rule"
					break
				}
				err := smPolicy.RemovePccRule(req.PccRuleId)
				if err != nil {
					errCause = err.Error()
				}
			case models.RuleOperation_MODIFY_PCC_RULE_AND_ADD_PACKET_FILTERS,
				models.RuleOperation_MODIFY_PCC_RULE_AND_REPLACE_PACKET_FILTERS,
				models.RuleOperation_MODIFY_PCC_RULE_AND_DELETE_PACKET_FILTERS,
				models.RuleOperation_MODIFY_PCC_RULE_WITHOUT_MODIFY_PACKET_FILTERS:
				if req.PccRuleId == "" || (req.RuleOp != models.RuleOperation_MODIFY_PCC_RULE_WITHOUT_MODIFY_PACKET_FILTERS && len(req.PackFiltInfo) < 1) {
					errCause = "Parameter Erroneous/Missing in Modify Pcc Rule"
					break
				}
				if rule, exist := smPolicyDecision.PccRules[req.PccRuleId]; exist {
					// Modify Qos if included
					rule.Precedence = req.Precedence
					if req.ReqQos != nil && len(rule.RefQosData) != 0 {
						qosId := rule.RefQosData[0]
						if qosData, exist := smPolicyDecision.QosDecs[qosId]; exist {
							origUl, origDl := smPolicy.IncreaseRemainGBR(qosId)
							gbrDl, gbrUl, err := smPolicy.DecreaseRemainGBR(req.ReqQos)
							if err != nil {
								smPolicy.RemainGbrDL = origDl
								smPolicy.RemainGbrUL = origUl
								rsp := pcf_util.GetProblemDetail(err.Error(), pcf_util.ERROR_TRAFFIC_MAPPING_INFO_REJECTED)
								logger.SMpolicylog.Warnf(rsp.Detail)
								pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
								return
							}
							qosData.Var5qi = req.ReqQos.Var5qi
							qosData.GbrDl = gbrDl
							qosData.GbrUl = gbrUl
							if qosData.GbrDl != "" {
								logger.SMpolicylog.Tracef("SM Policy Dnn[%s] Data Aggregate decrease %s and then DL GBR remain[%.2f Kbps]", smPolicyContext.Dnn, qosData.GbrDl, *smPolicy.RemainGbrDL)
							}
							if qosData.GbrUl != "" {
								logger.SMpolicylog.Tracef("SM Policy Dnn[%s] Data Aggregate decrease %s and then UL GBR remain[%.2f Kbps]", smPolicyContext.Dnn, qosData.GbrUl, *smPolicy.RemainGbrUL)
							}
							smPolicyDecision.QosDecs[qosId] = qosData
						} else {
							errCause = "Parameter Erroneous/Missing in Modify Pcc Rule"
							break
						}
					}
					infos := pcf_util.ConvertPacketInfoToFlowInformation(req.PackFiltInfo)
					switch req.RuleOp {
					case models.RuleOperation_MODIFY_PCC_RULE_AND_ADD_PACKET_FILTERS:
						// Set PackFiltId
						for i := range infos {
							infos[i].PackFiltId = pcf_util.GetPackFiltId(smPolicy.PackFiltIdGenarator)
							smPolicy.PackFiltMapToPccRuleId[infos[i].PackFiltId] = req.PccRuleId
							smPolicy.PackFiltIdGenarator++
						}
						rule.FlowInfos = append(rule.FlowInfos, infos...)
					case models.RuleOperation_MODIFY_PCC_RULE_AND_REPLACE_PACKET_FILTERS:
						// Replace all Packet Filters
						for _, info := range rule.FlowInfos {
							delete(smPolicy.PackFiltMapToPccRuleId, info.PackFiltId)
						}
						// Set PackFiltId
						for i := range infos {
							infos[i].PackFiltId = pcf_util.GetPackFiltId(smPolicy.PackFiltIdGenarator)
							smPolicy.PackFiltMapToPccRuleId[infos[i].PackFiltId] = req.PccRuleId
							smPolicy.PackFiltIdGenarator++
						}
						rule.FlowInfos = infos
					case models.RuleOperation_MODIFY_PCC_RULE_AND_DELETE_PACKET_FILTERS:
						removeId := make(map[string]bool)
						for _, info := range infos {
							delete(smPolicy.PackFiltMapToPccRuleId, info.PackFiltId)
							removeId[info.PackFiltId] = true
						}
						result := []models.FlowInformation{}
						for _, info := range rule.FlowInfos {
							if _, exist := removeId[info.PackFiltId]; !exist {
								result = append(result, info)
							}
						}
						rule.FlowInfos = result
					}
					smPolicyDecision.PccRules[req.PccRuleId] = rule
				} else {
					errCause = fmt.Sprintf("Can't find the pccRuleId[%s] in Session[%d]", req.PccRuleId, smPolicyContext.PduSessionId)
				}

			}

		case models.PolicyControlRequestTrigger_AC_TY_CH: // UE Access Type Change (subsclause 4.2.4.8 in TS29512)
			if request.AccessType == "" {
				errCause = "Access Type is empty in Trigger AC_TY_CH"
				break
			}
			if request.AccessType == models.AccessType__3_GPP_ACCESS && smPolicyContext.Var3gppPsDataOffStatus {
				// TODO: Handle Data off Status
				// Block Session Service except for Exempt Serice which is described in TS22011, TS 23221
			}
			smPolicyContext.AccessType = request.AccessType
			afEventsNotification.AccessType = request.AccessType
			if request.RatType != "" {
				smPolicyContext.RatType = request.RatType
				afEventsNotification.RatType = request.RatType
			}
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_ACCESS_TYPE_CHANGE,
			}
			afEventsNotification.EvNotifs = append(afEventsNotification.EvNotifs, afNotif)
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_UE_IP_CH: // SMF notice PCF "ipv4Address" & ipv6AddressPrefix (always)
			// TODO: Decide new Session Rule / Pcc rule
			if request.RelIpv4Address == smPolicyContext.Ipv4Address {
				smPolicyContext.Ipv4Address = ""
			}
			if request.RelIpv6AddressPrefix == smPolicyContext.Ipv6AddressPrefix {
				smPolicyContext.Ipv6AddressPrefix = ""
			}
			if request.Ipv4Address != "" {
				smPolicyContext.Ipv4Address = request.Ipv4Address
			}
			if request.Ipv6AddressPrefix != "" {
				smPolicyContext.Ipv6AddressPrefix = request.Ipv6AddressPrefix
			}
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_UE_MAC_CH: // SMF notice PCF when SMF detect new UE MAC
		case models.PolicyControlRequestTrigger_AN_CH_COR: // Access Network Charging Correlation Info (subsclause 4.2.6.5.1, 4.2.4.13 in TS29512)
			// request.AccNetChIds
		case models.PolicyControlRequestTrigger_US_RE: // UMC (subsclause 4.2.4.10, 5.8 in TS29512)
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_USAGE_REPORT,
			}
			afEventsNotification.EvNotifs = append(afEventsNotification.EvNotifs, afNotif)
		case models.PolicyControlRequestTrigger_APP_STA: // ADC (subsclause 4.2.4.6, 5.8 in TS29512)
			// request.AppDetectionInfos
		case models.PolicyControlRequestTrigger_APP_STO: // ADC (subsclause 4.2.4.6, 5.8 in TS29512)
			// request.AppDetectionInfos
		case models.PolicyControlRequestTrigger_AN_INFO: // NetLoc (subsclause 4.2.4.9, 5.8 in TS29512)
		case models.PolicyControlRequestTrigger_CM_SES_FAIL: // Credit Management Session Failure
			// request.CreditManageStatus
		case models.PolicyControlRequestTrigger_PS_DA_OFF: // 3GPP PS Data Off status changed (subsclause 4.2.4.8, 5.8 in TS29512) (always)
			if smPolicyContext.Var3gppPsDataOffStatus != request.Var3gppPsDataOffStatus {
				// TODO: Handle Data off Status
				if request.Var3gppPsDataOffStatus {
					// Block Session Service except for Exempt Serice which is described in TS22011, TS 23221
				} else {
					// UnBlock Session Sevice
				}
				smPolicyContext.Var3gppPsDataOffStatus = request.Var3gppPsDataOffStatus
			}
		case models.PolicyControlRequestTrigger_DEF_QOS_CH: // Default QoS Change (subsclause 4.2.4.5 in TS29512) (always)
			if request.SubsDefQos == nil {
				errCause = "SubsDefQos  is nil in Trigger DEF_QOS_CH"
				break
			}
			smPolicyContext.SubsDefQos = request.SubsDefQos
			sessRuleId := fmt.Sprintf("SessRuleId-%d", smPolicyContext.PduSessionId)
			if smPolicyDecision.SessRules[sessRuleId].AuthDefQos == nil {
				tmp := smPolicyDecision.SessRules[sessRuleId]
				tmp.AuthDefQos = new(models.AuthorizedDefaultQos)
				smPolicyDecision.SessRules[sessRuleId] = tmp
			}
			authQos := smPolicyDecision.SessRules[sessRuleId].AuthDefQos
			authQos.Var5qi = request.SubsDefQos.Var5qi
			authQos.Arp = request.SubsDefQos.Arp
			authQos.PriorityLevel = request.SubsDefQos.PriorityLevel
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_SE_AMBR_CH: // Session Ambr Change (subsclause 4.2.4.4 in TS29512) (always)
			if request.SubsSessAmbr == nil {
				errCause = "SubsSessAmbr  is nil in Trigger SE_AMBR_CH"
				break
			}
			smPolicyContext.SubsSessAmbr = request.SubsSessAmbr
			sessRuleId := fmt.Sprintf("SessRuleId-%d", smPolicyContext.PduSessionId)
			if smPolicyDecision.SessRules[sessRuleId].AuthSessAmbr == nil {
				tmp := smPolicyDecision.SessRules[sessRuleId]
				tmp.AuthSessAmbr = new(models.Ambr)
				smPolicyDecision.SessRules[sessRuleId] = tmp
			}
			*smPolicyDecision.SessRules[sessRuleId].AuthSessAmbr = *request.SubsSessAmbr
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_QOS_NOTIF: // SMF notify PCF when receiving from RAN that QoS can/can't be guaranteed (subsclause 4.2.4.20 in TS29512) (always)
			// request.QncReports
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_QOS_NOTIF,
			}
			afEventsNotification.EvNotifs = append(afEventsNotification.EvNotifs, afNotif)
			afEventsNotification.QncReports = request.QncReports
		case models.PolicyControlRequestTrigger_NO_CREDIT: // Out of Credit
		case models.PolicyControlRequestTrigger_PRA_CH: // Presence Reporting (subsclause 4.2.6.5.6, 4.2.4.16, 5.8 in TS29512)
			// request.RepPraInfos
		case models.PolicyControlRequestTrigger_SAREA_CH: // Change Of Service Area
			if request.UserLocationInfo == nil {
				errCause = "UserLocationInfo  is nil in Trigger SAREA_CH"
				break
			}
			smPolicyContext.UserLocationInfo = request.UserLocationInfo
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_SCNN_CH: // Change of Serving Network Function
			if request.ServNfId == nil {
				errCause = "ServNfId  is nil in Trigger SCNN_CH"
				break
			}
			smPolicyContext.ServNfId = request.ServNfId
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_RE_TIMEOUT: // Revalidation TimeOut (subsclause 4.2.4.13 in TS29512)
			// formatTimeStr := time.Now()
			// formatTimeStr = formatTimeStr.Add(time.Second * 60)
			// formatTimeStrAdd := formatTimeStr.Format(pcf_context.GetTimeformat())
			// formatTime, err := time.Parse(pcf_context.GetTimeformat(), formatTimeStrAdd)
			// if err == nil {
			// 	smPolicyDecision.RevalidationTime = &formatTime
			// }
		case models.PolicyControlRequestTrigger_RES_RELEASE: // Outcome of request Pcc rule removal (subsclause 4.2.6.5.2, 5.8 in TS29512)
		case models.PolicyControlRequestTrigger_SUCC_RES_ALLO: // Successful resource allocation (subsclause 4.2.6.5.5, 4.2.4.14 in TS29512)
			afNotif := models.AfEventNotification{
				Event: models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION,
			}
			afEventsNotification.EvNotifs = append(afEventsNotification.EvNotifs, afNotif)
		case models.PolicyControlRequestTrigger_RAT_TY_CH: // Change of RatType
			if request.RatType == "" {
				errCause = "RatType is empty in Trigger RAT_TY_CH"
				break
			}
			smPolicyContext.RatType = request.RatType
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_REF_QOS_IND_CH: // Change of reflective Qos Indication from UE
			smPolicyContext.RefQosIndication = request.RefQosIndication
			// TODO: modify Decision about RefQos in Pcc rule
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		case models.PolicyControlRequestTrigger_NUM_OF_PACKET_FILTER: // Interworking Only (always)
		case models.PolicyControlRequestTrigger_UE_STATUS_RESUME: // UE State Resume
			// TODO
		case models.PolicyControlRequestTrigger_UE_TZ_CH: // UE TimeZome Change
			if request.UeTimeZone == "" {
				errCause = "Ue TimeZone is empty in Trigger UE_TZ_CH"
				break
			}
			smPolicyContext.UeTimeZone = request.UeTimeZone
			logger.SMpolicylog.Tracef("SM Policy Update(%s) Successfully", trigger)
		}
	}

	var successRules, failRules []models.RuleReport
	for _, rule := range request.RuleReports {
		if rule.RuleStatus == models.RuleStatus_ACTIVE {
			successRules = append(successRules, rule)
		} else {
			failRules = append(failRules, rule)
			// release fail pccRules in SmPolicy
			for _, pccRuleId := range rule.PccRuleIds {
				if err := smPolicy.RemovePccRule(pccRuleId); err != nil {
					logger.SMpolicylog.Warnf("SM Policy Notification about failed installing PccRule[%s]", err.Error())
				}
			}
		}
	}
	if len(failRules) > 0 {
		afNotif := models.AfEventNotification{
			Event: models.AfEvent_FAILED_RESOURCES_ALLOCATION,
		}
		afEventsNotification.EvNotifs = append(afEventsNotification.EvNotifs, afNotif)
	}
	if afEventsNotification.EvNotifs != nil {
		sendSmPolicyRelatedAppSessionNotification(smPolicy, afEventsNotification, request.AccuUsageReports, successRules, failRules)
	}

	if errCause != "" {
		rsp := pcf_util.GetProblemDetail(errCause, pcf_util.ERROR_TRIGGER_EVENT)
		logger.SMpolicylog.Warnf(errCause)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	logger.SMpolicylog.Tracef("SMPolicy SmPolicyId[%s] Update", smPolicyId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, *smPolicyDecision)

}

func sendSmPolicyRelatedAppSessionNotification(smPolicy *pcf_context.UeSmPolicyData, notification models.EventsNotification, usageReports []models.AccuUsageReport, successRules, failRules []models.RuleReport) {

	for appSessionId := range smPolicy.AppSessions {
		if appSession, exist := pcf_context.PCF_Self().AppSessionPool[appSessionId]; exist {
			if len(appSession.Events) == 0 {
				continue
			}
			sessionNotif := models.EventsNotification{}
			for _, notif := range notification.EvNotifs {
				if _, found := appSession.Events[notif.Event]; found {
					switch notif.Event {
					case models.AfEvent_ACCESS_TYPE_CHANGE:
						sessionNotif.AccessType = notification.AccessType
						sessionNotif.RatType = notification.RatType
					case models.AfEvent_FAILED_RESOURCES_ALLOCATION:
						failItem := models.ResourcesAllocationInfo{
							McResourcStatus: models.MediaComponentResourcesStatus_INACTIVE,
						}
						flows := make(map[int32]models.Flows)
						for _, report := range failRules {
							for _, pccRuleId := range report.PccRuleIds {
								if key, exist := appSession.PccRuleIdMapToCompId[pccRuleId]; exist {
									items := strings.Split(key, "-")
									if items[0] != "appId" {
										compN, _ := strconv.Atoi(items[0])
										compN32 := int32(compN)
										if len(items) == 1 {
											// Comp
											flow := models.Flows{
												MedCompN: compN32,
											}
											failItem.Flows = append(failItem.Flows, flow)
										} else if len(items) == 2 {
											// have subComp
											fNum, _ := strconv.Atoi(items[1])
											fNum32 := int32(fNum)

											flow, exist := flows[compN32]
											if !exist {
												flow = models.Flows{
													MedCompN: compN32,
													FNums:    []int32{fNum32},
												}
											} else {
												flow.FNums = append(flow.FNums, fNum32)
											}
											flows[compN32] = flow
										}
									}
									// Release related resource
									delete(appSession.PccRuleIdMapToCompId, pccRuleId)
									delete(appSession.RelatedPccRuleIds, key)
								}
							}
						}
						for _, flow := range flows {
							failItem.Flows = append(failItem.Flows, flow)
						}
						if failItem.Flows != nil {
							sessionNotif.FailedResourcAllocReports = append(sessionNotif.FailedResourcAllocReports, failItem)
						} else {
							continue
						}
					case models.AfEvent_PLMN_CHG:
						sessionNotif.PlmnId = notification.PlmnId
					case models.AfEvent_QOS_NOTIF:
						for _, report := range sessionNotif.QncReports {
							for _, pccRuleId := range report.RefPccRuleIds {
								if _, exist := appSession.PccRuleIdMapToCompId[pccRuleId]; exist {
									sessionNotif.QncReports = append(sessionNotif.QncReports, report)
									break
								}
							}
						}
						if sessionNotif.QncReports == nil {
							continue
						}
					case models.AfEvent_SUCCESSFUL_RESOURCES_ALLOCATION:
						// Subscription to resources allocation outcome
						if successRules == nil {
							continue
						}
						flows := make(map[int32]models.Flows)
						for _, report := range successRules {
							for _, pccRuleId := range report.PccRuleIds {
								if key, exist := appSession.PccRuleIdMapToCompId[pccRuleId]; exist {
									items := strings.Split(key, "-")
									if items[0] != "appId" {
										compN, _ := strconv.Atoi(items[0])
										compN32 := int32(compN)
										if len(items) == 1 {
											// Comp
											flow := models.Flows{
												MedCompN: compN32,
											}
											notif.Flows = append(notif.Flows, flow)
										} else if len(items) == 2 {
											// have subComp
											fNum, _ := strconv.Atoi(items[1])
											fNum32 := int32(fNum)
											flow, exist := flows[compN32]
											if !exist {
												flow = models.Flows{
													MedCompN: compN32,
													FNums:    []int32{fNum32},
												}
											} else {
												flow.FNums = append(flow.FNums, fNum32)
											}
											flows[compN32] = flow
										}
									}

								}
							}
						}
						for _, flow := range flows {
							notif.Flows = append(notif.Flows, flow)
						}
						if notif.Flows == nil {
							continue
						}
					case models.AfEvent_USAGE_REPORT:
						for _, report := range usageReports {
							for _, pccRuleId := range appSession.RelatedPccRuleIds {
								if pccRule, exist := appSession.SmPolicyData.PolicyDecision.PccRules[pccRuleId]; exist {
									if pccRule.RefUmData != nil && pccRule.RefUmData[0] == report.RefUmIds {
										sessionNotif.UsgRep = &models.AccumulatedUsage{
											Duration:       report.TimeUsage,
											TotalVolume:    report.VolUsage,
											UplinkVolume:   report.VolUsageUplink,
											DownlinkVolume: report.VolUsageDownlink,
										}
										break
									}
								}
							}
							if sessionNotif.UsgRep != nil {
								sessionNotif.EvNotifs = append(sessionNotif.EvNotifs, notif)
								break
							}
						}
						fallthrough
					default:
						continue
					}
					sessionNotif.EvNotifs = append(sessionNotif.EvNotifs, notif)
				}
			}
			if sessionNotif.EvNotifs != nil {
				SendAppSessionEventNotification(appSession, sessionNotif)
			}

		}
	}
}

func SendSMPolicyUpdateNotification(ue *pcf_context.UeContext, smPolId string, request models.SmPolicyNotification) {
	logger.SMpolicylog.Tracef("Send SM Policy Update Notification")
	if ue == nil {
		logger.SMpolicylog.Warnln("SM Policy Update Notification Error[Ue is nil]")
		return
	}
	smPolicyData := ue.SmPolicyData[smPolId]
	if smPolicyData == nil || smPolicyData.PolicyContext == nil {
		logger.SMpolicylog.Warnf("SM Policy Update Notification Error[Can't find smPolId[%s] in UE(%s)]", smPolId, ue.Supi)
		return
	}
	client := pcf_util.GetNpcfSMPolicyCallbackClient()
	uri := smPolicyData.PolicyContext.NotificationUri
	if uri != "" {
		_, httpResponse, err := client.DefaultCallbackApi.SmPolicyUpdateNotification(context.Background(), uri, request)
		if err != nil {
			if httpResponse != nil {
				logger.SMpolicylog.Warnf("SM Policy Update Notification Error[%s]", httpResponse.Status)
			} else {
				logger.SMpolicylog.Warnf("SM Policy Update Notification Failed[%s]", err.Error())
			}
			return
		} else if httpResponse == nil {
			logger.SMpolicylog.Warnln("SM Policy Update Notification Failed[HTTP Response is nil]")
			return
		}
		if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
			logger.SMpolicylog.Warnf("SM Policy Update Notification Failed")
		} else {
			logger.SMpolicylog.Tracef("SM Policy Update Notification Success")
		}
	}

}

func SendSMPolicyTerminationRequestNotification(ue *pcf_context.UeContext, smPolId string, request models.TerminationNotification) {
	logger.SMpolicylog.Tracef("Send SM Policy Termination Request Notification")
	if ue == nil {
		logger.SMpolicylog.Warnln("SM Policy Termination Request Notification Error[Ue is nil]")
		return
	}
	smPolicyData := ue.SmPolicyData[smPolId]
	if smPolicyData == nil || smPolicyData.PolicyContext == nil {
		logger.SMpolicylog.Warnf("SM Policy Update Notification Error[Can't find smPolId[%s] in UE(%s)]", smPolId, ue.Supi)
		return
	}
	client := pcf_util.GetNpcfSMPolicyCallbackClient()
	uri := smPolicyData.PolicyContext.NotificationUri
	if uri != "" {
		rsp, err := client.DefaultCallbackApi.SmPolicyControlTerminationRequestNotification(context.Background(), uri, request)
		if err != nil {
			if rsp != nil {
				logger.AMpolicylog.Warnf("SM Policy Termination Request Notification Error[%s]", rsp.Status)
			} else {
				logger.AMpolicylog.Warnf("SM Policy Termination Request Notification Error[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.AMpolicylog.Warnln("SM Policy Termination Request Notification Error[HTTP Response is nil]")
			return
		}
		if rsp.StatusCode != http.StatusNoContent {
			logger.SMpolicylog.Warnf("SM Policy Termination Request Notification  Failed")
		} else {
			logger.SMpolicylog.Tracef("SM Policy Termination Request Notification Success")
		}
	}

}
