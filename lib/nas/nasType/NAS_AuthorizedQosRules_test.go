//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthorizedQosRules(t *testing.T) {}

var nasTypeAuthenticationRequestAuthorizedQosRulesIeiTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationCommandAuthorizedQosRulesType, nasMessage.PDUSessionModificationCommandAuthorizedQosRulesType},
}

func TestNasTypeAuthorizedQosRulesGetSetIei(t *testing.T) {}

var nasTypeAuthorizedQosRulesLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeAuthorizedQosRulesGetSetLen(t *testing.T) {}

type nasTypetAuthorizedQosRulesQosRule struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeAuthorizedQosRulesTable = []nasTypetAuthorizedQosRulesQosRule{
	{2, []uint8{0x00, 0x01}, []uint8{0x00, 0x1}},
}

func TestNasTypeAuthorizedQosRulesGetSetAuthorizedQosRules(t *testing.T) {}

type testAuthorizedQosRulesDataTemplate struct {
	in  nasType.AuthorizedQosRules
	out nasType.AuthorizedQosRules
}

var AuthorizedQosRulesTestData = []nasType.AuthorizedQosRules{
	{nasMessage.PDUSessionModificationCommandAuthorizedQosRulesType, 2, []byte{0x00, 0x00}}, //AuthenticationResult
}

var AuthorizedQosRulesExpectedData = []nasType.AuthorizedQosRules{
	{nasMessage.PDUSessionModificationCommandAuthorizedQosRulesType, 2, []byte{0x00, 0x00}}, //AuthenticationResult
}

var AuthorizedQosRulesTestTable = []testAuthorizedQosRulesDataTemplate{
	{AuthorizedQosRulesTestData[0], AuthorizedQosRulesExpectedData[0]},
}

func TestNasTypeAuthorizedQosRules(t *testing.T) {}
