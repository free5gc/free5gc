//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNetworkDaylightSavingTime(t *testing.T) {}

var nasTypeNetworkDaylightSavingTimeConfigurationUpdateCommandNetworkDaylightSavingTimeable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType, nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType},
}

func TestNasTypeNetworkDaylightSavingTimeGetSetIei(t *testing.T) {}

var nasTypeNetworkDaylightSavingTimeLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeNetworkDaylightSavingTimeGetSetLen(t *testing.T) {}

type nasTypeNetworkDaylightSavingTimevalueData struct {
	in  uint8
	out uint8
}

var nasTypeNetworkDaylightSavingTimevalueTable = []nasTypeNetworkDaylightSavingTimevalueData{
	{0x03, 0x03},
}

func TestNasTypeNetworkDaylightSavingTimeGetSetvalue(t *testing.T) {}

type testNetworkDaylightSavingTimeDataTemplate struct {
	inIei    uint8
	inLen    uint8
	invalue  uint8
	outIei   uint8
	outLen   uint8
	outvalue uint8
}

var testNetworkDaylightSavingTimeTestTable = []testNetworkDaylightSavingTimeDataTemplate{
	{nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType, 2, 0x03,
		nasMessage.ConfigurationUpdateCommandNetworkDaylightSavingTimeType, 2, 0x03},
}

func TestNasTypeNetworkDaylightSavingTime(t *testing.T) {}
