package smf_context

import (
	"fmt"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"net"
	"reflect"

	"github.com/google/uuid"
)

var upfPool map[string]*UPF

func init() {
	upfPool = make(map[string]*UPF)
}

type UPTunnel struct {
	Node  *UPF
	ULPDR *PDR
	DLPDR *PDR

	ULTEID uint32
	DLTEID uint32

	UpfRoot  *DataPathNode
	ULCLRoot *DataPathNode
}

type UPFStatus int

const (
	NotAssociated          UPFStatus = 0
	AssociatedSettingUp    UPFStatus = 1
	AssociatedSetUpSuccess UPFStatus = 2
)

type UPF struct {
	uuid      uuid.UUID
	NodeID    pfcpType.NodeID
	UPIPInfo  pfcpType.UserPlaneIPResourceInformation
	UPFStatus UPFStatus

	pdrPool         map[uint16]*PDR
	farPool         map[uint32]*FAR
	barPool         map[uint8]*BAR
	urrPool         map[uint32]*URR
	qerPool         map[uint32]*QER
	teidPool        map[uint32]bool
	pdrCount        uint16
	farCount        uint32
	barCount        uint8
	urrCount        uint32
	qerCount        uint32
	TEIDCount       uint32
	pdrIDReuseQueue *IDQueue
	farIDReuseQueue *IDQueue
	barIDReuseQueue *IDQueue
}

// UUID return this UPF UUID (allocate by SMF in this time)
// Maybe allocate by UPF in future
func (upf *UPF) UUID() string {
	return upf.uuid.String()
}

// NewUPF returns a new UPF context in SMF
func NewUPF(nodeID *pfcpType.NodeID) (upf *UPF) {
	upf = new(UPF)
	upf.uuid = uuid.New()

	upfPool[upf.UUID()] = upf

	// Initialize context
	upf.UPFStatus = NotAssociated
	upf.NodeID = *nodeID
	upf.pdrPool = make(map[uint16]*PDR)
	upf.farPool = make(map[uint32]*FAR)
	upf.barPool = make(map[uint8]*BAR)
	upf.qerPool = make(map[uint32]*QER)
	upf.urrPool = make(map[uint32]*URR)
	upf.teidPool = make(map[uint32]bool)
	upf.pdrIDReuseQueue = NewIDQueue(PDRType)
	upf.farIDReuseQueue = NewIDQueue(FARType)
	upf.barIDReuseQueue = NewIDQueue(BARType)

	return upf
}

func (upf *UPF) GenerateTEID() (id uint32, err error) {
	if upf.UPFStatus != AssociatedSetUpSuccess {
		err := fmt.Errorf("this upf not associate with smf")
		return 0, err
	}
	id = uint32(upf.GetValidID(TEIDType))
	upf.teidPool[id] = true
	return
}

func (upf *UPF) PFCPAddr() *net.UDPAddr {
	return &net.UDPAddr{
		IP:   upf.NodeID.ResolveNodeIdToIp(),
		Port: pfcpUdp.PFCP_PORT,
	}
}

func RetrieveUPFNodeByNodeID(nodeID pfcpType.NodeID) (upf *UPF) {
	for _, upf := range upfPool {
		if reflect.DeepEqual(upf.NodeID, nodeID) {
			return upf
		}
	}
	return nil
}

func RemoveUPFNodeByNodeId(nodeID pfcpType.NodeID) {
	for upfID, upf := range upfPool {
		if reflect.DeepEqual(upf.NodeID, nodeID) {
			delete(upfPool, upfID)
			break
		}
	}
}

func SelectUPFByDnn(Dnn string) *UPF {
	for _, upf := range upfPool {
		if upf.UPIPInfo.Assoni && string(upf.UPIPInfo.NetworkInstance) == Dnn {
			return upf
		}
	}
	return nil
}

func (upf *UPF) GetUPFIP() string {
	return upf.NodeID.ResolveNodeIdToIp().String()
}

func (upf *UPF) GetUPFID() string {

	upInfo := GetUserPlaneInformation()
	upfIP := upf.NodeID.ResolveNodeIdToIp().String()
	return upInfo.GetUPFIDByIP(upfIP)

}

func (upf *UPF) pdrID() (pdrID uint16, err error) {
	if upf.UPFStatus != AssociatedSetUpSuccess {
		err := fmt.Errorf("this upf not associate with smf")
		return 0, err
	}

	if upf.pdrIDReuseQueue.IsEmpty() {
		id := upf.GetValidID(PDRType)
		pdrID = uint16(id)
	} else {
		id, err := upf.pdrIDReuseQueue.Pop()

		if err != nil {
			logger.CtxLog.Errorln("allocate id error:", err)
		}

		pdrID = uint16(id)
	}

	return
}

func (upf *UPF) farID() (farID uint32, err error) {
	if upf.UPFStatus != AssociatedSetUpSuccess {
		err := fmt.Errorf("this upf not associate with smf")
		return 0, err
	}

	if upf.farIDReuseQueue.IsEmpty() {

		id := upf.GetValidID(FARType)
		farID = uint32(id)
	} else {
		id, err := upf.farIDReuseQueue.Pop()

		if err != nil {
			logger.CtxLog.Errorln("allocate id error:", err)
		}
		farID = uint32(id)
	}

	return
}

func (upf *UPF) barID() (barID uint8, err error) {
	if upf.UPFStatus != AssociatedSetUpSuccess {
		err := fmt.Errorf("this upf not associate with smf")
		return 0, err
	}

	if upf.barIDReuseQueue.IsEmpty() {

		id := upf.GetValidID(BARType)
		barID = uint8(id)
	} else {
		id, err := upf.barIDReuseQueue.Pop()

		if err != nil {
			logger.CtxLog.Errorln("allocate id error:", err)
		}
		barID = uint8(id)
	}

	return
}

func (upf *UPF) AddPDR() (pdr *PDR, err error) {

	if upf.UPFStatus != AssociatedSetUpSuccess {
		err = fmt.Errorf("this upf not associate with smf")
		return nil, err
	}

	pdr = new(PDR)
	PDRID, _ := upf.pdrID()
	pdr.PDRID = uint16(PDRID)
	upf.pdrPool[pdr.PDRID] = pdr
	pdr.FAR, _ = upf.AddFAR()
	return pdr, nil
}

func (upf *UPF) AddFAR() (far *FAR, err error) {

	if upf.UPFStatus != AssociatedSetUpSuccess {
		err = fmt.Errorf("this upf not associate with smf")
		return nil, err
	}

	far = new(FAR)
	far.FARID, _ = upf.farID()
	upf.farPool[far.FARID] = far
	return far, nil
}

func (upf *UPF) AddBAR() (bar *BAR, err error) {

	if upf.UPFStatus != AssociatedSetUpSuccess {
		err = fmt.Errorf("this upf not associate with smf")
		return nil, err
	}

	bar = new(BAR)
	BARID, _ := upf.barID()
	bar.BARID = uint8(BARID)
	upf.barPool[bar.BARID] = bar
	return bar, nil
}

func (pdr *PDR) InitializePDR(smContext *SMContext) {

	tunnel := smContext.Tunnel
	pdr.State = RULE_INITIAL
	pdr.Precedence = 32
	pdr.PDI = PDI{
		SourceInterface: pfcpType.SourceInterface{
			InterfaceValue: pfcpType.SourceInterfaceAccess,
		},
		LocalFTeid: &pfcpType.FTEID{
			V4:          true,
			Teid:        tunnel.ULTEID,
			Ipv4Address: tunnel.Node.UPIPInfo.Ipv4Address,
		},
		NetworkInstance: []byte(smContext.Dnn),
		UEIPAddress: &pfcpType.UEIPAddress{
			V4:          true,
			Ipv4Address: smContext.PDUAddress.To4(),
		},
	}
	pdr.OuterHeaderRemoval = new(pfcpType.OuterHeaderRemoval)
	pdr.OuterHeaderRemoval.OuterHeaderRemovalDescription = pfcpType.OuterHeaderRemovalGtpUUdpIpv4

	pdr.FAR.InitializeFAR(smContext)
}

func (far *FAR) InitializeFAR(smContext *SMContext) {

	far.ApplyAction.Forw = true
	far.ForwardingParameters = &ForwardingParameters{
		DestinationInterface: pfcpType.DestinationInterface{
			InterfaceValue: pfcpType.DestinationInterfaceCore,
		},
		NetworkInstance: []byte(smContext.Dnn),
	}
}

func (upf *UPF) RemovePDR(pdr *PDR) (err error) {

	if upf.UPFStatus != AssociatedSetUpSuccess {
		err = fmt.Errorf("this upf not associate with smf")
		return err
	}

	upf.pdrIDReuseQueue.Push(int(pdr.PDRID))
	delete(upf.pdrPool, pdr.PDRID)
	return nil
}

func (upf *UPF) RemoveFAR(far *FAR) (err error) {

	upf.farIDReuseQueue.Push(int(far.FARID))
	delete(upf.farPool, far.FARID)
	return nil
}

func (upf *UPF) RemoveBAR(bar *BAR) (err error) {

	if upf.UPFStatus != AssociatedSetUpSuccess {
		err = fmt.Errorf("this upf not associate with smf")
		return err
	}

	upf.barIDReuseQueue.Push(int(bar.BARID))
	delete(upf.barPool, bar.BARID)
	return nil
}

func (upf *UPF) GetValidID(idType IDType) (id int) {

	switch idType {
	case PDRType:
		for {
			upf.pdrCount++
			if _, exist := upf.pdrPool[upf.pdrCount]; !exist { // valid id
				break
			}
		}

		id = int(upf.pdrCount)
	case FARType:
		for {
			upf.farCount++
			if _, exist := upf.farPool[upf.farCount]; !exist { // valid id
				break
			}
		}

		id = int(upf.farCount)
	case BARType:
		for {
			upf.barCount++
			if _, exist := upf.barPool[upf.barCount]; !exist { // valid id
				break
			}
		}

		id = int(upf.barCount)
	case TEIDType:
		for {
			upf.TEIDCount++
			if _, exist := upf.teidPool[upf.TEIDCount]; !exist { // valid id
				break
			}
		}

		id = int(upf.TEIDCount)
	}
	return
}

func (upf *UPF) CheckPDRIDExist(id int) (exist bool) {
	_, exist = upf.pdrPool[uint16(id)]
	return
}

func (upf *UPF) CheckFARIDExist(id int) (exist bool) {
	_, exist = upf.farPool[uint32(id)]
	return
}

func (upf *UPF) CheckBARIDExist(id int) (exist bool) {
	_, exist = upf.barPool[uint8(id)]
	return
}
