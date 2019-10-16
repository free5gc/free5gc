//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type IdentityResponse struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.IdentityResponseMessageIdentity
	nasType.MobileIdentity
}

func NewIdentityResponse(iei uint8) (identityResponse *IdentityResponse) {}

func (a *IdentityResponse) EncodeIdentityResponse(buffer *bytes.Buffer) {}

func (a *IdentityResponse) DecodeIdentityResponse(byteArray *[]byte) {}
