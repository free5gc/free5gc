//go:binary-only-package

package nasConvert

import (
	"free5gc/lib/nas/nasMessage"
)

// TS 24.008 10.5.7.4a
func GPRSTimer3ToNas(timerValue int) (timerValueNas uint8) {}
