//go:binary-only-package

package ngapConvert

import (
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
)

func RanIdToModels(ranNodeId ngapType.GlobalRANNodeID) (ranId models.GlobalRanNodeId) {}

func RanIDToNgap(modelsRanNodeId models.GlobalRanNodeId) (ngapRanNodeId ngapType.GlobalRANNodeID) {}
