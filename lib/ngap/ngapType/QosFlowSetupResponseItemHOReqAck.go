//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type QosFlowSetupResponseItemHOReqAck struct {
	QosFlowIdentifier      QosFlowIdentifier
	DataForwardingAccepted *DataForwardingAccepted                                           `aper:"optional"`
	IEExtensions           *ProtocolExtensionContainerQosFlowSetupResponseItemHOReqAckExtIEs `aper:"optional"`
}
