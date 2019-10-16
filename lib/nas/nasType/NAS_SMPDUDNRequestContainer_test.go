//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSMPDUDNRequestContainer(t *testing.T) {}

var nasTypeSMPDUDNRequestContainerTable = []NasTypeIeiData{
	{nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType, nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType},
}

func TestNasTypeSMPDUDNRequestContainerGetSetIei(t *testing.T) {}

var nasTypeSMPDUDNRequestContainerLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeSMPDUDNRequestContainerGetSetLen(t *testing.T) {}

type nasTypeSMPDUDNRequestContainerDNSpecificIdentityData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeSMPDUDNRequestContainerDNSpecificIdentityTable = []nasTypeSMPDUDNRequestContainerDNSpecificIdentityData{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeSMPDUDNRequestContainerGetSetDNSpecificIdentity(t *testing.T) {}

type testSMPDUDNRequestContainerDataTemplate struct {
	in  nasType.SMPDUDNRequestContainer
	out nasType.SMPDUDNRequestContainer
}

var SMPDUDNRequestContainerTestData = []nasType.SMPDUDNRequestContainer{
	{nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType, 2, []uint8{}},
}

var SMPDUDNRequestContainerExpectedTestData = []nasType.SMPDUDNRequestContainer{
	{nasMessage.PDUSessionEstablishmentRequestSMPDUDNRequestContainerType, 2, []uint8{0x01, 0x01}},
}

var SMPDUDNRequestContainerTestTable = []testSMPDUDNRequestContainerDataTemplate{
	{SMPDUDNRequestContainerTestData[0], SMPDUDNRequestContainerExpectedTestData[0]},
}

func TestNasTypeSMPDUDNRequestContainer(t *testing.T) {}
