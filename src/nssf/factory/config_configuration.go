/*
 * NSSF Configuration Factory
 */

package factory

type Configuration struct {
	SupportedNssaiInPlmnList []SupportedNssaiInPlmn `yaml:"supportedNssaiInPlmnList"`

	NsiList []NsiConfig `yaml:"nsiList,omitempty"`

	AmfSetList []AmfSetConfig `yaml:"amfSetList"`

	AmfList []AmfConfig `yaml:"amfList"`

	TaList []TaConfig `yaml:"taList"`

	MappingListFromPlmn []MappingFromPlmnConfig `yaml:"mappingListFromPlmn"`
}
