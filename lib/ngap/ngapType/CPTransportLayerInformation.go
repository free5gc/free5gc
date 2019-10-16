//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	CPTransportLayerInformationPresentNothing int = iota /* No components present */
	CPTransportLayerInformationPresentEndpointIPAddress
	CPTransportLayerInformationPresentChoiceExtensions
)

type CPTransportLayerInformation struct {
	Present           int
	EndpointIPAddress *TransportLayerAddress
	ChoiceExtensions  *ProtocolIESingleContainerCPTransportLayerInformationExtIEs
}
