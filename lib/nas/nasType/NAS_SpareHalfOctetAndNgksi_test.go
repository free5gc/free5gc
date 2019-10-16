//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeNewNgksiAndSpareHalfOctetData struct {
	inTsc                  uint8
	outTsc                 uint8
	inNASKeySetIdentifier  uint8
	outNASKeySetIdentifier uint8
	inSpareHalfOctet       uint8
	outSpareHalfOctet      uint8
}

var nasTypeNewNgksiAndSpareHalfOctetTable = []nasTypeNewNgksiAndSpareHalfOctetData{
	{0x1, 0x1, 0x7, 0x7, 0x7, 0x7},
}

func TestNasTypeNewSpareHalfOctetAndNgksi(t *testing.T) {}

func TestNasTypeGetSetSpareHalfOctetAndNgksi(t *testing.T) {}
