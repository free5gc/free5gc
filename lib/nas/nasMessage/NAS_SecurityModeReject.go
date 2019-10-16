//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type SecurityModeReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.SecurityModeRejectMessageIdentity
	nasType.Cause5GMM
}

func NewSecurityModeReject(iei uint8) (securityModeReject *SecurityModeReject) {}

func (a *SecurityModeReject) EncodeSecurityModeReject(buffer *bytes.Buffer) {}

func (a *SecurityModeReject) DecodeSecurityModeReject(byteArray *[]byte) {}
