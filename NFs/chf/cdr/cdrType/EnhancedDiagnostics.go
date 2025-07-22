package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type EnhancedDiagnostics struct { /* Sequence Type */
	/* Sequence of = 35, FULL Name = struct EnhancedDiagnostics__rANNASCause */
	/* RANNASCause */
	RANNASCause []RANNASCause `ber:"tagNum:0"`
}
