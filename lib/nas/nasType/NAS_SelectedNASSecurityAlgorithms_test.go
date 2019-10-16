//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSelectedNASSecurityAlgorithms(t *testing.T) {}

var nasTypePDUSessionReleaseCompleteSelectedNASSecurityAlgorithmsTable = []NasTypeIeiData{
	{nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType, nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType},
}

func TestNasTypeSelectedNASSecurityAlgorithmsGetSetIei(t *testing.T) {}

var nasTypeSelectedNASSecurityAlgorithmsTypeOfCipheringAlgorithmTable = []NasTypeLenuint8Data{
	{0x01, 0x01},
}

func TestNasTypeSelectedNASSecurityAlgorithmsGetSetTypeOfCipheringAlgorithm(t *testing.T) {}

type nasTypeSelectedNASSecurityAlgorithmsTypeOfIntegrityProtectionAlgorithmData struct {
	in  uint8
	out uint8
}

var nasTypeSelectedNASSecurityAlgorithmsTypeOfIntegrityProtectionAlgorithmTable = []nasTypeSelectedNASSecurityAlgorithmsTypeOfIntegrityProtectionAlgorithmData{
	{0x01, 0x01},
}

func TestNasTypeSelectedNASSecurityAlgorithmsGetSetTypeOfIntegrityProtectionAlgorithm(t *testing.T) {}

type testSelectedNASSecurityAlgorithmsDataTemplate struct {
	inTypeOfCipheringAlgorithm           uint8
	inTypeOfIntegrityProtectionAlgorithm uint8
	in                                   nasType.SelectedNASSecurityAlgorithms
	out                                  nasType.SelectedNASSecurityAlgorithms
}

var SelectedNASSecurityAlgorithmsTestData = []nasType.SelectedNASSecurityAlgorithms{
	{nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType, 0x01},
}

var SelectedNASSecurityAlgorithmsExpectedTestData = []nasType.SelectedNASSecurityAlgorithms{
	{nasMessage.SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType, 0x11},
}

var SelectedNASSecurityAlgorithmsTestTable = []testSelectedNASSecurityAlgorithmsDataTemplate{
	{0x01, 0x01, SelectedNASSecurityAlgorithmsTestData[0], SelectedNASSecurityAlgorithmsExpectedTestData[0]},
}

func TestNasTypeSelectedNASSecurityAlgorithms(t *testing.T) {}
