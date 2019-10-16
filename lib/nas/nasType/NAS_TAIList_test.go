//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewTAIList(t *testing.T) {}

var nasTypeTAIListTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptTAIListType, nasMessage.RegistrationAcceptTAIListType},
}

func TestNasTypeTAIListGetSetIei(t *testing.T) {}

var nasTypeTAIListLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeTAIListGetSetLen(t *testing.T) {}

type nasTypeTAIListPartialTrackingAreaIdentityListData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeTAIListPartialTrackingAreaIdentityListTable = []nasTypeTAIListPartialTrackingAreaIdentityListData{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeTAIListGetSetPartialTrackingAreaIdentityList(t *testing.T) {}

type testTAIListDataTemplate struct {
	in  nasType.TAIList
	out nasType.TAIList
}

var TAIListTestData = []nasType.TAIList{
	{nasMessage.RegistrationAcceptTAIListType, 2, []uint8{}},
}

var TAIListExpectedTestData = []nasType.TAIList{
	{nasMessage.RegistrationAcceptTAIListType, 2, []uint8{0x01, 0x01}},
}

var TAIListTestTable = []testTAIListDataTemplate{
	{TAIListTestData[0], TAIListExpectedTestData[0]},
}

func TestNasTypeTAIList(t *testing.T) {}
