package udr_handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/logger"
	"free5gc/src/udr/udr_context"
	"free5gc/src/udr/udr_handler/udr_message"
	"free5gc/src/udr/udr_producer"
	"strconv"
	"strings"
	"time"
)

var HandlerLog *logrus.Entry

func init() {
	// init Pool
	HandlerLog = logger.HandlerLog
}

func Handle() {
	for {
		select {
		case msg, ok := <-udr_message.UdrChannel:
			if ok {
				udr_producer.CurrentResourceUri = udr_context.UDR_Self().GetIPv4Uri() + msg.HTTPRequest.URL.EscapedPath()
				switch msg.Event {
				case udr_message.EventCreateAccessAndMobilityData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateAccessAndMobilityData(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.AccessAndMobilityData))
				case udr_message.EventDeleteAccessAndMobilityData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleDeleteAccessAndMobilityData(msg.ResponseChan, ueId)
				case udr_message.EventQueryAccessAndMobilityData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryAccessAndMobilityData(msg.ResponseChan, ueId)
				case udr_message.EventQueryAmData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]
					udr_producer.HandleQueryAmData(msg.ResponseChan, ueId, servingPlmnId)
				case udr_message.EventAmfContext3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleAmfContext3gpp(msg.ResponseChan, ueId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventCreateAmfContext3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateAmfContext3gpp(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.Amf3GppAccessRegistration))
				case udr_message.EventQueryAmfContext3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryAmfContext3gpp(msg.ResponseChan, ueId)
				case udr_message.EventAmfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleAmfContextNon3gpp(msg.ResponseChan, ueId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventCreateAmfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateAmfContextNon3gpp(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.AmfNon3GppAccessRegistration))
				case udr_message.EventQueryAmfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryAmfContextNon3gpp(msg.ResponseChan, ueId)
				case udr_message.EventModifyAmfSubscriptionInfo:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleModifyAmfSubscriptionInfo(msg.ResponseChan, ueId, subsId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventModifyAuthentication:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleModifyAuthentication(msg.ResponseChan, ueId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventQueryAuthSubsData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryAuthSubsData(msg.ResponseChan, ueId)
				case udr_message.EventCreateAuthenticationSoR:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateAuthenticationSoR(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.SorData))
				case udr_message.EventQueryAuthSoR:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryAuthSoR(msg.ResponseChan, ueId)
				case udr_message.EventCreateAuthenticationStatus:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateAuthenticationStatus(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.AuthEvent))
				case udr_message.EventQueryAuthenticationStatus:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryAuthenticationStatus(msg.ResponseChan, ueId)
				case udr_message.EventApplicationDataInfluenceDataGet:
					// TODO
					udr_producer.HandleApplicationDataInfluenceDataGet(msg.ResponseChan)
				case udr_message.EventApplicationDataInfluenceDataInfluenceIdDelete:
					// TODO
					influenceId := msg.HTTPRequest.Params["influenceId"]
					udr_producer.HandleApplicationDataInfluenceDataInfluenceIdDelete(msg.ResponseChan, influenceId)
				case udr_message.EventApplicationDataInfluenceDataInfluenceIdPatch:
					// TODO
					influenceId := msg.HTTPRequest.Params["influenceId"]
					udr_producer.HandleApplicationDataInfluenceDataInfluenceIdPatch(msg.ResponseChan, influenceId, msg.HTTPRequest.Body.(models.TrafficInfluDataPatch))
				case udr_message.EventApplicationDataInfluenceDataInfluenceIdPut:
					// TODO
					influenceId := msg.HTTPRequest.Params["influenceId"]
					udr_producer.HandleApplicationDataInfluenceDataInfluenceIdPut(msg.ResponseChan, influenceId, msg.HTTPRequest.Body.(models.TrafficInfluData))
				case udr_message.EventApplicationDataInfluenceDataSubsToNotifyGet:
					// TODO
					udr_producer.HandleApplicationDataInfluenceDataSubsToNotifyGet(msg.ResponseChan)
				case udr_message.EventApplicationDataInfluenceDataSubsToNotifyPost:
					// TODO
					udr_producer.HandleApplicationDataInfluenceDataSubsToNotifyPost(msg.ResponseChan, msg.HTTPRequest.Body.(models.TrafficInfluSub))
				case udr_message.EventApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete:
					// TODO
					subscriptionId := msg.HTTPRequest.Params["subscriptionId"]
					udr_producer.HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete(msg.ResponseChan, subscriptionId)
				case udr_message.EventApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet:
					// TODO
					subscriptionId := msg.HTTPRequest.Params["subscriptionId"]
					udr_producer.HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet(msg.ResponseChan, subscriptionId)
				case udr_message.EventApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut:
					// TODO
					subscriptionId := msg.HTTPRequest.Params["subscriptionId"]
					udr_producer.HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut(msg.ResponseChan, subscriptionId, msg.HTTPRequest.Body.(models.TrafficInfluSub))
				case udr_message.EventApplicationDataPfdsAppIdDelete:
					// TODO
					appId := msg.HTTPRequest.Params["appId"]
					udr_producer.HandleApplicationDataPfdsAppIdDelete(msg.ResponseChan, appId)
				case udr_message.EventApplicationDataPfdsAppIdGet:
					// TODO
					appId := msg.HTTPRequest.Params["appId"]
					udr_producer.HandleApplicationDataPfdsAppIdGet(msg.ResponseChan, appId)
				case udr_message.EventApplicationDataPfdsAppIdPut:
					// TODO
					appId := msg.HTTPRequest.Params["appId"]
					udr_producer.HandleApplicationDataPfdsAppIdPut(msg.ResponseChan, appId, msg.HTTPRequest.Body.(models.PfdDataForApp))
				case udr_message.EventApplicationDataPfdsGet:
					// TODO
					udr_producer.HandleApplicationDataPfdsGet(msg.ResponseChan)
				case udr_message.EventExposureDataSubsToNotifyPost:
					// TODO
					udr_producer.HandleExposureDataSubsToNotifyPost(msg.ResponseChan, msg.HTTPRequest.Body.(models.ExposureDataSubscription))
				case udr_message.EventExposureDataSubsToNotifySubIdDelete:
					// TODO
					subId := msg.HTTPRequest.Params["subId"]
					udr_producer.HandleExposureDataSubsToNotifySubIdDelete(msg.ResponseChan, subId)
				case udr_message.EventExposureDataSubsToNotifySubIdPut:
					// TODO
					subId := msg.HTTPRequest.Params["subId"]
					udr_producer.HandleExposureDataSubsToNotifySubIdPut(msg.ResponseChan, subId, msg.HTTPRequest.Body.(models.ExposureDataSubscription))
				case udr_message.EventPolicyDataBdtDataBdtReferenceIdDelete:
					// TODO
					bdtReferenceId := msg.HTTPRequest.Params["bdtReferenceId"]
					udr_producer.HandlePolicyDataBdtDataBdtReferenceIdDelete(msg.ResponseChan, bdtReferenceId)
				case udr_message.EventPolicyDataBdtDataBdtReferenceIdGet:
					// TODO
					bdtReferenceId := msg.HTTPRequest.Params["bdtReferenceId"]
					udr_producer.HandlePolicyDataBdtDataBdtReferenceIdGet(msg.ResponseChan, bdtReferenceId)
				case udr_message.EventPolicyDataBdtDataBdtReferenceIdPut:
					// TODO
					bdtReferenceId := msg.HTTPRequest.Params["bdtReferenceId"]
					udr_producer.HandlePolicyDataBdtDataBdtReferenceIdPut(msg.ResponseChan, bdtReferenceId, msg.HTTPRequest.Body.(models.BdtData))
				case udr_message.EventPolicyDataBdtDataGet:
					// TODO
					udr_producer.HandlePolicyDataBdtDataGet(msg.ResponseChan)
				case udr_message.EventPolicyDataPlmnsPlmnIdUePolicySetGet:
					// TODO
					plmnId := msg.HTTPRequest.Params["plmnId"]
					udr_producer.HandlePolicyDataPlmnsPlmnIdUePolicySetGet(msg.ResponseChan, plmnId)
				case udr_message.EventPolicyDataSponsorConnectivityDataSponsorIdGet:
					// TODO
					sponsorId := msg.HTTPRequest.Params["sponsorId"]
					udr_producer.HandlePolicyDataSponsorConnectivityDataSponsorIdGet(msg.ResponseChan, sponsorId)
				case udr_message.EventPolicyDataSubsToNotifyPost:
					// TODO
					udr_producer.HandlePolicyDataSubsToNotifyPost(msg.ResponseChan, msg.HTTPRequest.Body.(models.PolicyDataSubscription))
				case udr_message.EventPolicyDataSubsToNotifySubsIdDelete:
					// TODO
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandlePolicyDataSubsToNotifySubsIdDelete(msg.ResponseChan, subsId)
				case udr_message.EventPolicyDataSubsToNotifySubsIdPut:
					// TODO
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandlePolicyDataSubsToNotifySubsIdPut(msg.ResponseChan, subsId, msg.HTTPRequest.Body.(models.PolicyDataSubscription))
				case udr_message.EventPolicyDataUesUeIdAmDataGet:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdAmDataGet(msg.ResponseChan, ueId)
				case udr_message.EventPolicyDataUesUeIdOperatorSpecificDataGet:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdOperatorSpecificDataGet(msg.ResponseChan, ueId)
				case udr_message.EventPolicyDataUesUeIdOperatorSpecificDataPatch:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdOperatorSpecificDataPatch(msg.ResponseChan, ueId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventPolicyDataUesUeIdOperatorSpecificDataPut:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdOperatorSpecificDataPut(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(map[string]models.OperatorSpecificDataContainer))
				case udr_message.EventPolicyDataUesUeIdSmDataGet:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]

					sNssai := models.Snssai{}
					sNssaiQuery := msg.HTTPRequest.Query.Get("snssai")
					err := json.Unmarshal([]byte(sNssaiQuery), &sNssai)
					if err != nil {
						HandlerLog.Warnln(err)
					}

					dnn := msg.HTTPRequest.Query.Get("dnn")

					udr_producer.HandlePolicyDataUesUeIdSmDataGet(msg.ResponseChan, ueId, sNssai, dnn)
				case udr_message.EventPolicyDataUesUeIdSmDataPatch:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdSmDataPatch(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(map[string]models.UsageMonData))
				case udr_message.EventPolicyDataUesUeIdSmDataUsageMonIdDelete:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					usageMonId := msg.HTTPRequest.Params["usageMonId"]
					udr_producer.HandlePolicyDataUesUeIdSmDataUsageMonIdDelete(msg.ResponseChan, ueId, usageMonId)
				case udr_message.EventPolicyDataUesUeIdSmDataUsageMonIdGet:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					usageMonId := msg.HTTPRequest.Params["usageMonId"]
					udr_producer.HandlePolicyDataUesUeIdSmDataUsageMonIdGet(msg.ResponseChan, ueId, usageMonId)
				case udr_message.EventPolicyDataUesUeIdSmDataUsageMonIdPut:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					usageMonId := msg.HTTPRequest.Params["usageMonId"]
					udr_producer.HandlePolicyDataUesUeIdSmDataUsageMonIdPut(msg.ResponseChan, ueId, usageMonId, msg.HTTPRequest.Body.(models.UsageMonData))
				case udr_message.EventPolicyDataUesUeIdUePolicySetGet:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdUePolicySetGet(msg.ResponseChan, ueId)
				case udr_message.EventPolicyDataUesUeIdUePolicySetPatch:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdUePolicySetPatch(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.UePolicySet))
				case udr_message.EventPolicyDataUesUeIdUePolicySetPut:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePolicyDataUesUeIdUePolicySetPut(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.UePolicySet))
				case udr_message.EventCreateAMFSubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleCreateAMFSubscriptions(msg.ResponseChan, ueId, subsId, msg.HTTPRequest.Body.([]models.AmfSubscriptionInfo))
				case udr_message.EventRemoveAmfSubscriptionsInfo:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleRemoveAmfSubscriptionsInfo(msg.ResponseChan, ueId, subsId)
				case udr_message.EventQueryEEData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryEEData(msg.ResponseChan, ueId)
				case udr_message.EventRemoveEeGroupSubscriptions:
					// TODO
					ueGroupId := msg.HTTPRequest.Params["ueGroupId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleRemoveEeGroupSubscriptions(msg.ResponseChan, ueGroupId, subsId)
				case udr_message.EventUpdateEeGroupSubscriptions:
					// TODO
					ueGroupId := msg.HTTPRequest.Params["ueGroupId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleUpdateEeGroupSubscriptions(msg.ResponseChan, ueGroupId, subsId, msg.HTTPRequest.Body.(models.EeSubscription))
				case udr_message.EventCreateEeGroupSubscriptions:
					// TODO
					ueGroupId := msg.HTTPRequest.Params["ueGroupId"]
					udr_producer.HandleCreateEeGroupSubscriptions(msg.ResponseChan, ueGroupId, msg.HTTPRequest.Body.(models.EeSubscription))
				case udr_message.EventQueryEeGroupSubscriptions:
					// TODO
					ueGroupId := msg.HTTPRequest.Params["ueGroupId"]
					udr_producer.HandleQueryEeGroupSubscriptions(msg.ResponseChan, ueGroupId)
				case udr_message.EventRemoveeeSubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleRemoveeeSubscriptions(msg.ResponseChan, ueId, subsId)
				case udr_message.EventUpdateEesubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleUpdateEesubscriptions(msg.ResponseChan, ueId, subsId, msg.HTTPRequest.Body.(models.EeSubscription))
				case udr_message.EventCreateEeSubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateEeSubscriptions(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.EeSubscription))
				case udr_message.EventQueryeesubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryeesubscriptions(msg.ResponseChan, ueId)
				case udr_message.EventPatchOperSpecData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandlePatchOperSpecData(msg.ResponseChan, ueId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventQueryOperSpecData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQueryOperSpecData(msg.ResponseChan, ueId)
				case udr_message.EventGetppData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleGetppData(msg.ResponseChan, ueId)
				case udr_message.EventCreateSessionManagementData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					pduSessionId, _ := strconv.ParseInt(msg.HTTPRequest.Params["pduSessionId"], 10, 64)
					udr_producer.HandleCreateSessionManagementData(msg.ResponseChan, ueId, int32(pduSessionId), msg.HTTPRequest.Body.(models.PduSessionManagementData))
				case udr_message.EventDeleteSessionManagementData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					pduSessionId, _ := strconv.ParseInt(msg.HTTPRequest.Params["pduSessionId"], 10, 64)
					udr_producer.HandleDeleteSessionManagementData(msg.ResponseChan, ueId, int32(pduSessionId))
				case udr_message.EventQuerySessionManagementData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					pduSessionId, _ := strconv.ParseInt(msg.HTTPRequest.Params["pduSessionId"], 10, 64)
					udr_producer.HandleQuerySessionManagementData(msg.ResponseChan, ueId, int32(pduSessionId))
				case udr_message.EventQueryProvisionedData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]
					udr_producer.HandleQueryProvisionedData(msg.ResponseChan, ueId, servingPlmnId)
				case udr_message.EventModifyPpData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleModifyPpData(msg.ResponseChan, ueId, msg.HTTPRequest.Body.([]models.PatchItem))
				case udr_message.EventGetAmfSubscriptionInfo:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleGetAmfSubscriptionInfo(msg.ResponseChan, ueId, subsId)
				case udr_message.EventGetIdentityData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleGetIdentityData(msg.ResponseChan, ueId)
				case udr_message.EventGetOdbData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleGetOdbData(msg.ResponseChan, ueId)
				case udr_message.EventGetSharedData:
					// TODO
					var sharedDataIds []string
					if len(msg.HTTPRequest.Query["shared-data-ids"]) != 0 {
						sharedDataIds = msg.HTTPRequest.Query["shared-data-ids"]
						if strings.Contains(sharedDataIds[0], ",") {
							sharedDataIds = strings.Split(sharedDataIds[0], ",")
						}
					}
					udr_producer.HandleGetSharedData(msg.ResponseChan, sharedDataIds)
				case udr_message.EventRemovesdmSubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleRemovesdmSubscriptions(msg.ResponseChan, ueId, subsId)
				case udr_message.EventUpdatesdmsubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleUpdatesdmsubscriptions(msg.ResponseChan, ueId, subsId, msg.HTTPRequest.Body.(models.SdmSubscription))
				case udr_message.EventCreateSdmSubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateSdmSubscriptions(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.SdmSubscription))
				case udr_message.EventQuerysdmsubscriptions:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQuerysdmsubscriptions(msg.ResponseChan, ueId)
				case udr_message.EventQuerySmData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]

					singleNssai := models.Snssai{}
					singleNssaiQuery := msg.HTTPRequest.Query.Get("single-nssai")
					err := json.Unmarshal([]byte(singleNssaiQuery), &singleNssai)
					if err != nil {
						HandlerLog.Warnln(err)
					}

					dnn := msg.HTTPRequest.Query.Get("dnn")

					udr_producer.HandleQuerySmData(msg.ResponseChan, ueId, servingPlmnId, singleNssai, dnn)
				case udr_message.EventCreateSmfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					pduSessionId, _ := strconv.ParseInt(msg.HTTPRequest.Params["pduSessionId"], 10, 64)
					udr_producer.HandleCreateSmfContextNon3gpp(msg.ResponseChan, ueId, int32(pduSessionId), msg.HTTPRequest.Body.(models.SmfRegistration))
				case udr_message.EventDeleteSmfContext:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					pduSessionId := msg.HTTPRequest.Params["pduSessionId"]
					udr_producer.HandleDeleteSmfContext(msg.ResponseChan, ueId, pduSessionId)
				case udr_message.EventQuerySmfRegistration:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					pduSessionId := msg.HTTPRequest.Params["pduSessionId"]
					udr_producer.HandleQuerySmfRegistration(msg.ResponseChan, ueId, pduSessionId)
				case udr_message.EventQuerySmfRegList:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQuerySmfRegList(msg.ResponseChan, ueId)
				case udr_message.EventQuerySmfSelectData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]
					udr_producer.HandleQuerySmfSelectData(msg.ResponseChan, ueId, servingPlmnId)
				case udr_message.EventCreateSmsfContext3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateSmsfContext3gpp(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.SmsfRegistration))
				case udr_message.EventDeleteSmsfContext3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleDeleteSmsfContext3gpp(msg.ResponseChan, ueId)
				case udr_message.EventQuerySmsfContext3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQuerySmsfContext3gpp(msg.ResponseChan, ueId)
				case udr_message.EventCreateSmsfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleCreateSmsfContextNon3gpp(msg.ResponseChan, ueId, msg.HTTPRequest.Body.(models.SmsfRegistration))
				case udr_message.EventDeleteSmsfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleDeleteSmsfContextNon3gpp(msg.ResponseChan, ueId)
				case udr_message.EventQuerySmsfContextNon3gpp:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					udr_producer.HandleQuerySmsfContextNon3gpp(msg.ResponseChan, ueId)
				case udr_message.EventQuerySmsMngData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]
					udr_producer.HandleQuerySmsMngData(msg.ResponseChan, ueId, servingPlmnId)
				case udr_message.EventQuerySmsData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]
					udr_producer.HandleQuerySmsData(msg.ResponseChan, ueId, servingPlmnId)
				case udr_message.EventPostSubscriptionDataSubscriptions:
					// TODO
					udr_producer.HandlePostSubscriptionDataSubscriptions(msg.ResponseChan, msg.HTTPRequest.Body.(models.SubscriptionDataSubscriptions))
				case udr_message.EventRemovesubscriptionDataSubscriptions:
					// TODO
					subsId := msg.HTTPRequest.Params["subsId"]
					udr_producer.HandleRemovesubscriptionDataSubscriptions(msg.ResponseChan, subsId)
				case udr_message.EventQueryTraceData:
					// TODO
					ueId := msg.HTTPRequest.Params["ueId"]
					servingPlmnId := msg.HTTPRequest.Params["servingPlmnId"]
					udr_producer.HandleQueryTraceData(msg.ResponseChan, ueId, servingPlmnId)
				default:
					HandlerLog.Warnf("Event[%d] has not been implemented", msg.Event)
				}
			} else {
				HandlerLog.Errorln("Channel closed!")
			}

		case <-time.After(time.Second * 1):

		}
	}
}
