//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
)

func SnssaiToModels(nasSnssai *nasType.SNSSAI) (snssai models.Snssai) {}

func SnssaiToNas(snssai models.Snssai) (buf []uint8) {}

func RejectedSnssaiToNas(snssai models.Snssai, rejectCause uint8) (rejectedSnssai []uint8) {}
