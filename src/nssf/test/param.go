/*
 * NSSF Testing Utility
 */

package test

import (
	"flag"

	"free5gc/lib/path_util"
	. "free5gc/src/nssf/plugin"
)

var (
	ConfigFileFromArgs string
	DefaultConfigFile  string = path_util.Gofree5gcPath("free5gc/src/nssf/test/conf/test_nssf_config.yaml")
)

type TestingUtil struct {
	ConfigFile string
}

type TestingNsselection struct {
	ConfigFile string

	GenerateNonRoamingQueryParameter func() NsselectionQueryParameter

	GenerateRoamingQueryParameter func() NsselectionQueryParameter
}

type TestingNssaiavailability struct {
	ConfigFile string

	NfId string

	SubscriptionId string

	NfNssaiAvailabilityUri string
}

func init() {
	flag.StringVar(&ConfigFileFromArgs, "config-file", DefaultConfigFile, "Configuration file")
	flag.Parse()
}
