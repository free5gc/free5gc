//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type ULNASTransport struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.ULNASTRANSPORTMessageIdentity
	nasType.SpareHalfOctetAndPayloadContainerType
	nasType.PayloadContainer
	*nasType.PduSessionID2Value
	*nasType.OldPDUSessionID
	*nasType.RequestType
	*nasType.SNSSAI
	*nasType.DNN
	*nasType.AdditionalInformation
}

func NewULNASTransport(iei uint8) (uLNASTransport *ULNASTransport) {}

const (
	ULNASTransportPduSessionID2ValueType    uint8 = 0x12
	ULNASTransportOldPDUSessionIDType       uint8 = 0x59
	ULNASTransportRequestTypeType           uint8 = 0x08
	ULNASTransportSNSSAIType                uint8 = 0x22
	ULNASTransportDNNType                   uint8 = 0x25
	ULNASTransportAdditionalInformationType uint8 = 0x24
)

func (a *ULNASTransport) EncodeULNASTransport(buffer *bytes.Buffer) {}

func (a *ULNASTransport) DecodeULNASTransport(byteArray *[]byte) {}
