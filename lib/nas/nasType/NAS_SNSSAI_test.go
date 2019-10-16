//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSNSSAI(t *testing.T) {}

var nasTypeServiceRequestSNSSAITable = []NasTypeIeiData{
	{nasMessage.PDUSessionEstablishmentAcceptSNSSAIType, nasMessage.PDUSessionEstablishmentAcceptSNSSAIType},
}

func TestNasTypeSNSSAIGetSetIei(t *testing.T) {}

var nasTypeServiceRequestSNSSAILenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeSNSSAIGetSetLen(t *testing.T) {}

type nasTypeSNSSAISST struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeSNSSAISSTTable = []nasTypeSNSSAISST{
	{2, 0x01, 0x01},
}

func TestNasTypeSNSSAIGetSetSST(t *testing.T) {}

type nasTypeSNSSAISD struct {
	inLen uint8
	in    [3]uint8
	out   [3]uint8
}

var nasTypeSNSSAISDTable = []nasTypeSNSSAISD{
	{3, [3]uint8{0x01, 0x01, 0x01}, [3]uint8{0x01, 0x01, 0x01}},
}

func TestNasTypeSNSSAIGetSetSD(t *testing.T) {}

type nasTypeSNSSAIMappedHPLMNSST struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeSNSSAIMappedHPLMNSSTTable = []nasTypeSNSSAIMappedHPLMNSST{
	{2, 0x01, 0x01},
}

func TestNasTypeSNSSAIGetSetMappedHPLMNSST(t *testing.T) {}

type nasTypeSNSSAIMappedHPLMNSD struct {
	inLen uint8
	in    [3]uint8
	out   [3]uint8
}

var nasTypeSNSSAIMappedHPLMNSDTable = []nasTypeSNSSAIMappedHPLMNSD{
	{3, [3]uint8{0x01, 0x01, 0x01}, [3]uint8{0x01, 0x01, 0x01}},
}

func TestNasTypeSNSSAIGetSetMappedHPLMNSD(t *testing.T) {}

type testSNSSAIDataTemplate struct {
	in  nasType.SNSSAI
	out nasType.SNSSAI
}

var SNSSAITestData = []nasType.SNSSAI{
	{nasMessage.PDUSessionEstablishmentAcceptSNSSAIType, 8, [8]uint8{}},
}

var SNSSAIExpectedData = []nasType.SNSSAI{
	{nasMessage.PDUSessionEstablishmentAcceptSNSSAIType, 8, [8]uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
}

var SNSSAITable = []testSNSSAIDataTemplate{
	{SNSSAITestData[0], SNSSAIExpectedData[0]},
}

func TestNasTypeSNSSAI(t *testing.T) {}
