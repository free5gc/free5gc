package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type SingleNSSAI struct { /* Sequence Type */
	SST SliceServiceType     `ber:"tagNum:0"`
	SD  *SliceDifferentiator `ber:"tagNum:1,optional"`
}
