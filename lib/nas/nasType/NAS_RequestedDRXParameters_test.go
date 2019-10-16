//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRequestedDRXParameters(t *testing.T) {}

var nasTypeRequestedDRXParametersServiceRejectT3346ValueTypeTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestRequestedDRXParametersType, nasMessage.RegistrationRequestRequestedDRXParametersType},
}

func TestNasTypeRequestedDRXParametersGetSetIei(t *testing.T) {}

var nasTypeRequestedDRXParametersLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeRequestedDRXParametersGetSetLen(t *testing.T) {}

type nasTypeRequestedDRXParametersDRXValueData struct {
	in  uint8
	out uint8
}

var nasTypeRequestedDRXParametersDRXValueTable = []nasTypeRequestedDRXParametersDRXValueData{
	{0x0f, 0x0f},
}

func TestNasTypeRequestedDRXParametersGetSetGPRSTimer2Value(t *testing.T) {}

type testRequestedDRXParametersDataTemplate struct {
	inIei       uint8
	inLen       uint8
	inDRXValue  uint8
	outIei      uint8
	outLen      uint8
	outDRXValue uint8
}

var testRequestedDRXParametersTestTable = []testRequestedDRXParametersDataTemplate{
	{nasMessage.RegistrationRequestRequestedDRXParametersType, 2, 0x0f,
		nasMessage.RegistrationRequestRequestedDRXParametersType, 2, 0x0f},
}

func TestNasTypeRequestedDRXParameters(t *testing.T) {}
