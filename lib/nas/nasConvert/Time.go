//go:binary-only-package

package nasConvert

import (
	"free5gc/lib/nas/nasType"
	"strings"
)

func LocalTimeZoneToNas(timezone string) (nasTimezone nasType.LocalTimeZone) {}

func DaylightSavingTimeToNas(timezone string) (nasDaylightSavingTimeToNas nasType.NetworkDaylightSavingTime) {}
