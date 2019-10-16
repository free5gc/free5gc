//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"fmt"
	"free5gc/lib/openapi/models"
)

//  subclause 9.11.3.53A in 3GPP TS 24.501
func UpuInfoToNas(upuInfo models.UpuInfo) (buf []uint8) {}

func upuInfoGetHeader(reg bool, ack bool) (buf uint8) {}

func UpuAckToModels(buf []uint8) (string, error) {}
