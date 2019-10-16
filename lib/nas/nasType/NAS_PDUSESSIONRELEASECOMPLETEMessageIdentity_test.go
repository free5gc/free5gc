//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONRELEASECOMPLETEMessageIdentity(t *testing.T) {}

type nasTypePDUSESSIONRELEASECOMPLETEMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONRELEASECOMPLETEMessageIdentityTable = []nasTypePDUSESSIONRELEASECOMPLETEMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypePDUSESSIONRELEASECOMPLETEMessageIdentityGetSetMessageType(t *testing.T) {}

type PDUSESSIONRELEASECOMPLETEMessageIdentityTestDataTemplate struct {
	in  nasType.PDUSESSIONRELEASECOMPLETEMessageIdentity
	out nasType.PDUSESSIONRELEASECOMPLETEMessageIdentity
}

var PDUSESSIONRELEASECOMPLETEMessageIdentityTestData = []nasType.PDUSESSIONRELEASECOMPLETEMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASECOMPLETEMessageIdentityExpectedTestData = []nasType.PDUSESSIONRELEASECOMPLETEMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASECOMPLETEMessageIdentityTable = []PDUSESSIONRELEASECOMPLETEMessageIdentityTestDataTemplate{
	{PDUSESSIONRELEASECOMPLETEMessageIdentityTestData[0], PDUSESSIONRELEASECOMPLETEMessageIdentityExpectedTestData[0]},
}

func TestNasTypePDUSESSIONRELEASECOMPLETEMessageIdentity(t *testing.T) {}
