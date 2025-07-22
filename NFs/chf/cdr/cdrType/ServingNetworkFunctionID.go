package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type ServingNetworkFunctionID struct { /* Sequence Type */
	ServingNetworkFunctionInformation NetworkFunctionInformation `ber:"tagNum:0"`
	AMFIdentifier                     *AMFID                     `ber:"tagNum:1,optional"`
}
