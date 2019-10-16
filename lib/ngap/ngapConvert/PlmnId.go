//go:binary-only-package

package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"strings"
)

func PlmnIdToModels(ngapPlmnId ngapType.PLMNIdentity) (modelsPlmnid models.PlmnId) {}
func PlmnIdToNgap(modelsPlmnid models.PlmnId) (ngapPlmnId ngapType.PLMNIdentity) {}
