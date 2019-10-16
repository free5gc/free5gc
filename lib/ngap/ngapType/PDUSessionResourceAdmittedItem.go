//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

type PDUSessionResourceAdmittedItem struct {
	PDUSessionID                       PDUSessionID
	HandoverRequestAcknowledgeTransfer aper.OctetString
	IEExtensions                       *ProtocolExtensionContainerPDUSessionResourceAdmittedItemExtIEs `aper:"optional"`
}
