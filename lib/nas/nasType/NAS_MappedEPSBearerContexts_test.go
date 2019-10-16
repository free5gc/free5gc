//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewMappedEPSBearerContexts(t *testing.T) {}

var nasTypeRegistrationRequestMappedEPSBearerContextsTable = []NasTypeIeiData{
	{nasMessage.PDUSessionModificationRequestMappedEPSBearerContextsType, nasMessage.PDUSessionModificationRequestMappedEPSBearerContextsType},
}

func TestNasTypeMappedEPSBearerContextsGetSetIei(t *testing.T) {}

var nasTypeMappedEPSBearerContextsLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeMappedEPSBearerContextsGetSetLen(t *testing.T) {}

type nasTypeMappedEPSBearerContextsMappedEPSBearerContextData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeMappedEPSBearerContextsMappedEPSBearerContextTable = []nasTypeMappedEPSBearerContextsMappedEPSBearerContextData{
	{2, []uint8{0xff, 0xff}, []uint8{0xff, 0xff}},
}

func TestNasTypeMappedEPSBearerContextsGetSetMappedEPSBearerContext(t *testing.T) {}

type testMappedEPSBearerContextsDataTemplate struct {
	inIei                     uint8
	inLen                     uint16
	inMappedEPSBearerContext  []uint8
	outIei                    uint8
	outLen                    uint16
	outMappedEPSBearerContext []uint8
}

var testMappedEPSBearerContextsTestTable = []testMappedEPSBearerContextsDataTemplate{
	{nasMessage.PDUSessionModificationRequestMappedEPSBearerContextsType, 2, []uint8{0xff, 0xff},
		nasMessage.PDUSessionModificationRequestMappedEPSBearerContextsType, 2, []uint8{0xff, 0xff}},
}

func TestNasTypeMappedEPSBearerContexts(t *testing.T) {}
