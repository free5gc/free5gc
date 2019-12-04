//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewReplayedS1UESecurityCapabilities(t *testing.T) {}

var nasTypeServiceRequestReplayedS1UESecurityCapabilitiesTable = []NasTypeIeiData{
	{nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType, nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetIei(t *testing.T) {}

var nasTypeServiceRequestReplayedS1UESecurityCapabilitiesLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetLen(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA0Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA0{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA0(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA1_128Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA1_128{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA1_128(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA2_128Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA2_128{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA2_128(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA3_128Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA3_128{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA3_128(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA4Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA4{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA4(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA5Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA5{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA5(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA6Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA6{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA6(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEEA7Table = []nasTypeReplayedS1UESecurityCapabilitiesEEA7{
	{1, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEEA7(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA0Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA0{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA0(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA1_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA1_128Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA1_128{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA1_128(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA2_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA2_128Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA2_128{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA2_128(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA3_128 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA3_128Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA3_128{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA3_128(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA4Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA4(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA5Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA5(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA6Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA6{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA6(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesEIA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesEIA7Table = []nasTypeReplayedS1UESecurityCapabilitiesEIA7{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetEIA7(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA0Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA0{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA0(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA1Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA1{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA1(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA2Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA2{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA2(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA3Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA3{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA3(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA4Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA4(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA5Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA5(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA6Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA6{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA6(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUEA7Table = []nasTypeReplayedS1UESecurityCapabilitiesUEA7{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUEA7(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA1Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA1{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA1(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA2Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA2{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA2(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA3Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA3{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA3(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA4Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA4(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA5Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA5(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA6Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA6{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA6(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesUIA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesUIA7Table = []nasTypeReplayedS1UESecurityCapabilitiesUIA7{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetUIA7(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA1Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA1{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA1(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA2Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA2{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA2(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA3Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA3{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA3(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA4Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA4(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA5Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA4{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA5(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA6Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA6{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA6(t *testing.T) {}

type nasTypeReplayedS1UESecurityCapabilitiesGEA7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeReplayedS1UESecurityCapabilitiesGEA7Table = []nasTypeReplayedS1UESecurityCapabilitiesGEA7{
	{5, 0x01, 0x01},
}

func TestNasTypeReplayedS1UESecurityCapabilitiesGetSetGEA7(t *testing.T) {}

type testReplayedS1UESecurityCapabilitiesDataTemplate struct {
	in  nasType.ReplayedS1UESecurityCapabilities
	out nasType.ReplayedS1UESecurityCapabilities
}

var ReplayedS1UESecurityCapabilitiesTestData = []nasType.ReplayedS1UESecurityCapabilities{
	{nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType, 5, []uint8{0x01, 0x01, 0x01, 0x01, 0x01}},
}

var ReplayedS1UESecurityCapabilitiesExpectedData = []nasType.ReplayedS1UESecurityCapabilities{
	{nasMessage.SecurityModeCommandReplayedS1UESecurityCapabilitiesType, 5, []uint8{0xff, 0xff, 0xff, 0x7f, 0x7f}},
}

var ReplayedS1UESecurityCapabilitiesTable = []testReplayedS1UESecurityCapabilitiesDataTemplate{
	{ReplayedS1UESecurityCapabilitiesTestData[0], ReplayedS1UESecurityCapabilitiesExpectedData[0]},
}

func TestNasTypeReplayedS1UESecurityCapabilities(t *testing.T) {}
