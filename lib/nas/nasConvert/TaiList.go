//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"free5gc/lib/openapi/models"
	"reflect"
)

// TS 24.501 9.11.3.9
func TaiListToNas(taiList []models.Tai) (taiListNas []uint8) {}
