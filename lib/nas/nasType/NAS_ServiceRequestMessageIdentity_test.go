//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewServiceRequestMessageIdentity(t *testing.T) {}

type nasTypeServiceRequestMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeServiceRequestMessageIdentityTable = []nasTypeServiceRequestMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeServiceRequestMessageIdentityGetSetMessageType(t *testing.T) {}

type ServiceRequestMessageIdentityTestDataTemplate struct {
	in  nasType.ServiceRequestMessageIdentity
	out nasType.ServiceRequestMessageIdentity
}

var ServiceRequestMessageIdentityTestData = []nasType.ServiceRequestMessageIdentity{
	{0x03},
}

var ServiceRequestMessageIdentityExpectedTestData = []nasType.ServiceRequestMessageIdentity{
	{0x03},
}

var ServiceRequestMessageIdentityTable = []ServiceRequestMessageIdentityTestDataTemplate{
	{ServiceRequestMessageIdentityTestData[0], ServiceRequestMessageIdentityExpectedTestData[0]},
}

func TestNasTypeServiceRequestMessageIdentity(t *testing.T) {}
