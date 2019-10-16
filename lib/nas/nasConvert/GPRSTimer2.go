//go:binary-only-package

package nasConvert

import (
	"free5gc/lib/nas/logger"
)

// TS 24.008 10.5.7.4, TS 24.501 9.11.2.4
// the unit of timerValue is second
func GPRSTimer2ToNas(timerValue int) (timerValueNas uint8) {}
