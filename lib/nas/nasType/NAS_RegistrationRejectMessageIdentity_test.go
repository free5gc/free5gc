//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRegistrationRejectMessageIdentity(t *testing.T) {}

type nasTypeRegistrationRejectMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeRegistrationRejectMessageIdentityTable = []nasTypeRegistrationRejectMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeRegistrationRejectMessageIdentityGetSetMessageType(t *testing.T) {}

type RegistrationRejectMessageIdentityTestDataTemplate struct {
	in  nasType.RegistrationRejectMessageIdentity
	out nasType.RegistrationRejectMessageIdentity
}

var RegistrationRejectMessageIdentityTestData = []nasType.RegistrationRejectMessageIdentity{
	{0x03},
}

var RegistrationRejectMessageIdentityExpectedTestData = []nasType.RegistrationRejectMessageIdentity{
	{0x03},
}

var RegistrationRejectMessageIdentityTable = []RegistrationRejectMessageIdentityTestDataTemplate{
	{RegistrationRejectMessageIdentityTestData[0], RegistrationRejectMessageIdentityExpectedTestData[0]},
}

func TestNasTypeRegistrationRejectMessageIdentity(t *testing.T) {}
