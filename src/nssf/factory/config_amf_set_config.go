/*
 * NSSF Configuration Factory
 */

package factory

import (
	. "free5gc/lib/openapi/models"
)

type AmfSetConfig struct {
	AmfSetId string `yaml:"amfSetId"`

	AmfList []string `yaml:"amfList,omitempty"`

	NrfAmfSet string `yaml:"nrfAmfSet,omitempty"`

	SupportedNssaiAvailabilityData []SupportedNssaiAvailabilityData `yaml:"supportedNssaiAvailabilityData"`
}
