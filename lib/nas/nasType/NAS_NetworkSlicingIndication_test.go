//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ConfigurationUpdateCommandNetworkSlicingIndicationTypeIeiInput uint8 = 0x09

func TestNasTypeNewNetworkSlicingIndication(t *testing.T) {}

var nasTypeConfigurationUpdateCommandNetworkSlicingIndicationTable = []NasTypeIeiData{
	{ConfigurationUpdateCommandNetworkSlicingIndicationTypeIeiInput, 0x09},
}

func TestNasTypeNetworkSlicingIndicationGetSetIei(t *testing.T) {}

type nasTypeNetworkSlicingIndication struct {
	inDCNI   uint8
	outDCNI  uint8
	inNSSCI  uint8
	outNSSCI uint8
	outIei   uint8
}

var nasTypeNetworkSlicingIndicationTable = []nasTypeNetworkSlicingIndication{
	{0x01, 0x01, 0x01, 0x01, 0x09},
}

func TestNasTypeNetworkSlicingIndication(t *testing.T) {}
