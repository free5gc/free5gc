//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionAuthenticationCommand struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity
	nasType.EAPMessage
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionAuthenticationCommand(iei uint8) (pDUSessionAuthenticationCommand *PDUSessionAuthenticationCommand) {}

const (
	PDUSessionAuthenticationCommandExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionAuthenticationCommand) EncodePDUSessionAuthenticationCommand(buffer *bytes.Buffer) {}

func (a *PDUSessionAuthenticationCommand) DecodePDUSessionAuthenticationCommand(byteArray *[]byte) {}
