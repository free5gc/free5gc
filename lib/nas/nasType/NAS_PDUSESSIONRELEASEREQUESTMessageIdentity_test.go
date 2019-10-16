//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONRELEASEREQUESTMessageIdentity(t *testing.T) {}

type nasTypePDUSESSIONRELEASEREQUESTMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONRELEASEREQUESTMessageIdentityTable = []nasTypePDUSESSIONRELEASEREQUESTMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypePDUSESSIONRELEASEREQUESTMessageIdentityGetSetMessageType(t *testing.T) {}

type PDUSESSIONRELEASEREQUESTMessageIdentityTestDataTemplate struct {
	in  nasType.PDUSESSIONRELEASEREQUESTMessageIdentity
	out nasType.PDUSESSIONRELEASEREQUESTMessageIdentity
}

var PDUSESSIONRELEASEREQUESTMessageIdentityTestData = []nasType.PDUSESSIONRELEASEREQUESTMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASEREQUESTMessageIdentityExpectedTestData = []nasType.PDUSESSIONRELEASEREQUESTMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASEREQUESTMessageIdentityTable = []PDUSESSIONRELEASEREQUESTMessageIdentityTestDataTemplate{
	{PDUSESSIONRELEASEREQUESTMessageIdentityTestData[0], PDUSESSIONRELEASEREQUESTMessageIdentityExpectedTestData[0]},
}

func TestNasTypePDUSESSIONRELEASEREQUESTMessageIdentity(t *testing.T) {}
