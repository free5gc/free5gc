//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewEmergencyNumberList(t *testing.T) {}

var nasTypeRegistrationAcceptEmergencyNumberListIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptEmergencyNumberListType, nasMessage.RegistrationAcceptEmergencyNumberListType},
}

func TestNasTypeEmergencyNumberListGetSetIei(t *testing.T) {}

var nasTypeEmergencyNumberListLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeEmergencyNumberListGetSetLen(t *testing.T) {}

type nasTypetEmergencyNumberListLengthof1EmergencyNumberInformation struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeEmergencyNumberListLengthof1EmergencyNumberInformationTable = []nasTypetEmergencyNumberListLengthof1EmergencyNumberInformation{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeEmergencyNumberListGetSetLengthof1EmergencyNumberInformation(t *testing.T) {}

type nasTypetEmergencyNumberListEmergencyServiceCategoryValue struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeEmergencyNumberListEmergencyServiceCategoryValueTable = []nasTypetEmergencyNumberListEmergencyServiceCategoryValue{
	{2, []uint8{0x01, 0x1f}, []uint8{0x01, 0x1f}},
}

func TestNasTypeEmergencyNumberListGetSetEmergencyServiceCategoryValue(t *testing.T) {}

type nasTypetEmergencyNumberListEmergencyInformation struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeEmergencyNumberListEmergencyInformationTable = []nasTypetEmergencyNumberListEmergencyInformation{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeEmergencyNumberListGetSetEmergencyInformation(t *testing.T) {}

type testEmergencyNumberListDataTemplate struct {
	in  nasType.EmergencyNumberList
	out nasType.EmergencyNumberList
}

var emergencyNumberListTestData = []nasType.EmergencyNumberList{
	{nasMessage.RegistrationAcceptEmergencyNumberListType, 3, []byte{0x02, 0x1f, 0x22}},
}

var emergencyNumberListExpectedData = []nasType.EmergencyNumberList{
	{nasMessage.RegistrationAcceptEmergencyNumberListType, 3, []byte{0x02, 0x1f, 0x22}},
}

var emergencyNumberListTestTable = []testEmergencyNumberListDataTemplate{
	{emergencyNumberListTestData[0], emergencyNumberListExpectedData[0]},
}

func TestNasTypeEmergencyNumberList(t *testing.T) {}
