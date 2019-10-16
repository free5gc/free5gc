//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type ServiceAccept struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.ServiceAcceptMessageIdentity
	*nasType.PDUSessionStatus
	*nasType.PDUSessionReactivationResult
	*nasType.PDUSessionReactivationResultErrorCause
	*nasType.EAPMessage
}

func NewServiceAccept(iei uint8) (serviceAccept *ServiceAccept) {}

const (
	ServiceAcceptPDUSessionStatusType                       uint8 = 0x50
	ServiceAcceptPDUSessionReactivationResultType           uint8 = 0x26
	ServiceAcceptPDUSessionReactivationResultErrorCauseType uint8 = 0x72
	ServiceAcceptEAPMessageType                             uint8 = 0x78
)

func (a *ServiceAccept) EncodeServiceAccept(buffer *bytes.Buffer) {}

func (a *ServiceAccept) DecodeServiceAccept(byteArray *[]byte) {}
