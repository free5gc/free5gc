//go:binary-only-package

package nasType

// ConfigurationUpdateCompleteMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
type ConfigurationUpdateCompleteMessageIdentity struct {
	Octet uint8
}

func NewConfigurationUpdateCompleteMessageIdentity() (configurationUpdateCompleteMessageIdentity *ConfigurationUpdateCompleteMessageIdentity) {}

// ConfigurationUpdateCompleteMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ConfigurationUpdateCompleteMessageIdentity) GetMessageType() (messageType uint8) {}

// ConfigurationUpdateCompleteMessageIdentity 9.7
// MessageType Row, sBit, len = [0, 0], 8 , 8
func (a *ConfigurationUpdateCompleteMessageIdentity) SetMessageType(messageType uint8) {}
