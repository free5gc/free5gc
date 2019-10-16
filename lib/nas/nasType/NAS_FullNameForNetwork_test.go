//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewFullNameForNetwork(t *testing.T) {}

var nasTypeConfigurationUpdateCommandFullNameForNetworkIeiTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandFullNameForNetworkType, nasMessage.ConfigurationUpdateCommandFullNameForNetworkType},
}

func TestNasTypeFullNameForNetworkGetSetIei(t *testing.T) {}

var nasTypeFullNameForNetworkLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeFullNameForNetworkGetSetLen(t *testing.T) {}

type nasTypetFullNameForNetworkExt struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeFullNameForNetworkExtTable = []nasTypetFullNameForNetworkExt{
	{2, 0x01, 0x01},
}

func TestNasTypeFullNameForNetworkGetSetExt(t *testing.T) {}

type nasTypetFullNameForNetworkCodingScheme struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeFullNameForNetworkCodingSchemeTable = []nasTypetFullNameForNetworkCodingScheme{
	{2, 0x07, 0x07},
}

func TestNasTypeFullNameForNetworkGetSetCodingScheme(t *testing.T) {}

type nasTypetFullNameForNetworkAddCI struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeFullNameForNetworkAddCITable = []nasTypetFullNameForNetworkAddCI{
	{2, 0x01, 0x01},
}

func TestNasTypeFullNameForNetworkGetSetAddCI(t *testing.T) {}

type nasTypetFullNameForNetworkNumberOfSpareBitsInLastOctet struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeFullNameForNetworkNumberOfSpareBitsInLastOctetTable = []nasTypetFullNameForNetworkNumberOfSpareBitsInLastOctet{
	{2, 0x07, 0x07},
}

func TestNasTypeFullNameForNetworkGetSetNumberOfSpareBitsInLastOctet(t *testing.T) {}

type nasTypetFullNameForNetworkTextString struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeFullNameForNetworkTextStringTable = []nasTypetFullNameForNetworkTextString{
	{3, []uint8{0x07, 0x07}, []uint8{0x07, 0x07}},
}

func TestNasTypeFullNameForNetworkGetSetTextString(t *testing.T) {}

type testFullNameForNetworkDataTemplate struct {
	inIei                           uint8
	inLen                           uint8
	inExt                           uint8
	inCodingScheme                  uint8
	inAddCI                         uint8
	inNumberOfSpareBitsInLastOctet  uint8
	inTextString                    []uint8
	outIei                          uint8
	outLen                          uint8
	outExt                          uint8
	outCodingScheme                 uint8
	outAddCI                        uint8
	outNumberOfSpareBitsInLastOctet uint8
	outTextString                   []uint8
}

var fullNameForNetworkestTable = []testFullNameForNetworkDataTemplate{
	{nasMessage.ConfigurationUpdateCommandFullNameForNetworkType, 3, 0x01, 0x01, 0x01, 0x01, []uint8{0x01, 0x01}, nasMessage.ConfigurationUpdateCommandFullNameForNetworkType, 3, 0x01, 0x01, 0x01, 0x01, []uint8{0x01, 0x01}},
}

func TestNasTypeFullNameForNetwork(t *testing.T) {}
