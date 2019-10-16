//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRQTimerValue(t *testing.T) {}

var nasTypePDUSessionReleaseCompleteRQTimerValueTable = []NasTypeIeiData{
	{nasMessage.PDUSessionEstablishmentAcceptRQTimerValueType, nasMessage.PDUSessionEstablishmentAcceptRQTimerValueType},
}

func TestNasTypeRQTimerValueGetSetIei(t *testing.T) {}

var nasTypeRQTimerValueUintTable = []NasTypeLenuint8Data{
	{0x01, 0x01},
}

func TestNasTypeRQTimerValueGetSetUint(t *testing.T) {}

type nasTypeRQTimerValueTimerValueData struct {
	in  uint8
	out uint8
}

var nasTypeRQTimerValueTimerValueTable = []nasTypeRQTimerValueTimerValueData{
	{0x01, 0x01},
}

func TestNasTypeRQTimerValueGetSetTimerValue(t *testing.T) {}

type testRQTimerValueDataTemplate struct {
	inUnit       uint8
	inTimerValue uint8
	in           nasType.RQTimerValue
	out          nasType.RQTimerValue
}

var rQTimerValueTestData = []nasType.RQTimerValue{
	{nasMessage.PDUSessionEstablishmentAcceptRQTimerValueType, 0x01},
}

var rQTimerValueExpectedTestData = []nasType.RQTimerValue{
	{nasMessage.PDUSessionEstablishmentAcceptRQTimerValueType, 0x21},
}

var rQTimerValueTestTable = []testRQTimerValueDataTemplate{
	{0x01, 0x01, rQTimerValueTestData[0], rQTimerValueExpectedTestData[0]},
}

func TestNasTypeRQTimerValue(t *testing.T) {}
