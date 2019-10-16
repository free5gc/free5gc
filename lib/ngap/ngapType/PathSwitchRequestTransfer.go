//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type PathSwitchRequestTransfer struct {
	DLNGUUPTNLInformation        UPTransportLayerInformation   `aper:"valueLB:0,valueUB:1"`
	DLNGUTNLInformationReused    *DLNGUTNLInformationReused    `aper:"optional"`
	UserPlaneSecurityInformation *UserPlaneSecurityInformation `aper:"valueExt,optional"`
	QosFlowAcceptedList          QosFlowAcceptedList
	IEExtensions                 *ProtocolExtensionContainerPathSwitchRequestTransferExtIEs `aper:"optional"`
}
