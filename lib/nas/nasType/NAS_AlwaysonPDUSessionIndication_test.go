//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var PDUSessionEstablishmentAcceptAlwaysonPDUSessionIndicationTypeIeiInput uint8 = 0x08

func TestNasTypeNewAlwaysonPDUSessionIndication(t *testing.T) {}

var nasTypePDUSessionEstablishmentRequestAlwaysonPDUSessionIndicationTable = []NasTypeIeiData{
	{PDUSessionEstablishmentAcceptAlwaysonPDUSessionIndicationTypeIeiInput, 0x08},
}

func TestNasTypeAlwaysonPDUSessionIndicationGetSetIei(t *testing.T) {}

type nasTypeAlwaysonPDUSessionIndicationAPSI struct {
	in  uint8
	out uint8
}

var nasTypeAlwaysonPDUSessionIndicationAPSITable = []nasTypeAlwaysonPDUSessionIndicationAPSI{
	{0x01, 0x01},
}

func TestNasTypeAlwaysonPDUSessionIndicationGetSetAPSI(t *testing.T) {}

type testAlwaysonPDUSessionIndicationDataTemplate struct {
	in  nasType.AlwaysonPDUSessionIndication
	out nasType.AlwaysonPDUSessionIndication
}

var alwaysonPDUSessionIndicationTestData = []nasType.AlwaysonPDUSessionIndication{
	{(0x80 + 0x01)},
}

var alwaysonPDUSessionIndicationExpectedTestData = []nasType.AlwaysonPDUSessionIndication{
	{(0x80 + 0x01)},
}

var alwaysonPDUSessionIndicationTestTable = []testAlwaysonPDUSessionIndicationDataTemplate{
	{alwaysonPDUSessionIndicationTestData[0], alwaysonPDUSessionIndicationExpectedTestData[0]},
}

func TestNasTypeAlwaysonPDUSessionIndication(t *testing.T) {}
