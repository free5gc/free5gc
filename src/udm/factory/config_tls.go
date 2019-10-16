/*
 * UDM Configuration Factory
 */

package factory

type Tls struct {
	Log string `yaml:"log,omitempty"`

	Pem string `yaml:"pem,omitempty"`

	Key string `yaml:"key,omitempty"`
}
