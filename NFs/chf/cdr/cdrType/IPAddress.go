package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	IPAddressPresentNothing int = iota /* No components present */
	IPAddressPresentIPBinV4Address
	IPAddressPresentIPBinV6Address
	IPAddressPresentIPTextV4Address
	IPAddressPresentIPTextV6Address
	IPAddressPresentIPBinV6AddressWithPrefix
)

type IPAddress struct {
	Present                  int                             /* Choice Type */
	IPBinV4Address           *IPBinV4Address                 `ber:"tagNum:0"`
	IPBinV6Address           *IPBinV6Address                 `ber:"tagNum:1"`
	IPTextV4Address          *asn.IA5String                  `ber:"tagNum:2"`
	IPTextV6Address          *asn.IA5String                  `ber:"tagNum:3"`
	IPBinV6AddressWithPrefix *IPBinV6AddressWithPrefixLength `ber:"tagNum:4"`
}
