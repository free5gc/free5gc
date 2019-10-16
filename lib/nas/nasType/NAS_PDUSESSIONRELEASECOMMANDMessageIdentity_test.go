//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONRELEASECOMMANDMessageIdentity(t *testing.T) {}

type nasTypePDUSESSIONRELEASECOMMANDMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONRELEASECOMMANDMessageIdentityTable = []nasTypePDUSESSIONRELEASECOMMANDMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypePDUSESSIONRELEASECOMMANDMessageIdentityGetSetMessageType(t *testing.T) {}

type PDUSESSIONRELEASECOMMANDMessageIdentityTestDataTemplate struct {
	in  nasType.PDUSESSIONRELEASECOMMANDMessageIdentity
	out nasType.PDUSESSIONRELEASECOMMANDMessageIdentity
}

var pDUSESSIONRELEASECOMMANDMessageIdentityTestData = []nasType.PDUSESSIONRELEASECOMMANDMessageIdentity{
	{0x03},
}

var pDUSESSIONRELEASECOMMANDMessageIdentityExpectedTestData = []nasType.PDUSESSIONRELEASECOMMANDMessageIdentity{
	{0x03},
}

var pDUSESSIONRELEASECOMMANDMessageIdentityTable = []PDUSESSIONRELEASECOMMANDMessageIdentityTestDataTemplate{
	{pDUSESSIONRELEASECOMMANDMessageIdentityTestData[0], pDUSESSIONRELEASECOMMANDMessageIdentityExpectedTestData[0]},
}

func TestNasTypePDUSESSIONRELEASECOMMANDMessageIdentity(t *testing.T) {}
