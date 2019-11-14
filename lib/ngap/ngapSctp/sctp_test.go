//go:binary-only-package

package ngapSctp_test

import (
	"encoding/binary"
	"net"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ishidawataru/sctp"

	"free5gc/lib/ngap/logger"
	"free5gc/lib/ngap/ngapSctp"
)

var testClientNum = 2

func TestSCTP(t *testing.T) {}

func sendSCTPMsg(conn net.Conn, readChan chan ngapSctp.ConnData) {}

func testClient() {}
