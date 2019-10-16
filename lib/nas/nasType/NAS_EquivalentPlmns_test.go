//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewEquivalentPlmns(t *testing.T) {}

var nasTypeRegistrationRequestEquivalentPlmnsIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptEquivalentPlmnsType, nasMessage.RegistrationAcceptEquivalentPlmnsType},
}

func TestNasTypeEquivalentPlmnsGetSetIei(t *testing.T) {}

var nasTypeEquivalentPlmnsLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeEquivalentPlmnsGetSetLen(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN1Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN1{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN1(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN1Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN1{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN1(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN1Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN1{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN1(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN1Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN1{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN1(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN1Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN1{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN1(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN1 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN1Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN1{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN1(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN2Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN2{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN2(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN2Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN2{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN2(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN2Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN2{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN2(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN2Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN2{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN2(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN2Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN2{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN2(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN2 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN2Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN2{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN2(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN3Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN3{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN3(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN3Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN3{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN3(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN3Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN3{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN3(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN3Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN3{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN3(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN3Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN3{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN3(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN3 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN3Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN3{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN3(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN4Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN4{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN4(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN4Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN4{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN4(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN4Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN4{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN4(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN4Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN4{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN4(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN4Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN4{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN4(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN4 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN4Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN4{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN4(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN5Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN5{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN5(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN5Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN5{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN5(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN5Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN5{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN5(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN5Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN5{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN5(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN5Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN5{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN5(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN5 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN5Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN5{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN5(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN6Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN6{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN6(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN6Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN6{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN6(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN6Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN6{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN6(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN6Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN6{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN6(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN6Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN6{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN6(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN6 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN6Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN6{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN6(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN7Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN7{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN7(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN7Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN7{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN7(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN7Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN7{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN7(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN7Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN7{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN7(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN7Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN7{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN7(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN7 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN7Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN7{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN7(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN8Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN8{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN8(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN8Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN8{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN8(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN8Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN8{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN8(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN8Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN8{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN8(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN8Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN8{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN8(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN8 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN8Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN8{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN8(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN9Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN9{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN9(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN9Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN9{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN9(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN9Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN9{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN9(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN9Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN9{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN9(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN9Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN9{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN9(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN9 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN9Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN9{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN9(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN10Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN10{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN10(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN10Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN10{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN10(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN10Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN10{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN10(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN10Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN10{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN10(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN10Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN10{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN10(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN10 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN10Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN10{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN10(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN11Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN11{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN11(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN11Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN11{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN11(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN11Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN11{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN11(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN11Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN11{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN11(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN11Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN11{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN11(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN11 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN11Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN11{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN11(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN12Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN12{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN12(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN12Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN12{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN12(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN12Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN12{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN12(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN12Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN12{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN12(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN12Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN12{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN12(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN12 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN12Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN12{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN12(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN13Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN13{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN13(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN13Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN13{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN13(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN13Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN13{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN13(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN13Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN13{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN13(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN13Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN13{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN13(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN13 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN13Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN13{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN13(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN14Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN14{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN14(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN14Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN14{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN14(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN14Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN14{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN14(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN14Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN14{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN14(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN14Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN14{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN14(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN14 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN14Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN14{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN14(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit2PLMN15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit2PLMN15Table = []nasTypeEquivalentPlmnsMCCDigit2PLMN15{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit2PLMN15(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit1PLMN15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit1PLMN15Table = []nasTypeEquivalentPlmnsMCCDigit1PLMN15{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit1PLMN15(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit3PLMN15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit3PLMN15Table = []nasTypeEquivalentPlmnsMNCDigit3PLMN15{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit3PLMN15(t *testing.T) {}

type nasTypeEquivalentPlmnsMCCDigit3PLMN15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMCCDigit3PLMN15Table = []nasTypeEquivalentPlmnsMCCDigit3PLMN15{
	{2, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMCCDigit3PLMN15(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit2PLMN15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit2PLMN15Table = []nasTypeEquivalentPlmnsMNCDigit2PLMN15{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit2PLMN15(t *testing.T) {}

type nasTypeEquivalentPlmnsMNCDigit1PLMN15 struct {
	inLen uint8
	in    uint8
	out   uint8
}

var nasTypeEquivalentPlmnsMNCDigit1PLMN15Table = []nasTypeEquivalentPlmnsMNCDigit1PLMN15{
	{3, 0x01, 0x01},
}

func TestNasTypeEquivalentPlmnsGetSetMNCDigit1PLMN15(t *testing.T) {}

type testEquivalentPlmnsDataTemplate struct {
	inMCCDigit2PLMN1  uint8
	inMCCDigit1PLMN1  uint8
	inMNCDigit3PLMN1  uint8
	inMCCDigit3PLMN1  uint8
	inMNCDigit2PLMN1  uint8
	inMNCDigit1PLMN1  uint8
	inMCCDigit2PLMN2  uint8
	inMCCDigit1PLMN2  uint8
	inMNCDigit3PLMN2  uint8
	inMCCDigit3PLMN2  uint8
	inMNCDigit2PLMN2  uint8
	inMNCDigit1PLMN2  uint8
	inMCCDigit2PLMN3  uint8
	inMCCDigit1PLMN3  uint8
	inMNCDigit3PLMN3  uint8
	inMCCDigit3PLMN3  uint8
	inMNCDigit2PLMN3  uint8
	inMNCDigit1PLMN3  uint8
	inMCCDigit2PLMN4  uint8
	inMCCDigit1PLMN4  uint8
	inMNCDigit3PLMN4  uint8
	inMCCDigit3PLMN4  uint8
	inMNCDigit2PLMN4  uint8
	inMNCDigit1PLMN4  uint8
	inMCCDigit2PLMN5  uint8
	inMCCDigit1PLMN5  uint8
	inMNCDigit3PLMN5  uint8
	inMCCDigit3PLMN5  uint8
	inMNCDigit2PLMN5  uint8
	inMNCDigit1PLMN5  uint8
	inMCCDigit2PLMN6  uint8
	inMCCDigit1PLMN6  uint8
	inMNCDigit3PLMN6  uint8
	inMCCDigit3PLMN6  uint8
	inMNCDigit2PLMN6  uint8
	inMNCDigit1PLMN6  uint8
	inMCCDigit2PLMN7  uint8
	inMCCDigit1PLMN7  uint8
	inMNCDigit3PLMN7  uint8
	inMCCDigit3PLMN7  uint8
	inMNCDigit2PLMN7  uint8
	inMNCDigit1PLMN7  uint8
	inMCCDigit2PLMN8  uint8
	inMCCDigit1PLMN8  uint8
	inMNCDigit3PLMN8  uint8
	inMCCDigit3PLMN8  uint8
	inMNCDigit2PLMN8  uint8
	inMNCDigit1PLMN8  uint8
	inMCCDigit2PLMN9  uint8
	inMCCDigit1PLMN9  uint8
	inMNCDigit3PLMN9  uint8
	inMCCDigit3PLMN9  uint8
	inMNCDigit2PLMN9  uint8
	inMNCDigit1PLMN9  uint8
	inMCCDigit2PLMN10 uint8
	inMCCDigit1PLMN10 uint8
	inMNCDigit3PLMN10 uint8
	inMCCDigit3PLMN10 uint8
	inMNCDigit2PLMN10 uint8
	inMNCDigit1PLMN10 uint8
	inMCCDigit2PLMN11 uint8
	inMCCDigit1PLMN11 uint8
	inMNCDigit3PLMN11 uint8
	inMCCDigit3PLMN11 uint8
	inMNCDigit2PLMN11 uint8
	inMNCDigit1PLMN11 uint8
	inMCCDigit2PLMN12 uint8
	inMCCDigit1PLMN12 uint8
	inMNCDigit3PLMN12 uint8
	inMCCDigit3PLMN12 uint8
	inMNCDigit2PLMN12 uint8
	inMNCDigit1PLMN12 uint8
	inMCCDigit2PLMN13 uint8
	inMCCDigit1PLMN13 uint8
	inMNCDigit3PLMN13 uint8
	inMCCDigit3PLMN13 uint8
	inMNCDigit2PLMN13 uint8
	inMNCDigit1PLMN13 uint8
	inMCCDigit2PLMN14 uint8
	inMCCDigit1PLMN14 uint8
	inMNCDigit3PLMN14 uint8
	inMCCDigit3PLMN14 uint8
	inMNCDigit2PLMN14 uint8
	inMNCDigit1PLMN14 uint8
	inMCCDigit2PLMN15 uint8
	inMCCDigit1PLMN15 uint8
	inMNCDigit3PLMN15 uint8
	inMCCDigit3PLMN15 uint8
	inMNCDigit2PLMN15 uint8
	inMNCDigit1PLMN15 uint8
	in                nasType.EquivalentPlmns
	out               nasType.EquivalentPlmns
}

var equivalentPlmnsTestData = []nasType.EquivalentPlmns{
	{nasMessage.RegistrationAcceptEquivalentPlmnsType, 3, [45]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
}

var equivalentPlmnsExpectedData = []nasType.EquivalentPlmns{
	{nasMessage.RegistrationAcceptEquivalentPlmnsType, 3, [45]byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}},
}

var equivalentPlmnsTestTable = []testEquivalentPlmnsDataTemplate{
	{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, equivalentPlmnsTestData[0], equivalentPlmnsExpectedData[0]},
}

func TestNasTypeEquivalentPlmns(t *testing.T) {}
