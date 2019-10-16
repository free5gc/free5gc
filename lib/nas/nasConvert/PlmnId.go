//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"free5gc/lib/openapi/models"
	"strconv"
)

func PlmnIDToNas(plmnID models.PlmnId) (plmnNas []uint8) {}

func PlmnIDToString(nasBuf []byte) (plmnID string) {}
