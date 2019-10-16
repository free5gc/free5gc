//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewCause5GMM(t *testing.T) {}

var nasTypeDeregistrationRequestUETerminatedDeregistrationCause5GMMTable = []NasTypeIeiData{
	{nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType, nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType},
}

func TestNasTypeCause5GMMGetSetIei(t *testing.T) {}

type nasTypeCause5GMMCauseValueData struct {
	in  uint8
	out uint8
}

var nasTypeCause5GMMOctetTable = []nasTypeCause5GMMCauseValueData{
	{0xff, 0xff},
}

func TestNasTypeCause5GMMGetSetCauseValue(t *testing.T) {}

type testCause5GMMDataTemplate struct {
	in  nasType.Cause5GMM
	out nasType.Cause5GMM
}

var cause5GMMTestData = []nasType.Cause5GMM{
	{nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType, 0xff},
}

var cause5GMMExpectedTestData = []nasType.Cause5GMM{
	{nasMessage.DeregistrationRequestUETerminatedDeregistrationCause5GMMType, 0xff},
}

var cause5GMMTestTable = []testCause5GMMDataTemplate{
	{cause5GMMTestData[0], cause5GMMExpectedTestData[0]},
}

func TestNasTypeCause5GMM(t *testing.T) {}
