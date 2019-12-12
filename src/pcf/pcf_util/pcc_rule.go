package pcf_util

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"time"
	// 	"free5gc/lib/path_util"
	// 	"free5gc/src/pcf/factory"
)

// TODO: use config file to pass default pcc rule
func CreateDefalutPccRules(id int32) models.PccRule {
	flowInfo := []models.FlowInformation{
		{
			FlowDescription:   "permit out ip from any to assigned",
			FlowDirection:     models.FlowDirectionRm_DOWNLINK,
			PacketFilterUsage: true,
			PackFiltId:        "PackFiltId-0",
		},
		{
			FlowDescription:   "permit out ip from any to assigned",
			FlowDirection:     models.FlowDirectionRm_DOWNLINK,
			PacketFilterUsage: true,
			PackFiltId:        "PackFiltId-1",
		},
	}
	return CreatePccRule(id, 10, flowInfo, true)
}

func GetPccRuleId(id int32) string {
	return fmt.Sprintf("PccRuleId-%d", id)
}

func GetQosId(id int32) string {
	return fmt.Sprintf("QosId-%d", id)
}

func GetCondId(id int32) string {
	return fmt.Sprintf("CondId-%d", id)
}

func GetPackFiltId(id int32) string {
	return fmt.Sprintf("PackFiltId-%d", id)
}

func CreatePccRule(id, precedence int32, flowInfo []models.FlowInformation, cond bool) models.PccRule {
	rule := models.PccRule{
		FlowInfos:  flowInfo,
		PccRuleId:  GetPccRuleId(id),
		Precedence: precedence,
		RefQosData: []string{
			GetQosId(id),
		},
	}
	if cond {
		rule.RefCondData = GetCondId(id)
	}
	return rule
}

func CreateCondData(id int32) models.ConditionData {
	activationTime := time.Now()
	return models.ConditionData{
		CondId:         GetCondId(id),
		ActivationTime: &activationTime,
	}
}

func CreateQosData(id, var5qi, arp int32) models.QosData {
	return models.QosData{
		QosId:  GetQosId(id),
		Var5qi: var5qi,
		Arp: &models.Arp{
			PriorityLevel: arp,
		},
	}
}

// EthDescription is Not Supported
func ConvertPacketInfoToFlowInformation(infos []models.PacketFilterInfo) (flowInfos []models.FlowInformation) {
	for _, info := range infos {
		flowInfo := models.FlowInformation{
			FlowDescription:   info.PackFiltCont,
			PackFiltId:        info.PackFiltId,
			PacketFilterUsage: true,
			TosTrafficClass:   info.TosTrafficClass,
			Spi:               info.Spi,
			FlowLabel:         info.FlowLabel,
			FlowDirection:     models.FlowDirectionRm(info.FlowDirection),
		}
		flowInfos = append(flowInfos, flowInfo)
	}
	return
}
