//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionAuthenticationComplete struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONAUTHENTICATIONCOMPLETEMessageIdentity
	nasType.EAPMessage
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionAuthenticationComplete(iei uint8) (pDUSessionAuthenticationComplete *PDUSessionAuthenticationComplete) {}

const (
	PDUSessionAuthenticationCompleteExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionAuthenticationComplete) EncodePDUSessionAuthenticationComplete(buffer *bytes.Buffer) {}

func (a *PDUSessionAuthenticationComplete) DecodePDUSessionAuthenticationComplete(byteArray *[]byte) {}
