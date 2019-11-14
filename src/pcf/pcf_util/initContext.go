package pcf_util

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/factory"
	"free5gc/src/pcf/logger"
	"free5gc/src/pcf/pcf_context"

	"github.com/google/uuid"
)

func InitpcfContext(context *pcf_context.PCFContext) {
	config := factory.PcfConfig
	logger.UtilLog.Infof("pcfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.PcfName != "" {
		context.Name = configuration.PcfName
	}
	sbi := configuration.Sbi
	context.NrfUri = configuration.NrfUri
	context.UriScheme = models.UriScheme(sbi.Scheme)
	context.HttpIPv4Address = "127.0.0.1" // default localhost
	context.HttpIpv4Port = 29507          // default port
	if sbi != nil {
		if sbi.IPv4Addr != "" {
			context.HttpIPv4Address = sbi.IPv4Addr
		}
		if sbi.Port != 0 {
			context.HttpIpv4Port = sbi.Port
		}
		if sbi.Scheme == "https" {
			context.UriScheme = models.UriScheme_HTTPS
		} else {
			context.UriScheme = models.UriScheme_HTTP
		}
	}
	serviceNameList := configuration.ServiceNameList
	nfService := context.InitNFService(serviceNameList, config.Info.Version)
	context.TimeFormat = configuration.TimeFormat
	context.DefaultBdtRefId = configuration.DefaultBdtRefId
	pcfServiceUris := make(map[models.ServiceName]string, len(serviceNameList))
	for index, nameString := range serviceNameList {
		name := models.ServiceName(nameString)
		version := *nfService[index].Versions
		if name == nfService[index].ServiceName {
			pcfServiceUris[name] = nfService[index].ApiPrefix + "/" + fmt.Sprint(nfService[index].ServiceName) + "/" + version[0].ApiVersionInUri
		}
	}
}
