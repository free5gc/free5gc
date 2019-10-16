//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewULNASTRANSPORTMessageIdentity(t *testing.T) {}

type nasTypeULNASTRANSPORTMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypeULNASTRANSPORTMessageIdentityTable = []nasTypeULNASTRANSPORTMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypeULNASTRANSPORTMessageIdentityGetSetMessageType(t *testing.T) {}

type ULNASTRANSPORTMessageIdentityTestDataTemplate struct {
	in  nasType.ULNASTRANSPORTMessageIdentity
	out nasType.ULNASTRANSPORTMessageIdentity
}

var ULNASTRANSPORTMessageIdentityTestData = []nasType.ULNASTRANSPORTMessageIdentity{
	{0x03},
}

var ULNASTRANSPORTMessageIdentityExpectedTestData = []nasType.ULNASTRANSPORTMessageIdentity{
	{0x03},
}

var ULNASTRANSPORTMessageIdentityTable = []ULNASTRANSPORTMessageIdentityTestDataTemplate{
	{ULNASTRANSPORTMessageIdentityTestData[0], ULNASTRANSPORTMessageIdentityExpectedTestData[0]},
}

func TestNasTypeULNASTRANSPORTMessageIdentity(t *testing.T) {}
