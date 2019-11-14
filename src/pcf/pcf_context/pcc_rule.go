package pcf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"time"
	// 	"free5gc/lib/path_util"
	// 	"free5gc/src/pcf/factory"
)

// TODO: use config file to pass default pcc rule
func CreateDefalutPccRules(id int32) models.PccRule {
	flowInfos := []models.FlowInformation{
		{
			FlowDescription:   "permit out ip from any to assigned",
			FlowDirection:     models.UPLINKFlowDirectionRm,
			PacketFilterUsage: true,
			PackFiltId:        "PackFiltId-0",
		},
		{
			FlowDescription:   "permit out ip from any to assigned",
			FlowDirection:     models.DOWNLINKFlowDirectionRm,
			PacketFilterUsage: true,
			PackFiltId:        "PackFiltId-1",
		},
	}
	return CreatePccRule(id, 10, flowInfos)
}

func CreatePccRule(id, precedence int32, flowInfo []models.FlowInformation) models.PccRule {
	idString := fmt.Sprintf("%d", id)
	return models.PccRule{
		FlowInfos:  flowInfo,
		PccRuleId:  "PccRuleId-" + idString,
		Precedence: precedence,
		RefQosData: []string{
			"QosId-" + idString,
		},
		RefCondData: "CondId-" + idString,
	}
}
func CreateCondData(id int32) models.ConditionData {
	idString := fmt.Sprintf("%d", id)
	activationTime := time.Now()
	return models.ConditionData{
		CondId:         "CondId-" + idString,
		ActivationTime: &activationTime,
	}
}

func CreateQosData(id, var5qi, arp int32) models.QosData {
	idString := fmt.Sprintf("%d", id)
	return models.QosData{
		QosId:  "QosId-" + idString,
		Var5qi: var5qi,
		Arp: &models.Arp{
			PriorityLevel: arp,
		},
	}
}
