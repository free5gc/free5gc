//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNewSpareHalfOctetAndDeregistrationType(t *testing.T) {}

type nasTypeDeregistrationTypeAndSpareHalfOctetSwitchOff struct {
	in  uint8
	out uint8
}

var nasTypeDeregistrationTypeAndSpareHalfOctetSwitchOffTable = []nasTypeDeregistrationTypeAndSpareHalfOctetSwitchOff{
	{0x01, 0x01},
}

func TestNasTypeDeregistrationTypeAndSpareHalfOctetGetSetSwitchOff(t *testing.T) {}

type nasTypeDeregistrationTypeAndSpareHalfOctetReRegistrationRequired struct {
	in  uint8
	out uint8
}

var nasTypeDeregistrationTypeAndSpareHalfOctetReRegistrationRequiredTable = []nasTypeDeregistrationTypeAndSpareHalfOctetReRegistrationRequired{
	{0x01, 0x01},
}

func TestNasTypeDeregistrationTypeAndSpareHalfOctetGetSetReRegistrationRequired(t *testing.T) {}

type nasTypeDeregistrationTypeAndSpareHalfOctetAccessType struct {
	in  uint8
	out uint8
}

var nasTypeDeregistrationTypeAndSpareHalfOctetAccessTypeTable = []nasTypeDeregistrationTypeAndSpareHalfOctetAccessType{
	{0x03, 0x3},
}

func TestNasTypeDeregistrationTypeAndSpareHalfOctetGetSetAccessType(t *testing.T) {}

type testDeregistrationTypeAndSpareHalfOctetDataTemplate struct {
	inSwitchOff              uint8
	inReRegistrationRequired uint8
	inAccessType             uint8
	in                       nasType.SpareHalfOctetAndDeregistrationType
	out                      nasType.SpareHalfOctetAndDeregistrationType
}

var deregistrationTypeAndSpareHalfOctetTestData = []nasType.SpareHalfOctetAndDeregistrationType{
	{0x0f},
}

var deregistrationTypeAndSpareHalfOctetExpectedData = []nasType.SpareHalfOctetAndDeregistrationType{
	{0xf},
}

var deregistrationTypeAndSpareHalfOctetTestTable = []testDeregistrationTypeAndSpareHalfOctetDataTemplate{
	{0x01, 0x01, 0x03, deregistrationTypeAndSpareHalfOctetTestData[0], deregistrationTypeAndSpareHalfOctetExpectedData[0]},
}

func TestNasTypeDeregistrationTypeAndSpareHalfOctet(t *testing.T) {}
