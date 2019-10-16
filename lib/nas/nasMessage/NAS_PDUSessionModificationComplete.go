//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionModificationComplete struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONMODIFICATIONCOMPLETEMessageIdentity
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionModificationComplete(iei uint8) (pDUSessionModificationComplete *PDUSessionModificationComplete) {}

const (
	PDUSessionModificationCompleteExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionModificationComplete) EncodePDUSessionModificationComplete(buffer *bytes.Buffer) {}

func (a *PDUSessionModificationComplete) DecodePDUSessionModificationComplete(byteArray *[]byte) {}
