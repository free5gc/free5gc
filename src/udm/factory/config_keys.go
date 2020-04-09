/*
 * UDM Configuration Factory
 */

package factory

type Keys struct {
	UdmProfileAHNPrivateKey string `yaml:"udmProfileAHNPrivateKey,omitempty"`
	UdmProfileAHNPublicKey  string `yaml:"udmProfileAHNPublicKey,omitempty"`
	UdmProfileBHNPrivateKey string `yaml:"udmProfileBHNPrivateKey,omitempty"`
	UdmProfileBHNPublicKey  string `yaml:"udmProfileBHNPublicKey,omitempty"`
}
