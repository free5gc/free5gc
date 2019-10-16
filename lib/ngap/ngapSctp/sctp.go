//go:binary-only-package

package ngapSctp

import (
	"flag"
	"net"
	"strings"

	"github.com/ishidawataru/sctp"

	"free5gc/lib/ngap/logger"
)

// ConnData structure that pair the connection and data together
type ConnData struct {
	remoteAddr string
	data       []byte
	err        error
}

var clientNum int

// GetError Return the Error of ConnData structure
func (cd ConnData) GetError() error {}

// GetRAddr Return the Remote Addr of ConnData structure
func (cd ConnData) GetRAddr() string {}

// GetData Return the Data of ConnData structure
func (cd ConnData) GetData() []byte {}

// Server - Init SCTP Server, Set initial value / resource
func Server(addrStr string) *sctp.SCTPListener {}

// Accept - Accepting SCTP socket
func Accept(sctpLn *sctp.SCTPListener) (conn net.Conn, err error) {}

// Start - Start SCTP read channel
func Start(conn net.Conn, readChan chan ConnData) {}

// SendMsg - used to send out message to SCTP connection
func SendMsg(conn net.Conn, msg []byte) {}

func closeConnection(conn net.Conn) {}

// Destroy - Destroy the SCTP Server Resource
func Destroy(ln *sctp.SCTPListener) error {}
