//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPTI(t *testing.T) {}

type nasTypePTI struct {
	in  uint8
	out uint8
}

var nasTypePTITable = []nasTypePTI{
	{0x03, 0x03},
}

func TestNasTypePTIGetSetPDUSessionIdentity(t *testing.T) {}

type PTITestDataTemplate struct {
	in  nasType.PTI
	out nasType.PTI
}

var PTITestData = []nasType.PTI{
	{0x03},
}

var PTIExpectedTestData = []nasType.PTI{
	{0x03},
}

var PTITable = []PTITestDataTemplate{
	{PTITestData[0], PTIExpectedTestData[0]},
}

func TestNasTypePTI(t *testing.T) {}
