/*
 * UDR Configuration Factory
 */

package factory

type Config struct {
	Info *Info `yaml:"info"`

	Configuration *Configuration `yaml:"configuration"`
}
