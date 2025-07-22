package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type Throughput struct { /* Sequence Type */
	GuaranteedThpt Bitrate `ber:"tagNum:0"`
	MaximumThpt    Bitrate `ber:"tagNum:1"`
}
