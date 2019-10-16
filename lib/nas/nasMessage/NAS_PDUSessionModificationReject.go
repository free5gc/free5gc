//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionModificationReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONMODIFICATIONREJECTMessageIdentity
	nasType.Cause5GSM
	*nasType.BackoffTimerValue
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionModificationReject(iei uint8) (pDUSessionModificationReject *PDUSessionModificationReject) {}

const (
	PDUSessionModificationRejectBackoffTimerValueType                    uint8 = 0x37
	PDUSessionModificationRejectExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionModificationReject) EncodePDUSessionModificationReject(buffer *bytes.Buffer) {}

func (a *PDUSessionModificationReject) DecodePDUSessionModificationReject(byteArray *[]byte) {}
