//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type DeregistrationRequestUETerminatedDeregistration struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.DeregistrationRequestMessageIdentity
	nasType.SpareHalfOctetAndDeregistrationType
	*nasType.Cause5GMM
	*nasType.T3346Value
}

func NewDeregistrationRequestUETerminatedDeregistration(iei uint8) (deregistrationRequestUETerminatedDeregistration *DeregistrationRequestUETerminatedDeregistration) {}

const (
	DeregistrationRequestUETerminatedDeregistrationCause5GMMType  uint8 = 0x58
	DeregistrationRequestUETerminatedDeregistrationT3346ValueType uint8 = 0x5F
)

func (a *DeregistrationRequestUETerminatedDeregistration) EncodeDeregistrationRequestUETerminatedDeregistration(buffer *bytes.Buffer) {}

func (a *DeregistrationRequestUETerminatedDeregistration) DecodeDeregistrationRequestUETerminatedDeregistration(byteArray *[]byte) {}
