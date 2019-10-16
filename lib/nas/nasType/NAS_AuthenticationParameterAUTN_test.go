//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthenticationParameterAUTN(t *testing.T) {}

var nasTypeAuthenticationResultAuthenticationParameterAUTNTable = []NasTypeIeiData{
	{nasMessage.AuthenticationRequestAuthenticationParameterAUTNType, nasMessage.AuthenticationRequestAuthenticationParameterAUTNType},
}

func TestNasTypeAuthenticationParameterAUTNGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultAuthenticationParameterAUTNLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeAuthenticationParameterAUTNGetSetLen(t *testing.T) {}

type nasTypeAuthenticationParameterAUTNOctetData struct {
	inLen uint8
	in    [16]uint8
	out   [16]uint8
}

var nasTypeAuthenticationParameterAUTNOctetTable = []nasTypeAuthenticationParameterAUTNOctetData{
	{16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestNasTypeAuthenticationParameterAUTNGetSetAUTN(t *testing.T) {}

type testAuthenticationParameterAUTNDataTemplate struct {
	in  nasType.AuthenticationParameterAUTN
	out nasType.AuthenticationParameterAUTN
}

var authenticationParameterAUTNTestData = []nasType.AuthenticationParameterAUTN{
	{nasMessage.AuthenticationRequestAuthenticationParameterAUTNType, 16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationParameterAUTNExpectedTestData = []nasType.AuthenticationParameterAUTN{
	{nasMessage.AuthenticationRequestAuthenticationParameterAUTNType, 16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
}

var authenticationParameterAUTNTestTable = []testAuthenticationParameterAUTNDataTemplate{
	{authenticationParameterAUTNTestData[0], authenticationParameterAUTNExpectedTestData[0]},
}

func TestNasTypeAuthenticationParameterAUTN(t *testing.T) {}
