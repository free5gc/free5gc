//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewServiceAcceptMessageIdentity(t *testing.T) {}

type nasTypeServiceAcceptMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeServiceAcceptMessageIdentityTable = []nasTypeServiceAcceptMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeServiceAcceptMessageIdentityGetSetMessageType(t *testing.T) {}

type ServiceAcceptMessageIdentityTestDataTemplate struct {
	in  nasType.ServiceAcceptMessageIdentity
	out nasType.ServiceAcceptMessageIdentity
}

var ServiceAcceptMessageIdentityTestData = []nasType.ServiceAcceptMessageIdentity{
	{0x03},
}

var ServiceAcceptMessageIdentityExpectedTestData = []nasType.ServiceAcceptMessageIdentity{
	{0x03},
}

var ServiceAcceptMessageIdentityTable = []ServiceAcceptMessageIdentityTestDataTemplate{
	{ServiceAcceptMessageIdentityTestData[0], ServiceAcceptMessageIdentityExpectedTestData[0]},
}

func TestNasTypeServiceAcceptMessageIdentity(t *testing.T) {}
