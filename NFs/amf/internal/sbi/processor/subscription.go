package processor

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/amf/internal/context"
	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/openapi/models"
)

// TS 29.518 5.2.2.5.1
func (p *Processor) HandleAMFStatusChangeSubscribeRequest(c *gin.Context,
	subscriptionDataReq models.AmfCommunicationSubscriptionData,
) {
	logger.CommLog.Info("Handle AMF Status Change Subscribe Request")

	subscriptionDataRsp, locationHeader, problemDetails := p.AMFStatusChangeSubscribeProcedure(subscriptionDataReq)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	c.Header("Location", locationHeader)
	c.JSON(http.StatusCreated, subscriptionDataRsp)
}

func (p *Processor) AMFStatusChangeSubscribeProcedure(subscriptionDataReq models.AmfCommunicationSubscriptionData) (
	subscriptionDataRsp models.AmfCommunicationSubscriptionData, locationHeader string,
	problemDetails *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	for _, guami := range subscriptionDataReq.GuamiList {
		for _, servedGumi := range amfSelf.ServedGuamiList {
			if reflect.DeepEqual(guami, servedGumi) {
				// AMF status is available
				subscriptionDataRsp.GuamiList = append(subscriptionDataRsp.GuamiList, guami)
			}
		}
	}

	if subscriptionDataRsp.GuamiList != nil {
		newSubscriptionID := amfSelf.NewAMFStatusSubscription(subscriptionDataReq)
		locationHeader = subscriptionDataReq.AmfStatusUri + "/" + newSubscriptionID
		logger.CommLog.Infof("new AMF Status Subscription[%s]", newSubscriptionID)
		return
	} else {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "UNSPECIFIED",
		}
		return
	}
}

// TS 29.518 5.2.2.5.2
func (p *Processor) HandleAMFStatusChangeUnSubscribeRequest(c *gin.Context) {
	logger.CommLog.Info("Handle AMF Status Change UnSubscribe Request")

	subscriptionID := c.Param("subscriptionId")

	problemDetails := p.AMFStatusChangeUnSubscribeProcedure(subscriptionID)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func (p *Processor) AMFStatusChangeUnSubscribeProcedure(subscriptionID string) (problemDetails *models.ProblemDetails) {
	amfSelf := context.GetSelf()

	if _, ok := amfSelf.FindAMFStatusSubscription(subscriptionID); !ok {
		problemDetails = &models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "SUBSCRIPTION_NOT_FOUND",
		}
	} else {
		logger.CommLog.Debugf("Delete AMF status subscription[%s]", subscriptionID)
		amfSelf.DeleteAMFStatusSubscription(subscriptionID)
	}
	return
}

// TS 29.518 5.2.2.5.1.3
func (p *Processor) HandleAMFStatusChangeSubscribeModify(c *gin.Context,
	updateSubscriptionData models.AmfCommunicationSubscriptionData,
) {
	logger.CommLog.Info("Handle AMF Status Change Subscribe Modify Request")

	subscriptionID := c.Param("subscriptionId")

	updatedSubscriptionData, problemDetails := p.
		AMFStatusChangeSubscribeModifyProcedure(subscriptionID, updateSubscriptionData)
	if problemDetails != nil {
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	c.JSON(http.StatusAccepted, updatedSubscriptionData)
}

func (p *Processor) AMFStatusChangeSubscribeModifyProcedure(subscriptionID string,
	subscriptionData models.AmfCommunicationSubscriptionData) (
	*models.AmfCommunicationSubscriptionData, *models.ProblemDetails,
) {
	amfSelf := context.GetSelf()

	if currentSubscriptionData, ok := amfSelf.FindAMFStatusSubscription(subscriptionID); !ok {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "Forbidden",
		}
		return nil, problemDetails
	} else {
		logger.CommLog.Debugf("Modify AMF status subscription[%s]", subscriptionID)

		currentSubscriptionData.GuamiList = currentSubscriptionData.GuamiList[:0]

		currentSubscriptionData.GuamiList = append(currentSubscriptionData.GuamiList, subscriptionData.GuamiList...)
		currentSubscriptionData.AmfStatusUri = subscriptionData.AmfStatusUri

		amfSelf.AMFStatusSubscriptions.Store(subscriptionID, currentSubscriptionData)
		return currentSubscriptionData, nil
	}
}
