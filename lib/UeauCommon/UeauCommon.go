//go:binary-only-package

package UeauCommon

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

const (
	FC_FOR_CK_PRIME_IK_PRIME_DERIVATION  = "20"
	FC_FOR_KSEAF_DERIVATION              = "6C"
	FC_FOR_RES_STAR_XRES_STAR_DERIVATION = "6B"
	FC_FOR_KAUSF_DERIVATION              = "6A"
	FC_FOR_KAMF_DERIVATION               = "6D"
	FC_FOR_KGNB_KN3IWF_DERIVATION        = "6E"
	FC_FOR_NH_DERIVATION                 = "6F"
	FC_FOR_ALGORITHM_KEY_DERIVATION      = "69"
)

func KDFLen(input []byte) []byte {}

// This function implements the KDF defined in TS.33220 cluase B.2.0.
//
// For P0-Pn, the ones that will be used directly as a string (e.g. "WLAN") should be type-casted by []byte(),
// and the ones originially in hex (e.g. "bb52e91c747a") should be converted by using hex.DecodeString().
//
// For L0-Ln, use KDFLen() function to calculate them (e.g. KDFLen(P0)).
func GetKDFValue(key []byte, FC string, param ...[]byte) []byte {}
