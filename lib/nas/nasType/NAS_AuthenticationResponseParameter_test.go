//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthenticationResponseParameter(t *testing.T) {}

var nasTypeAuthenticationResponseAuthenticationResponseParameterTable = []NasTypeIeiData{
	{nasMessage.AuthenticationResponseAuthenticationResponseParameterType, nasMessage.AuthenticationResponseAuthenticationResponseParameterType},
}

func TestNasTypeAuthenticationResponseParameterGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResponseAuthenticationResponseParameterLenTable = []NasTypeLenuint8Data{
	{16, 16},
}

func TestNasTypeAuthenticationResponseParameterGetSetLen(t *testing.T) {}

type nasTypeAuthenticationResponseParameterOctetData struct {
	inLen uint8
	in    [16]uint8
	out   [16]uint8
}

var nasTypeAuthenticationResponseParameterOctetTable = []nasTypeAuthenticationResponseParameterOctetData{
	{16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestNasTypeAuthenticationResponseParameterGetSetRES(t *testing.T) {}

type testAuthenticationResponseParameterDataTemplate struct {
	in  nasType.AuthenticationResponseParameter
	out nasType.AuthenticationResponseParameter
}

var authenticationResponseParameterTestData = []nasType.AuthenticationResponseParameter{
	{nasMessage.AuthenticationResponseAuthenticationResponseParameterType, 16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationResponseParameterExpectedTestData = []nasType.AuthenticationResponseParameter{
	{nasMessage.AuthenticationResponseAuthenticationResponseParameterType, 16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationResponseParameterTestTable = []testAuthenticationResponseParameterDataTemplate{
	{authenticationResponseParameterTestData[0], authenticationResponseParameterExpectedTestData[0]},
}

func TestNasTypeAuthenticationResponseParameter(t *testing.T) {}
