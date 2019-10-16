//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type DeregistrationAcceptUETerminatedDeregistration struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.DeregistrationAcceptMessageIdentity
}

func NewDeregistrationAcceptUETerminatedDeregistration(iei uint8) (deregistrationAcceptUETerminatedDeregistration *DeregistrationAcceptUETerminatedDeregistration) {}

func (a *DeregistrationAcceptUETerminatedDeregistration) EncodeDeregistrationAcceptUETerminatedDeregistration(buffer *bytes.Buffer) {}

func (a *DeregistrationAcceptUETerminatedDeregistration) DecodeDeregistrationAcceptUETerminatedDeregistration(byteArray *[]byte) {}
