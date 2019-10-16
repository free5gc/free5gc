//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewServiceRejectMessageIdentity(t *testing.T) {}

type nasTypeServiceRejectMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeServiceRejectMessageIdentityTable = []nasTypeServiceRejectMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeServiceRejectMessageIdentityGetSetMessageType(t *testing.T) {}

type ServiceRejectMessageIdentityTestDataTemplate struct {
	in  nasType.ServiceRejectMessageIdentity
	out nasType.ServiceRejectMessageIdentity
}

var ServiceRejectMessageIdentityTestData = []nasType.ServiceRejectMessageIdentity{
	{0x03},
}

var ServiceRejectMessageIdentityExpectedTestData = []nasType.ServiceRejectMessageIdentity{
	{0x03},
}

var ServiceRejectMessageIdentityTable = []ServiceRejectMessageIdentityTestDataTemplate{
	{ServiceRejectMessageIdentityTestData[0], ServiceRejectMessageIdentityExpectedTestData[0]},
}

func TestNasTypeServiceRejectMessageIdentity(t *testing.T) {}
