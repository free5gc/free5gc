package context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/nrf/factory"
	"free5gc/src/nrf/logger"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var NrfNfProfile models.NfProfile
var Nrf_NfInstanceID string

func InitNrfContext() {

	config := factory.NrfConfig
	logger.InitLog.Infof("nrfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration

	NrfNfProfile.NfInstanceId = uuid.New().String()
	NrfNfProfile.NfType = models.NfType_NRF
	NrfNfProfile.NfStatus = models.NfStatus_REGISTERED

	serviceNameList := configuration.ServiceNameList

	NFServices := InitNFService(serviceNameList, config.Info.Version)
	NrfNfProfile.NfServices = &NFServices
}

func InitNFService(serivceName []string, version string) []models.NfService {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	NFServices := make([]models.NfService, len(serivceName))
	for index, nameString := range serivceName {
		name := models.ServiceName(nameString)
		NFServices[index] = models.NfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          models.UriScheme(factory.NrfConfig.Configuration.Sbi.Scheme),
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       GetIPv4Uri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: factory.NrfConfig.Configuration.Sbi.IPv4Addr,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(factory.NrfConfig.Configuration.Sbi.Port),
				},
			},
		}
	}
	return NFServices
}

func GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", factory.NrfConfig.Configuration.Sbi.Scheme,
		factory.NrfConfig.Configuration.Sbi.IPv4Addr, factory.NrfConfig.Configuration.Sbi.Port)
}

func GetServiceIp() string {
	return factory.NrfConfig.Configuration.DefaultServiceIP
}
