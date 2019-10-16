//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPayloadContainer(t *testing.T) {}

var nasTypePayloadContainerRegistrationRequestPayloadContainerTypeTable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestPayloadContainerType, nasMessage.RegistrationRequestPayloadContainerType},
}

func TestNasTypePayloadContainerGetSetIei(t *testing.T) {}

var nasTypePayloadContainerLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypePayloadContainerGetSetLen(t *testing.T) {}

type nasTypePayloadContainerPayloadContainerContentsData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypePayloadContainerPayloadContainerContentsTable = []nasTypePayloadContainerPayloadContainerContentsData{
	{2, []uint8{0x0f, 0x0f}, []uint8{0x0f, 0x0f}},
}

func TestNasTypePayloadContainerGetSetPayloadContainerContents(t *testing.T) {}

type testPayloadContainerDataTemplate struct {
	inIei                       uint8
	inLen                       uint16
	inPayloadContainerContents  []uint8
	outIei                      uint8
	outLen                      uint16
	outPayloadContainerContents []uint8
}

var testPayloadContainerTestTable = []testPayloadContainerDataTemplate{
	{nasMessage.RegistrationRequestPayloadContainerType, 2, []uint8{0x0f, 0x0f},
		nasMessage.RegistrationRequestPayloadContainerType, 2, []uint8{0x0f, 0x0f}},
}

func TestNasTypePayloadContainer(t *testing.T) {}
