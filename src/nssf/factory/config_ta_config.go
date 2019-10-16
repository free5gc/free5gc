/*
 * NSSF Configuration Factory
 */

package factory

import (
	. "free5gc/lib/openapi/models"
)

type TaConfig struct {
	Tai *Tai `yaml:"tai"`

	AccessType *AccessType `yaml:"accessType"`

	SupportedSnssaiList []Snssai `yaml:"supportedSnssaiList"`

	RestrictedSnssaiList []RestrictedSnssai `yaml:"restrictedSnssaiList,omitempty"`
}
