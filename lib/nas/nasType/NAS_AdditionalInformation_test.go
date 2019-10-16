//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAdditionalInformation(t *testing.T) {}

var nasTypeULNASTransportAdditionalInformationTable = []NasTypeIeiData{
	{nasMessage.ULNASTransportAdditionalInformationType, nasMessage.ULNASTransportAdditionalInformationType},
}

func TestNasTypeAdditionalInformationGetSetIei(t *testing.T) {}

var nasTypeULNASTransportAdditionalInformationTLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeAdditionalInformationGetSetLen(t *testing.T) {}

type AdditionalInformationValue struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeAdditionalInformationValueTable = []AdditionalInformationValue{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x01}},
}

func TestNasTypeAdditionalInformationGetSetAdditionalInformationValue(t *testing.T) {}

type testAdditionalInformationDataTemplate struct {
	in  nasType.AdditionalInformation
	out nasType.AdditionalInformation
}

var additionalInformationTestData = []nasType.AdditionalInformation{
	{nasMessage.ULNASTransportAdditionalInformationType, 2, []uint8{0x00, 0x01}},
}

var additionalInformationExpectedTestData = []nasType.AdditionalInformation{
	{nasMessage.ULNASTransportAdditionalInformationType, 2, []uint8{0x00, 0x01}},
}

var additionalInformationTable = []testAdditionalInformationDataTemplate{
	{additionalInformationTestData[0], additionalInformationExpectedTestData[0]},
}

func TestNasTypeAdditionalInformation(t *testing.T) {}
