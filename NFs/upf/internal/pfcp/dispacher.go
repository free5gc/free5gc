package pfcp

import (
	"net"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-pfcp/message"
)

func (s *PfcpServer) reqDispacher(msg message.Message, addr net.Addr) error {
	switch req := msg.(type) {
	case *message.HeartbeatRequest:
		s.handleHeartbeatRequest(req, addr)
	case *message.AssociationSetupRequest:
		s.handleAssociationSetupRequest(req, addr)
	case *message.AssociationUpdateRequest:
		s.handleAssociationUpdateRequest(req, addr)
	case *message.AssociationReleaseRequest:
		s.handleAssociationReleaseRequest(req, addr)
	case *message.SessionEstablishmentRequest:
		s.handleSessionEstablishmentRequest(req, addr)
	case *message.SessionModificationRequest:
		s.handleSessionModificationRequest(req, addr)
	case *message.SessionDeletionRequest:
		s.handleSessionDeletionRequest(req, addr)
	default:
		return errors.Errorf("pfcp reqDispacher unknown msg type: %d", msg.MessageType())
	}
	return nil
}

func (s *PfcpServer) rspDispacher(msg message.Message, addr net.Addr, req message.Message) error {
	switch rsp := msg.(type) {
	case *message.SessionReportResponse:
		s.handleSessionReportResponse(rsp, addr, req)
	default:
		return errors.Errorf("pfcp rspDispacher unknown msg type: %d", msg.MessageType())
	}
	return nil
}

func (s *PfcpServer) txtoDispacher(msg message.Message, addr net.Addr) error {
	switch req := msg.(type) {
	case *message.SessionReportRequest:
		s.handleSessionReportRequestTimeout(req, addr)
	default:
		return errors.Errorf("pfcp txtoDispacher unknown msg type: %d", msg.MessageType())
	}
	return nil
}
