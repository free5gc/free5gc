//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSelectedSSCModeAndSelectedPDUSessionType(t *testing.T) {}

type nasTypeSelectedSSCModeAndSelectedPDUSessionTypeSSCModeData struct {
	in  uint8
	out uint8
}

var nasTypeSelectedSSCModeAndSelectedPDUSessionTypeSSCModeTable = []nasTypeSelectedSSCModeAndSelectedPDUSessionTypeSSCModeData{
	{0x01, 0x01},
}

func TestNasTypeSelectedSSCModeAndSelectedPDUSessionTypeGetSetSSCMode(t *testing.T) {}

type nasTypeSelectedPDUSessionTypeAndSelectedPDUSessionTypePDUSessionTypeData struct {
	in  uint8
	out uint8
}

var nasTypeSelectedPDUSessionTypeAndSelectedPDUSessionTypePDUSessionTypeTable = []nasTypeSelectedPDUSessionTypeAndSelectedPDUSessionTypePDUSessionTypeData{
	{0x01, 0x01},
}

func TestNasTypeSelectedPDUSessionTypeAndSelectedPDUSessionTypeGetSetPDUSessionType(t *testing.T) {}

type SelectedSSCModeAndSelectedPDUSessionTypeTestDataTemplate struct {
	in  nasType.SelectedSSCModeAndSelectedPDUSessionType
	out nasType.SelectedSSCModeAndSelectedPDUSessionType
}

var SelectedSSCModeAndSelectedPDUSessionTypeTestData = []nasType.SelectedSSCModeAndSelectedPDUSessionType{
	{0x00},
}

var SelectedSSCModeAndSelectedPDUSessionTypeExpectedTestData = []nasType.SelectedSSCModeAndSelectedPDUSessionType{
	{0x11},
}

var SelectedSSCModeAndSelectedPDUSessionTypeTable = []SelectedSSCModeAndSelectedPDUSessionTypeTestDataTemplate{
	{SelectedSSCModeAndSelectedPDUSessionTypeTestData[0], SelectedSSCModeAndSelectedPDUSessionTypeExpectedTestData[0]},
}

func TestNasTypeSelectedSSCModeAndSelectedPDUSessionType(t *testing.T) {}
