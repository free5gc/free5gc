//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewEAPMessage(t *testing.T) {}

var nasTypeAuthenticationRequestEAPMessageIeiTable = []NasTypeIeiData{
	{0, 0},
}

func TestNasTypeEAPMessageGetSetIei(t *testing.T) {}

var nasTypeEAPMessageLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeEAPMessageGetSetLen(t *testing.T) {}

type nasTypetEAPMessageData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeEAPMessageTable = []nasTypetEAPMessageData{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeEAPMessageGetSetEAPMessage(t *testing.T) {}

type testEAPDataTemplate struct {
	in  nasType.EAPMessage
	out nasType.EAPMessage
}

var EAPMessageTestData = []nasType.EAPMessage{
	{0, 2, []byte{0x00, 0x00}}, //AuthenticationResult
}

var EAPMessageExpectedTestData = []nasType.EAPMessage{
	{0, 2, []byte{0x00, 0x00}}, //AuthenticationResult
}

var EAPMessageTestTable = []testEAPDataTemplate{
	{EAPMessageTestData[0], EAPMessageExpectedTestData[0]},
}

func TestNasTypeEAPMessage(t *testing.T) {}
