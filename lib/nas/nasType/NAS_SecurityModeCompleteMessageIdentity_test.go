//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSecurityModeCompleteMessageIdentity(t *testing.T) {}

type nasTypeSecurityModeCompleteMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeSecurityModeCompleteMessageIdentityTable = []nasTypeSecurityModeCompleteMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeSecurityModeCompleteMessageIdentityGetSetMessageType(t *testing.T) {}

type SecurityModeCompleteMessageIdentityTestDataTemplate struct {
	in  nasType.SecurityModeCompleteMessageIdentity
	out nasType.SecurityModeCompleteMessageIdentity
}

var SecurityModeCompleteMessageIdentityTestData = []nasType.SecurityModeCompleteMessageIdentity{
	{0x03},
}

var SecurityModeCompleteMessageIdentityExpectedTestData = []nasType.SecurityModeCompleteMessageIdentity{
	{0x03},
}

var SecurityModeCompleteMessageIdentityTable = []SecurityModeCompleteMessageIdentityTestDataTemplate{
	{SecurityModeCompleteMessageIdentityTestData[0], SecurityModeCompleteMessageIdentityExpectedTestData[0]},
}

func TestNasTypeSecurityModeCompleteMessageIdentity(t *testing.T) {}
