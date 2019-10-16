//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthenticationParameterRAND(t *testing.T) {}

var nasTypeAuthenticationRequestAuthenticationParameterRANDTable = []NasTypeIeiData{
	{nasMessage.AuthenticationRequestAuthenticationParameterRANDType, nasMessage.AuthenticationRequestAuthenticationParameterRANDType},
}

func TestNasTypeAuthenticationParameterRANDGetSetIei(t *testing.T) {}

type nasTypeAuthenticationParameterRANDOctetData struct {
	in  [16]uint8
	out [16]uint8
}

var nasTypeAuthenticationParameterRANDOctetTable = []nasTypeAuthenticationParameterRANDOctetData{
	{[16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestNasTypeAuthenticationParameterRANDGetSetRANDValue(t *testing.T) {}

type testAuthenticationParameterRANDDataTemplate struct {
	in  nasType.AuthenticationParameterRAND
	out nasType.AuthenticationParameterRAND
}

var authenticationParameterRANDTestData = []nasType.AuthenticationParameterRAND{
	{nasMessage.AuthenticationRequestAuthenticationParameterRANDType, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationParameterRANDExpectedTestData = []nasType.AuthenticationParameterRAND{
	{nasMessage.AuthenticationRequestAuthenticationParameterRANDType, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationParameterRANDTestTable = []testAuthenticationParameterRANDDataTemplate{
	{authenticationParameterRANDTestData[0], authenticationParameterRANDExpectedTestData[0]},
}

func TestNasTypeAuthenticationParameterRAND(t *testing.T) {}
