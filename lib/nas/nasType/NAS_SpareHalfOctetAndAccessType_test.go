//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSpareHalfOctetAndAccessType(t *testing.T) {}

type nasTypeAccessType struct {
	in  uint8
	out uint8
}

var nasTypeAccessTypeTable = []nasTypeAccessType{
	{0x03, 0x03},
}

func TestNasTypeGetSetAccessType(t *testing.T) {}

type AccessTypeAndSpareHalfOctetTestDataTemplate struct {
	in  nasType.SpareHalfOctetAndAccessType
	out nasType.SpareHalfOctetAndAccessType
}

var accessTypeAndSpareHalfOctetTestData = []nasType.SpareHalfOctetAndAccessType{
	{0x03},
}

var accessTypeAndSpareHalfOctetExpectedTestData = []nasType.SpareHalfOctetAndAccessType{
	{0x03},
}

var accessTypeAndSpareHalfOctetTable = []AccessTypeAndSpareHalfOctetTestDataTemplate{
	{accessTypeAndSpareHalfOctetTestData[0], accessTypeAndSpareHalfOctetExpectedTestData[0]},
}

func TestNasTypeAccessTypeAndSpareHalfOctet(t *testing.T) {}
