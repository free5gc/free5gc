//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionModificationCommandReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity
	nasType.Cause5GSM
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionModificationCommandReject(iei uint8) (pDUSessionModificationCommandReject *PDUSessionModificationCommandReject) {}

const (
	PDUSessionModificationCommandRejectExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionModificationCommandReject) EncodePDUSessionModificationCommandReject(buffer *bytes.Buffer) {}

func (a *PDUSessionModificationCommandReject) DecodePDUSessionModificationCommandReject(byteArray *[]byte) {}
