package context

import (
	"math"
	"os"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/sm"
	"github.com/google/uuid"

	"github.com/free5gc/chf/internal/logger"
	"github.com/free5gc/chf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/idgenerator"
)

// Init CHF Context from config flie
func InitChfContext(context *CHFContext) {
	config := factory.ChfConfig
	logger.InitLog.Infof("chfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)

	configuration := config.Configuration
	sbi := configuration.Sbi

	context.NfId = uuid.New().String()
	context.Name = "CHF"
	context.NrfUri = configuration.NrfUri
	context.NrfCertPem = configuration.NrfCertPem
	context.UriScheme = models.UriScheme(configuration.Sbi.Scheme)
	context.RatingSessionIdGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
	context.AccountSessionIdGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
	context.RegisterIPv4 = factory.ChfSbiDefaultIPv4 // default localhost
	context.SBIPort = factory.ChfSbiDefaultPort      // default port
	if sbi != nil {
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
			logger.InitLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			context.BindingIPv4 = sbi.BindingIPv4
			if context.BindingIPv4 == "" {
				logger.InitLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				context.BindingIPv4 = "0.0.0.0"
			}
		}
	}

	rfDiameter := configuration.RfDiameter
	abmfDiameter := configuration.AbmfDiameter

	context.RatingCfg = &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("client"),
		OriginRealm:      datatype.DiameterIdentity("go-diameter"),
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(rfDiameter.HostIPv4),
		},
	}
	context.AbmfCfg = &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("client"),
		OriginRealm:      datatype.DiameterIdentity("go-diameter"),
		VendorID:         13,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(abmfDiameter.HostIPv4),
		},
	}

	context.Url = string(context.UriScheme) + "://" + context.RegisterIPv4 + ":" + strconv.Itoa(context.SBIPort)

	context.NfService = make(map[models.ServiceName]models.NrfNfManagementNfService)
	AddNfServices(&context.NfService, config, context)
}

func AddNfServices(
	serviceMap *map[models.ServiceName]models.NrfNfManagementNfService, config *factory.Config, context *CHFContext,
) {
	var nfService models.NrfNfManagementNfService
	var ipEndPoints []models.IpEndPoint
	var nfServiceVersions []models.NfServiceVersion
	services := *serviceMap

	nfService.ServiceInstanceId = context.NfId
	nfService.ServiceName = models.ServiceName_NCHF_CONVERGEDCHARGING
	nfService.ApiPrefix = context.Url
	var ipEndPoint models.IpEndPoint
	ipEndPoint.Ipv4Address = context.RegisterIPv4
	ipEndPoint.Port = int32(context.SBIPort)
	ipEndPoint.Transport = models.NrfNfManagementTransportProtocol_TCP
	ipEndPoints = append(ipEndPoints, ipEndPoint)

	var nfServiceVersion models.NfServiceVersion
	nfServiceVersion.ApiFullVersion = config.Info.Version
	nfServiceVersion.ApiVersionInUri = "v1"
	nfServiceVersions = append(nfServiceVersions, nfServiceVersion)

	nfService.Scheme = context.UriScheme
	nfService.NfServiceStatus = models.NfServiceStatus_REGISTERED

	nfService.IpEndPoints = ipEndPoints
	nfService.Versions = nfServiceVersions
	services[models.ServiceName_NCHF_CONVERGEDCHARGING] = nfService
}
