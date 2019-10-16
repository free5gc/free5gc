package ausf_context

import (
	// "fmt"
	"free5gc/lib/openapi/models"
	"regexp"
)

type AUSFContext struct {
	NfId            string
	GroupId         string
	HttpIpv4Port    int
	HttpIPv4Address string
	Url             string
	UriScheme       models.UriScheme
	NrfUri          string
	NfService       map[models.ServiceName]models.NfService
	PlmnList        []models.PlmnId
	UdmUeauUrl      string
}

type AusfUeContext struct {
	Supi               string
	Kausf              string
	Kseaf              string
	ServingNetworkName string
	AuthStatus         models.AuthResult
	UdmUeauUrl         string

	// for 5G AKA
	XresStar string

	// for EAP-AKA'
	K_aut string
	XRES  string
}

const (
	EAP_AKA_PRIME_TYPENUM = 50
)

// Attribute Types for EAP-AKA'
const (
	AT_RAND_ATTRIBUTE         = 1
	AT_AUTN_ATTRIBUTE         = 2
	AT_RES_ATTRIBUTE          = 3
	AT_MAC_ATTRIBUTE          = 11
	AT_NOTIFICATION_ATTRIBUTE = 12
	AT_IDENTITY_ATTRIBUTE     = 14
	AT_KDF_INPUT_ATTRIBUTE    = 23
	AT_KDF_ATTRIBUTE          = 24
)

var ausfContext AUSFContext
var ausfUeContextPool map[string]*AusfUeContext
var snRegex *regexp.Regexp

func Init() {
	ausfUeContextPool = make(map[string]*AusfUeContext)
	snRegex, _ = regexp.Compile("5G:mnc[0-9]{3}[.]mcc[0-9]{3}[.]3gppnetwork[.]org")
	InitAusfContext(&ausfContext)
}

func NewAusfUeContext(identifier string) (ausfUeContext *AusfUeContext) {
	ausfUeContext = new(AusfUeContext)
	ausfUeContext.Supi = identifier // supi
	return ausfUeContext
}

func AddAusfUeContextToPool(ausfUeContext *AusfUeContext) {
	ausfUeContextPool[ausfUeContext.Supi] = ausfUeContext
}

func CheckIfAusfUeContextExists(ref string) bool {
	return (ausfUeContextPool[ref] != nil)
}

func GetAusfUeContext(ref string) (ausfUeContext *AusfUeContext) {
	ausfUeContext = ausfUeContextPool[ref]
	return ausfUeContext
}

func IsServingNetworkAuthorized(lookup string) bool {
	if snRegex.MatchString(lookup) {
		return true
	} else {
		return false
	}
}

func GetSelf() *AUSFContext {
	return &ausfContext
}

func (a AUSFContext) GetSelfID() string {
	return a.NfId
}
