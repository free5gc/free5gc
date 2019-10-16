//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type IdentityRequest struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.IdentityRequestMessageIdentity
	nasType.SpareHalfOctetAndIdentityType
}

func NewIdentityRequest(iei uint8) (identityRequest *IdentityRequest) {}

func (a *IdentityRequest) EncodeIdentityRequest(buffer *bytes.Buffer) {}

func (a *IdentityRequest) DecodeIdentityRequest(byteArray *[]byte) {}
