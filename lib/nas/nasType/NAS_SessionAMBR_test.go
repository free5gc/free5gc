//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSessionAMBR(t *testing.T) {}

var nasTypeSessionAMBRPDUSessionEstablishmentAcceptSessionAMBRTypeTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationCommandSessionAMBRType, nasMessage.PDUSessionModificationCommandSessionAMBRType},
}

func TestNasTypeSessionAMBRGetSetIei(t *testing.T) {}

var nasTypeSessionAMBRLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeSessionAMBRGetSetLen(t *testing.T) {}

type nasTypeSessionAMBRUnitForSessionAMBRForDownlinkValueData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeSessionAMBRUnitForSessionAMBRForDownlinkValueTable = []nasTypeSessionAMBRUnitForSessionAMBRForDownlinkValueData{
	{2, 0x01, 0x01},
}

func TestNasTypeSessionAMBRGetSetUnitForSessionAMBRForDownlink(t *testing.T) {}

type nasTypeSessionAMBRSessionAMBRForDownlinkData struct {
	inLen uint8
	in    [2]uint8
	out   [2]uint8
}

var nasTypeSessionAMBRSessionAMBRForDownlinkTable = []nasTypeSessionAMBRSessionAMBRForDownlinkData{
	{2, [2]uint8{0x01, 0x01}, [2]uint8{0x01, 0x01}},
}

func TestNasTypeSessionAMBRGetSetSessionAMBRForDownlink(t *testing.T) {}

type nasTypeSessionAMBRUnitForSessionAMBRForUplinkValueData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeSessionAMBRUnitForSessionAMBRForUplinkValueTable = []nasTypeSessionAMBRUnitForSessionAMBRForUplinkValueData{
	{2, 0x01, 0x01},
}

func TestNasTypeSessionAMBRGetSetUnitForSessionAMBRForUplink(t *testing.T) {}

type nasTypeSessionAMBRSessionAMBRForUplinkData struct {
	inLen uint8
	in    [2]uint8
	out   [2]uint8
}

var nasTypeSessionAMBRSessionAMBRForUplinkTable = []nasTypeSessionAMBRSessionAMBRForUplinkData{
	{2, [2]uint8{0x01, 0x01}, [2]uint8{0x01, 0x01}},
}

func TestNasTypeSessionAMBRGetSetSessionAMBRForUplink(t *testing.T) {}

type testSessionAMBRDataTemplate struct {
	in  nasType.SessionAMBR
	out nasType.SessionAMBR
}

var sessionAMBRTestData = []nasType.SessionAMBR{
	{nasMessage.PDUSessionModificationCommandSessionAMBRType, 6, [6]uint8{}},
}

var sessionAMBRExpectedTestData = []nasType.SessionAMBR{
	{nasMessage.PDUSessionModificationCommandSessionAMBRType, 6, [6]uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
}

var sessionAMBRTable = []testSessionAMBRDataTemplate{
	{sessionAMBRTestData[0], sessionAMBRExpectedTestData[0]},
}

func TestNasTypeSessionAMBR(t *testing.T) {}
