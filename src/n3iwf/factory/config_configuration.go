/*
 * N3IWF Configuration Factory
 */

package factory

import (
	"free5gc/src/n3iwf/n3iwf_context"
)

type Configuration struct {
	N3IWFInfo            n3iwf_context.N3IWFNFInfo `yaml:"N3IWFInformation"`
	AMFAddress           []ConfigAMFAddr           `yaml:"AMFAddress"`
	IKEBindAddr          string                    `yaml:"IKEBindAddress"`
	IPSecGatewayAddr     string                    `yaml:"IPSecInterfaceAddress"`
	GTPBindAddr          string                    `yaml:"GTPBindAddress"`
	TCPPort              uint16                    `yaml:"NASTCPPort"`
	FQDN                 string                    `yaml:"FQDN"`                 // e.g. n3iwf.free5gc.org
	PrivateKey           string                    `yaml:"PrivateKey"`           // file path
	CertificateAuthority string                    `yaml:"CertificateAuthority"` // file path
	Certificate          string                    `yaml:"Certificate"`          // file path
	UEIPAddressRange     string                    `yaml:"UEIPAddressRange"`     // e.g. 10.0.1.0/24
	InterfaceMark        uint32                    `yaml:"IPSecInterfaceMark"`   // must != 0, if not specified, random one
}

type ConfigAMFAddr struct {
	NetworkAddress string `yaml:"IP"`
	Port           int    `yaml:"Port,omitempty"`
}
