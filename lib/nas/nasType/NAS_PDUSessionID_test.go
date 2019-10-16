//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSessionID(t *testing.T) {}

var nasTypePDUSessionIDULNASTransportOldPDUSessionIDTypeTable = []NasTypeIeiData{
	{nasMessage.ULNASTransportOldPDUSessionIDType, nasMessage.ULNASTransportOldPDUSessionIDType},
}

func TestNasTypePDUSessionIDGetSetIei(t *testing.T) {}

type nasTypePDUSessionIDPduSessionIdentity2ValueData struct {
	in  uint8
	out uint8
}

var nasTypePDUSessionIDPduSessionIdentity2ValueTable = []nasTypePDUSessionIDPduSessionIdentity2ValueData{
	{0xff, 0xff},
}

func TestNasTypePDUSessionIDGetSetPduSessionIdentity2Value(t *testing.T) {}

type testPDUSessionIDDataTemplate struct {
	inPduSessionIdentity2Value  uint8
	outPduSessionIdentity2Value uint8
}

var testPDUSessionIDTestTable = []testPDUSessionIDDataTemplate{
	{0x0f, 0x0f},
}

func TestNasTypePDUSessionID(t *testing.T) {}
