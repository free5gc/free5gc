/*
 * UDR Configuration Factory
 */

package factory

type Configuration struct {
	Sbi *Sbi `yaml:"sbi"`

	Mongodb *Mongodb `yaml:"mongodb"`

	NrfUri string `yaml:"nrfUri"`
}
