//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type PDUSessionEstablishmentAccept struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.PDUSessionID
	nasType.PTI
	nasType.PDUSESSIONESTABLISHMENTACCEPTMessageIdentity
	nasType.SelectedSSCModeAndSelectedPDUSessionType
	nasType.AuthorizedQosRules
	nasType.SessionAMBR
	*nasType.Cause5GSM
	*nasType.PDUAddress
	*nasType.RQTimerValue
	*nasType.SNSSAI
	*nasType.AlwaysonPDUSessionIndication
	*nasType.MappedEPSBearerContexts
	*nasType.EAPMessage
	*nasType.AuthorizedQosFlowDescriptions
	*nasType.ExtendedProtocolConfigurationOptions
	*nasType.DNN
}

func NewPDUSessionEstablishmentAccept(iei uint8) (pDUSessionEstablishmentAccept *PDUSessionEstablishmentAccept) {}

const (
	PDUSessionEstablishmentAcceptCause5GSMType                            uint8 = 0x59
	PDUSessionEstablishmentAcceptPDUAddressType                           uint8 = 0x29
	PDUSessionEstablishmentAcceptRQTimerValueType                         uint8 = 0x56
	PDUSessionEstablishmentAcceptSNSSAIType                               uint8 = 0x22
	PDUSessionEstablishmentAcceptAlwaysonPDUSessionIndicationType         uint8 = 0x08
	PDUSessionEstablishmentAcceptMappedEPSBearerContextsType              uint8 = 0x75
	PDUSessionEstablishmentAcceptEAPMessageType                           uint8 = 0x78
	PDUSessionEstablishmentAcceptAuthorizedQosFlowDescriptionsType        uint8 = 0x79
	PDUSessionEstablishmentAcceptExtendedProtocolConfigurationOptionsType uint8 = 0x7B
	PDUSessionEstablishmentAcceptDNNType                                  uint8 = 0x25
)

func (a *PDUSessionEstablishmentAccept) EncodePDUSessionEstablishmentAccept(buffer *bytes.Buffer) {}

func (a *PDUSessionEstablishmentAccept) DecodePDUSessionEstablishmentAccept(byteArray *[]byte) {}
