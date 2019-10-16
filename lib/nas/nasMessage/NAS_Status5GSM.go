//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type Status5GSM struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.STATUSMessageIdentity5GSM
	nasType.Cause5GSM
}

func NewStatus5GSM(iei uint8) (status5GSM *Status5GSM) {}

func (a *Status5GSM) EncodeStatus5GSM(buffer *bytes.Buffer) {}

func (a *Status5GSM) DecodeStatus5GSM(byteArray *[]byte) {}
