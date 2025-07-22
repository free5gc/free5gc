package message

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpType"
	"github.com/free5gc/pfcp/pfcpUdp"
	"github.com/free5gc/smf/internal/context"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/internal/pfcp/udp"
)

var seq uint32

func getSeqNumber() uint32 {
	return atomic.AddUint32(&seq, 1)
}

func SendPfcpAssociationSetupRequest(upNodeID pfcpType.NodeID) (resMsg *pfcpUdp.Message, err error) {
	pfcpMsg, err := BuildPfcpAssociationSetupRequest()
	if err != nil {
		return nil, fmt.Errorf("build PFCP Association Setup Request failed: %v", err)
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_NOT_PRESENT,
			MessageType:    pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
			SequenceNumber: getSeqNumber(),
		},
		Body: pfcpMsg,
	}

	addr := &net.UDPAddr{
		IP:   upNodeID.ResolveNodeIdToIp(),
		Port: pfcpUdp.PFCP_PORT,
	}

	resMsg, err = udp.SendPfcpRequest(message, addr)
	if err != nil {
		return nil, err
	}

	if resMsg.MessageType() != pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE {
		return resMsg, fmt.Errorf("received unexpected response message")
	}

	return resMsg, nil
}

func SendPfcpAssociationSetupResponse(addr *net.UDPAddr, cause pfcpType.Cause) {
	pfcpMsg, err := BuildPfcpAssociationSetupResponse(cause)
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Association Setup Response failed: %v", err)
		return
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_NOT_PRESENT,
			MessageType:    pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE,
			SequenceNumber: 1,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}

func SendPfcpAssociationReleaseRequest(upNodeID pfcpType.NodeID) (resMsg *pfcpUdp.Message, err error) {
	pfcpMsg, err := BuildPfcpAssociationReleaseRequest()
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Association Release Request failed: %v", err)
		return nil, err
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_NOT_PRESENT,
			MessageType:    pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST,
			SequenceNumber: 1,
		},
		Body: pfcpMsg,
	}

	addr := &net.UDPAddr{
		IP:   upNodeID.ResolveNodeIdToIp(),
		Port: pfcpUdp.PFCP_PORT,
	}

	resMsg, err = udp.SendPfcpRequest(message, addr)
	if err != nil {
		return nil, err
	}

	if resMsg.MessageType() != pfcp.PFCP_ASSOCIATION_RELEASE_RESPONSE {
		return resMsg, fmt.Errorf("received unexpected response message")
	}

	return resMsg, nil
}

func SendPfcpAssociationReleaseResponse(addr *net.UDPAddr, cause pfcpType.Cause) {
	pfcpMsg, err := BuildPfcpAssociationReleaseResponse(cause)
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Association Release Response failed: %v", err)
		return
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_NOT_PRESENT,
			MessageType:    pfcp.PFCP_ASSOCIATION_RELEASE_RESPONSE,
			SequenceNumber: 1,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}

func SendPfcpSessionEstablishmentRequest(
	upf *context.UPF,
	ctx *context.SMContext,
	pdrList []*context.PDR,
	farList []*context.FAR,
	barList []*context.BAR,
	qerList []*context.QER,
	urrList []*context.URR,
) (resMsg *pfcpUdp.Message, err error) {
	nodeIDtoIP := upf.NodeID.ResolveNodeIdToIp()
	if err = upf.IsAssociated(); err != nil {
		return nil, err
	}

	pfcpMsg, err := BuildPfcpSessionEstablishmentRequest(upf.NodeID, nodeIDtoIP.String(),
		ctx, pdrList, farList, barList, qerList, urrList)
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Establishment Request failed: %v", err)
		return nil, err
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:         pfcp.PfcpVersion,
			MP:              1,
			S:               pfcp.SEID_PRESENT,
			MessageType:     pfcp.PFCP_SESSION_ESTABLISHMENT_REQUEST,
			SEID:            0,
			SequenceNumber:  getSeqNumber(),
			MessagePriority: 0,
		},
		Body: pfcpMsg,
	}

	upaddr := &net.UDPAddr{
		IP:   nodeIDtoIP,
		Port: pfcpUdp.PFCP_PORT,
	}
	logger.PduSessLog.Traceln("[SMF] Send SendPfcpSessionEstablishmentRequest")
	logger.PduSessLog.Traceln("Send to addr ", upaddr.String())

	resMsg, err = udp.SendPfcpRequest(message, upaddr)
	if err != nil {
		return nil, err
	}

	if resMsg.MessageType() != pfcp.PFCP_SESSION_ESTABLISHMENT_RESPONSE {
		return resMsg, fmt.Errorf("received unexpected type response message: %+v", resMsg.PfcpMessage.Header)
	}

	localSEID := ctx.PFCPContext[nodeIDtoIP.String()].LocalSEID
	if resMsg.PfcpMessage.Header.SEID != localSEID {
		return resMsg, fmt.Errorf("received unexpected SEID response message: %+v, exptcted: %d",
			resMsg.PfcpMessage.Header, localSEID)
	}

	return resMsg, nil
}

// Deprecated: PFCP Session Establishment Procedure should be initiated by the CP function
func SendPfcpSessionEstablishmentResponse(addr *net.UDPAddr) {
	pfcpMsg, err := BuildPfcpSessionEstablishmentResponse()
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Establishment Response failed: %v", err)
		return
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:         pfcp.PfcpVersion,
			MP:              1,
			S:               pfcp.SEID_PRESENT,
			MessageType:     pfcp.PFCP_SESSION_ESTABLISHMENT_RESPONSE,
			SEID:            123456789123456789,
			SequenceNumber:  1,
			MessagePriority: 12,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}

func SendPfcpSessionModificationRequest(
	upf *context.UPF,
	ctx *context.SMContext,
	pdrList []*context.PDR,
	farList []*context.FAR,
	barList []*context.BAR,
	qerList []*context.QER,
	urrList []*context.URR,
) (resMsg *pfcpUdp.Message, err error) {
	nodeIDtoIP := upf.NodeID.ResolveNodeIdToIp()
	if err = upf.IsAssociated(); err != nil {
		return nil, err
	}

	pfcpMsg, err := BuildPfcpSessionModificationRequest(upf.NodeID, nodeIDtoIP.String(),
		ctx, pdrList, farList, barList, qerList, urrList)
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Modification Request failed: %v", err)
		return nil, err
	}

	seqNum := getSeqNumber()
	remoteSEID := ctx.PFCPContext[nodeIDtoIP.String()].RemoteSEID
	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:         pfcp.PfcpVersion,
			MP:              1,
			S:               pfcp.SEID_PRESENT,
			MessageType:     pfcp.PFCP_SESSION_MODIFICATION_REQUEST,
			SEID:            remoteSEID,
			SequenceNumber:  seqNum,
			MessagePriority: 12,
		},
		Body: pfcpMsg,
	}

	upaddr := &net.UDPAddr{
		IP:   nodeIDtoIP,
		Port: pfcpUdp.PFCP_PORT,
	}

	resMsg, err = udp.SendPfcpRequest(message, upaddr)
	if err != nil {
		return nil, err
	}

	if resMsg.MessageType() != pfcp.PFCP_SESSION_MODIFICATION_RESPONSE {
		return resMsg, fmt.Errorf("received unexpected type response message: %+v", resMsg.PfcpMessage.Header)
	}

	localSEID := ctx.PFCPContext[nodeIDtoIP.String()].LocalSEID
	if resMsg.PfcpMessage.Header.SEID != localSEID {
		return resMsg, fmt.Errorf("received unexpected SEID response message: %+v, exptcted: %d",
			resMsg.PfcpMessage.Header, localSEID)
	}

	return resMsg, nil
}

// Deprecated: PFCP Session Modification Procedure should be initiated by the CP function
func SendPfcpSessionModificationResponse(addr *net.UDPAddr) {
	pfcpMsg, err := BuildPfcpSessionModificationResponse()
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Modification Response failed: %v", err)
		return
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:         pfcp.PfcpVersion,
			MP:              1,
			S:               pfcp.SEID_PRESENT,
			MessageType:     pfcp.PFCP_SESSION_MODIFICATION_RESPONSE,
			SEID:            123456789123456789,
			SequenceNumber:  1,
			MessagePriority: 12,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}

func SendPfcpSessionDeletionRequest(
	upf *context.UPF,
	ctx *context.SMContext,
) (resMsg *pfcpUdp.Message, err error) {
	nodeIDtoIP := upf.NodeID.ResolveNodeIdToIp()
	if err = upf.IsAssociated(); err != nil {
		return nil, err
	}

	pfcpMsg, err := BuildPfcpSessionDeletionRequest()
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Deletion Request failed: %v", err)
		return nil, err
	}
	seqNum := getSeqNumber()
	remoteSEID := ctx.PFCPContext[nodeIDtoIP.String()].RemoteSEID
	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:         pfcp.PfcpVersion,
			MP:              1,
			S:               pfcp.SEID_PRESENT,
			MessageType:     pfcp.PFCP_SESSION_DELETION_REQUEST,
			SEID:            remoteSEID,
			SequenceNumber:  seqNum,
			MessagePriority: 12,
		},
		Body: pfcpMsg,
	}

	upaddr := &net.UDPAddr{
		IP:   nodeIDtoIP,
		Port: pfcpUdp.PFCP_PORT,
	}

	resMsg, err = udp.SendPfcpRequest(message, upaddr)
	if err != nil {
		return nil, err
	}

	if resMsg.MessageType() != pfcp.PFCP_SESSION_DELETION_RESPONSE {
		return resMsg, fmt.Errorf("received unexpected type response message: %+v", resMsg.PfcpMessage.Header)
	}

	localSEID := ctx.PFCPContext[nodeIDtoIP.String()].LocalSEID
	if resMsg.PfcpMessage.Header.SEID != localSEID {
		return resMsg, fmt.Errorf("received unexpected SEID response message: %+v, exptcted: %d",
			resMsg.PfcpMessage.Header, localSEID)
	}

	return resMsg, nil
}

// Deprecated: PFCP Session Deletion Procedure should be initiated by the CP function
func SendPfcpSessionDeletionResponse(addr *net.UDPAddr) {
	pfcpMsg, err := BuildPfcpSessionDeletionResponse()
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Deletion Response failed: %v", err)
		return
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:         pfcp.PfcpVersion,
			MP:              1,
			S:               pfcp.SEID_PRESENT,
			MessageType:     pfcp.PFCP_SESSION_DELETION_RESPONSE,
			SEID:            123456789123456789,
			SequenceNumber:  1,
			MessagePriority: 12,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}

func SendPfcpSessionReportResponse(addr *net.UDPAddr, cause pfcpType.Cause, seqFromUPF uint32, seid uint64) {
	pfcpMsg, err := BuildPfcpSessionReportResponse(cause)
	if err != nil {
		logger.PfcpLog.Errorf("build PFCP Session Report Response failed: %v", err)
		return
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_PRESENT,
			MessageType:    pfcp.PFCP_SESSION_REPORT_RESPONSE,
			SequenceNumber: seqFromUPF,
			SEID:           seid,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}

func SendPfcpHeartbeatRequest(upf *context.UPF) (resMsg *pfcpUdp.Message, err error) {
	pfcpMsg, err := BuildPfcpHeartbeatRequest()
	if err != nil {
		return nil, fmt.Errorf("build PFCP Heartbeat Request failed: %w", err)
	}

	reqMsg := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_NOT_PRESENT,
			MessageType:    pfcp.PFCP_HEARTBEAT_REQUEST,
			SequenceNumber: getSeqNumber(),
		},
		Body: pfcpMsg,
	}

	upfAddr := &net.UDPAddr{
		IP:   upf.NodeID.ResolveNodeIdToIp(),
		Port: pfcpUdp.PFCP_PORT,
	}

	resMsg, err = udp.SendPfcpRequest(reqMsg, upfAddr)
	if err != nil {
		return nil, err
	}

	if resMsg.MessageType() != pfcp.PFCP_HEARTBEAT_RESPONSE {
		return resMsg, fmt.Errorf("received unexpected response message")
	}

	return resMsg, nil
}

func SendHeartbeatResponse(addr *net.UDPAddr, seq uint32) {
	pfcpMsg := pfcp.HeartbeatResponse{
		RecoveryTimeStamp: &pfcpType.RecoveryTimeStamp{
			RecoveryTimeStamp: udp.ServerStartTime,
		},
	}

	message := &pfcp.Message{
		Header: pfcp.Header{
			Version:        pfcp.PfcpVersion,
			MP:             0,
			S:              pfcp.SEID_NOT_PRESENT,
			MessageType:    pfcp.PFCP_HEARTBEAT_RESPONSE,
			SequenceNumber: seq,
		},
		Body: pfcpMsg,
	}

	udp.SendPfcpResponse(message, addr)
}
