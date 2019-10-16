//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type Status5GMM struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.STATUSMessageIdentity5GMM
	nasType.Cause5GMM
}

func NewStatus5GMM(iei uint8) (status5GMM *Status5GMM) {}

func (a *Status5GMM) EncodeStatus5GMM(buffer *bytes.Buffer) {}

func (a *Status5GMM) DecodeStatus5GMM(byteArray *[]byte) {}
