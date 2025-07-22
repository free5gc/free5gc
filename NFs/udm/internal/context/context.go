package context

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nnrf_NFDiscovery "github.com/free5gc/openapi/nrf/NFDiscovery"
	"github.com/free5gc/openapi/oauth"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/udm/pkg/factory"
	"github.com/free5gc/udm/pkg/suci"
	"github.com/free5gc/util/idgenerator"
)

var udmContext = UDMContext{}

const (
	LocationUriAmf3GppAccessRegistration int = iota
	LocationUriAmfNon3GppAccessRegistration
	LocationUriSmfRegistration
	LocationUriSdmSubscription
	LocationUriSharedDataSubscription
)

func Init() {
	GetSelf().NfService = make(map[models.ServiceName]models.NrfNfManagementNfService)
	GetSelf().EeSubscriptionIDGenerator = idgenerator.NewGenerator(1, math.MaxInt32)
	InitUdmContext(GetSelf())
}

type NFContext interface {
	AuthorizationCheck(token string, serviceName models.ServiceName) error
}

var _ NFContext = &UDMContext{}

type UDMContext struct {
	NfId                           string
	GroupId                        string
	SBIPort                        int
	RegisterIPv4                   string // IP register to NRF
	BindingIPv4                    string
	UriScheme                      models.UriScheme
	NfService                      map[models.ServiceName]models.NrfNfManagementNfService
	NFDiscoveryClient              *Nnrf_NFDiscovery.APIClient
	UdmUePool                      sync.Map // map[supi]*UdmUeContext
	NrfUri                         string
	NrfCertPem                     string
	GpsiSupiList                   models.IdentityData
	SharedSubsDataMap              map[string]models.UdmSdmSharedData // sharedDataIds as key
	SubscriptionOfSharedDataChange sync.Map                           // subscriptionID as key
	SuciProfiles                   []suci.SuciProfile
	EeSubscriptionIDGenerator      *idgenerator.IDGenerator
	OAuth2Required                 bool
}

type UdmUeContext struct {
	Supi                              string
	Gpsi                              string
	ExternalGroupID                   string
	Nssai                             *models.Nssai
	Amf3GppAccessRegistration         *models.Amf3GppAccessRegistration
	AmfNon3GppAccessRegistration      *models.AmfNon3GppAccessRegistration
	AccessAndMobilitySubscriptionData *models.AccessAndMobilitySubscriptionData
	SmfSelSubsData                    *models.SmfSelectionSubscriptionData
	UeCtxtInSmfData                   *models.UeContextInSmfData
	TraceDataResponse                 models.TraceDataResponse
	TraceData                         *models.TraceData
	SessionManagementSubsData         map[string]models.SessionManagementSubscriptionData
	SubsDataSets                      *models.UdmSdmSubscriptionDataSets
	SubscribeToNotifChange            map[string]*models.SdmSubscription
	SubscribeToNotifSharedDataChange  *models.SdmSubscription
	PduSessionID                      string
	UdrUri                            string
	UdmSubsToNotify                   map[string]*models.SubscriptionDataSubscriptions
	EeSubscriptions                   map[string]*models.UdmEeEeSubscription // subscriptionID as key
	amSubsDataLock                    sync.Mutex
	smfSelSubsDataLock                sync.Mutex
	SmSubsDataLock                    sync.RWMutex
}

func (ue *UdmUeContext) Init() {
	ue.UdmSubsToNotify = make(map[string]*models.SubscriptionDataSubscriptions)
	ue.EeSubscriptions = make(map[string]*models.UdmEeEeSubscription)
	ue.SubscribeToNotifChange = make(map[string]*models.SdmSubscription)
}

type UdmNFContext struct {
	SubscriptionID                   string
	SubscribeToNotifChange           *models.SdmSubscription // SubscriptionID as key
	SubscribeToNotifSharedDataChange *models.SdmSubscription // SubscriptionID as key
}

func InitUdmContext(context *UDMContext) {
	config := factory.UdmConfig
	logger.UtilLog.Info("udmconfig Info: Version[", config.Info.Version, "] Description[", config.Info.Description, "]")
	configuration := config.Configuration
	udmContext.NfId = uuid.New().String()
	sbi := configuration.Sbi
	udmContext.UriScheme = ""
	udmContext.SBIPort = factory.UdmSbiDefaultPort
	udmContext.RegisterIPv4 = factory.UdmSbiDefaultIPv4
	if sbi != nil {
		if sbi.Scheme != "" {
			udmContext.UriScheme = models.UriScheme(sbi.Scheme)
		}
		if sbi.RegisterIPv4 != "" {
			udmContext.RegisterIPv4 = sbi.RegisterIPv4
		}
		if sbi.Port != 0 {
			udmContext.SBIPort = sbi.Port
		}

		udmContext.BindingIPv4 = os.Getenv(sbi.BindingIPv4)
		if udmContext.BindingIPv4 != "" {
			logger.UtilLog.Info("Parsing ServerIPv4 address from ENV Variable.")
		} else {
			udmContext.BindingIPv4 = sbi.BindingIPv4
			if udmContext.BindingIPv4 == "" {
				logger.UtilLog.Warn("Error parsing ServerIPv4 address as string. Using the 0.0.0.0 address as default.")
				udmContext.BindingIPv4 = "0.0.0.0"
			}
		}
	}
	udmContext.NrfUri = configuration.NrfUri
	context.NrfCertPem = configuration.NrfCertPem
	servingNameList := configuration.ServiceNameList

	udmContext.SuciProfiles = configuration.SuciProfiles

	udmContext.InitNFService(servingNameList, config.Info.Version)
}

func (context *UDMContext) ManageSmData(smDatafromUDR []models.SessionManagementSubscriptionData, snssaiFromReq string,
	dnnFromReq string) (mp map[string]models.SessionManagementSubscriptionData, ind string,
	Dnns []models.DnnConfiguration, allDnns []map[string]models.DnnConfiguration,
) {
	smDataMap := make(map[string]models.SessionManagementSubscriptionData)
	sNssaiList := make([]string, len(smDatafromUDR))
	// to obtain all DNN configurations identified by "dnn" for all network slices where such DNN is available
	AllDnnConfigsbyDnn := make([]models.DnnConfiguration, len(sNssaiList))
	// to obtain all DNN configurations for all network slice(s)
	AllDnns := make([]map[string]models.DnnConfiguration, len(smDatafromUDR))
	var snssaikey string // Required snssai to obtain all DNN configurations

	for idx, smSubscriptionData := range smDatafromUDR {
		singleNssaiStr := openapi.MarshToJsonString(smSubscriptionData.SingleNssai)[0]
		smDataMap[singleNssaiStr] = smSubscriptionData
		// sNssaiList = append(sNssaiList, singleNssaiStr)
		AllDnns[idx] = smSubscriptionData.DnnConfigurations
		if strings.Contains(singleNssaiStr, snssaiFromReq) {
			snssaikey = singleNssaiStr
		}

		if _, ok := smSubscriptionData.DnnConfigurations[dnnFromReq]; ok {
			AllDnnConfigsbyDnn = append(AllDnnConfigsbyDnn, smSubscriptionData.DnnConfigurations[dnnFromReq])
		}
	}

	return smDataMap, snssaikey, AllDnnConfigsbyDnn, AllDnns
}

// HandleGetSharedData related functions
func MappingSharedData(sharedDatafromUDR []models.UdmSdmSharedData) (mp map[string]models.UdmSdmSharedData) {
	sharedSubsDataMap := make(map[string]models.UdmSdmSharedData)
	for i := 0; i < len(sharedDatafromUDR); i++ {
		sharedSubsDataMap[sharedDatafromUDR[i].SharedDataId] = sharedDatafromUDR[i]
	}
	return sharedSubsDataMap
}

func ObtainRequiredSharedData(Sharedids []string, response []models.UdmSdmSharedData) (
	sharedDatas []models.UdmSdmSharedData,
) {
	sharedSubsDataMap := MappingSharedData(response)
	Allkeys := make([]string, len(sharedSubsDataMap))
	MatchedKeys := make([]string, len(Sharedids))
	counter := 0
	for k := range sharedSubsDataMap {
		Allkeys = append(Allkeys, k)
	}

	for j := 0; j < len(Sharedids); j++ {
		for i := 0; i < len(Allkeys); i++ {
			if strings.Contains(Allkeys[i], Sharedids[j]) {
				MatchedKeys[counter] = Allkeys[i]
			}
		}
		counter += 1
	}

	shared_Data := make([]models.UdmSdmSharedData, len(MatchedKeys))
	if len(MatchedKeys) != 1 {
		for i := 0; i < len(MatchedKeys); i++ {
			shared_Data[i] = sharedSubsDataMap[MatchedKeys[i]]
		}
	} else {
		shared_Data[0] = sharedSubsDataMap[MatchedKeys[0]]
	}
	return shared_Data
}

// Returns the  SUPI from the SUPI list (SUPI list contains either a SUPI or a NAI)
func GetCorrespondingSupi(list models.IdentityData) (id string) {
	var identifier string
	for i := 0; i < len(list.SupiList); i++ {
		if strings.Contains(list.SupiList[i], "imsi") {
			identifier = list.SupiList[i]
		}
	}
	return identifier
}

// functions related to Retrieval of multiple datasets(GetSupi)
func (context *UDMContext) CreateSubsDataSetsForUe(supi string, body models.UdmSdmSubscriptionDataSets) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.SubsDataSets = &body
}

// Functions related to the trace data configuration
func (context *UDMContext) CreateTraceDataforUe(supi string, body models.TraceData) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.TraceData = &body
}

// functions related to sdmSubscription (subscribe to notification of data change)
func (udmUeContext *UdmUeContext) CreateSubscriptiontoNotifChange(subscriptionID string, body *models.SdmSubscription) {
	if _, exist := udmUeContext.SubscribeToNotifChange[subscriptionID]; !exist {
		udmUeContext.SubscribeToNotifChange[subscriptionID] = body
	}
}

// TODO: this function has wrong UE pool key with subscriptionID
func (context *UDMContext) CreateSubstoNotifSharedData(subscriptionID string, body *models.SdmSubscription) {
	context.SubscriptionOfSharedDataChange.Store(subscriptionID, body)
}

// functions related UecontextInSmfData
func (context *UDMContext) CreateUeContextInSmfDataforUe(supi string, body models.UeContextInSmfData) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.UeCtxtInSmfData = &body
}

// functions for SmfSelectionSubscriptionData
func (context *UDMContext) CreateSmfSelectionSubsDataforUe(supi string, body models.SmfSelectionSubscriptionData) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.SmfSelSubsData = &body
}

// SetSmfSelectionSubsData ... functions to set SmfSelectionSubscriptionData
func (udmUeContext *UdmUeContext) SetSmfSelectionSubsData(smfSelSubsData *models.SmfSelectionSubscriptionData) {
	udmUeContext.smfSelSubsDataLock.Lock()
	defer udmUeContext.smfSelSubsDataLock.Unlock()
	udmUeContext.SmfSelSubsData = smfSelSubsData
}

// SetSMSubsData ... functions to set SessionManagementSubsData
func (udmUeContext *UdmUeContext) SetSMSubsData(smSubsData map[string]models.SessionManagementSubscriptionData) {
	udmUeContext.SmSubsDataLock.Lock()
	defer udmUeContext.SmSubsDataLock.Unlock()
	udmUeContext.SessionManagementSubsData = smSubsData
}

func (context *UDMContext) NewUdmUe(supi string) *UdmUeContext {
	ue := new(UdmUeContext)
	ue.Init()
	ue.Supi = supi
	context.UdmUePool.Store(supi, ue)
	return ue
}

func (context *UDMContext) UdmUeFindBySupi(supi string) (*UdmUeContext, bool) {
	if value, ok := context.UdmUePool.Load(supi); ok {
		return value.(*UdmUeContext), ok
	} else {
		return nil, false
	}
}

func (context *UDMContext) UdmUeFindByGpsi(gpsi string) (*UdmUeContext, bool) {
	var ue *UdmUeContext
	ok := false
	context.UdmUePool.Range(func(key, value interface{}) bool {
		candidate := value.(*UdmUeContext)
		if candidate.Gpsi == gpsi {
			ue = candidate
			ok = true
			return false
		}
		return true
	})
	return ue, ok
}

// Function to create the AccessAndMobilitySubscriptionData for Ue
func (context *UDMContext) CreateAccessMobilitySubsDataForUe(supi string,
	body models.AccessAndMobilitySubscriptionData,
) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.AccessAndMobilitySubscriptionData = &body
}

// Function to set the AccessAndMobilitySubscriptionData for Ue
func (udmUeContext *UdmUeContext) SetAMSubsriptionData(amData *models.AccessAndMobilitySubscriptionData) {
	udmUeContext.amSubsDataLock.Lock()
	defer udmUeContext.amSubsDataLock.Unlock()
	udmUeContext.AccessAndMobilitySubscriptionData = amData
}

func (context *UDMContext) UdmAmf3gppRegContextExists(supi string) bool {
	if ue, ok := context.UdmUeFindBySupi(supi); ok {
		return ue.Amf3GppAccessRegistration != nil
	} else {
		return false
	}
}

func (context *UDMContext) UdmAmfNon3gppRegContextExists(supi string) bool {
	if ue, ok := context.UdmUeFindBySupi(supi); ok {
		return ue.AmfNon3GppAccessRegistration != nil
	} else {
		return false
	}
}

func (context *UDMContext) UdmSmfRegContextNotExists(supi string) bool {
	if ue, ok := context.UdmUeFindBySupi(supi); ok {
		return ue.PduSessionID == ""
	} else {
		return true
	}
}

func (context *UDMContext) CreateAmf3gppRegContext(supi string, body models.Amf3GppAccessRegistration) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.Amf3GppAccessRegistration = &body
}

func (context *UDMContext) CreateAmfNon3gppRegContext(supi string, body models.AmfNon3GppAccessRegistration) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	ue.AmfNon3GppAccessRegistration = &body
}

func (context *UDMContext) CreateSmfRegContext(supi string, pduSessionID string) {
	ue, ok := context.UdmUeFindBySupi(supi)
	if !ok {
		ue = context.NewUdmUe(supi)
	}
	if ue.PduSessionID == "" {
		ue.PduSessionID = pduSessionID
	}
}

func (context *UDMContext) GetAmf3gppRegContext(supi string) *models.Amf3GppAccessRegistration {
	if ue, ok := context.UdmUeFindBySupi(supi); ok {
		return ue.Amf3GppAccessRegistration
	} else {
		return nil
	}
}

func (context *UDMContext) GetAmfNon3gppRegContext(supi string) *models.AmfNon3GppAccessRegistration {
	if ue, ok := context.UdmUeFindBySupi(supi); ok {
		return ue.AmfNon3GppAccessRegistration
	} else {
		return nil
	}
}

func (ue *UdmUeContext) GetLocationURI(types int) string {
	switch types {
	case LocationUriAmf3GppAccessRegistration:
		return GetSelf().GetIPv4Uri() + factory.UdmUecmResUriPrefix + "/" + ue.Supi + "/registrations/amf-3gpp-access"
	case LocationUriAmfNon3GppAccessRegistration:
		return GetSelf().GetIPv4Uri() + factory.UdmUecmResUriPrefix + "/" + ue.Supi + "/registrations/amf-non-3gpp-access"
	case LocationUriSmfRegistration:

		return GetSelf().GetIPv4Uri() +
			factory.UdmUecmResUriPrefix + "/" + ue.Supi + "/registrations/smf-registrations/" + ue.PduSessionID
	}
	return ""
}

func (ue *UdmUeContext) GetLocationURI2(types int, supi string) string {
	switch types {
	case LocationUriSharedDataSubscription:
		// return GetSelf().GetIPv4Uri() + UdmSdmResUriPrefix +"/shared-data-subscriptions/" + nf.SubscriptionID
	case LocationUriSdmSubscription:
		return GetSelf().GetIPv4Uri() + factory.UdmSdmResUriPrefix + "/" + supi + "/sdm-subscriptions/"
	}
	return ""
}

func (ue *UdmUeContext) SameAsStoredGUAMI3gpp(inGuami models.Guami) bool {
	if ue.Amf3GppAccessRegistration == nil {
		return false
	}
	ug := ue.Amf3GppAccessRegistration.Guami
	if ug != nil {
		if (ug.PlmnId == nil) == (inGuami.PlmnId == nil) {
			if ug.PlmnId != nil && ug.PlmnId.Mcc == inGuami.PlmnId.Mcc && ug.PlmnId.Mnc == inGuami.PlmnId.Mnc {
				if ug.AmfId == inGuami.AmfId {
					return true
				}
			}
		}
	}
	return false
}

func (ue *UdmUeContext) SameAsStoredGUAMINon3gpp(inGuami models.Guami) bool {
	if ue.AmfNon3GppAccessRegistration == nil {
		return false
	}
	ug := ue.AmfNon3GppAccessRegistration.Guami
	if ug != nil {
		if (ug.PlmnId == nil) == (inGuami.PlmnId == nil) {
			if ug.PlmnId != nil && ug.PlmnId.Mcc == inGuami.PlmnId.Mcc && ug.PlmnId.Mnc == inGuami.PlmnId.Mnc {
				if ug.AmfId == inGuami.AmfId {
					return true
				}
			}
		}
	}
	return false
}

func (context *UDMContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
}

// GetSDMUri ... get subscriber data management service uri
func (context *UDMContext) GetSDMUri() string {
	return context.GetIPv4Uri() + factory.UdmSdmResUriPrefix
}

func (context *UDMContext) InitNFService(serviceName []string, version string) {
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	for index, nameString := range serviceName {
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

func (c *UDMContext) GetTokenCtx(serviceName models.ServiceName, targetNF models.NrfNfManagementNfType) (
	context.Context, *models.ProblemDetails, error,
) {
	if !c.OAuth2Required {
		return context.TODO(), nil, nil
	}
	return oauth.GetTokenCtx(models.NrfNfManagementNfType_UDM, targetNF,
		c.NfId, c.NrfUri, string(serviceName))
}

func GetSelf() *UDMContext {
	return &udmContext
}

func (context *UDMContext) AuthorizationCheck(token string, serviceName models.ServiceName) error {
	if !context.OAuth2Required {
		logger.UtilLog.Debugf("UDMContext::AuthorizationCheck: OAuth2 not required\n")
		return nil
	}
	logger.UtilLog.Debugf("UDMContext::AuthorizationCheck: token[%s] serviceName[%s]\n", token, serviceName)
	err := oauth.VerifyOAuth(token, string(serviceName), context.NrfCertPem)
	if err != nil {
		return err
	}
	return nil
}
