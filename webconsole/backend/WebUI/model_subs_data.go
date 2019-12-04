package WebUI

import "free5gc/lib/openapi/models"

type SubsData struct {
	PlmnID                            string                                   `json:"plmnID"`
	UeId                              string                                   `json:"ueId"`
	AuthenticationSubscription        models.AuthenticationSubscription        `json:"AuthenticationSubscription"`
	AccessAndMobilitySubscriptionData models.AccessAndMobilitySubscriptionData `json:"AccessAndMobilitySubscriptionData"`
	SmfSelectionSubscriptionData      models.SmfSelectionSubscriptionData      `json:"SmfSelectionSubscriptionData"`
}
