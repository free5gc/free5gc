/*
 * UDM Configuration Factory
 */

package factory

type Sbi struct {
	Scheme   string `yaml:"scheme"`
	IPv4Addr string `yaml:"ipv4Addr,omitempty"`
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	Port int  `yaml:"port,omitempty"`
	Tls  *Tls `yaml:"tls,omitempty"`
}
