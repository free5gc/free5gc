//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewGUTI5G(t *testing.T) {}

var nasTypeConfigurationUpdateCommandGUTI5GTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandGUTI5GType, nasMessage.ConfigurationUpdateCommandGUTI5GType},
}

func TestNasTypeGUTI5GGetSetIei(t *testing.T) {}

var nasTypeGUTI5GLenTable = []NasTypeLenUint16Data{
	{12, 12},
}

func TestNasTypeGUTI5GGetSetLen(t *testing.T) {}

type nasTypeGUTI5GSpare struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GSpareTable = []nasTypeGUTI5GSpare{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetSpare(t *testing.T) {}

type nasTypeGUTI5GTypeOfIdentity struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GTypeOfIdentityTable = []nasTypeGUTI5GTypeOfIdentity{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetTypeOfIdentity(t *testing.T) {}

type nasTypeGUTI5GMCCDigit2 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GMCCDigit2Table = []nasTypeGUTI5GMCCDigit2{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetMCCDigit2(t *testing.T) {}

type nasTypeGUTI5GMCCDigit1 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GMCCDigit1Table = []nasTypeGUTI5GMCCDigit1{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetMCCDigit1(t *testing.T) {}

type nasTypeGUTI5GMNCDigit3 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GMNCDigit3Table = []nasTypeGUTI5GMNCDigit3{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetMNCDigit3(t *testing.T) {}

type nasTypeGUTI5GMCCDigit3 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GMCCDigit3Table = []nasTypeGUTI5GMCCDigit3{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetMCCDigit3(t *testing.T) {}

type nasTypeGUTI5GMNCDigit2 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GMNCDigit2Table = []nasTypeGUTI5GMNCDigit2{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetMNCDigit2(t *testing.T) {}

type nasTypeGUTI5GMNCDigit1 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GMNCDigit1Table = []nasTypeGUTI5GMNCDigit1{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetMNCDigit1(t *testing.T) {}

type nasTypeGUTI5GAMFRegionID struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GAMFRegionIDTable = []nasTypeGUTI5GAMFRegionID{
	{12, 0x01, 0x01},
}

func TestNasTypeGUTI5GGetSetAMFRegionID(t *testing.T) {}

type nasTypeGUTI5GAMFSetID struct {
	inLen uint16
	in    uint16
	out   uint16
}

var nasTypeGUTI5GAMFSetIDTable = []nasTypeGUTI5GAMFSetID{
	{12, 0x101, 0x101},
}

func TestNasTypeGUTI5GGetSetAMFSetID(t *testing.T) {}

type nasTypeGUTI5GAMFPointer struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeGUTI5GAMFPointerTable = []nasTypeGUTI5GAMFPointer{
	{12, 0x0f, 0x0f},
}

func TestNasTypeGUTI5GGetSetAMFPointer(t *testing.T) {}

type nasTypeGUTI5GTMSI5G struct {
	inLen uint16
	in    [4]uint8
	out   [4]uint8
}

var nasTypeGUTI5GTMSI5GTable = []nasTypeGUTI5GTMSI5G{
	{12, [4]uint8{0x01, 0x01, 0x01, 0x01}, [4]uint8{0x01, 0x01, 0x01, 0x01}},
}

func TestNasTypeGUTI5GGetSetTMSI5G(t *testing.T) {}

type testGUTI5GDataTemplate struct {
	inIei            uint8
	inLen            uint16
	inSpare          uint8
	inTypeOfIdentity uint8
	inMCCDigit2      uint8
	inMCCDigit1      uint8
	inMNCDigit3      uint8
	inMCCDigit3      uint8
	inMNCDigit2      uint8
	inMNCDigit1      uint8
	inAMFRegionID    uint8
	inAMFSetID       uint16
	inAMFPointer     uint8
	inTMSI5G         [4]uint8

	outIei            uint8
	outLen            uint16
	outSpare          uint8
	outTypeOfIdentity uint8
	outMCCDigit2      uint8
	outMCCDigit1      uint8
	outMNCDigit3      uint8
	outMCCDigit3      uint8
	outMNCDigit2      uint8
	outMNCDigit1      uint8
	outAMFRegionID    uint8
	outAMFSetID       uint16
	outAMFPointer     uint8
	outTMSI5G         [4]uint8
}

var gUTI5GTestTable = []testGUTI5GDataTemplate{
	{nasMessage.ConfigurationUpdateCommandGUTI5GType, 12, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x101, 0x01, [4]uint8{0x01, 0x01, 0x01, 0x01},
		nasMessage.ConfigurationUpdateCommandGUTI5GType, 12, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x101, 0x01, [4]uint8{0x01, 0x01, 0x01, 0x01}},
}

func TestNasTypeGUTI5G(t *testing.T) {}
