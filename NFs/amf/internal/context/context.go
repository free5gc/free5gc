package context

import (
	"context"
	"fmt"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/free5gc/amf/internal/logger"
	"github.com/free5gc/amf/pkg/factory"
	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/security"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/openapi/oauth"
	"github.com/free5gc/util/idgenerator"
)

var (
	amfContext                       AMFContext
	tmsiGenerator                    *idgenerator.IDGenerator = nil
	amfUeNGAPIDGenerator             *idgenerator.IDGenerator = nil
	amfStatusSubscriptionIDGenerator *idgenerator.IDGenerator = nil
)

func init() {
	GetSelf().LadnPool = make(map[string]factory.Ladn)
	GetSelf().EventSubscriptionIDGenerator = idgenerator.NewGenerator(1, math.MaxInt32)
	GetSelf().Name = "amf"
	GetSelf().UriScheme = models.UriScheme_HTTPS
	GetSelf().RelativeCapacity = 0xff
	GetSelf().ServedGuamiList = make([]models.Guami, 0, MaxNumOfServedGuamiList)
	GetSelf().PlmnSupportList = make([]factory.PlmnSupportItem, 0, MaxNumOfPLMNs)
	GetSelf().NfService = make(map[models.ServiceName]models.NrfNfManagementNfService)
	GetSelf().NetworkName.Full = "free5GC"
	tmsiGenerator = idgenerator.NewGenerator(1, math.MaxInt32)
	amfStatusSubscriptionIDGenerator = idgenerator.NewGenerator(1, math.MaxInt32)
	amfUeNGAPIDGenerator = idgenerator.NewGenerator(1, MaxValueOfAmfUeNgapId)
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &AMFContext{}

type AMFContext struct {
	EventSubscriptionIDGenerator *idgenerator.IDGenerator
	EventSubscriptions           sync.Map
	UePool                       sync.Map                // map[supi]*AmfUe
	RanUePool                    sync.Map                // map[AmfUeNgapID]*RanUe
	AmfRanPool                   sync.Map                // map[net.Conn]*AmfRan
	LadnPool                     map[string]factory.Ladn // dnn as key
	SupportTaiLists              []models.Tai
	ServedGuamiList              []models.Guami
	PlmnSupportList              []factory.PlmnSupportItem
	RelativeCapacity             int64
	NfId                         string
	Name                         string
	NfService                    map[models.ServiceName]models.NrfNfManagementNfService // nfservice that amf support
	UriScheme                    models.UriScheme
	BindingIPv4                  string
	SBIPort                      int
	RegisterIPv4                 string
	HttpIPv6Address              string
	TNLWeightFactor              int64
	SupportDnnLists              []string
	AMFStatusSubscriptions       sync.Map // map[subscriptionID]models.SubscriptionData
	NrfUri                       string
	NrfCertPem                   string
	SecurityAlgorithm            SecurityAlgorithm
	NetworkName                  factory.NetworkName
	NgapIpList                   []string // NGAP Server IP
	NgapPort                     int
	T3502Value                   int    // unit is second
	T3512Value                   int    // unit is second
	Non3gppDeregTimerValue       int    // unit is second
	TimeZone                     string // "[+-]HH:MM[+][1-2]", Refer to TS 29.571 - 5.2.2 Simple Data Types
	// read-only fields
	T3513Cfg factory.TimerValue
	T3522Cfg factory.TimerValue
	T3550Cfg factory.TimerValue
	T3560Cfg factory.TimerValue
	T3565Cfg factory.TimerValue
	T3570Cfg factory.TimerValue
	T3555Cfg factory.TimerValue
	Locality string

	OAuth2Required bool
}

type AMFContextEventSubscription struct {
	IsAnyUe           bool
	IsGroupUe         bool
	UeSupiList        []string
	Expiry            *time.Time
	EventSubscription models.AmfEventSubscription
}

type SecurityAlgorithm struct {
	IntegrityOrder []uint8 // slice of security.AlgIntegrityXXX
	CipheringOrder []uint8 // slice of security.AlgCipheringXXX
}

func InitAmfContext(context *AMFContext) {
	config := factory.AmfConfig
	logger.UtilLog.Infof("amfconfig Info: Version[%s]", config.GetVersion())
	configuration := config.Configuration
	context.NfId = uuid.New().String()
	if configuration.AmfName != "" {
		context.Name = configuration.AmfName
	}
	if configuration.NgapIpList != nil {
		context.NgapIpList = configuration.NgapIpList
	} else {
		context.NgapIpList = []string{"127.0.0.1"} // default localhost
	}
	context.NgapPort = config.GetNgapPort()
	context.UriScheme = models.UriScheme(config.GetSbiScheme())
	context.RegisterIPv4 = config.GetSbiRegisterIP()
	context.SBIPort = config.GetSbiPort()
	context.BindingIPv4 = config.GetSbiBindingIP()

	context.InitNFService(config.GetServiceNameList(), config.GetVersion())
	context.ServedGuamiList = configuration.ServedGumaiList
	context.SupportTaiLists = configuration.SupportTAIList
	context.PlmnSupportList = configuration.PlmnSupportList
	context.SupportDnnLists = configuration.SupportDnnList
	for _, ladn := range configuration.SupportLadnList {
		context.LadnPool[ladn.Dnn] = ladn
	}
	context.NrfUri = config.GetNrfUri()
	context.NrfCertPem = configuration.NrfCertPem
	security := configuration.Security
	if security != nil {
		context.SecurityAlgorithm.IntegrityOrder = getIntAlgOrder(security.IntegrityOrder)
		context.SecurityAlgorithm.CipheringOrder = getEncAlgOrder(security.CipheringOrder)
	}
	context.NetworkName = configuration.NetworkName
	context.TimeZone = nasConvert.GetTimeZone(time.Now())
	context.T3502Value = configuration.T3502Value
	context.T3512Value = configuration.T3512Value
	context.Non3gppDeregTimerValue = configuration.Non3gppDeregTimerValue
	context.T3513Cfg = configuration.T3513
	context.T3522Cfg = configuration.T3522
	context.T3550Cfg = configuration.T3550
	context.T3560Cfg = configuration.T3560
	context.T3565Cfg = configuration.T3565
	context.T3570Cfg = configuration.T3570
	context.T3555Cfg = configuration.T3555
	context.Locality = configuration.Locality
}

func getIntAlgOrder(integrityOrder []string) (intOrder []uint8) {
	for _, intAlg := range integrityOrder {
		switch intAlg {
		case "NIA0":
			intOrder = append(intOrder, security.AlgIntegrity128NIA0)
		case "NIA1":
			intOrder = append(intOrder, security.AlgIntegrity128NIA1)
		case "NIA2":
			intOrder = append(intOrder, security.AlgIntegrity128NIA2)
		case "NIA3":
			intOrder = append(intOrder, security.AlgIntegrity128NIA3)
		default:
			logger.UtilLog.Errorf("Unsupported algorithm: %s", intAlg)
		}
	}
	return
}

func getEncAlgOrder(cipheringOrder []string) (encOrder []uint8) {
	for _, encAlg := range cipheringOrder {
		switch encAlg {
		case "NEA0":
			encOrder = append(encOrder, security.AlgCiphering128NEA0)
		case "NEA1":
			encOrder = append(encOrder, security.AlgCiphering128NEA1)
		case "NEA2":
			encOrder = append(encOrder, security.AlgCiphering128NEA2)
		case "NEA3":
			encOrder = append(encOrder, security.AlgCiphering128NEA3)
		default:
			logger.UtilLog.Errorf("Unsupported algorithm: %s", encAlg)
		}
	}
	return
}

func NewPlmnSupportItem() (item factory.PlmnSupportItem) {
	item.SNssaiList = make([]models.Snssai, 0, MaxNumOfSlice)
	return
}

func (context *AMFContext) TmsiAllocate() int32 {
	tmsi, err := tmsiGenerator.Allocate()
	if err != nil {
		logger.CtxLog.Errorf("Allocate TMSI error: %+v", err)
		return -1
	}
	return int32(tmsi)
}

func (context *AMFContext) FreeTmsi(tmsi int64) {
	tmsiGenerator.FreeID(tmsi)
}

func (context *AMFContext) AllocateAmfUeNgapID() (int64, error) {
	return amfUeNGAPIDGenerator.Allocate()
}

func (context *AMFContext) AllocateGutiToUe(ue *AmfUe) {
	servedGuami := context.ServedGuamiList[0]
	ue.Tmsi = context.TmsiAllocate()

	plmnID := servedGuami.PlmnId.Mcc + servedGuami.PlmnId.Mnc
	tmsiStr := fmt.Sprintf("%08x", ue.Tmsi)
	ue.Guti = plmnID + servedGuami.AmfId + tmsiStr
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

func (context *AMFContext) NewAMFStatusSubscription(subscriptionData models.AmfCommunicationSubscriptionData) (
	subscriptionID string,
) {
	id, err := amfStatusSubscriptionIDGenerator.Allocate()
	if err != nil {
		logger.CtxLog.Errorf("Allocate subscriptionID error: %+v", err)
		return ""
	}

	subscriptionID = strconv.Itoa(int(id))
	context.AMFStatusSubscriptions.Store(subscriptionID, subscriptionData)
	return
}

// Return Value: (subscriptionData *models.SubScriptionData, ok bool)
func (context *AMFContext) FindAMFStatusSubscription(subscriptionID string) (
	*models.AmfCommunicationSubscriptionData, bool,
) {
	if value, ok := context.AMFStatusSubscriptions.Load(subscriptionID); ok {
		subscriptionData := value.(models.AmfCommunicationSubscriptionData)
		return &subscriptionData, ok
	} else {
		return nil, false
	}
}

func (context *AMFContext) DeleteAMFStatusSubscription(subscriptionID string) {
	context.AMFStatusSubscriptions.Delete(subscriptionID)
	if id, err := strconv.ParseInt(subscriptionID, 10, 64); err != nil {
		logger.CtxLog.Error(err)
	} else {
		amfStatusSubscriptionIDGenerator.FreeID(id)
	}
}

func (context *AMFContext) NewEventSubscription(subscriptionID string, subscription *AMFContextEventSubscription) {
	context.EventSubscriptions.Store(subscriptionID, subscription)
}

func (context *AMFContext) FindEventSubscription(subscriptionID string) (*AMFContextEventSubscription, bool) {
	if value, ok := context.EventSubscriptions.Load(subscriptionID); ok {
		return value.(*AMFContextEventSubscription), ok
	} else {
		return nil, false
	}
}

func (context *AMFContext) DeleteEventSubscription(subscriptionID string) {
	context.EventSubscriptions.Delete(subscriptionID)
	if id, err := strconv.ParseInt(subscriptionID, 10, 32); err != nil {
		logger.CtxLog.Error(err)
	} else {
		context.EventSubscriptionIDGenerator.FreeID(id)
	}
}

func (context *AMFContext) AddAmfUeToUePool(ue *AmfUe, supi string) {
	if len(supi) == 0 {
		logger.CtxLog.Errorf("Supi is nil")
	}
	ue.Supi = supi
	context.UePool.Store(ue.Supi, ue)
}

func (context *AMFContext) NewAmfUe(supi string) *AmfUe {
	ue := AmfUe{}
	ue.init()

	if supi != "" {
		context.AddAmfUeToUePool(&ue, supi)
	}

	context.AllocateGutiToUe(&ue)

	logger.CtxLog.Infof("New AmfUe [supi:%s][guti:%s]", supi, ue.Guti)
	return &ue
}

func (context *AMFContext) AmfUeFindByUeContextID(ueContextID string) (*AmfUe, bool) {
	if strings.HasPrefix(ueContextID, "imsi") {
		return context.AmfUeFindBySupi(ueContextID)
	}
	if strings.HasPrefix(ueContextID, "imei") {
		return context.AmfUeFindByPei(ueContextID)
	}
	if strings.HasPrefix(ueContextID, "5g-guti") {
		guti := ueContextID[strings.LastIndex(ueContextID, "-")+1:]
		return context.AmfUeFindByGuti(guti)
	}
	return nil, false
}

func (context *AMFContext) AmfUeFindBySupi(supi string) (*AmfUe, bool) {
	if value, ok := context.UePool.Load(supi); ok {
		return value.(*AmfUe), ok
	}
	return nil, false
}

func (context *AMFContext) AmfUeFindBySuci(suci string) (ue *AmfUe, ok bool) {
	context.UePool.Range(func(key, value interface{}) bool {
		candidate := value.(*AmfUe)
		if ok = (candidate.Suci == suci); ok {
			ue = candidate
			return false
		}
		return true
	})
	return
}

func (context *AMFContext) AmfUeFindByPei(pei string) (*AmfUe, bool) {
	var ue *AmfUe
	var ok bool
	context.UePool.Range(func(key, value interface{}) bool {
		candidate := value.(*AmfUe)
		if ok = (candidate.Pei == pei); ok {
			ue = candidate
			return false
		}
		return true
	})
	return ue, ok
}

func (context *AMFContext) NewAmfRan(conn net.Conn) *AmfRan {
	ran := AmfRan{}
	ran.SupportedTAList = make([]SupportedTAI, 0, MaxNumOfTAI*MaxNumOfBroadcastPLMNs)
	ran.Conn = conn
	addr := conn.RemoteAddr()
	if addr != nil {
		ran.Log = logger.NgapLog.WithField(logger.FieldRanAddr, addr.String())
	} else {
		ran.Log = logger.NgapLog.WithField(logger.FieldRanAddr, "(nil)")
	}

	context.AmfRanPool.Store(conn, &ran)
	return &ran
}

// use net.Conn to find RAN context, return *AmfRan and ok bit
func (context *AMFContext) AmfRanFindByConn(conn net.Conn) (*AmfRan, bool) {
	if value, ok := context.AmfRanPool.Load(conn); ok {
		return value.(*AmfRan), ok
	}
	return nil, false
}

// use ranNodeID to find RAN context, return *AmfRan and ok bit
func (context *AMFContext) AmfRanFindByRanID(ranNodeID models.GlobalRanNodeId) (*AmfRan, bool) {
	var ran *AmfRan
	var ok bool
	context.AmfRanPool.Range(func(key, value interface{}) bool {
		amfRan := value.(*AmfRan)
		if amfRan.RanId == nil {
			return true
		}

		switch amfRan.RanPresent {
		case RanPresentGNbId:
			if amfRan.RanId.GNbId != nil && ranNodeID.GNbId != nil &&
				amfRan.RanId.GNbId.GNBValue == ranNodeID.GNbId.GNBValue {
				ran = amfRan
				ok = true
				return false
			}
		case RanPresentNgeNbId:
			if amfRan.RanId.NgeNbId == ranNodeID.NgeNbId {
				ran = amfRan
				ok = true
				return false
			}
		case RanPresentN3IwfId:
			if amfRan.RanId.N3IwfId == ranNodeID.N3IwfId {
				ran = amfRan
				ok = true
				return false
			}
		}
		return true
	})
	return ran, ok
}

func (context *AMFContext) DeleteAmfRan(conn net.Conn) {
	context.AmfRanPool.Delete(conn)
}

func (context *AMFContext) InSupportDnnList(targetDnn string) bool {
	for _, dnn := range context.SupportDnnLists {
		if dnn == targetDnn {
			return true
		}
	}
	return false
}

func (context *AMFContext) InPlmnSupportList(snssai models.Snssai) bool {
	for _, plmnSupportItem := range context.PlmnSupportList {
		for _, supportSnssai := range plmnSupportItem.SNssaiList {
			if openapi.SnssaiEqualFold(supportSnssai, snssai) {
				return true
			}
		}
	}
	return false
}

func (context *AMFContext) AmfUeFindByGuti(guti string) (*AmfUe, bool) {
	var ue *AmfUe
	var ok bool
	context.UePool.Range(func(key, value interface{}) bool {
		candidate := value.(*AmfUe)
		if ok = (candidate.Guti == guti); ok {
			ue = candidate
			return false
		}
		return true
	})
	return ue, ok
}

func (context *AMFContext) AmfUeFindByPolicyAssociationID(polAssoId string) (*AmfUe, bool) {
	var ue *AmfUe
	var ok bool
	context.UePool.Range(func(key, value interface{}) bool {
		candidate := value.(*AmfUe)
		if ok = (candidate.PolicyAssociationId == polAssoId); ok {
			ue = candidate
			return false
		}
		return true
	})
	return ue, ok
}

func (context *AMFContext) RanUeFindByAmfUeNgapID(amfUeNgapID int64) *RanUe {
	if value, ok := context.RanUePool.Load(amfUeNgapID); ok {
		return value.(*RanUe)
	}
	return nil
}

func (context *AMFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
}

func (context *AMFContext) InitNFService(serivceName []string, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, nameString := range serivceName {
		name := models.ServiceName(nameString)
		context.NfService[name] = models.NrfNfManagementNfService{
			ServiceInstanceId: strconv.Itoa(index),
			ServiceName:       name,
			Versions: []models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          context.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       context.GetIPv4Uri(),
			IpEndPoints: []models.IpEndPoint{
				{
					Ipv4Address: context.RegisterIPv4,
					Transport:   models.NrfNfManagementTransportProtocol_TCP,
					Port:        int32(context.SBIPort),
				},
			},
		}
	}
}

// Reset AMF Context
func (context *AMFContext) Reset() {
	context.AmfRanPool.Range(func(key, value interface{}) bool {
		context.UePool.Delete(key)
		return true
	})
	for key := range context.LadnPool {
		delete(context.LadnPool, key)
	}
	context.RanUePool.Range(func(key, value interface{}) bool {
		context.RanUePool.Delete(key)
		return true
	})
	context.UePool.Range(func(key, value interface{}) bool {
		context.UePool.Delete(key)
		return true
	})
	context.EventSubscriptions.Range(func(key, value interface{}) bool {
		context.DeleteEventSubscription(key.(string))
		return true
	})
	for key := range context.NfService {
		delete(context.NfService, key)
	}
	context.SupportTaiLists = context.SupportTaiLists[:0]
	context.PlmnSupportList = context.PlmnSupportList[:0]
	context.ServedGuamiList = context.ServedGuamiList[:0]
	context.RelativeCapacity = 0xff
	context.NfId = ""
	context.UriScheme = models.UriScheme_HTTPS
	context.SBIPort = 0
	context.BindingIPv4 = ""
	context.RegisterIPv4 = ""
	context.HttpIPv6Address = ""
	context.Name = "amf"
	context.NrfUri = ""
	context.NrfCertPem = ""
	context.OAuth2Required = false
}

// Create new AMF context
func GetSelf() *AMFContext {
	return &amfContext
}

func (c *AMFContext) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	if !c.OAuth2Required {
		return context.TODO(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType_AMF, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}

func (c *AMFContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !c.OAuth2Required {
		logger.UtilLog.Debugf("AMFContext::AuthorizationCheck: OAuth2 not required\n")
		return nil
	}

	logger.UtilLog.Debugf("AMFContext::AuthorizationCheck: token[%s] serviceName[%s]\n", token, serviceName)
	return oauth.VerifyOAuth(token, string(serviceName), c.NrfCertPem)
}
