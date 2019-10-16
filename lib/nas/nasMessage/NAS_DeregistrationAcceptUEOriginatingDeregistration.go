//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type DeregistrationAcceptUEOriginatingDeregistration struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.DeregistrationAcceptMessageIdentity
}

func NewDeregistrationAcceptUEOriginatingDeregistration(iei uint8) (deregistrationAcceptUEOriginatingDeregistration *DeregistrationAcceptUEOriginatingDeregistration) {}

func (a *DeregistrationAcceptUEOriginatingDeregistration) EncodeDeregistrationAcceptUEOriginatingDeregistration(buffer *bytes.Buffer) {}

func (a *DeregistrationAcceptUEOriginatingDeregistration) DecodeDeregistrationAcceptUEOriginatingDeregistration(byteArray *[]byte) {}
