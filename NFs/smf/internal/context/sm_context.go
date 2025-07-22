package context

import (
	"fmt"
	"math"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/free5gc/nas/nasConvert"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/ngap/ngapType"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/pkg/factory"
	"github.com/free5gc/util/idgenerator"
)

var (
	smContextPool    sync.Map
	canonicalRef     sync.Map
	seidSMContextMap sync.Map
)

type DLForwardingType int

const (
	IndirectForwarding DLForwardingType = iota
	DirectForwarding
	NoForwarding
)

type UrrType int

// Reserved URR report for ID = 0 ~ 6
const (
	N3N6_MBQE_URR UrrType = iota
	N3N6_MAQE_URR
	N3N9_MBQE_URR
	N3N9_MAQE_URR
	N9N6_MBQE_URR
	N9N6_MAQE_URR
	NOT_FOUND_URR
)

func (t UrrType) String() string {
	urrTypeList := []string{"N3N6_MBQE", "N3N6_MAQE", "N3N9_MBQE", "N3N9_MAQE", "N9N6_MBQE", "N9N6_MAQE"}
	return urrTypeList[t]
}

func (t UrrType) IsBeforeQos() bool {
	urrTypeList := []bool{true, false, true, false, true, false}
	return urrTypeList[t]
}

func (t UrrType) Direct() string {
	urrTypeList := []string{"N3N6", "N3N6", "N3N9", "N3N9", "N9N6", "N9N6"}
	return urrTypeList[t]
}

var smContextCount uint64

type SMContextState uint32

const (
	InActive SMContextState = iota
	ActivePending
	Active
	InActivePending
	ModificationPending
	PFCPModification
)

const DefaultPrecedence uint32 = 255

func init() {
}

func GetSMContextCount() uint64 {
	return atomic.AddUint64(&smContextCount, 1)
}

type EventExposureNotification struct {
	*models.NsmfEventExposureNotification

	Uri string
}

type UsageReport struct {
	UrrId uint32
	UpfId string

	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64

	TotalPktNum    uint64
	UplinkPktNum   uint64
	DownlinkPktNum uint64

	ReportTpye models.ChfConvergedChargingTriggerType
}

var TeidGenerator *idgenerator.IDGenerator

type SMContext struct {
	*models.SmfPduSessionSmContextCreateData

	Ref string

	LocalSEID  uint64
	RemoteSEID uint64

	UnauthenticatedSupi bool

	Pei          string
	Identifier   string
	PDUSessionID int32

	LocalULTeid                   uint32
	LocalDLTeid                   uint32
	LocalULTeidForSplitPDUSession uint32
	LocalDLTeidForSplitPDUSession uint32
	NrdcIndicator                 bool

	UpCnxState models.UpCnxState

	HoState models.HoState

	SelectionParam         *UPFSelectionParams
	PDUAddress             net.IP
	UseStaticIP            bool
	SelectedPDUSessionType uint8

	DnnConfiguration models.DnnConfiguration

	SMPolicyID string

	// Handover related
	DLForwardingType         DLForwardingType
	DLDirectForwardingTunnel *ngapType.UPTransportLayerInformation
	IndirectForwardingTunnel *DataPath

	// UP Security support TS 29.502 R16 6.1.6.2.39
	UpSecurity                                                     *models.UpSecurity
	MaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink   models.MaxIntegrityProtectedDataRate
	MaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink models.MaxIntegrityProtectedDataRate
	// SMF verified UP security result of Xn-handover TS 33.501 6.6.1
	UpSecurityFromPathSwitchRequestSameAsLocalStored bool

	// Client
	CommunicationClientApiPrefix string

	AMFProfile         models.NrfNfDiscoveryNfProfile
	SelectedPCFProfile models.NrfNfDiscoveryNfProfile
	SelectedCHFProfile models.NrfNfDiscoveryNfProfile
	SmStatusNotifyUri  string

	Tunnel      *UPTunnel
	DCTunnel    *UPTunnel
	SelectedUPF *UPNode
	BPManager   *BPManager
	// NodeID(string form) to PFCP Session Context
	PFCPContext                         map[string]*PFCPSessionContext
	PDUSessionRelease_DUE_TO_DUP_PDU_ID bool

	DNNInfo *SnssaiSmfDnnInfo

	// SM Policy related
	PCCRules            map[string]*PCCRule
	DCPCCRules          map[string]*PCCRule
	SessionRules        map[string]*SessionRule
	TrafficControlDatas map[string]*TrafficControlData
	ChargingData        map[string]*models.ChargingData
	QosDatas            map[string]*models.QosData

	UpPathChgEarlyNotification map[string]*EventExposureNotification // Key: Uri+NotifId
	UpPathChgLateNotification  map[string]*EventExposureNotification // Key: Uri+NotifId
	DataPathToBeRemoved        map[int64]*DataPath                   // Key: pathID

	SelectedSessionRuleID string

	// QoS
	defRuleID               uint8
	QoSRuleIDGenerator      *idgenerator.IDGenerator
	PacketFilterIDGenerator *idgenerator.IDGenerator
	QFIGenerator            *idgenerator.IDGenerator
	PCCRuleIDToQoSRuleID    map[string]uint8
	qosDataToQFI            map[string]uint8
	PacketFilterIDToNASPFID map[string]uint8
	AMBRQerMap              map[uuid.UUID]uint32
	QerUpfMap               map[string]uint32
	AdditonalQosFlows       map[uint8]*QoSFlow // Key: qfi

	// URR
	UrrIDGenerator     *idgenerator.IDGenerator
	UrrIdMap           map[UrrType]uint32
	UrrUpfMap          map[string]*URR
	UrrReportTime      time.Duration
	UrrReportThreshold uint64
	// Cache the usage reports, sent from UPF
	// Those information will be included in CDR.
	UrrReports []UsageReport

	// Charging Related
	ChargingDataRef string
	// Each PDU session has a unique charging id
	ChargingID    int32
	RequestedUnit int32
	// key = urrid
	// All urr can map to a rating group
	// However, a rating group may map to more than one urr
	// e.g. In UL CL case, the rating group for recoreding PDU Session volume may map to two URR
	//		one is for PSA 1, the other is for PSA 2.
	ChargingInfo map[uint32]*ChargingInfo
	// NAS
	Pti                     uint8
	EstAcceptCause5gSMValue uint8

	// PCO Related
	ProtocolConfigurationOptions *ProtocolConfigurationOptions

	// State
	state SMContextState

	UeCmRegistered bool

	// Loggers
	Log *logrus.Entry

	// 5GSM Timers
	// T3591 is PDU SESSION MODIFICATION COMMAND timer
	T3591 *Timer
	// T3592 is PDU SESSION RELEASE COMMAND timer
	T3592 *Timer

	// lock
	SMLock sync.Mutex
}

func GenerateTEID() (uint32, error) {
	var id uint32
	if tmpID, err := TeidGenerator.Allocate(); err != nil {
		return 0, err
	} else {
		id = uint32(tmpID)
	}

	return id, nil
}

func ReleaseTEID(teid uint32) {
	TeidGenerator.FreeID(int64(teid))
}

func canonicalName(id string, pduSessID int32) string {
	return fmt.Sprintf("%s-%d", id, pduSessID)
}

func ResolveRef(id string, pduSessID int32) (string, error) {
	if value, ok := canonicalRef.Load(canonicalName(id, pduSessID)); ok {
		ref := value.(string)
		return ref, nil
	} else {
		return "", fmt.Errorf("UE[%s] - PDUSessionID[%d] not found in SMContext", id, pduSessID)
	}
}

func NewSMContext(id string, pduSessID int32) *SMContext {
	smContext := new(SMContext)
	// Create Ref and identifier
	smContext.Ref = uuid.New().URN()
	smContextPool.Store(smContext.Ref, smContext)
	canonicalRef.Store(canonicalName(id, pduSessID), smContext.Ref)

	smContext.Log = logger.PduSessLog.WithFields(logrus.Fields{
		logger.FieldSupi:         id,
		logger.FieldPDUSessionID: fmt.Sprintf("%d", pduSessID),
	})

	smContext.SetState(InActive)
	smContext.Identifier = id
	smContext.PDUSessionID = pduSessID
	smContext.PFCPContext = make(map[string]*PFCPSessionContext)
	smContext.LocalSEID = GetSMContextCount()

	// initialize SM Policy Data
	smContext.PCCRules = make(map[string]*PCCRule)
	smContext.DCPCCRules = make(map[string]*PCCRule)
	smContext.SessionRules = make(map[string]*SessionRule)
	smContext.TrafficControlDatas = make(map[string]*TrafficControlData)
	smContext.QosDatas = make(map[string]*models.QosData)
	smContext.UpPathChgEarlyNotification = make(map[string]*EventExposureNotification)
	smContext.UpPathChgLateNotification = make(map[string]*EventExposureNotification)
	smContext.DataPathToBeRemoved = make(map[int64]*DataPath)

	smContext.ProtocolConfigurationOptions = &ProtocolConfigurationOptions{}

	smContext.BPManager = NewBPManager(id)
	smContext.Tunnel = NewUPTunnel()
	smContext.DCTunnel = NewUPTunnel()

	smContext.QoSRuleIDGenerator = idgenerator.NewGenerator(1, 255)
	if defRuleID, err := smContext.QoSRuleIDGenerator.Allocate(); err != nil {
		logger.CtxLog.Warnln("Create Default QoS rule error ", err)
	} else {
		smContext.defRuleID = uint8(defRuleID)
	}

	smContext.PacketFilterIDGenerator = idgenerator.NewGenerator(1, 255)
	smContext.QFIGenerator = idgenerator.NewGenerator(2, 63) // 1 always reserve for default Qos
	smContext.PCCRuleIDToQoSRuleID = make(map[string]uint8)
	smContext.PacketFilterIDToNASPFID = make(map[string]uint8)
	smContext.qosDataToQFI = make(map[string]uint8)
	smContext.AMBRQerMap = make(map[uuid.UUID]uint32)
	smContext.QerUpfMap = make(map[string]uint32)
	smContext.AdditonalQosFlows = make(map[uint8]*QoSFlow)
	smContext.UrrIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
	smContext.UrrIdMap = make(map[UrrType]uint32)
	smContext.GenerateUrrId()
	smContext.UrrUpfMap = make(map[string]*URR)

	smContext.ChargingInfo = make(map[uint32]*ChargingInfo)
	smContext.ChargingID = GenerateChargingID()

	if factory.SmfConfig != nil &&
		factory.SmfConfig.Configuration != nil {
		smContext.UrrReportTime = time.Duration(factory.SmfConfig.Configuration.UrrPeriod) * time.Second
		smContext.UrrReportThreshold = factory.SmfConfig.Configuration.UrrThreshold
		logger.CtxLog.Infof("UrrPeriod: %v", smContext.UrrReportTime)
		logger.CtxLog.Infof("UrrThreshold: %d", smContext.UrrReportThreshold)
		if factory.SmfConfig.Configuration.RequestedUnit != 0 {
			smContext.RequestedUnit = factory.SmfConfig.Configuration.RequestedUnit
		} else {
			smContext.RequestedUnit = 1000
		}
	}

	var err error
	smContext.LocalDLTeid, err = GenerateTEID()
	if err != nil {
		return nil
	}

	smContext.LocalULTeid, err = GenerateTEID()
	if err != nil {
		return nil
	}

	smContext.LocalULTeidForSplitPDUSession, err = GenerateTEID()
	if err != nil {
		return nil
	}

	smContext.LocalDLTeidForSplitPDUSession, err = GenerateTEID()
	if err != nil {
		return nil
	}

	smContext.NrdcIndicator = false

	return smContext
}

// *** add unit test ***//
func GetSMContextByRef(ref string) *SMContext {
	var smCtx *SMContext
	if value, ok := smContextPool.Load(ref); ok {
		smCtx = value.(*SMContext)
	}
	return smCtx
}

func GetSMContextById(id string, pduSessID int32) *SMContext {
	var smCtx *SMContext
	ref, err := ResolveRef(id, pduSessID)
	if err != nil {
		return nil
	}
	if value, ok := smContextPool.Load(ref); ok {
		smCtx = value.(*SMContext)
	}
	return smCtx
}

// *** add unit test ***//
func RemoveSMContext(ref string) {
	var smContext *SMContext
	if value, ok := smContextPool.Load(ref); ok {
		smContext = value.(*SMContext)
	} else {
		return
	}

	if smContext.SelectedUPF != nil && smContext.PDUAddress != nil {
		logger.PduSessLog.Infof("UE[%s] PDUSessionID[%d] Release IP[%s]",
			smContext.Supi, smContext.PDUSessionID, smContext.PDUAddress.String())
		GetUserPlaneInformation().
			ReleaseUEIP(smContext.SelectedUPF, smContext.PDUAddress, smContext.UseStaticIP)
		smContext.SelectedUPF = nil
	}

	for _, pfcpSessionContext := range smContext.PFCPContext {
		seidSMContextMap.Delete(pfcpSessionContext.LocalSEID)
	}

	ReleaseTEID(smContext.LocalULTeid)
	ReleaseTEID(smContext.LocalDLTeid)
	ReleaseTEID(smContext.LocalULTeidForSplitPDUSession)
	ReleaseTEID(smContext.LocalDLTeidForSplitPDUSession)

	smContext.NrdcIndicator = false

	smContextPool.Delete(ref)
	canonicalRef.Delete(canonicalName(smContext.Supi, smContext.PDUSessionID))
	smContext.Log.Infof("smContext[%s] is deleted from pool", ref)
}

// *** add unit test ***//
func GetSMContextBySEID(seid uint64) *SMContext {
	if value, ok := seidSMContextMap.Load(seid); ok {
		smContext := value.(*SMContext)
		return smContext
	}
	return nil
}

func (smContext *SMContext) GenerateUrrId() {
	if id, err := smContext.UrrIDGenerator.Allocate(); err == nil {
		smContext.UrrIdMap[N3N6_MBQE_URR] = uint32(id)
	}
	if id, err := smContext.UrrIDGenerator.Allocate(); err == nil {
		smContext.UrrIdMap[N3N6_MAQE_URR] = uint32(id)
	}
	if id, err := smContext.UrrIDGenerator.Allocate(); err == nil {
		smContext.UrrIdMap[N9N6_MBQE_URR] = uint32(id)
	}
	if id, err := smContext.UrrIDGenerator.Allocate(); err == nil {
		smContext.UrrIdMap[N9N6_MAQE_URR] = uint32(id)
	}
	if id, err := smContext.UrrIDGenerator.Allocate(); err == nil {
		smContext.UrrIdMap[N3N9_MBQE_URR] = uint32(id)
	}
	if id, err := smContext.UrrIDGenerator.Allocate(); err == nil {
		smContext.UrrIdMap[N3N9_MAQE_URR] = uint32(id)
	}
}

func (smContext *SMContext) BuildCreatedData() *models.SmfPduSessionSmContextCreatedData {
	return &models.SmfPduSessionSmContextCreatedData{
		SNssai: smContext.SNssai,
	}
}

func (smContext *SMContext) SetState(state SMContextState) {
	oldState := SMContextState(atomic.LoadUint32((*uint32)(&smContext.state)))

	atomic.StoreUint32((*uint32)(&smContext.state), uint32(state))
	smContext.Log.Tracef("State[%s] -> State[%s]", oldState, state)
}

func (smContext *SMContext) CheckState(state SMContextState) bool {
	curState := SMContextState(atomic.LoadUint32((*uint32)(&smContext.state)))
	check := curState == state
	if !check {
		smContext.Log.Warnf("Unexpected state, expect: [%s], actual:[%s]", state, curState)
	}
	return check
}

func (smContext *SMContext) State() SMContextState {
	return SMContextState(atomic.LoadUint32((*uint32)(&smContext.state)))
}

func (smContext *SMContext) PDUAddressToNAS() ([12]byte, uint8) {
	var addr [12]byte
	var addrLen uint8
	copy(addr[:], smContext.PDUAddress)
	switch smContext.SelectedPDUSessionType {
	case nasMessage.PDUSessionTypeIPv4:
		var addrLenBuf uint8 = 4 + 1
		addrLen = addrLenBuf
	case nasMessage.PDUSessionTypeIPv6:
	case nasMessage.PDUSessionTypeIPv4IPv6:
		var addrLenBuf uint8 = 12 + 1
		addrLen = addrLenBuf
	}
	return addr, addrLen
}

func (smContext *SMContext) GetNodeIDByLocalSEID(seid uint64) pfcpType.NodeID {
	for _, pfcpCtx := range smContext.PFCPContext {
		if pfcpCtx.LocalSEID == seid {
			return pfcpCtx.NodeID
		}
	}

	return pfcpType.NodeID{}
}

func (smContext *SMContext) AllocateLocalSEIDForUPPath(path UPPath) {
	for _, upNode := range path {
		NodeIDtoIP := upNode.NodeID.ResolveNodeIdToIp().String()
		if _, exist := smContext.PFCPContext[NodeIDtoIP]; !exist {
			allocatedSEID := AllocateLocalSEID()

			smContext.PFCPContext[NodeIDtoIP] = &PFCPSessionContext{
				PDRs:      make(map[uint16]*PDR),
				NodeID:    upNode.NodeID,
				LocalSEID: allocatedSEID,
			}

			seidSMContextMap.Store(allocatedSEID, smContext)
		}
	}
}

func (smContext *SMContext) AllocateLocalSEIDForDataPath(dataPath *DataPath) {
	logger.PduSessLog.Traceln("In AllocateLocalSEIDForDataPath")
	for node := dataPath.FirstDPNode; node != nil; node = node.Next() {
		NodeIDtoIP := node.UPF.NodeID.ResolveNodeIdToIp().String()
		logger.PduSessLog.Traceln("NodeIDtoIP: ", NodeIDtoIP)
		if _, exist := smContext.PFCPContext[NodeIDtoIP]; !exist {
			allocatedSEID := AllocateLocalSEID()
			smContext.PFCPContext[NodeIDtoIP] = &PFCPSessionContext{
				PDRs:      make(map[uint16]*PDR),
				NodeID:    node.UPF.NodeID,
				LocalSEID: allocatedSEID,
			}

			seidSMContextMap.Store(allocatedSEID, smContext)
		}
	}
}

func (smContext *SMContext) PutPDRtoPFCPSession(nodeID pfcpType.NodeID, pdr *PDR) error {
	NodeIDtoIP := nodeID.ResolveNodeIdToIp().String()
	if pfcpSessCtx, exist := smContext.PFCPContext[NodeIDtoIP]; exist {
		smContext.Log.Tracef("PutPDRtoPFCPSession [%+v]", pdr)
		pfcpSessCtx.PDRs[pdr.PDRID] = pdr
	} else {
		return fmt.Errorf("can't find PFCPContext[%s] to put PDR(%d)", NodeIDtoIP, pdr.PDRID)
	}
	return nil
}

func (c *SMContext) findPSAandAllocUeIP(param *UPFSelectionParams) error {
	c.Log.Traceln("findPSAandAllocUeIP")
	if param == nil {
		return fmt.Errorf("UPFSelectionParams is nil")
	}

	upi := GetUserPlaneInformation()
	if GetSelf().ULCLSupport && CheckUEHasPreConfig(c.Supi) {
		groupName := GetULCLGroupNameFromSUPI(c.Supi)
		preConfigPathPool := GetUEDefaultPathPool(groupName)
		if preConfigPathPool != nil {
			selectedUPFName := ""
			selectedUPFName, c.PDUAddress, c.UseStaticIP = preConfigPathPool.SelectUPFAndAllocUEIPForULCL(
				upi, param)
			c.SelectedUPF = upi.UPFs[selectedUPFName]
		}
	} else {
		c.SelectedUPF, c.PDUAddress, c.UseStaticIP = upi.SelectUPFAndAllocUEIP(param)
		c.Log.Infof("Allocated PDUAdress[%s]", c.PDUAddress.String())
	}
	if c.PDUAddress == nil {
		return fmt.Errorf("fail to allocate PDU address, Selection Parameter: %s",
			param.String())
	}
	return nil
}

func (c *SMContext) AllocUeIP() error {
	c.SelectionParam = &UPFSelectionParams{
		Dnn: c.Dnn,
		SNssai: &SNssai{
			Sst: c.SNssai.Sst,
			Sd:  c.SNssai.Sd,
		},
	}

	if len(c.DnnConfiguration.StaticIpAddress) > 0 {
		staticIPConfig := c.DnnConfiguration.StaticIpAddress[0]
		if staticIPConfig.Ipv4Addr != "" {
			c.SelectionParam.PDUAddress = net.ParseIP(staticIPConfig.Ipv4Addr).To4()
		}
	}

	if err := c.findPSAandAllocUeIP(c.SelectionParam); err != nil {
		return err
	}
	return nil
}

// This function create ULCL data paths.
func (c *SMContext) SelectULCLDataPaths() error {
	if c.SelectionParam == nil || c.SelectedUPF == nil {
		return fmt.Errorf("SelectULCLDataPath err: SelectionParam or SelectedUPF is nil")
	}

	if GetSelf().ULCLSupport && CheckUEHasPreConfig(c.Supi) {
		c.Log.Infof("Has pre-config ULCL paths")
		uePreConfigPaths := GetUEPreConfigPaths(c.Supi, c.SelectedUPF.Name)
		for _, dp := range uePreConfigPaths.DataPathPool {
			if !dp.IsDefaultPath {
				c.Tunnel.AddDataPath(dp)
			}
		}
	}
	return nil
}

// This function create a data path to be default data path.
func (c *SMContext) SelectDefaultDataPath() error {
	if c.SelectionParam == nil || c.SelectedUPF == nil {
		return fmt.Errorf("SelectDefaultDataPath err: SelectionParam or SelectedUPF is nil")
	}

	defaultPath := c.Tunnel.DataPathPool.GetDefaultPath()
	if defaultPath != nil {
		// A default path already exists.
		// Use this one.
		c.Log.Infof("Has default path")
		defaultPath = c.Tunnel.DataPathPool.GetDefaultPath()
	} else if GetSelf().ULCLSupport && CheckUEHasPreConfig(c.Supi) {
		// Fallback on pre-config default path
		c.Log.Infof("Has pre-config default path")
		uePreConfigPaths := GetUEPreConfigPaths(c.Supi, c.SelectedUPF.Name)
		defaultPath = uePreConfigPaths.DataPathPool.GetDefaultPath()
		c.Tunnel.AddDataPath(defaultPath)
	} else {
		// UE has no pre-config path and default path
		// Use default route
		c.Log.Infof("Has no pre-config route. Has no default path")
		defaultUPPath := GetUserPlaneInformation().GetDefaultUserPlanePathByDNNAndUPF(
			c.SelectionParam, c.SelectedUPF)
		defaultPath = GenerateDataPath(defaultUPPath)
		if defaultPath != nil {
			defaultPath.IsDefaultPath = true
			c.Tunnel.AddDataPath(defaultPath)
		}
	}

	if defaultPath == nil {
		return fmt.Errorf("data path not found, Selection Parameter: %s",
			c.SelectionParam.String())
	}

	if !defaultPath.Activated {
		defaultPath.ActivateTunnelAndPDR(c, DefaultPrecedence)
	}

	return nil
}

func (c *SMContext) CreatePccRuleDataPath(pccRule *PCCRule,
	tcData *TrafficControlData, qosData *models.QosData,
	chgData *models.ChargingData,
) error {
	var targetRoute models.RouteToLocation
	if tcData != nil && len(tcData.RouteToLocs) > 0 {
		targetRoute = *tcData.RouteToLocs[0]
	}
	param := &UPFSelectionParams{
		Dnn: c.Dnn,
		SNssai: &SNssai{
			Sst: c.SNssai.Sst,
			Sd:  c.SNssai.Sd,
		},
		Dnai: targetRoute.Dnai,
	}
	createdUpPath := GetUserPlaneInformation().GetDefaultUserPlanePathByDNN(param)
	createdDataPath := GenerateDataPath(createdUpPath)
	if createdDataPath == nil {
		return fmt.Errorf("fail to create data path for pcc rule[%s]", pccRule.PccRuleId)
	}
	c.Log.Tracef("CreatePccRuleDataPath: pcc rule: %+v", pccRule)

	// Try to use a default pcc rule as default data path
	if c.Tunnel.DataPathPool.GetDefaultPath() == nil &&
		pccRule.Precedence == 255 {
		createdDataPath.IsDefaultPath = true
	}

	createdDataPath.GBRFlow = isGBRFlow(qosData)
	createdDataPath.ActivateTunnelAndPDR(c, uint32(pccRule.Precedence))
	c.Tunnel.AddDataPath(createdDataPath)
	pccRule.Datapath = createdDataPath
	pccRule.AddDataPathForwardingParameters(c, &targetRoute)

	if chgLevel, err := pccRule.IdentifyChargingLevel(); err != nil {
		c.Log.Warnf("fail to identify charging level[%+v] for pcc rule[%s]", err, pccRule.PccRuleId)
	} else {
		pccRule.Datapath.AddChargingRules(c, chgLevel, chgData)
	}

	if pccRule.RefQosDataID() != "" {
		pccRule.Datapath.AddQoS(c, pccRule.QFI, qosData)
		c.AddQosFlow(pccRule.QFI, qosData)
	}
	return nil
}

func (c *SMContext) CreateDcPccRuleDataPathOnDcTunnel(pccRule *PCCRule,
	tcData *TrafficControlData, qosData *models.QosData,
	chgData *models.ChargingData,
) error {
	var targetRoute models.RouteToLocation
	if tcData != nil && len(tcData.RouteToLocs) > 0 {
		targetRoute = *tcData.RouteToLocs[0]
	}
	param := &UPFSelectionParams{
		Dnn: c.Dnn,
		SNssai: &SNssai{
			Sst: c.SNssai.Sst,
			Sd:  c.SNssai.Sd,
		},
		Dnai: targetRoute.Dnai,
	}
	createdUpPath := GetUserPlaneInformation().GetDefaultUserPlanePathByDNN(param)
	createdDataPath := GenerateDataPath(createdUpPath)
	if createdDataPath == nil {
		return fmt.Errorf("fail to create data path on DCTunnel for pcc rule[%s]", pccRule.PccRuleId)
	}
	c.Log.Infof("CreatePccRuleDataPathOnDctunnel: pcc rule: %+v", pccRule)

	// Try to use a default pcc rule as default data path
	if c.DCTunnel.DataPathPool.GetDefaultPath() == nil &&
		pccRule.Precedence == 255 {
		createdDataPath.IsDefaultPath = true
	}

	createdDataPath.GBRFlow = isGBRFlow(qosData)
	createdDataPath.ActivateDcTunnelAndPDR(c, uint32(pccRule.Precedence))
	c.DCTunnel.AddDataPath(createdDataPath)
	pccRule.Datapath = createdDataPath
	pccRule.AddDataPathForwardingParametersOnDcTunnel(c, &targetRoute)

	if chgLevel, err := pccRule.IdentifyChargingLevel(); err != nil {
		c.Log.Warnf("fail to identify charging level[%+v] for pcc rule[%s]", err, pccRule.PccRuleId)
	} else {
		pccRule.Datapath.AddChargingRules(c, chgLevel, chgData)
	}

	if pccRule.RefQosDataID() != "" {
		pccRule.Datapath.AddQoS(c, pccRule.QFI, qosData)
		c.AddQosFlow(pccRule.QFI, qosData)
	}

	return nil
}

func (c *SMContext) BuildUpPathChgEventExposureNotification(
	chgEvent *models.UpPathChgEvent,
	srcRoute, tgtRoute *models.RouteToLocation,
) {
	if chgEvent == nil {
		return
	}
	if chgEvent.NotificationUri == "" {
		c.Log.Warnf("No NotificationUri [%s]", chgEvent.NotificationUri)
		return
	}

	en := models.SmfEventExposureEventNotification{
		Event:            models.SmfEvent_UP_PATH_CH,
		SourceTraRouting: srcRoute,
		TargetTraRouting: tgtRoute,
	}
	if srcRoute.Dnai != tgtRoute.Dnai {
		en.SourceDnai = srcRoute.Dnai
		en.TargetDnai = tgtRoute.Dnai
	}
	// TODO: sourceUeIpv4Addr, sourceUeIpv6Prefix, targetUeIpv4Addr, targetUeIpv6Prefix

	k := chgEvent.NotificationUri + chgEvent.NotifCorreId
	if strings.Contains(string(chgEvent.DnaiChgType), "EARLY") {
		en.DnaiChgType = models.DnaiChangeType("EARLY")
		v, ok := c.UpPathChgEarlyNotification[k]
		if ok {
			v.EventNotifs = append(v.EventNotifs, en)
		} else {
			c.UpPathChgEarlyNotification[k] = newEventExposureNotification(
				chgEvent.NotificationUri, chgEvent.NotifCorreId, &en)
		}
	}
	if strings.Contains(string(chgEvent.DnaiChgType), "LATE") {
		en.DnaiChgType = models.DnaiChangeType("LATE")
		v, ok := c.UpPathChgLateNotification[k]
		if ok {
			v.EventNotifs = append(v.EventNotifs, en)
		} else {
			c.UpPathChgLateNotification[k] = newEventExposureNotification(
				chgEvent.NotificationUri, chgEvent.NotifCorreId, &en)
		}
	}
}

func newEventExposureNotification(
	uri, id string,
	en *models.SmfEventExposureEventNotification,
) *EventExposureNotification {
	return &EventExposureNotification{
		NsmfEventExposureNotification: &models.NsmfEventExposureNotification{
			NotifId:     id,
			EventNotifs: []models.SmfEventExposureEventNotification{*en},
		},
		Uri: uri,
	}
}

type NotifCallback func(uri string,
	notification *models.NsmfEventExposureNotification)

func (c *SMContext) SendUpPathChgNotification(chgType string, notifCb NotifCallback) {
	var notifications map[string]*EventExposureNotification
	switch chgType {
	case "EARLY":
		notifications = c.UpPathChgEarlyNotification
	case "LATE":
		notifications = c.UpPathChgLateNotification
	default:
		return
	}
	for k, n := range notifications {
		c.Log.Infof("Send UpPathChg Event Exposure Notification [%s][%s] to NEF/AF", chgType, n.NotifId)
		go notifCb(n.Uri, n.NsmfEventExposureNotification)
		delete(notifications, k)
	}
}

func (smContext *SMContext) RemovePDRfromPFCPSession(nodeID pfcpType.NodeID, pdr *PDR) {
	NodeIDtoIP := nodeID.ResolveNodeIdToIp().String()
	pfcpSessCtx := smContext.PFCPContext[NodeIDtoIP]
	delete(pfcpSessCtx.PDRs, pdr.PDRID)
}

func (smContext *SMContext) IsAllowedPDUSessionType(requestedPDUSessionType uint8) error {
	dnnPDUSessionType := smContext.DnnConfiguration.PduSessionTypes
	if dnnPDUSessionType == nil {
		return fmt.Errorf("this SMContext[%s] has no subscription pdu session type info", smContext.Ref)
	}

	allowIPv4 := false
	allowIPv6 := false
	allowEthernet := false

	for _, allowedPDUSessionType := range smContext.DnnConfiguration.PduSessionTypes.AllowedSessionTypes {
		switch allowedPDUSessionType {
		case models.PduSessionType_IPV4:
			allowIPv4 = true
		case models.PduSessionType_IPV6:
			allowIPv6 = true
		case models.PduSessionType_IPV4_V6:
			allowIPv4 = true
			allowIPv6 = true
		case models.PduSessionType_ETHERNET:
			allowEthernet = true
		}
	}

	supportedPDUSessionType := GetSelf().SupportedPDUSessionType
	switch supportedPDUSessionType {
	case "IPv4":
		if !allowIPv4 {
			return fmt.Errorf(
				"no SupportedPDUSessionType[%q] in DNN[%s] configuration",
				supportedPDUSessionType,
				smContext.Dnn,
			)
		}
	case "IPv6":
		if !allowIPv6 {
			return fmt.Errorf(
				"no SupportedPDUSessionType[%q] in DNN[%s] configuration",
				supportedPDUSessionType,
				smContext.Dnn,
			)
		}
	case "IPv4v6":
		if !allowIPv4 && !allowIPv6 {
			return fmt.Errorf(
				"no SupportedPDUSessionType[%q] in DNN[%s] configuration",
				supportedPDUSessionType,
				smContext.Dnn,
			)
		}
	case "Ethernet":
		if !allowEthernet {
			return fmt.Errorf(
				"no SupportedPDUSessionType[%q] in DNN[%s] configuration",
				supportedPDUSessionType,
				smContext.Dnn,
			)
		}
	}

	smContext.EstAcceptCause5gSMValue = 0
	switch nasConvert.PDUSessionTypeToModels(requestedPDUSessionType) {
	case models.PduSessionType_IPV4:
		if allowIPv4 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV4)
		} else {
			return fmt.Errorf("PduSessionType_IPV4 is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	case models.PduSessionType_IPV6:
		if allowIPv6 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV6)
		} else {
			return fmt.Errorf("PduSessionType_IPV6 is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	case models.PduSessionType_IPV4_V6:
		if allowIPv4 && allowIPv6 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV4_V6)
		} else if allowIPv4 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV4)
			smContext.EstAcceptCause5gSMValue = nasMessage.Cause5GSMPDUSessionTypeIPv4OnlyAllowed
		} else if allowIPv6 {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_IPV6)
			smContext.EstAcceptCause5gSMValue = nasMessage.Cause5GSMPDUSessionTypeIPv6OnlyAllowed
		} else {
			return fmt.Errorf("PduSessionType_IPV4_V6 is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	case models.PduSessionType_ETHERNET:
		if allowEthernet {
			smContext.SelectedPDUSessionType = nasConvert.ModelsToPDUSessionType(models.PduSessionType_ETHERNET)
		} else {
			return fmt.Errorf("PduSessionType_ETHERNET is not allowed in DNN[%s] configuration", smContext.Dnn)
		}
	default:
		return fmt.Errorf("requested PDU Sesstion type[%d] is not supported", requestedPDUSessionType)
	}
	return nil
}

func (smContext *SMContext) StopT3591() {
	if smContext.T3591 != nil {
		smContext.T3591.Stop()
		smContext.T3591 = nil
	}
}

func (smContext *SMContext) StopT3592() {
	if smContext.T3592 != nil {
		smContext.T3592.Stop()
		smContext.T3592 = nil
	}
}

func (smContextState SMContextState) String() string {
	switch smContextState {
	case InActive:
		return "InActive"
	case ActivePending:
		return "ActivePending"
	case Active:
		return "Active"
	case InActivePending:
		return "InActivePending"
	case ModificationPending:
		return "ModificationPending"
	case PFCPModification:
		return "PFCPModification"
	default:
		return "Unknown State"
	}
}

func (smContext *SMContext) AssignQFI(qosId string) uint8 {
	qfi, ok := smContext.qosDataToQFI[qosId]
	if !ok {
		newId, err := smContext.QFIGenerator.Allocate()
		if err != nil {
			return 0
		}
		smContext.qosDataToQFI[qosId] = uint8(newId)
		return uint8(newId)
	}
	return qfi
}

func (smContext *SMContext) RemoveQFI(qosId string) {
	qfi, ok := smContext.qosDataToQFI[qosId]
	if ok {
		smContext.QFIGenerator.FreeID(int64(qfi))
		delete(smContext.qosDataToQFI, qosId)
		smContext.RemoveQosFlow(qfi)
	}
}
