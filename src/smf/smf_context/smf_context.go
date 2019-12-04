package smf_context

import (
	"fmt"
	"net"

	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/Nnrf_NFManagement"
	"free5gc/src/smf/factory"
	"free5gc/src/smf/logger"

	"github.com/google/uuid"
	"free5gc/lib/openapi/models"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
)

func init() {
	smfContext.NfInstanceID = uuid.New().String()
}

var smfContext SMFContext

type SMFContext struct {
	Name         string
	NfInstanceID string

	URIScheme   models.UriScheme
	HTTPAddress string
	HTTPPort    int

	CPNodeID pfcpType.NodeID

	UDMProfiles []models.NfProfile
	PCFProfiles []models.NfProfile

	UPNodeIDs []pfcpType.NodeID
	Key       string
	PEM       string
	KeyLog    string

	UESubNet      *net.IPNet
	UEAddressTemp net.IP

	NrfUri             string
	NFManagementClient *Nnrf_NFManagement.APIClient
	NFDiscoveryClient  *Nnrf_NFDiscovery.APIClient
}

func AllocUEIP() net.IP {
	smfContext.UEAddressTemp[3]++
	return smfContext.UEAddressTemp
}

func InitSmfContext(config *factory.Config) {
	if config == nil {
		logger.CtxLog.Infof("Config is nil")
	}

	logger.CtxLog.Infof("smfconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	if configuration.SmfName != "" {
		smfContext.Name = configuration.SmfName
	}

	sbi := configuration.Sbi
	smfContext.URIScheme = models.UriScheme(sbi.Scheme)
	smfContext.HTTPAddress = "127.0.0.1" // default localhost
	smfContext.HTTPPort = 29502          // default port
	if sbi != nil {
		if sbi.IPv4Addr != "" {
			smfContext.HTTPAddress = sbi.IPv4Addr
		}
		if sbi.Port != 0 {
			smfContext.HTTPPort = sbi.Port
		}

		if tls := sbi.TLS; tls != nil {
			smfContext.Key = tls.Key
			smfContext.PEM = tls.PEM
		}
	}
	if configuration.NrfUri != "" {
		smfContext.NrfUri = configuration.NrfUri
	} else {
		smfContext.NrfUri = fmt.Sprintf("%s://%s:%d", smfContext.URIScheme, smfContext.HTTPAddress, 29510)
	}

	if pfcp := configuration.PFCP; pfcp != nil {
		if pfcp.Port == 0 {
			pfcp.Port = pfcpUdp.PFCP_PORT
		}
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", pfcp.Addr, pfcp.Port))
		if err != nil {
			logger.CtxLog.Warnf("PFCP Parse Addr Fail: %v", err)
		}

		smfContext.CPNodeID.NodeIdType = 0
		smfContext.CPNodeID.NodeIdValue = addr.IP.To4()
	}

	_, ipNet, err := net.ParseCIDR(configuration.UESubnet)
	if err != nil {
		logger.InitLog.Errorln(err)
	}
	smfContext.UESubNet = ipNet
	smfContext.UEAddressTemp = ipNet.IP

	// Set client and set url
	ManagementConfig := Nnrf_NFManagement.NewConfiguration()
	ManagementConfig.SetBasePath(SMF_Self().NrfUri)
	smfContext.NFManagementClient = Nnrf_NFManagement.NewAPIClient(ManagementConfig)

	NFDiscovryConfig := Nnrf_NFDiscovery.NewConfiguration()
	NFDiscovryConfig.SetBasePath(SMF_Self().NrfUri)
	smfContext.NFDiscoveryClient = Nnrf_NFDiscovery.NewAPIClient(NFDiscovryConfig)

	SetupNFProfile()
}

func SMF_Self() *SMFContext {
	return &smfContext
}
