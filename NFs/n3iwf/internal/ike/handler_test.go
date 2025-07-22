package ike

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	ike_message "github.com/free5gc/ike/message"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/pkg/factory"
	"github.com/free5gc/util/ippool"
)

func TestRemoveIkeUe(t *testing.T) {
	n3iwf, err := NewN3iwfTestApp(&factory.Config{})
	require.NoError(t, err)

	n3iwf.ikeServer, err = NewServer(n3iwf)
	require.NoError(t, err)

	n3iwfCtx := n3iwf.n3iwfCtx
	ikeSA := n3iwfCtx.NewIKESecurityAssociation()
	ikeUe := n3iwfCtx.NewN3iwfIkeUe(ikeSA.LocalSPI)
	ikeUe.N3IWFIKESecurityAssociation = ikeSA
	ikeUe.IPSecInnerIP = net.ParseIP("10.0.0.1")
	ikeSA.IsUseDPD = false

	n3iwfCtx.IPSecInnerIPPool, err = ippool.NewIPPool("10.0.0.0/24")
	require.NoError(t, err)
	_, err = n3iwfCtx.IPSecInnerIPPool.Allocate(nil)
	require.NoError(t, err)

	ikeUe.CreateHalfChildSA(1, 123, 1)

	ikeAuth := &ike_message.SecurityAssociation{}

	proposal := ikeAuth.Proposals.BuildProposal(1, 1, []byte{0, 1, 2, 3})
	var attributeType uint16 = ike_message.AttributeTypeKeyLength
	var attributeValue uint16 = 256
	proposal.EncryptionAlgorithm.BuildTransform(ike_message.TypeEncryptionAlgorithm,
		ike_message.ENCR_AES_CBC, &attributeType, &attributeValue, nil)

	proposal.IntegrityAlgorithm.BuildTransform(ike_message.TypeIntegrityAlgorithm,
		ike_message.AUTH_HMAC_SHA1_96, nil, nil, nil)

	proposal.ExtendedSequenceNumbers.BuildTransform(
		ike_message.TypeExtendedSequenceNumbers, ike_message.ESN_DISABLE, nil, nil, nil)

	childSA, err := ikeUe.CompleteChildSA(1, 456, ikeAuth)
	require.NoError(t, err)

	err = n3iwf.ikeServer.removeIkeUe(ikeSA.LocalSPI)
	require.NoError(t, err)

	_, ok := n3iwfCtx.IkeUePoolLoad(ikeSA.LocalSPI)
	require.False(t, ok)

	_, ok = n3iwfCtx.IKESALoad(ikeSA.LocalSPI)
	require.False(t, ok)

	_, ok = ikeUe.N3IWFChildSecurityAssociation[childSA.InboundSPI]
	require.False(t, ok)
}

func TestGenerateNATDetectHash(t *testing.T) {
	n3iwf, err := NewN3iwfTestApp(&factory.Config{})
	require.NoError(t, err)

	n3iwf.ikeServer, err = NewServer(n3iwf)
	require.NoError(t, err)

	tests := []struct {
		name         string
		initiatorSPI uint64
		responderSPI uint64
		Addr         net.UDPAddr
		expectedData []byte
	}{
		{
			name:         "Generate NAT-D hash",
			initiatorSPI: 0x1122334455667788,
			responderSPI: 0xaabbeeddeeff1122,
			Addr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 4500,
			},
			expectedData: []byte{
				0xd2, 0xee, 0x40, 0x2d, 0x5d, 0x53, 0xe4, 0x4a,
				0x01, 0x2d, 0x44, 0x2a, 0x90, 0x05, 0xc1, 0xea,
				0x38, 0x8a, 0x81, 0x7e,
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			data, err := n3iwf.ikeServer.generateNATDetectHash(tt.initiatorSPI, tt.responderSPI, &tt.Addr)
			require.NoError(t, err)

			require.Equal(t, tt.expectedData, data)
		})
	}
}

func TestBuildNATDetectMsg(t *testing.T) {
	n3iwf, err := NewN3iwfTestApp(&factory.Config{})
	require.NoError(t, err)

	n3iwf.ikeServer, err = NewServer(n3iwf)
	require.NoError(t, err)

	remoteSPI := uint64(0x1234567890abcdef)
	localSPI := uint64(0xfedcba0987654321)
	ikeSA := &n3iwf_context.IKESecurityAssociation{
		LocalSPI:  localSPI,
		RemoteSPI: remoteSPI,
	}
	payload := &ike_message.IKEPayloadContainer{}

	ueAddr := net.UDPAddr{
		IP:   net.ParseIP("192.168.1.1"),
		Port: 4500,
	}
	n3iwfAddr := net.UDPAddr{
		IP:   net.ParseIP("192.168.1.2"),
		Port: 4500,
	}

	err = n3iwf.ikeServer.buildNATDetectNotifPayload(ikeSA, payload, &ueAddr, &n3iwfAddr)
	require.NoError(t, err)

	var notifications []*ike_message.Notification
	for _, ikePayload := range *payload {
		switch ikePayload.Type() {
		case ike_message.TypeN:
			notifications = append(notifications, ikePayload.(*ike_message.Notification))
		default:
			require.Fail(t, "Get unexpected IKE payload type : %v", ikePayload.Type())
		}
	}

	for _, notification := range notifications {
		switch notification.NotifyMessageType {
		case ike_message.NAT_DETECTION_SOURCE_IP:
			expectedData := []byte{
				0x13, 0xd8, 0x9e, 0xdc, 0xfa, 0x39, 0xe4, 0xc0,
				0x06, 0x80, 0x5f, 0xde, 0x11, 0x62, 0xd8, 0x76,
				0xee, 0xe8, 0xf2, 0x00,
			}
			require.Equal(t, expectedData, notification.NotificationData)
		case ike_message.NAT_DETECTION_DESTINATION_IP:
			expectedData := []byte{
				0x0d, 0x36, 0x26, 0x71, 0xaf, 0x7f, 0x0b, 0x19,
				0x32, 0xec, 0xf8, 0xf3, 0xe1, 0x84, 0x87, 0xf0,
				0x47, 0x76, 0x83, 0x04,
			}
			require.Equal(t, expectedData, notification.NotificationData)
		}
	}
}

func TestHandleNATDetect(t *testing.T) {
	n3iwf, err := NewN3iwfTestApp(&factory.Config{})
	require.NoError(t, err)

	n3iwf.ikeServer, err = NewServer(n3iwf)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		initiatorSPI           uint64
		responderSPI           uint64
		notification           []*ike_message.Notification
		ueAddr                 net.UDPAddr
		n3iwfAddr              net.UDPAddr
		expectedUeBehindNAT    bool
		expectedN3iwfBehindNAT bool
	}{
		{
			name:         "UE and N3IWF is not behind NAT",
			initiatorSPI: 0x1234567890abcdef,
			responderSPI: 0xfedcba0987654321,
			notification: []*ike_message.Notification{
				{
					NotifyMessageType: ike_message.NAT_DETECTION_SOURCE_IP,
					NotificationData: []byte{
						0x0d, 0x36, 0x26, 0x71, 0xaf, 0x7f, 0x0b, 0x19,
						0x32, 0xec, 0xf8, 0xf3, 0xe1, 0x84, 0x87, 0xf0,
						0x47, 0x76, 0x83, 0x04,
					},
				},
				{
					NotifyMessageType: ike_message.NAT_DETECTION_DESTINATION_IP,
					NotificationData: []byte{
						0x13, 0xd8, 0x9e, 0xdc, 0xfa, 0x39, 0xe4, 0xc0,
						0x06, 0x80, 0x5f, 0xde, 0x11, 0x62, 0xd8, 0x76,
						0xee, 0xe8, 0xf2, 0x00,
					},
				},
			},
			ueAddr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 4500,
			},
			n3iwfAddr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.2"),
				Port: 4500,
			},
			expectedUeBehindNAT:    false,
			expectedN3iwfBehindNAT: false,
		},
		{
			name:         "UE is behind NAT and N3IWF is not behind NAT",
			initiatorSPI: 0x1234567890abcdef,
			responderSPI: 0xfedcba0987654321,
			notification: []*ike_message.Notification{
				{
					NotifyMessageType: ike_message.NAT_DETECTION_SOURCE_IP,
					NotificationData: []byte{
						0x0b, 0x17, 0x2d, 0x42, 0xaf, 0x7f, 0x0b, 0x19,
						0x32, 0xec, 0xf8, 0xf3, 0xe1, 0x84, 0x87, 0xf0,
						0x47, 0x76, 0x83, 0x04,
					},
				},
				{
					NotifyMessageType: ike_message.NAT_DETECTION_DESTINATION_IP,
					NotificationData: []byte{
						0x13, 0xd8, 0x9e, 0xdc, 0xfa, 0x39, 0xe4, 0xc0,
						0x06, 0x80, 0x5f, 0xde, 0x11, 0x62, 0xd8, 0x76,
						0xee, 0xe8, 0xf2, 0x00,
					},
				},
			},
			ueAddr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 4500,
			},
			n3iwfAddr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.2"),
				Port: 4500,
			},
			expectedUeBehindNAT:    true,
			expectedN3iwfBehindNAT: false,
		},
		{
			name:         "UE and N3IWF is behind NAT",
			initiatorSPI: 0x1234567890abcdef,
			responderSPI: 0xfedcba0987654321,
			notification: []*ike_message.Notification{
				{
					NotifyMessageType: ike_message.NAT_DETECTION_SOURCE_IP,
					NotificationData: []byte{
						0x0b, 0x16, 0x26, 0x71, 0xaf, 0x7f, 0x0b, 0x19,
						0x32, 0xec, 0xf8, 0xf3, 0xe1, 0x84, 0x87, 0xf0,
						0x47, 0x76, 0x83, 0x04,
					},
				},
				{
					NotifyMessageType: ike_message.NAT_DETECTION_DESTINATION_IP,
					NotificationData: []byte{
						0x0f, 0xd9, 0x9e, 0xdc, 0xfa, 0x39, 0xe4, 0xc0,
						0x06, 0x80, 0x5f, 0xde, 0x11, 0x62, 0xd8, 0x76,
						0xee, 0xe8, 0xf2, 0x00,
					},
				},
			},
			ueAddr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 4500,
			},
			n3iwfAddr: net.UDPAddr{
				IP:   net.ParseIP("192.168.1.2"),
				Port: 4500,
			},
			expectedUeBehindNAT:    true,
			expectedN3iwfBehindNAT: true,
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			ueBehindNAT, n3iwfBehindNAT, err := n3iwf.ikeServer.handleNATDetect(
				tt.initiatorSPI, tt.responderSPI,
				tt.notification, &tt.ueAddr, &tt.n3iwfAddr)
			require.NoError(t, err)

			require.Equal(t, tt.expectedUeBehindNAT, ueBehindNAT)
			require.Equal(t, tt.expectedN3iwfBehindNAT, n3iwfBehindNAT)
		})
	}
}
