/*
 * AUSF Configuration Factory
 */

package factory

import (
	"free5gc/lib/openapi/models"
)

type Configuration struct {
	Sbi *Sbi `yaml:"sbi,omitempty"`

	ServiceNameList []string `yaml:"serviceNameList,omitempty"`

	NrfUri string `yaml:"nrfUri,omitempty"`

	PlmnSupportList []models.PlmnId `yaml:"plmnSupportList,omitempty"`

	GroupId string `yaml:"groupId,omitempty"`
}
