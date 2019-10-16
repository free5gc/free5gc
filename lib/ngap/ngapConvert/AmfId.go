//go:binary-only-package

package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/logger"
)

func AmfIdToNgap(amfId string) (regionId, setId, ptrId aper.BitString) {}

func AmfIdToModels(regionId, setId, ptrId aper.BitString) (amfId string) {}
