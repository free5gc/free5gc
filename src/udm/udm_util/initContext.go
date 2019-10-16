package udm_util

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/factory"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"

	"github.com/google/uuid"
)

func InitUDMContext(context *udm_context.UDMContext) {
	config := factory.UdmConfig
	logger.UtilLog.Infof("udmconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.UdmName != "" {
		context.Name = configuration.UdmName
	}
	nrfclient := config.Configuration.Nrfclient
	context.NrfUri = fmt.Sprintf("%s://%s:%d", nrfclient.Scheme, nrfclient.Ipv4Adrr, nrfclient.Port)
	sbi := configuration.Sbi
	context.UriScheme = models.UriScheme(sbi.Scheme)
	context.HttpIpv4Port = 29503
	context.HttpIPv4Address = "127.0.0.1"
	if sbi != nil {
		if sbi.IPv4Addr != "" {
			context.HttpIPv4Address = sbi.IPv4Addr
		}
		if sbi.Port != 0 {
			context.HttpIpv4Port = sbi.Port
		}
	}
	servingNameList := configuration.ServiceNameList
	context.InitNFService(servingNameList, config.Info.Version)
}
