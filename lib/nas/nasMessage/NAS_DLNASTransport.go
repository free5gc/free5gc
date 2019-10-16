//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type DLNASTransport struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.DLNASTRANSPORTMessageIdentity
	nasType.SpareHalfOctetAndPayloadContainerType
	nasType.PayloadContainer
	*nasType.PduSessionID2Value
	*nasType.AdditionalInformation
	*nasType.Cause5GMM
	*nasType.BackoffTimerValue
}

func NewDLNASTransport(iei uint8) (dLNASTransport *DLNASTransport) {}

const (
	DLNASTransportPduSessionID2ValueType    uint8 = 0x12
	DLNASTransportAdditionalInformationType uint8 = 0x24
	DLNASTransportCause5GMMType             uint8 = 0x58
	DLNASTransportBackoffTimerValueType     uint8 = 0x37
)

func (a *DLNASTransport) EncodeDLNASTransport(buffer *bytes.Buffer) {}

func (a *DLNASTransport) DecodeDLNASTransport(byteArray *[]byte) {}
