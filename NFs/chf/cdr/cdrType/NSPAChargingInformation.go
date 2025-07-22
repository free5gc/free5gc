package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type NSPAChargingInformation struct { /* Set Type */
	SingelNSSAI SingleNSSAI `ber:"tagNum:0"`
}
