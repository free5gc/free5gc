//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewIMEISV(t *testing.T) {}

var nasTypeSecurityModeCompleteIMEISVTypeTable = []NasTypeIeiData{
	{nasMessage.SecurityModeCompleteIMEISVType, nasMessage.SecurityModeCompleteIMEISVType},
}

func TestNasTypeIMEISVGetSetIei(t *testing.T) {}

var nasTypeIMEISVLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeIMEISVGetSetLen(t *testing.T) {}

type nasTypeIMEISVIdentityDigit1 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeIMEISVIdentityDigit1Table = []nasTypeIMEISVIdentityDigit1{
	{2, 0x01, 0x01},
}

func TestNasTypeIMEISVGetSetIdentityDigit1(t *testing.T) {}

type nasTypeIMEISVOddEvenIdic struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeIMEISVOddEvenIdicTable = []nasTypeIMEISVOddEvenIdic{
	{2, 0x01, 0x01},
}

func TestNasTypeIMEISVGetSetOddEvenIdic(t *testing.T) {}

type nasTypeIMEISVTypeOfIdentity struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeIMEISVTypeOfIdentityTable = []nasTypeIMEISVTypeOfIdentity{
	{2, 0x07, 0x07},
}

func TestNasTypeIMEISVGetSetTypeOfIdentity(t *testing.T) {}

type nasTypeIMEISVIdentityDigitP_1 struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeIMEISVIdentityDigitP_1Table = []nasTypeIMEISVIdentityDigitP_1{
	{2, 0x01, 0x01},
}

func TestNasTypeIMEISVGetSetIdentityDigitP_1(t *testing.T) {}

type nasTypeIMEISVGetIdentityDigitP struct {
	inLen uint16
	in    uint8
	out   uint8
}

var nasTypeIMEISVGetIdentityDigitPTable = []nasTypeIMEISVGetIdentityDigitP{
	{2, 0x0f, 0x0f},
}

func TestNasTypeIMEISVGetSetGetIdentityDigitP(t *testing.T) {}

type testIMEISVDataTemplate struct {
	inIei              uint8
	inLen              uint16
	inIdentityDigit1   uint8
	inOddEvenIdic      uint8
	inTypeOfIdentity   uint8
	inIdentityDigitP_1 uint8
	inIdentityDigitP   uint8

	outIei              uint8
	outLen              uint16
	outIdentityDigit1   uint8
	outOddEvenIdic      uint8
	outTypeOfIdentity   uint8
	outIdentityDigitP_1 uint8
	outIdentityDigitP   uint8
}

var iMEISVTestTable = []testIMEISVDataTemplate{
	{nasMessage.SecurityModeCompleteIMEISVType, 2, 0x01, 0x01, 0x01, 0x01, 0x01,
		nasMessage.SecurityModeCompleteIMEISVType, 2, 0x01, 0x01, 0x01, 0x01, 0x01},
}

func TestNasTypeIMEISV(t *testing.T) {}
