//go:binary-only-package

package milenage

import (
	"fmt"
	"free5gc/lib/aes"
	"free5gc/lib/core_aes"
	"reflect"
	"strconv"
)

func aes128EncryptBlock(key, in, out []uint8) int {}

/*
int aes_128_encrypt_block(const c_uint8_t *key,
const c_uint8_t *in, c_uint8_t *out)
{
const int key_bits = 128;
unsigned int rk[RKLENGTH(128)];
int nrounds;

nrounds = aes_setup_enc(rk, key, key_bits);
aes_encrypt(rk, nrounds, in, out);

return 0;
}*/

/**
 * milenage_f1 - Milenage f1 and f1* algorithms
 * @opc: OPc = 128-bit value derived from OP and K
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * @sqn: SQN = 48-bit sequence number
 * @amf: AMF = 16-bit authentication management field
 * @mac_a: Buffer for MAC-A = 64-bit network authentication code, or %NULL
 * @mac_s: Buffer for MAC-S = 64-bit resync authentication code, or %NULL
 * Returns: 0 on success, -1 on failure
 */
func milenageF1(opc, k, _rand, sqn, amf, mac_a, mac_s []uint8) int {}

/**
 * milenage_f2345 - Milenage f2, f3, f4, f5, f5* algorithms
 * @opc: OPc = 128-bit value derived from OP and K
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * @res: Buffer for RES = 64-bit signed response (f2), or %NULL
 * @ck: Buffer for CK = 128-bit confidentiality key (f3), or %NULL
 * @ik: Buffer for IK = 128-bit integrity key (f4), or %NULL
 * @ak: Buffer for AK = 48-bit anonymity key (f5), or %NULL
 * @akstar: Buffer for AK = 48-bit anonymity key (f5*), or %NULL
 * Returns: 0 on success, -1 on failure
 */
func milenageF2345(opc, k, _rand, res, ck, ik, ak, akstar []uint8) int {}

func MilenageGenerate(opc, amf, k, sqn, _rand, autn, ik, ck, ak, res []uint8, res_len *uint) {}

/**
 * milenage_auts - Milenage AUTS validation
 * @opc: OPc = 128-bit operator variant algorithm configuration field (encr.)
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * @auts: AUTS = 112-bit authentication token from client
 * @sqn: Buffer for SQN = 48-bit sequence number
 * Returns: 0 = success (sqn filled), -1 on failure
 */
//int milenage_auts(const c_uint8_t *opc, const c_uint8_t *k, const c_uint8_t *_rand, const c_uint8_t *auts, c_uint8_t *sqn)
func Milenage_auts(opc, k, _rand, auts, sqn []uint8) int {}

/**
 * gsm_milenage - Generate GSM-Milenage (3GPP TS 55.205) authentication triplet
 * @opc: OPc = 128-bit operator variant algorithm configuration field (encr.)
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * @sres: Buffer for SRES = 32-bit SRES
 * @kc: Buffer for Kc = 64-bit Kc
 * Returns: 0 on success, -1 on failure
 */
func Gsm_milenage(opc, k, _rand, sres, kc []uint8) int {}

/**
 * milenage_generate - Generate AKA AUTN,IK,CK,RES
 * @opc: OPc = 128-bit operator variant algorithm configuration field (encr.)
 * @k: K = 128-bit subscriber key
 * @sqn: SQN = 48-bit sequence number
 * @_rand: RAND = 128-bit random challenge
 * @autn: AUTN = 128-bit authentication token
 * @ik: Buffer for IK = 128-bit integrity key (f4), or %NULL
 * @ck: Buffer for CK = 128-bit confidentiality key (f3), or %NULL
 * @res: Buffer for RES = 64-bit signed response (f2), or %NULL
 * @res_len: Variable that will be set to RES length
 * @auts: 112-bit buffer for AUTS
 * Returns: 0 on success, -1 on failure, or -2 on synchronization failure
 */
func Milenage_check(opc, k, sqn, _rand, autn, ik, ck, res []uint8, res_len *uint, auts []uint8) int {}

func milenage_opc(k, op, opc []uint8) {}

// implementation of os_memcmp
func os_memcmp(a, b []uint8, num int) int {}

func F1_Test(opc, k, _rand, sqn, amf, mac_a, mac_s []uint8) int {}

func F2345_Test(opc, k, _rand, res, ck, ik, ak, akstar []uint8) int {}

func GenerateOPC(k, op, opc []uint8) {}

func InsertData(op, k, _rand, sqn, amf []uint8, OP, K, RAND, SQN, AMF string) {}
