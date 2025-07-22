package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type MultipleUnitUsage struct { /* Sequence Type */
	RatingGroup RatingGroupId `ber:"tagNum:0"`
	/* Sequence of = 35, FULL Name = struct MultipleUnitUsage__usedUnitContainers */
	/* UsedUnitContainer */
	UsedUnitContainers   []UsedUnitContainer  `ber:"tagNum:1,optional"`
	UPFID                *NetworkFunctionName `ber:"tagNum:2,optional"`
	MultihomedPDUAddress *PDUAddress          `ber:"tagNum:3,optional"`
}
