package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type IPBinV6AddressWithPrefixLength struct { /* Sequence Type */
	IPBinV6Address         IPBinV6Address
	PDPAddressPrefixLength *PDPAddressPrefixLength `ber:"optional,default:64"`
}
