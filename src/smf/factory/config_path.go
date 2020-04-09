package factory

type Path struct {
	DestinationIP string `yaml:"DestinationIP,omitempty"`

	DestinationPort string `yaml:"DestinationPort,omitempty"`

	UPF []string `yaml:"UPF,omitempty"`
}
