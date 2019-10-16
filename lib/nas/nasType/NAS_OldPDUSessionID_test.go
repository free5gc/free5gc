//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewOldPDUSessionID(t *testing.T) {}

var nasTypeOldPDUSessionIDULNASTransportOldPDUSessionIDTypeTable = []NasTypeIeiData{
	{nasMessage.ULNASTransportOldPDUSessionIDType, nasMessage.ULNASTransportOldPDUSessionIDType},
}

func TestNasTypeOldPDUSessionIDGetSetIei(t *testing.T) {}

type nasTypeOldPDUSessionIDPduSessionIdentity2Value struct {
	in  uint8
	out uint8
}

var nasTypeOldPDUSessionIDPduSessionIdentity2ValueTable = []nasTypeOldPDUSessionIDPduSessionIdentity2Value{
	{0xff, 0xff},
}

func TestNasTypeOldPDUSessionIDGetSetOldPDUSessionID(t *testing.T) {}

type nasTypeOldPDUSessionID struct {
	inIei                       uint8
	inPduSessionIdentity2Value  uint8
	outIei                      uint8
	outPduSessionIdentity2Value uint8
}

var nasTypeOldPDUSessionIDTable = []nasTypeOldPDUSessionID{
	{nasMessage.ULNASTransportOldPDUSessionIDType, 0xff,
		nasMessage.ULNASTransportOldPDUSessionIDType, 0xff},
}

func TestNasTypeOldPDUSessionID(t *testing.T) {}
