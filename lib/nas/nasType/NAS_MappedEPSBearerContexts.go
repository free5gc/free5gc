//go:binary-only-package

package nasType

// MappedEPSBearerContexts 9.11.4.8
// MappedEPSBearerContext Row, sBit, len = [0, 0], 8 , INF
type MappedEPSBearerContexts struct {
	Iei    uint8
	Len    uint16
	Buffer []uint8
}

func NewMappedEPSBearerContexts(iei uint8) (mappedEPSBearerContexts *MappedEPSBearerContexts) {}

// MappedEPSBearerContexts 9.11.4.8
// Iei Row, sBit, len = [], 8, 8
func (a *MappedEPSBearerContexts) GetIei() (iei uint8) {}

// MappedEPSBearerContexts 9.11.4.8
// Iei Row, sBit, len = [], 8, 8
func (a *MappedEPSBearerContexts) SetIei(iei uint8) {}

// MappedEPSBearerContexts 9.11.4.8
// Len Row, sBit, len = [], 8, 16
func (a *MappedEPSBearerContexts) GetLen() (len uint16) {}

// MappedEPSBearerContexts 9.11.4.8
// Len Row, sBit, len = [], 8, 16
func (a *MappedEPSBearerContexts) SetLen(len uint16) {}

// MappedEPSBearerContexts 9.11.4.8
// MappedEPSBearerContext Row, sBit, len = [0, 0], 8 , INF
func (a *MappedEPSBearerContexts) GetMappedEPSBearerContext() (mappedEPSBearerContext []uint8) {}

// MappedEPSBearerContexts 9.11.4.8
// MappedEPSBearerContext Row, sBit, len = [0, 0], 8 , INF
func (a *MappedEPSBearerContexts) SetMappedEPSBearerContext(mappedEPSBearerContext []uint8) {}
