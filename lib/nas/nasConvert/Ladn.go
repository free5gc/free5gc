//go:binary-only-package

package nasConvert

import (
	"free5gc/src/amf/amf_context"
)

func LadnToModels(buf []uint8) (dnnValues []string) {}

func LadnToNas(ladn amf_context.LADN) (ladnNas []uint8) {}
