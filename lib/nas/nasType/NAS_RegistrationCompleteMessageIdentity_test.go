//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRegistrationCompleteMessageIdentity(t *testing.T) {}

type nasTypeRegistrationCompleteMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeRegistrationCompleteMessageIdentityTable = []nasTypeRegistrationCompleteMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeRegistrationCompleteMessageIdentityGetSetMessageType(t *testing.T) {}

type RegistrationCompleteMessageIdentityTestDataTemplate struct {
	in  nasType.RegistrationCompleteMessageIdentity
	out nasType.RegistrationCompleteMessageIdentity
}

var RegistrationCompleteMessageIdentityTestData = []nasType.RegistrationCompleteMessageIdentity{
	{0x03},
}

var RegistrationCompleteMessageIdentityExpectedTestData = []nasType.RegistrationCompleteMessageIdentity{
	{0x03},
}

var RegistrationCompleteMessageIdentityTable = []RegistrationCompleteMessageIdentityTestDataTemplate{
	{RegistrationCompleteMessageIdentityTestData[0], RegistrationCompleteMessageIdentityExpectedTestData[0]},
}

func TestNasTypeRegistrationCompleteMessageIdentity(t *testing.T) {}
