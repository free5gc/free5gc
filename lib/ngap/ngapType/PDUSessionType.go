//go:binary-only-package

package ngapType

import "free5gc/lib/aper"

// Need to import "free5gc/lib/aper" if it uses "aper"

const (
	PDUSessionTypePresentIpv4         aper.Enumerated = 0
	PDUSessionTypePresentIpv6         aper.Enumerated = 1
	PDUSessionTypePresentIpv4v6       aper.Enumerated = 2
	PDUSessionTypePresentEthernet     aper.Enumerated = 3
	PDUSessionTypePresentUnstructured aper.Enumerated = 4
)

type PDUSessionType struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:4"`
}
