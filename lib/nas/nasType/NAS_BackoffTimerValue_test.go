//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewBackoffTimerValue(t *testing.T) {}

var nasTypeAuthenticationRequestBackoffTimerValueIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptT3512ValueType, nasMessage.RegistrationAcceptT3512ValueType},
}

func TestNasTypeBackoffTimerValueGetSetIei(t *testing.T) {}

var nasTypeBackoffTimerValueLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeBackoffTimerValueGetSetLen(t *testing.T) {}

type nasTypeBackoffTimerValueUintTimerValue struct {
	in  uint8
	out uint8
}

var nasTypeBackoffTimerValueUintTimerValueTable = []nasTypeBackoffTimerValueUintTimerValue{
	{0x07, 0x07},
}

func TestNasTypeBackoffTimerValueGetSetUintTimerValue(t *testing.T) {}

type nasTypeBackoffTimerValueTimerValue struct {
	in  uint8
	out uint8
}

var nasTypeBackoffTimerValueTimerValueTable = []nasTypeBackoffTimerValueTimerValue{
	{0x07, 0x07},
}

func TestNasTypeBackoffTimerValueGetSetTimerValue(t *testing.T) {}

type testBackoffTimerValueDataTemplate struct {
	inUnitTimerValue uint8
	inTimerValue     uint8
	in               nasType.BackoffTimerValue
	out              nasType.BackoffTimerValue
}

var BackoffTimerValueTestData = []nasType.BackoffTimerValue{
	{nasMessage.RegistrationAcceptT3512ValueType, 1, 0xff},
}
var BackoffTimerValueExpectedData = []nasType.BackoffTimerValue{
	{nasMessage.RegistrationAcceptT3512ValueType, 1, 0xff},
}

var BackoffTimerValueDataTestTable = []testBackoffTimerValueDataTemplate{
	{0x07, 0x1F, BackoffTimerValueTestData[0], BackoffTimerValueExpectedData[0]},
}

func TestNasTypeBackoffTimer(t *testing.T) {}
