//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionAuthenticationResult struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONAUTHENTICATIONRESULTMessageIdentity
	*nasType.EAPMessage
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionAuthenticationResult(iei uint8) (pDUSessionAuthenticationResult *PDUSessionAuthenticationResult) {}

const (
	PDUSessionAuthenticationResultEAPMessageType                           uint8 = 0x78
	PDUSessionAuthenticationResultExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionAuthenticationResult) EncodePDUSessionAuthenticationResult(buffer *bytes.Buffer) {}

func (a *PDUSessionAuthenticationResult) DecodePDUSessionAuthenticationResult(byteArray *[]byte) {}
