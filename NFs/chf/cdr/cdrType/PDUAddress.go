package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type PDUAddress struct { /* Sequence Type */
	PDUIPv4Address           *IPAddress          `ber:"tagNum:0,optional"`
	PDUIPv6AddresswithPrefix *IPAddress          `ber:"tagNum:1,optional"`
	IPV4dynamicAddressFlag   *DynamicAddressFlag `ber:"tagNum:2,optional"`
	IPV6dynamicPrefixFlag    *DynamicAddressFlag `ber:"tagNum:3,optional"`
	/* Sequence of = 35, FULL Name = struct PDUAddress__additionalPDUIPv6Prefixes */
	/* IPAddress */
	AdditionalPDUIPv6Prefixes []IPAddress `ber:"tagNum:4,optional"`
}
