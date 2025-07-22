package cdrType

import "github.com/free5gc/chf/cdr/asn"

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type PDUSessionChargingInformation struct { /* Set Type */
	PDUSessionChargingID      ChargingID                 `ber:"tagNum:0"`
	UserIdentifier            *InvolvedParty             `ber:"tagNum:1,optional"`
	UserEquipmentInfo         *SubscriberEquipmentNumber `ber:"tagNum:2,optional"`
	UserLocationInformation   *UserLocationInformation   `ber:"tagNum:3,optional"`
	UserRoamerInOut           *RoamerInOut               `ber:"tagNum:4,optional"`
	PresenceReportingAreaInfo *PresenceReportingAreaInfo `ber:"tagNum:5,optional"`
	PDUSessionId              PDUSessionId               `ber:"tagNum:6"`
	NetworkSliceInstanceID    *SingleNSSAI               `ber:"tagNum:7,optional"`
	PDUType                   *PDUSessionType            `ber:"tagNum:8,optional"`
	SSCMode                   *SSCMode                   `ber:"tagNum:9,optional"`
	SUPIPLMNIdentifier        *PLMNId                    `ber:"tagNum:10,optional"`
	/* Sequence of = 35, FULL Name = struct PDUSessionChargingInformation__servingNetworkFunctionID */
	/* ServingNetworkFunctionID */
	ServingNetworkFunctionID  []ServingNetworkFunctionID `ber:"tagNum:11,optional"`
	RATType                   *RATType                   `ber:"tagNum:12,optional"`
	DataNetworkNameIdentifier *DataNetworkNameIdentifier `ber:"tagNum:13,optional"`
	PDUAddress                *PDUAddress                `ber:"tagNum:14,optional"`
	AuthorizedQoSInformation  *AuthorizedQoSInformation  `ber:"tagNum:15,optional"`
	UETimeZone                *MSTimeZone                `ber:"tagNum:16,optional"`
	PDUSessionstartTime       *TimeStamp                 `ber:"tagNum:17,optional"`
	PDUSessionstopTime        *TimeStamp                 `ber:"tagNum:18,optional"`
	Diagnostics               *Diagnostics               `ber:"tagNum:19,optional"`
	ChargingCharacteristics   *ChargingCharacteristics   `ber:"tagNum:20,optional"`
	ChChSelectionMode         *ChChSelectionMode         `ber:"tagNum:21,optional"`
	ThreeGPPPSDataOffStatus   *ThreeGPPPSDataOffStatus   `ber:"tagNum:22,optional"`
	/* Sequence of = 35, FULL Name = struct PDUSessionChargingInformation__rANSecondaryRATUsageReport */
	/* NGRANSecondaryRATUsageReport */
	RANSecondaryRATUsageReport           []NGRANSecondaryRATUsageReport     `ber:"tagNum:23,optional"`
	SubscribedQoSInformation             *SubscribedQoSInformation          `ber:"tagNum:24,optional"`
	AuthorizedSessionAMBR                *SessionAMBR                       `ber:"tagNum:25,optional"`
	SubscribedSessionAMBR                *SessionAMBR                       `ber:"tagNum:26,optional"`
	ServingCNPLMNID                      *PLMNId                            `ber:"tagNum:27,optional"`
	SUPIunauthenticatedFlag              *asn.NULL                          `ber:"tagNum:28,optional"`
	DnnSelectionMode                     *DNNSelectionMode                  `ber:"tagNum:29,optional"`
	HomeProvidedChargingID               *ChargingID                        `ber:"tagNum:30,optional"`
	MAPDUNonThreeGPPUserLocationInfo     *UserLocationInformation           `ber:"tagNum:31,optional"`
	MAPDUNonThreeGPPRATType              *RATType                           `ber:"tagNum:32,optional"`
	MAPDUSessionInformation              *MAPDUSessionInformation           `ber:"tagNum:33,optional"`
	EnhancedDiagnostics                  *EnhancedDiagnostics5G             `ber:"tagNum:34,optional"`
	UserLocationInformationASN1          *UserLocationInformationStructured `ber:"tagNum:35,optional"`
	MAPDUNonThreeGPPUserLocationInfoASN1 *UserLocationInformationStructured `ber:"tagNum:36,optional"`
}
