package udr_util

import (
	"fmt"
	"github.com/google/uuid"
	"free5gc/lib/openapi/models"
	"free5gc/src/udr/factory"
	"free5gc/src/udr/logger"
	"free5gc/src/udr/udr_context"
)

func InitUdrContext(context *udr_context.UDRContext) {
	config := factory.UdrConfig
	logger.UtilLog.Infof("udrconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	sbi := configuration.Sbi
	context.UriScheme = models.UriScheme(sbi.Scheme)
	context.HttpIPv4Address = "127.0.0.1" // default localhost
	context.HttpIpv4Port = 29504          // default port
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
}
