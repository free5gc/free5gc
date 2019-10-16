//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewLastVisitedRegisteredTAI(t *testing.T) {}

var nasTypeRegistrationRequestLastVisitedRegisteredTAITable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestLastVisitedRegisteredTAIType, nasMessage.RegistrationRequestLastVisitedRegisteredTAIType},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetIei(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAIMCCDigit2Data struct {
	in  uint8
	out uint8
}

var nasTypeLastVisitedRegisteredTAIMCCDigit2Table = []nasTypeLastVisitedRegisteredTAIMCCDigit2Data{
	{0x01, 0x01},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetMCCDigit2(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAIMCCDigit1Data struct {
	in  uint8
	out uint8
}

var nasTypeLastVisitedRegisteredTAIMCCDigit1Table = []nasTypeLastVisitedRegisteredTAIMCCDigit1Data{
	{0x01, 0x01},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetMCCDigit1(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAIMNCDigit3Data struct {
	in  uint8
	out uint8
}

var nasTypeLastVisitedRegisteredTAIMNCDigit3Table = []nasTypeLastVisitedRegisteredTAIMNCDigit3Data{
	{0x01, 0x01},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetMNCDigit3(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAIMCCDigit3Data struct {
	in  uint8
	out uint8
}

var nasTypeLastVisitedRegisteredTAIMCCDigit3Table = []nasTypeLastVisitedRegisteredTAIMCCDigit3Data{
	{0x01, 0x01},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetMCCDigit3(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAIMNCDigit2Data struct {
	in  uint8
	out uint8
}

var nasTypeLastVisitedRegisteredTAIMNCDigit2Table = []nasTypeLastVisitedRegisteredTAIMNCDigit2Data{
	{0x01, 0x01},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetMNCDigit2(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAIMNCDigit1Data struct {
	in  uint8
	out uint8
}

var nasTypeLastVisitedRegisteredTAIMNCDigit1Table = []nasTypeLastVisitedRegisteredTAIMNCDigit1Data{
	{0x01, 0x01},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetMNCDigit1(t *testing.T) {}

type nasTypeLastVisitedRegisteredTAITACData struct {
	in  [3]uint8
	out [3]uint8
}

var nasTypeLastVisitedRegisteredTAITACTable = []nasTypeLastVisitedRegisteredTAITACData{
	{[3]uint8{0x01, 0x01, 0x01}, [3]uint8{0x01, 0x01, 0x01}},
}

func TestNasTypeLastVisitedRegisteredTAIGetSetTAC(t *testing.T) {}

type testLastVisitedRegisteredTAIDataTemplate struct {
	inIei       uint8
	inMCCDigit2 uint8
	inMCCDigit1 uint8
	inMNCDigit3 uint8
	inMCCDigit3 uint8
	inMNCDigit2 uint8
	inMNCDigit1 uint8
	inTAC       [3]uint8

	outIei       uint8
	outMCCDigit2 uint8
	outMCCDigit1 uint8
	outMNCDigit3 uint8
	outMCCDigit3 uint8
	outMNCDigit2 uint8
	outMNCDigit1 uint8
	outTAC       [3]uint8
}

var testLastVisitedRegisteredTAITestTable = []testLastVisitedRegisteredTAIDataTemplate{
	{nasMessage.RegistrationRequestLastVisitedRegisteredTAIType, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, [3]uint8{0x01, 0x01, 0x01},
		nasMessage.RegistrationRequestLastVisitedRegisteredTAIType, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, [3]uint8{0x01, 0x01, 0x01}},
}

func TestNasTypeLastVisitedRegisteredTAI(t *testing.T) {}
