//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewExtendedEmergencyNumberList(t *testing.T) {}

var nasTypeRegistrationAcceptExtendedEmergencyNumberListIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptExtendedEmergencyNumberListType, nasMessage.RegistrationAcceptExtendedEmergencyNumberListType},
}

func TestNasTypeExtendedEmergencyNumberListGetSetIei(t *testing.T) {}

var nasTypeExtendedEmergencyNumberListLenTable = []NasTypeLenUint16Data{
	{4, 4},
}

func TestNasTypeExtendedEmergencyNumberListGetSetLen(t *testing.T) {}

type nasTypetExtendedEmergencyNumberListEENL struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeExtendedEmergencyNumberListEENLTable = []nasTypetExtendedEmergencyNumberListEENL{
	{2, 0x01, 0x01},
}

func TestNasTypeExtendedEmergencyNumberListGetSetEENL(t *testing.T) {}

type nasTypetExtendedEmergencyNumberListEmergencyInformation struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeExtendedEmergencyNumberListEmergencyInformationTable = []nasTypetExtendedEmergencyNumberListEmergencyInformation{
	{3, []uint8{0x00, 0x00, 0x01}, []uint8{0x00, 0x00, 0x01}},
}

func TestNasTypeExtendedEmergencyNumberListGetSetExtendedEmergencyNumberList(t *testing.T) {}
