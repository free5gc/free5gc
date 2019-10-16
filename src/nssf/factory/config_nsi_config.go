/*
 * NSSF Configuration Factory
 */

package factory

import (
	. "free5gc/lib/openapi/models"
)

type NsiConfig struct {
	Snssai *Snssai `yaml:"snssai"`

	NsiInformationList []NsiInformation `yaml:"nsiInformationList"`
}
