package n3iwf_data_relay_test

import (
	"free5gc/lib/ngap/ngapType"
	"free5gc/src/n3iwf/n3iwf_context"
	"free5gc/src/n3iwf/n3iwf_data_relay"
	"free5gc/src/n3iwf/n3iwf_handler"
	"sync"
	"testing"
	"time"
)

func TestUserPlaneRelay(t *testing.T) {
	// Prepare N3IWF context
	n3iwfSelf := n3iwf_context.N3IWFSelf()

	n3iwfSelf.IPSecGatewayAddress = "10.0.0.1"

	// N3IWF handler
	go n3iwf_handler.Handle()

	// Listen N1 interface
	if err := n3iwf_data_relay.ListenN1UPTraffic(); err != nil {
		t.Fatal(err)
	}

	// Prepare UE context
	ue := n3iwfSelf.NewN3iwfUe()

	// Add UE inner IP to context
	ue.IPSecInnerIP = "10.0.0.2"
	n3iwfSelf.AllocatedUEIPAddress["10.0.0.2"] = ue

	// GTP address is acquired from SMF
	n3iwfSelf.GTPBindAddress = "172.31.0.153"

	userPlaneConn, remoteAddr, err := n3iwf_data_relay.SetupGTPTunnelWithUPF("172.31.0.152")
	if err != nil {
		t.Fatal(err)
	}

	ueTEID := n3iwfSelf.NewTEID(ue)

	// New PDU session
	pduSession, err := ue.CreatePDUSession(1, ngapType.SNSSAI{})
	if err != nil {
		t.Fatal(err)
	}

	ueAssociatedGTPConnection := &n3iwf_context.GTPConnectionInfo{
		UPFUDPAddr:          remoteAddr,
		IncomingTEID:        ueTEID,
		OutgoingTEID:        1,
		UserPlaneConnection: userPlaneConn,
	}

	pduSession.GTPConnection = ueAssociatedGTPConnection

	// Listen GTP
	if err := n3iwf_data_relay.ListenGTP(userPlaneConn); err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func TestControlPlaneRelay(t *testing.T) {
	// Prepare N3IWF context
	n3iwfSelf := n3iwf_context.N3IWFSelf()
	n3iwfSelf.IPSecGatewayAddress = "10.0.0.1"
	n3iwfSelf.TCPPort = 20000

	// N3IWF handler
	go n3iwf_handler.Handle()

	// Set up N1 NAS TCP server
	if err := n3iwf_data_relay.SetupNASTCPServer(); err != nil {
		t.Fatal(err)
	}

	// Prepare UE context
	ue := n3iwfSelf.NewN3iwfUe()
	ue.IPSecInnerIP = "10.0.0.2"

	n3iwfSelf.AllocatedUEIPAddress[ue.IPSecInnerIP] = ue

	for {
		if ue.TCPConnection != nil {
			n3iwf_data_relay.ForwardCPTrafficFromN2(ue, []byte("Hello UE\n"))
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
