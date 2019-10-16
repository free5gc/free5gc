//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAdditional5GSecurityInformation(t *testing.T) {}

var nasTypeSecurityModeCommandAdditional5GSecurityInformationTable = []NasTypeIeiData{
	{0x36, 0x36},
}

func TestNasTypeAdditional5GSecurityInformationGetSetIei(t *testing.T) {}

var nasTypeSecurityModeCommandAdditional5GSecurityInformationLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeAdditional5GSecurityInformationGetSetLen(t *testing.T) {}

type NasTypeRINMRuint8Data struct {
	in  uint8
	out uint8
}

var nasTypeAdditional5GSecurityInformationRINMR = []NasTypeRINMRuint8Data{
	{0x1, 0x1},
	{0x0, 0x0},
}

func TestNasTypeAdditional5GSecurityInformationGetSetRINMR(t *testing.T) {}

type NasTypeHDPuint8Data struct {
	in  uint8
	out uint8
}

var nasTypeAdditional5GSecurityInformationHDP = []NasTypeHDPuint8Data{
	{0x1, 0x1},
	{0x0, 0x0},
}

func TestNasTypeAdditional5GSecurityInformationGetSetHDP(t *testing.T) {}
