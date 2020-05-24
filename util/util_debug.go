//+build debug

package util

import (
	"free5gc/lib/path_util"
)

// Path of HTTP2 key and log file

var NrfLogPath = path_util.Gofree5gcPath("free5gc/nrfsslkey.log")
var NrfPemPath = path_util.Gofree5gcPath("free5gc/support/TLS/_debug.pem")
var NrfKeyPath = path_util.Gofree5gcPath("free5gc/support/TLS/_debug.key")
