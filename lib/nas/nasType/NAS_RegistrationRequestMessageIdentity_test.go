//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRegistrationRequestMessageIdentity(t *testing.T) {}

type nasTypeRegistrationRequestMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeRegistrationRequestMessageIdentityTable = []nasTypeRegistrationRequestMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeRegistrationRequestMessageIdentityGetSetMessageType(t *testing.T) {}

type RegistrationRequestMessageIdentityTestDataTemplate struct {
	in  nasType.RegistrationRequestMessageIdentity
	out nasType.RegistrationRequestMessageIdentity
}

var RegistrationRequestMessageIdentityTestData = []nasType.RegistrationRequestMessageIdentity{
	{0x03},
}

var RegistrationRequestMessageIdentityExpectedTestData = []nasType.RegistrationRequestMessageIdentity{
	{0x03},
}

var RegistrationRequestMessageIdentityTable = []RegistrationRequestMessageIdentityTestDataTemplate{
	{RegistrationRequestMessageIdentityTestData[0], RegistrationRequestMessageIdentityExpectedTestData[0]},
}

func TestNasTypeRegistrationRequestMessageIdentity(t *testing.T) {}
