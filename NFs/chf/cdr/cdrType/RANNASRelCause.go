package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type RANNASRelCause struct { /* Sequence Type */
	NgApCause    *NgApCause    `ber:"tagNum:0,optional"`
	FivegMmCause *FiveGMmCause `ber:"tagNum:1,optional"`
	FivegSmCause *FiveGSmCause `ber:"tagNum:2,optional"`
	EpsCause     *RANNASCause  `ber:"tagNum:3,optional"`
}
