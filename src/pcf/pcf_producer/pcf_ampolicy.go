package pcf_producer

import (
	"context"
	"fmt"
	"github.com/mohae/deepcopy"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_consumer"
	"free5gc/src/pcf/pcf_context"
	"free5gc/src/pcf/pcf_handler/pcf_message"
	"free5gc/src/pcf/pcf_util"
	"net/http"
	"reflect"
)

func DeletePoliciesPolAssoId(httpChannel chan pcf_message.HttpResponseMessage, polAssoId string) {

	logger.AMpolicylog.Traceln("Handle Policy Association Delete")

	ue := pcf_context.PCF_Self().PCFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		rsp := pcf_util.GetProblemDetail("polAssoId not found  in PCF", pcf_util.CONTEXT_NOT_FOUND)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	delete(ue.AMPolicyData, polAssoId)
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusNoContent, nil)
}

// PoliciesPolAssoIdGet -
func GetPoliciesPolAssoId(httpChannel chan pcf_message.HttpResponseMessage, polAssoId string) {

	logger.AMpolicylog.Traceln("Handle Policy Association Get")

	ue := pcf_context.PCF_Self().PCFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		rsp := pcf_util.GetProblemDetail("polAssoId not found  in PCF", pcf_util.CONTEXT_NOT_FOUND)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	amPolicyData := ue.AMPolicyData[polAssoId]
	rsp := models.PolicyAssociation{
		SuppFeat: amPolicyData.SuppFeat,
	}
	if amPolicyData.Rfsp != 0 {
		rsp.Rfsp = amPolicyData.Rfsp
	}
	if amPolicyData.ServAreaRes != nil {
		rsp.ServAreaRes = amPolicyData.ServAreaRes
	}
	if amPolicyData.Triggers != nil {
		rsp.Triggers = amPolicyData.Triggers
		for _, trigger := range amPolicyData.Triggers {
			if trigger == models.RequestTrigger_PRA_CH {
				rsp.Pras = amPolicyData.Pras
				break
			}
		}
	}
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, rsp)

}
func UpdatePostPoliciesPolAssoId(httpChannel chan pcf_message.HttpResponseMessage, polAssoId string, request models.PolicyAssociationUpdateRequest) {

	logger.AMpolicylog.Traceln("Handle Policy Association Update")

	ue := pcf_context.PCF_Self().PCFUeFindByPolicyId(polAssoId)
	if ue == nil || ue.AMPolicyData[polAssoId] == nil {
		rsp := pcf_util.GetProblemDetail("polAssoId not found  in PCF", pcf_util.CONTEXT_NOT_FOUND)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}

	amPolicyData := ue.AMPolicyData[polAssoId]
	var rsp models.PolicyUpdate
	if request.NotificationUri != "" {
		amPolicyData.NotificationUri = request.NotificationUri
	}
	if request.AltNotifIpv4Addrs != nil {
		amPolicyData.AltNotifIpv4Addrs = request.AltNotifIpv4Addrs
	}
	if request.AltNotifIpv6Addrs != nil {
		amPolicyData.AltNotifIpv6Addrs = request.AltNotifIpv6Addrs
	}
	for _, trigger := range request.Triggers {
		//TODO: Modify the value according to policies
		switch trigger {
		case models.RequestTrigger_LOC_CH:
			//TODO: report to AF subscriber
			if request.UserLoc == nil {
				rsp := pcf_util.GetProblemDetail("UserLoc are nli", pcf_util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("UserLoc doesn't exist in Policy Association Requset Update while Triggers include LOC_CH")
				pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
				return
			}
			amPolicyData.UserLoc = request.UserLoc
			logger.AMpolicylog.Infof("Ue[%s] UserLocation %+v", ue.Supi, amPolicyData.UserLoc)
		case models.RequestTrigger_PRA_CH:
			if request.PraStatuses == nil {
				rsp := pcf_util.GetProblemDetail("PraStatuses are nli", pcf_util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("PraStatuses doesn't exist in Policy Association Requset Update while Triggers include PRA_CH")
				pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
				return
			}
			for praId, praInfo := range request.PraStatuses {
				//TODO: report to AF subscriber
				logger.AMpolicylog.Infof("Policy Association Presence Id[%s] change state to %s", praId, praInfo.PresenceState)
			}
		case models.RequestTrigger_SERV_AREA_CH:
			if request.ServAreaRes == nil {
				rsp := pcf_util.GetProblemDetail("ServAreaRes are nli", pcf_util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("ServAreaRes doesn't exist in Policy Association Requset Update while Triggers include SERV_AREA_CH")
				pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
				return
			} else {
				amPolicyData.ServAreaRes = request.ServAreaRes
				rsp.ServAreaRes = request.ServAreaRes
			}
		case models.RequestTrigger_RFSP_CH:
			if request.Rfsp == 0 {
				rsp := pcf_util.GetProblemDetail("Rfsp are nli", pcf_util.ERROR_REQUEST_PARAMETERS)
				logger.AMpolicylog.Warnln("Rfsp doesn't exist in Policy Association Requset Update while Triggers include RFSP_CH")
				pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
				return
			} else {
				amPolicyData.Rfsp = request.Rfsp
				rsp.Rfsp = request.Rfsp
			}
		}
	}
	//TODO: handle TraceReq
	//TODO: Change Request Trigger Policies if needed
	rsp.Triggers = amPolicyData.Triggers
	//TODO: Change Policies if needed
	// rsp.Pras
	pcf_message.SendHttpResponseMessage(httpChannel, nil, http.StatusOK, rsp)

}

// Create AM Policy
func PostPolicies(httpChannel chan pcf_message.HttpResponseMessage, request models.PolicyAssociationRequest) {
	var rsp models.PolicyAssociation
	var err error

	logger.AMpolicylog.Traceln("Handle Policy Association Request")

	pcfSelf := pcf_context.PCF_Self()
	ue := pcfSelf.UePool[request.Supi]
	if ue == nil {
		ue, err = pcfSelf.NewPCFUe(request.Supi)
		if err != nil {
			// supi format dose not match "imsi-..."
			rsp := pcf_util.GetProblemDetail("Supi Format Error", pcf_util.ERROR_REQUEST_PARAMETERS)
			logger.AMpolicylog.Errorln(err.Error())
			pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
			return
		}
	}
	udrUri := getUdrUri(ue)
	if udrUri == "" {
		// Can't find any UDR support this Ue
		delete(pcfSelf.UePool, ue.Supi)
		rsp := pcf_util.GetProblemDetail("Ue is not supported in PCF", pcf_util.USER_UNKNOWN)
		logger.AMpolicylog.Errorf("Ue[%s] is not supported in PCF", ue.Supi)
		pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
		return
	}
	ue.UdrUri = udrUri

	rsp.Request = deepcopy.Copy(&request).(*models.PolicyAssociationRequest)
	assolId := fmt.Sprintf("%s-%d", ue.Supi, ue.PolAssociationIDGenerator)
	amPolicy := ue.AMPolicyData[assolId]

	if amPolicy == nil || amPolicy.AmPolicyData == nil {
		client := pcf_util.GetNudrClient(udrUri)
		var response *http.Response
		amData, response, err := client.DefaultApi.PolicyDataUesUeIdAmDataGet(context.Background(), ue.Supi)
		if err != nil || response == nil || response.StatusCode != http.StatusOK {
			rsp := pcf_util.GetProblemDetail("Can't find UE AM Policy Data in UDR", pcf_util.USER_UNKNOWN)
			logger.AMpolicylog.Errorf("Can't find UE[%s] AM Policy Data in UDR", ue.Supi)
			pcf_message.SendHttpResponseMessage(httpChannel, nil, int(rsp.Status), rsp)
			return
		}
		if amPolicy == nil {
			amPolicy = ue.NewUeAMPolicyData(assolId, request)
		}
		amPolicy.AmPolicyData = &amData
	}

	// TODO: according to PCF Policy to determine ServAreaRes, Rfsp, SuppFeat
	// amPolicy.ServAreaRes =
	// amPolicy.Rfsp =
	amPolicy.SuppFeat = pcf_util.GetNegotiateSuppFeat(request.SuppFeat, pcfSelf.PcfSuppFeats[models.ServiceName_NPCF_AM_POLICY_CONTROL])
	if amPolicy.Rfsp != 0 {
		rsp.Rfsp = amPolicy.Rfsp
	}
	rsp.SuppFeat = amPolicy.SuppFeat
	// TODO: add Reports
	// rsp.Triggers
	// rsp.Pras
	ue.PolAssociationIDGenerator++
	// Create location header for update, delete, get
	locationHeader := pcf_util.GetResourceUri(models.ServiceName_NPCF_AM_POLICY_CONTROL, assolId)
	headers := http.Header{
		"Location": {locationHeader},
	}
	logger.AMpolicylog.Tracef("AMPolicy association Id[%s] Create", assolId)
	pcf_message.SendHttpResponseMessage(httpChannel, headers, http.StatusCreated, rsp)

	if request.Guami != nil {
		// if consumer is AMF then subscribe this AMF Status
		for _, statusSubsData := range pcfSelf.AMFStatusSubsData {
			for _, guami := range statusSubsData.GuamiList {
				if reflect.DeepEqual(guami, request.Guami) {
					amPolicy.AmfStatusChangeSubscription = &statusSubsData
				}
			}
		}
	}
}

// Send AM Policy Update to AMF if policy has changed
func SendAMPolicyUpdateNotification(ue *pcf_context.UeContext, PolId string, request models.PolicyUpdate) {
	if ue == nil {
		logger.AMpolicylog.Warnln("Policy Update Notification Error[Ue is nil]")
		return
	}
	amPolicyData := ue.AMPolicyData[PolId]
	if amPolicyData == nil {
		logger.AMpolicylog.Warnf("Policy Update Notification Error[Can't find polAssoId[%s] in UE(%s)]", PolId, ue.Supi)
		return
	}
	client := pcf_util.GetNpcfAMPolicyCallbackClient()
	uri := amPolicyData.NotificationUri
	for uri != "" {

		rsp, err := client.DefaultCallbackApi.PolicyUpdateNotification(context.Background(), uri, request)
		if err != nil {
			if rsp != nil && rsp.StatusCode != http.StatusNoContent {
				logger.AMpolicylog.Warnf("Policy Update Notification Error[%s]", rsp.Status)
			} else {
				logger.AMpolicylog.Warnf("Policy Update Notification Failed[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.AMpolicylog.Warnln("Policy Update Notification Failed[HTTP Response is nil]")
			return
		}
		if rsp.StatusCode == http.StatusTemporaryRedirect {
			// for redirect case, resend the notification to redirect target
			uRI, err := rsp.Location()
			if err != nil {
				logger.AMpolicylog.Warnln("Policy Update Notification Redirect Need Supply URI")
				return
			}
			uri = uRI.String()
			continue
		}

		logger.AMpolicylog.Infoln("Policy Update Notification Success")
		return
	}

}

// Send AM Policy Update to AMF if policy has been terminated
func SendAMPolicyTerminationRequestNotification(ue *pcf_context.UeContext, PolId string, request models.TerminationNotification) {
	if ue == nil {
		logger.AMpolicylog.Warnln("Policy Assocition Termination Request Notification Error[Ue is nil]")
		return
	}
	amPolicyData := ue.AMPolicyData[PolId]
	if amPolicyData == nil {
		logger.AMpolicylog.Warnf("Policy Assocition Termination Request Notification Error[Can't find polAssoId[%s] in UE(%s)]", PolId, ue.Supi)
		return
	}
	client := pcf_util.GetNpcfAMPolicyCallbackClient()
	uri := amPolicyData.NotificationUri
	for uri != "" {

		rsp, err := client.DefaultCallbackApi.PolicyAssocitionTerminationRequestNotification(context.Background(), uri, request)
		if err != nil {
			if rsp != nil && rsp.StatusCode != http.StatusNoContent {
				logger.AMpolicylog.Warnf("Policy Assocition Termination Request Notification Error[%s]", rsp.Status)
			} else {
				logger.AMpolicylog.Warnf("Policy Assocition Termination Request Notification Failed[%s]", err.Error())
			}
			return
		} else if rsp == nil {
			logger.AMpolicylog.Warnln("Policy Assocition Termination Request Notification Failed[HTTP Response is nil]")
			return
		}
		if rsp.StatusCode == http.StatusTemporaryRedirect {
			// for redirect case, resend the notification to redirect target
			uRI, err := rsp.Location()
			if err != nil {
				logger.AMpolicylog.Warnln("Policy Assocition Termination Request Notification Redirect Need Supply URI")
				return
			}
			uri = uRI.String()
			continue
		}
		return
	}

}

// returns UDR Uri of Ue, if ue.UdrUri dose not exist, query NRF to get supported Udr Uri
func getUdrUri(ue *pcf_context.UeContext) string {
	if ue.UdrUri != "" {
		return ue.UdrUri
	}
	return pcf_consumer.SendNFIntancesUDR(pcf_context.PCF_Self().NrfUri, ue.Supi)
}
