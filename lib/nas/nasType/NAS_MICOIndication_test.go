//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ConfigurationUpdateCommandMICOIndicationTypeIeiInput uint8 = 0x0B

func TestNasTypeNewMICOIndication(t *testing.T) {}

var nasTypeConfigurationUpdateCommandMICOIndicationTable = []NasTypeIeiData{
	{ConfigurationUpdateCommandMICOIndicationTypeIeiInput, 0x0B},
}

func TestNasTypeMICOIndicationGetSetIei(t *testing.T) {}

type nasTypeMICOIndicationRAAI struct {
	in  uint8
	out uint8
}

var nasTypeMICOIndicationRAAITable = []nasTypeMICOIndicationRAAI{
	{0x01, 0x01},
}

func TestNasTypeMICOIndicationGetSetRAAI(t *testing.T) {}

type testMICOIndicationDataTemplate struct {
	inRAAI uint8
	in     nasType.MICOIndication
	out    nasType.MICOIndication
}

var mICOIndicationTestData = []nasType.MICOIndication{
	{(0xB0 + 0x01)},
}

var mICOIndicationExpectedData = []nasType.MICOIndication{
	{(0xB0 + 0x01)},
}

var mICOIndicationTestTable = []testMICOIndicationDataTemplate{
	{0x01, mICOIndicationTestData[0], mICOIndicationExpectedData[0]},
}

func TestNasTypeMICOIndication(t *testing.T) {}
