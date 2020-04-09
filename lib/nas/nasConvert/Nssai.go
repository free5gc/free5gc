//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"fmt"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
)

func RequestedNssaiToModels(nasNssai *nasType.RequestedNSSAI) (nssai []models.Snssai) {}

func requestedSnssaiToModels(buf []byte) (snssai models.Snssai, length int) {}

func RejectedNssaiToNas(rejectedNssaiInPlmn []models.Snssai, rejectedNssaiInTa []models.Snssai) (rejectedNssaiNas nasType.RejectedNSSAI) {}
