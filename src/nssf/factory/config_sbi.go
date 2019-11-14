/*
 * NSSF Configuration Factory
 */

package factory

import (
	"free5gc/lib/openapi/models"
)

type Sbi struct {
	Scheme models.UriScheme `yaml:"scheme"`

	// Currently only support IPv4 and thus `Ipv4Addr` field shall not be empty
	Ipv4Addr string `yaml:"ipv4Addr"`

	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`

	Port int `yaml:"port"`
}
