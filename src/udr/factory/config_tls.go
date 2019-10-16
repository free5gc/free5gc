/*
 * UDR Configuration Factory
 */

package factory

type Tls struct {
	Log string `yaml:"log"`

	Pem string `yaml:"pem"`

	Key string `yaml:"key"`
}
