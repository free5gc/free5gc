package amf_context

import (
	"fmt"
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/logger"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var amfContext = AMFContext{}
var TmsiGenerator int32 = 0
var amfUeNgapIdGenerator int64 = 0

func init() {
	AMF_Self().EventSubscriptions = make(map[string]*AMFContextEventSubscription)
	AMF_Self().UePool = make(map[string]*AmfUe)
	AMF_Self().GutiPool = make(map[string]*AmfUe)
	AMF_Self().LadnPool = make(map[string]*LADN)
	AMF_Self().TmsiPool = make(map[int32]*AmfUe)
	AMF_Self().RanUePool = make(map[int64]*RanUe)
	AMF_Self().AmfRanPool = make(map[string]*AmfRan)
	AMF_Self().RanIdPool = make(map[models.GlobalRanNodeId]*AmfRan)
	AMF_Self().EventSubscriptionIDGenerator = 1
	AMF_Self().Name = "amf"
	AMF_Self().UriScheme = models.UriScheme_HTTPS
	AMF_Self().RelativeCapacity = 0xff
	AMF_Self().ServedGuamiList = make([]models.Guami, 0, MaxNumOfServedGuamiList)
	AMF_Self().PlmnSupportList = make([]PlmnSupportItem, 0, MaxNumOfPLMNs)
	AMF_Self().AMFStatusSubscriptionIDGenerator = 1
	AMF_Self().AMFStatusSubscriptions = make(map[string]*models.SubscriptionData)
	AMF_Self().NfService = make(map[models.ServiceName]models.NfService)
	AMF_Self().NetworkName.Full = "free5GC"
}

type AMFContext struct {
	EventSubscriptionIDGenerator     int
	EventSubscriptions               map[string]*AMFContextEventSubscription
	UePool                           map[string]*AmfUe // use imsi as key
	GutiPool                         map[string]*AmfUe
	TmsiPool                         map[int32]*AmfUe // tmsi as key
	RanIdPool                        map[models.GlobalRanNodeId]*AmfRan
	RanUePool                        map[int64]*RanUe   // AmfUeNgapId as key
	AmfRanPool                       map[string]*AmfRan // use remote Addr String as key
	LadnPool                         map[string]*LADN   // ladn as key
	SupportTaiLists                  []models.Tai
	ServedGuamiList                  []models.Guami
	PlmnSupportList                  []PlmnSupportItem
	RelativeCapacity                 int64
	NfId                             string
	Name                             string
	NfService                        map[models.ServiceName]models.NfService // use ServiceName as key, nfservice that amf support
	UriScheme                        models.UriScheme
	HttpIpv4Port                     int
	HttpIPv4Address                  string
	HttpIPv6Address                  string
	TNLWeightFactor                  int64
	SupportDnnLists                  []string
	AMFStatusSubscriptionIDGenerator int
	AMFStatusSubscriptions           map[string]*models.SubscriptionData
	NrfUri                           string
	SecurityAlgorithm                SecurityAlgorithm
	NetworkName                      NetworkName
	NgapIpList                       []string // NGAP Server IP
	T3502Value                       int      // unit is second
	T3512Value                       int      // unit is second
	Non3gppDeregistrationTimerValue  int      // unit is second
}

type AMFContextEventSubscription struct {
	IsAnyUe           bool
	IsGroupUe         bool
	UeSupiList        []string
	Expiry            *time.Time
	EventSubscription models.AmfEventSubscription
}

type PlmnSupportItem struct {
	PlmnId     models.PlmnId   `yaml:"plmnId"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty"`
}

type NetworkName struct {
	Full  string `yaml:"full"`
	Short string `yaml:"short,omitempty"`
}

type SecurityAlgorithm struct {
	IntegrityOrder []uint8 // 8bits(NIA1, NIA2, NIA3 , EIA0, EIA1, EIA2, EIA3, ..)
	CipheringOrder []uint8 // 8bits(NEA1, NEA2, NEA3 , EEA0, EEA1, EEA2, EEA3, ..)
}

func NewPlmnSupportItem() (item PlmnSupportItem) {
	item.SNssaiList = make([]models.Snssai, 0, MaxNumOfSlice)
	return
}

func (context *AMFContext) TmsiAlloc() int32 {
	TmsiGenerator %= math.MaxInt32
	TmsiGenerator++
	for {
		if _, double := context.TmsiPool[TmsiGenerator]; double {
			TmsiGenerator++
		} else {
			break
		}
	}
	return TmsiGenerator
}

func (context *AMFContext) AmfUeNgapIdAlloc() int64 {
	amfUeNgapIdGenerator %= MaxValueOfAmfUeNgapId
	amfUeNgapIdGenerator++
	for {
		if _, double := context.RanUePool[amfUeNgapIdGenerator]; double {
			amfUeNgapIdGenerator++
		} else {
			break
		}
	}
	return amfUeNgapIdGenerator
}

func (context *AMFContext) AllocateGutiToUe(ue *AmfUe) {

	// if ue has a previous tmsi/guti, remove it first
	if ue.Tmsi != 0 {
		delete(context.TmsiPool, ue.Tmsi)
		delete(context.GutiPool, ue.Guti)
	}

	servedGuami := context.ServedGuamiList[0]
	ue.Tmsi = context.TmsiAlloc()

	plmnID := servedGuami.PlmnId.Mcc + servedGuami.PlmnId.Mnc
	tmsiStr := fmt.Sprintf("%08x", ue.Tmsi)
	ue.Guti = plmnID + servedGuami.AmfId + tmsiStr

	context.TmsiPool[ue.Tmsi] = ue
	context.GutiPool[ue.Guti] = ue
}

func (context *AMFContext) AllocateRegistrationArea(ue *AmfUe, anType models.AccessType) {

	// clear the previous registration area if need
	if len(ue.RegistrationArea[anType]) > 0 {
		ue.RegistrationArea[anType] = nil
	}

	// allocate a new tai list as a registration area to ue
	// TODO: algorithm to choose TAI list
	for _, supportTai := range context.SupportTaiLists {
		if reflect.DeepEqual(supportTai, ue.Tai) {
			ue.RegistrationArea[anType] = append(ue.RegistrationArea[anType], supportTai)
			break
		}
	}
}

func (context *AMFContext) AddAmfUeToUePool(ue *AmfUe, supi string) {
	if len(supi) == 0 {
		logger.ContextLog.Errorf("Supi is nil")
	}
	ue.Supi = supi
	context.UePool[ue.Supi] = ue
}

func (context *AMFContext) NewAmfUe(supi string) *AmfUe {
	ue := AmfUe{}
	ue.init()

	if supi != "" {
		context.AddAmfUeToUePool(&ue, supi)
	}

	context.AllocateGutiToUe(&ue)

	return &ue
}

func (context *AMFContext) NewAmfRan(conn net.Conn) *AmfRan {
	ran := AmfRan{}
	ran.SupportedTAList = make([]SupportedTAI, 0, MaxNumOfTAI*MaxNumOfBroadcastPLMNs)
	context.AmfRanPool[conn.RemoteAddr().String()] = &ran
	ran.Conn = conn
	return &ran
}

func (context *AMFContext) InSupportDnnList(targetDnn string) bool {
	for _, dnn := range context.SupportDnnLists {
		if dnn == targetDnn {
			return true
		}
	}
	return false
}

func (context *AMFContext) AmfUeFindByGuti(targetGuti string) *AmfUe {
	if ue, ok := context.GutiPool[targetGuti]; ok {
		return ue
	}
	return nil
}

func (context *AMFContext) AmfRanFindByRanId(ranNodeId models.GlobalRanNodeId) *AmfRan {

	for _, amfRan := range context.AmfRanPool { // amfRan = context.AmfRanPool[i]
		switch amfRan.RanPresent {
		case RanPresentGNbId:
			if amfRan.RanId.GNbId.GNBValue == ranNodeId.GNbId.GNBValue {
				return amfRan
			}
		case RanPresentNgeNbId:
			if amfRan.RanId.NgeNbId == ranNodeId.NgeNbId {
				return amfRan
			}
		case RanPresentN3IwfId:
			if amfRan.RanId.N3IwfId == ranNodeId.N3IwfId {
				return amfRan
			}
		}
	}

	return nil
}

func (context *AMFContext) RanUeFindByAmfUeNgapID(amfUeNgapID int64) *RanUe {
	if ue, ok := context.RanUePool[amfUeNgapID]; ok {
		return ue
	}
	return nil
}

func (context *AMFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
}

func (context *AMFContext) InitNFService(serivceName []string, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, nameString := range serivceName {
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

// Reset AMF Context
func (context *AMFContext) Reset() {
	for key := range context.AmfRanPool {
		delete(context.AmfRanPool, key)
	}
	for key := range context.GutiPool {
		delete(context.GutiPool, key)
	}
	for key := range context.LadnPool {
		delete(context.LadnPool, key)
	}
	for key := range context.RanUePool {
		delete(context.RanUePool, key)
	}
	for key := range context.UePool {
		delete(context.UePool, key)
	}
	for key := range context.TmsiPool {
		delete(context.TmsiPool, key)
	}
	for key := range context.RanIdPool {
		delete(context.RanIdPool, key)
	}
	for key := range context.EventSubscriptions {
		delete(context.EventSubscriptions, key)
	}
	for key := range context.NfService {
		delete(context.NfService, key)
	}
	context.SupportTaiLists = context.SupportTaiLists[:0]
	context.PlmnSupportList = context.PlmnSupportList[:0]
	context.ServedGuamiList = context.ServedGuamiList[:0]
	context.EventSubscriptionIDGenerator = 1
	context.RelativeCapacity = 0xff
	context.NfId = ""
	context.UriScheme = models.UriScheme_HTTPS
	context.HttpIpv4Port = 0
	context.HttpIPv4Address = ""
	context.HttpIPv6Address = ""
	context.Name = "amf"
	context.NrfUri = ""
	TmsiGenerator = 0
	amfUeNgapIdGenerator = 0
}

// Create new AMF context
func AMF_Self() *AMFContext {
	return &amfContext
}
