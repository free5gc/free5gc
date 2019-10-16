//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewSORTransparentContainer(t *testing.T) {}

var nasTypeSORTransparentContainerTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptSORTransparentContainerType, nasMessage.RegistrationAcceptSORTransparentContainerType},
}

func TestNasTypeSORTransparentContainerGetSetIei(t *testing.T) {}

type nasTypeSORTransparentContainerLenData struct {
	in  uint16
	out uint16
}

var nasTypeSORTransparentContainerLenTable = []nasTypeSORTransparentContainerLenData{
	{2, 2},
}

func TestNasTypeSORTransparentContainerGetSetLen(t *testing.T) {}

type nasTypeSORTransparentContainerSORContentData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeSORTransparentContainerSORContentTable = []nasTypeSORTransparentContainerSORContentData{
	{2, []uint8{0x01, 0x01}, []uint8{0x01, 0x01}},
}

func TestNasTypeSORTransparentContainerGetSetSORContent(t *testing.T) {}

type testSORTransparentContainerDataTemplate struct {
	in  nasType.SORTransparentContainer
	out nasType.SORTransparentContainer
}

var SORTransparentContainerTestData = []nasType.SORTransparentContainer{
	{nasMessage.RegistrationAcceptSORTransparentContainerType, 2, []uint8{}},
}

var SORTransparentContainerExpectedTestData = []nasType.SORTransparentContainer{
	{nasMessage.RegistrationAcceptSORTransparentContainerType, 2, []uint8{0x01, 0x01}},
}

var SORTransparentContainerTestTable = []testSORTransparentContainerDataTemplate{
	{SORTransparentContainerTestData[0], SORTransparentContainerExpectedTestData[0]},
}

func TestNasTypeSORTransparentContainer(t *testing.T) {}
