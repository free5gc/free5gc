package pfcp_handler_test

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpType"
	"free5gc/lib/pfcp/pfcpUdp"
	"free5gc/src/smf/smf_handler"
	"free5gc/src/smf/smf_pfcp/pfcp_udp"
)

const testPfcpClientPort = 12345

func init() {
	pfcp_udp.ServerNodeId = pfcpType.NodeID{
		NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
		NodeIdValue: net.ParseIP("127.0.0.1").To4(),
	}

	pfcp_udp.Run()

	// Reset start time of PFCP server
	pfcp_udp.ServerStartTime = time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC)

	go smf_handler.Handle()
}

func TestHandlePfcpAssociationSetupRequest(t *testing.T) {
	testReq := pfcp.Message{
		Header: pfcp.Header{
			Version:        1,
			MP:             0,
			S:              0,
			MessageType:    pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
			SequenceNumber: 1,
		},
		Body: pfcp.PFCPAssociationSetupRequest{
			NodeID: &pfcpType.NodeID{
				NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
				NodeIdValue: net.ParseIP("192.168.0.4").To4(),
			},
			RecoveryTimeStamp: &pfcpType.RecoveryTimeStamp{
				RecoveryTimeStamp: time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			UPFunctionFeatures: &pfcpType.UPFunctionFeatures{
				SupportedFeatures: pfcpType.UpFunctionFeaturesBucp | pfcpType.UpFunctionFeaturesPdiu,
			},
		},
	}

	testRsp := pfcp.Message{
		Header: pfcp.Header{
			Version:        1,
			MP:             0,
			S:              0,
			MessageType:    pfcp.PFCP_ASSOCIATION_SETUP_RESPONSE,
			MessageLength:  31,
			SequenceNumber: 1,
		},
		Body: pfcp.PFCPAssociationSetupResponse{
			NodeID: &pfcp_udp.ServerNodeId,
			Cause: &pfcpType.Cause{
				CauseValue: pfcpType.CauseRequestAccepted,
			},
			RecoveryTimeStamp: &pfcpType.RecoveryTimeStamp{
				RecoveryTimeStamp: pfcp_udp.ServerStartTime,
			},
			CPFunctionFeatures: &pfcpType.CPFunctionFeatures{
				SupportedFeatures: 0,
			},
		},
	}

	srcAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: testPfcpClientPort,
	}
	dstAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: pfcpUdp.PFCP_PORT,
	}

	err := pfcpUdp.SendPfcpMessage(testReq, srcAddr, dstAddr)
	assert.Nil(t, err)

	var rspMsg pfcp.Message
	err = pfcpUdp.ReceivePfcpMessage(&rspMsg, srcAddr, dstAddr)
	assert.Nil(t, err)

	assert.Equal(t, testRsp, rspMsg)
}

func TestHandlePfcpAssociationReleaseRequest(t *testing.T) {
	testReq := pfcp.Message{
		Header: pfcp.Header{
			Version:        1,
			MP:             0,
			S:              0,
			MessageType:    pfcp.PFCP_ASSOCIATION_RELEASE_REQUEST,
			SequenceNumber: 1,
		},
		Body: pfcp.PFCPAssociationReleaseRequest{
			NodeID: &pfcpType.NodeID{
				NodeIdType:  pfcpType.NodeIdTypeIpv4Address,
				NodeIdValue: net.ParseIP("192.168.0.4").To4(),
			},
		},
	}

	testRsp := pfcp.Message{
		Header: pfcp.Header{
			Version:        1,
			MP:             0,
			S:              0,
			MessageType:    pfcp.PFCP_ASSOCIATION_RELEASE_RESPONSE,
			MessageLength:  18,
			SequenceNumber: 1,
		},
		Body: pfcp.PFCPAssociationReleaseResponse{
			NodeID: &pfcp_udp.ServerNodeId,
			Cause: &pfcpType.Cause{
				// Cause value would be No Established PFCP Association if TestHandlePfcpAssociationReleaseRequest is run alone
				// CauseValue: pfcpType.CauseNoEstablishedPfcpAssociation,
				CauseValue: pfcpType.CauseRequestAccepted,
			},
		},
	}

	srcAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: testPfcpClientPort,
	}
	dstAddr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: pfcpUdp.PfcpUdpDestinationPort,
	}

	err := pfcpUdp.SendPfcpMessage(testReq, srcAddr, dstAddr)
	assert.Nil(t, err)

	var rspMsg pfcp.Message
	err = pfcpUdp.ReceivePfcpMessage(&rspMsg, srcAddr, dstAddr)
	assert.Nil(t, err)

	assert.Equal(t, testRsp, rspMsg)
}
