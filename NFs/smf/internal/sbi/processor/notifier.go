package processor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/smf/EventExposure"
	smf_context "github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
)

func (p *Processor) HandleChargingNotification(
	c *gin.Context,
	chargingNotifyRequest models.ChargingNotifyRequest,
	smContextRef string,
) {
	logger.ChargingLog.Info("Handle Charging Notification")

	problemDetails := p.chargingNotificationProcedure(chargingNotifyRequest, smContextRef)
	if problemDetails == nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(int(problemDetails.Status), problemDetails)
}

// While receive Charging Notification from CHF, SMF will send Charging Information to CHF and update UPF
// The Charging Notification will be sent when CHF found the changes of the quota file.
func (p *Processor) chargingNotificationProcedure(
	req models.ChargingNotifyRequest, smContextRef string,
) *models.ProblemDetails {
	if smContext := smf_context.GetSMContextByRef(smContextRef); smContext != nil {
		smContext.SMLock.Lock()
		defer smContext.SMLock.Unlock()
		upfUrrMap := make(map[string][]*smf_context.URR)
		for _, reauthorizeDetail := range req.ReauthorizationDetails {
			rg := reauthorizeDetail.RatingGroup
			logger.ChargingLog.Infof("Force update charging information for rating group %d", rg)
			for _, urr := range smContext.UrrUpfMap {
				chgInfo := smContext.ChargingInfo[urr.URRID]
				if chgInfo.RatingGroup == rg ||
					chgInfo.ChargingLevel == smf_context.PduSessionCharging {
					logger.ChargingLog.Tracef("Query URR (%d) for Rating Group (%d)", urr.URRID, rg)
					upfId := smContext.ChargingInfo[urr.URRID].UpfId
					upfUrrMap[upfId] = append(upfUrrMap[upfId], urr)
				}
			}
		}
		for upfId, urrList := range upfUrrMap {
			upf := smf_context.GetUpfById(upfId)
			if upf == nil {
				logger.ChargingLog.Warnf("Cound not find upf %s", upfId)
				continue
			}
			QueryReport(smContext, upf, urrList, models.ChfConvergedChargingTriggerType_FORCED_REAUTHORISATION)
		}
		p.ReportUsageAndUpdateQuota(smContext)
	} else {
		detail := fmt.Sprintf("SM Context [%s] Not Found ", smContextRef)
		return openapi.ProblemDetailsDataNotFound(detail)
	}

	return nil
}

func (p *Processor) HandleSMPolicyUpdateNotify(
	c *gin.Context,
	request models.SmPolicyNotification,
	smContextRef string,
) {
	logger.PduSessLog.Infoln("In HandleSMPolicyUpdateNotify")
	decision := request.SmPolicyDecision
	smContext := smf_context.GetSMContextByRef(smContextRef)

	if smContext == nil {
		logger.PduSessLog.Errorf("SMContext[%s] not found", smContextRef)
		c.Status(http.StatusBadRequest)
		return
	}

	smContext.SMLock.Lock()
	defer smContext.SMLock.Unlock()

	smContext.CheckState(smf_context.Active)
	// Wait till the state becomes Active again
	// TODO: implement waiting in concurrent architecture

	smContext.SetState(smf_context.ModificationPending)

	// Update SessionRule from decision
	if err := smContext.ApplySessionRules(decision); err != nil {
		// TODO: Fill the error body
		smContext.Log.Errorf("SMPolicyUpdateNotify err: %v", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// TODO: Response data type -
	// [200 OK] UeCampingRep
	// [200 OK] array(PartialSuccessReport)
	// [400 Bad Request] ErrorReport
	if err := smContext.ApplyPccRules(decision); err != nil {
		smContext.Log.Errorf("apply sm policy decision error: %+v", err)
		// TODO: Fill the error body
		c.Status(http.StatusBadRequest)
		return
	}

	smContext.SendUpPathChgNotification("EARLY", SendUpPathChgEventExposureNotification)

	ActivateUPFSession(smContext, nil)

	smContext.SendUpPathChgNotification("LATE", SendUpPathChgEventExposureNotification)

	smContext.PostRemoveDataPath()

	c.Status(http.StatusNoContent)
}

func SendUpPathChgEventExposureNotification(
	uri string, notification *models.NsmfEventExposureNotification,
) {
	configuration := EventExposure.NewConfiguration()
	client := EventExposure.NewAPIClient(configuration)
	request := &EventExposure.CreateIndividualSubcriptionMyNotificationPostRequest{
		NsmfEventExposureNotification: notification,
	}
	_, err := client.
		SubscriptionsCollectionApi.
		CreateIndividualSubcriptionMyNotificationPost(context.Background(), uri, request)

	switch err := err.(type) {
	case openapi.GenericOpenAPIError:
		logger.PduSessLog.Warnf("SMF Event Exposure Notification Error[%s]", err.Error())
	case error:
		logger.PduSessLog.Warnf("SMF Event Exposure Notification Failed[%s]", err.Error())
	case nil:
		logger.PduSessLog.Tracef("SMF Event Exposure Notification Success")
	default:
		logger.PduSessLog.Warnf("SMF Event Exposure Notification Unknown Error: %+v", err)
	}
}
