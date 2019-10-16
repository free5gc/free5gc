//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
)

func AmfIdToNas(amfId string) (amfRegionId uint8, amfSetId uint16, amfPointer uint8) {}

func AmfIdToModels(amfRegionId uint8, amfSetId uint16, amfPointer uint8) (amfId string) {}
