package test

import (
	"encoding/binary"
	"encoding/hex"
	"free5gc/lib/UeauCommon"
	"free5gc/lib/milenage"
	"free5gc/lib/openapi/models"
)

type RanUeContext struct {
	Supi               string
	RanUeNgapId        int64
	AmfUeNgapId        int64
	ULCount            uint32
	DLOverflow         uint16
	DLCountSQN         uint8
	CipheringAlg       uint8
	IntegrityAlg       uint8
	KnasEnc            []uint8
	KnasInt            []uint8
	Kamf               []uint8
	AuthenticationSubs models.AuthenticationSubscription
}

func NewRanUeContext(supi string, ranUeNgapId int64, cipheringAlg, integrityAlg uint8) *RanUeContext {
	ue := RanUeContext{}
	ue.RanUeNgapId = ranUeNgapId
	ue.Supi = supi
	ue.CipheringAlg = cipheringAlg
	ue.IntegrityAlg = integrityAlg
	return &ue
}
func (ue *RanUeContext) GetSecurityULCount() []byte {
	var r = make([]byte, 4)
	binary.BigEndian.PutUint32(r, ue.ULCount&0xffffff)
	return r
}

func (ue *RanUeContext) GetSecurityDLCount() []byte {
	var r = make([]byte, 4)
	binary.BigEndian.PutUint16(r, ue.DLOverflow)
	r[3] = ue.DLCountSQN
	r[2] = r[1]
	r[1] = r[0]
	r[0] = 0x00
	return r
}

func (ue *RanUeContext) DeriveRESstarAndSetKey(authSubs models.AuthenticationSubscription, RAND []byte, snNmae string) []byte {

	SQN, _ := hex.DecodeString(authSubs.SequenceNumber)

	AMF, _ := hex.DecodeString(authSubs.AuthenticationManagementField)

	// Run milenage
	MAC_A, MAC_S := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)
	OPC, _ := hex.DecodeString(authSubs.Opc.OpcValue)
	K, _ := hex.DecodeString(authSubs.PermanentKey.PermanentKeyValue)
	// Generate MAC_A, MAC_S
	milenage.F1_Test(OPC, K, RAND, SQN, AMF, MAC_A, MAC_S)

	// Generate RES, CK, IK, AK, AKstar
	milenage.F2345_Test(OPC, K, RAND, RES, CK, IK, AK, AKstar)

	// derive RES*
	key := append(CK, IK...)
	FC := UeauCommon.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
	P0 := []byte(snNmae)
	P1 := RAND
	P2 := RES

	ue.DerivateKamf(key, snNmae, SQN, AK)
	ue.DerivateAlgKey()
	kdfVal_for_resStar := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1), P2, UeauCommon.KDFLen(P2))
	return kdfVal_for_resStar[len(kdfVal_for_resStar)/2:]

}

func (ue *RanUeContext) DerivateKamf(key []byte, snName string, SQN, AK []byte) {

	FC := UeauCommon.FC_FOR_KAUSF_DERIVATION
	P0 := []byte(snName)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(SQN); i++ {
		SQNxorAK[i] = SQN[i] ^ AK[i]
	}
	P1 := SQNxorAK
	Kausf := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1))
	P0 = []byte(snName)
	Kseaf := UeauCommon.GetKDFValue(Kausf, UeauCommon.FC_FOR_KSEAF_DERIVATION, P0, UeauCommon.KDFLen(P0))

	P0, _ = hex.DecodeString(ue.Supi)
	L0 := UeauCommon.KDFLen(P0)
	P1 = []byte{0x00, 0x00}
	L1 := UeauCommon.KDFLen(P1)

	ue.Kamf = UeauCommon.GetKDFValue(Kseaf, UeauCommon.FC_FOR_KAMF_DERIVATION, P0, L0, P1, L1)
}

// Algorithm key Derivation function defined in TS 33.501 Annex A.9
func (ue *RanUeContext) DerivateAlgKey() {
	// Security Key
	P0 := []byte{N_NAS_ENC_ALG}
	L0 := UeauCommon.KDFLen(P0)
	P1 := []byte{ue.CipheringAlg}
	L1 := UeauCommon.KDFLen(P1)

	kenc := UeauCommon.GetKDFValue(ue.Kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	ue.KnasEnc = kenc[16:32]

	// Integrity Key
	P0 = []byte{N_NAS_INT_ALG}
	L0 = UeauCommon.KDFLen(P0)
	P1 = []byte{ue.IntegrityAlg}
	L1 = UeauCommon.KDFLen(P1)

	kint := UeauCommon.GetKDFValue(ue.Kamf, UeauCommon.FC_FOR_ALGORITHM_KEY_DERIVATION, P0, L0, P1, L1)
	ue.KnasInt = kint[16:32]
}
