//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSSCMode(t *testing.T) {}

var nasTypeSSCModeIeiTable = []NasTypeIeiData{
	{0x01, 0x01},
}

func TestNasTypeSSCModeGetSetIei(t *testing.T) {}

type nasTypeRequestSpareType struct {
	in  uint8
	out uint8
}

var nasTypeSSCModeSpareTable = []nasTypeRequestSpareType{
	{0x01, 0x01},
}

func TestNasTypeSSCModeGetSetSpare(t *testing.T) {}

type nasTypeRequestSSCModeType struct {
	in  uint8
	out uint8
}

var nasTypeSSCModeSSCModeTable = []nasTypeRequestSSCModeType{
	{0x01, 0x01},
}

func TestNasTypeSSCModeGetSetSSCMode(t *testing.T) {}

type SSCModeTestDataTemplate struct {
	in  nasType.SSCMode
	out nasType.SSCMode
}

var SSCModeTestData = []nasType.SSCMode{
	{nasMessage.PDUSessionEstablishmentRequestSSCModeType},
}

var SSCModeExpectedTestData = []nasType.SSCMode{
	{0x19},
}

var SSCModeTable = []SSCModeTestDataTemplate{
	{SSCModeTestData[0], SSCModeExpectedTestData[0]},
}

func TestNasTypeSSCMode(t *testing.T) {}
