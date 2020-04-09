package pcf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/pcf/factory"
	"strconv"
	"strings"
)

var pcfContext = PCFContext{}

func init() {
	PCF_Self().Name = "pcf"
	PCF_Self().UriScheme = models.UriScheme_HTTPS
	PCF_Self().TimeFormat = "2006-01-02 15:04:05"
	PCF_Self().DefaultBdtRefId = "BdtPolicyId-"
	PCF_Self().NfService = make(map[models.ServiceName]models.NfService)
	PCF_Self().PcfServiceUris = make(map[models.ServiceName]string)
	PCF_Self().PcfSuppFeats = make(map[models.ServiceName][]byte)
	PCF_Self().UePool = make(map[string]*UeContext)
	PCF_Self().BdtPolicyPool = make(map[string]models.BdtPolicy)
	PCF_Self().BdtPolicyIdGenerator = 1
	PCF_Self().AppSessionPool = make(map[string]*AppSessionData)
	PCF_Self().AMFStatusSubsData = make(map[string]AMFStatusSubscriptionData)
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
	PcfServiceUris  map[models.ServiceName]string
	PcfSuppFeats    map[models.ServiceName][]byte
	NrfUri          string
	DefaultUdrUri   string
	UePool          map[string]*UeContext
	// Bdt Policy related
	BdtPolicyPool        map[string]models.BdtPolicy // use BdtPolicyId as key
	BdtPolicyIdGenerator uint64
	// App Session related
	AppSessionPool map[string]*AppSessionData // use AppSessionId(ue.Supi-%d) or (BdtRefId-%d) as key
	// AMF Status Change Subscription related
	AMFStatusSubsData map[string]AMFStatusSubscriptionData // subscriptionId as key
}

type AMFStatusSubscriptionData struct {
	AmfUri       string
	AmfStatusUri string
	GuamiList    []models.Guami
}

type AppSessionData struct {
	AppSessionId      string
	AppSessionContext *models.AppSessionContext
	// (compN/compN-subCompN/appId-%s) map to PccRule
	RelatedPccRuleIds    map[string]string
	PccRuleIdMapToCompId map[string]string
	// EventSubscription
	Events   map[models.AfEvent]models.AfNotifMethod
	EventUri string
	// related Session
	SmPolicyData *UeSmPolicyData
}

// Create new PCF context
func PCF_Self() *PCFContext {
	return &pcfContext
}
func GetTimeformat() string {
	return pcfContext.TimeFormat
}
func GetPcfContext() PCFContext {
	return pcfContext
}
func GetUri(name models.ServiceName) string {
	return pcfContext.PcfServiceUris[name]
}

var PolicyAuthorizationUri = "/npcf-policyauthorization/v1/app-sessions/"
var SmUri = "/npcf-smpolicycontrol/v1"
var IPv4Address = "192.168."
var IPv6Address = "ffab::"
var CheckNotifiUri = "/npcf-callback/v1/nudr-notify/"
var Ipv4_pool = make(map[string]string)
var Ipv6_pool = make(map[string]string)

// BdtPolicy default value
const DefaultBdtRefId = "BdtPolicyId-"

func (context *PCFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
}

// Init NfService with supported service list ,and version of services
func (context *PCFContext) InitNFService(serviceList []factory.Service, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, service := range serviceList {
		name := models.ServiceName(service.ServiceName)
		context.NfService[name] = models.NfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          context.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       context.GetIPv4Uri(),
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.HttpIPv4Address,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.HttpIpv4Port),
				},
			},
			SupportedFeatures: service.SuppFeat,
		}
	}
}

// Allocate PCF Ue with supi and add to pcf Context and returns allocated ue
func (context *PCFContext) NewPCFUe(Supi string) (*UeContext, error) {
	if strings.HasPrefix(Supi, "imsi-") {
		context.UePool[Supi] = &UeContext{}
		context.UePool[Supi].SmPolicyData = make(map[string]*UeSmPolicyData)
		context.UePool[Supi].AMPolicyData = make(map[string]*UeAMPolicyData)
		context.UePool[Supi].PolAssociationIDGenerator = 1
		context.UePool[Supi].AppSessionIdGenerator = 1
		context.UePool[Supi].Supi = Supi
		return context.UePool[Supi], nil
	} else {
		return nil, fmt.Errorf(" add Ue context fail ")
	}
}

// Return Bdt Policy Id with format "BdtPolicyId-%d" which be allocated
func (context *PCFContext) AllocBdtPolicyId() string {
	bdtPolicyId := fmt.Sprintf("BdtPolicyId-%d", context.BdtPolicyIdGenerator)
	_, exist := context.BdtPolicyPool[bdtPolicyId]
	for exist {
		context.BdtPolicyIdGenerator++
		bdtPolicyId := fmt.Sprintf("BdtPolicyId-%d", context.BdtPolicyIdGenerator)
		_, exist = context.BdtPolicyPool[bdtPolicyId]
	}
	context.BdtPolicyIdGenerator++
	return bdtPolicyId
}

// Find PcfUe which the policyId belongs to
func (context *PCFContext) PCFUeFindByPolicyId(PolicyId string) *UeContext {
	index := strings.LastIndex(PolicyId, "-")
	if index == -1 {
		return nil
	}
	supi := PolicyId[:index]
	if supi != "" {
		return context.UePool[supi]
	}
	return nil
}

// Find PcfUe which the AppSessionId belongs to
func (context *PCFContext) PCFUeFindByAppSessionId(appSessionId string) *UeContext {
	index := strings.LastIndex(appSessionId, "-")
	if index == -1 {
		return nil
	}
	supi := appSessionId[:index]
	if supi != "" {
		return context.UePool[supi]
	}
	return nil
}

// Find PcfUe which Ipv4 belongs to
func (context *PCFContext) PcfUeFindByIPv4(v4 string) *UeContext {
	for _, ue := range context.UePool {
		if ue.SMPolicyFindByIpv6(v4) != nil {
			return ue
		}
	}
	return nil
}

// Find PcfUe which Ipv6 belongs to
func (context *PCFContext) PcfUeFindByIPv6(v6 string) *UeContext {
	for _, ue := range context.UePool {
		if ue.SMPolicyFindByIpv6(v6) != nil {
			return ue
		}
	}
	return nil
}

// Session Binding from application request to get corresponding Sm policy
func (context *PCFContext) SessionBinding(req *models.AppSessionContextReqData) (policy *UeSmPolicyData, err error) {
	// TODO: support dnn, snssai, ... because Ip Address is not enough with same ip address in different ip domains, details in subclause 4.2.2.2 of TS 29514
	if ue, exist := context.UePool[req.Supi]; exist {
		if req.UeIpv4 != "" {
			policy = ue.SMPolicyFindByIpv4(req.UeIpv4)
		} else if req.UeIpv6 != "" {
			policy = ue.SMPolicyFindByIpv6(req.UeIpv6)
		} else {
			err = fmt.Errorf("Ue finding by MAC address does not support")
		}
	} else if req.UeIpv4 != "" {
		policy = ue.SMPolicyFindByIpv4(req.UeIpv4)
	} else if req.UeIpv6 != "" {
		policy = ue.SMPolicyFindByIpv6(req.UeIpv6)
	} else {
		err = fmt.Errorf("Ue finding by MAC address does not support")
	}
	if err == nil && policy == nil {
		if req.UeIpv4 != "" {
			err = fmt.Errorf("Can't find Ue with Ipv4[%s]", req.UeIpv4)
		} else {
			err = fmt.Errorf("Can't find Ue with Ipv6[%s]", req.UeIpv6)
		}
	}
	return
}

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
