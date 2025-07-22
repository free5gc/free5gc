package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	IPBinaryAddressPresentNothing int = iota /* No components present */
	IPBinaryAddressPresentIPBinV4Address
	IPBinaryAddressPresentIPBinV6Address
)

type IPBinaryAddress struct {
	Present        int             /* Choice Type */
	IPBinV4Address *IPBinV4Address `ber:"tagNum:0"`
	IPBinV6Address *IPBinV6AddressWithOrWithoutPrefixLength
}
