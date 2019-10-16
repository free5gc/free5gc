//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAdditionalGUTI(t *testing.T) {}

var nasTypeRegistrationRequestAdditionalGUTITable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestAdditionalGUTIType, nasMessage.RegistrationRequestAdditionalGUTIType},
}

func TestNasTypeAdditionalGUTIGetSetIei(t *testing.T) {}

var nasTypeRegistrationRequestAdditionalGUTILenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeAdditionalGUTIGetSetLen(t *testing.T) {}

var nasTypeSpareTable = []NasTypeLenuint8Data{
	{0x1, 0x1},
}

func TestNasTypeAdditionalGUTIGetSetSpare(t *testing.T) {}

var nasTypeTypeOfIdentityTable = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3, 0x3},
	{0x4, 0x4},
	{0x5, 0x5},
}

func TestNasTypeAdditionalGUTIGetSetTypeOfIdentity(t *testing.T) {}

var nasTypeMCCDigit2Table = []NasTypeLenuint8Data{
	{0x0, 0x00},
	{0x1, 0x01},
	{0x2, 0x02},
	{0x3, 0x03},
	{0x4, 0x04},
	{0x5, 0x05},
}

func TestNasTypeAdditionalGUTIGetSetMCCDigit2(t *testing.T) {}

var nasTypeMCCDigit1Table = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3, 0x3},
	{0x4, 0x4},
	{0x5, 0x5},
}

func TestNasTypeAdditionalGUTIGetSetMCCDigit1(t *testing.T) {}

var nasTypeMNCDigit3Table = []NasTypeLenuint8Data{
	{0x0, 0x00},
	{0x1, 0x01},
	{0x2, 0x02},
	{0x3, 0x03},
	{0x4, 0x04},
	{0x5, 0x05},
}

func TestNasTypeAdditionalGUTIGetSetMNCDigit3(t *testing.T) {}

var nasTypeMCCDigit3Table = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3, 0x3},
	{0x4, 0x4},
	{0x5, 0x5},
}

func TestNasTypeAdditionalGUTIGetSetMCCDigit3(t *testing.T) {}

var nasTypeMNCDigit2Table = []NasTypeLenuint8Data{
	{0x0, 0x00},
	{0x1, 0x01},
	{0x2, 0x02},
	{0x3, 0x03},
	{0x4, 0x04},
	{0x5, 0x05},
}

func TestNasTypeAdditionalGUTIGetSetMNCDigit2(t *testing.T) {}

var nasTypeMNCDigit1Table = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3, 0x3},
	{0x4, 0x4},
	{0x5, 0x5},
}

func TestNasTypeAdditionalGUTIGetSetMNCDigit1(t *testing.T) {}

var nasTypeAMFRegionIDTable = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3, 0x3},
	{0x4, 0x4},
	{0x5, 0x5},
}

func TestNasTypeAdditionalGUTIGetSetAMFRegionID(t *testing.T) {}

var nasTypeAMFSetIDTable = []NasTypeLenUint16Data{
	{0x0000, 0x0000},
	{0x01FF, 0x01FF},
	// {0x200, 0x200},
}

func TestNasTypeAdditionalGUTIGetSetAMFSetID(t *testing.T) {}

var nasTypeAMFPointerTable = []NasTypeLenuint8Data{
	{0x0, 0x0},
	{0x1, 0x1},
	{0x2, 0x2},
	{0x3F, 0x3F},
	{0x4, 0x4},
	{0x1F, 0x1F},
}

func TestNasTypeAdditionalGUTIGetSetAMFPointer(t *testing.T) {}

type nasTypeTMSI5G struct {
	in  [4]uint8
	out [4]uint8
}

var nasTypeTMSI5GTable = []nasTypeTMSI5G{
	{[4]uint8{0x0, 0x0, 0x0, 0x0}, [4]uint8{0x0, 0x0, 0x0, 0x0}},
	{[4]uint8{0xFF, 0xFF, 0xFF, 0xFF}, [4]uint8{0xFF, 0xFF, 0xFF, 0xFF}},
}

func TestNasTypeAdditionalGUTIGetSetTMSI5G(t *testing.T) {}

type nasTypeAdditionalGUTI struct {
	in  nasType.AdditionalGUTI
	out nasType.AdditionalGUTI
}

var additionalGUTITestData = []nasType.AdditionalGUTI{
	{nasMessage.RegistrationRequestAdditionalGUTIType, 0xffff, [11]uint8{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
}

var additionalGUTIExpectedData = []nasType.AdditionalGUTI{
	{nasMessage.RegistrationRequestAdditionalGUTIType, 0xffff, [11]uint8{0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
}

var nasTypeAdditionalGUTITable = []nasTypeAdditionalGUTI{
	{additionalGUTITestData[0], additionalGUTIExpectedData[0]},
}

func TestNasTypeAdditionalGUTI(t *testing.T) {}
