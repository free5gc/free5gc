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

	ServiceList []Service `yaml:"serviceList,omitempty"`
}

type Service struct {
	ServiceName string `yaml:"serviceName"`
	SuppFeat    string `yaml:"suppFeat,omitempty"`
}
