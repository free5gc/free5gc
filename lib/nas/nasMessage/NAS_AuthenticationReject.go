//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type AuthenticationReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.AuthenticationRejectMessageIdentity
	*nasType.EAPMessage
}

func NewAuthenticationReject(iei uint8) (authenticationReject *AuthenticationReject) {}

const (
	AuthenticationRejectEAPMessageType uint8 = 0x78
)

func (a *AuthenticationReject) EncodeAuthenticationReject(buffer *bytes.Buffer) {}

func (a *AuthenticationReject) DecodeAuthenticationReject(byteArray *[]byte) {}
