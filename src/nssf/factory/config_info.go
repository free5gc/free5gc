/*
 * NSSF Configuration Factory
 */

package factory

type Info struct {
	Version string `yaml:"version"`

	Description string `yaml:"description,omitempty"`
}
