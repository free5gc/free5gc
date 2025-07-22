package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type TAI struct { /* Sequence Type */
	PLMNId PLMNId `ber:"tagNum:0"`
	Tac    TAC    `ber:"tagNum:1"`
}
