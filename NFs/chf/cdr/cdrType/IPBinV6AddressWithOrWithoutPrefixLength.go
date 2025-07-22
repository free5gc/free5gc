package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

const (
	IPBinV6AddressWithOrWithoutPrefixLengthPresentNothing int = iota /* No components present */
	IPBinV6AddressWithOrWithoutPrefixLengthPresentIPBinV6Address
	IPBinV6AddressWithOrWithoutPrefixLengthPresentIPBinV6AddressWithPrefix
)

type IPBinV6AddressWithOrWithoutPrefixLength struct {
	Present                  int                             /* Choice Type */
	IPBinV6Address           *IPBinV6Address                 `ber:"tagNum:1"`
	IPBinV6AddressWithPrefix *IPBinV6AddressWithPrefixLength `ber:"tagNum:4"`
}
