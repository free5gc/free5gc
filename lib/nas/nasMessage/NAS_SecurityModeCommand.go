//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type SecurityModeCommand struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.SecurityModeCommandMessageIdentity
	nasType.SelectedNASSecurityAlgorithms
	nasType.SpareHalfOctetAndNgksi
	nasType.ReplayedUESecurityCapabilities
	*nasType.IMEISVRequest
	*nasType.SelectedEPSNASSecurityAlgorithms
	*nasType.Additional5GSecurityInformation
	*nasType.EAPMessage
	*nasType.ABBA
	*nasType.ReplayedS1UESecurityCapabilities
}

func NewSecurityModeCommand(iei uint8) (securityModeCommand *SecurityModeCommand) {}

const (
	SecurityModeCommandIMEISVRequestType                    uint8 = 0x0E
	SecurityModeCommandSelectedEPSNASSecurityAlgorithmsType uint8 = 0x57
	SecurityModeCommandAdditional5GSecurityInformationType  uint8 = 0x36
	SecurityModeCommandEAPMessageType                       uint8 = 0x78
	SecurityModeCommandABBAType                             uint8 = 0x38
	SecurityModeCommandReplayedS1UESecurityCapabilitiesType uint8 = 0x19
)

func (a *SecurityModeCommand) EncodeSecurityModeCommand(buffer *bytes.Buffer) {}

func (a *SecurityModeCommand) DecodeSecurityModeCommand(byteArray *[]byte) {}
