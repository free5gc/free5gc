//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRequestedQosRules(t *testing.T) {}

var nasTypeAuthenticationResultRequestedQosRulesTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationRequestRequestedQosRulesType, nasMessage.PDUSessionModificationRequestRequestedQosRulesType},
}

func TestNasTypeRequestedQosRulesGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultRequestedQosRulesLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeRequestedQosRulesGetSetLen(t *testing.T) {}

type nasTypeRequestedQosRulesData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeRequestedQosRulesTable = []nasTypeRequestedQosRulesData{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeRequestedQosRulesGetSetContent(t *testing.T) {}

type testRequestedQosRulesDataTemplate struct {
	in  nasType.RequestedQosRules
	out nasType.RequestedQosRules
}

var RequestedQosRulesTestData = []nasType.RequestedQosRules{
	{nasMessage.PDUSessionModificationRequestRequestedQosRulesType, 2, []byte{0x01, 0x02}},
}

var RequestedQosRulesExpectedTestData = []nasType.RequestedQosRules{
	{nasMessage.PDUSessionModificationRequestRequestedQosRulesType, 2, []byte{0x01, 0x02}},
}

var RequestedQosRulesTestTable = []testRequestedQosRulesDataTemplate{
	{RequestedQosRulesTestData[0], RequestedQosRulesExpectedTestData[0]},
}

func TestNasTypeRequestedQosRules(t *testing.T) {}
