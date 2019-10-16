//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type AuthenticationResponse struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.AuthenticationResponseMessageIdentity
	*nasType.AuthenticationResponseParameter
	*nasType.EAPMessage
}

func NewAuthenticationResponse(iei uint8) (authenticationResponse *AuthenticationResponse) {}

const (
	AuthenticationResponseAuthenticationResponseParameterType uint8 = 0x2D
	AuthenticationResponseEAPMessageType                      uint8 = 0x78
)

func (a *AuthenticationResponse) EncodeAuthenticationResponse(buffer *bytes.Buffer) {}

func (a *AuthenticationResponse) DecodeAuthenticationResponse(byteArray *[]byte) {}
