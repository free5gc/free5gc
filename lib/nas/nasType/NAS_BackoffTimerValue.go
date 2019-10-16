//go:binary-only-package

package nasType

// BackoffTimerValue 9.11.2.5
// UnitTimerValue Row, sBit, len = [0, 0], 8 , 3
// TimerValue Row, sBit, len = [0, 0], 5 , 5
type BackoffTimerValue struct {
	Iei   uint8
	Len   uint8
	Octet uint8
}

func NewBackoffTimerValue(iei uint8) (backoffTimerValue *BackoffTimerValue) {}

// BackoffTimerValue 9.11.2.5
// Iei Row, sBit, len = [], 8, 8
func (a *BackoffTimerValue) GetIei() (iei uint8) {}

// BackoffTimerValue 9.11.2.5
// Iei Row, sBit, len = [], 8, 8
func (a *BackoffTimerValue) SetIei(iei uint8) {}

// BackoffTimerValue 9.11.2.5
// Len Row, sBit, len = [], 8, 8
func (a *BackoffTimerValue) GetLen() (len uint8) {}

// BackoffTimerValue 9.11.2.5
// Len Row, sBit, len = [], 8, 8
func (a *BackoffTimerValue) SetLen(len uint8) {}

// BackoffTimerValue 9.11.2.5
// UnitTimerValue Row, sBit, len = [0, 0], 8 , 3
func (a *BackoffTimerValue) GetUnitTimerValue() (unitTimerValue uint8) {}

// BackoffTimerValue 9.11.2.5
// UnitTimerValue Row, sBit, len = [0, 0], 8 , 3
func (a *BackoffTimerValue) SetUnitTimerValue(unitTimerValue uint8) {}

// BackoffTimerValue 9.11.2.5
// TimerValue Row, sBit, len = [0, 0], 5 , 5
func (a *BackoffTimerValue) GetTimerValue() (timerValue uint8) {}

// BackoffTimerValue 9.11.2.5
// TimerValue Row, sBit, len = [0, 0], 5 , 5
func (a *BackoffTimerValue) SetTimerValue(timerValue uint8) {}
