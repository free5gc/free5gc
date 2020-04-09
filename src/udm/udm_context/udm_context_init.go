package udm_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/udm/factory"
	"free5gc/src/udm/logger"

	"github.com/google/uuid"
)

func TestInit() {
	DefaultUDMConfigPath := path_util.Gofree5gcPath("free5gc/config/udmcfg.conf")
	factory.InitConfigFactory(DefaultUDMConfigPath)
	Init()
}

func InitUDMContext(context *UDMContext) {
	config := factory.UdmConfig
	logger.UtilLog.Info("udmconfig Info: Version[", config.Info.Version, "] Description[", config.Info.Description, "]")
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.UdmName != "" {
		context.Name = configuration.UdmName
	}
	nrfclient := config.Configuration.Nrfclient
	context.NrfUri = fmt.Sprintf("%s://%s:%d", nrfclient.Scheme, nrfclient.Ipv4Addr, nrfclient.Port)
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
	if configuration.NrfUri != "" {
		context.NrfUri = configuration.NrfUri
	} else {
		context.NrfUri = fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, 29510)
	}
	servingNameList := configuration.ServiceNameList

	context.Keys = configuration.Keys

	context.InitNFService(servingNameList, config.Info.Version)
}
