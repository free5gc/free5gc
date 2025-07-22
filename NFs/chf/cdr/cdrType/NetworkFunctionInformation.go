package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NetworkFunctionInformation struct { /* Sequence Type */
	NetworkFunctionality          NetworkFunctionality `ber:"tagNum:0"`
	NetworkFunctionName           *NetworkFunctionName `ber:"tagNum:1,optional"`
	NetworkFunctionIPv4Address    *IPAddress           `ber:"tagNum:2,optional"`
	NetworkFunctionPLMNIdentifier *PLMNId              `ber:"tagNum:3,optional"`
	NetworkFunctionIPv6Address    *IPAddress           `ber:"tagNum:4,optional"`
	NetworkFunctionFQDN           *NodeAddress         `ber:"tagNum:5,optional"`
}
