package factory

type Configuration struct {
	UdmName string `yaml:"udmName,omitempty"`

	Sbi *Sbi `yaml:"sbi,omitempty"`

	ServiceNameList []string `yaml:"serviceNameList,omitempty"`

	Udrclient *Udrclient `yaml:"udrclient,omitempty"`

	Nrfclient *Nrfclient `yaml:"nrfclient,omitempty"`
}
