//go:binary-only-package

package ngapConvert

import (
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"strconv"
)

func UEAmbrToModels(ngapAmbr ngapType.UEAggregateMaximumBitRate) (modelAmbr models.Ambr) {}

func UEAmbrToInt64(modelAmbr models.Ambr) (ueAmbrUL, ueAmbrDL int64) {}
