//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type RegistrationComplete struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.RegistrationCompleteMessageIdentity
	*nasType.SORTransparentContainer
}

func NewRegistrationComplete(iei uint8) (registrationComplete *RegistrationComplete) {}

const (
	RegistrationCompleteSORTransparentContainerType uint8 = 0x73
)

func (a *RegistrationComplete) EncodeRegistrationComplete(buffer *bytes.Buffer) {}

func (a *RegistrationComplete) DecodeRegistrationComplete(byteArray *[]byte) {}
