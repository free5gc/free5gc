/*
 * PCF Configuration Factory
 */

package factory

type Configuration struct {
	PcfName string `yaml:"pcfName,omitempty"`

	Sbi *Sbi `yaml:"sbi,omitempty"`

	TimeFormat string `yaml:"timeFormat,omitempty"`

	DefaultBdtRefId string `yaml:"defaultBdtRefId,omitempty"`

	NrfUri string `yaml:"nrfUri,omitempty"`

	ServiceNameList []string `yaml:"serviceNameList,omitempty"`
}
