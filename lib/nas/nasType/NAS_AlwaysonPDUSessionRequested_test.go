//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var AlwaysonPDUSessionRequestedIeiInput uint8 = 0x0B

func TestNasTypeNewAlwaysonPDUSessionRequested(t *testing.T) {}

var nasTypePDUSessionEstablishmentRequestAlwaysonPDUSessionRequestedTable = []NasTypeIeiData{
	{AlwaysonPDUSessionRequestedIeiInput, 0x0B},
}

func TestNasTypeAlwaysonPDUSessionRequestedGetSetIei(t *testing.T) {}

type nasTypeAlwaysonPDUSessionRequestedAPSI struct {
	in  uint8
	out uint8
}

var nasTypeAlwaysonPDUSessionRequestedAPSRTable = []nasTypeAlwaysonPDUSessionRequestedAPSI{
	{0x01, 0x01},
}

func TestNasTypeAlwaysonPDUSessionRequestedGetSetAPSR(t *testing.T) {}

type testAlwaysonPDUSessionRequestedDataTemplate struct {
	in  nasType.AlwaysonPDUSessionRequested
	out nasType.AlwaysonPDUSessionRequested
}

var alwaysonPDUSessionRequestedTestData = []nasType.AlwaysonPDUSessionRequested{
	{(0xB0 + 0x01)},
}

var alwaysonPDUSessionRequestedExpectedTestData = []nasType.AlwaysonPDUSessionRequested{
	{(0xB0 + 0x01)},
}

var alwaysonPDUSessionRequestedTestTable = []testAlwaysonPDUSessionRequestedDataTemplate{
	{alwaysonPDUSessionRequestedTestData[0], alwaysonPDUSessionRequestedExpectedTestData[0]},
}

func TestNasTypeAlwaysonPDUSessionRequested(t *testing.T) {}
