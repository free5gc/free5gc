//go:binary-only-package

package nasConvert

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/openapi/models"
)

func PDUSessionTypeToModels(nasPduSessType uint8) (pduSessType models.PduSessionType) {}

func ModelsToPDUSessionType(pduSessType models.PduSessionType) (nasPduSessType uint8) {}
