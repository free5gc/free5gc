//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewIntegrityProtectionMaximumDataRate(t *testing.T) {}

var nasTypePDUSessionModificationRequestIntegrityProtectionMaximumDataRateTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationRequestIntegrityProtectionMaximumDataRateType, nasMessage.PDUSessionModificationRequestIntegrityProtectionMaximumDataRateType},
}

func TestNasTypeIntegrityProtectionMaximumDataRateGetSetIei(t *testing.T) {}

type nasTypeIntegrityProtectionMaximumDataRateMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLinkData struct {
	in  uint8
	out uint8
}

var nasTypeIntegrityProtectionMaximumDataRateMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLinkTable = []nasTypeIntegrityProtectionMaximumDataRateMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLinkData{
	{0xff, 0xff},
}

func TestNasTypeIntegrityProtectionMaximumDataRateGetSetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink(t *testing.T) {}

type nasTypeIntegrityProtectionMaximumDataRateMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLinkData struct {
	in  uint8
	out uint8
}

var nasTypeIntegrityProtectionMaximumDataRateMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLinkTable = []nasTypeIntegrityProtectionMaximumDataRateMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLinkData{
	{0xff, 0xff},
}

func TestNasTypeIntegrityProtectionMaximumDataRateGetSetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink(t *testing.T) {}

type testIntegrityProtectionMaximumDataRateDataTemplate struct {
	inIei                                                             uint8
	inMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink    uint8
	inMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink  uint8
	outIei                                                            uint8
	outMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink   uint8
	outMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink uint8
}

var integrityProtectionMaximumDataRateTestTable = []testIntegrityProtectionMaximumDataRateDataTemplate{
	{nasMessage.PDUSessionModificationRequestIntegrityProtectionMaximumDataRateType, 0xff, 0x11,
		nasMessage.PDUSessionModificationRequestIntegrityProtectionMaximumDataRateType, 0xff, 0x11},
}

func TestNasTypeIntegrityProtectionMaximumDataRate(t *testing.T) {}
