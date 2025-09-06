/*
 * BSF Subscriptions Processor
 */

package processor

import (
	"net/http"

	"github.com/gin-gonic/gin"

	bsfContext "github.com/free5gc/bsf/internal/context"
	"github.com/free5gc/bsf/internal/logger"
	"github.com/free5gc/bsf/internal/util"
	"github.com/free5gc/openapi/models"
)

// CreateIndividualSubcription handles POST /subscriptions
func CreateIndividualSubcription(c *gin.Context) {
	logger.ProcLog.Infof("Handle CreateIndividualSubcription")

	var request models.BsfSubscription
	if err := c.ShouldBindJSON(&request); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Convert to internal representation
	subscription := &bsfContext.BsfSubscription{
		Events:            request.Events,
		NotifUri:          request.NotifUri,
		NotifCorreId:      request.NotifCorreId,
		Supi:              request.Supi,
		Gpsi:              util.StringToPtr(request.Gpsi),
		SnssaiDnnPairs:    request.SnssaiDnnPairs,
		AddSnssaiDnnPairs: request.AddSnssaiDnnPairs,
		SuppFeat:          util.StringToPtr(request.SuppFeat),
	}

	// Create new subscription
	subId := bsfContext.BsfSelf.CreateSubscription(subscription)

	// Convert back to response format
	response := models.BsfSubscriptionResp{
		Events:            subscription.Events,
		NotifUri:          subscription.NotifUri,
		NotifCorreId:      subscription.NotifCorreId,
		Supi:              subscription.Supi,
		Gpsi:              util.PtrToString(subscription.Gpsi),
		SnssaiDnnPairs:    subscription.SnssaiDnnPairs,
		AddSnssaiDnnPairs: subscription.AddSnssaiDnnPairs,
		SuppFeat:          util.PtrToString(subscription.SuppFeat),
	}

	locationHeader := "/nbsf-management/v1/subscriptions/" + subId
	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, response)
}

// ReplaceIndividualSubcription handles PUT /subscriptions/{subId}
func ReplaceIndividualSubcription(c *gin.Context) {
	logger.ProcLog.Infof("Handle ReplaceIndividualSubcription")

	subId := c.Param("subId")

	var request models.BsfSubscription
	if err := c.ShouldBindJSON(&request); err != nil {
		problemDetail := models.ProblemDetails{
			Status: http.StatusBadRequest,
			Cause:  "MALFORMED_REQUEST",
		}
		c.JSON(http.StatusBadRequest, problemDetail)
		return
	}

	// Check if subscription exists
	_, exists := bsfContext.BsfSelf.GetSubscription(subId)

	// Convert to internal representation
	subscription := &bsfContext.BsfSubscription{
		Events:            request.Events,
		NotifUri:          request.NotifUri,
		NotifCorreId:      request.NotifCorreId,
		Supi:              request.Supi,
		Gpsi:              util.StringToPtr(request.Gpsi),
		SnssaiDnnPairs:    request.SnssaiDnnPairs,
		AddSnssaiDnnPairs: request.AddSnssaiDnnPairs,
		SuppFeat:          util.StringToPtr(request.SuppFeat),
	}

	if exists {
		// Update existing subscription
		bsfContext.BsfSelf.UpdateSubscription(subId, subscription)

		// Return updated subscription
		response := models.BsfSubscriptionResp{
			Events:            subscription.Events,
			NotifUri:          subscription.NotifUri,
			NotifCorreId:      subscription.NotifCorreId,
			Supi:              subscription.Supi,
			Gpsi:              util.PtrToString(subscription.Gpsi),
			SnssaiDnnPairs:    subscription.SnssaiDnnPairs,
			AddSnssaiDnnPairs: subscription.AddSnssaiDnnPairs,
			SuppFeat:          util.PtrToString(subscription.SuppFeat),
		}
		c.JSON(http.StatusOK, response)
	} else {
		// Create new subscription with given ID
		subscription.SubId = subId
		bsfContext.BsfSelf.Subscriptions[subId] = subscription

		c.Status(http.StatusNoContent)
	}
}

// DeleteIndividualSubcription handles DELETE /subscriptions/{subId}
func DeleteIndividualSubcription(c *gin.Context) {
	logger.ProcLog.Infof("Handle DeleteIndividualSubcription")

	subId := c.Param("subId")

	if bsfContext.BsfSelf.DeleteSubscription(subId) {
		c.Status(http.StatusNoContent)
	} else {
		problemDetail := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "RESOURCE_NOT_FOUND",
		}
		c.JSON(http.StatusNotFound, problemDetail)
	}
}
