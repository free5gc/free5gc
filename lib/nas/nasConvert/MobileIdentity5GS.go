//go:binary-only-package

package nasConvert

import (
	"encoding/hex"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"free5gc/lib/openapi/models"
	"strconv"
)

func GetTypeOfIdentity(buf byte) uint8 {}

// TS 24.501 9.11 .3.4
func GetPlmnIdAndSupiFromSuci(buf []byte) (plmnId, supi string) {}

// nasType: TS 24.501 9.11.3.4
func GutiToString(buf []byte) (guami models.Guami, guti string) {}

func GutiToNas(guti string) (gutiNas nasType.GUTI5G) {}

// PEI: ^(imei-[0-9]{15}|imeisv-[0-9]{16}|.+)$
func PeiToString(buf []byte) (pei string) {}
