//go:binary-only-package

package pfcp_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/pfcp"
	"testing"
)

const (
	NodeRelatedHeaderHex    = "2005000000000100"
	SessionRelatedHeaderHex = "233400000000000000000001000000F0"
)

var NodeRelatedHeader = pfcp.Header{
	Version:        1,
	S:              0,
	MessageType:    pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
	MessageLength:  0,
	SequenceNumber: 1,
}

var SessionRelatedHeader = pfcp.Header{
	Version:         1,
	MP:              1,
	S:               1,
	MessageType:     pfcp.PFCP_SESSION_MODIFICATION_REQUEST,
	MessageLength:   0,
	SEID:            1,
	SequenceNumber:  0,
	MessagePriority: 15,
}

func TestPFCPHeader_MarshalBinary(t *testing.T) {}

func TestPFCPHeader_UnmarshalBinary(t *testing.T) {}
