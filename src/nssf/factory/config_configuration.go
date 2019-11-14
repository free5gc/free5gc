/*
 * NSSF Configuration Factory
 */

package factory

import (
	"free5gc/lib/openapi/models"
)

type Configuration struct {
	NssfName string `yaml:"nssfName,omitempty"`

	Sbi *Sbi `yaml:"sbi"`

	ServiceNameList []models.ServiceName `yaml:"serviceNameList"`

	NrfUri string `yaml:"nrfUri"`

	SupportedPlmnList []models.PlmnId `yaml:"supportedPlmnList,omitempty"`

	SupportedNssaiInPlmnList []SupportedNssaiInPlmn `yaml:"supportedNssaiInPlmnList"`

	NsiList []NsiConfig `yaml:"nsiList,omitempty"`

	AmfSetList []AmfSetConfig `yaml:"amfSetList"`

	AmfList []AmfConfig `yaml:"amfList"`

	TaList []TaConfig `yaml:"taList"`

	MappingListFromPlmn []MappingFromPlmnConfig `yaml:"mappingListFromPlmn"`
}
