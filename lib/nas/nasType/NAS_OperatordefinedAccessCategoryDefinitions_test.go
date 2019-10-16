//go:binary-only-package

package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewOperatordefinedAccessCategoryDefinitions(t *testing.T) {}

var nasTypeOperatordefinedAccessCategoryDefinitionsConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsTypeTable = []NasTypeIeiData{
	{nasMessage.ConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsType, nasMessage.ConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsType},
}

func TestNasTypeOperatordefinedAccessCategoryDefinitionsGetSetIei(t *testing.T) {}

var nasTypeOperatordefinedAccessCategoryDefinitionsLenTable = []NasTypeLenUint16Data{
	{2, 2},
}

func TestNasTypeOperatordefinedAccessCategoryDefinitionsGetSetLen(t *testing.T) {}

type nasTypeOperatordefinedAccessCategoryDefinitionsOperatorDefinedAccessCategoryDefintiionData struct {
	inLen uint16
	in    []uint8
	out   []uint8
}

var nasTypeOperatordefinedAccessCategoryDefinitionsOperatorDefinedAccessCategoryDefintiionTable = []nasTypeOperatordefinedAccessCategoryDefinitionsOperatorDefinedAccessCategoryDefintiionData{
	{2, []uint8{0x0f, 0x0f}, []uint8{0x0f, 0x0f}},
}

func TestNasTypeOperatordefinedAccessCategoryDefinitionsGetSetOperatorDefinedAccessCategoryDefintiion(t *testing.T) {}

type testOperatordefinedAccessCategoryDefinitionsDataTemplate struct {
	inIei                                      uint8
	inLen                                      uint16
	inOperatorDefinedAccessCategoryDefintiion  []uint8
	outIei                                     uint8
	outLen                                     uint16
	outOperatorDefinedAccessCategoryDefintiion []uint8
}

var testOperatordefinedAccessCategoryDefinitionsTestTable = []testOperatordefinedAccessCategoryDefinitionsDataTemplate{
	{nasMessage.ConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsType, 2, []uint8{0x0f, 0x0f},
		nasMessage.ConfigurationUpdateCommandOperatordefinedAccessCategoryDefinitionsType, 2, []uint8{0x0f, 0x0f}},
}

func TestNasTypeOperatordefinedAccessCategoryDefinitions(t *testing.T) {}
