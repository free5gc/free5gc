//go:binary-only-package

package nasType

// RQTimerValue 9.11.2.3
// Unit Row, sBit, len = [0, 0], 8 , 3
// TimerValue Row, sBit, len = [0, 0], 5 , 5
type RQTimerValue struct {
	Iei   uint8
	Octet uint8
}

func NewRQTimerValue(iei uint8) (rQTimerValue *RQTimerValue) {}

// RQTimerValue 9.11.2.3
// Iei Row, sBit, len = [], 8, 8
func (a *RQTimerValue) GetIei() (iei uint8) {}

// RQTimerValue 9.11.2.3
// Iei Row, sBit, len = [], 8, 8
func (a *RQTimerValue) SetIei(iei uint8) {}

// RQTimerValue 9.11.2.3
// Unit Row, sBit, len = [0, 0], 8 , 3
func (a *RQTimerValue) GetUnit() (unit uint8) {}

// RQTimerValue 9.11.2.3
// Unit Row, sBit, len = [0, 0], 8 , 3
func (a *RQTimerValue) SetUnit(unit uint8) {}

// RQTimerValue 9.11.2.3
// TimerValue Row, sBit, len = [0, 0], 5 , 5
func (a *RQTimerValue) GetTimerValue() (timerValue uint8) {}

// RQTimerValue 9.11.2.3
// TimerValue Row, sBit, len = [0, 0], 5 , 5
func (a *RQTimerValue) SetTimerValue(timerValue uint8) {}
