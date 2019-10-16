//go:binary-only-package

package ngapConvert

import (
	"free5gc/lib/aper"
	"free5gc/lib/ngap/logger"
	"free5gc/lib/ngap/ngapType"
	"net"
)

func IPAddressToString(ipAddr ngapType.TransportLayerAddress) (ipv4Addr, ipv6Addr string) {}

func IPAddressToNgap(ipv4Addr, ipv6Addr string) (ipAddr ngapType.TransportLayerAddress) {}
