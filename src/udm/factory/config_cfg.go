/*
 * UDM Configuration Factory
 */

package factory

type Udrclient struct {
	Scheme   string `yaml:"scheme"`
	Ipv4Adrr string `yaml:"ipv4Adrr,omitempty"`
	Port     int    `yaml:"port,omitempty"`
}
