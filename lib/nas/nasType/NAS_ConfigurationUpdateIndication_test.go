//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ConfigurationUpdateCommandConfigurationUpdateIndicationTypeIeiInput uint8 = 0x0D

func TestNasTypeNewConfigurationUpdateIndication(t *testing.T) {}

var nasTypePDUSessionEstablishmentRequestConfigurationUpdateIndicationTable = []NasTypeIeiData{
	{ConfigurationUpdateCommandConfigurationUpdateIndicationTypeIeiInput, 0x0D},
}

func TestNasTypeConfigurationUpdateIndicationGetSetIei(t *testing.T) {}

type nasTypeConfigurationUpdateIndicationRED struct {
	in  uint8
	out uint8
}

var nasTypeConfigurationUpdateIndicationREDTable = []nasTypeConfigurationUpdateIndicationRED{
	{0x01, 0x01},
}

func TestNasTypeConfigurationUpdateIndicationGetSetRED(t *testing.T) {}

type nasTypeConfigurationUpdateIndicationACK struct {
	in  uint8
	out uint8
}

var nasTypeConfigurationUpdateIndicationACKTable = []nasTypeConfigurationUpdateIndicationACK{
	{0x01, 0x01},
}

func TestNasTypeConfigurationUpdateIndicationGetSetACK(t *testing.T) {}

type testConfigurationUpdateIndicationDataTemplate struct {
	inRED uint8
	inACK uint8
	in    nasType.ConfigurationUpdateIndication
	out   nasType.ConfigurationUpdateIndication
}

var configurationUpdateIndicationTestData = []nasType.ConfigurationUpdateIndication{
	{(0xD0 + 0x03)},
}

var configurationUpdateIndicationExpectedData = []nasType.ConfigurationUpdateIndication{
	{(0xD0 + 0x03)},
}

var configurationUpdateIndicationTestTable = []testConfigurationUpdateIndicationDataTemplate{
	{0x01, 0x01, configurationUpdateIndicationTestData[0], configurationUpdateIndicationExpectedData[0]},
}

func TestNasTypeConfigurationUpdateIndication(t *testing.T) {}
