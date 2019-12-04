//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRequestedQosFlowDescriptions(t *testing.T) {}

var nasTypeAuthenticationResultRequestedQosFlowDescriptionsTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationRequestRequestedQosFlowDescriptionsType, nasMessage.PDUSessionModificationRequestRequestedQosFlowDescriptionsType},
}

func TestNasTypeRequestedQosFlowDescriptionsGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultRequestedQosFlowDescriptionsLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeRequestedQosFlowDescriptionsGetSetLen(t *testing.T) {}

type nasTypeRequestedQosFlowDescriptionsData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeRequestedQosFlowDescriptionsTable = []nasTypeRequestedQosFlowDescriptionsData{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeRequestedQosFlowDescriptionsGetSetContent(t *testing.T) {}

type testRequestedQosFlowDescriptionsDataTemplate struct {
	in  nasType.RequestedQosFlowDescriptions
	out nasType.RequestedQosFlowDescriptions
}

var RequestedQosFlowDescriptionsTestData = []nasType.RequestedQosFlowDescriptions{
	{nasMessage.PDUSessionModificationRequestRequestedQosFlowDescriptionsType, 2, []byte{0x01, 0x02}},
}

var RequestedQosFlowDescriptionsExpectedTestData = []nasType.RequestedQosFlowDescriptions{
	{nasMessage.PDUSessionModificationRequestRequestedQosFlowDescriptionsType, 2, []byte{0x01, 0x02}},
}

var RequestedQosFlowDescriptionsTestTable = []testRequestedQosFlowDescriptionsDataTemplate{
	{RequestedQosFlowDescriptionsTestData[0], RequestedQosFlowDescriptionsExpectedTestData[0]},
}

func TestNasTypeRequestedQosFlowDescriptions(t *testing.T) {}
