//go:binary-only-package

package nasType

// ConfigurationUpdateIndication 9.11.3.18
// Iei Row, sBit, len = [0, 0], 8 , 4
// RED Row, sBit, len = [0, 0], 2 , 1
// ACK Row, sBit, len = [0, 0], 1 , 1
type ConfigurationUpdateIndication struct {
	Octet uint8
}

func NewConfigurationUpdateIndication(iei uint8) (configurationUpdateIndication *ConfigurationUpdateIndication) {}

// ConfigurationUpdateIndication 9.11.3.18
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *ConfigurationUpdateIndication) GetIei() (iei uint8) {}

// ConfigurationUpdateIndication 9.11.3.18
// Iei Row, sBit, len = [0, 0], 8 , 4
func (a *ConfigurationUpdateIndication) SetIei(iei uint8) {}

// ConfigurationUpdateIndication 9.11.3.18
// RED Row, sBit, len = [0, 0], 2 , 1
func (a *ConfigurationUpdateIndication) GetRED() (rED uint8) {}

// ConfigurationUpdateIndication 9.11.3.18
// RED Row, sBit, len = [0, 0], 2 , 1
func (a *ConfigurationUpdateIndication) SetRED(rED uint8) {}

// ConfigurationUpdateIndication 9.11.3.18
// ACK Row, sBit, len = [0, 0], 1 , 1
func (a *ConfigurationUpdateIndication) GetACK() (aCK uint8) {}

// ConfigurationUpdateIndication 9.11.3.18
// ACK Row, sBit, len = [0, 0], 1 , 1
func (a *ConfigurationUpdateIndication) SetACK(aCK uint8) {}
