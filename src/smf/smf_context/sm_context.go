package smf_context

import (
	"fmt"
	"github.com/google/uuid"
	"free5gc/lib/Namf_Communication"
	"free5gc/lib/Npcf_SMPolicyControl"
	"free5gc/lib/openapi/models"
	"net"
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

	PDUAddress net.IP

	// Client
	SMPolicyClient      *Npcf_SMPolicyControl.APIClient
	CommunicationClient *Namf_Communication.APIClient

	AMFProfile models.NfProfile

	SMState SMState

	Tunnel *UPTunnel
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
