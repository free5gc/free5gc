//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type AMFTNLAssociationToUpdateItem struct {
	AMFTNLAssociationAddress CPTransportLayerInformation                                    `aper:"valueLB:0,valueUB:1"`
	TNLAssociationUsage      *TNLAssociationUsage                                           `aper:"optional"`
	TNLAddressWeightFactor   *TNLAddressWeightFactor                                        `aper:"optional"`
	IEExtensions             *ProtocolExtensionContainerAMFTNLAssociationToUpdateItemExtIEs `aper:"optional"`
}
