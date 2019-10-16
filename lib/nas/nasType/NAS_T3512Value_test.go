//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewT3512Value(t *testing.T) {}

var nasTypeT3512ValueIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptT3512ValueType, nasMessage.RegistrationAcceptT3512ValueType},
}

func TestNasTypeT3512ValueGetSetIei(t *testing.T) {}

var nasTypeT3512ValueLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeT3512ValueGetSetLen(t *testing.T) {}

type NasTypeT3512ValueUintValueData struct {
	in  uint8
	out uint8
}

var nasTypeT3512ValueUintValueTable = []NasTypeT3512ValueUintValueData{
	{0x1, 0x01},
}

func TestNasTypeT3512ValueGetSetUintValue(t *testing.T) {}

type NasTypeT3512ValueTimerValueData struct {
	in  uint8
	out uint8
}

var nasTypeT3512ValueTimerValueTable = []NasTypeT3512ValueTimerValueData{
	{0x1, 0x01},
}

func TestNasTypeT3512ValueGetSetTimerValue(t *testing.T) {}

type testT3512ValueDataTemplate struct {
	in  nasType.T3512Value
	out nasType.T3512Value
}

var T3512ValueTestData = []nasType.T3512Value{
	{nasMessage.RegistrationAcceptT3512ValueType, 1, 0x00},
}
var T3512ValueExpectedData = []nasType.T3512Value{
	{nasMessage.RegistrationAcceptT3512ValueType, 1, 0x21},
}

var T3512ValueDataTestTable = []testT3512ValueDataTemplate{
	{T3512ValueTestData[0], T3512ValueExpectedData[0]},
}

func TestNasTypeT3512Value(t *testing.T) {}
