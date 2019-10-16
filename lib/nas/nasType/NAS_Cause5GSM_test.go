//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewCause5GSM(t *testing.T) {}

var nasTypePDUSessionReleaseCompleteCause5GSMTable = []NasTypeIeiData{
	{nasMessage.PDUSessionReleaseCompleteCause5GSMType, nasMessage.PDUSessionReleaseCompleteCause5GSMType},
}

func TestNasTypeCause5GSMGetSetIei(t *testing.T) {}

type nasTypeCause5GSMCauseValueData struct {
	in  uint8
	out uint8
}

var nasTypeCause5GSMOctetTable = []nasTypeCause5GSMCauseValueData{
	{0xff, 0xff},
}

func TestNasTypeCause5GSMGetSetCauseValue(t *testing.T) {}

type testCause5GSMDataTemplate struct {
	in  nasType.Cause5GSM
	out nasType.Cause5GSM
}

var cause5GSMTestData = []nasType.Cause5GSM{
	{nasMessage.PDUSessionReleaseCompleteCause5GSMType, 0xff},
}

var cause5GSMExpectedTestData = []nasType.Cause5GSM{
	{nasMessage.PDUSessionReleaseCompleteCause5GSMType, 0xff},
}

var cause5GSMTestTable = []testCause5GSMDataTemplate{
	{cause5GSMTestData[0], cause5GSMExpectedTestData[0]},
}

func TestNasTypeCause5GSM(t *testing.T) {}
