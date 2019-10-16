/*
 * UDR Configuration Factory
 */

package factory

type Mongodb struct {
	Name string `yaml:"name"`

	Url string `yaml:"url"`
}
