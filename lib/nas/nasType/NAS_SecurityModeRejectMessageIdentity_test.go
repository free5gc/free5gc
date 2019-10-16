//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSecurityModeRejectMessageIdentity(t *testing.T) {}

type nasTypeSecurityModeRejectMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeSecurityModeRejectMessageIdentityTable = []nasTypeSecurityModeRejectMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeSecurityModeRejectMessageIdentityGetSetMessageType(t *testing.T) {}

type SecurityModeRejectMessageIdentityTestDataTemplate struct {
	in  nasType.SecurityModeRejectMessageIdentity
	out nasType.SecurityModeRejectMessageIdentity
}

var SecurityModeRejectMessageIdentityTestData = []nasType.SecurityModeRejectMessageIdentity{
	{0x03},
}

var SecurityModeRejectMessageIdentityExpectedTestData = []nasType.SecurityModeRejectMessageIdentity{
	{0x03},
}

var SecurityModeRejectMessageIdentityTable = []SecurityModeRejectMessageIdentityTestDataTemplate{
	{SecurityModeRejectMessageIdentityTestData[0], SecurityModeRejectMessageIdentityExpectedTestData[0]},
}

func TestNasTypeSecurityModeRejectMessageIdentity(t *testing.T) {}
