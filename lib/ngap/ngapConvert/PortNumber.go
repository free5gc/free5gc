//go:binary-only-package

package ngapConvert

import (
	"encoding/binary"
	"free5gc/lib/ngap/ngapType"
)

func PortNumberToInt(port ngapType.PortNumber) (portInt32 int32) {}

func PortNumberToNgap(portInt32 int32) (port ngapType.PortNumber) {}
