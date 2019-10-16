//go:binary-only-package

package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/logger"
	"free5gc/lib/ngap/ngapType"
	"free5gc/lib/openapi/models"
	"strings"
)

func TraceDataToModels(traceActivation ngapType.TraceActivation) (traceData models.TraceData) {}

func TraceDataToNgap(traceData models.TraceData, trsr string) (traceActivation ngapType.TraceActivation) {}
