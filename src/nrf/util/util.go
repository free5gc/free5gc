package util

import (
	"free5gc/lib/path_util"
)

// Path of HTTP2 key and log file
var (
	NRF_LOG_PATH = path_util.Gofree5gcPath("free5gc/nrfsslkey.log")
	NRF_PEM_PATH = path_util.Gofree5gcPath("free5gc/support/TLS/nrf.pem")
	NRF_KEY_PATH = path_util.Gofree5gcPath("free5gc/support/TLS/nrf.key")
)
