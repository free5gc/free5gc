/*
 * AMF Configuration Factory
 */

package factory

import (
	"free5gc/lib/openapi/models"
	"free5gc/src/amf/amf_context"
)

type Configuration struct {
	AmfName string `yaml:"amfName,omitempty"`

	NgapIpList []string `yaml:"ngapIpList,omitempty"`

	Sbi *Sbi `yaml:"sbi,omitempty"`

	ServiceNameList []string `yaml:"serviceNameList,omitempty"`

	ServedGumaiList []models.Guami `yaml:"servedGuamiList,omitempty"`

	SupportTAIList []models.Tai `yaml:"supportTaiList,omitempty"`

	PlmnSupportList []amf_context.PlmnSupportItem `yaml:"plmnSupportList,omitempty"`

	SupportDnnList []string `yaml:"supportDnnList,omitempty"`

	NrfUri string `yaml:"nrfUri,omitempty"`

	Security *Security `yaml:"security,omitempty"`

	NetworkName amf_context.NetworkName `yaml:"networkName,omitempty"`

	T3502 int `yaml:"t3502,omitempty"`

	T3512 int `yaml:"t3512,omitempty"`

	Non3gppDeregistrationTimer int `yaml:"mon3gppDeregistrationTimer,omitempty"`
}
