//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewUplinkDataStatus(t *testing.T) {}

var nasTypeServiceRequestUplinkDataStatusTable = []NasTypeIeiData{
	{nasMessage.ServiceRequestUplinkDataStatusType, nasMessage.ServiceRequestUplinkDataStatusType},
}

func TestNasTypeUplinkDataStatusGetSetIei(t *testing.T) {}

var nasTypeServiceRequestUplinkDataStatusLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeUplinkDataStatusGetSetLen(t *testing.T) {}

type nasTypeUplinkDataStatusPSI7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI7Table = []nasTypeUplinkDataStatusPSI7{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI7(t *testing.T) {}

type nasTypeUplinkDataStatusPSI6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI6Table = []nasTypeUplinkDataStatusPSI6{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI6(t *testing.T) {}

type nasTypeUplinkDataStatusPSI5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI5Table = []nasTypeUplinkDataStatusPSI5{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI5(t *testing.T) {}

type nasTypeUplinkDataStatusPSI4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI4Table = []nasTypeUplinkDataStatusPSI4{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI4(t *testing.T) {}

type nasTypeUplinkDataStatusPSI3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI3Table = []nasTypeUplinkDataStatusPSI3{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI3(t *testing.T) {}

type nasTypeUplinkDataStatusPSI2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI2Table = []nasTypeUplinkDataStatusPSI2{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI2(t *testing.T) {}

type nasTypeUplinkDataStatusPSI1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI1Table = []nasTypeUplinkDataStatusPSI1{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI1(t *testing.T) {}

type nasTypeUplinkDataStatusPSI0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI0Table = []nasTypeUplinkDataStatusPSI0{
	{1, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI0(t *testing.T) {}

type nasTypeUplinkDataStatusPSI15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI15Table = []nasTypeUplinkDataStatusPSI15{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI15(t *testing.T) {}

type nasTypeUplinkDataStatusPSI14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI14Table = []nasTypeUplinkDataStatusPSI14{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI14(t *testing.T) {}

type nasTypeUplinkDataStatusPSI13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI13Table = []nasTypeUplinkDataStatusPSI13{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI13(t *testing.T) {}

type nasTypeUplinkDataStatusPSI12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI12Table = []nasTypeUplinkDataStatusPSI12{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI12(t *testing.T) {}

type nasTypeUplinkDataStatusPSI11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI11Table = []nasTypeUplinkDataStatusPSI11{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI11(t *testing.T) {}

type nasTypeUplinkDataStatusPSI10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI10Table = []nasTypeUplinkDataStatusPSI11{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI10(t *testing.T) {}

type nasTypeUplinkDataStatusPSI9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI9Table = []nasTypeUplinkDataStatusPSI9{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI9(t *testing.T) {}

type nasTypeUplinkDataStatusPSI8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeUplinkDataStatusPSI8Table = []nasTypeUplinkDataStatusPSI8{
	{2, 0x01, 0x01},
}

func TestNasTypeUplinkDataStatusGetSetPSI8(t *testing.T) {}

type nasTypeUplinkDataStatusSpare struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypeUplinkDataStatusSpareTable = []nasTypeUplinkDataStatusSpare{
	// Spare Len: 5-2 = 3, since size of 2 are reserved for IE and LENGTH
	{5, []uint8{0x11, 0x12, 0x13}, []uint8{0x11, 0x12, 0x13}},
	{5, []uint8{0x12, 0x11, 0x13}, []uint8{0x12, 0x11, 0x13}},
}

func TestNasTypeUplinkDataStatusGetSetSpare(t *testing.T) {}

type testUplinkDataStatusDataTemplate struct {
	in  nasType.UplinkDataStatus
	out nasType.UplinkDataStatus
}

/*
	For the 1st testcase with len 2, our input for SetSpare function will not be read
	as the len size is too small(< 3). However, SetSpare function won't raise any error
	since the "make" function in golang will create a zero-length slice instead of nil slice.

	REFERENCE: https://programming.guide/go/nil-slice-vs-empty-slice.html
*/
var UplinkDataStatusTestData = []nasType.UplinkDataStatus{
	{nasMessage.ServiceRequestUplinkDataStatusType, 2, []uint8{}},
	{nasMessage.ServiceRequestUplinkDataStatusType, 3, []uint8{}},
	{nasMessage.ServiceRequestUplinkDataStatusType, 4, []uint8{}},
	{nasMessage.ServiceRequestUplinkDataStatusType, 5, []uint8{}},
}

var UplinkDataStatusExpectedData = []nasType.UplinkDataStatus{
	{nasMessage.ServiceRequestUplinkDataStatusType, 2, []uint8{0xFF, 0xFF}},
	{nasMessage.ServiceRequestUplinkDataStatusType, 3, []uint8{0xFF, 0xFF, 0x14}},
	{nasMessage.ServiceRequestUplinkDataStatusType, 4, []uint8{0xFF, 0xFF, 0x14, 0x15}},
	{nasMessage.ServiceRequestUplinkDataStatusType, 5, []uint8{0xFF, 0xFF, 0x14, 0x15, 0x16}},
}

var UplinkDataStatusTable = []testUplinkDataStatusDataTemplate{
	{UplinkDataStatusTestData[0], UplinkDataStatusExpectedData[0]},
	{UplinkDataStatusTestData[1], UplinkDataStatusExpectedData[1]},
	{UplinkDataStatusTestData[2], UplinkDataStatusExpectedData[2]},
	{UplinkDataStatusTestData[3], UplinkDataStatusExpectedData[3]},
}

func TestNasTypeUplinkDataStatus(t *testing.T) {}
