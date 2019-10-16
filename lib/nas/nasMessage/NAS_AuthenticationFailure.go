//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type AuthenticationFailure struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.AuthenticationFailureMessageIdentity
	nasType.Cause5GMM
	*nasType.AuthenticationFailureParameter
}

func NewAuthenticationFailure(iei uint8) (authenticationFailure *AuthenticationFailure) {}

const (
	AuthenticationFailureAuthenticationFailureParameterType uint8 = 0x30
)

func (a *AuthenticationFailure) EncodeAuthenticationFailure(buffer *bytes.Buffer) {}

func (a *AuthenticationFailure) DecodeAuthenticationFailure(byteArray *[]byte) {}
