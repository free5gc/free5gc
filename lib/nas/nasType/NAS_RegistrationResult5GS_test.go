//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewRegistrationResult5GS(t *testing.T) {}

var nasTypeRegistrationAcceptNetworkFeatureSupport5GSTable = []NasTypeIeiData{
	{nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType, nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType},
}

func TestNasTypeRegistrationResult5GSGetSetIei(t *testing.T) {}

var nasTypeRegistrationAcceptNetworkFeatureSupport5GSLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeRegistrationResult5GSGetSetLen(t *testing.T) {}

type NasTypeSMSAlloweduint8Data struct {
	in  uint8
	out uint8
}

var nasTypeRegistrationResult5GSSMSAllowed = []NasTypeSMSAlloweduint8Data{
	{0x01, 0x01},
	// {0x0, 0x0},
}

func TestNasTypeRegistrationResult5GSGetSetSMSAllowed(t *testing.T) {}

type NasTypeRegistrationResultValue5GSuint8Data struct {
	in  uint8
	out uint8
}

var nasTypeRegistrationResult5GSRegistrationResultValue5GS = []NasTypeRegistrationResultValue5GSuint8Data{
	{0x1, 0x1},
	{0x0, 0x0},
}

func TestNasTypeRegistrationResult5GSGetSetRegistrationResultValue5GS(t *testing.T) {}

type testRegistrationResult5GSDataTemplate struct {
	inSMSAllowed                 uint8
	inRegistrationResultValue5GS uint8
	in                           nasType.RegistrationResult5GS
	out                          nasType.RegistrationResult5GS
}

var registrationResult5GSTestData = []nasType.RegistrationResult5GS{
	{nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType, 1, 0x05},
}
var registrationResult5GSExpectedData = []nasType.RegistrationResult5GS{
	{nasMessage.RegistrationAcceptNetworkFeatureSupport5GSType, 1, 0x0f},
}

var registrationResult5GSDataTestTable = []testRegistrationResult5GSDataTemplate{
	{0x07, 0x1F, registrationResult5GSTestData[0], registrationResult5GSExpectedData[0]},
}

func TestNasTypeRegistrationResult5GS(t *testing.T) {}
