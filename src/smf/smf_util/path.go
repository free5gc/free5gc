package smf_util

import (
	"free5gc/lib/path_util"
)

var SmfLogPath = path_util.Gofree5gcPath("free5gc/smfsslkey.log")
var SmfPemPath = path_util.Gofree5gcPath("free5gc/support/TLS/smf.pem")
var SmfKeyPath = path_util.Gofree5gcPath("free5gc/support/TLS/smf.key")
var DefaultSmfConfigPath = path_util.Gofree5gcPath("free5gc/config/smfcfg.conf")
