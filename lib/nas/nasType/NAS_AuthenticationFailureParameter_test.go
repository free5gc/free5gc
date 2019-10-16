//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthenticationFailureParameter(t *testing.T) {}

var nasTypeAuthenticationResultAuthenticationFailureParameterTable = []NasTypeIeiData{
	{nasMessage.AuthenticationFailureAuthenticationFailureParameterType, nasMessage.AuthenticationFailureAuthenticationFailureParameterType},
}

func TestNasTypeAuthenticationFailureParameterGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultAuthenticationFailureParameterLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeAuthenticationFailureParameterGetSetLen(t *testing.T) {}

type nasTypeOctetData struct {
	inLen uint8
	in    [14]uint8
	out   [14]uint8
}

var nasTypeOctetTable = []nasTypeOctetData{
	{14, [14]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, [14]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestNasTypeAuthenticationFailureParameterGetSetOctet(t *testing.T) {}

type testAuthenticationFailureParameterDataTemplate struct {
	in  nasType.AuthenticationFailureParameter
	out nasType.AuthenticationFailureParameter
}

var authenticationFailureParameterTestData = []nasType.AuthenticationFailureParameter{
	{nasMessage.AuthenticationFailureAuthenticationFailureParameterType, 14, [14]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationFailureParameterExpectedTestData = []nasType.AuthenticationFailureParameter{
	{nasMessage.AuthenticationFailureAuthenticationFailureParameterType, 14, [14]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationFailureParameterTestTable = []testAuthenticationFailureParameterDataTemplate{
	{authenticationFailureParameterTestData[0], authenticationFailureParameterExpectedTestData[0]},
}

func TestNasTypeAuthenticationFailureParameter(t *testing.T) {}
