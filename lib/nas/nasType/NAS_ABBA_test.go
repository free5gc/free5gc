//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewABBA(t *testing.T) {}

var nasTypeAuthenticationResultABBATable = []NasTypeIeiData{
	{nasMessage.AuthenticationResultABBAType, nasMessage.AuthenticationResultABBAType},
}

func TestNasTypeABBAGetSetIei(t *testing.T) {}

var nasTypeAuthenticationResultABBALenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeABBAGetSetLen(t *testing.T) {}

type nasTypeContentData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeContentTable = []nasTypeContentData{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeABBAGetSetContent(t *testing.T) {}

type testABBADataTemplate struct {
	in  nasType.ABBA
	out nasType.ABBA
}

var aBBATestData = []nasType.ABBA{
	{nasMessage.AuthenticationResultABBAType, 2, []byte{0x00, 0x00}},
}

var aBBAExpectedTestData = []nasType.ABBA{
	{nasMessage.AuthenticationResultABBAType, 2, []byte{0x00, 0x00}},
}

var aBBATestTable = []testABBADataTemplate{
	{aBBATestData[0], aBBAExpectedTestData[0]},
}

func TestNasTypeABBA(t *testing.T) {}
