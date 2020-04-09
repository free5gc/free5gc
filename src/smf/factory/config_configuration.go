/*
 * AMF Configuration Factory
 */

package factory

import "free5gc/lib/openapi/models"

type Configuration struct {
	SmfName string `yaml:"smfName,omitempty"`

	Sbi *Sbi `yaml:"sbi,omitempty"`

	PFCP *PFCP `yaml:"pfcp,omitempty"`

	NrfUri string `yaml:"nrfUri,omitempty"`

	UserPlaneInformation UserPlaneInformation `yaml:"userplane_information"`

	UESubnet string `yaml:"ue_subnet"`

	ServiceNameList []string `yaml:"serviceNameList,omitempty"`

	SNssaiInfo []models.SnssaiSmfInfoItem `yaml:"snssai_info,omitempty"`

	ULCL bool `yaml:"ulcl,omitempty"`
}
