//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSTATUSMessageIdentity5GMM(t *testing.T) {}

type nasTypeSTATUSMessageIdentity5GMM struct {
	in  uint8
	out uint8
}

var nasTypeSTATUSMessageIdentity5GMMTable = []nasTypeSTATUSMessageIdentity5GMM{
	{0x03, 0x03},
}

func TestNasTypeSTATUSMessageIdentity5GMMGetSetMessageType(t *testing.T) {}

type STATUSMessageIdentity5GMMTestDataTemplate struct {
	in  nasType.STATUSMessageIdentity5GMM
	out nasType.STATUSMessageIdentity5GMM
}

var STATUSMessageIdentity5GMMTestData = []nasType.STATUSMessageIdentity5GMM{
	{0x03},
}

var STATUSMessageIdentity5GMMExpectedTestData = []nasType.STATUSMessageIdentity5GMM{
	{0x03},
}

var STATUSMessageIdentity5GMMTable = []STATUSMessageIdentity5GMMTestDataTemplate{
	{STATUSMessageIdentity5GMMTestData[0], STATUSMessageIdentity5GMMExpectedTestData[0]},
}

func TestNasTypeSTATUSMessageIdentity5GMM(t *testing.T) {}
