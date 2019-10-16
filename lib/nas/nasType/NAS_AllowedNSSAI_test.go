//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAllowedNSSAI(t *testing.T) {}

var nasTypeConfigurationUpdateCommandConfiguredNSSAITable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType, nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType},
}

func TestNasTypeAllowedNSSAIGetSetIei(t *testing.T) {}

var nasTypeConfigurationUpdateCommandConfiguredNSSAILenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeAllowedNSSAIGetSetLen(t *testing.T) {}

type SNSSAIValue struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeSNSSAIValueTable = []SNSSAIValue{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x01}},
}

func TestNasTypeAllowedNSSAIGetSetSNSSAIValue(t *testing.T) {}

type testAllowedNSSAIDataTemplate struct {
	in  nasType.AllowedNSSAI
	out nasType.AllowedNSSAI
}

var AllowedNSSAITestData = []nasType.AllowedNSSAI{
	{nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType, 2, []uint8{0x00, 0x01}},
}

var AllowedNSSAIExpectedTestData = []nasType.AllowedNSSAI{
	{nasMessage.ConfigurationUpdateCommandConfiguredNSSAIType, 2, []uint8{0x00, 0x01}},
}

var AllowedNSSAITable = []testAllowedNSSAIDataTemplate{
	{AllowedNSSAITestData[0], AllowedNSSAIExpectedTestData[0]},
}

func TestNasTypeAllowedNSSAI(t *testing.T) {}
