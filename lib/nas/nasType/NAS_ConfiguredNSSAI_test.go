//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewConfiguredNSSAI(t *testing.T) {}

var nasTypeAuthenticationRequestConfiguredNSSAIIeiTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandAllowedNSSAIType, nasMessage.ConfigurationUpdateCommandAllowedNSSAIType},
}

func TestNasTypeConfiguredNSSAIGetSetIei(t *testing.T) {}

var nasTypeConfiguredNSSAILenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeConfiguredNSSAIGetSetLen(t *testing.T) {}

type nasTypetConfiguredNSSAISNSSAIValue struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeConfiguredNSSAISNSSAIValueTable = []nasTypetConfiguredNSSAISNSSAIValue{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeConfiguredNSSAIGetSetSNSSAIValue(t *testing.T) {}

type testConfiguredNSSAIDataTemplate struct {
	in  nasType.ConfiguredNSSAI
	out nasType.ConfiguredNSSAI
}

var configuredNSSAITestData = []nasType.ConfiguredNSSAI{
	{nasMessage.ConfigurationUpdateCommandAllowedNSSAIType, 2, []byte{0x00, 0x00}},
}

var configuredNSSAIExpectedData = []nasType.ConfiguredNSSAI{
	{nasMessage.ConfigurationUpdateCommandAllowedNSSAIType, 2, []byte{0x00, 0x00}},
}

var configuredNSSAITestTable = []testConfiguredNSSAIDataTemplate{
	{configuredNSSAITestData[0], configuredNSSAIExpectedData[0]},
}

func TestNasTypeConfiguredNSSAI(t *testing.T) {}
