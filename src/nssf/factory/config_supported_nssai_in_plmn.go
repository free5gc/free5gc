/*
 * NSSF Configuration Factory
 */

package factory

import (
	. "free5gc/lib/openapi/models"
)

type SupportedNssaiInPlmn struct {
	PlmnId *PlmnId `yaml:"plmnId"`

	SupportedSnssaiList []Snssai `yaml:"supportedSnssaiList"`
}
