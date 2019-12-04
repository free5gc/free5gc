//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewReplayedUESecurityCapabilities(t *testing.T) {}

var nasTypeServiceRequestReplayedUESecurityCapabilitiesTable = []NasTypeIeiData{
	{nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType, nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIei(t *testing.T) {}

var nasTypeServiceRequestReplayedUESecurityCapabilitiesLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetLen(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA0_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA0{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA0_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA1_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA1_128_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA1_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA1_128_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA2_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA2_128_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA2_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA2_128_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA3_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA3_128_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA3_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA3_128_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA4_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA4_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA4_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA5_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA5_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA5_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA6_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA6_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA6_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA6_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEA7_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEA7_5GTable = []nasTypeReplayedUESecurityCapabilitiesEA7_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEA7_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA0_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA0_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA0_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA0_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA1_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA1_128_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA1_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA1_128_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA2_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA2_128_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA2_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA2_128_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA3_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA3_128_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA3_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA3_128_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA4_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA4_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA4_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA5_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA5_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA5_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA6_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA6_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA6_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA6_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesIA7_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesIA7_5GTable = []nasTypeReplayedUESecurityCapabilitiesIA7_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetIA7_5G(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA0Table = []nasTypeReplayedUESecurityCapabilitiesEEA0{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA0(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA1_128Table = []nasTypeReplayedUESecurityCapabilitiesEEA1_128{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA1_128(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA2_128Table = []nasTypeReplayedUESecurityCapabilitiesEEA2_128{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA2_128(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA3_128Table = []nasTypeReplayedUESecurityCapabilitiesEEA3_128{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA3_128(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA4Table = []nasTypeReplayedUESecurityCapabilitiesEEA4{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA4(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA5Table = []nasTypeReplayedUESecurityCapabilitiesEEA5{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA5(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA6Table = []nasTypeReplayedUESecurityCapabilitiesEEA6{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA6(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEEA7Table = []nasTypeReplayedUESecurityCapabilitiesEEA7{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEEA7(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA0Table = []nasTypeReplayedUESecurityCapabilitiesEIA0{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA0(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA1_128Table = []nasTypeReplayedUESecurityCapabilitiesEIA1_128{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA1_128(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA2_128Table = []nasTypeReplayedUESecurityCapabilitiesEIA2_128{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA2_128(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA3_128Table = []nasTypeReplayedUESecurityCapabilitiesEIA3_128{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA3_128(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA4Table = []nasTypeReplayedUESecurityCapabilitiesEIA4{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA4(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA5Table = []nasTypeReplayedUESecurityCapabilitiesEIA4{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA5(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA6Table = []nasTypeReplayedUESecurityCapabilitiesEIA6{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA6(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesEIA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedUESecurityCapabilitiesEIA7Table = []nasTypeReplayedUESecurityCapabilitiesEIA7{
	{2, 0x01, 0x01},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetEIA7(t *testing.T) {}

type nasTypeReplayedUESecurityCapabilitiesSpare struct {
	in  [4]uint8
	out [4]uint8
}

var nasTypeReplayedUESecurityCapabilitiesSpareTable = []nasTypeReplayedUESecurityCapabilitiesSpare{
	// last 2 value of input will be replaced by 0
	{[4]uint8{0x11, 0x12, 0x13, 0x14}, [4]uint8{0x11, 0x12, 0x13, 0x14}},
}

func TestNasTypeReplayedUESecurityCapabilitiesGetSetSpare(t *testing.T) {}

type testReplayedUESecurityCapabilitiesDataTemplate struct {
	in  nasType.ReplayedUESecurityCapabilities
	out nasType.ReplayedUESecurityCapabilities
}

var replayedUESecurityCapabilitiesTestData = []nasType.ReplayedUESecurityCapabilities{
	{nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType, 8, []uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
}

var replayedUESecurityCapabilitiesExpectedData = []nasType.ReplayedUESecurityCapabilities{
	{nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType, 8, []uint8{0xff, 0xff, 0xff, 0xff, 0x11, 0x12, 0x13, 0x14}},
}

var replayedUESecurityCapabilitiesTable = []testReplayedUESecurityCapabilitiesDataTemplate{
	{replayedUESecurityCapabilitiesTestData[0], replayedUESecurityCapabilitiesExpectedData[0]},
}

func TestNasTypeReplayedUESecurityCapabilities(t *testing.T) {}
