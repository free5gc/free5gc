//go:binary-only-package

package ngap

import (
	"fmt"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/ngapType"
	"reflect"
	"strings"
)

// Decoder is to decode raw data to NGAP pdu pointer with PER Aligned
func Decoder(b []byte) (pdu *ngapType.NGAPPDU, err error) {}

// Encoder is to NGAP pdu to raw data with PER Aligned
func Encoder(pdu ngapType.NGAPPDU) ([]byte, error) {}

func PrintResult(v reflect.Value, layer int) (s string) {}
