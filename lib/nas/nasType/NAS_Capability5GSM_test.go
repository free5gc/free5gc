//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewCapability5GSM(t *testing.T) {}

var nasTypePDUSessionModificationRequestCapability5GSMTypeTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationRequestCapability5GSMType, nasMessage.PDUSessionModificationRequestCapability5GSMType},
}

func TestNasTypeCapability5GSMGetSetIei(t *testing.T) {}

var nasTypeCapability5GSMLenTable = []NasTypeLenuint8Data{
	{13, 13},
}

func TestNasTypeCapability5GSMGetSetLen(t *testing.T) {}

type nasTypeCapability5GSMMH6PDUData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeCapability5GSMMH6PDUTable = []nasTypeCapability5GSMMH6PDUData{
	{13, 0x01, 0x01},
}

func TestNasTypeCapability5GSMGetSetMH6PDU(t *testing.T) {}

type nasTypeCapability5GSMRqoSData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeCapability5GSMRqoSTable = []nasTypeCapability5GSMRqoSData{
	{12, 0x01, 0x01},
}

func TestNasTypeCapability5GSMGetSetRqoS(t *testing.T) {}

type nasTypeCapability5GSMSpareData struct {
	inLen uint8
	in    [12]uint8
	out   [12]uint8
}

var nasTypeCapability5GSMSpareTable = []nasTypeCapability5GSMSpareData{
	{12, [12]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, [12]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

func TestNasTypeCapability5GSMGetSetSpare(t *testing.T) {}

type testCapability5GSMDataTemplate struct {
	inLen    uint8
	inMH6PDU uint8
	inRqoS   uint8
	inSpare  [12]uint8
	in       nasType.Capability5GSM
	out      nasType.Capability5GSM
}

var capability5GSMTestData = []nasType.Capability5GSM{
	{nasMessage.PDUSessionModificationRequestCapability5GSMType, 13, [13]uint8{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

var capability5GSMExpectedTestData = []nasType.Capability5GSM{
	{nasMessage.PDUSessionModificationRequestCapability5GSMType, 13, [13]uint8{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
}

var capability5GSMTestTable = []testCapability5GSMDataTemplate{
	{13, 0x01, 0x01, [12]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, capability5GSMTestData[0], capability5GSMExpectedTestData[0]},
}

func TestNasTypeCapability5GSM(t *testing.T) {}
