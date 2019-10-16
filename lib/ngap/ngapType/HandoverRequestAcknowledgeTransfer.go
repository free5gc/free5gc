//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type HandoverRequestAcknowledgeTransfer struct {
	DLNGUUPTNLInformation         UPTransportLayerInformation  `aper:"valueLB:0,valueUB:1"`
	DLForwardingUPTNLInformation  *UPTransportLayerInformation `aper:"valueLB:0,valueUB:1,optional"`
	SecurityResult                *SecurityResult              `aper:"valueExt,optional"`
	QosFlowSetupResponseList      QosFlowSetupResponseListHOReqAck
	QosFlowFailedToSetupList      *QosFlowList                                                        `aper:"optional"`
	DataForwardingResponseDRBList *DataForwardingResponseDRBList                                      `aper:"optional"`
	IEExtensions                  *ProtocolExtensionContainerHandoverRequestAcknowledgeTransferExtIEs `aper:"optional"`
}
