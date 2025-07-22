package context

import "net"

// SnssaiSmfInfo records the SMF S-NSSAI related information
type SnssaiSmfInfo struct {
	Snssai   SNssai
	DnnInfos map[string]*SnssaiSmfDnnInfo
}

// SnssaiSmfDnnInfo records the SMF per S-NSSAI DNN information
type SnssaiSmfDnnInfo struct {
	DNS   DNS
	PCSCF PCSCF
}

type DNS struct {
	IPv4Addr net.IP
	IPv6Addr net.IP
}

type PCSCF struct {
	IPv4Addr net.IP
}
