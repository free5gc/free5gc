//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/models"
)

// TS 24.501 9.11.3.49
func PartialServiceAreaListToNas(plmnID models.PlmnId, serviceAreaRestriction models.ServiceAreaRestriction) (partialServiceAreaList []byte) {}
