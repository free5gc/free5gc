//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSTATUSMessageIdentity5GSM(t *testing.T) {}

type nasTypeSTATUSMessageIdentity5GSM struct {
	in  uint8
	out uint8
}

var nasTypeSTATUSMessageIdentity5GSMTable = []nasTypeSTATUSMessageIdentity5GSM{
	{0x03, 0x03},
}

func TestNasTypeSTATUSMessageIdentity5GSMGetSetMessageType(t *testing.T) {}

type STATUSMessageIdentity5GSMTestDataTemplate struct {
	in  nasType.STATUSMessageIdentity5GSM
	out nasType.STATUSMessageIdentity5GSM
}

var STATUSMessageIdentity5GSMTestData = []nasType.STATUSMessageIdentity5GSM{
	{0x03},
}

var STATUSMessageIdentity5GSMExpectedTestData = []nasType.STATUSMessageIdentity5GSM{
	{0x03},
}

var STATUSMessageIdentity5GSMTable = []STATUSMessageIdentity5GSMTestDataTemplate{
	{STATUSMessageIdentity5GSMTestData[0], STATUSMessageIdentity5GSMExpectedTestData[0]},
}

func TestNasTypeSTATUSMessageIdentity5GSM(t *testing.T) {}
