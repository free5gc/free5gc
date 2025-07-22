package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type MAPDUSteeringMode struct { /* Sequence Type */
	SteerModeValue *SteerModeValue `ber:"tagNum:0,optional"`
	Active         *AccessType     `ber:"tagNum:1,optional"`
	Standby        *AccessType     `ber:"tagNum:2,optional"`
	ThreegLoad     *int64          `ber:"tagNum:3,optional"`
	PrioAcc        *AccessType     `ber:"tagNum:4,optional"`
}
