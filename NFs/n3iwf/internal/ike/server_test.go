package ike

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/require"

	ike_message "github.com/free5gc/ike/message"
	n3iwf_context "github.com/free5gc/n3iwf/internal/context"
	"github.com/free5gc/n3iwf/internal/ngap"
	"github.com/free5gc/n3iwf/pkg/factory"
)

type n3iwfTestApp struct {
	cfg        *factory.Config
	n3iwfCtx   *n3iwf_context.N3IWFContext
	ngapServer *ngap.Server
	ikeServer  *Server
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
}

func (a *n3iwfTestApp) Config() *factory.Config {
	return a.cfg
}

func (a *n3iwfTestApp) Context() *n3iwf_context.N3IWFContext {
	return a.n3iwfCtx
}

func (a *n3iwfTestApp) CancelContext() context.Context {
	return a.ctx
}

func (a *n3iwfTestApp) SendNgapEvt(evt n3iwf_context.NgapEvt) {
	a.ngapServer.SendNgapEvt(evt)
}

func NewN3iwfTestApp(cfg *factory.Config) (*n3iwfTestApp, error) {
	var err error
	ctx, cancel := context.WithCancel(context.Background())

	n3iwfApp := &n3iwfTestApp{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}

	n3iwfApp.n3iwfCtx, err = n3iwf_context.NewTestContext(n3iwfApp)
	if err != nil {
		return nil, err
	}
	return n3iwfApp, err
}

func NewTestCfg() *factory.Config {
	return &factory.Config{
		Configuration: &factory.Configuration{},
	}
}

func TestHandleNattMsg(t *testing.T) {
	initiatorSPI := uint64(0x123)
	ikeMessage := ike_message.NewMessage(initiatorSPI, 0, ike_message.IKE_SA_INIT,
		true, false, 0, nil)
	pkt, err := ikeMessage.Encode()
	require.NoError(t, err)

	NonESPPkt := append([]byte{0, 0, 0, 0}, pkt...)

	tests := []struct {
		name         string
		conn         *net.UDPConn
		rcvPkt       []byte
		lAddr, rAddr *net.UDPAddr
		msg          *ike_message.IKEMessage
		expectedErr  bool
	}{
		{
			name:   "Received NAT-T Keepalive",
			rcvPkt: []byte{0xff},
			lAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.2"),
				Port: 4500,
			},
			rAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.1"),
				Port: 4500,
			},
			expectedErr: false,
		},
		{
			name:   "Received NAT-T Msg is too short",
			rcvPkt: []byte{0x01, 0x02},
			lAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.2"),
				Port: 4500,
			},
			rAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.1"),
				Port: 4500,
			},
			expectedErr: true,
		},
		{
			name:   "Received IKE packet from port 4500, and no need to drop",
			rcvPkt: NonESPPkt,
			lAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.2"),
				Port: 4500,
			},
			rAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.1"),
				Port: 4500,
			},
			expectedErr: false,
		},
		{
			name:   "Received ESP packet from port 4500",
			rcvPkt: []byte{0x1, 0x2, 0x3, 0x4, 0x5},
			lAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.2"),
				Port: 4500,
			},
			rAddr: &net.UDPAddr{
				IP:   net.ParseIP("10.100.100.1"),
				Port: 4500,
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handleNattMsg(tt.rcvPkt, tt.rAddr, tt.lAddr, nil)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckIKEMessage(t *testing.T) {
	n3iwf, err := NewN3iwfTestApp(NewTestCfg())
	require.NoError(t, err)

	n3iwf.ikeServer, err = NewServer(n3iwf)
	require.NoError(t, err)
	ikeServer := n3iwf.ikeServer

	srcIP := &net.UDPAddr{
		IP:   net.ParseIP("10.100.100.1"),
		Port: 500,
	}
	dstIP := &net.UDPAddr{
		IP:   net.ParseIP("10.100.100.2"),
		Port: 500,
	}

	mockConn, err := net.DialUDP("udp", nil, dstIP)
	require.NoError(t, err)

	initiatorSPI := uint64(0x123)
	nonceData := []byte("randomNonce")
	payload := new(ike_message.IKEPayloadContainer)
	payload.BuildNonce(nonceData)

	ikeMsg := ike_message.NewMessage(initiatorSPI, 0, ike_message.IKE_SA_INIT,
		true, false, 0, *payload)

	tests := []struct {
		name        string
		conn        *net.UDPConn
		localAddr   *net.UDPAddr
		remoteAddr  *net.UDPAddr
		msg         *ike_message.IKEMessage
		expectedErr bool
	}{
		{
			name:       "Receive packet has IKE version error",
			conn:       mockConn,
			localAddr:  dstIP,
			remoteAddr: srcIP,
			msg: &ike_message.IKEMessage{
				IKEHeader: &ike_message.IKEHeader{
					InitiatorSPI: initiatorSPI,
					ExchangeType: ike_message.IKE_SA_INIT,
					Flags:        0,
					MajorVersion: 3,
					MinorVersion: 0,
				},
			},
			expectedErr: true,
		},
		{
			name:        "Decode IKE_SA_INIT msg",
			conn:        mockConn,
			localAddr:   dstIP,
			remoteAddr:  srcIP,
			msg:         ikeMsg,
			expectedErr: false,
		},
		{
			name:       "SPI not found from IKE header",
			conn:       mockConn,
			localAddr:  dstIP,
			remoteAddr: srcIP,
			msg: &ike_message.IKEMessage{
				IKEHeader: &ike_message.IKEHeader{
					InitiatorSPI: initiatorSPI,
					ExchangeType: ike_message.IKE_AUTH,
					Flags:        ike_message.ResponseBitCheck,
					MajorVersion: 2,
					MinorVersion: 0,
				},
				Payloads: ikeMsg.Payloads,
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := tt.msg.Encode()
			require.NoError(t, err)

			_, _, err = ikeServer.checkIKEMessage(
				msg, tt.conn, tt.localAddr, tt.remoteAddr)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConstructPacketWithESP(t *testing.T) {
	srcIP := &net.UDPAddr{
		IP: net.IPv4(192, 168, 0, 1),
	}
	dstIP := &net.UDPAddr{
		IP: net.IPv4(192, 168, 0, 2),
	}

	espPacket := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}

	packet, err := constructPacketWithESP(srcIP, dstIP, espPacket)
	require.NoError(t, err)

	packetParsed := gopacket.NewPacket(packet, layers.LayerTypeIPv4, gopacket.Default)
	ipLayer := packetParsed.Layer(layers.LayerTypeIPv4)
	require.NotNil(t, ipLayer)

	ipv4, _ := ipLayer.(*layers.IPv4)
	require.Equal(t, ipv4.SrcIP.To4().String(), srcIP.IP.String())
	require.Equal(t, ipv4.DstIP.To4().String(), dstIP.IP.String())
	require.Equal(t, ipv4.Protocol, layers.IPProtocolESP)
}
