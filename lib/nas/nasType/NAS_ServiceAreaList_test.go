//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewServiceAreaList(t *testing.T) {}

var nasTypeServiceAreaListTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptServiceAreaListType, nasMessage.RegistrationAcceptServiceAreaListType},
}

func TestNasTypeServiceAreaListGetSetIei(t *testing.T) {}

var nasTypeServiceAreaListLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeServiceAreaListGetSetLen(t *testing.T) {}

type nasTypeServiceAreaListPartialServiceAreaListData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeServiceAreaListPartialServiceAreaListTable = []nasTypeServiceAreaListPartialServiceAreaListData{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeServiceAreaListGetSetPartialServiceAreaList(t *testing.T) {}

type testServiceAreaListDataTemplate struct {
	in  nasType.ServiceAreaList
	out nasType.ServiceAreaList
}

var ServiceAreaListTestData = []nasType.ServiceAreaList{
	{nasMessage.RegistrationAcceptServiceAreaListType, 2, []uint8{}},
}

var ServiceAreaListExpectedTestData = []nasType.ServiceAreaList{
	{nasMessage.RegistrationAcceptServiceAreaListType, 2, []uint8{0x01, 0x01}},
}

var ServiceAreaListTestTable = []testServiceAreaListDataTemplate{
	{ServiceAreaListTestData[0], ServiceAreaListExpectedTestData[0]},
}

func TestNasTypeServiceAreaList(t *testing.T) {}
