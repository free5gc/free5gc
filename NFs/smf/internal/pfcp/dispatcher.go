package pfcp

import (
	"github.com/free5gc/pfcp"
	"github.com/free5gc/pfcp/pfcpUdp"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/internal/pfcp/handler"
)

func Dispatch(msg *pfcpUdp.Message) {
	switch msg.PfcpMessage.Header.MessageType {
	case pfcp.PFCP_HEARTBEAT_REQUEST:
		handler.HandlePfcpHeartbeatRequest(msg)
	case pfcp.PFCP_PFD_MANAGEMENT_REQUEST:
		handler.HandlePfcpPfdManagementRequest(msg)
	case pfcp.PFCP_ASSOCIATION_SETUP_REQUEST:
		handler.HandlePfcpAssociationSetupRequest(msg)
	case pfcp.PFCP_ASSOCIATION_UPDATE_REQUEST:
		handler.HandlePfcpAssociationUpdateRequest(msg)
	case pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST:
		handler.HandlePfcpAssociationReleaseRequest(msg)
	case pfcp.PFCP_NODE_REPORT_REQUEST:
		handler.HandlePfcpNodeReportRequest(msg)
	case pfcp.PFCP_SESSION_SET_DELETION_REQUEST:
		handler.HandlePfcpSessionSetDeletionRequest(msg)
	case pfcp.PFCP_SESSION_REPORT_REQUEST:
		handler.HandlePfcpSessionReportRequest(msg)
	default:
		logger.PfcpLog.Errorf("Unknown PFCP message type: %d", msg.PfcpMessage.Header.MessageType)
		return
	}
}
