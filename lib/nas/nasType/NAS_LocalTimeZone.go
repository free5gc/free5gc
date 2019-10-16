//go:binary-only-package

package nasType

// LocalTimeZone 9.11.3.52
// TimeZone Row, sBit, len = [0, 0], 8 , 8
type LocalTimeZone struct {
	Iei   uint8
	Octet uint8
}

func NewLocalTimeZone(iei uint8) (localTimeZone *LocalTimeZone) {}

// LocalTimeZone 9.11.3.52
// Iei Row, sBit, len = [], 8, 8
func (a *LocalTimeZone) GetIei() (iei uint8) {}

// LocalTimeZone 9.11.3.52
// Iei Row, sBit, len = [], 8, 8
func (a *LocalTimeZone) SetIei(iei uint8) {}

// LocalTimeZone 9.11.3.52
// TimeZone Row, sBit, len = [0, 0], 8 , 8
func (a *LocalTimeZone) GetTimeZone() (timeZone uint8) {}

// LocalTimeZone 9.11.3.52
// TimeZone Row, sBit, len = [0, 0], 8 , 8
func (a *LocalTimeZone) SetTimeZone(timeZone uint8) {}
