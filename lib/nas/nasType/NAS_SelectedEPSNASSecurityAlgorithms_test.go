//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSelectedEPSNASSecurityAlgorithms(t *testing.T) {}

var nasTypePDUSessionReleaseCompleteSelectedEPSNASSecurityAlgorithmsTable = []NasTypeIeiData{
	{nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType, nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType},
}

func TestNasTypeSelectedEPSNASSecurityAlgorithmsGetSetIei(t *testing.T) {}

var nasTypeSelectedEPSNASSecurityAlgorithmsTypeOfCipheringAlgorithmTable = []NasTypeLenuint8Data{
	{0x01, 0x01},
}

func TestNasTypeSelectedEPSNASSecurityAlgorithmsGetSetTypeOfCipheringAlgorithm(t *testing.T) {}

type nasTypeSelectedEPSNASSecurityAlgorithmsTypeOfIntegrityProtectionAlgorithmData struct {
	in  uint8
	out uint8
}

var nasTypeSelectedEPSNASSecurityAlgorithmsTypeOfIntegrityProtectionAlgorithmTable = []nasTypeSelectedEPSNASSecurityAlgorithmsTypeOfIntegrityProtectionAlgorithmData{
	{0x01, 0x01},
}

func TestNasTypeSelectedEPSNASSecurityAlgorithmsGetSetTypeOfIntegrityProtectionAlgorithm(t *testing.T) {}

type testSelectedEPSNASSecurityAlgorithmsDataTemplate struct {
	inTypeOfCipheringAlgorithm           uint8
	inTypeOfIntegrityProtectionAlgorithm uint8
	in                                   nasType.SelectedEPSNASSecurityAlgorithms
	out                                  nasType.SelectedEPSNASSecurityAlgorithms
}

var SelectedEPSNASSecurityAlgorithmsTestData = []nasType.SelectedEPSNASSecurityAlgorithms{
	{nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType, 0x01},
}

var SelectedEPSNASSecurityAlgorithmsExpectedTestData = []nasType.SelectedEPSNASSecurityAlgorithms{
	{nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType, 0x11},
}

var SelectedEPSNASSecurityAlgorithmsTestTable = []testSelectedEPSNASSecurityAlgorithmsDataTemplate{
	{0x01, 0x01, SelectedEPSNASSecurityAlgorithmsTestData[0], SelectedEPSNASSecurityAlgorithmsExpectedTestData[0]},
}

func TestNasTypeSelectedEPSNASSecurityAlgorithms(t *testing.T) {}
