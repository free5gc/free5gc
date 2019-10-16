//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSMSIndication(t *testing.T) {}

var nasTypeSMSIndicationIeiTable = []NasTypeIeiData{
	{0x01, 0x01},
}

func TestNasTypeSMSIndicationGetSetIei(t *testing.T) {}

type nasTypeSMSIndicationSAIType struct {
	in  uint8
	out uint8
}

var nasTypeSMSIndicationSAITable = []nasTypeSMSIndicationSAIType{
	{0x01, 0x01},
}

func TestNasTypeSMSIndicationGetSetSAI(t *testing.T) {}

type SMSIndicationTestDataTemplate struct {
	in  nasType.SMSIndication
	out nasType.SMSIndication
}

var SMSIndicationTestData = []nasType.SMSIndication{
	{},
}

var SMSIndicationExpectedTestData = []nasType.SMSIndication{
	{0x11},
}

var SMSIndicationTable = []SMSIndicationTestDataTemplate{
	{SMSIndicationTestData[0], SMSIndicationExpectedTestData[0]},
}

func TestNasTypeSMSIndication(t *testing.T) {}
