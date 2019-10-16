//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPduSessionID2Value(t *testing.T) {}

var nasTypePDUSessionIDULNASTransportPduSessionID2ValueTypeTypeTable = []NasTypeIeiData{
	{nasMessage.ULNASTransportPduSessionID2ValueType, nasMessage.ULNASTransportPduSessionID2ValueType},
}

func TestNasTypePduSessionID2ValueGetSetIei(t *testing.T) {}

type nasTypePDUSessionIDPduSessionID2ValueData struct {
	in  uint8
	out uint8
}

var nasTypePduSessionIdentity2ValueTable = []nasTypePDUSessionIDPduSessionID2ValueData{
	{0xff, 0xff},
}

func TestNasTypeGetSetPduSessionIdentity2Value(t *testing.T) {}

type testPduSessionIdentity2ValueDataTemplate struct {
	inIei                       uint8
	inPduSessionIdentity2Value  uint8
	outIei                      uint8
	outPduSessionIdentity2Value uint8
}

var testPduSessionIdentity2ValueTestTable = []testPduSessionIdentity2ValueDataTemplate{
	{nasMessage.ULNASTransportPduSessionID2ValueType, 0x0f,
		nasMessage.ULNASTransportPduSessionID2ValueType, 0x0f},
}

func TestNasTypePDUSessionID2Value(t *testing.T) {}
