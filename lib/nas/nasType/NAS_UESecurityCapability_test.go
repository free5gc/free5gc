//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewUESecurityCapability(t *testing.T) {}

var nasTypeServiceRequestUESecurityCapabilityTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestUESecurityCapabilityType, nasMessage.RegistrationRequestUESecurityCapabilityType},
}

func TestNasTypeUESecurityCapabilityGetSetIei(t *testing.T) {}

var nasTypeServiceRequestUESecurityCapabilityLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeUESecurityCapabilityGetSetLen(t *testing.T) {}

type nasTypeUESecurityCapabilityEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA0_5GTable = []nasTypeUESecurityCapabilityEA0{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA0_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA1_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA1_128_5GTable = []nasTypeUESecurityCapabilityEA1_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA1_128_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA2_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA2_128_5GTable = []nasTypeUESecurityCapabilityEA2_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA2_128_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA3_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA3_128_5GTable = []nasTypeUESecurityCapabilityEA3_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA3_128_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA4_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA4_5GTable = []nasTypeUESecurityCapabilityEA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA4_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA5_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA5_5GTable = []nasTypeUESecurityCapabilityEA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA5_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA6_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA6_5GTable = []nasTypeUESecurityCapabilityEA6_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA6_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEA7_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEA7_5GTable = []nasTypeUESecurityCapabilityEA7_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEA7_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA0_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA0_5GTable = []nasTypeUESecurityCapabilityIA0_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA0_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA1_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA1_128_5GTable = []nasTypeUESecurityCapabilityIA1_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA1_128_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA2_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA2_128_5GTable = []nasTypeUESecurityCapabilityIA2_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA2_128_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA3_128_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA3_128_5GTable = []nasTypeUESecurityCapabilityIA3_128_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA3_128_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA4_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA4_5GTable = []nasTypeUESecurityCapabilityIA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA4_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA5_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA5_5GTable = []nasTypeUESecurityCapabilityIA4_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA5_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA6_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA6_5GTable = []nasTypeUESecurityCapabilityIA6_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA6_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityIA7_5G struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityIA7_5GTable = []nasTypeUESecurityCapabilityIA7_5G{
	{2, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetIA7_5G(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA0Table = []nasTypeUESecurityCapabilityEEA0{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA0(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA1_128Table = []nasTypeUESecurityCapabilityEEA1_128{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA1_128(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA2_128Table = []nasTypeUESecurityCapabilityEEA2_128{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA2_128(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA3_128Table = []nasTypeUESecurityCapabilityEEA3_128{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA3_128(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA4Table = []nasTypeUESecurityCapabilityEEA4{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA4(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA5Table = []nasTypeUESecurityCapabilityEEA5{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA5(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA6Table = []nasTypeUESecurityCapabilityEEA6{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA6(t *testing.T) {}

type nasTypeUESecurityCapabilityEEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEEA7Table = []nasTypeUESecurityCapabilityEEA7{
	{3, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEEA7(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA0Table = []nasTypeUESecurityCapabilityEIA0{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA0(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA1_128Table = []nasTypeUESecurityCapabilityEIA1_128{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA1_128(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA2_128Table = []nasTypeUESecurityCapabilityEIA2_128{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA2_128(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA3_128Table = []nasTypeUESecurityCapabilityEIA3_128{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA3_128(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA4Table = []nasTypeUESecurityCapabilityEIA4{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA4(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA5Table = []nasTypeUESecurityCapabilityEIA4{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA5(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA6Table = []nasTypeUESecurityCapabilityEIA6{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA6(t *testing.T) {}

type nasTypeUESecurityCapabilityEIA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUESecurityCapabilityEIA7Table = []nasTypeUESecurityCapabilityEIA7{
	{4, 0x01, 0x01},
}

func TestNasTypeUESecurityCapabilityGetSetEIA7(t *testing.T) {}

type nasTypeUESecurityCapabilitySpare struct {
	in  [4]uint8
	out [4]uint8
}

var nasTypeUESecurityCapabilitySpareTable = []nasTypeUESecurityCapabilitySpare{
	{[4]uint8{0x11, 0x12, 0x13, 0x14}, [4]uint8{0x11, 0x12, 0x13, 0x14}},
}

func TestNasTypeUESecurityCapabilityGetSetSpare(t *testing.T) {}

type testUESecurityCapabilityDataTemplate struct {
	in  nasType.UESecurityCapability
	out nasType.UESecurityCapability
}

var UESecurityCapabilityTestData = []nasType.UESecurityCapability{
	{nasMessage.RegistrationRequestUESecurityCapabilityType, 8, []uint8{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}},
}

var UESecurityCapabilityExpectedData = []nasType.UESecurityCapability{
	{nasMessage.RegistrationRequestUESecurityCapabilityType, 8, []uint8{0xff, 0xff, 0xff, 0xff, 0x11, 0x12, 0x13, 0x14}},
}

var UESecurityCapabilityTable = []testUESecurityCapabilityDataTemplate{
	{UESecurityCapabilityTestData[0], UESecurityCapabilityExpectedData[0]},
}

func TestNasTypeUESecurityCapability(t *testing.T) {}
