//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type RegistrationReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.RegistrationRejectMessageIdentity
	nasType.Cause5GMM
	*nasType.T3346Value
	*nasType.T3502Value
	*nasType.EAPMessage
}

func NewRegistrationReject(iei uint8) (registrationReject *RegistrationReject) {}

const (
	RegistrationRejectT3346ValueType uint8 = 0x5F
	RegistrationRejectT3502ValueType uint8 = 0x16
	RegistrationRejectEAPMessageType uint8 = 0x78
)

func (a *RegistrationReject) EncodeRegistrationReject(buffer *bytes.Buffer) {}

func (a *RegistrationReject) DecodeRegistrationReject(byteArray *[]byte) {}
