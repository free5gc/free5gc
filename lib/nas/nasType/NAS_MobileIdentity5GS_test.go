//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewMobileIdentity5GS(t *testing.T) {}

var nasTypeMobileIdentity5GSRegistrationRequestAdditionalGUTITable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestAdditionalGUTIType, nasMessage.RegistrationRequestAdditionalGUTIType},
}

func TestNasTypeMobileIdentity5GSGetSetIei(t *testing.T) {}

var nasTypeMobileIdentity5GSLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeMobileIdentity5GSGetSetLen(t *testing.T) {}

type nasTypeMobileIdentity5GSMobileIdentity5GSContentsData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeMobileIdentity5GSMobileIdentity5GSContentsTable = []nasTypeMobileIdentity5GSMobileIdentity5GSContentsData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeMobileIdentity5GSGetSetMobileIdentity5GSContents(t *testing.T) {}

type testMobileIdentity5GSDataTemplate struct {
	inIei                        uint8
	inLen                        uint16
	inMobileIdentity5GSContents  []uint8
	outIei                       uint8
	outLen                       uint16
	outMobileIdentity5GSContents []uint8
}

var testMobileIdentity5GSTestTable = []testMobileIdentity5GSDataTemplate{
	{nasMessage.RegistrationRequestAdditionalGUTIType, 2, []uint8{0xff, 0xff},
		nasMessage.RegistrationRequestAdditionalGUTIType, 2, []uint8{0xff, 0xff}},
}

func TestNasTypeMobileIdentity5GS(t *testing.T) {}
