package smf_context

import (
	"context"
	"encoding/json"
	"fmt"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/Nnrf_NFDiscovery"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/lib/nas/nasConvert"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/smf/logger"
	"net"
	"net/http"

	"github.com/google/uuid"
)

var smContextPool map[string]*SMContext
var canonicalRef map[string]string

var smContextCount uint64

type SMState int

const (
	PDUSessionInactive SMState = 0
	PDUSessionActive   SMState = 1
)

func init() {
	smContextPool = make(map[string]*SMContext)
	canonicalRef = make(map[string]string)
}

func GetSMContextCount() uint64 {
	smContextCount++
	return smContextCount
}

type SMContext struct {
	Ref string

	LocalSEID  uint64
	RemoteSEID uint64

	UnauthenticatedSupi bool
	// SUPI or PEI
	Supi           string
	Pei            string
	Identifier     string
	Gpsi           string
	PDUSessionID   int32
	Dnn            string
	Snssai         *models.Snssai
	HplmnSnssai    *models.Snssai
	ServingNetwork *models.PlmnId
	ServingNfId    string

	UpCnxState models.UpCnxState

	AnType          models.AccessType
	RatType         models.RatType
	PresenceInLadn  models.PresenceState
	UeLocation      *models.UserLocation
	UeTimeZone      string
	AddUeLocation   *models.UserLocation
	OldPduSessionId int32
	HoState         models.HoState

	PDUAddress             net.IP
	SelectedPDUSessionType uint8

	DnnConfiguration models.DnnConfiguration
	SessionRule      models.SessionRule

	// Client
	SMPolicyClient      *Npcf_SMPolicyControl.APIClient
	CommunicationClient *Namf_Communication.APIClient

	AMFProfile         models.NfProfile
	SelectedPCFProfile models.NfProfile
	SmStatusNotifyUri  string

	SMState SMState

	Tunnel    *UPTunnel
	BPManager *BPManager
}

func canonicalName(identifier string, pduSessID int32) (canonical string) {
	return fmt.Sprintf("%s-%d", identifier, pduSessID)
}

func ResolveRef(identifier string, pduSessID int32) (ref string, err error) {
	ref, ok := canonicalRef[canonicalName(identifier, pduSessID)]
	if ok {
		err = nil
	} else {
		err = fmt.Errorf(
			"UE '%s' - PDUSessionID '%d' not found in SMContext", identifier, pduSessID)
	}
	return
}

func NewSMContext(identifier string, pduSessID int32) (smContext *SMContext) {
	smContext = new(SMContext)
	// Create Ref and identifier
	smContext.Ref = uuid.New().URN()
	smContextPool[smContext.Ref] = smContext
	canonicalRef[canonicalName(identifier, pduSessID)] = smContext.Ref

	smContext.Identifier = identifier
	smContext.PDUSessionID = pduSessID
	smContext.LocalSEID = GetSMContextCount()
	return smContext
}

func GetSMContext(ref string) (smContext *SMContext) {
	smContext = smContextPool[ref]
	return smContext
}

func RemoveSMContext(ref string) {
	delete(smContextPool, ref)
}

func GetSMContextBySEID(SEID uint64) *SMContext {
	for _, smCtx := range smContextPool {
		if smCtx.LocalSEID == SEID {
			return smCtx
		}
	}
	return nil
}

func (smContext *SMContext) SetCreateData(createData *models.SmContextCreateData) {

	smContext.Gpsi = createData.Gpsi
	smContext.Supi = createData.Supi
	smContext.Dnn = createData.Dnn
	smContext.Snssai = createData.SNssai
	smContext.HplmnSnssai = createData.HplmnSnssai
	smContext.ServingNetwork = createData.ServingNetwork
	smContext.AnType = createData.AnType
	smContext.RatType = createData.RatType
	smContext.PresenceInLadn = createData.PresenceInLadn
	smContext.UeLocation = createData.UeLocation
	smContext.UeTimeZone = createData.UeTimeZone
	smContext.AddUeLocation = createData.AddUeLocation
	smContext.OldPduSessionId = createData.OldPduSessionId
	smContext.ServingNfId = createData.ServingNfId
}

func (smContext *SMContext) BuildCreatedData() (createdData *models.SmContextCreatedData) {
	createdData = new(models.SmContextCreatedData)
	createdData.SNssai = smContext.Snssai
	return
}

func (smContext *SMContext) PDUAddressToNAS() (addr [12]byte, addrLen uint8) {
	copy(addr[:], smContext.PDUAddress)
	switch smContext.SelectedPDUSessionType {
	case nasMessage.PDUSessionTypeIPv4:
		addrLen = 4 + 1
	case nasMessage.PDUSessionTypeIPv6:
	case nasMessage.PDUSessionTypeIPv4IPv6:
		addrLen = 12 + 1
	}
	return
}

// PCFSelection will select PCF for this SM Context
func (smContext *SMContext) PCFSelection() (err error) {

	// Send NFDiscovery for find PCF
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{}

	rep, res, err := SMF_Self().NFDiscoveryClient.NFInstancesStoreApi.SearchNFInstances(context.TODO(), models.NfType_PCF, models.NfType_SMF, &localVarOptionals)
	if err != nil {
		return
	}

	if res != nil {
		if status := res.StatusCode; status != http.StatusOK {
			apiError := err.(common.GenericOpenAPIError)
			problemDetails := apiError.Model().(models.ProblemDetails)

			logger.CtxLog.Warningf("NFDiscovery PCF return status: %d\n", status)
			logger.CtxLog.Warningf("Detail: %v\n", problemDetails.Title)
		}
	}

	// Select PCF from available PCF

	smContext.SelectedPCFProfile = rep.NfInstances[0]

	SelectedPCFProfileString, _ := json.MarshalIndent(smContext.SelectedPCFProfile, "", "  ")
	logger.CtxLog.Tracef("Select PCF Profile: %s\n", SelectedPCFProfileString)

	// Create SMPolicyControl Client for this SM Context
	for _, service := range *smContext.SelectedPCFProfile.NfServices {
		if service.ServiceName == models.ServiceName_NPCF_SMPOLICYCONTROL {
			SmPolicyControlConf := Npcf_SMPolicyControl.NewConfiguration()
			SmPolicyControlConf.SetBasePath(service.ApiPrefix)
			smContext.SMPolicyClient = Npcf_SMPolicyControl.NewAPIClient(SmPolicyControlConf)
		}
	}

	return
}

func (smContext *SMContext) isAllowedPDUSessionType(nasPDUSessionType uint8) bool {
	dnnPDUSessionType := smContext.DnnConfiguration.PduSessionTypes
	if dnnPDUSessionType == nil {
		logger.CtxLog.Errorf("this SMContext[%s] has no subscription pdu session type info\n", smContext.Ref)
		return false
	}

	for _, allowedPDUSessionType := range smContext.DnnConfiguration.PduSessionTypes.AllowedSessionTypes {
		if allowedPDUSessionType == nasConvert.PDUSessionTypeToModels(nasPDUSessionType) {
			return true
		}
	}
	return false
}
