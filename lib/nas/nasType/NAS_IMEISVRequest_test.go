//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var SecurityModeCommandIMEISVRequestTypeIeiInput uint8 = 0x0E

func TestNasTypeNewIMEISVRequest(t *testing.T) {}

var nasTypePDUSessionEstablishmentRequestIMEISVRequestTable = []NasTypeIeiData{
	{SecurityModeCommandIMEISVRequestTypeIeiInput, 0x0E},
}

func TestNasTypeIMEISVRequestGetSetIei(t *testing.T) {}

type nasTypeIMEISVRequestIMEISVRequestValue struct {
	in  uint8
	out uint8
}

var nasTypeIMEISVRequestIMEISVRequestValueTable = []nasTypeIMEISVRequestIMEISVRequestValue{
	{0x07, 0x07},
}

func TestNasTypeIMEISVRequestGetSetIMEISVRequestValue(t *testing.T) {}

type testIMEISVRequestDataTemplate struct {
	inIei                uint8
	inIMEISVRequestValue uint8

	outIei                uint8
	outIMEISVRequestValue uint8
}

var iMEISVRequestTestTable = []testIMEISVRequestDataTemplate{
	{SecurityModeCommandIMEISVRequestTypeIeiInput, 0x07,
		0x0E, 0x07},
}

func TestNasTypeIMEISVRequest(t *testing.T) {}
