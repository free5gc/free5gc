package pfcp

import (
	"net"

	"github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"
)

func (s *PfcpServer) handleAssociationSetupRequest(
	req *message.AssociationSetupRequest,
	addr net.Addr,
) {
	s.log.Infoln("handleAssociationSetupRequest")

	if req.NodeID == nil {
		s.log.Errorln("not found NodeID")
		return
	}
	rnodeid, err := req.NodeID.NodeID()
	if err != nil {
		s.log.Errorln(err)
		return
	}
	s.log.Debugf("remote nodeid: %v\n", rnodeid)

	// deleting the existing PFCP association and associated PFCP sessions,
	// if a PFCP association was already established for the Node ID
	// received in the request, regardless of the Recovery Timestamp
	// received in the request.
	if node, ok := s.rnodes[rnodeid]; ok {
		s.log.Infof("delete node: %#+v\n", node)
		node.Reset()
		delete(s.rnodes, rnodeid)
	}
	node := s.NewNode(rnodeid, addr, s.driver)
	s.rnodes[rnodeid] = node

	rsp := message.NewAssociationSetupResponse(
		req.Header.SequenceNumber,
		newIeNodeID(s.nodeID),
		ie.NewCause(ie.CauseRequestAccepted),
		ie.NewRecoveryTimeStamp(s.recoveryTime),
		// TODO:
		// ie.NewUPFunctionFeatures(),
	)

	err = s.sendRspTo(rsp, addr)
	if err != nil {
		s.log.Errorln(err)
		return
	}
}

func (s *PfcpServer) handleAssociationUpdateRequest(
	req *message.AssociationUpdateRequest,
	addr net.Addr,
) {
	s.log.Infoln("handleAssociationUpdateRequest not supported")
}

func (s *PfcpServer) handleAssociationReleaseRequest(
	req *message.AssociationReleaseRequest,
	addr net.Addr,
) {
	s.log.Infoln("handleAssociationReleaseRequest not supported")
}

func newIeNodeID(nodeID string) *ie.IE {
	ip := net.ParseIP(nodeID)
	if ip != nil {
		if ip.To4() != nil {
			return ie.NewNodeID(nodeID, "", "")
		}
		return ie.NewNodeID("", nodeID, "")
	}
	return ie.NewNodeID("", "", nodeID)
}
