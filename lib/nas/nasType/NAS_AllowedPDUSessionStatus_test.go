//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAllowedPDUSessionStatus(t *testing.T) {}

var nasTypeServiceRequestAllowedPDUSessionStatusTable = []NasTypeIeiData{
	{nasMessage.ServiceRequestAllowedPDUSessionStatusType, nasMessage.ServiceRequestAllowedPDUSessionStatusType},
}

func TestNasTypeAllowedPDUSessionStatusGetSetIei(t *testing.T) {}

var nasTypeServiceRequestAllowedPDUSessionStatusLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeAllowedPDUSessionStatusGetSetLen(t *testing.T) {}

type nasTypePSI7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI7Table = []nasTypePSI7{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI7(t *testing.T) {}

type nasTypePSI6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI6Table = []nasTypePSI6{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI6(t *testing.T) {}

type nasTypePSI5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI5Table = []nasTypePSI5{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI5(t *testing.T) {}

type nasTypePSI4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI4Table = []nasTypePSI4{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI4(t *testing.T) {}

type nasTypePSI3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI3Table = []nasTypePSI3{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI3(t *testing.T) {}

type nasTypePSI2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI2Table = []nasTypePSI2{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI2(t *testing.T) {}

type nasTypePSI1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI1Table = []nasTypePSI1{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI1(t *testing.T) {}

type nasTypePSI0 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI0Table = []nasTypePSI0{
	{1, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI0(t *testing.T) {}

type nasTypePSI15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI15Table = []nasTypePSI15{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI15(t *testing.T) {}

type nasTypePSI14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI14Table = []nasTypePSI14{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI14(t *testing.T) {}

type nasTypePSI13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI13Table = []nasTypePSI13{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI13(t *testing.T) {}

type nasTypePSI12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI12Table = []nasTypePSI12{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI12(t *testing.T) {}

type nasTypePSI11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI11Table = []nasTypePSI11{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI11(t *testing.T) {}

type nasTypePSI10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI10Table = []nasTypePSI11{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI10(t *testing.T) {}

type nasTypePSI9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI9Table = []nasTypePSI9{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI9(t *testing.T) {}

type nasTypePSI8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypePSI8Table = []nasTypePSI8{
	{2, 0x01, 0x01},
}

func TestNasTypeAllowedPDUSessionStatusGetSetPSI8(t *testing.T) {}

type testAllowedPDUSessionStatusDataTemplate struct {
	in  nasType.AllowedPDUSessionStatus
	out nasType.AllowedPDUSessionStatus
}

var AllowedPDUSessionStatusTestData = []nasType.AllowedPDUSessionStatus{
	{nasMessage.ServiceRequestAllowedPDUSessionStatusType, 3, []uint8{0xFF, 0xFF, 0xFF}},
}

var AllowedPDUSessionStatusExpectedData = []nasType.AllowedPDUSessionStatus{
	{nasMessage.ServiceRequestAllowedPDUSessionStatusType, 3, []uint8{0xFF, 0xFF, 0xFF}},
}

var allowedPDUSessionStatusTable = []testAllowedPDUSessionStatusDataTemplate{
	{AllowedPDUSessionStatusTestData[0], AllowedPDUSessionStatusExpectedData[0]},
}

func TestNasTypeAllowedPDUSessionStatus(t *testing.T) {}
