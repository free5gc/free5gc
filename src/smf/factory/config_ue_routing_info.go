/*
 * SMF Routing Configuration Factory
 */

package factory

type UERoutingInfo struct {
	SUPI string `yaml:"SUPI,omitempty"`

	AN string `yaml:"AN,omitempty"`

	PathList []Path `yaml:"PathList,omitempty"`
}
