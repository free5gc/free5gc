package smf_pfcp

import (
	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/logger"
	"free5gc/src/smf/smf_handler/smf_message"
	"free5gc/src/smf/smf_pfcp/pfcp_handler"
)

func Dispatch(msg *pfcpUdp.Message, ResponseQueue *smf_message.ResponseQueue) {
	switch msg.PfcpMessage.Header.MessageType {
	case pfcp.PFCP_HEARTBEAT_REQUEST:
		pfcp_handler.HandlePfcpHeartbeatRequest(msg)
	case pfcp.PFCP_HEARTBEAT_RESPONSE:
		pfcp_handler.HandlePfcpHeartbeatResponse(msg)
	case pfcp.PFCP_PFD_MANAGEMENT_REQUEST:
		pfcp_handler.HandlePfcpPfdManagementRequest(msg)
	case pfcp.PFCP_PFD_MANAGEMENT_RESPONSE:
		pfcp_handler.HandlePfcpPfdManagementResponse(msg)
	case pfcp.PFCP_ASSOCIATION_SETUP_REQUEST:
		pfcp_handler.HandlePfcpAssociationSetupRequest(msg)
	case pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE:
		pfcp_handler.HandlePfcpAssociationSetupResponse(msg)
	case pfcp.PFCP_ASSOCIATION_UPDATE_REQUEST:
		pfcp_handler.HandlePfcpAssociationUpdateRequest(msg)
	case pfcp.PFCP_ASSOCIATION_UPDATE_RESPONSE:
		pfcp_handler.HandlePfcpAssociationUpdateResponse(msg)
	case pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST:
		pfcp_handler.HandlePfcpAssociationReleaseRequest(msg)
	case pfcp.PFCP_ASSOCIATION_RELEASE_RESPONSE:
		pfcp_handler.HandlePfcpAssociationReleaseResponse(msg)
	case pfcp.PFCP_VERSION_NOT_SUPPORTED_RESPONSE:
		pfcp_handler.HandlePfcpVersionNotSupportedResponse(msg)
	case pfcp.PFCP_NODE_REPORT_REQUEST:
		pfcp_handler.HandlePfcpNodeReportRequest(msg)
	case pfcp.PFCP_NODE_REPORT_RESPONSE:
		pfcp_handler.HandlePfcpNodeReportResponse(msg)
	case pfcp.PFCP_SESSION_SET_DELETION_REQUEST:
		pfcp_handler.HandlePfcpSessionSetDeletionRequest(msg)
	case pfcp.PFCP_SESSION_SET_DELETION_RESPONSE:
		pfcp_handler.HandlePfcpSessionSetDeletionResponse(msg)
	case pfcp.PFCP_SESSION_ESTABLISHMENT_RESPONSE:
		pfcp_handler.HandlePfcpSessionEstablishmentResponse(msg)
	case pfcp.PFCP_SESSION_MODIFICATION_RESPONSE:
		pfcp_handler.HandlePfcpSessionModificationResponse(msg, ResponseQueue)
	case pfcp.PFCP_SESSION_DELETION_RESPONSE:
		pfcp_handler.HandlePfcpSessionDeletionResponse(msg)
	case pfcp.PFCP_SESSION_REPORT_REQUEST:
		pfcp_handler.HandlePfcpSessionReportRequest(msg)
	case pfcp.PFCP_SESSION_REPORT_RESPONSE:
		pfcp_handler.HandlePfcpSessionReportResponse(msg)
	default:
		logger.PfcpLog.Errorf("Unknown PFCP message type: %d", msg.PfcpMessage.Header.MessageType)
		return
	}
}
