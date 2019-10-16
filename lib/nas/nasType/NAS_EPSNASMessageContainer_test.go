//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewEPSNASMessageContainer(t *testing.T) {}

var nasTypeRegistrationRequestEPSNASMessageContainerIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestEPSNASMessageContainerType, nasMessage.RegistrationRequestEPSNASMessageContainerType},
}

func TestNasTypeEPSNASMessageContainerGetSetIei(t *testing.T) {}

var nasTypeEPSNASMessageContainerLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeEPSNASMessageContainerGetSetLen(t *testing.T) {}

type nasTypeEPSNASMessageContainerEPANASMessageContainer struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeEPSNASMessageContainerEPANASMessageContainerTable = []nasTypeEPSNASMessageContainerEPANASMessageContainer{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeEPSNASMessageContainerGetSetEPANASMessageContainer(t *testing.T) {}

type testEPSNASMessageContainerDataTemplate struct {
	in  nasType.EPSNASMessageContainer
	out nasType.EPSNASMessageContainer
}

var ePSNASMessageContainerTestData = []nasType.EPSNASMessageContainer{
	{nasMessage.RegistrationRequestEPSNASMessageContainerType, 3, []byte{0x02, 0x1f, 0x22}},
}

var ePSNASMessageContainerExpectedData = []nasType.EPSNASMessageContainer{
	{nasMessage.RegistrationRequestEPSNASMessageContainerType, 3, []byte{0x02, 0x1f, 0x22}},
}

var ePSNASMessageContainerTestTable = []testEPSNASMessageContainerDataTemplate{
	{ePSNASMessageContainerTestData[0], ePSNASMessageContainerExpectedData[0]},
}

func TestNasTypeEPSNASMessageContainer(t *testing.T) {}
