package cdrType

// Need to import "gofree5gc/lib/aper" if it uses "aper"

type MultipleQFIContainer struct { /* Sequence Type */
	QosFlowId *QoSFlowId `ber:"tagNum:0,optional"`
	/* Sequence of = 35, FULL Name = struct MultipleQFIContainer__triggers */
	/* Trigger */
	Triggers                  []Trigger                  `ber:"tagNum:1,optional"`
	TriggerTimeStamp          *TimeStamp                 `ber:"tagNum:2,optional"`
	DataTotalVolume           *DataVolumeOctets          `ber:"tagNum:3,optional"`
	DataVolumeUplink          *DataVolumeOctets          `ber:"tagNum:4,optional"`
	DataVolumeDownlink        *DataVolumeOctets          `ber:"tagNum:5,optional"`
	LocalSequenceNumber       *LocalSequenceNumber       `ber:"tagNum:6,optional"`
	TimeOfFirstUsage          *TimeStamp                 `ber:"tagNum:8,optional"`
	TimeOfLastUsage           *TimeStamp                 `ber:"tagNum:9,optional"`
	QoSInformation            *FiveGQoSInformation       `ber:"tagNum:10,optional"`
	UserLocationInformation   *UserLocationInformation   `ber:"tagNum:11,optional"`
	UETimeZone                *MSTimeZone                `ber:"tagNum:12,optional"`
	PresenceReportingAreaInfo *PresenceReportingAreaInfo `ber:"tagNum:13,optional"`
	RATType                   *RATType                   `ber:"tagNum:14,optional"`
	ReportTime                TimeStamp                  `ber:"tagNum:15"`
	/* Sequence of = 35, FULL Name = struct MultipleQFIContainer__servingNetworkFunctionID */
	/* ServingNetworkFunctionID */
	ServingNetworkFunctionID    []ServingNetworkFunctionID         `ber:"tagNum:16,optional"`
	ThreeGPPPSDataOffStatus     *ThreeGPPPSDataOffStatus           `ber:"tagNum:17,optional"`
	ThreeGPPChargingID          *ChargingID                        `ber:"tagNum:18,optional"`
	Diagnostics                 *Diagnostics                       `ber:"tagNum:19,optional"`
	ExtensionDiagnostics        *EnhancedDiagnostics               `ber:"tagNum:20,optional"`
	QoSCharacteristics          *QoSCharacteristics                `ber:"tagNum:21,optional"`
	Time                        *CallDuration                      `ber:"tagNum:22,optional"`
	UserLocationInformationASN1 *UserLocationInformationStructured `ber:"tagNum:23,optional"`
}
