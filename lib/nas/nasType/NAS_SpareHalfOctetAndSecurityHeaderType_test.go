//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeSpareHalfOctetAndSecurityHeaderType(t *testing.T) {}

type nasTypeSecurityHeaderTypeAndSpareHalfOctetData struct {
	inSecurityHeader  uint8
	inSpareHalfOctet  uint8
	outSecurityHeader uint8
	outSpareHalfOctet uint8
}

var nasTypeSecurityHeaderTypeAndSpareHalfOctetTable = []nasTypeSecurityHeaderTypeAndSpareHalfOctetData{
	{0x8, 0x1, 0x8, 0x01},
}

func TestNasTypeGetSetSpareHalfOctetAndSecurityHeaderType(t *testing.T) {}
