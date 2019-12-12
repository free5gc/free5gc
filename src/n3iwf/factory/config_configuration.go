/*
 * N3IWF Configuration Factory
 */

package factory

import (
	"free5gc/src/n3iwf/n3iwf_context"
)

type Configuration struct {
	N3IWFInfo  n3iwf_context.N3IWFNFInfo `yaml:"N3IWFInformation"`
	AMFAddress []ConfigAMFAddr           `yaml:"AMFAddress"`
}

type ConfigAMFAddr struct {
	NetworkAddress string `yaml:"IP"`
	Port           int    `yaml:"Port,omitempty"`
}
