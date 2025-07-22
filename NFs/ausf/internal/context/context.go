package context

import (
	"context"
	"regexp"
	"sync"

	"github.com/free5gc/ausf/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
)

type AUSFContext struct {
	suciSupiMap          sync.Map
	UePool               sync.Map
	NfId                 string
	GroupID              string
	SBIPort              int
	RegisterIPv4         string
	BindingIPv4          string
	Url                  string
	UriScheme            models.UriScheme
	NrfUri               string
	NrfCertPem           string
	NfService            map[models.ServiceName]models.NrfNfManagementNfService
	PlmnList             []models.PlmnId
	UdmUeauUrl           string
	snRegex              *regexp.Regexp
	EapAkaSupiImsiPrefix bool
	OAuth2Required       bool
}

type AusfUeContext struct {
	Supi               string
	Kausf              string
	Kseaf              string
	ServingNetworkName string
	AuthStatus         models.AusfUeAuthenticationAuthResult
	UdmUeauUrl         string

	// for 5G AKA
	XresStar string

	// for EAP-AKA'
	K_aut    string
	XRES     string
	Rand     string
	EapID    uint8
	Resynced bool
}

type SuciSupiMap struct {
	SupiOrSuci string
	Supi       string
}

type EapAkaPrimeAttribute struct {
	Type   uint8
	Length uint8
	Value  []byte
}

type EapAkaPrimePkt struct {
	Subtype    uint8
	Attributes map[uint8]EapAkaPrimeAttribute
	MACInput   []byte
}

const (
	EAP_AKA_PRIME_TYPENUM = 50
)

// Attribute Types for EAP-AKA'
const (
	AT_RAND_ATTRIBUTE              = 1
	AT_AUTN_ATTRIBUTE              = 2
	AT_RES_ATTRIBUTE               = 3
	AT_AUTS_ATTRIBUTE              = 4
	AT_MAC_ATTRIBUTE               = 11
	AT_NOTIFICATION_ATTRIBUTE      = 12
	AT_IDENTITY_ATTRIBUTE          = 14
	AT_CLIENT_ERROR_CODE_ATTRIBUTE = 22
	AT_KDF_INPUT_ATTRIBUTE         = 23
	AT_KDF_ATTRIBUTE               = 24
)

// Subtypes for EAP-AKA'
const (
	AKA_CHALLENGE_SUBTYPE               = 1
	AKA_AUTHENTICATION_REJECT_SUBTYPE   = 2
	AKA_SYNCHRONIZATION_FAILURE_SUBTYPE = 4
	AKA_NOTIFICATION_SUBTYPE            = 12
	AKA_CLIENT_ERROR_SUBTYPE            = 14
)

var ausfContext AUSFContext

func Init() {
	if snRegex, err := regexp.Compile("5G:mnc[0-9]{3}[.]mcc[0-9]{3}[.]3gppnetwork[.]org"); err != nil {
		logger.CtxLog.Warnf("SN compile error: %+v", err)
	} else {
		ausfContext.snRegex = snRegex
	}
	InitAusfContext(&ausfContext)
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &AUSFContext{}

func NewAusfUeContext(identifier string) (ausfUeContext *AusfUeContext) {
	ausfUeContext = new(AusfUeContext)
	ausfUeContext.Supi = identifier // supi
	return ausfUeContext
}

func AddAusfUeContextToPool(ausfUeContext *AusfUeContext) {
	ausfContext.UePool.Store(ausfUeContext.Supi, ausfUeContext)
}

func CheckIfAusfUeContextExists(ref string) bool {
	_, ok := ausfContext.UePool.Load(ref)
	return ok
}

func GetAusfUeContext(ref string) *AusfUeContext {
	context, _ := ausfContext.UePool.Load(ref)
	ausfUeContext := context.(*AusfUeContext)
	return ausfUeContext
}

func AddSuciSupiPairToMap(supiOrSuci string, supi string) {
	newPair := new(SuciSupiMap)
	newPair.SupiOrSuci = supiOrSuci
	newPair.Supi = supi
	ausfContext.suciSupiMap.Store(supiOrSuci, newPair)
}

func CheckIfSuciSupiPairExists(ref string) bool {
	_, ok := ausfContext.suciSupiMap.Load(ref)
	return ok
}

func GetSupiFromSuciSupiMap(ref string) (supi string) {
	val, _ := ausfContext.suciSupiMap.Load(ref)
	suciSupiMap := val.(*SuciSupiMap)
	supi = suciSupiMap.Supi
	return supi
}

func IsServingNetworkAuthorized(lookup string) bool {
	if ausfContext.snRegex.MatchString(lookup) {
		return true
	} else {
		return false
	}
}

func GetSelf() *AUSFContext {
	return &ausfContext
}

func (a *AUSFContext) GetSelfID() string {
	return a.NfId
}

func (c *AUSFContext) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	if !c.OAuth2Required {
		return context.TODO(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType_AUSF, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}

func (c *AUSFContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !c.OAuth2Required {
		logger.UtilLog.Debugf("AUSFContext::AuthorizationCheck: OAuth2 not required\n")
		return nil
	}

	logger.UtilLog.Debugf("AUSFContext::AuthorizationCheck: token[%s] serviceName[%s]\n", token, serviceName)
	return oauth.VerifyOAuth(token, string(serviceName), c.NrfCertPem)
}
