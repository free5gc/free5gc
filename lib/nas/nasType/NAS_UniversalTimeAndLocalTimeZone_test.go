//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewUniversalTimeAndLocalTimeZone(t *testing.T) {}

var nasTypeServiceRequestUniversalTimeAndLocalTimeZoneTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType, nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetIei(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneYear struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneYearTable = []nasTypeUniversalTimeAndLocalTimeZoneYear{
	{2, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetYear(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneMonth struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneMonthTable = []nasTypeUniversalTimeAndLocalTimeZoneMonth{
	{3, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetMonth(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneDay struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneDayTable = []nasTypeUniversalTimeAndLocalTimeZoneDay{
	{2, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetDay(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneHour struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneHourTable = []nasTypeUniversalTimeAndLocalTimeZoneHour{
	{3, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetHour(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneMinute struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneMinuteTable = []nasTypeUniversalTimeAndLocalTimeZoneMinute{
	{2, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetMinute(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneSecond struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneSecondTable = []nasTypeUniversalTimeAndLocalTimeZoneSecond{
	{2, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetSecond(t *testing.T) {}

type nasTypeUniversalTimeAndLocalTimeZoneTimeZone struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUniversalTimeAndLocalTimeZoneTimeZoneTable = []nasTypeUniversalTimeAndLocalTimeZoneTimeZone{
	{2, 0x01, 0x01},
}

func TestNasTypeUniversalTimeAndLocalTimeZoneGetSetTimeZone(t *testing.T) {}

type testUniversalTimeAndLocalTimeZoneDataTemplate struct {
	in  nasType.UniversalTimeAndLocalTimeZone
	out nasType.UniversalTimeAndLocalTimeZone
}

var UniversalTimeAndLocalTimeZoneTestData = []nasType.UniversalTimeAndLocalTimeZone{
	{nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType, [7]uint8{}},
}

var UniversalTimeAndLocalTimeZoneExpectedData = []nasType.UniversalTimeAndLocalTimeZone{
	{nasMessage.ConfigurationUpdateCommandUniversalTimeAndLocalTimeZoneType, [7]uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
}

var UniversalTimeAndLocalTimeZoneTable = []testUniversalTimeAndLocalTimeZoneDataTemplate{
	{UniversalTimeAndLocalTimeZoneTestData[0], UniversalTimeAndLocalTimeZoneExpectedData[0]},
}

func TestNasTypeUniversalTimeAndLocalTimeZone(t *testing.T) {}
