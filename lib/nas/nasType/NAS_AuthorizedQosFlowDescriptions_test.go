//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthorizedQosFlowDescriptions(t *testing.T) {}

var nasTypePDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsTable = []NasTypeIeiData{
	{nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType, nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType},
}

func TestNasTypeAuthorizedQosFlowDescriptionsGetSetIei(t *testing.T) {}

var nasTypePDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsLenTable = []NasTypeLenUint16Data{
	{12, 12},
}

func TestNasTypeAuthorizedQosFlowDescriptionsGetSetLen(t *testing.T) {}

type nasTypeQoSFlowDescription struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeQoSFlowDescriptionTable = []nasTypeQoSFlowDescription{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x01}},
}

func TestNasTypeAuthorizedQosFlowDescriptionsGetSetQoSFlowDescription(t *testing.T) {}

type testAuthorizedQosFlowDescriptionsDataTemplate struct {
	in  nasType.AuthorizedQosFlowDescriptions
	out nasType.AuthorizedQosFlowDescriptions
}

var AuthorizedQosFlowDescriptionsTestData = []nasType.AuthorizedQosFlowDescriptions{
	{nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType, 2, []uint8{0x00, 0x01}},
}

var AuthorizedQosFlowDescriptionsExpectedTestData = []nasType.AuthorizedQosFlowDescriptions{
	{nasMessage.PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType, 2, []uint8{0x00, 0x01}},
}

var AuthorizedQosFlowDescriptionsTable = []testAuthorizedQosFlowDescriptionsDataTemplate{
	{AuthorizedQosFlowDescriptionsTestData[0], AuthorizedQosFlowDescriptionsExpectedTestData[0]},
}

func TestNasTypeAuthorizedQosFlowDescriptions(t *testing.T) {}
