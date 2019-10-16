//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeNgksiAndServiceTypeData struct {
	inTsc                  uint8
	outTsc                 uint8
	inNASKeySetIdentifier  uint8
	outNASKeySetIdentifier uint8
	inServiceTypeValue     uint8
	outServiceTypeValue    uint8
}

var nasTypeNgksiAndServiceTypeTable = []nasTypeNgksiAndServiceTypeData{
	{0x01, 0x01, 0x07, 0x07, 0x7, 0x07},
}

func TestNasTypeNewServiceTypeAndNgksi(t *testing.T) {}

func TestNasTypeGetSetNgksiAndServiceType(t *testing.T) {}
