//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pDUSessionEstablishmentRejectAllowedSSCModeIeiInput uint8 = 0xf

func TestNasTypeNewAllowedSSCMode(t *testing.T) {}

//var nasTypePDUSessionEstablishmentRejectAllowedSSCModeOut = (nasMessage.PDUSessionEstablishmentRejectAllowedSSCModeType & 15) << 4
var nasTypePDUSessionEstablishmentRejectAllowedSSCModeTable = []NasTypeIeiData{
	{pDUSessionEstablishmentRejectAllowedSSCModeIeiInput, pDUSessionEstablishmentRejectAllowedSSCModeIeiInput},
}

func TestNasTypeAllowedSSCModeGetSetIei(t *testing.T) {}

var AllowedSSCModeSSC1Table = []NasTypeLenuint8Data{
	{0x01, 0x01},
}

func TestNasTypeAllowedSSCModeGetSetSSC1(t *testing.T) {}

var AllowedSSCModeSSC2Table = []NasTypeLenuint8Data{
	{0x01, 0x01},
}

func TestNasTypeAllowedSSCModeGetSetSSC2(t *testing.T) {}

var AllowedSSCModeSSC3Table = []NasTypeLenuint8Data{
	{0x01, 0x01},
}

func TestNasTypeAllowedSSCModeGetSetSSC3(t *testing.T) {}

type testAllowedSSCModeDataTemplate struct {
	in  nasType.AllowedSSCMode
	out nasType.AllowedSSCMode
}

var allowedSSCModeTestData = []nasType.AllowedSSCMode{
	{0xF0 + 0x07},
}

var allowedSSCModeExpectedTestData = []nasType.AllowedSSCMode{
	{0xF0 + 0x07},
}

var allowedSSCModeTestTable = []testAllowedSSCModeDataTemplate{
	{allowedSSCModeTestData[0], allowedSSCModeExpectedTestData[0]},
}

func TestNasTypeAllowedSSCMode(t *testing.T) {}
