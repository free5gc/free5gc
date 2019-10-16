//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var PDUSessionEstablishmentRequestPDUSessionTypeTypeIeiInput uint8 = 0x09

func TestNasTypeNewPDUSessionType(t *testing.T) {}

var nasTypePDUSessionEstablishmentRequestPDUSessionTypeTable = []NasTypeIeiData{
	{PDUSessionEstablishmentRequestPDUSessionTypeTypeIeiInput, 0x09},
}

func TestNasTypePDUSessionTypeGetSetIei(t *testing.T) {}

var nasTypePDUSessionEstablishmentRequestPDUSessionTypeSpareTable = []NasTypeLenuint8Data{
	{0x1, 0x1},
}

func TestNasTypePDUSessionTypeGetSetSpare(t *testing.T) {}

var nasTypePDUSessionTypeValue = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3, 0x3},
	{0x4, 0x4},
	{0x5, 0x5},
}

func TestNasTypePDUSessionTypeGetSetPDUSessionTypeValue(t *testing.T) {}

type testPDUSessionTypeDataTemplate struct {
	inPDUSessionTypeValue uint8
	in                    nasType.PDUSessionType
	out                   nasType.PDUSessionType
}

var pDUSessionTypeTestData = []nasType.PDUSessionType{
	{(nasMessage.PDUSessionEstablishmentRequestPDUSessionTypeType)},
}

var pDUSessionTypeExpectedData = []nasType.PDUSessionType{
	{(0x90 + 0x01)},
}

var pDUSessionTypeTestTable = []testPDUSessionTypeDataTemplate{
	{0x01, pDUSessionTypeTestData[0], pDUSessionTypeExpectedData[0]},
}

func TestNasTypePDUSessionType(t *testing.T) {}
