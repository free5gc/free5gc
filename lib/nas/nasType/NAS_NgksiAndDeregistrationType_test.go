//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNgksiAndDeregistrationType(t *testing.T) {}

type nasTypeNgksiAndDeregistrationTypeTSC struct {
	in  uint8
	out uint8
}

var nasTypeNgksiAndDeregistrationTypeTSCTable = []nasTypeNgksiAndDeregistrationTypeTSC{
	{0x01, 0x01},
}

func TestNasTypeNgksiAndDeregistrationTypeGetSetTSC(t *testing.T) {}

type nasTypeNgksiAndDeregistrationTypeNasKeySetIdentifiler struct {
	in  uint8
	out uint8
}

var nasTypeNgksiAndDeregistrationTypeNasKeySetIdentifilerTable = []nasTypeNgksiAndDeregistrationTypeNasKeySetIdentifiler{
	{0x07, 0x07},
}

func TestNasTypeNgksiAndDeregistrationTypeGetSetNasKeySetIdentifiler(t *testing.T) {}

type nasTypeNgksiAndDeregistrationTypeSwitchOff struct {
	in  uint8
	out uint8
}

var nasTypeNgksiAndDeregistrationTypeSwitchOffTable = []nasTypeNgksiAndDeregistrationTypeSwitchOff{
	{0x01, 0x01},
}

func TestNasTypeNgksiAndDeregistrationTypeGetSetSwitchOff(t *testing.T) {}

type nasTypeNgksiAndDeregistrationTypeReRegistrationRequired struct {
	in  uint8
	out uint8
}

var nasTypeNgksiAndDeregistrationTypeReRegistrationRequiredTable = []nasTypeNgksiAndDeregistrationTypeReRegistrationRequired{
	{0x01, 0x01},
}

func TestNasTypeNgksiAndDeregistrationTypeGetSetReRegistrationRequired(t *testing.T) {}

type nasTypeNgksiAndDeregistrationTypeAccessType struct {
	in  uint8
	out uint8
}

var nasTypeNgksiAndDeregistrationTypeAccessTypeTable = []nasTypeNgksiAndDeregistrationTypeAccessType{
	{0x03, 0x03},
}

func TestNasTypeNgksiAndDeregistrationTypeGetSetAccessType(t *testing.T) {}

type testNgksiAndDeregistrationTypeDataTemplate struct {
	inTSC                     uint8
	inNasKeySetIdentifiler    uint8
	inSwitchOff               uint8
	inReRegistrationRequired  uint8
	inAccessType              uint8
	outTSC                    uint8
	outNasKeySetIdentifiler   uint8
	outSwitchOff              uint8
	outReRegistrationRequired uint8
	outAccessType             uint8
}

var NgksiAndDeregistrationTypeTestTable = []testNgksiAndDeregistrationTypeDataTemplate{
	{0x01, 0x07, 0x01, 0x01, 0x03,
		0x01, 0x07, 0x01, 0x01, 0x03},
}

func TestNasTypeNgksiAndDeregistrationType(t *testing.T) {}
