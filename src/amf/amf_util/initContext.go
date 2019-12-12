package amf_util

import (
	"fmt"
	"github.com/google/uuid"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
	"free5gc/src/amf/factory"
	"free5gc/src/amf/logger"
)

func InitAmfContext(context *amf_context.AMFContext) {
	config := factory.AmfConfig
	logger.UtilLog.Infof("amfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.AmfName != "" {
		context.Name = configuration.AmfName
	}
	if configuration.NgapIpList != nil {
		context.NgapIpList = configuration.NgapIpList
	} else {
		context.NgapIpList = []string{"127.0.0.1"} // default localhost
	}
	sbi := configuration.Sbi
	context.UriScheme = models.UriScheme(sbi.Scheme)
	context.HttpIPv4Address = "127.0.0.1" // default localhost
	context.HttpIpv4Port = 29518          // default port
	if sbi != nil {
		if sbi.IPv4Addr != "" {
			context.HttpIPv4Address = sbi.IPv4Addr
		}
		if sbi.Port != 0 {
			context.HttpIpv4Port = sbi.Port
		}
	}
	serviceNameList := configuration.ServiceNameList
	context.InitNFService(serviceNameList, config.Info.Version)
	context.ServedGuamiList = configuration.ServedGumaiList
	context.SupportTaiLists = configuration.SupportTAIList
	for i := range context.SupportTaiLists {
		context.SupportTaiLists[i].Tac = TACConfigToModels(context.SupportTaiLists[i].Tac)
	}
	context.PlmnSupportList = configuration.PlmnSupportList
	context.SupportDnnLists = configuration.SupportDnnList
	if configuration.NrfUri != "" {
		context.NrfUri = configuration.NrfUri
	} else {
		context.NrfUri = fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, 29510)
	}
	security := configuration.Security
	if security != nil {
		context.SecurityAlgorithm.IntegrityOrder = getIntAlgOrder(security.IntegrityOrder)
		context.SecurityAlgorithm.CipheringOrder = getEncAlgOrder(security.CipheringOrder)
	}
	context.NetworkName = configuration.NetworkName
	context.T3502Value = configuration.T3502
	context.T3512Value = configuration.T3512
	context.Non3gppDeregistrationTimerValue = configuration.Non3gppDeregistrationTimer
}

func getIntAlgOrder(integrityOrder []string) (intOrder []uint8) {
	for _, intAlg := range integrityOrder {
		switch intAlg {
		case "NIA0":
			intOrder = append(intOrder, 0x00)
		case "NIA1":
			intOrder = append(intOrder, 0x80)
		case "NIA2":
			intOrder = append(intOrder, 0x40)
		case "NIA3":
			intOrder = append(intOrder, 0x20)
		case "EIA0":
			intOrder = append(intOrder, 0x10)
		case "EIA1":
			intOrder = append(intOrder, 0x08)
		case "EIA2":
			intOrder = append(intOrder, 0x04)
		case "EIA3":
			intOrder = append(intOrder, 0x02)
		}
	}
	return
}
func getEncAlgOrder(cipheringOrder []string) (encOrder []uint8) {
	for _, encAlg := range cipheringOrder {
		switch encAlg {
		case "NEA0":
			encOrder = append(encOrder, 0x00)
		case "NEA1":
			encOrder = append(encOrder, 0x80)
		case "NEA2":
			encOrder = append(encOrder, 0x40)
		case "NEA3":
			encOrder = append(encOrder, 0x20)
		case "EEA0":
			encOrder = append(encOrder, 0x10)
		case "EEA1":
			encOrder = append(encOrder, 0x08)
		case "EEA2":
			encOrder = append(encOrder, 0x04)
		case "EEA3":
			encOrder = append(encOrder, 0x02)
		}
	}
	return
}
