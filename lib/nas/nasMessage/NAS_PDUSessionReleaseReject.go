//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionReleaseReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONRELEASEREJECTMessageIdentity
	nasType.Cause5GSM
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionReleaseReject(iei uint8) (pDUSessionReleaseReject *PDUSessionReleaseReject) {}

const (
	PDUSessionReleaseRejectExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionReleaseReject) EncodePDUSessionReleaseReject(buffer *bytes.Buffer) {}

func (a *PDUSessionReleaseReject) DecodePDUSessionReleaseReject(byteArray *[]byte) {}
