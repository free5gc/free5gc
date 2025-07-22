package context

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/pfcp/pfcpUdp"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/pkg/factory"
	"github.com/free5gc/util/idgenerator"
)

var upfPool sync.Map

type UPTunnel struct {
	PathIDGenerator *idgenerator.IDGenerator
	DataPathPool    DataPathPool
	ANInformation   struct {
		IPAddress net.IP
		TEID      uint32
	}
}

func (t *UPTunnel) UpdateANInformation(ip net.IP, teid uint32) {
	t.ANInformation.IPAddress = ip
	t.ANInformation.TEID = teid

	for _, dataPath := range t.DataPathPool {
		if dataPath.Activated {
			ANUPF := dataPath.FirstDPNode
			DLPDR := ANUPF.DownLinkTunnel.PDR

			if DLPDR.FAR.ForwardingParameters.OuterHeaderCreation != nil {
				// Old AN tunnel exists
				DLPDR.FAR.ForwardingParameters.SendEndMarker = true
			}

			DLPDR.FAR.ForwardingParameters.OuterHeaderCreation = new(pfcpType.OuterHeaderCreation)
			dlOuterHeaderCreation := DLPDR.FAR.ForwardingParameters.OuterHeaderCreation
			dlOuterHeaderCreation.OuterHeaderCreationDescription = pfcpType.OuterHeaderCreationGtpUUdpIpv4
			dlOuterHeaderCreation.Teid = t.ANInformation.TEID
			dlOuterHeaderCreation.Ipv4Address = t.ANInformation.IPAddress.To4()
			DLPDR.FAR.State = RULE_UPDATE
		}
	}
}

type UPFStatus int

const (
	NotAssociated          UPFStatus = 0
	AssociatedSettingUp    UPFStatus = 1
	AssociatedSetUpSuccess UPFStatus = 2
)

type UPF struct {
	uuid              uuid.UUID
	NodeID            pfcpType.NodeID
	Addr              string
	RecoveryTimeStamp time.Time

	AssociationContext context.Context
	CancelAssociation  context.CancelFunc

	SNssaiInfos  []*SnssaiUPFInfo
	N3Interfaces []*UPFInterfaceInfo
	N9Interfaces []*UPFInterfaceInfo

	pdrPool sync.Map
	farPool sync.Map
	barPool sync.Map
	qerPool sync.Map
	urrPool sync.Map

	pdrIDGenerator *idgenerator.IDGenerator
	farIDGenerator *idgenerator.IDGenerator
	barIDGenerator *idgenerator.IDGenerator
	urrIDGenerator *idgenerator.IDGenerator
	qerIDGenerator *idgenerator.IDGenerator
}

// UPFSelectionParams ... parameters for upf selection
type UPFSelectionParams struct {
	Dnn        string
	SNssai     *SNssai
	Dnai       string
	PDUAddress net.IP
}

// UPFInterfaceInfo store the UPF interface information
type UPFInterfaceInfo struct {
	NetworkInstances      []string
	IPv4EndPointAddresses []net.IP
	IPv6EndPointAddresses []net.IP
	EndpointFQDN          string
}

func GetUpfById(uuid string) *UPF {
	upf, ok := upfPool.Load(uuid)
	if ok {
		return upf.(*UPF)
	}
	return nil
}

// NewUPFInterfaceInfo parse the InterfaceUpfInfoItem to generate UPFInterfaceInfo
func NewUPFInterfaceInfo(i *factory.InterfaceUpfInfoItem) *UPFInterfaceInfo {
	interfaceInfo := new(UPFInterfaceInfo)

	interfaceInfo.IPv4EndPointAddresses = make([]net.IP, 0)
	interfaceInfo.IPv6EndPointAddresses = make([]net.IP, 0)

	logger.CtxLog.Infoln("Endpoints:", i.Endpoints)

	for _, endpoint := range i.Endpoints {
		eIP := net.ParseIP(endpoint)
		if eIP == nil {
			interfaceInfo.EndpointFQDN = endpoint
		} else if eIPv4 := eIP.To4(); eIPv4 == nil {
			interfaceInfo.IPv6EndPointAddresses = append(interfaceInfo.IPv6EndPointAddresses, eIP)
		} else {
			interfaceInfo.IPv4EndPointAddresses = append(interfaceInfo.IPv4EndPointAddresses, eIPv4)
		}
	}

	interfaceInfo.NetworkInstances = make([]string, len(i.NetworkInstances))
	copy(interfaceInfo.NetworkInstances, i.NetworkInstances)

	return interfaceInfo
}

// *** add unit test ***//
// IP returns the IP of the user plane IP information of the pduSessType
func (i *UPFInterfaceInfo) IP(pduSessType uint8) (net.IP, error) {
	if (pduSessType == nasMessage.PDUSessionTypeIPv4 ||
		pduSessType == nasMessage.PDUSessionTypeIPv4IPv6) && len(i.IPv4EndPointAddresses) != 0 {
		return i.IPv4EndPointAddresses[0], nil
	}

	if (pduSessType == nasMessage.PDUSessionTypeIPv6 ||
		pduSessType == nasMessage.PDUSessionTypeIPv4IPv6) && len(i.IPv6EndPointAddresses) != 0 {
		return i.IPv6EndPointAddresses[0], nil
	}

	if i.EndpointFQDN != "" {
		if resolvedAddr, err := net.ResolveIPAddr("ip", i.EndpointFQDN); err != nil {
			logger.CtxLog.Errorf("resolve addr [%s] failed", i.EndpointFQDN)
		} else {
			switch pduSessType {
			case nasMessage.PDUSessionTypeIPv4:
				return resolvedAddr.IP.To4(), nil
			case nasMessage.PDUSessionTypeIPv6:
				return resolvedAddr.IP.To16(), nil
			default:
				v4addr := resolvedAddr.IP.To4()
				if v4addr != nil {
					return v4addr, nil
				} else {
					return resolvedAddr.IP.To16(), nil
				}
			}
		}
	}

	return nil, errors.New("not matched ip address")
}

func (upfSelectionParams *UPFSelectionParams) String() string {
	str := ""
	Dnn := upfSelectionParams.Dnn
	if Dnn != "" {
		str += fmt.Sprintf("Dnn: %s\n", Dnn)
	}

	SNssai := upfSelectionParams.SNssai
	if SNssai != nil {
		str += fmt.Sprintf("Sst: %d, Sd: %s\n", int(SNssai.Sst), SNssai.Sd)
	}

	Dnai := upfSelectionParams.Dnai
	if Dnai != "" {
		str += fmt.Sprintf("DNAI: %s\n", Dnai)
	}

	pduAddress := upfSelectionParams.PDUAddress
	if pduAddress != nil {
		str += fmt.Sprintf("PDUAddress: %s\n", pduAddress)
	}

	return str
}

// UUID return this UPF UUID (allocate by SMF in this time)
// Maybe allocate by UPF in future
func (upf *UPF) UUID() string {
	uuid := upf.uuid.String()
	return uuid
}

func NewUPTunnel() (tunnel *UPTunnel) {
	tunnel = &UPTunnel{
		DataPathPool:    make(DataPathPool),
		PathIDGenerator: idgenerator.NewGenerator(1, 2147483647),
	}

	return
}

// *** add unit test ***//
func (t *UPTunnel) AddDataPath(dataPath *DataPath) {
	pathID, err := t.PathIDGenerator.Allocate()
	if err != nil {
		logger.CtxLog.Warnf("Allocate pathID error: %+v", err)
		return
	}

	dataPath.PathID = pathID
	t.DataPathPool[pathID] = dataPath
}

func (t *UPTunnel) RemoveDataPath(pathID int64) {
	delete(t.DataPathPool, pathID)
	t.PathIDGenerator.FreeID(pathID)
}

// *** add unit test ***//
// NewUPF returns a new UPF context in SMF
func NewUPF(nodeID *pfcpType.NodeID, ifaces []*factory.InterfaceUpfInfoItem) (upf *UPF) {
	upf = new(UPF)
	upf.uuid = uuid.New()

	upfPool.Store(upf.UUID(), upf)

	// Initialize context
	upf.AssociationContext, upf.CancelAssociation = context.WithCancel(context.Background())
	upf.CancelAssociation() // necessary to avoid nil pointer for checks of AssociationContext before UPF is associated

	upf.NodeID = *nodeID
	upf.pdrIDGenerator = idgenerator.NewGenerator(1, math.MaxUint16)
	upf.farIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
	upf.barIDGenerator = idgenerator.NewGenerator(1, math.MaxUint8)
	upf.qerIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)
	upf.urrIDGenerator = idgenerator.NewGenerator(1, math.MaxUint32)

	upf.N3Interfaces = make([]*UPFInterfaceInfo, 0)
	upf.N9Interfaces = make([]*UPFInterfaceInfo, 0)

	for _, iface := range ifaces {
		upIface := NewUPFInterfaceInfo(iface)

		switch iface.InterfaceType {
		case models.UpInterfaceType_N3:
			upf.N3Interfaces = append(upf.N3Interfaces, upIface)
		case models.UpInterfaceType_N9:
			upf.N9Interfaces = append(upf.N9Interfaces, upIface)
		}
	}

	return upf
}

// *** add unit test ***//
// GetInterface return the UPFInterfaceInfo that match input cond
func (upf *UPF) GetInterface(interfaceType models.UpInterfaceType, dnn string) *UPFInterfaceInfo {
	switch interfaceType {
	case models.UpInterfaceType_N3:
		for i, iface := range upf.N3Interfaces {
			for _, nwInst := range iface.NetworkInstances {
				if nwInst == dnn {
					return upf.N3Interfaces[i]
				}
			}
		}
	case models.UpInterfaceType_N9:
		for i, iface := range upf.N9Interfaces {
			for _, nwInst := range iface.NetworkInstances {
				if nwInst == dnn {
					return upf.N9Interfaces[i]
				}
			}
		}
	}
	return nil
}

func (upf *UPF) PFCPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   upf.NodeID.ResolveNodeIdToIp(),
		Port: pfcpUdp.PFCP_PORT,
	}
}

// *** add unit test ***//
func RetrieveUPFNodeByNodeID(nodeID pfcpType.NodeID) *UPF {
	var targetUPF *UPF = nil
	upfPool.Range(func(key, value interface{}) bool {
		curUPF := value.(*UPF)
		if curUPF.NodeID.NodeIdType != nodeID.NodeIdType &&
			(curUPF.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn || nodeID.NodeIdType == pfcpType.NodeIdTypeFqdn) {
			curUPFNodeIdIP := curUPF.NodeID.ResolveNodeIdToIp().To4()
			nodeIdIP := nodeID.ResolveNodeIdToIp().To4()
			logger.CtxLog.Tracef("RetrieveUPF - upfNodeIdIP:[%+v], nodeIdIP:[%+v]", curUPFNodeIdIP, nodeIdIP)
			if reflect.DeepEqual(curUPFNodeIdIP, nodeIdIP) {
				targetUPF = curUPF
				return false
			}
		} else if reflect.DeepEqual(curUPF.NodeID, nodeID) {
			targetUPF = curUPF
			return false
		}
		return true
	})

	return targetUPF
}

// *** add unit test ***//
func RemoveUPFNodeByNodeID(nodeID pfcpType.NodeID) bool {
	upfID := ""
	upfPool.Range(func(key, value interface{}) bool {
		upfID = key.(string)
		upf := value.(*UPF)
		if upf.NodeID.NodeIdType != nodeID.NodeIdType &&
			(upf.NodeID.NodeIdType == pfcpType.NodeIdTypeFqdn || nodeID.NodeIdType == pfcpType.NodeIdTypeFqdn) {
			upfNodeIdIP := upf.NodeID.ResolveNodeIdToIp().To4()
			nodeIdIP := nodeID.ResolveNodeIdToIp().To4()
			logger.CtxLog.Tracef("RemoveUPF - upfNodeIdIP:[%+v], nodeIdIP:[%+v]", upfNodeIdIP, nodeIdIP)
			if reflect.DeepEqual(upfNodeIdIP, nodeIdIP) {
				return false
			}
		} else if reflect.DeepEqual(upf.NodeID, nodeID) {
			return false
		}
		upfID = ""
		return true
	})

	if upfID != "" {
		upfPool.Delete(upfID)
		return true
	}
	return false
}

func (upf *UPF) GetUPFIP() string {
	upfIP := upf.NodeID.ResolveNodeIdToIp().String()
	return upfIP
}

func (upf *UPF) GetUPFID() string {
	upInfo := GetUserPlaneInformation()
	upfIP := upf.NodeID.ResolveNodeIdToIp().String()
	return upInfo.GetUPFIDByIP(upfIP)
}

func (upf *UPF) pdrID() (pdrID uint16, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	tmpID, err := upf.pdrIDGenerator.Allocate()
	if err != nil {
		return 0, err
	}
	pdrID = uint16(tmpID)
	return
}

func (upf *UPF) farID() (farID uint32, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	tmpID, err := upf.farIDGenerator.Allocate()
	if err != nil {
		return 0, err
	}
	farID = uint32(tmpID)
	return
}

func (upf *UPF) barID() (barID uint8, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	tmpID, err := upf.barIDGenerator.Allocate()
	if err != nil {
		return 0, err
	}
	barID = uint8(tmpID)
	return
}

func (upf *UPF) qerID() (qerID uint32, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	tmpID, err := upf.qerIDGenerator.Allocate()
	if err != nil {
		return 0, err
	}
	qerID = uint32(tmpID)
	return
}

func (upf *UPF) urrID() (urrID uint32, err error) {
	tmpID, err := upf.urrIDGenerator.Allocate()
	if err != nil {
		return 0, err
	}
	urrID = uint32(tmpID)
	return
}

func (upf *UPF) AddPDR() (pdr *PDR, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	pdrID, err := upf.pdrID()
	if err != nil {
		return
	}

	newFAR, err := upf.AddFAR()
	if err != nil {
		return
	}

	pdr = &PDR{
		PDRID: pdrID,
		FAR:   newFAR,
	}
	upf.pdrPool.Store(pdr.PDRID, pdr)
	return
}

func (upf *UPF) AddFAR() (far *FAR, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	farID, err := upf.farID()
	if err != nil {
		return
	}
	far = &FAR{
		FARID: farID,
	}
	upf.farPool.Store(far.FARID, far)
	return far, nil
}

func (upf *UPF) AddBAR() (bar *BAR, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	barID, err := upf.barID()
	if err != nil {
		return
	}
	bar = &BAR{
		BARID: barID,
	}
	upf.barPool.Store(bar.BARID, bar)
	return
}

func (upf *UPF) AddQER() (qer *QER, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	qerID, err := upf.qerID()
	if err != nil {
		return
	}
	qer = &QER{
		QERID: qerID,
	}
	upf.qerPool.Store(qer.QERID, qer)
	return
}

func (upf *UPF) AddURR(urrID uint32, opts ...UrrOpt) (urr *URR, err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	if urrID == 0 {
		urrID, err = upf.urrID()
		if err != nil {
			return
		}
	}

	urr = &URR{
		URRID:                  urrID,
		MeasureMethod:          MesureMethodVol,
		MeasurementInformation: MeasureInformation(true, false),
	}

	for _, opt := range opts {
		opt(urr)
	}

	upf.urrPool.Store(urr.URRID, urr)
	return
}

func (upf *UPF) GetUUID() uuid.UUID {
	return upf.uuid
}

func (upf *UPF) GetQERById(qerId uint32) *QER {
	qer, ok := upf.qerPool.Load(qerId)
	if ok {
		return qer.(*QER)
	}
	return nil
}

// *** add unit test ***//
func (upf *UPF) RemovePDR(pdr *PDR) (err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	upf.pdrIDGenerator.FreeID(int64(pdr.PDRID))
	upf.pdrPool.Delete(pdr.PDRID)
	return
}

// *** add unit test ***//
func (upf *UPF) RemoveFAR(far *FAR) (err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	upf.farIDGenerator.FreeID(int64(far.FARID))
	upf.farPool.Delete(far.FARID)
	return
}

// *** add unit test ***//
func (upf *UPF) RemoveBAR(bar *BAR) (err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	upf.barIDGenerator.FreeID(int64(bar.BARID))
	upf.barPool.Delete(bar.BARID)
	return
}

// *** add unit test ***//
func (upf *UPF) RemoveQER(qer *QER) (err error) {
	if err = upf.IsAssociated(); err != nil {
		return
	}

	upf.qerIDGenerator.FreeID(int64(qer.QERID))
	upf.qerPool.Delete(qer.QERID)
	return
}

func (upf *UPF) isSupportSnssai(snssai *SNssai) bool {
	for _, snssaiInfo := range upf.SNssaiInfos {
		if snssaiInfo.SNssai.Equal(snssai) {
			return true
		}
	}
	return false
}

func (upf *UPF) ProcEachSMContext(procFunc func(*SMContext)) {
	smContextPool.Range(func(key, value interface{}) bool {
		smContext := value.(*SMContext)
		if smContext.SelectedUPF != nil && smContext.SelectedUPF.UPF == upf {
			procFunc(smContext)
		}
		return true
	})
}

func (upf *UPF) IsAssociated() error {
	select {
	case <-upf.AssociationContext.Done():
		return fmt.Errorf("UPF[%s] not associated with SMF",
			upf.NodeID.ResolveNodeIdToIp().String())
	default:
		return nil
	}
}
