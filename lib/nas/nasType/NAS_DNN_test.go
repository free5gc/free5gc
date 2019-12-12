//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewDNN(t *testing.T) {}

var nasTypeDNNIeiTable = []NasTypeIeiData{
	{0, 0},
}

func TestNasTypDNNGetSetIei(t *testing.T) {}

var nasTypeDNNLenTable = []NasTypeLenuint8Data{
	{1, 1},
}

func TestNasTypeDNNGetSetLen(t *testing.T) {}

type nasTypetDNNData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeDNNTable = []nasTypetDNNData{
	{8, []uint8{0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74}, []uint8{0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74}},
}

func TestNasTypeDNNGetSetDNNValue(t *testing.T) {}

type testDNNDataTemplate struct {
	in  nasType.DNN
	out nasType.DNN
}

var DNNTestData = []nasType.DNN{
	{0, 7, []byte{0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74}}, //AuthenticationResult
}

var DNNExpectedTestData = []nasType.DNN{
	{0, 8, []byte{0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74}}, //AuthenticationResult
}

var DNNTestTable = []testDNNDataTemplate{
	{DNNTestData[0], DNNExpectedTestData[0]},
}

func TestNasTypeDNN(t *testing.T) {}
