//go:binary-only-package

package ngapConvert

import (
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
)

// TS 38.413 9.3.1.85
func RATRestrictionInformationToNgap(ratType models.RatType) (ratResInfo ngapType.RATRestrictionInformation) {}
