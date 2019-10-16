//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewCapability5GMM(t *testing.T) {}

var nasTypeRegistrationRequestCapability5GMMTypeTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestCapability5GMMType, nasMessage.RegistrationRequestCapability5GMMType},
}

func TestNasTypeCapability5GMMGetSetIei(t *testing.T) {}

var nasTypeCapability5GMMLenTable = []NasTypeLenuint8Data{
	{12, 12},
}

func TestNasTypeCapability5GMMGetSetLen(t *testing.T) {}

type nasTypeCapability5GMMLPPData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeCapability5GMMLPPTable = []nasTypeCapability5GMMLPPData{
	{12, 0x01, 0x01},
}

func TestNasTypeCapability5GMMGetSetLPP(t *testing.T) {}

type nasTypeCapability5GMMHOAttachData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeCapability5GMMHOAttachTable = []nasTypeCapability5GMMHOAttachData{
	{12, 0x01, 0x01},
}

func TestNasTypeCapability5GMMGetSetHOAttach(t *testing.T) {}

type nasTypeCapability5GMMS1ModeData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeCapability5GMMS1ModeTable = []nasTypeCapability5GMMS1ModeData{
	{12, 0x01, 0x01},
}

func TestNasTypeCapability5GMMGetSetS1Mode(t *testing.T) {}

type nasTypeCapability5GMMSpareData struct {
	inLen uint8
	in    [12]uint8
	out   [12]uint8
}

var nasTypeCapability5GMMSpareTable = []nasTypeCapability5GMMSpareData{
	{12, [12]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, [12]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

func TestNasTypeCapability5GMMGetSetSpare(t *testing.T) {}

type testCapability5GMMDataTemplate struct {
	inLPP      uint8
	inHOAttach uint8
	inS1Mode   uint8
	inSpare    [12]uint8
	in         nasType.Capability5GMM
	out        nasType.Capability5GMM
}

var capability5GMMTestData = []nasType.Capability5GMM{
	{nasMessage.RegistrationRequestCapability5GMMType, 13, [13]uint8{0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

var capability5GMMExpectedTestData = []nasType.Capability5GMM{
	{nasMessage.RegistrationRequestCapability5GMMType, 13, [13]uint8{0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

var capability5GMMTestTable = []testCapability5GMMDataTemplate{
	{0x01, 0x01, 0x01, [12]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, capability5GMMTestData[0], capability5GMMExpectedTestData[0]},
}

func TestNasTypeCapability5GMM(t *testing.T) {}
