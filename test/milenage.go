package test

import (
	"github.com/free5gc/util/milenage"
)

func MilenageF1(opc, k, rand, sqn, amf []byte, macA, macS []byte) error {

	mac_a, mac_s, err := milenage.F1(opc, k, rand, sqn, amf)
	if err != nil {
		return err
	}

	if macA != nil {
		copy(macA, mac_a)
	}

	if macS != nil {
		copy(macS, mac_s)
	}

	return nil
}

func MilenageF2345(opc, k, rand []byte, res, ck, ik, ak, akstar []byte) error {

	resOut, ckOut, ikOut, akOut, akstarOut, err := milenage.F2345(opc, k, rand)
	if err != nil {
		return err
	}

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
		copy(akstar, akstarOut)
	}

	return nil
}

