//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRequestType(t *testing.T) {}

var nasTypeRequestTypeIeiTable = []NasTypeIeiData{
	{0x08, 0x08},
}

func TestNasTypeRequestTypeGetSetIei(t *testing.T) {}

type nasTypeRequestRequestTypeValueType struct {
	in  uint8
	out uint8
}

var nasTypeRequestTypeRequestTypeValueTable = []nasTypeRequestRequestTypeValueType{
	{0x03, 0x03},
}

func TestNasTypeRequestTypeGetSetRequestTypeValue(t *testing.T) {}

type RequestTypeTestDataTemplate struct {
	in  nasType.RequestType
	out nasType.RequestType
}

var RequestTypeTestData = []nasType.RequestType{
	{nasMessage.ULNASTransportRequestTypeType + 0x01},
}

var RequestTypeExpectedTestData = []nasType.RequestType{
	{0x81},
}

var RequestTypeTable = []RequestTypeTestDataTemplate{
	{RequestTypeTestData[0], RequestTypeExpectedTestData[0]},
}

func TestNasTypeRequestType(t *testing.T) {}
