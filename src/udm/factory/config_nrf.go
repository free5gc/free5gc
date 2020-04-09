/*
 * UDM Configuration Factory
 */

package factory

type Nrfclient struct {
	Scheme   string `yaml:"scheme"`
	Ipv4Addr string `yaml:"ipv4Addr,omitempty"`
	Port     int    `yaml:"port,omitempty"`
}
