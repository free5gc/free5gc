//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewUEStatus(t *testing.T) {}

var nasTypeUEStatusIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestUEStatusType, nasMessage.RegistrationRequestUEStatusType},
}

func TestNasTypeUEStatusGetSetIei(t *testing.T) {}

var nasTypeUEStatusLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeUEStatusGetSetLen(t *testing.T) {}

type NasTypeUEStatusN1ModeRegData struct {
	in  uint8
	out uint8
}

var nasTypeUEStatusN1ModeRegTable = []NasTypeUEStatusN1ModeRegData{
	{0x01, 0x01},
}

func TestNasTypeUEStatusGetSetN1ModeReg(t *testing.T) {}

type NasTypeUEStatusS1ModeRegData struct {
	in  uint8
	out uint8
}

var nasTypeUEStatusS1ModeRegTable = []NasTypeUEStatusS1ModeRegData{
	{0x01, 0x01},
}

func TestNasTypeUEStatusGetSetS1ModeReg(t *testing.T) {}

type testUEStatusDataTemplate struct {
	in  nasType.UEStatus
	out nasType.UEStatus
}

var UEStatusTestData = []nasType.UEStatus{
	{nasMessage.RegistrationRequestUEStatusType, 1, 0x05},
}
var UEStatusExpectedData = []nasType.UEStatus{
	{nasMessage.RegistrationRequestUEStatusType, 1, 0x03},
}

var UEStatusDataTestTable = []testUEStatusDataTemplate{
	{UEStatusTestData[0], UEStatusExpectedData[0]},
}

func TestNasTypeUEStatus(t *testing.T) {}
