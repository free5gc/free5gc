//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewLADNIndication(t *testing.T) {}

var nasTypeRegistrationRequestLADNIndicationTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestLADNIndicationType, nasMessage.RegistrationRequestLADNIndicationType},
}

func TestNasTypeLADNIndicationGetSetIei(t *testing.T) {}

var nasTypeLADNIndicationLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeLADNIndicationGetSetLen(t *testing.T) {}

type nasTypeLADNIndicationLADNDNNValueData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeLADNIndicationLADNDNNValueTable = []nasTypeLADNIndicationLADNDNNValueData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeLADNIndicationGetSetLADNDNNValue(t *testing.T) {}

type testLADNIndicationDataTemplate struct {
	inIei           uint8
	inLen           uint16
	inLADNDNNValue  []uint8
	outIei          uint8
	outLen          uint16
	outLADNDNNValue []uint8
}

var testLADNIndicationTestTable = []testLADNIndicationDataTemplate{
	{nasMessage.RegistrationRequestLADNIndicationType, 2, []uint8{0xff, 0xff},
		nasMessage.RegistrationRequestLADNIndicationType, 2, []uint8{0xff, 0xff}},
}

func TestNasTypeLADNIndication(t *testing.T) {}
