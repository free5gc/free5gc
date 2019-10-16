//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

var RegistrationAcceptNSSAIInclusionModeTypeIeiInput uint8 = 0x0A

func TestNasTypeNewNSSAIInclusionMode(t *testing.T) {}

var nasTypeNSSAIInclusionModeRegistrationAcceptNSSAIInclusionModeTypeTable = []NasTypeIeiData{
	{RegistrationAcceptNSSAIInclusionModeTypeIeiInput, 0x0A},
}

func TestNasTypeNSSAIInclusionModeGetSetIei(t *testing.T) {}

type nasTypeNSSAIInclusionMode struct {
	inIei                 uint8
	inNSSAIInclusionMode  uint8
	outIei                uint8
	outNSSAIInclusionMode uint8
}

var nasTypeNSSAIInclusionModeTable = []nasTypeNSSAIInclusionMode{
	{RegistrationAcceptNSSAIInclusionModeTypeIeiInput, 0x03,
		0x0A, 0x03},
}

func TestNasTypeNSSAIInclusionMode(t *testing.T) {}
