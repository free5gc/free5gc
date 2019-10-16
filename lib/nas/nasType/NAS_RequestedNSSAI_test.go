//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRequestedNSSAI(t *testing.T) {}

var nasTypeAuthenticationResultRequestedNSSAITable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestRequestedNSSAIType, nasMessage.RegistrationRequestRequestedNSSAIType},
}

func TestNasTypeRequestedNSSAIGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultRequestedNSSAILenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeRequestedNSSAIGetSetLen(t *testing.T) {}

type nasTypeRequestedNSSAIData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeRequestedNSSAITable = []nasTypeRequestedNSSAIData{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeRequestedNSSAIGetSetContent(t *testing.T) {}

type testRequestedNSSAIDataTemplate struct {
	in  nasType.RequestedNSSAI
	out nasType.RequestedNSSAI
}

var RequestedNSSAITestData = []nasType.RequestedNSSAI{
	{nasMessage.RegistrationRequestRequestedNSSAIType, 2, []byte{0x01, 0x02}},
}

var RequestedNSSAIExpectedTestData = []nasType.RequestedNSSAI{
	{nasMessage.RegistrationRequestRequestedNSSAIType, 2, []byte{0x01, 0x02}},
}

var RequestedNSSAITestTable = []testRequestedNSSAIDataTemplate{
	{RequestedNSSAITestData[0], RequestedNSSAIExpectedTestData[0]},
}

func TestNasTypeRequestedNSSAI(t *testing.T) {}
