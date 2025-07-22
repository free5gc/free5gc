package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NetworkAreaInfo struct { /* Sequence Type */
	/* Sequence of = 35, FULL Name = struct NetworkAreaInfo__ecgis */
	/* Ecgi */
	Ecgis []Ecgi `ber:"tagNum:0,optional"`
	/* Sequence of = 35, FULL Name = struct NetworkAreaInfo__ncgis */
	/* Ncgi */
	Ncgis []Ncgi `ber:"tagNum:1,optional"`
	/* Sequence of = 35, FULL Name = struct NetworkAreaInfo__gRanNodeIds */
	/* GlobalRanNodeId */
	GRanNodeIds []GlobalRanNodeId `ber:"tagNum:2,optional"`
	/* Sequence of = 35, FULL Name = struct NetworkAreaInfo__tais */
	/* TAI */
	Tais []TAI `ber:"tagNum:3,optional"`
}
