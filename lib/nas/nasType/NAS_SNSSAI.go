//go:binary-only-package

package nasType

// SNSSAI 9.11.2.8
// SST Row, sBit, len = [0, 0], 8 , 8
// SD Row, sBit, len = [1, 3], 8 , 24
// MappedHPLMNSST Row, sBit, len = [4, 4], 8 , 8
// MappedHPLMNSD Row, sBit, len = [5, 7], 8 , 24
type SNSSAI struct {
	Iei   uint8
	Len   uint8
	Octet [8]uint8
}

func NewSNSSAI(iei uint8) (sNSSAI *SNSSAI) {}

// SNSSAI 9.11.2.8
// Iei Row, sBit, len = [], 8, 8
func (a *SNSSAI) GetIei() (iei uint8) {}

// SNSSAI 9.11.2.8
// Iei Row, sBit, len = [], 8, 8
func (a *SNSSAI) SetIei(iei uint8) {}

// SNSSAI 9.11.2.8
// Len Row, sBit, len = [], 8, 8
func (a *SNSSAI) GetLen() (len uint8) {}

// SNSSAI 9.11.2.8
// Len Row, sBit, len = [], 8, 8
func (a *SNSSAI) SetLen(len uint8) {}

// SNSSAI 9.11.2.8
// SST Row, sBit, len = [0, 0], 8 , 8
func (a *SNSSAI) GetSST() (sST uint8) {}

// SNSSAI 9.11.2.8
// SST Row, sBit, len = [0, 0], 8 , 8
func (a *SNSSAI) SetSST(sST uint8) {}

// SNSSAI 9.11.2.8
// SD Row, sBit, len = [1, 3], 8 , 24
func (a *SNSSAI) GetSD() (sD [3]uint8) {}

// SNSSAI 9.11.2.8
// SD Row, sBit, len = [1, 3], 8 , 24
func (a *SNSSAI) SetSD(sD [3]uint8) {}

// SNSSAI 9.11.2.8
// MappedHPLMNSST Row, sBit, len = [4, 4], 8 , 8
func (a *SNSSAI) GetMappedHPLMNSST() (mappedHPLMNSST uint8) {}

// SNSSAI 9.11.2.8
// MappedHPLMNSST Row, sBit, len = [4, 4], 8 , 8
func (a *SNSSAI) SetMappedHPLMNSST(mappedHPLMNSST uint8) {}

// SNSSAI 9.11.2.8
// MappedHPLMNSD Row, sBit, len = [5, 7], 8 , 24
func (a *SNSSAI) GetMappedHPLMNSD() (mappedHPLMNSD [3]uint8) {}

// SNSSAI 9.11.2.8
// MappedHPLMNSD Row, sBit, len = [5, 7], 8 , 24
func (a *SNSSAI) SetMappedHPLMNSD(mappedHPLMNSD [3]uint8) {}
