package factory

// UserPlaneInformation describe core network userplane information
type UserPlaneInformation struct {
	UPNodes map[string]UPNode `yaml:"up_nodes"`
	Links   []UPLink          `yaml:"links"`
}

// UPNode represent the user plane node
type UPNode struct {
	Type         string `yaml:"type"`
	NodeID       string `yaml:"node_id"`
	UPResourceIP string `yaml:"node"`
	ANIP         string `yaml:"an_ip"`
	Dnn          string `yaml:"dnn"`
}

type UPLink struct {
	A string `yaml:"A"`
	B string `yaml:"B"`
}
