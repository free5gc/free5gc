//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNgksiAndRegistrationType5GS(t *testing.T) {}

var RegistrationType5GSAndNgksiFORTable = []NasTypeLenuint8Data{
	{0x1, 0x1},
}

func TestNasTypeRegistrationType5GSAndNgksiGetSetFOR(t *testing.T) {}

var RegistrationType5GSAndNgksiRegistrationType5GSTable = []NasTypeLenuint8Data{
	{0x07, 0x7},
}

func TestNasTypeRegistrationType5GSAndNgksiGetSetRegistrationType5GS(t *testing.T) {}

var RegistrationType5GSAndNgksiTSCTable = []NasTypeLenuint8Data{
	{0x1, 0x01},
}

func TestNasTypeRegistrationType5GSAndNgksiGetSetTSC(t *testing.T) {}

var RegistrationType5GSAndNgksiNasKeySetIdentifilerTable = []NasTypeLenuint8Data{
	{0x04, 0x04},
}

func TestNasTypeRegistrationType5GSAndNgksiGetSetNasKeySetIdentifiler(t *testing.T) {}

type testRegistrationType5GSAndNgksiDataTemplate struct {
	inFOR                  uint8
	inRegistrationType5GS  uint8
	inTSC                  uint8
	inNasKeySetIdentifiler uint8
	in                     nasType.NgksiAndRegistrationType5GS
	out                    nasType.NgksiAndRegistrationType5GS
}

var registrationType5GSAndNgksiTestData = []nasType.NgksiAndRegistrationType5GS{
	{0x01},
}

var registrationType5GSAndNgksiExpectedTestData = []nasType.NgksiAndRegistrationType5GS{
	{0x99},
}

var registrationType5GSAndNgksiTestTable = []testRegistrationType5GSAndNgksiDataTemplate{
	{0x1, 0x1, 0x1, 0x1, registrationType5GSAndNgksiTestData[0], registrationType5GSAndNgksiExpectedTestData[0]},
}

func TestNasTypeRegistrationType5GSAndNgksi(t *testing.T) {}
