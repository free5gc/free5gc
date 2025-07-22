package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type Ncgi struct { /* Sequence Type */
	PlmnId   PLMNId   `ber:"tagNum:0"`
	NrCellId NrCellId `ber:"tagNum:1"`
	Nid      *Nid     `ber:"tagNum:2,optional"`
}
