//go:binary-only-package

package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
)

func SNssaiToModels(ngapSnssai ngapType.SNSSAI) (modelsSnssai models.Snssai) {}

func SNssaiToNgap(modelsSnssai models.Snssai) (ngapSnssai ngapType.SNSSAI) {}
