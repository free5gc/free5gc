//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNegotiatedDRXParameters(t *testing.T) {}

var nasTypeNegotiatedDRXParametersRegistrationRequestAdditionalGUTITable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptNegotiatedDRXParametersType, nasMessage.RegistrationAcceptNegotiatedDRXParametersType},
}

func TestNasTypeNegotiatedDRXParametersGetSetIei(t *testing.T) {}

var nasTypeNegotiatedDRXParametersLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeNegotiatedDRXParametersGetSetLen(t *testing.T) {}

type nasTypeNegotiatedDRXParametersDRXValueData struct {
	in  uint8
	out uint8
}

var nasTypeNegotiatedDRXParametersDRXValueTable = []nasTypeNegotiatedDRXParametersDRXValueData{
	{0x0f, 0x0f},
}

func TestNasTypeNegotiatedDRXParametersGetSetDRXValue(t *testing.T) {}

type testNegotiatedDRXParametersDataTemplate struct {
	inIei       uint8
	inLen       uint8
	inDRXValue  uint8
	outIei      uint8
	outLen      uint8
	outDRXValue uint8
}

var testNegotiatedDRXParametersTestTable = []testNegotiatedDRXParametersDataTemplate{
	{nasMessage.RegistrationAcceptNegotiatedDRXParametersType, 2, 0x0f,
		nasMessage.RegistrationAcceptNegotiatedDRXParametersType, 2, 0x0f},
}

func TestNasTypeNegotiatedDRXParameters(t *testing.T) {}
