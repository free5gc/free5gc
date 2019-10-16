package pfcp_message_test

import (
	"free5gc/src/smf/smf_context"
	"net"
	"testing"
	"time"

	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/smf_pfcp/pfcp_message"
	"free5gc/src/smf/smf_pfcp/pfcp_udp"
)

var testAddr *net.UDPAddr

// Adjust waiting time in millisecond if PFCP packets are not captured
var testWaitingTime int = 500

var dummyContext *smf_context.SMContext

func init() {
	pfcp_udp.ServerNodeId = pfcpType.NodeID{
		NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
		NodeIdValue: net.ParseIP("127.0.0.1").To4(),
	}

	pfcp_udp.Run()

	testAddr = &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: pfcpUdp.PFCP_PORT,
	}
	dummyContext = smf_context.NewSMContext("imsi-20893000000001", 3)

}

func TestSendPfcpAssociationSetupRequest(t *testing.T) {
	pfcp_message.SendPfcpAssociationSetupRequest(testAddr)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpAssociationSetupResponse(t *testing.T) {
	cause := pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}
	pfcp_message.SendPfcpAssociationSetupResponse(testAddr, cause)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpAssociationReleaseRequest(t *testing.T) {
	pfcp_message.SendPfcpAssociationReleaseRequest(testAddr)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpAssociationReleaseResponse(t *testing.T) {
	cause := pfcpType.Cause{
		CauseValue: pfcpType.CauseRequestAccepted,
	}
	pfcp_message.SendPfcpAssociationReleaseResponse(testAddr, cause)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpSessionEstablishmentRequest(t *testing.T) {
	pfcp_message.SendPfcpSessionEstablishmentRequest(testAddr, dummyContext)
	time.Sleep(time.Duration(testWaitingTime) * time.Millisecond)
}

func TestSendPfcpSessionEstablishmentResponse(t *testing.T) {
	pfcp_message.SendPfcpSessionEstablishmentResponse(testAddr)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpSessionModificationRequest(t *testing.T) {
	pfcp_message.SendPfcpSessionModificationRequest(testAddr)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpSessionModificationResponse(t *testing.T) {
	pfcp_message.SendPfcpSessionModificationResponse(testAddr)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpSessionDeletionRequest(t *testing.T) {
	pfcp_message.SendPfcpSessionDeletionRequest(testAddr)
	time.Sleep(1000 * time.Millisecond)
}

func TestSendPfcpSessionDeletionResponse(t *testing.T) {
	pfcp_message.SendPfcpSessionDeletionResponse(testAddr)
	time.Sleep(1000 * time.Millisecond)
}
