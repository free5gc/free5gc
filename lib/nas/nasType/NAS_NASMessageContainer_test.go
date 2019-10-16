//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewNASMessageContainer(t *testing.T) {}

var nasTypeNASMessageContainerRegistrationRequestAdditionalGUTITable = []NasTypeIeiData{
	{nasMessage.SecurityModeCompleteNASMessageContainerType, nasMessage.SecurityModeCompleteNASMessageContainerType},
}

func TestNasTypeNASMessageContainerGetSetIei(t *testing.T) {}

var nasTypeNASMessageContainerLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeNASMessageContainerGetSetLen(t *testing.T) {}

type nasTypeNASMessageContainerNASMessageContainerContentsData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeNASMessageContainerNASMessageContainerContentsTable = []nasTypeNASMessageContainerNASMessageContainerContentsData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeNASMessageContainerGetSetNASMessageContainerContents(t *testing.T) {}

type testNASMessageContainerDataTemplate struct {
	inIei                          uint8
	inLen                          uint16
	inNASMessageContainerContents  []uint8
	outIei                         uint8
	outLen                         uint16
	outNASMessageContainerContents []uint8
}

var testNASMessageContainerTestTable = []testNASMessageContainerDataTemplate{
	{nasMessage.SecurityModeCompleteNASMessageContainerType, 2, []uint8{0xff, 0xff},
		nasMessage.SecurityModeCompleteNASMessageContainerType, 2, []uint8{0xff, 0xff}},
}

func TestNasTypeNASMessageContainer(t *testing.T) {}
