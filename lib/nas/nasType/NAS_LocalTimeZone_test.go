//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewLocalTimeZone(t *testing.T) {}

var nasTypeConfigurationUpdateCommandLocalTimeZoneTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandLocalTimeZoneType, nasMessage.ConfigurationUpdateCommandLocalTimeZoneType},
}

func TestNasTypeLocalTimeZoneGetSetIei(t *testing.T) {}

type nasTypeLocalTimeZoneTimeZoneData struct {
	in  uint8
	out uint8
}

var nasTypeLocalTimeZoneOctetTable = []nasTypeLocalTimeZoneTimeZoneData{
	{0xff, 0xff},
}

func TestNasTypeLocalTimeZoneGetSetTimeZone(t *testing.T) {}

type testLocalTimeZoneDataTemplate struct {
	in  nasType.LocalTimeZone
	out nasType.LocalTimeZone
}

var LocalTimeZoneTestData = []nasType.LocalTimeZone{
	{nasMessage.ConfigurationUpdateCommandLocalTimeZoneType, 0xff},
}

var LocalTimeZoneExpectedTestData = []nasType.LocalTimeZone{
	{nasMessage.ConfigurationUpdateCommandLocalTimeZoneType, 0xff},
}

var LocalTimeZoneTestTable = []testLocalTimeZoneDataTemplate{
	{LocalTimeZoneTestData[0], LocalTimeZoneExpectedTestData[0]},
}

func TestNasTypeLocalTimeZone(t *testing.T) {}
