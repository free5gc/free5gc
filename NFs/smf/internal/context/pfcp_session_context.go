package context

import (
	"fmt"
	"net"

	"github.com/free5gc/pfcp/pfcpType"
)

type PFCPSessionResponseStatus int

const (
	SessionEstablishSuccess PFCPSessionResponseStatus = iota
	SessionEstablishFailed
	SessionUpdateSuccess
	SessionUpdateFailed
	SessionReleaseSuccess
	SessionReleaseFailed
)

type FSEID struct {
	IP   net.IP
	SEID uint64
}

type PFCPSessionContext struct {
	PDRs       map[uint16]*PDR
	NodeID     pfcpType.NodeID
	LocalSEID  uint64
	RemoteSEID uint64
}

func (pfcpSessionContext *PFCPSessionContext) String() string {
	str := "\n"
	for pdrID, pdr := range pfcpSessionContext.PDRs {
		str += fmt.Sprintln("PDR ID: ", pdrID)
		str += fmt.Sprintf("PDR: %v\n", pdr)
	}

	str += fmt.Sprintln("Node ID: ", pfcpSessionContext.NodeID.ResolveNodeIdToIp().String())
	str += fmt.Sprintln("LocalSEID: ", pfcpSessionContext.LocalSEID)
	str += fmt.Sprintln("RemoteSEID: ", pfcpSessionContext.RemoteSEID)
	str += "\n"

	return str
}

func (pfcpSessionResponseStatus PFCPSessionResponseStatus) String() string {
	switch pfcpSessionResponseStatus {
	case SessionUpdateSuccess:
		return "SessionUpdateSuccess"
	case SessionUpdateFailed:
		return "SessionUpdateFailed"
	case SessionReleaseSuccess:
		return "SessionReleaseSuccess"
	case SessionReleaseFailed:
		return "SessionReleaseFailed"
	default:
		return "Unknown PFCP Session Response Status"
	}
}
