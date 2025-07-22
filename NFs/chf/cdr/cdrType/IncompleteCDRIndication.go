package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type IncompleteCDRIndication struct { /* Sequence Type */
	InitialLost     *bool `ber:"tagNum:0,optional"`
	UpdateLost      *bool `ber:"tagNum:1,optional"`
	TerminationLost *bool `ber:"tagNum:2,optional"`
}
