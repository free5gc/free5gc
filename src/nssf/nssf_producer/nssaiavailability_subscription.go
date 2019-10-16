/*
 * NSSF NSSAI Availability
 *
 * NSSF NSSAI Availability Service
 */

package nssf_producer

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	. "free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/logger"
	"free5gc/src/nssf/util"
)

// Get available subscription ID from configuration
// In this implementation, string converted from 32-bit integer is used as subscription ID
func getUnusedSubscriptionId() (string, error) {
	var idx uint32 = 1
	for _, subscription := range factory.NssfConfig.Subscriptions {
		tempId, _ := strconv.Atoi(subscription.SubscriptionId)
		if uint32(tempId) == idx {
			if idx == math.MaxUint32 {
				return "", fmt.Errorf("No available subscription ID")
			}
			idx = idx + 1
		} else {
			break
		}
	}
	return strconv.Itoa(int(idx)), nil
}

// NSSAIAvailability subscription POST method
func subscriptionPost(n NssfEventSubscriptionCreateData, s *NssfEventSubscriptionCreatedData, d *ProblemDetails) (status int) {
	var subscription factory.Subscription
	tempId, err := getUnusedSubscriptionId()
	if err != nil {
		logger.Nssaiavailability.Warnf(err.Error())

		*d = ProblemDetails{
			Title:  util.UNSUPPORTED_RESOURCE,
			Status: http.StatusNotFound,
			Detail: err.Error(),
		}

		status = http.StatusNotFound
		return
	}

	subscription.SubscriptionId = tempId
	subscription.SubscriptionData = new(NssfEventSubscriptionCreateData)
	*subscription.SubscriptionData = n

	factory.NssfConfig.Subscriptions = append(factory.NssfConfig.Subscriptions, subscription)

	s.SubscriptionId = subscription.SubscriptionId
	if !subscription.SubscriptionData.Expiry.IsZero() {
		s.Expiry = new(time.Time)
		*s.Expiry = *subscription.SubscriptionData.Expiry
	}
	s.AuthorizedNssaiAvailabilityData = util.AuthorizeOfTaListFromConfig(subscription.SubscriptionData.TaiList)

	status = http.StatusCreated
	return
}

func subscriptionDelete(subscriptionId string, d *ProblemDetails) (status int) {
	for i, subscription := range factory.NssfConfig.Subscriptions {
		if subscription.SubscriptionId == subscriptionId {
			factory.NssfConfig.Subscriptions = append(factory.NssfConfig.Subscriptions[:i],
				factory.NssfConfig.Subscriptions[i+1:]...)

			status = http.StatusNoContent
			return
		}
	}

	// No specific subscription ID exists
	problemDetail := fmt.Sprintf("Subscription ID '%s' is not available", subscriptionId)
	*d = ProblemDetails{
		Title:  util.UNSUPPORTED_RESOURCE,
		Status: http.StatusNotFound,
		Detail: problemDetail,
	}

	status = http.StatusNotFound
	return
}
