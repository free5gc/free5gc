//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSessionStatus(t *testing.T) {}

var nasTypeServiceRequestPDUSessionStatusTable = []NasTypeIeiData{
	{nasMessage.ServiceRequestPDUSessionStatusType, nasMessage.ServiceRequestPDUSessionStatusType},
}

func TestNasTypePDUSessionStatusGetSetIei(t *testing.T) {}

var nasTypeServiceRequestPDUSessionStatusLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypePDUSessionStatusGetSetLen(t *testing.T) {}

type nasTypePDUSessionStatusPSI7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI7Table = []nasTypePDUSessionStatusPSI7{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI7(t *testing.T) {}

type nasTypePDUSessionStatusPSI6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI6Table = []nasTypePDUSessionStatusPSI6{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI6(t *testing.T) {}

type nasTypePDUSessionStatusPSI5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI5Table = []nasTypePDUSessionStatusPSI5{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI5(t *testing.T) {}

type nasTypePDUSessionStatusPSI4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI4Table = []nasTypePDUSessionStatusPSI4{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI4(t *testing.T) {}

type nasTypePDUSessionStatusPSI3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI3Table = []nasTypePDUSessionStatusPSI3{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI3(t *testing.T) {}

type nasTypePDUSessionStatusPSI2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI2Table = []nasTypePDUSessionStatusPSI2{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI2(t *testing.T) {}

type nasTypePDUSessionStatusPSI1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI1Table = []nasTypePDUSessionStatusPSI1{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI1(t *testing.T) {}

type nasTypePDUSessionStatusPSI0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI0Table = []nasTypePDUSessionStatusPSI0{
	{1, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI0(t *testing.T) {}

type nasTypePDUSessionStatusPSI15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI15Table = []nasTypePDUSessionStatusPSI15{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI15(t *testing.T) {}

type nasTypePDUSessionStatusPSI14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI14Table = []nasTypePDUSessionStatusPSI14{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI14(t *testing.T) {}

type nasTypePDUSessionStatusPSI13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI13Table = []nasTypePDUSessionStatusPSI13{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI13(t *testing.T) {}

type nasTypePDUSessionStatusPSI12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI12Table = []nasTypePDUSessionStatusPSI12{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI12(t *testing.T) {}

type nasTypePDUSessionStatusPSI11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI11Table = []nasTypePDUSessionStatusPSI11{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI11(t *testing.T) {}

type nasTypePDUSessionStatusPSI10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI10Table = []nasTypePDUSessionStatusPSI11{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI10(t *testing.T) {}

type nasTypePDUSessionStatusPSI9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI9Table = []nasTypePDUSessionStatusPSI9{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI9(t *testing.T) {}

type nasTypePDUSessionStatusPSI8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePDUSessionStatusPSI8Table = []nasTypePDUSessionStatusPSI8{
	{2, 0x01, 0x01},
}

func TestNasTypePDUSessionStatusGetSetPSI8(t *testing.T) {}

type nasTypePDUSessionStatusSpare struct {
	inLen uint8
	in    []uint8
	out   []uint8
}

var nasTypePDUSessionStatusSpareTable = []nasTypePDUSessionStatusSpare{
	// Spare Len: 5-2 = 3, since size of 2 are reserved for IE and LENGTH
	{5, []uint8{0x11, 0x12, 0x13}, []uint8{0x11, 0x12, 0x13}},
	{5, []uint8{0x12, 0x11, 0x13}, []uint8{0x12, 0x11, 0x13}},
}

func TestNasTypePDUSessionStatusGetSetSpare(t *testing.T) {}

type testPDUSessionStatusDataTemplate struct {
	in  nasType.PDUSessionStatus
	out nasType.PDUSessionStatus
}

/*
	For the 1st testcase with len 2, our input for SetSpare function will not be read
	as the len size is too small(< 3). However, SetSpare function won't raise any error
	since the "make" function in golang will create a zero-length slice instead of nil slice.

	REFERENCE: https://programming.guide/go/nil-slice-vs-empty-slice.html
*/
var PDUSessionStatusTestData = []nasType.PDUSessionStatus{
	{nasMessage.ServiceRequestPDUSessionStatusType, 2, []uint8{}},
	{nasMessage.ServiceRequestPDUSessionStatusType, 3, []uint8{}},
	{nasMessage.ServiceRequestPDUSessionStatusType, 4, []uint8{}},
	{nasMessage.ServiceRequestPDUSessionStatusType, 5, []uint8{}},
}

var PDUSessionStatusExpectedData = []nasType.PDUSessionStatus{
	{nasMessage.ServiceRequestPDUSessionStatusType, 2, []uint8{0xFF, 0xFF}},
	{nasMessage.ServiceRequestPDUSessionStatusType, 3, []uint8{0xFF, 0xFF, 0x14}},
	{nasMessage.ServiceRequestPDUSessionStatusType, 4, []uint8{0xFF, 0xFF, 0x14, 0x15}},
	{nasMessage.ServiceRequestPDUSessionStatusType, 5, []uint8{0xFF, 0xFF, 0x14, 0x15, 0x16}},
}

var PDUSessionStatusTable = []testPDUSessionStatusDataTemplate{
	{PDUSessionStatusTestData[0], PDUSessionStatusExpectedData[0]},
	{PDUSessionStatusTestData[1], PDUSessionStatusExpectedData[1]},
	{PDUSessionStatusTestData[2], PDUSessionStatusExpectedData[2]},
	{PDUSessionStatusTestData[3], PDUSessionStatusExpectedData[3]},
}

func TestNasTypePDUSessionStatus(t *testing.T) {}
