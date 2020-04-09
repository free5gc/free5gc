/*
 * SMF Routing Configuration Factory
 */

package factory

type RoutingConfig struct {
	Info *Info `yaml:"info"`

	UERoutingInfo []*UERoutingInfo `yaml:"ueRoutingInfo"`
}
