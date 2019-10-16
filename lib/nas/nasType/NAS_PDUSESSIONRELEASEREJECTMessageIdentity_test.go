//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONRELEASEREJECTMessageIdentity(t *testing.T) {}

type nasTypePDUSESSIONRELEASEREJECTMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONRELEASEREJECTMessageIdentityTable = []nasTypePDUSESSIONRELEASEREJECTMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypePDUSESSIONRELEASEREJECTMessageIdentityGetSetMessageType(t *testing.T) {}

type PDUSESSIONRELEASEREJECTMessageIdentityTestDataTemplate struct {
	in  nasType.PDUSESSIONRELEASEREJECTMessageIdentity
	out nasType.PDUSESSIONRELEASEREJECTMessageIdentity
}

var PDUSESSIONRELEASEREJECTMessageIdentityTestData = []nasType.PDUSESSIONRELEASEREJECTMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASEREJECTMessageIdentityExpectedTestData = []nasType.PDUSESSIONRELEASEREJECTMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASEREJECTMessageIdentityTable = []PDUSESSIONRELEASEREJECTMessageIdentityTestDataTemplate{
	{PDUSESSIONRELEASEREJECTMessageIdentityTestData[0], PDUSESSIONRELEASEREJECTMessageIdentityExpectedTestData[0]},
}

func TestNasTypePDUSESSIONRELEASEREJECTMessageIdentity(t *testing.T) {}
