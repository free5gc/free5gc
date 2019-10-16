//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionReleaseRequest struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONRELEASEREQUESTMessageIdentity
	*nasType.Cause5GSM
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionReleaseRequest(iei uint8) (pDUSessionReleaseRequest *PDUSessionReleaseRequest) {}

const (
	PDUSessionReleaseRequestCause5GSMType                            uint8 = 0x59
	PDUSessionReleaseRequestExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionReleaseRequest) EncodePDUSessionReleaseRequest(buffer *bytes.Buffer) {}

func (a *PDUSessionReleaseRequest) DecodePDUSessionReleaseRequest(byteArray *[]byte) {}
