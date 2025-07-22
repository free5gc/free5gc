package processor

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/udr/DataRepository"
	udr_context "github.com/free5gc/udr/internal/context"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/internal/util"
)

var CurrentResourceUri string

func PreHandleOnDataChangeNotify(ueId string, resourceId string, patchItems []models.PatchItem,
	origValue map[string]interface{}, newValue map[string]interface{},
) {
	notifyItems := []models.NotifyItem{}
	changes := []models.ChangeItem{}

	for _, patchItem := range patchItems {
		change := models.ChangeItem{
			Op:        models.ChangeType(patchItem.Op),
			Path:      patchItem.Path,
			From:      patchItem.From,
			OrigValue: origValue,
			NewValue:  newValue,
		}
		changes = append(changes, change)
	}

	notifyItem := models.NotifyItem{
		ResourceId: resourceId,
		Changes:    changes,
	}

	notifyItems = append(notifyItems, notifyItem)

	go SendOnDataChangeNotify(ueId, notifyItems)
}

func PreHandlePolicyDataChangeNotification(ueId string, dataId string, value interface{}) {
	policyDataChangeNotification := models.PolicyDataChangeNotification{}

	if ueId != "" {
		policyDataChangeNotification.UeId = ueId
	}

	switch v := value.(type) {
	case models.AmPolicyData:
		policyDataChangeNotification.AmPolicyData = &v
	case models.UePolicySet:
		policyDataChangeNotification.UePolicySet = &v
	case models.SmPolicyData:
		policyDataChangeNotification.SmPolicyData = &v
	case models.UsageMonData:
		policyDataChangeNotification.UsageMonId = dataId
		policyDataChangeNotification.UsageMonData = &v
	case models.SponsorConnectivityData:
		policyDataChangeNotification.SponsorId = dataId
		policyDataChangeNotification.SponsorConnectivityData = &v
	case models.BdtData:
		policyDataChangeNotification.BdtRefId = dataId
		policyDataChangeNotification.BdtData = &v
	default:
		return
	}

	go SendPolicyDataChangeNotification(policyDataChangeNotification)
}

func PreHandleInfluenceDataUpdateNotification(influenceId string, original, modified *models.TrafficInfluData) {
	resUri := fmt.Sprintf("%s/application-data/influenceData/%s",
		udr_context.GetSelf().GetIPv4GroupUri(udr_context.NUDR_DR), influenceId)

	go SendInfluenceDataUpdateNotification(resUri, original, modified)
}

func SendOnDataChangeNotify(ueId string, notifyItems []models.NotifyItem) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.HttpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	udrSelf := udr_context.GetSelf()
	configuration := DataRepository.NewConfiguration()
	client := DataRepository.NewAPIClient(configuration)

	for _, subscriptionDataSubscription := range udrSelf.SubscriptionDataSubscriptions {
		if ueId == subscriptionDataSubscription.UeId {
			onDataChangeNotifyUrl := subscriptionDataSubscription.CallbackReference

			dataChangeReq := DataRepository.SubscriptionDataSubscriptionsOnDataChangePostRequest{}
			dataChangeReq.DataChangeNotify.UeId = ueId
			dataChangeReq.DataChangeNotify.OriginalCallbackReference = []string{
				subscriptionDataSubscription.OriginalCallbackReference,
			}
			dataChangeReq.DataChangeNotify.NotifyItems = notifyItems
			rsp, err := client.SubsToNotifyCollectionApi.SubscriptionDataSubscriptionsOnDataChangePost(
				context.TODO(), onDataChangeNotifyUrl, &dataChangeReq)

			if err != nil {
				logger.SBILog.Errorln(err.Error())
			} else if rsp == nil {
				logger.SBILog.Errorln("Empty SubscriptionDataSubscriptionsOnDataChangePost response")
			}
		}
	}
}

func SendPolicyDataChangeNotification(policyDataChangeNotification models.PolicyDataChangeNotification) {
	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.HttpLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	udrSelf := udr_context.GetSelf()

	for _, policyDataSubscription := range udrSelf.PolicyDataSubscriptions {
		policyDataChangeNotificationUrl := policyDataSubscription.NotificationUri

		configuration := DataRepository.NewConfiguration()
		client := DataRepository.NewAPIClient(configuration)

		req := DataRepository.CreateIndividualPolicyDataSubscriptionPolicyDataChangeNotificationPostRequest{
			PolicyDataChangeNotification: []models.PolicyDataChangeNotification{
				policyDataChangeNotification,
			},
		}

		rsp, err := client.PolicyDataSubscriptionsCollectionApi.
			CreateIndividualPolicyDataSubscriptionPolicyDataChangeNotificationPost(context.TODO(),
				policyDataChangeNotificationUrl, &req)

		if err != nil {
			logger.SBILog.Errorln(err.Error())
		} else if rsp == nil {
			logger.SBILog.Errorln("Empty CreateIndividualPolicyDataSubscriptionPolicyDataChangeNotificationPost response")
		}
	}
}

func SendInfluenceDataUpdateNotification(resUri string, original, modified *models.TrafficInfluData) {
	udrSelf := udr_context.GetSelf()

	configuration := DataRepository.NewConfiguration()
	client := DataRepository.NewAPIClient(configuration)

	var trafficInfluDataNotif models.TrafficInfluDataNotif
	trafficInfluDataNotif.ResUri = resUri
	udrSelf.InfluenceDataSubscriptions.Range(func(key, value interface{}) bool {
		influenceDataSubscription, ok := value.(*models.TrafficInfluSub)
		if !ok {
			logger.HttpLog.Errorf("Failed to load influenceData subscription ID [%+v]", key)
			return true
		}
		influenceDataChangeNotificationUrl := influenceDataSubscription.NotificationUri

		// Check if the modified data is subscribed
		// If positive, send notification about the update
		if checkInfluenceDataSubscription(modified, influenceDataSubscription) {
			logger.HttpLog.Tracef("Send notification about update of influence data")
			trafficInfluDataNotif.TrafficInfluData = modified

			req := DataRepository.CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPostRequest{
				RequestBody: []interface{}{trafficInfluDataNotif},
			}

			rsp, err := client.InfluenceDataSubscriptionsCollectionApi.
				CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPost(
					context.TODO(), influenceDataChangeNotificationUrl, &req)

			if err != nil {
				logger.SBILog.Errorln(err.Error())
			} else if rsp == nil {
				logger.SBILog.Errorln(
					"Empty CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPost response")
			}
		} else if checkInfluenceDataSubscription(original, influenceDataSubscription) {
			// If the modified data is not subscribed or nil, check if the original data is subscribed
			// If positive, send notification about the removal
			logger.HttpLog.Tracef("Send notification about removal of influence data")
			trafficInfluDataNotif.TrafficInfluData = nil
			req := DataRepository.CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPostRequest{
				RequestBody: []interface{}{trafficInfluDataNotif},
			}

			rsp, err := client.InfluenceDataSubscriptionsCollectionApi.
				CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPost(
					context.TODO(), influenceDataChangeNotificationUrl, &req)

			if err != nil {
				logger.SBILog.Errorln(err.Error())
			} else if rsp == nil {
				logger.SBILog.Errorln(
					"Empty CreateIndividualInfluenceDataSubscriptionTrafficInfluenceDataChangeNotificationPost response")
			}
		}
		return true
	})
}

func checkInfluenceDataSubscription(data *models.TrafficInfluData, sub *models.TrafficInfluSub) bool {
	if data == nil || sub == nil {
		return false
	}
	if data.Dnn != "" && !util.Contain(data.Dnn, sub.Dnns) {
		return false
	} else if data.Snssai != nil && !util.Contain(*data.Snssai, sub.Snssais) {
		return false
	} else if data.InterGroupId != "AnyUE" {
		if data.InterGroupId != "" && !util.Contain(data.InterGroupId, sub.InternalGroupIds) {
			return false
		} else if data.Supi != "" && !util.Contain(data.Supi, sub.Supis) {
			return false
		}
	}
	return true
}
