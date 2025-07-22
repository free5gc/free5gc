package util

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
	udr_context "github.com/free5gc/udr/internal/context"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/pkg/factory"
)

func InitUdrContext(context *udr_context.UDRContext) {
	config := factory.UdrConfig
	logger.UtilLog.Infof("udrconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	context.RegisterIPv4 = factory.UDR_DEFAULT_IPV4 // default localhost
	context.SBIPort = factory.UDR_DEFAULT_PORT_INT  // default port
	if sbi := configuration.Sbi; sbi != nil {
		context.UriScheme = models.UriScheme(sbi.Scheme)
		if sbi.RegisterIPv4 != "" {
			context.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			context.SBIPort = sbi.Port
		}

		context.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if context.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			context.BindingIPv4 = sbi.BindingIPv4
			if context.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				context.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	if configuration.NrfUri != "" {
		context.NrfUri = configuration.NrfUri
	} else {
		logger.UtilLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		context.NrfUri = fmt.Sprintf("%s://%s:%d", context.UriScheme, "127.0.0.1", 29510)
	}
}
