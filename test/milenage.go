package test

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
	// Directly call F2345 to get all keys including AK and AK*
	resOut, ckOut, ikOut, akOut, akstarOut, err := milenage.F2345(opc, k, rand)
	if err != nil {
		return err
	}

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
		copy(ak, akOut) // AK from f5 - used for AUTN
	}
	if akstar != nil {
		copy(akstar, akstarOut) // AK* from f5* - used for AUTS
	}

	return nil
}