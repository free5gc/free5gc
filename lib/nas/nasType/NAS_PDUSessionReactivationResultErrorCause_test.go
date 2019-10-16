//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSessionReactivationResultErrorCause(t *testing.T) {}

var nasTypeRegistrationAcceptPDUSessionReactivationResultErrorCauseTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType, nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType},
}

func TestNasTypePDUSessionReactivationResultErrorCauseGetSetIei(t *testing.T) {}

var nasTypeRegistrationAcceptPDUSessionReactivationResultErrorCauseLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypePDUSessionReactivationResultErrorCauseGetSetLen(t *testing.T) {}

type PDUSessionIDAndCauseValue struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypePDUSessionIDAndCauseValueTable = []PDUSessionIDAndCauseValue{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x01}},
}

func TestNasTypePDUSessionReactivationResultErrorCauseGetSetPDUSessionIDAndCauseValue(t *testing.T) {}

type testPDUSessionReactivationResultErrorCauseDataTemplate struct {
	in  nasType.PDUSessionReactivationResultErrorCause
	out nasType.PDUSessionReactivationResultErrorCause
}

var pDUSessionReactivationResultErrorCauseTestData = []nasType.PDUSessionReactivationResultErrorCause{
	{nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType, 2, []uint8{0x00, 0x01}},
}

var pDUSessionReactivationResultErrorCauseExpectedTestData = []nasType.PDUSessionReactivationResultErrorCause{
	{nasMessage.RegistrationAcceptPDUSessionReactivationResultErrorCauseType, 2, []uint8{0x00, 0x01}},
}

var pDUSessionReactivationResultErrorCauseInformationTable = []testPDUSessionReactivationResultErrorCauseDataTemplate{
	{pDUSessionReactivationResultErrorCauseTestData[0], pDUSessionReactivationResultErrorCauseExpectedTestData[0]},
}

func TestNasTypePDUSessionReactivationResultErrorCauseData(t *testing.T) {}
