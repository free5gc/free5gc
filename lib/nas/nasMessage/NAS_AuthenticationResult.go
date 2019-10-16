//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type AuthenticationResult struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.AuthenticationResultMessageIdentity
	nasType.SpareHalfOctetAndNgksi
	nasType.EAPMessage
	*nasType.ABBA
}

func NewAuthenticationResult(iei uint8) (authenticationResult *AuthenticationResult) {}

const (
	AuthenticationResultABBAType uint8 = 0x38
)

func (a *AuthenticationResult) EncodeAuthenticationResult(buffer *bytes.Buffer) {}

func (a *AuthenticationResult) DecodeAuthenticationResult(byteArray *[]byte) {}
