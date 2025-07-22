package processor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/internal/util"
	"github.com/free5gc/util/mongoapi"
)

func (p *Processor) HandleAmfStatusChangeNotify(
	c *gin.Context,
	amfStatusChangeNotification models.AmfStatusChangeNotification,
) {
	logger.CallbackLog.Warnf("[PCF] Handle Amf Status Change Notify is not implemented.")

	// TODO: handle AMF Status Change Notify
	logger.CallbackLog.Debugf("receive AMF status change notification[%+v]", amfStatusChangeNotification)

	c.JSON(http.StatusNoContent, nil)
}

func (p *Processor) HandlePolicyDataChangeNotify(
	c *gin.Context,
	supi string,
	policyDataChangeNotification models.PolicyDataChangeNotification,
) {
	logger.CallbackLog.Warnf("[PCF] Handle Policy Data Change Notify is not implemented.")

	PolicyDataChangeNotifyProcedure(supi, policyDataChangeNotification)

	c.JSON(http.StatusNotImplemented, gin.H{})
}

// TODO: handle Policy Data Change Notify
func PolicyDataChangeNotifyProcedure(supi string, notification models.PolicyDataChangeNotification) {
}

func (p *Processor) HandleInfluenceDataUpdateNotify(
	c *gin.Context,
	supi string,
	pduSessionId string,
	trafficInfluDataNotif []models.TrafficInfluDataNotif,
) {
	logger.CallbackLog.Infof("[PCF] Handle Influence Data Update Notify")

	smPolicyID := fmt.Sprintf("%s-%s", supi, pduSessionId)
	ue := p.Context().PCFUeFindByPolicyId(smPolicyID)
	if ue == nil || ue.SmPolicyData[smPolicyID] == nil {
		problemDetail := util.GetProblemDetail("smPolicyID not found in PCF", util.CONTEXT_NOT_FOUND)
		logger.CallbackLog.Error(problemDetail.Detail)
		c.JSON(int(problemDetail.Status), problemDetail)
		return
	}
	smPolicy := ue.SmPolicyData[smPolicyID]
	decision := smPolicy.PolicyDecision
	influenceDataToPccRule := smPolicy.InfluenceDataToPccRule
	precedence := getAvailablePrecedence(smPolicy.PolicyDecision.PccRules)
	for _, notification := range trafficInfluDataNotif {
		influenceID := getInfluenceID(notification.ResUri)
		if influenceID == "" {
			continue
		}
		// notifying deletion
		if notification.TrafficInfluData == nil {
			pccRuleID := influenceDataToPccRule[influenceID]
			decision = &models.SmPolicyDecision{}
			if err := smPolicy.RemovePccRule(pccRuleID, decision); err != nil {
				logger.CallbackLog.Errorf("Remove PCC rule error: %+v", err)
			}
			delete(influenceDataToPccRule, influenceID)
		} else {
			var chgData *models.ChargingData
			var chargingInterface map[string]interface{}

			trafficInfluData := *notification.TrafficInfluData

			filterCharging := bson.M{
				"ueId":   ue.Supi,
				"snssai": util.SnssaiModelsToHex(*trafficInfluData.Snssai),
				"dnn":    "",
				"filter": "",
			}
			chargingInterface, err := mongoapi.RestfulAPIGetOne(chargingDataColl, filterCharging, 2)
			if err != nil {
				logger.SmPolicyLog.Errorf("Fail to get charging data to mongoDB err: %+v", err)
				chgData = nil
			} else if chargingInterface != nil {
				rg, err1 := p.Context().RatingGroupIdGenerator.Allocate()
				if err1 != nil {
					logger.SmPolicyLog.Error("rating group allocate error")
					problemDetails := util.GetProblemDetail("rating group allocate error", util.ERROR_IDGENERATOR)
					c.JSON(int(problemDetails.Status), problemDetails)
					return
				}
				chgData = &models.ChargingData{
					ChgId:          util.GetChgId(smPolicy.ChargingIdGenerator),
					RatingGroup:    int32(rg),
					ReportingLevel: models.ReportingLevel_RAT_GR_LEVEL,
					MeteringMethod: models.MeteringMethod_VOLUME,
				}

				switch chargingInterface["chargingMethod"].(string) {
				case "Online":
					chgData.Online = true
					chgData.Offline = false
				case "Offline":
					chgData.Online = false
					chgData.Offline = true
				}

				if decision.ChgDecs == nil {
					decision.ChgDecs = make(map[string]*models.ChargingData)
				}

				chargingInterface["ratingGroup"] = chgData.RatingGroup
				logger.SmPolicyLog.Tracef("put ratingGroup[%+v] for [%+v] to MongoDB", chgData.RatingGroup, ue.Supi)
				if _, err = mongoapi.RestfulAPIPutOne(
					chargingDataColl, chargingInterface, chargingInterface, 2); err != nil {
					logger.SmPolicyLog.Errorf("Fail to put charging data to mongoDB err: %+v", err)
				} else {
					smPolicy.ChargingIdGenerator++
				}
				if ue.RatingGroupData == nil {
					ue.RatingGroupData = make(map[string][]int32)
				}
				ue.RatingGroupData[smPolicyID] = append(ue.RatingGroupData[smPolicyID], chgData.RatingGroup)
			}

			if pccRuleID, ok := influenceDataToPccRule[influenceID]; ok {
				// notifying Individual Influence Data update
				pccRule := decision.PccRules[pccRuleID]
				util.SetSmPolicyDecisionByTrafficInfluData(decision, pccRule, trafficInfluData, chgData)
			} else {
				// notifying Individual Influence Data creation

				pccRule := util.CreatePccRule(smPolicy.PccRuleIdGenerator, precedence, nil, trafficInfluData.AfAppId)
				util.SetSmPolicyDecisionByTrafficInfluData(decision, pccRule, trafficInfluData, chgData)
				influenceDataToPccRule[influenceID] = pccRule.PccRuleId
				smPolicy.PccRuleIdGenerator++
				if precedence < Precedence_Maximum {
					precedence++
				}
			}
		}
	}
	smPolicyNotification := models.SmPolicyNotification{
		ResourceUri:      util.GetResourceUri(models.ServiceName_NPCF_SMPOLICYCONTROL, smPolicyID),
		SmPolicyDecision: decision,
	}
	go p.SendSMPolicyUpdateNotification(smPolicy.PolicyContext.NotificationUri, &smPolicyNotification)
	c.JSON(http.StatusNoContent, nil)
}

func getInfluenceID(resUri string) string {
	temp := strings.Split(resUri, "/")
	return temp[len(temp)-1]
}
