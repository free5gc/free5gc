//go:binary-only-package

package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
)

func TaiToModels(tai ngapType.TAI) (modelsTai models.Tai) {}

func TaiToNgap(tai models.Tai) (ngapTai ngapType.TAI) {}
