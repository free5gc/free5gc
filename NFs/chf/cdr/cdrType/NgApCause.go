package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NgApCause struct { /* Sequence Type */
	Group int64 `ber:"tagNum:0"`
	Value int64 `ber:"tagNum:1"`
}
