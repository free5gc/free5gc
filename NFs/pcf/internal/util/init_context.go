package util

import (
	"os"

	"github.com/google/uuid"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pcf/internal/context"
	"github.com/free5gc/pcf/internal/logger"
	"github.com/free5gc/pcf/pkg/factory"
	"github.com/free5gc/util/mongoapi"
)

// Init PCF Context from config flie
func InitpcfContext(context *context.PCFContext) {
	config := factory.PcfConfig
	logger.UtilLog.Infof("pcfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.PcfName != "" {
		context.Name = configuration.PcfName
	}

	mongodb := config.Configuration.Mongodb
	// Connect to MongoDB
	if err := mongoapi.SetMongoDB(mongodb.Name, mongodb.Url); err != nil {
		logger.UtilLog.Errorf("InitpcfContext err: %+v", err)
		return
	}

	sbi := configuration.Sbi
	context.NrfUri = configuration.NrfUri
	context.UriScheme = ""
	context.RegisterIPv4 = factory.PcfSbiDefaultIPv4 // default localhost
	context.SBIPort = factory.PcfSbiDefaultPort      // default port
	if sbi != nil {
		if sbi.Scheme != "" {
			context.UriScheme = models.UriScheme(sbi.Scheme)
		}
		if sbi.RegisterIPv4 != "" {
			context.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			context.SBIPort = sbi.Port
		}
		if sbi.Scheme == "https" {
			context.UriScheme = models.UriScheme_HTTPS
		} else {
			context.UriScheme = models.UriScheme_HTTP
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
	serviceList := configuration.ServiceList
	context.InitNFService(serviceList, config.Info.Version)
	context.TimeFormat = configuration.TimeFormat
	context.DefaultBdtRefId = configuration.DefaultBdtRefId
	for _, service := range context.NfService {
		var err error
		context.PcfServiceUris[service.ServiceName] = service.ApiPrefix +
			"/" + string(service.ServiceName) + "/" + (service.Versions)[0].ApiVersionInUri
		context.PcfSuppFeats[service.ServiceName], err = openapi.NewSupportedFeature(service.SupportedFeatures)
		if err != nil {
			logger.UtilLog.Errorf("openapi NewSupportedFeature error: %+v", err)
		}
	}
	context.Locality = configuration.Locality
}
