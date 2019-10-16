//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNon3GppDeregistrationTimerValue(t *testing.T) {}

var nasTypeNon3GppDeregistrationTimerValueServiceRejectT3346ValueTypeTable = []NasTypeIeiData{
	{nasMessage.ServiceRejectT3346ValueType, nasMessage.ServiceRejectT3346ValueType},
}

func TestNasTypeNon3GppDeregistrationTimerValueGetSetIei(t *testing.T) {}

var nasTypeNon3GppDeregistrationTimerValueLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeNon3GppDeregistrationTimerValueGetSetLen(t *testing.T) {}

type nasTypeNon3GppDeregistrationTimerValueGPRSTimer2ValueData struct {
	in  uint8
	out uint8
}

var nasTypeNon3GppDeregistrationTimerValueGPRSTimer2ValueTable = []nasTypeNon3GppDeregistrationTimerValueGPRSTimer2ValueData{
	{0x0f, 0x0f},
}

func TestNasTypeNon3GppDeregistrationTimerValueGetSetGPRSTimer2Value(t *testing.T) {}

type testNon3GppDeregistrationTimerValueDataTemplate struct {
	inIei              uint8
	inLen              uint8
	inGPRSTimer2Value  uint8
	outIei             uint8
	outLen             uint8
	outGPRSTimer2Value uint8
}

var testNon3GppDeregistrationTimerValueTestTable = []testNon3GppDeregistrationTimerValueDataTemplate{
	{nasMessage.ServiceRejectT3346ValueType, 2, 0x0f,
		nasMessage.ServiceRejectT3346ValueType, 2, 0x0f},
}

func TestNasTypeNon3GppDeregistrationTimerValue(t *testing.T) {}
