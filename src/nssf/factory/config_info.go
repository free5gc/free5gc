/*
 * NSSF Configuration Factory
 */

package factory

import (
	. "free5gc/lib/openapi/models"
)

type Info struct {
	Version string `yaml:"version,omitempty"`

	Description string `yaml:"description,omitempty"`

	Host string `yaml:"host"`

	ServingPlmnIdList []PlmnId `yaml:"servingPlmnIdList"`
}
