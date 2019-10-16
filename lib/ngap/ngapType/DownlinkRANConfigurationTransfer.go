//go:binary-only-package

package ngapType

// Need to import "free5gc/lib/aper" if it uses "aper"

type DownlinkRANConfigurationTransfer struct {
	ProtocolIEs ProtocolIEContainerDownlinkRANConfigurationTransferIEs
}
