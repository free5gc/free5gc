//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUAddress(t *testing.T) {}

var nasTypePDUAddressPDUSessionEstablishmentAcceptPDUAddressTypeTable = []NasTypeIeiData{
	{nasMessage.PDUSessionEstablishmentAcceptPDUAddressType, nasMessage.PDUSessionEstablishmentAcceptPDUAddressType},
}

func TestNasTypePDUAddressGetSetIei(t *testing.T) {}

var nasTypePDUAddressLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypePDUAddressGetSetLen(t *testing.T) {}

type nasTypePDUAddressPDUSessionTypeValueData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUAddressPDUSessionTypeValueTable = []nasTypePDUAddressPDUSessionTypeValueData{
	{2, 0xff, 0x07},
}

func TestNasTypePDUAddressGetSetPDUSessionTypeValue(t *testing.T) {}

type nasTypePDUAddressPDUAddressInformationData struct {
	inLen uint8
	in    [12]uint8
	out   [12]uint8
}

var nasTypePDUAddressPDUAddressInformationTable = []nasTypePDUAddressPDUAddressInformationData{
	{12, [12]uint8{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f}, [12]uint8{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f}},
}

func TestNasTypePDUAddressGetSetPDUAddressInformation(t *testing.T) {}

type testPDUAddressDataTemplate struct {
	inIei                    uint8
	inLen                    uint8
	inPDUSessionTypeValue    uint8
	inPDUAddressInformation  [12]uint8
	outIei                   uint8
	outLen                   uint8
	outPDUSessionTypeValue   uint8
	outPDUAddressInformation [12]uint8
}

var testPDUAddressTestTable = []testPDUAddressDataTemplate{
	{nasMessage.PDUSessionEstablishmentAcceptPDUAddressType, 12, 0x07, [12]uint8{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f},
		nasMessage.PDUSessionEstablishmentAcceptPDUAddressType, 12, 0x07, [12]uint8{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f}},
}

func TestNasTypePDUAddress(t *testing.T) {}
