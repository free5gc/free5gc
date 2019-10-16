//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewLADNInformation(t *testing.T) {}

var nasTypeRegistrationRequestLADNInformationTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandLADNInformationType, nasMessage.ConfigurationUpdateCommandLADNInformationType},
}

func TestNasTypeLADNInformationGetSetIei(t *testing.T) {}

var nasTypeLADNInformationLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeLADNInformationGetSetLen(t *testing.T) {}

type nasTypeLADNInformationLADNDData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeLADNInformationLADNDTable = []nasTypeLADNInformationLADNDData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeLADNInformationGetSetLADND(t *testing.T) {}

type testLADNInformationDataTemplate struct {
	inIei    uint8
	inLen    uint16
	inLADND  []uint8
	outIei   uint8
	outLen   uint16
	outLADND []uint8
}

var testLADNInformationTestTable = []testLADNInformationDataTemplate{
	{nasMessage.ConfigurationUpdateCommandLADNInformationType, 2, []uint8{0xff, 0xff},
		nasMessage.ConfigurationUpdateCommandLADNInformationType, 2, []uint8{0xff, 0xff}},
}

func TestNasTypeLADNInformation(t *testing.T) {}
