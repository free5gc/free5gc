//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewShortNameForNetwork(t *testing.T) {}

var nasTypeShortNameForNetworkTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandShortNameForNetworkType, nasMessage.ConfigurationUpdateCommandShortNameForNetworkType},
}

func TestNasTypeShortNameForNetworkGetSetIei(t *testing.T) {}

var nasTypeShortNameForNetworkLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeShortNameForNetworkGetSetLen(t *testing.T) {}

type nasTypeShortNameForNetworkExtData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeShortNameForNetworkExtTable = []nasTypeShortNameForNetworkExtData{
	{2, 0x01, 0x01},
}

func TestNasTypeShortNameForNetworkGetSetExt(t *testing.T) {}

type nasTypeShortNameForNetworkCodingSchemeData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeShortNameForNetworkCodingSchemeTable = []nasTypeShortNameForNetworkCodingSchemeData{
	{2, 0x01, 0x01},
}

func TestNasTypeShortNameForNetworkGetSetCodingScheme(t *testing.T) {}

type nasTypeShortNameForNetworkAddCIData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeShortNameForNetworkAddCITable = []nasTypeShortNameForNetworkAddCIData{
	{2, 0x01, 0x01},
}

func TestNasTypeShortNameForNetworkGetSetAddCI(t *testing.T) {}

type nasTypeShortNameForNetworkNumberOfSpareBitsInLastOctettData struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeShortNameForNetworkNumberOfSpareBitsInLastOctetTable = []nasTypeShortNameForNetworkNumberOfSpareBitsInLastOctettData{
	{2, 0x01, 0x01},
}

func TestNasTypeShortNameForNetworkGetSetNumberOfSpareBitsInLastOctet(t *testing.T) {}

type nasTypeShortNameForNetworkTextStringData struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeShortNameForNetworkTextStringTable = []nasTypeShortNameForNetworkTextStringData{
	{3, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeShortNameForNetworkGetSetTextString(t *testing.T) {}

type testShortNameForNetworkDataTemplate struct {
	in  nasType.ShortNameForNetwork
	out nasType.ShortNameForNetwork
}

var ShortNameForNetworkTestData = []nasType.ShortNameForNetwork{
	{nasMessage.ConfigurationUpdateCommandShortNameForNetworkType, 3, []uint8{}},
}

var ShortNameForNetworkExpectedTestData = []nasType.ShortNameForNetwork{
	{nasMessage.ConfigurationUpdateCommandShortNameForNetworkType, 3, []uint8{0x99, 0x01, 0x01}},
}

var ShortNameForNetworkTestTable = []testShortNameForNetworkDataTemplate{
	{ShortNameForNetworkTestData[0], ShortNameForNetworkExpectedTestData[0]},
}

func TestNasTypeShortNameForNetwork(t *testing.T) {}
