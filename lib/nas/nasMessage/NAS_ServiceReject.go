//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type ServiceReject struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.ServiceRejectMessageIdentity
	nasType.Cause5GMM
	*nasType.PDUSessionStatus
	*nasType.T3346Value
	*nasType.EAPMessage
}

func NewServiceReject(iei uint8) (serviceReject *ServiceReject) {}

const (
	ServiceRejectPDUSessionStatusType uint8 = 0x50
	ServiceRejectT3346ValueType       uint8 = 0x5F
	ServiceRejectEAPMessageType       uint8 = 0x78
)

func (a *ServiceReject) EncodeServiceReject(buffer *bytes.Buffer) {}

func (a *ServiceReject) DecodeServiceReject(byteArray *[]byte) {}
