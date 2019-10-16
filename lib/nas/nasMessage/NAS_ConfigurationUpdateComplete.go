//go:binary-only-package

package nasMessage

import (
	"bytes"
	"encoding/binary"
	"free5gc/lib/nas/nasType"
)

type ConfigurationUpdateComplete struct {
	nasType.ExtendedProtocolDiscriminator
	nasType.SpareHalfOctetAndSecurityHeaderType
	nasType.ConfigurationUpdateCompleteMessageIdentity
}

func NewConfigurationUpdateComplete(iei uint8) (configurationUpdateComplete *ConfigurationUpdateComplete) {}

func (a *ConfigurationUpdateComplete) EncodeConfigurationUpdateComplete(buffer *bytes.Buffer) {}

func (a *ConfigurationUpdateComplete) DecodeConfigurationUpdateComplete(byteArray *[]byte) {}
