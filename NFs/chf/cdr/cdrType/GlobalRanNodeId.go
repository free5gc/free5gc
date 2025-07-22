package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type GlobalRanNodeId struct { /* Sequence Type */
	PLMNId  *PLMNId  `ber:"tagNum:0,optional"`
	N3IwfId *N3IwFId `ber:"tagNum:1,optional"`
	GNbId   *GNbId   `ber:"tagNum:2,optional"`
	NgeNbId *NgeNbId `ber:"tagNum:3,optional"`
	WagfId  *WAgfId  `ber:"tagNum:4,optional"`
	TngfId  *TngfId  `ber:"tagNum:5,optional"`
	Nid     *Nid     `ber:"tagNum:6,optional"`
	ENbId   *ENbId   `ber:"tagNum:7,optional"`
}
