//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type QosFlowLevelQosParameters struct {
	QosCharacteristics             QosCharacteristics                                         `aper:"valueLB:0,valueUB:2"`
	AllocationAndRetentionPriority AllocationAndRetentionPriority                             `aper:"valueExt"`
	GBRQosInformation              *GBRQosInformation                                         `aper:"valueExt,optional"`
	ReflectiveQosAttribute         *ReflectiveQosAttribute                                    `aper:"optional"`
	AdditionalQosFlowInformation   *AdditionalQosFlowInformation                              `aper:"optional"`
	IEExtensions                   *ProtocolExtensionContainerQosFlowLevelQosParametersExtIEs `aper:"optional"`
}
