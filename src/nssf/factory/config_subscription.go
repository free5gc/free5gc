/*
 * NSSF Configuration Factory
 */

package factory

import (
	. "free5gc/lib/openapi/models"
)

type Subscription struct {
	SubscriptionId string `yaml:"subscriptionId"`

	SubscriptionData *NssfEventSubscriptionCreateData `yaml:"subscriptionData"`
}
