//go:binary-only-package

package nasConvert

import (
	"free5gc/lib/nas/nasType"
)

// TS 24.501 9.11.3.35, TS 24.008 10.5.3.5a
func FullNetworkNameToNas(name string) (fullNetworkName nasType.FullNameForNetwork) {}

func ShortNetworkNameToNas(name string) (shortNetworkName nasType.ShortNameForNetwork) {}
