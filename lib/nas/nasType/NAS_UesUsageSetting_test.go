//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewUesUsageSetting(t *testing.T) {}

var nasTypeUesUsageSettingIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestUesUsageSettingType, nasMessage.RegistrationRequestUesUsageSettingType},
}

func TestNasTypeUesUsageSettingGetSetIei(t *testing.T) {}

var nasTypeUesUsageSettingLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeUesUsageSettingGetSetLen(t *testing.T) {}

type NasTypeUesUsageSettingUesUsageSettingData struct {
	in  uint8
	out uint8
}

var nasTypeUesUsageSettingUesUsageSettingTable = []NasTypeUesUsageSettingUesUsageSettingData{
	{0x1, 0x1},
}

func TestNasTypeUesUsageSettingGetSetUesUsageSetting(t *testing.T) {}

type testUesUsageSettingDataTemplate struct {
	in  nasType.UesUsageSetting
	out nasType.UesUsageSetting
}

var UesUsageSettingTestData = []nasType.UesUsageSetting{
	{nasMessage.RegistrationRequestUesUsageSettingType, 1, 0x01},
}
var UesUsageSettingExpectedData = []nasType.UesUsageSetting{
	{nasMessage.RegistrationRequestUesUsageSettingType, 1, 0x01},
}

var UesUsageSettingDataTestTable = []testUesUsageSettingDataTemplate{
	{UesUsageSettingTestData[0], UesUsageSettingExpectedData[0]},
}

func TestNasTypeUesUsageSetting(t *testing.T) {}
