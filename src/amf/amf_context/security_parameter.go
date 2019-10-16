package amf_context

// TS 33501 Annex A.8 Algorithm distinguisher For Knas_int Knas_enc
const (
	N_NAS_ENC_ALG uint8 = 0x01
	N_NAS_INT_ALG uint8 = 0x02
	N_RRC_ENC_ALG uint8 = 0x03
	N_RRC_INT_ALG uint8 = 0x04
	N_UP_ENC_alg  uint8 = 0x05
	N_UP_INT_alg  uint8 = 0x06
)

// TS 33501 Annex D Algorithm identifier values For Knas_int
const (
	ALG_INTEGRITY_128_NIA0 uint8 = 0x00 // NULL
	ALG_INTEGRITY_128_NIA1 uint8 = 0x01 // 128-Snow3G
	ALG_INTEGRITY_128_NIA2 uint8 = 0x02 // 128-AES
	ALG_INTEGRITY_128_NIA3 uint8 = 0x03 // 128-ZUC
)

// TS 33501 Annex D Algorithm identifier values For Knas_enc
const (
	ALG_CIPHERING_128_NEA0 uint8 = 0x00 // NULL
	ALG_CIPHERING_128_NEA1 uint8 = 0x01 // 128-Snow3G
	ALG_CIPHERING_128_NEA2 uint8 = 0x02 // 128-AES
	ALG_CIPHERING_128_NEA3 uint8 = 0x03 // 128-ZUC
)

// 1bit
const (
	SECURITY_DIRECTION_UPLINK   uint8 = 0x00
	SECURITY_DIRECTION_DOWNLINK uint8 = 0x01
)

// 5bits
const (
	SECURITY_ONLY_ONE_BEARER uint8 = 0x00
	SECURITY_BEARER_3GPP     uint8 = 0x01
	SECURITY_BEARER_NON_3GPP uint8 = 0x02
)

// TS 33501 Annex A.0 Access type distinguisher For Kgnb Kn3iwf
const (
	ACCESS_TYPE_3GPP     uint8 = 0x01
	ACCESS_TYPE_NON_3GPP uint8 = 0x02
)
