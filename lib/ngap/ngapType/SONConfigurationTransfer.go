//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type SONConfigurationTransfer struct {
	TargetRANNodeID        TargetRANNodeID                                           `aper:"valueExt"`
	SourceRANNodeID        SourceRANNodeID                                           `aper:"valueExt"`
	SONInformation         SONInformation                                            `aper:"valueLB:0,valueUB:2"`
	XnTNLConfigurationInfo XnTNLConfigurationInfo                                    `aper:"valueExt"`
	IEExtensions           *ProtocolExtensionContainerSONConfigurationTransferExtIEs `aper:"optional"`
}
