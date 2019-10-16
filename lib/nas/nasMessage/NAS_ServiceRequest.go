//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type ServiceRequest struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.ServiceRequestMessageIdentity
	nasType.ServiceTypeAndNgksi
	nasType.TMSI5GS
	*nasType.UplinkDataStatus
	*nasType.PDUSessionStatus
	*nasType.AllowedPDUSessionStatus
	*nasType.NASMessageContainer
}

func NewServiceRequest(iei uint8) (serviceRequest *ServiceRequest) {}

const (
	ServiceRequestUplinkDataStatusType        uint8 = 0x40
	ServiceRequestPDUSessionStatusType        uint8 = 0x50
	ServiceRequestAllowedPDUSessionStatusType uint8 = 0x25
	ServiceRequestNASMessageContainerType     uint8 = 0x71
)

func (a *ServiceRequest) EncodeServiceRequest(buffer *bytes.Buffer) {}

func (a *ServiceRequest) DecodeServiceRequest(byteArray *[]byte) {}
