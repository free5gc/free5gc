package pfcp

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/free5gc/go-upf/internal/forwarder"
	"github.com/free5gc/go-upf/internal/report"
	logger_util "github.com/free5gc/util/logger"
)

type PfcpServerMock struct {
	PfcpServer
}

func (p *PfcpServerMock) GetRNodes() map[string]*RemoteNode {
	return p.rnodes
}

func (p *PfcpServerMock) AddRNode(rnodeid string, node *RemoteNode) {
	p.rnodes[rnodeid] = node
}

func TestStart(t *testing.T) {
}

func TestStop(t *testing.T) {
	s := &PfcpServer{
		log: logrus.WithField(logger_util.FieldControlPlaneNodeID, "127.0.0.1"),
	}

	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	if err != nil {
		t.Errorf("failed to resolve UDP address: %v", err)
		return
	}
	s.conn, err = net.ListenUDP("udp4", addr)
	if err != nil {
		t.Errorf("expected err to be nil, but got %v", err)
	}

	if s.conn == nil {
		t.Errorf("expected s.conn not to be nil")
		return
	}

	s.Stop()

	if !isConnClosed(s.conn) {
		t.Errorf("expected connection to be closed")
	}
}

func TestNewNode(t *testing.T) {
	s := &PfcpServer{
		log: logrus.WithField(logger_util.FieldControlPlaneNodeID, "127.0.0.1"),
	}

	id := "smf1"
	driver := forwarder.Empty{}
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8805")
	if err != nil {
		t.Errorf("failed to resolve UDP address: %v", err)
		return
	}

	newNode := s.NewNode(id, addr, driver)

	assert.NotNil(t, newNode)
	assert.Equal(t, id, newNode.ID)
}

func TestUpdateNodeID(t *testing.T) {
	s := &PfcpServerMock{
		PfcpServer: PfcpServer{
			log:    logrus.WithField(logger_util.FieldControlPlaneNodeID, "127.0.0.1"),
			rnodes: make(map[string]*RemoteNode),
		},
	}

	origNodeId := "127.0.0.1"
	node := s.NewNode(origNodeId, nil, nil)
	s.AddRNode(origNodeId, node)

	newNodeId := "192.168.56.101"
	s.UpdateNodeID(node, newNodeId)

	assert.Nil(t, s.GetRNodes()[origNodeId])
	assert.NotNil(t, s.GetRNodes()[newNodeId])
}

func TestNotifySessReport(t *testing.T) {
	s := &PfcpServer{
		srCh: make(chan report.SessReport),
	}

	reports := []report.Report{}

	sr := report.SessReport{
		SEID:    1,
		Reports: reports,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		receivedSr := <-s.srCh
		assert.EqualValues(t, sr, receivedSr)
	}()

	s.NotifySessReport(sr)

	wg.Wait()
}

func TestNotifyTransTimeout(t *testing.T) {
	s := &PfcpServer{
		trToCh: make(chan TransactionTimeout),
	}
	txId := "127.0.0.1-1"

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		receivedSeid := <-s.trToCh
		assert.EqualValues(t, TransactionTimeout{TrType: TX, TrID: txId}, receivedSeid)
	}()

	s.NotifyTransTimeout(TX, txId)

	wg.Wait()
}

func isConnClosed(conn *net.UDPConn) bool {
	oneByte := make([]byte, 1)
	err := conn.SetReadDeadline(time.Now())
	if err != nil {
		return true
	}
	_, err = conn.Read(oneByte)
	if err != nil {
		netErr, ok := err.(net.Error)
		if ok && netErr.Timeout() {
			// The read timed out, which means the connection is still open
			return false
		}
		// Any other error means the connection is closed
		return true
	}

	// If we were able to read a byte, the connection is definitely open
	return false
}
