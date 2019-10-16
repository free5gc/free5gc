//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionEstablishmentReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONESTABLISHMENTREJECTMessageIdentity
	nasType.Cause5GSM
	*nasType.BackoffTimerValue
	*nasType.AllowedSSCMode
	*nasType.EAPMessage
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionEstablishmentReject(iei uint8) (pDUSessionEstablishmentReject *PDUSessionEstablishmentReject) {}

const (
	PDUSessionEstablishmentRejectBackoffTimerValueType                    uint8 = 0x37
	PDUSessionEstablishmentRejectAllowedSSCModeType                       uint8 = 0x0F
	PDUSessionEstablishmentRejectEAPMessageType                           uint8 = 0x78
	PDUSessionEstablishmentRejectExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionEstablishmentReject) EncodePDUSessionEstablishmentReject(buffer *bytes.Buffer) {}

func (a *PDUSessionEstablishmentReject) DecodePDUSessionEstablishmentReject(byteArray *[]byte) {}
