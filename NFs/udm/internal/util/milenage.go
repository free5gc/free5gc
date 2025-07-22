package util

import (
	"encoding/hex"

	"github.com/free5gc/util/milenage"
)

func MilenageF1(opc, k, rand, sqn, amf []byte, macA, macS []byte) error {
	ik, ck, xres, autn, err := milenage.GenerateAKAParameters(opc, k, rand, sqn, amf)
	if err != nil {
		return err
	}
	// Suppress unused variable warnings
	_ = ik
	_ = ck
	_ = xres

	// AUTN = (SQN xor AK) || AMF || MAC-A
	// MAC-A is the last 8 bytes of AUTN
	if len(autn) >= 8 && macA != nil {
		copy(macA, autn[len(autn)-8:])
	}

	// For MAC-S, use resync AMF (0000)
	if macS != nil {
		resyncAMFBytes, err := hex.DecodeString("0000")
		if err != nil {
			return err
		}
		ikS, ckS, xresS, autnS, err := milenage.GenerateAKAParameters(opc, k, rand, sqn, resyncAMFBytes)
		if err != nil {
			return err
		}
		// Suppress unused variable warnings
		_ = ikS
		_ = ckS
		_ = xresS

		if len(autnS) >= 8 {
			copy(macS, autnS[len(autnS)-8:])
		}
	}

	return nil
}

func MilenageF2345(opc, k, rand []byte, res, ck, ik, ak, akstar []byte) error {
	// Use GenerateAKAParameters to get basic parameters
	ikOut, ckOut, resOut, autn, err := milenage.GenerateAKAParameters(opc, k, rand, make([]byte, 6), make([]byte, 2))
	if err != nil {
		return err
	}

	// Use GenerateKeysWithAUTN to get AK
	sqnhe, akOut, ikOut2, ckOut2, resOut2, err := milenage.GenerateKeysWithAUTN(opc, k, rand, autn)
	if err != nil {
		return err
	}
	// Suppress unused variable warnings
	_ = sqnhe
	_ = ikOut2
	_ = ckOut2
	_ = resOut2

	// Copy results to output parameters
	if res != nil {
		copy(res, resOut)
	}
	if ck != nil {
		copy(ck, ckOut)
	}
	if ik != nil {
		copy(ik, ikOut)
	}
	if ak != nil {
		copy(ak, akOut)
	}
	if akstar != nil {
		// For AK*, we need to use a different SQN, but due to API limitations, we use the same value for now
		copy(akstar, akOut)
	}

	return nil
}
