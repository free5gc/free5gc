//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var RegistrationRequestNoncurrentNativeNASKeySetIdentifierTypeIeiInput uint8 = 0x0C

func TestNasTypeNewNoncurrentNativeNASKeySetIdentifier(t *testing.T) {}

var nasTypeConfigurationUpdateCommandNoncurrentNativeNASKeySetIdentifierTable = []NasTypeIeiData{
	{RegistrationRequestNoncurrentNativeNASKeySetIdentifierTypeIeiInput, 0x0C},
}

func TestNasTypeNoncurrentNativeNASKeySetIdentifierGetSetIei(t *testing.T) {}

type nasTypeNoncurrentNativeNASKeySetIdentifier struct {
	inIei                   uint8
	inTsc                   uint8
	inNasKeySetIdentifiler  uint8
	outIei                  uint8
	outTsc                  uint8
	outNasKeySetIdentifiler uint8
}

var nasTypeNoncurrentNativeNASKeySetIdentifierTable = []nasTypeNoncurrentNativeNASKeySetIdentifier{
	{RegistrationRequestNoncurrentNativeNASKeySetIdentifierTypeIeiInput, 0x01, 0x01,
		0x0C, 0x01, 0x01},
}

func TestNasTypeNoncurrentNativeNASKeySetIdentifier(t *testing.T) {}
