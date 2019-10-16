package pcf_context

import (
	"errors"
	"fmt"
	"free5gc/lib/openapi/models"
)

var pcfContext = PCFContext{}

func init() {
	PCF_Self().Name = "pcf"
	PCF_Self().UriScheme = models.UriScheme_HTTPS
	PCF_Self().TimeFormat = "2006-01-02 15:04:05"
	PCF_Self().DefaultBdtRefId = "BdtPolicyId-"
}

type PCFContext struct {
	NfId            string
	Name            string
	UriScheme       models.UriScheme
	HttpIPv4Address string
	HttpIpv4Port    int
	TimeFormat      string
	DefaultBdtRefId string
	NfService       map[models.ServiceName]models.NfService
	AmpolicyUri     string
	NotifiUri       string
	BdtPolicyUri    string
	BdtUri          string
}

// Create new PCF context
func PCF_Self() *PCFContext {
	return &pcfContext
}
func GetTimeformat() string {
	return pcfContext.TimeFormat
}
func Getampolicyuri() string {
	return pcfContext.AmpolicyUri
}

var AmpolicyUri = "/npcf-am-policy-control/v1/policies/"
var SmpolicyUri = "/npcf-smpolicycontrol/v1/sm-policies/"
var SmUri = "/npcf-smpolicycontrol/v1"
var IPv4Address = "192.168."
var IPv6Address = "ffab::"
var NotifiUri = "https://localhost:29507/npcf-smpolicycontrol/v1/nudr-notify/"
var CheckNotifiUri = "/npcf-smpolicycontrol/v1/nudr-notify/"
var Ipv4_pool = make(map[string]string)
var Ipv6_pool = make(map[string]string)

// BdtPolicy default value
const DefaultBdtRefId = "BdtPolicyId-"

var BdtPolicyUri = "/npcf-bdtpolicycontrol/v1/bdtpolicies/"
var BdtUri = "/npcf-bdtpolicycontrol/v1"

// key is supi
var pcfUeContext = make(map[string]*PCFUeContext)

type PCFUeContext struct {
	Supi                        string
	AspId                       string
	AppSessionIdStore           *AppSessionIdStore
	PolAssociationIDStore       *PolAssociationIDStore
	SmPolicyControlStore        *models.SmPolicyControl
	PolicyDataSubscriptionStore *models.PolicyDataSubscription
	BdtPolicyTimeout            bool
	BdtPolicyStore              *models.BdtPolicy
	PolicyDataChangeStore       *models.PolicyDataChangeNotification
}

func CheckAspidOnPcfUeContext(aspId string) (key string, err error) {
	for key := range pcfUeContext {
		if aspId == pcfUeContext[key].AspId {
			return key, nil
		}
	}
	return "", errors.New(" Not found Aspid on PcfUeContext ")
}

func GetPCFUeContext() (make map[string]*PCFUeContext) {
	return pcfUeContext
}

func NewPCFUe(Supi string) error {
	if Supi != "" {
		pcfUeContext[Supi] = &PCFUeContext{}
		pcfUeContext[Supi].Supi = Supi
		return nil
	} else {
		return errors.New(" add Ue context fail ")
	}
}

func AddAspIdToUe(Supi string, aspId string) error {
	if Supi != "" {
		pcfUeContext[Supi] = &PCFUeContext{}
		pcfUeContext[Supi].AspId = aspId
		return nil
	} else {
		return errors.New(" add aspId to Ue context fail ")
	}
}

// polAssoidTemp -
type PolAssociationIDStore struct {
	PolAssoId                     string
	PolAssoidTemp                 models.PolicyAssociation
	PolAssoidUpdateTemp           models.PolicyUpdate
	PolAssoidSubcCatsTemp         models.AmPolicyData
	PolAssoidDataSubscriptionTemp models.PolicyDataSubscription
	PolAssoidSubscriptiondataTemp models.SubscriptionData
}

//var polAssociationContextStore []PolAssociationIDStore

// AppSessionIdStore -
type AppSessionIdStore struct {
	AppSessionId      string
	AppSessionContext models.AppSessionContext
}

var AppSessionContextStore []AppSessionIdStore

// BdtPolicyData_store -
var BdtPolicyData_store []models.BdtPolicyData
var CreateFailBdtDateStore []models.BdtData

func Ipv4Pool(ipindex int32) string {
	ipv4address := IPv4Address + fmt.Sprint((int(ipindex)/255)+1) + "." + fmt.Sprint(int(ipindex)%255)
	return ipv4address
}
func Ipv4Index() int32 {

	if len(Ipv4_pool) == 0 {
		Ipv4_pool["1"] = Ipv4Pool(1)
	} else {
		for i := 1; i <= len(Ipv4_pool); i++ {
			if Ipv4_pool[fmt.Sprint(i)] == "" {
				Ipv4_pool[fmt.Sprint(i)] = Ipv4Pool(int32(i))
				return int32(i)
			}
		}

		Ipv4_pool[fmt.Sprint(int32(len(Ipv4_pool)+1))] = Ipv4Pool(int32(len(Ipv4_pool) + 1))
		return int32(len(Ipv4_pool))
	}
	return 1
}
func GetIpv4Address(ipindex int32) string {
	return Ipv4_pool[fmt.Sprint(ipindex)]
}
func DeleteIpv4index(Ipv4index int32) {
	delete(Ipv4_pool, fmt.Sprint(Ipv4index))
}
func Ipv6Pool(ipindex int32) string {

	ipv6address := IPv6Address + fmt.Sprintf("%x\n", ipindex)
	return ipv6address
}
func Ipv6Index() int32 {

	if len(Ipv6_pool) == 0 {
		Ipv6_pool["1"] = Ipv6Pool(1)
	} else {
		for i := 1; i <= len(Ipv6_pool); i++ {
			if Ipv6_pool[fmt.Sprint(i)] == "" {
				Ipv6_pool[fmt.Sprint(i)] = Ipv6Pool(int32(i))
				return int32(i)
			}
		}

		Ipv6_pool[fmt.Sprint(int32(len(Ipv6_pool)+1))] = Ipv6Pool(int32(len(Ipv6_pool) + 1))
		return int32(len(Ipv6_pool))
	}
	return 1
}
func GetIpv6Address(ipindex int32) string {
	return Ipv6_pool[fmt.Sprint(ipindex)]
}
func DeleteIpv6index(Ipv6index int32) {
	delete(Ipv6_pool, fmt.Sprint(Ipv6index))
}
