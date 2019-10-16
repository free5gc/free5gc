//go:binary-only-package

package ngapConvert

import (
	"encoding/hex"
	"free5gc/lib/aper"
	"free5gc/lib/ngap/logger"
)

func BitStringToHex(bitString *aper.BitString) (hexString string) {}

func HexToBitString(hexString string, bitLength int) (bitString aper.BitString) {}

func ByteToBitString(byteArray []byte, bitLength int) (bitString aper.BitString) {}
