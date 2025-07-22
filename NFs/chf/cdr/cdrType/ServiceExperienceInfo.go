package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type ServiceExperienceInfo struct { /* Sequence Type */
	SvcExprc         *SvcExperience             `ber:"tagNum:0,optional"`
	SvcExprcVariance *int64                     `ber:"tagNum:1,optional"`
	Snssai           *SingleNSSAI               `ber:"tagNum:2,optional"`
	AppId            *asn.OctetString           `ber:"tagNum:3,optional"`
	Confidence       *int64                     `ber:"tagNum:4,optional"`
	Dnn              *DataNetworkNameIdentifier `ber:"tagNum:5,optional"`
	NetworkArea      *NetworkAreaInfo           `ber:"tagNum:6,optional"`
	NsiId            *asn.OctetString           `ber:"tagNum:7,optional"`
	Ratio            *int64                     `ber:"tagNum:8,optional"`
}
