//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRejectedNSSAI(t *testing.T) {}

var nasTypeAuthenticationResultRejectedNSSAITable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptRejectedNSSAIType, nasMessage.RegistrationAcceptRejectedNSSAIType},
}

func TestNasTypeRejectedNSSAIGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultRejectedNSSAILenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeRejectedNSSAIGetSetLen(t *testing.T) {}

type nasTypeRejectedNSSAIContentsData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeRejectedNSSAIContentsTable = []nasTypeRejectedNSSAIContentsData{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeRejectedNSSAIGetSetRejectedNSSAIContents(t *testing.T) {}

type testRejectedNSSAIDataTemplate struct {
	in  nasType.RejectedNSSAI
	out nasType.RejectedNSSAI
}

var RejectedNSSAITestData = []nasType.RejectedNSSAI{
	{nasMessage.RegistrationAcceptRejectedNSSAIType, 2, []byte{0x00, 0x00}},
}

var RejectedNSSAIExpectedTestData = []nasType.RejectedNSSAI{
	{nasMessage.RegistrationAcceptRejectedNSSAIType, 2, []byte{0x00, 0x00}},
}

var RejectedNSSAITestTable = []testRejectedNSSAIDataTemplate{
	{RejectedNSSAITestData[0], RejectedNSSAIExpectedTestData[0]},
}

func TestNasTypeRejectedNSSAI(t *testing.T) {}
