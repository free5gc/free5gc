package smf_context

import (
	"fmt"
	"net"
	"reflect"

	"free5gc/lib/pfcp/pfcpType"
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
}

type UPF struct {
	NodeID   pfcpType.NodeID
	UPIPInfo pfcpType.UserPlaneIPResourceInformation

	pdrPool   map[int]*PDR
	farPool   map[int]*FAR
	barPool   map[int]*BAR
	urrPool   map[int]*URR
	qerPool   map[int]*QER
	pdrCount  uint32
	farCount  uint32
	barCount  uint32
	urrCount  uint32
	qerCount  uint32
	TEIDCount uint32
}

func AddUPF(nodeId *pfcpType.NodeID) (upf *UPF) {
	upf = new(UPF)
	key, _ := generateUpfIdFromNodeId(nodeId)
	upfPool[key] = upf
	upf.NodeID = *nodeId
	upf.pdrPool = make(map[int]*PDR)
	upf.farPool = make(map[int]*FAR)
	upf.barPool = make(map[int]*BAR)
	upf.qerPool = make(map[int]*QER)
	upf.urrPool = make(map[int]*URR)
	return
}

func generateUpfIdFromNodeId(nodeId *pfcpType.NodeID) (string, error) {
	switch nodeId.NodeIdType {
	case pfcpType.NodeIdTypeIpv4Address, pfcpType.NodeIdTypeIpv6Address:
		return net.IP(nodeId.NodeIdValue).String(), nil
	case pfcpType.NodeIdTypeFqdn:
		return string(nodeId.NodeIdValue), nil
	default:
		return "", fmt.Errorf("Invalid Node ID type: %v", nodeId.NodeIdType)
	}
}

func (upf *UPF) GenerateTEID() uint32 {
	upf.TEIDCount++
	return upf.TEIDCount
}

func RetrieveUPFNodeByNodeId(nodeId pfcpType.NodeID) (upf *UPF) {
	for _, upf := range upfPool {
		if reflect.DeepEqual(upf.NodeID, nodeId) {
			return upf
		}
	}
	return nil
}

func RemoveUPFNodeByNodeId(nodeId pfcpType.NodeID) {
	for upfID, upf := range upfPool {
		if reflect.DeepEqual(upf.NodeID, nodeId) {
			delete(upfPool, upfID)
			break
		}
	}
}

func SelectUPFByDnn(Dnn string) *UPF {
	for _, upf := range upfPool {
		if !upf.UPIPInfo.Assoni || string(upf.UPIPInfo.NetworkInstance) == Dnn {
			return upf
		}
	}
	return nil
}

func (upf *UPF) pdrID() uint32 {
	upf.pdrCount++
	return upf.pdrCount
}

func (upf *UPF) farID() uint32 {
	upf.farCount++
	return upf.farCount
}

func (upf *UPF) barID() uint32 {
	upf.barCount++
	return upf.barCount
}

func (upf *UPF) AddPDR() (pdr *PDR) {
	pdr = new(PDR)
	pdr.PDRID = uint16(upf.pdrID())
	upf.pdrPool[len(upf.pdrPool)] = pdr
	pdr.FAR = upf.AddFAR()
	return pdr
}

func (upf *UPF) AddFAR() (far *FAR) {
	far = new(FAR)
	far.FARID = upf.farID()
	upf.farPool[len(upf.farPool)] = far
	return far
}

func (upf *UPF) AddBAR() (bar *BAR) {
	bar = new(BAR)
	bar.BARID = uint8(upf.barID())
	upf.barPool[len(upf.barPool)] = bar
	return bar
}
