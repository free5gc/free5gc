//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionModificationCommand struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONMODIFICATIONCOMMANDMessageIdentity
	*nasType.Cause5GSM
	*nasType.SessionAMBR
	*nasType.RQTimerValue
	*nasType.AlwaysonPDUSessionIndication
	*nasType.AuthorizedQosRules
	*nasType.MappedEPSBearerContexts
	*nasType.AuthorizedQosFlowDescriptions
	*nasType.ExtendedProtocolConfigurationOptions
}

func NewPDUSessionModificationCommand(iei uint8) (pDUSessionModificationCommand *PDUSessionModificationCommand) {}

const (
	PDUSessionModificationCommandCause5GSMType                            uint8 = 0x59
	PDUSessionModificationCommandSessionAMBRType                          uint8 = 0x2A
	PDUSessionModificationCommandRQTimerValueType                         uint8 = 0x56
	PDUSessionModificationCommandAlwaysonPDUSessionIndicationType         uint8 = 0x08
	PDUSessionModificationCommandAuthorizedQosRulesType                   uint8 = 0x7A
	PDUSessionModificationCommandMappedEPSBearerContextsType              uint8 = 0x7F
	PDUSessionModificationCommandAuthorizedQosFlowDescriptionsType        uint8 = 0x79
	PDUSessionModificationCommandExtendedProtocolConfigurationOptionsType uint8 = 0x7B
)

func (a *PDUSessionModificationCommand) EncodePDUSessionModificationCommand(buffer *bytes.Buffer) {}

func (a *PDUSessionModificationCommand) DecodePDUSessionModificationCommand(byteArray *[]byte) {}
