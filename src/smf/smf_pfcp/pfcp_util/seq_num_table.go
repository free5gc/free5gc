package pfcp_util

import (

	//"free5gc/src/smf/smf_pfcp/pfcp_udp"
	"free5gc/lib/pfcp"
	"free5gc/src/smf/logger"
)

type SeqNumTableItem struct {
	PacketState PacketState
	MessageType pfcp.MessageType
}

type SeqNumTable struct {
	fromSMF map[uint32]*SeqNumTableItem
	toSMF   map[uint32]*SeqNumTableItem
}

func NewSeqNumTable() *SeqNumTable {
	var snt SeqNumTable
	snt.fromSMF = make(map[uint32]*SeqNumTableItem)
	snt.toSMF = make(map[uint32]*SeqNumTableItem)
	return &snt
}

func (snt SeqNumTable) RecvCheckAndPutItem(msg *pfcp.Message) (Success bool) {

	Item := new(SeqNumTableItem)
	seqNum := msg.Header.SequenceNumber
	Success = false
	switch msg.Header.MessageType {
	case pfcp.PFCP_HEARTBEAT_REQUEST:
		logger.PfcpLog.Warnf("PFCP Heartbeat Request handling is not implemented")
	case pfcp.PFCP_HEARTBEAT_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Heartbeat Response handling is not implemented")
	case pfcp.PFCP_PFD_MANAGEMENT_REQUEST:
		logger.PfcpLog.Warnf("PFCP PFD Management Request handling is not implemented")
	case pfcp.PFCP_PFD_MANAGEMENT_RESPONSE:
		logger.PfcpLog.Warnf("PFCP PFD Management Response handling is not implemented")
	case pfcp.PFCP_ASSOCIATION_SETUP_REQUEST:
		Item.PacketState = RECV_REQUEST
		Item.MessageType = pfcp.PFCP_ASSOCIATION_SETUP_REQUEST
	case pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(RECV_RESPONSE))
		return
	case pfcp.PFCP_ASSOCIATION_UPDATE_REQUEST:
		Item.PacketState = RECV_REQUEST
		Item.MessageType = pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE
	case pfcp.PFCP_ASSOCIATION_UPDATE_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Association Update Response handling is not implemented")
	case pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST:
		Item.PacketState = RECV_REQUEST
		Item.MessageType = pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST
	case pfcp.PFCP_ASSOCIATION_RELEASE_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(RECV_RESPONSE))
		return
	case pfcp.PFCP_VERSION_NOT_SUPPORTED_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Version Not Support Response handling is not implemented")
	case pfcp.PFCP_NODE_REPORT_REQUEST:
		logger.PfcpLog.Warnf("PFCP Node Report Request handling is not implemented")
	case pfcp.PFCP_NODE_REPORT_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Node Report Response handling is not implemented")
	case pfcp.PFCP_SESSION_SET_DELETION_REQUEST:
		logger.PfcpLog.Warnf("PFCP Session Set Deletion Request handling is not implemented")
	case pfcp.PFCP_SESSION_SET_DELETION_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Session Set Deletion Response handling is not implemented")
	case pfcp.PFCP_SESSION_ESTABLISHMENT_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(RECV_RESPONSE))
		return
	case pfcp.PFCP_SESSION_MODIFICATION_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(RECV_RESPONSE))
		return
	case pfcp.PFCP_SESSION_DELETION_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Session Deletion Response handling is not implemented")
	case pfcp.PFCP_SESSION_REPORT_REQUEST:
		logger.PfcpLog.Warnf("PFCP Session Report Response handling is not implemented")
	case pfcp.PFCP_SESSION_REPORT_RESPONSE:
		logger.PfcpLog.Warnf("PFCP Session Report Response handling is not implemented")
	default:
		logger.PfcpLog.Errorf("Unknown PFCP message type: %d", msg.Header.MessageType)
		return

	}

	_, exist := snt.toSMF[seqNum]
	if !exist {
		snt.toSMF[seqNum] = Item
		Success = true
	} else {
		logger.PfcpLog.Errorf("\n[SMF PFCP]Sequence Number %d already exists.\n", seqNum)
		logger.PfcpLog.Errorf("\n[SMF PFCP]Message Type %d\n", Item.MessageType)
	}

	return
}

func (snt SeqNumTable) SendCheckAndPutItem(msg *pfcp.Message) (Success bool) {
	Item := new(SeqNumTableItem)
	seqNum := msg.Header.SequenceNumber
	Success = false

	switch msg.Header.MessageType {
	case pfcp.PFCP_ASSOCIATION_SETUP_REQUEST:
		Item.PacketState = SEND_REQUEST
		Item.MessageType = pfcp.PFCP_ASSOCIATION_SETUP_REQUEST
	case pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(SEND_RESPONSE))
		return
	case pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST:
		Item.PacketState = SEND_REQUEST
		Item.MessageType = pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST
	case pfcp.PFCP_ASSOCIATION_RELEASE_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(SEND_RESPONSE))
		return
	case pfcp.PFCP_SESSION_ESTABLISHMENT_REQUEST:
		Item.PacketState = SEND_REQUEST
		Item.MessageType = pfcp.PFCP_SESSION_ESTABLISHMENT_REQUEST
	case pfcp.PFCP_SESSION_ESTABLISHMENT_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(SEND_RESPONSE))
		return
	case pfcp.PFCP_SESSION_MODIFICATION_REQUEST:
		Item.PacketState = SEND_REQUEST
		Item.MessageType = pfcp.PFCP_SESSION_MODIFICATION_REQUEST
	case pfcp.PFCP_SESSION_MODIFICATION_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(SEND_RESPONSE))
		return
	case pfcp.PFCP_SESSION_DELETION_REQUEST:
		Item.PacketState = SEND_REQUEST
		Item.MessageType = pfcp.PFCP_SESSION_DELETION_REQUEST
	case pfcp.PFCP_SESSION_DELETION_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(SEND_RESPONSE))
		return
	case pfcp.PFCP_SESSION_REPORT_RESPONSE:
		Success = snt.RemoveItem(seqNum, uint8(SEND_RESPONSE))
		return
	default:
		logger.PfcpLog.Errorf("\nUnknown PFCP message type: %d\n", msg.Header.MessageType)
		return

	}

	_, exist := snt.fromSMF[seqNum]
	if !exist {
		snt.fromSMF[seqNum] = Item
		Success = true
	} else {
		logger.PfcpLog.Errorf("\n[SMF PFCP]Sequence Number %d already exists.\n", seqNum)
		logger.PfcpLog.Errorf("\n[SMF PFCP]Message Type %d\n", Item.MessageType)
	}

	return
}

func (snt SeqNumTable) RemoveItem(seqNum uint32, newStateInInt uint8) (Success bool) {
	newState := PacketState(newStateInInt)
	Success = false

	var item *SeqNumTableItem
	var exist bool

	if newState == SEND_RESPONSE {
		item, exist = snt.toSMF[seqNum]

		if !exist {
			logger.PfcpLog.Warnf("\n[SMF PFCP] Can't send response without having corresponding request.\n")
			logger.PfcpLog.Warnf("\n[SMF PFCP] Packet sequence number: %d\n", seqNum)
			return
		}
	} else if newState == RECV_RESPONSE {
		item, exist = snt.fromSMF[seqNum]

		if !exist {
			logger.PfcpLog.Warnf("\n[SMF PFCP] Can't receive response without having corresponding request.\n")
			logger.PfcpLog.Warnf("\n[SMF PFCP] Packet sequence number: %d\n", seqNum)
			return
		}
	}

	switch newState {
	case SEND_RESPONSE:
		if item.PacketState != RECV_REQUEST {
			logger.PfcpLog.Warnf("\n[SMF PFCP] Wrong Packet State when sending response.\n")
			logger.PfcpLog.Warnf("\n[SMF PFCP] Respone message type %d\n", item.MessageType)
			logger.PfcpLog.Warnf("\n[SMF PFCP] Packet State: %d\n", item.PacketState)
			logger.PfcpLog.Warnf("\n[SMF PFCP] Packet sequence number: %d\n", seqNum)
			return
		}

		delete(snt.toSMF, seqNum)
		Success = true
	case RECV_RESPONSE:
		if item.PacketState != SEND_REQUEST {
			logger.PfcpLog.Warnf("\n[SMF PFCP] Wrong Packet State when receiving response.\n")
			logger.PfcpLog.Warnf("\n[SMF PFCP] Respone message type %d\n", item.MessageType)
			logger.PfcpLog.Warnf("\n[SMF PFCP] Packet State: %d\n", item.PacketState)
			logger.PfcpLog.Warnf("\n[SMF PFCP] Packet sequence number: %d\n", seqNum)
			return
		}

		delete(snt.fromSMF, seqNum)
		Success = true
	default:
		logger.PfcpLog.Errorf("\nWrong Packet State: %d\n", newState)
	}

	return
}
