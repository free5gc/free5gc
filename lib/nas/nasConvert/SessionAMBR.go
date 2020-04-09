//go:binary-only-package

package nasConvert

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
)

func ModelsToSessionAMBR(ambr *models.Ambr) (sessAmbr nasType.SessionAMBR) {}

func strToAMBRUnit(unit string) uint8 {}
