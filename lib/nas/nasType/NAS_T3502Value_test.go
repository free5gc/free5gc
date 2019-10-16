//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewT3502Value(t *testing.T) {}

var nasTypeT3502ValueIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationRejectT3502ValueType, nasMessage.RegistrationRejectT3502ValueType},
}

func TestNasTypeT3502ValueGetSetIei(t *testing.T) {}

var nasTypeT3502ValueLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeT3502ValueGetSetLen(t *testing.T) {}

type NasTypeT3502ValueGPRSTimer2ValueData struct {
	in  uint8
	out uint8
}

var nasTypeT3502ValueGPRSTimer2ValueTable = []NasTypeT3502ValueGPRSTimer2ValueData{
	{0x2, 0x2},
}

func TestNasTypeT3502ValueGetSetGPRSTimer2Value(t *testing.T) {}

type testT3502ValueDataTemplate struct {
	in  nasType.T3502Value
	out nasType.T3502Value
}

var T3502ValueTestData = []nasType.T3502Value{
	{nasMessage.RegistrationRejectT3502ValueType, 1, 0x05},
}
var T3502ValueExpectedData = []nasType.T3502Value{
	{nasMessage.RegistrationRejectT3502ValueType, 1, 0x05},
}

var T3502ValueDataTestTable = []testT3502ValueDataTemplate{
	{T3502ValueTestData[0], T3502ValueExpectedData[0]},
}

func TestNasTypeT3502Value(t *testing.T) {}
