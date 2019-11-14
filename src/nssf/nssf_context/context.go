/*
 * NF Context for NSSF
 *
 * Configuration of NSSF itself shall be accessed with NSSF context
 * Configuration of network slices shall be accessed with configuration factory
 */

package nssf_context

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"free5gc/lib/openapi/models"
	"free5gc/src/nssf/factory"
	"free5gc/src/nssf/logger"
)

var nssfContext = NSSFContext{}

// Initialize NSSF context with default value
func init() {
	nssfContext.NfId = uuid.New().String()

	nssfContext.Name = "NSSF"

	nssfContext.UriScheme = models.UriScheme_HTTPS
	// Default NSSF would open services at port 29531 on loopback interface
	nssfContext.HttpIpv4Address = "127.0.0.1"
	nssfContext.Port = 29531

	serviceName := []models.ServiceName{
		models.ServiceName_NNSSF_NSSELECTION,
		models.ServiceName_NNSSF_NSSAIAVAILABILITY,
	}
	nssfContext.NfService = initNfService(serviceName, "1.0.0")

	nssfContext.NrfUri = fmt.Sprintf("%s://%s:%d", models.UriScheme_HTTPS, nssfContext.HttpIpv4Address, 29510)
}

type NSSFContext struct {
	NfId            string
	Name            string
	UriScheme       models.UriScheme
	HttpIpv4Address string
	// HttpIpv6Address string
	Port              int
	NfService         map[models.ServiceName]models.NfService
	NrfUri            string
	SupportedPlmnList []models.PlmnId
}

// Initialize NSSF context with configuration factory
func InitNssfContext() {
	if !factory.Configured {
		logger.ContextLog.Warnf("NSSF is not configured")
		return
	}
	nssfConfig := factory.NssfConfig

	if nssfConfig.Configuration.NssfName != "" {
		nssfContext.Name = nssfConfig.Configuration.NssfName
	}

	nssfContext.UriScheme = nssfConfig.Configuration.Sbi.Scheme
	nssfContext.HttpIpv4Address = nssfConfig.Configuration.Sbi.Ipv4Addr
	nssfContext.Port = nssfConfig.Configuration.Sbi.Port

	nssfContext.NfService = initNfService(nssfConfig.Configuration.ServiceNameList, nssfConfig.Info.Version)

	nssfContext.NrfUri = nssfConfig.Configuration.NrfUri

	nssfContext.SupportedPlmnList = nssfConfig.Configuration.SupportedPlmnList
}

func initNfService(serviceName []models.ServiceName, version string) (nfService map[models.ServiceName]models.NfService) {
	versionUri := "v" + strings.Split(version, ".")[0]
	nfService = make(map[models.ServiceName]models.NfService)
	for idx, name := range serviceName {
		nfService[name] = models.NfService{
			ServiceInstanceId: strconv.Itoa(idx),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          nssfContext.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       GetIpv4Uri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: nssfContext.HttpIpv4Address,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(nssfContext.Port),
				},
			},
		}
	}

	return
}

func GetIpv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", nssfContext.UriScheme, nssfContext.HttpIpv4Address, nssfContext.Port)
}

func NSSF_Self() *NSSFContext {
	return &nssfContext
}
