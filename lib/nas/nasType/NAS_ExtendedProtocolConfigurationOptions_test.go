//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewExtendedProtocolConfigurationOptions(t *testing.T) {}

var nasTypeRegistrationAcceptExtendedProtocolConfigurationOptionsIeiTable = []NasTypeIeiData{
	{0x7B, 0x7B},
}

func TestNasTypeExtendedProtocolConfigurationOptionsGetSetIei(t *testing.T) {}

var nasTypeExtendedProtocolConfigurationOptionsLenTable = []NasTypeLenUint16Data{
	{4, 4},
}

func TestNasTypeExtendedProtocolConfigurationOptionsGetSetLen(t *testing.T) {}

type nasTypetExtendedProtocolConfigurationOptionsExtendedProtocolConfigurationOptionsContents struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeExtendedProtocolConfigurationOptionsExtendedProtocolConfigurationOptionsContentsTable = []nasTypetExtendedProtocolConfigurationOptionsExtendedProtocolConfigurationOptionsContents{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeExtendedProtocolConfigurationOptionsGetSetEENL(t *testing.T) {}
