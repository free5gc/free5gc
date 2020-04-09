package udr_producer

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/udr_producer/udr_producer_callback"
)

func PreHandleOnDataChangeNotify(ueId string, resourceId string, patchItems []models.PatchItem, origValue interface{}, newValue interface{}) {
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

	go udr_producer_callback.SendOnDataChangeNotify(ueId, notifyItems)
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

	go udr_producer_callback.SendPolicyDataChangeNotification(policyDataChangeNotification)
}
