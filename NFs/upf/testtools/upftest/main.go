package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/wmnsk/go-pfcp/ie"
	"github.com/wmnsk/go-pfcp/message"
)

func main() {
	var (
		server   = flag.String("-s", "127.0.0.8:8805", "server's addr/port")
		nodeid   = flag.String("-n", "127.0.0.8", "client's node id")
		boottime = time.Now()
		seq      uint32
		err      error
		buf      = make([]byte, 1500)
		waiting  bool
	)
	flag.Parse()

	raddr, err := net.ResolveUDPAddr("udp4", *server)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}

	seq += 1
	asreq, err := message.NewAssociationSetupRequest(
		seq, // Sequence Number
		ie.NewNodeID(*nodeid, "", ""),
		ie.NewRecoveryTimeStamp(boottime),
	).Marshal()
	if err != nil {
		log.Fatal(err)
	}

	if _, err = conn.Write(asreq); err != nil {
		log.Fatal(err)
	}
	log.Printf("sent PFCP Association Setup Request to: %s", raddr)

	if err = conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		log.Fatal(err)
	}

	waiting = true
	for waiting {
		n, addr, err1 := conn.ReadFrom(buf)
		if err1 != nil {
			log.Fatal(err1)
		}

		msg, err1 := message.Parse(buf[:n])
		if err1 != nil {
			log.Printf("ignored undecodable message: %x, error: %s", buf[:n], err1)
			continue
		}

		asres, ok := msg.(*message.AssociationSetupResponse)
		if !ok {
			log.Printf("got unexpected message: %s, from: %s", msg.MessageTypeName(), addr)
			continue
		}

		waiting = false
		if asres.Cause == nil {
			log.Printf("got non accepted response")
			return
		}
		if cause, err1 := asres.Cause.Cause(); cause != ie.CauseRequestAccepted || err1 != nil {
			log.Printf("got non accepted response")
			return
		}
	}

	seq += 1
	sereq, err := message.NewSessionEstablishmentRequest(
		1,   // MP(Message Priority) flag
		0,   // FO(Follow On) flag
		0,   // SEID(Session Endpoint Identifier)
		seq, // Sequence Number
		0,   // Message Priority
		ie.NewNodeID(*nodeid, "", ""),
		ie.NewFSEID(1, net.ParseIP(*nodeid), nil),
		ie.NewCreatePDR(
			ie.NewPDRID(1),
			ie.NewPrecedence(255),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceAccess),
				ie.NewFTEID(1 /* flags */, 1 /* teid */, net.ParseIP("172.16.1.1"), nil, 0 /* chid */),
				ie.NewNetworkInstance(""),
				ie.NewUEIPAddress(2, "60.60.0.6", "", 0, 0),
			),
			ie.NewOuterHeaderRemoval(0 /* desc */, 0 /* ext */),
			ie.NewFARID(1),
			ie.NewQERID(1),
		),
		ie.NewCreatePDR(
			ie.NewPDRID(2),
			ie.NewPrecedence(255),
			ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceCore),
				ie.NewNetworkInstance("internet"),
				ie.NewUEIPAddress(2, "60.60.0.6", "", 0, 0),
			),
			ie.NewFARID(2),
			ie.NewQERID(1),
		),
		ie.NewCreateFAR(
			ie.NewFARID(1),
			ie.NewApplyAction(2), // 2: FORW(Forward)
			ie.NewForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceSGiLANN6LAN),
				ie.NewNetworkInstance("internet"),
			),
		),
		ie.NewCreateFAR(
			ie.NewFARID(2),
			ie.NewApplyAction(2), // 2: FORW(Forward)
			ie.NewForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				ie.NewNetworkInstance("internet"),
				ie.NewOuterHeaderCreation(
					256,          // desc
					2,            // teid
					"172.16.1.3", // v4
					"",           // v6
					0,            // port
					0,            // ctag
					0,            // stag
				),
			),
		),
		ie.NewCreateQER(
			ie.NewQERID(1),
			ie.NewGateStatus(ie.GateStatusOpen, ie.GateStatusOpen),
			ie.NewMBR(2000000, 1000000),
			ie.NewQFI(1),
		),
		ie.NewPDNType(ie.PDNTypeIPv4),
	).Marshal()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.Write(sereq); err != nil {
		log.Fatal(err)
	}
	log.Printf("sent Session Establishment Request to: %s", raddr)

	if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		log.Fatal(err)
	}

	waiting = true
	for waiting {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		msg, err := message.Parse(buf[:n])
		if err != nil {
			log.Printf("ignored undecodable message: %x, error: %s", buf[:n], err)
			continue
		}

		seres, ok := msg.(*message.SessionEstablishmentResponse)
		if !ok {
			log.Printf("got unexpected message: %s, from: %s", msg.MessageTypeName(), addr)
			continue
		}

		waiting = false
		if seres.Cause == nil {
			log.Printf("got non accepted response")
			return
		}
		if cause, err := seres.Cause.Cause(); cause != ie.CauseRequestAccepted || err != nil {
			log.Printf("got non accepted response")
			return
		}
	}
}
