//go:binary-only-package

package nasType

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Year Row, sBit, len = [0, 0], 8 , 8
// Month Row, sBit, len = [1, 1], 8 , 8
// Day Row, sBit, len = [2, 2], 8 , 8
// Hour Row, sBit, len = [3, 3], 8 , 8
// Minute Row, sBit, len = [4, 4], 8 , 8
// Second Row, sBit, len = [5, 5], 8 , 8
// TimeZone Row, sBit, len = [6, 6], 8 , 8
type UniversalTimeAndLocalTimeZone struct {
	Iei   uint8
	Octet [7]uint8
}

func NewUniversalTimeAndLocalTimeZone(iei uint8) (universalTimeAndLocalTimeZone *UniversalTimeAndLocalTimeZone) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Iei Row, sBit, len = [], 8, 8
func (a *UniversalTimeAndLocalTimeZone) GetIei() (iei uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Iei Row, sBit, len = [], 8, 8
func (a *UniversalTimeAndLocalTimeZone) SetIei(iei uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Year Row, sBit, len = [0, 0], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetYear() (year uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Year Row, sBit, len = [0, 0], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetYear(year uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Month Row, sBit, len = [1, 1], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetMonth() (month uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Month Row, sBit, len = [1, 1], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetMonth(month uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Day Row, sBit, len = [2, 2], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetDay() (day uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Day Row, sBit, len = [2, 2], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetDay(day uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Hour Row, sBit, len = [3, 3], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetHour() (hour uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Hour Row, sBit, len = [3, 3], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetHour(hour uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Minute Row, sBit, len = [4, 4], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetMinute() (minute uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Minute Row, sBit, len = [4, 4], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetMinute(minute uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Second Row, sBit, len = [5, 5], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetSecond() (second uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// Second Row, sBit, len = [5, 5], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetSecond(second uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// TimeZone Row, sBit, len = [6, 6], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) GetTimeZone() (timeZone uint8) {}

// UniversalTimeAndLocalTimeZone 9.11.3.53
// TimeZone Row, sBit, len = [6, 6], 8 , 8
func (a *UniversalTimeAndLocalTimeZone) SetTimeZone(timeZone uint8) {}
