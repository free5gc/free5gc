//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionReleaseComplete struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONRELEASECOMPLETEMessageIdentity
	*nasType.Cause5GSM
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionReleaseComplete(iei uint8) (pDUSessionReleaseComplete *PDUSessionReleaseComplete) {}

const (
	PDUSessionReleaseCompleteCause5GSMType                            uint8 = 0x59
	PDUSessionReleaseCompleteExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionReleaseComplete) EncodePDUSessionReleaseComplete(buffer *bytes.Buffer) {}

func (a *PDUSessionReleaseComplete) DecodePDUSessionReleaseComplete(byteArray *[]byte) {}
