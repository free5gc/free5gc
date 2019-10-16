//go:binary-only-package

package nasType

// MaximumNumberOfSupportedPacketFilters 9.11.4.9
// MaximumNumberOfSupportedPacketFilters Row, sBit, len = [0, 1], 8 , 10
type MaximumNumberOfSupportedPacketFilters struct {
	Iei   uint8
	Octet [2]uint8
}

func NewMaximumNumberOfSupportedPacketFilters(iei uint8) (maximumNumberOfSupportedPacketFilters *MaximumNumberOfSupportedPacketFilters) {}

// MaximumNumberOfSupportedPacketFilters 9.11.4.9
// Iei Row, sBit, len = [], 8, 8
func (a *MaximumNumberOfSupportedPacketFilters) GetIei() (iei uint8) {}

// MaximumNumberOfSupportedPacketFilters 9.11.4.9
// Iei Row, sBit, len = [], 8, 8
func (a *MaximumNumberOfSupportedPacketFilters) SetIei(iei uint8) {}

// MaximumNumberOfSupportedPacketFilters 9.11.4.9
// MaximumNumberOfSupportedPacketFilters Row, sBit, len = [0, 1], 8 , 10
func (a *MaximumNumberOfSupportedPacketFilters) GetMaximumNumberOfSupportedPacketFilters() (maximumNumberOfSupportedPacketFilters uint16) {}

// MaximumNumberOfSupportedPacketFilters 9.11.4.9
// MaximumNumberOfSupportedPacketFilters Row, sBit, len = [0, 1], 8 , 10
func (a *MaximumNumberOfSupportedPacketFilters) SetMaximumNumberOfSupportedPacketFilters(maximumNumberOfSupportedPacketFilters uint16) {}
