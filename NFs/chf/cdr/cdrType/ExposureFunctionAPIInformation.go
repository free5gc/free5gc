package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type ExposureFunctionAPIInformation struct { /* Set Type */
	GroupIdentifier              *AddressString              `ber:"tagNum:0,optional"`
	APIDirection                 *APIDirection               `ber:"tagNum:1,optional"`
	APITargetNetworkFunction     *NetworkFunctionInformation `ber:"tagNum:2,optional"`
	APIResultCode                *APIResultCode              `ber:"tagNum:3,optional"`
	APIName                      asn.IA5String               `ber:"tagNum:4"`
	APIReference                 *asn.IA5String              `ber:"tagNum:5,optional"`
	APIContent                   *asn.OctetString            `ber:"tagNum:6,optional"`
	ExternalIndividualIdentifier *InvolvedParty              `ber:"tagNum:7,optional"`
	ExternalGroupIdentifier      *ExternalGroupIdentifier    `ber:"tagNum:8,optional"`
}
