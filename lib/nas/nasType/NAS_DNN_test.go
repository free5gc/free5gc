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
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeDNNGetSetDNNValue(t *testing.T) {}

type testDNNDataTemplate struct {
	in  nasType.DNN
	out nasType.DNN
}

var DNNTestData = []nasType.DNN{
	{0, 2, []byte{0x00, 0x00}}, //AuthenticationResult
}

var DNNExpectedTestData = []nasType.DNN{
	{0, 2, []byte{0x00, 0x00}}, //AuthenticationResult
}

var DNNTestTable = []testDNNDataTemplate{
	{DNNTestData[0], DNNExpectedTestData[0]},
}

func TestNasTypeDNN(t *testing.T) {}
