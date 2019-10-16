//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRegistrationAcceptMessageIdentity(t *testing.T) {}

type nasTypeRegistrationAcceptMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeRegistrationAcceptMessageIdentityTable = []nasTypeRegistrationAcceptMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeRegistrationAcceptMessageIdentityGetSetMessageType(t *testing.T) {}

type RegistrationAcceptMessageIdentityTestDataTemplate struct {
	in  nasType.RegistrationAcceptMessageIdentity
	out nasType.RegistrationAcceptMessageIdentity
}

var RegistrationAcceptMessageIdentityTestData = []nasType.RegistrationAcceptMessageIdentity{
	{0x03},
}

var RegistrationAcceptMessageIdentityExpectedTestData = []nasType.RegistrationAcceptMessageIdentity{
	{0x03},
}

var RegistrationAcceptMessageIdentityTable = []RegistrationAcceptMessageIdentityTestDataTemplate{
	{RegistrationAcceptMessageIdentityTestData[0], RegistrationAcceptMessageIdentityExpectedTestData[0]},
}

func TestNasTypeRegistrationAcceptMessageIdentity(t *testing.T) {}
