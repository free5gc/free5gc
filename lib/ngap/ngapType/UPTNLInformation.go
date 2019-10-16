//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	UPTNLInformationPresentNothing int = iota /* No components present */
	UPTNLInformationPresentSingleTNLInformation
	UPTNLInformationPresentMultipleTNLInformation
	UPTNLInformationPresentChoiceExtensions
)

type UPTNLInformation struct {
	Present                int
	SingleTNLInformation   *SingleTNLInformation   `aper:"valueExt"`
	MultipleTNLInformation *MultipleTNLInformation `aper:"valueExt"`
	ChoiceExtensions       *ProtocolIESingleContainerUPTNLInformationExtIEs
}
