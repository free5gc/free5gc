package udm_context

import (
	// "fmt"
	"fmt"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/openapi"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/factory"
	"strconv"
	"strings"
)

var udmContext UDMContext
var udmUeContext UdmUeContext

const (
	LocationUriAmf3GppAccessRegistration int = iota
	LocationUriAmfNon3GppAccessRegistration
	LocationUriSmfRegistration
	LocationUriSdmSubscription
	LocationUriSharedDataSubscription
)

func Init() {
	UDM_Self().UdmUePool = make(map[string]*UdmUeContext)
	UDM_Self().NfService = make(map[models.ServiceName]models.NfService)
	InitUDMContext(&udmContext)
}

type UDMContext struct {
	Name              string
	NfId              string
	GroupId           string
	HttpIpv4Port      int
	HttpIPv4Address   string
	UriScheme         models.UriScheme
	NfService         map[models.ServiceName]models.NfService
	NFDiscoveryClient *Nnrf_NFDiscovery.APIClient
	UdmUePool         map[string]*UdmUeContext // supi as key
	UdmNFPool         map[string]*UdmNFContext // SubscriptionID as key
	NrfUri            string
	GpsiSupiList      models.IdentityData
	SharedSubsDataMap map[string]models.SharedData // sharedDataIds as key
	Keys              *factory.Keys
}

type UdmUeContext struct {
	Supi                              string
	GpsiFromReq                       string
	Nssai                             *models.Nssai
	Amf3GppAccessRegistration         *models.Amf3GppAccessRegistration
	AmfNon3GppAccessRegistration      *models.AmfNon3GppAccessRegistration
	AccessAndMobilitySubscriptionData *models.AccessAndMobilitySubscriptionData
	SmfSelSubsData                    *models.SmfSelectionSubscriptionData
	UeCtxtInSmfData                   *models.UeContextInSmfData
	TraceDataResponse                 models.TraceDataResponse
	TraceData                         *models.TraceData
	SessionManagementSubsData         map[string]models.SessionManagementSubscriptionData
	SubsDataSets                      *models.SubscriptionDataSets
	SubscribeToNotifChange            *models.SdmSubscription
	SubscribeToNotifSharedDataChange  *models.SdmSubscription
	PduSessionID                      string
	UdrUri                            string
	CreatedEeSubscription             models.CreatedEeSubscription
	UdmSubsToNotify                   map[string]*models.SubscriptionDataSubscriptions
}

// Functions related to EE services
func CreateEeSusbContext(ueId string, body models.CreatedEeSubscription) {
	udmUe := UDM_Self().UdmUePool[ueId]
	if udmUe == nil {
		udmUe = CreateUdmUe(ueId)
	}
	udmUe.CreatedEeSubscription = body
}

type UdmNFContext struct {
	SubscriptionID                   string
	SubscribeToNotifChange           *models.SdmSubscription // SubscriptionID as key
	SubscribeToNotifSharedDataChange *models.SdmSubscription // SubscriptionID as key
}

func GetUdmProfileAHNPublicKey() string {
	return udmContext.Keys.UdmProfileAHNPublicKey
}

func GetUdmProfileAHNPrivateKey() string {
	return udmContext.Keys.UdmProfileAHNPrivateKey
}

func GetUdmProfileBHNPublicKey() string {
	return udmContext.Keys.UdmProfileBHNPublicKey
}

func GetUdmProfileBHNPrivateKey() string {
	return udmContext.Keys.UdmProfileBHNPrivateKey
}

func ManageSmData(smDatafromUDR []models.SessionManagementSubscriptionData, snssaiFromReq string, dnnFromReq string) (mp map[string]models.SessionManagementSubscriptionData, ind string,
	Dnns []models.DnnConfiguration, allDnns []map[string]models.DnnConfiguration) {

	smDataMap := make(map[string]models.SessionManagementSubscriptionData)
	sNssaiList := make([]string, len(smDatafromUDR))
	AllDnnConfigsbyDnn := make([]models.DnnConfiguration, 1, len(sNssaiList)) // to obtain all DNN configurations identified by "dnn" for all network slices where such DNN is available
	AllDnns := make([]map[string]models.DnnConfiguration, len(smDatafromUDR)) // to obtain all DNN configurations for all network slice(s)
	var snssaikey string                                                      // Required snssai to obtain all DNN configurations

	for idx, smSubscriptionData := range smDatafromUDR {
		singleNssaiStr := openapi.MarshToJsonString(smSubscriptionData.SingleNssai)[0]
		smDataMap[singleNssaiStr] = smSubscriptionData
		sNssaiList = append(sNssaiList, singleNssaiStr)
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

// Check if SessionMgtSubsData context exists
func UdmUeSessionMgtSubsDataExisting(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe.SessionManagementSubsData != nil {
		return udmUe.SessionManagementSubsData == nil
	}
	return true
}

// HandleGetSharedData related functions
func MappingSharedData(sharedDatafromUDR []models.SharedData) (mp map[string]models.SharedData) {
	sharedSubsDataMap := make(map[string]models.SharedData)
	for i := 0; i < len(sharedDatafromUDR); i++ {
		sharedSubsDataMap[sharedDatafromUDR[i].SharedDataId] = sharedDatafromUDR[i]
	}
	return sharedSubsDataMap
}

func ObtainRequiredSharedData(Sharedids []string, response []models.SharedData) (sharedDatas []models.SharedData) {
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

	shared_Data := make([]models.SharedData, len(MatchedKeys))
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
func CreateSubsDataSetsForUe(Supi string, body models.SubscriptionDataSets) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.SubsDataSets = &body
}

func UdmUeSubsDataSetsExisting(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.SubsDataSets == nil
	}
	return true
}

// Functions related to the trace data configuration
func CreateTraceDataforUe(Supi string, body models.TraceData) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.TraceData = &body
}

func UdmUeTraceDataExisting(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.TraceData == nil
	}
	return true
}

// functions related to sdmSubscription (subscribe to notification of data change)
func CreateSubscriptiontoNotifChange(SubscriptionID string, body *models.SdmSubscription) {
	udmUe := UDM_Self().UdmUePool[SubscriptionID]
	if udmUe == nil {
		udmUe = CreateUdmUe(SubscriptionID)
	}
	udmUe.SubscribeToNotifChange = body
}

func UdmNfCntxtNotExisting(SubscriptionID string) bool {
	udmNf := UDM_Self().UdmNFPool[SubscriptionID]
	if udmNf != nil {
		return udmNf.SubscribeToNotifChange == nil
	}
	return true
}

func CreateSubstoNotifSharedData(SubscriptionID string, body *models.SdmSubscription) {
	udmUe := UDM_Self().UdmUePool[SubscriptionID]
	if udmUe == nil {
		udmUe = CreateUdmUe(SubscriptionID)
	}
	udmUe.SubscribeToNotifSharedDataChange = body
}

func UdmNfCntxtSharedDataExisting(SubscriptionID string) bool {
	udmNf := UDM_Self().UdmNFPool[SubscriptionID]
	if udmNf != nil {
		return udmNf.SubscribeToNotifSharedDataChange == nil
	}
	return true
}

// functions related UecontextInSmfData
func CreateUeContextInSmfDataforUe(Supi string, body models.UeContextInSmfData) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.UeCtxtInSmfData = &body
}

func UdmUeCtxtInSmfDataExisting(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.UeCtxtInSmfData == nil
	}
	return true
}

// functions for SmfSelectionSubscriptionData
func CreateSmfSelectionSubsDatadforUe(Supi string, body models.SmfSelectionSubscriptionData) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.SmfSelSubsData = &body
}

func UdmueSmfSelSubsDataNotExisting(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.SmfSelSubsData == nil
	}
	return true
}

func CreateUdmUe(Supi string) (udmUe *UdmUeContext) {
	udmUe = new(UdmUeContext)
	udmUe.Supi = Supi
	UDM_Self().UdmUePool[Supi] = udmUe
	return
}
func CreateUdmNf(SubscriptionID string) (udmNf *UdmNFContext) {
	udmNf = new(UdmNFContext)
	udmNf.SubscriptionID = SubscriptionID
	UDM_Self().UdmNFPool[SubscriptionID] = udmNf
	return
}

// Check if Access and Mobility Subscription Data context exists
func UdmUeAccessMobilitySubsDataExisting(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe.AccessAndMobilitySubscriptionData != nil {
		return udmUe.AccessAndMobilitySubscriptionData == nil
	}
	return true
}

// Function to create the AccessAndMobilitySubscriptionData for Ue
func CreateAccessMobilitySubsDataForUe(Supi string, body models.AccessAndMobilitySubscriptionData) {
	UdmUe := UDM_Self().UdmUePool[Supi]
	if UdmUe == nil {
		UdmUe = CreateUdmUe(Supi)
	}
	UdmUe.AccessAndMobilitySubscriptionData = &body
}

func GetAccessMobilitySubsDataForUe(Supi string) *models.AccessAndMobilitySubscriptionData {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.AccessAndMobilitySubscriptionData
	}
	return nil
}

func UdmAmf3gppRegContextExists(Supi string) bool {
	if udmUe := UDM_Self().UdmUePool[Supi]; udmUe != nil {
		return udmUe.Amf3GppAccessRegistration != nil
	}
	return false
}

func UdmAmfNon3gppRegContextExists(Supi string) bool {
	if udmUe := UDM_Self().UdmUePool[Supi]; udmUe != nil {
		return udmUe.AmfNon3GppAccessRegistration != nil
	}
	return false
}

func UdmSmfRegContextNotExists(Supi string) bool {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.PduSessionID == ""
	}
	return true
}

func CreateAmf3gppRegContext(Supi string, body models.Amf3GppAccessRegistration) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.Amf3GppAccessRegistration = &body
}

func CreateAmfNon3gppRegContext(Supi string, body models.AmfNon3GppAccessRegistration) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	udmUe.AmfNon3GppAccessRegistration = &body
}

func CreateSmfRegContext(Supi string, pduSessionID string) {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe == nil {
		udmUe = CreateUdmUe(Supi)
	}
	if udmUe.PduSessionID == "" {
		udmUe.PduSessionID = pduSessionID
	}
}

func GetAmf3gppRegContext(Supi string) *models.Amf3GppAccessRegistration {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.Amf3GppAccessRegistration
	}
	return nil
}

func GetAmfNon3gppRegContext(Supi string) *models.AmfNon3GppAccessRegistration {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.AmfNon3GppAccessRegistration
	}
	return nil
}
func GetSmfRegContext(Supi string) string {
	udmUe := UDM_Self().UdmUePool[Supi]
	if udmUe != nil {
		return udmUe.PduSessionID
	}
	return ""
}

func (ue *UdmUeContext) GetLocationURI(types int) string {
	switch types {
	case LocationUriAmf3GppAccessRegistration:
		return UDM_Self().GetIPv4Uri() + "/nudm-uecm/v1/" + ue.Supi + "/registrations/amf-3gpp-access"
	case LocationUriAmfNon3GppAccessRegistration:
		return UDM_Self().GetIPv4Uri() + "/nudm-uecm/v1/" + ue.Supi + "/registrations/amf-non-3gpp-access"
	case LocationUriSmfRegistration:
		return UDM_Self().GetIPv4Uri() + "/nudm-uecm/v1/" + ue.Supi + "/registrations/smf-registrations/" + ue.PduSessionID
	}
	return ""
}

func (ue *UdmUeContext) GetLocationURI2(types int, supi string) string {
	switch types {
	case LocationUriSharedDataSubscription:
		// return UDM_Self().GetIPv4Uri() + "/nudm-sdm/v1/shared-data-subscriptions/" + nf.SubscriptionID
	case LocationUriSdmSubscription:
		return UDM_Self().GetIPv4Uri() + "/nudm-sdm/v1/" + supi + "/sdm-subscriptions/"
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
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.HttpIPv4Address, context.HttpIpv4Port)
}

func (context *UDMContext) InitNFService(serviceName []string, version string) {
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

func UDM_Self() *UDMContext {
	return &udmContext
}

func UdmUe_self() *UdmUeContext {
	return &udmUeContext
}
