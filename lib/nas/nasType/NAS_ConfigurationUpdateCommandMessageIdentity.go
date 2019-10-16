//go:binary-only-package

package nasType

// ConfigurationUpdateCommandMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type ConfigurationUpdateCommandMessageIdentity struct {
	Octet uint8
}

func NewConfigurationUpdateCommandMessageIdentity() (configurationUpdateCommandMessageIdentity *ConfigurationUpdateCommandMessageIdentity) {}

// ConfigurationUpdateCommandMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ConfigurationUpdateCommandMessageIdentity) GetMessageType() (messageType uint8) {}

// ConfigurationUpdateCommandMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ConfigurationUpdateCommandMessageIdentity) SetMessageType(messageType uint8) {}
