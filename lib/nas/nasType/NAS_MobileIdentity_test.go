//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewMobileIdentity(t *testing.T) {}

var nasTypeMobileIdentityRegistrationRequestAdditionalGUTITable = []NasTypeIeiData{
	{nasMessage.RegistrationRequestAdditionalGUTIType, nasMessage.RegistrationRequestAdditionalGUTIType},
}

func TestNasTypeMobileIdentityGetSetIei(t *testing.T) {}

var nasTypeMobileIdentityLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeMobileIdentityGetSetLen(t *testing.T) {}

type nasTypeMobileIdentityMobileIdentityContentsData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeMobileIdentityMobileIdentityContentsTable = []nasTypeMobileIdentityMobileIdentityContentsData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeMobileIdentityGetSetMobileIdentityContents(t *testing.T) {}

type testMobileIdentityDataTemplate struct {
	inIei                     uint8
	inLen                     uint16
	inMobileIdentityContents  []uint8
	outIei                    uint8
	outLen                    uint16
	outMobileIdentityContents []uint8
}

var testMobileIdentityTestTable = []testMobileIdentityDataTemplate{
	{nasMessage.RegistrationRequestAdditionalGUTIType, 2, []uint8{0xff, 0xff},
		nasMessage.RegistrationRequestAdditionalGUTIType, 2, []uint8{0xff, 0xff}},
}

func TestNasTypeMobileIdentity(t *testing.T) {}
