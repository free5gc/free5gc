package pcf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
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
	PCF_Self().UePool = make(map[string]*UeContext)
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
	NrfUri          string
	UdrUri          string
	NotifiUri       string
	BdtPolicyUri    string
	BdtUri          string
	UePool          map[string]*UeContext
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

var AmpolicyUri = "/npcf-am-policy-control/v1"
var SmpolicyUri = "/npcf-smpolicycontrol/v1/sm-policies/"
var PolicyAuthorizationUri = "/npcf-policyauthorization/v1/app-sessions/"
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

func (context *PCFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
}

func (context *PCFContext) InitNFService(serviceName []string, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, nameString := range serviceName {
		name := models.ServiceName(nameString)
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
		}
	}
}

func (context *PCFContext) NewPCFUe(Supi string) (*UeContext, error) {
	if strings.HasPrefix(Supi, "imsi-") {
		context.UePool[Supi] = &UeContext{}
		context.UePool[Supi].SmPolicyData = make(map[string]*UeSmPolicyData)
		context.UePool[Supi].AMPolicyData = make(map[string]*UeAMPolicyData)
		context.UePool[Supi].PolAssociationIDGenerator = 1
		context.UePool[Supi].Supi = Supi
		return context.UePool[Supi], nil
	} else {
		return nil, fmt.Errorf(" add Ue context fail ")
	}
}

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

func (context *PCFContext) PCFUeFindByIPv4(v4 string) *UeContext {
	for _, ue := range context.UePool {
		if ue.SmPolicyControlStore.Context.Ipv4Address == v4 {
			return ue
		}
	}
	return nil
}
func (context *PCFContext) PCFUeFindByIPv6(v6 string) *UeContext {
	for _, ue := range context.UePool {
		if ue.SmPolicyControlStore.Context.Ipv6AddressPrefix == v6 {
			return ue
		}
	}
	return nil
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
