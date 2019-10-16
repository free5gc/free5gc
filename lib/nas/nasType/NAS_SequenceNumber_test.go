//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSequenceNumber(t *testing.T) {}

type nasTypeSequenceNumber struct {
	in  uint8
	out uint8
}

var nasTypeSequenceNumberTable = []nasTypeSequenceNumber{
	{0x03, 0x03},
}

func TestNasTypeSequenceNumberGetSetSQN(t *testing.T) {}

type SequenceNumberTestDataTemplate struct {
	in  nasType.SequenceNumber
	out nasType.SequenceNumber
}

var SequenceNumberTestData = []nasType.SequenceNumber{
	{0x03},
}

var SequenceNumberExpectedTestData = []nasType.SequenceNumber{
	{0x03},
}

var SequenceNumberTable = []SequenceNumberTestDataTemplate{
	{SequenceNumberTestData[0], SequenceNumberExpectedTestData[0]},
}

func TestNasTypeSequenceNumber(t *testing.T) {}
