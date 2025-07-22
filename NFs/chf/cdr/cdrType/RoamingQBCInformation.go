package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type RoamingQBCInformation struct { /* Set Type */
	/* Sequence of = 35, FULL Name = struct RoamingQBCInformation__multipleQFIcontainer */
	/* MultipleQFIContainer */
	MultipleQFIcontainer   []MultipleQFIContainer  `ber:"tagNum:0,optional"`
	UPFID                  *NetworkFunctionName    `ber:"tagNum:1,optional"`
	RoamingChargingProfile *RoamingChargingProfile `ber:"tagNum:2,optional"`
}
