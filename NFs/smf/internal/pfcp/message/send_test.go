package message_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	smf_context "github.com/free5gc/smf/internal/context"
	smf_pfcp "github.com/free5gc/smf/internal/pfcp"
	"github.com/free5gc/smf/internal/pfcp/message"
	"github.com/free5gc/smf/internal/pfcp/udp"
)

func TestSendPfcpAssociationSetupRequest(t *testing.T) {
}

func TestSendPfcpSessionEstablishmentResponse(t *testing.T) {
}

func TestSendPfcpSessionEstablishmentRequest(t *testing.T) {
}

func TestSendHeartbeatResponse(t *testing.T) {
	smfContext := smf_context.GetSelf()
	smfContext.PfcpContext, smfContext.PfcpCancelFunc = context.WithCancel(context.Background())
	udp.Run(smf_pfcp.Dispatch)

	udp.ServerStartTime = time.Now()
	var seq uint32 = 1
	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 7001,
	}
	message.SendHeartbeatResponse(addr, seq)

	err := udp.ClosePfcp()
	require.NoError(t, err)
}
