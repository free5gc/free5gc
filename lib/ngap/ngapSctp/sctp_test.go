//go:binary-only-package

package ngapSctp_test

import (
	"github.com/stretchr/testify/assert"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/ishidawataru/sctp"

	"free5gc/lib/ngap/logger"
	"free5gc/lib/ngap/ngapSctp"
)

var testClientNum = 2

func TestSCTP(t *testing.T) {}

func sendSCTPMsg(conn net.Conn, readChan chan ngapSctp.ConnData) {}

func testClient() {}
