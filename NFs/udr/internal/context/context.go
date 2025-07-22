package context

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
	"github.com/free5gc/udr/internal/logger"
	"github.com/free5gc/udr/pkg/factory"
)

var udrContext = UDRContext{}

type subsId = string

type UDRServiceType int

const (
	NUDR_DR UDRServiceType = iota
)

func Init() {
	udrContext.Name = "udr"
	udrContext.EeSubscriptionIDGenerator = 1
	udrContext.SdmSubscriptionIDGenerator = 1
	udrContext.SubscriptionDataSubscriptionIDGenerator = 1
	udrContext.PolicyDataSubscriptionIDGenerator = 1
	udrContext.SubscriptionDataSubscriptions = make(map[subsId]*models.SubscriptionDataSubscriptions)
	udrContext.PolicyDataSubscriptions = make(map[subsId]*models.PolicyDataSubscription)
	udrContext.InfluenceDataSubscriptionIDGenerator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	serviceName := []models.ServiceName{
		models.ServiceName_NUDR_DR,
	}
	udrContext.NrfUri = fmt.Sprintf("%s://%s:%d", models.UriScheme_HTTPS, udrContext.RegisterIPv4, 29510)
	initUdrContext()

	config := factory.UdrConfig
	udrContext.NfService = initNfService(serviceName, config.Info.Version)
}

type UDRContext struct {
	Name                                    string
	UriScheme                               models.UriScheme
	BindingIPv4                             string
	SBIPort                                 int
	NfService                               map[models.ServiceName]models.NrfNfManagementNfService
	RegisterIPv4                            string // IP register to NRF
	HttpIPv6Address                         string
	NfId                                    string
	NrfUri                                  string
	NrfCertPem                              string
	EeSubscriptionIDGenerator               int
	SdmSubscriptionIDGenerator              int
	SubscriptionDataSubscriptionIDGenerator int
	PolicyDataSubscriptionIDGenerator       int
	InfluenceDataSubscriptionIDGenerator    *rand.Rand
	UESubsCollection                        sync.Map // map[ueId]*UESubsData
	UEGroupCollection                       sync.Map // map[ueGroupId]*UEGroupSubsData
	SubscriptionDataSubscriptions           map[subsId]*models.SubscriptionDataSubscriptions
	PolicyDataSubscriptions                 map[subsId]*models.PolicyDataSubscription
	InfluenceDataSubscriptions              sync.Map
	appDataInfluDataSubscriptionIdGenerator uint64
	mtx                                     sync.RWMutex
	OAuth2Required                          bool
}

type UESubsData struct {
	EeSubscriptionCollection map[subsId]*EeSubscriptionCollection
	SdmSubscriptions         map[subsId]*models.SdmSubscription
}

type UEGroupSubsData struct {
	EeSubscriptions map[subsId]*models.EeSubscription
}

type EeSubscriptionCollection struct {
	EeSubscriptions      *models.EeSubscription
	AmfSubscriptionInfos []models.AmfSubscriptionInfo
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &UDRContext{}

// Reset UDR Context
func (context *UDRContext) Reset() {
	context.UESubsCollection.Range(func(key, value interface{}) bool {
		context.UESubsCollection.Delete(key)
		return true
	})
	context.UEGroupCollection.Range(func(key, value interface{}) bool {
		context.UEGroupCollection.Delete(key)
		return true
	})
	for key := range context.SubscriptionDataSubscriptions {
		delete(context.SubscriptionDataSubscriptions, key)
	}
	for key := range context.PolicyDataSubscriptions {
		delete(context.PolicyDataSubscriptions, key)
	}
	context.InfluenceDataSubscriptions.Range(func(key, value interface{}) bool {
		context.InfluenceDataSubscriptions.Delete(key)
		return true
	})
	context.EeSubscriptionIDGenerator = 1
	context.SdmSubscriptionIDGenerator = 1
	context.SubscriptionDataSubscriptionIDGenerator = 1
	context.PolicyDataSubscriptionIDGenerator = 1
	context.InfluenceDataSubscriptionIDGenerator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	context.UriScheme = models.UriScheme_HTTPS
	context.Name = "udr"
}

func initUdrContext() {
	config := factory.UdrConfig
	logger.UtilLog.Infof("udrconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration
	udrContext.NfId = uuid.New().String()
	udrContext.RegisterIPv4 = factory.UDR_DEFAULT_IPV4 // default localhost
	udrContext.SBIPort = factory.UDR_DEFAULT_PORT_INT  // default port
	if sbi := configuration.Sbi; sbi != nil {
		udrContext.UriScheme = models.UriScheme(sbi.Scheme)
		if sbi.RegisterIPv4 != "" {
			udrContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			udrContext.SBIPort = sbi.Port
		}

		udrContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if udrContext.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			udrContext.BindingIPv4 = sbi.BindingIPv4
			if udrContext.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				udrContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	if configuration.NrfUri != "" {
		udrContext.NrfUri = configuration.NrfUri
	} else {
		logger.UtilLog.Warn("NRF Uri is empty! Using localhost as NRF IPv4 address.")
		udrContext.NrfUri = fmt.Sprintf("%s://%s:%d", udrContext.UriScheme, "127.0.0.1", 29510)
	}
	udrContext.NrfCertPem = configuration.NrfCertPem
}

func initNfService(serviceName []models.ServiceName, version string) (
	nfService map[models.ServiceName]models.NrfNfManagementNfService,
) {
	versionUri := "v" + strings.Split(version, ".")[0]
	nfService = make(map[models.ServiceName]models.NrfNfManagementNfService)
	for idx, name := range serviceName {
		nfService[name] = models.NrfNfManagementNfService{
			ServiceInstanceId: strconv.Itoa(idx),
			ServiceName:       name,
			Versions: []models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          udrContext.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       GetIPv4Uri(),
			IpEndPoints: []models.IpEndPoint{
				{
					Ipv4Address: udrContext.RegisterIPv4,
					Transport:   models.NrfNfManagementTransportProtocol_TCP,
					Port:        int32(udrContext.SBIPort),
				},
			},
		}
	}

	return
}

func GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", udrContext.UriScheme, udrContext.RegisterIPv4, udrContext.SBIPort)
}

func (context *UDRContext) GetIPv4GroupUri(udrServiceType UDRServiceType) string {
	var serviceUri string

	switch udrServiceType {
	case NUDR_DR:
		serviceUri = factory.UdrDrResUriPrefix
	default:
		serviceUri = ""
	}

	return fmt.Sprintf("%s://%s:%d%s", context.UriScheme, context.RegisterIPv4, context.SBIPort, serviceUri)
}

// Create new UDR context
func GetSelf() *UDRContext {
	return &udrContext
}

func (context *UDRContext) NewAppDataInfluDataSubscriptionID() uint64 {
	context.mtx.Lock()
	defer context.mtx.Unlock()
	context.appDataInfluDataSubscriptionIdGenerator++
	return context.appDataInfluDataSubscriptionIdGenerator
}

func NewInfluenceDataSubscriptionId() string {
	if GetSelf().InfluenceDataSubscriptionIDGenerator == nil {
		GetSelf().InfluenceDataSubscriptionIDGenerator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	}
	return fmt.Sprintf("%08x", GetSelf().InfluenceDataSubscriptionIDGenerator.Uint32())
}

func (c *UDRContext) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	if !c.OAuth2Required {
		return context.TODO(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType_UDR, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}

func (c *UDRContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !c.OAuth2Required {
		logger.UtilLog.Debugf("UDRContext::AuthorizationCheck: OAuth2 not required\n")
		return nil
	}

	logger.UtilLog.Debugf("UDRContext::AuthorizationCheck: token[%s] serviceName[%s]\n", token, serviceName)
	return oauth.VerifyOAuth(token, string(serviceName), c.NrfCertPem)
}
