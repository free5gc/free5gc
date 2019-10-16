//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSecurityModeCommandMessageIdentity(t *testing.T) {}

type nasTypeSecurityModeCommandMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeSecurityModeCommandMessageIdentityTable = []nasTypeSecurityModeCommandMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeSecurityModeCommandMessageIdentityGetSetMessageType(t *testing.T) {}

type SecurityModeCommandMessageIdentityTestDataTemplate struct {
	in  nasType.SecurityModeCommandMessageIdentity
	out nasType.SecurityModeCommandMessageIdentity
}

var SecurityModeCommandMessageIdentityTestData = []nasType.SecurityModeCommandMessageIdentity{
	{0x03},
}

var SecurityModeCommandMessageIdentityExpectedTestData = []nasType.SecurityModeCommandMessageIdentity{
	{0x03},
}

var SecurityModeCommandMessageIdentityTable = []SecurityModeCommandMessageIdentityTestDataTemplate{
	{SecurityModeCommandMessageIdentityTestData[0], SecurityModeCommandMessageIdentityExpectedTestData[0]},
}

func TestNasTypeSecurityModeCommandMessageIdentity(t *testing.T) {}
