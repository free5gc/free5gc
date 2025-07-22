package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type PDUContainerInformation struct { /* Sequence Type */
	ChargingRuleBaseName               *ChargingRuleBaseName      `ber:"tagNum:0,optional"`
	TimeOfFirstUsage                   *TimeStamp                 `ber:"tagNum:2,optional"`
	TimeOfLastUsage                    *TimeStamp                 `ber:"tagNum:3,optional"`
	QoSInformation                     *FiveGQoSInformation       `ber:"tagNum:4,optional"`
	UserLocationInformation            *UserLocationInformation   `ber:"tagNum:5,optional"`
	PresenceReportingAreaInfo          *PresenceReportingAreaInfo `ber:"tagNum:6,optional"`
	RATType                            *RATType                   `ber:"tagNum:7,optional"`
	SponsorIdentity                    *asn.OctetString           `ber:"tagNum:8,optional"`
	ApplicationServiceProviderIdentity *asn.OctetString           `ber:"tagNum:9,optional"`
	/* Sequence of = 35, FULL Name = struct PDUContainerInformation__servingNetworkFunctionID */
	/* ServingNetworkFunctionID */
	ServingNetworkFunctionID    []ServingNetworkFunctionID         `ber:"tagNum:10,optional"`
	UETimeZone                  *MSTimeZone                        `ber:"tagNum:11,optional"`
	ThreeGPPPSDataOffStatus     *ThreeGPPPSDataOffStatus           `ber:"tagNum:12,optional"`
	QoSCharacteristics          *QoSCharacteristics                `ber:"tagNum:13,optional"`
	AfChargingIdentifier        *ChargingID                        `ber:"tagNum:14,optional"`
	AfChargingIdString          *AFChargingID                      `ber:"tagNum:15,optional"`
	MAPDUSteeringFunctionality  *MAPDUSteeringFunctionality        `ber:"tagNum:16,optional"`
	MAPDUSteeringMode           *MAPDUSteeringMode                 `ber:"tagNum:17,optional"`
	UserLocationInformationASN1 *UserLocationInformationStructured `ber:"tagNum:18,optional"`
	/* Sequence of = 35, FULL Name = struct PDUContainerInformation__listOfPresenceReportingAreaInformation */
	/* PresenceReportingAreaInfo */
	ListOfPresenceReportingAreaInformation []PresenceReportingAreaInfo `ber:"tagNum:19,optional"`
}
