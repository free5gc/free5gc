package suci

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"
	"math"
	"math/big"
	"strings"

	"golang.org/x/crypto/curve25519"
)

// profile A
const ProfileAMacKeyLen = 32 // octets
const ProfileAEncKeyLen = 16 // octets
const ProfileAIcbLen = 16    // octets
const ProfileAMacLen = 8     // octets
const ProfileAHashLen = 32   // octets

// profile B
const ProfileBMacKeyLen = 32 // octets
const ProfileBEncKeyLen = 16 // octets
const ProfileBIcbLen = 16    // octets
const ProfileBMacLen = 8     // octets
const ProfileBHashLen = 32   // octets

func CompressKey(uncompressed []byte, y *big.Int) []byte {
	compressed := uncompressed[0:33]
	if y.Bit(0) == 1 { // 0x03
		compressed[0] = 0x03
	} else { // 0x02
		compressed[0] = 0x02
	}
	// fmt.Printf("compressed: %x\n", compressed)
	return compressed
}

// modified from https://stackoverflow.com/questions/46283760/how-to-uncompress-a-single-x9-62-compressed-point-on-an-ecdh-p256-curve-in-go
func uncompressKey(compressedBytes []byte, priv []byte) (*big.Int, *big.Int) {
	// Split the sign byte from the rest
	signByte := uint(compressedBytes[0])
	xBytes := compressedBytes[1:]

	x := new(big.Int).SetBytes(xBytes)
	three := big.NewInt(3)

	// The params for P256
	c := elliptic.P256().Params()

	// The equation is y^2 = x^3 - 3x + b
	// x^3, mod P
	xCubed := new(big.Int).Exp(x, three, c.P)

	// 3x, mod P
	threeX := new(big.Int).Mul(x, three)
	threeX.Mod(threeX, c.P)

	// x^3 - 3x + b mod P
	ySquared := new(big.Int).Sub(xCubed, threeX)
	ySquared.Add(ySquared, c.B)
	ySquared.Mod(ySquared, c.P)

	// find the square root mod P
	y := new(big.Int).ModSqrt(ySquared, c.P)
	if y == nil {
		// If this happens then you're dealing with an invalid point.
		logger.UeauLog.Errorln("Uncompressed key with invalid point")
		return nil, nil
	}

	// Finally, check if you have the correct root. If not you want -y mod P
	if y.Bit(0) != signByte&1 {
		y.Neg(y)
		y.Mod(y, c.P)
	}
	// fmt.Printf("xUncom: %x\nyUncon: %x\n", x, y)
	return x, y
}

func HmacSha256(input, macKey []byte, macLen int) []byte {
	h := hmac.New(sha256.New, macKey)
	_, _ = h.Write(input)
	macVal := h.Sum(nil)
	macTag := macVal[:macLen]
	// fmt.Printf("macVal: %x\nmacTag: %x\n", macVal, macTag)
	return macTag
}

func Aes128ctr(input, encKey, icb []byte) []byte {
	output := make([]byte, len(input))
	block, _ := aes.NewCipher(encKey)
	stream := cipher.NewCTR(block, icb)
	stream.XORKeyStream(output, input)
	// fmt.Printf("aes input: %x %x %x\naes output: %x\n", input, encKey, icb, output)
	return output
}

func AnsiX963KDF(sharedKey, publicKey []byte, profileEncKeyLen, profileMacKeyLen, profileHashLen int) []byte {
	var counter uint32 = 0x00000001
	var kdfKey []byte
	kdfRounds := int(math.Ceil(float64(profileEncKeyLen+profileMacKeyLen) / float64(profileHashLen)))
	for i := 1; i <= kdfRounds; i++ {
		counterBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(counterBytes, counter)
		// fmt.Printf("counterBytes: %x\n", counterBytes)
		tmpK := sha256.Sum256(append(append(sharedKey, counterBytes...), publicKey...))
		sliceK := tmpK[:]
		kdfKey = append(kdfKey, sliceK...)
		// fmt.Printf("kdfKey in round %d: %x\n", i, kdfKey)
		counter++
	}
	return kdfKey
}

func profileA(input string) (string, error) {
	logger.UeauLog.Infoln("SuciToSupi Profile A")
	s, hexDecodeErr := hex.DecodeString(input)
	if hexDecodeErr != nil {
		logger.UeauLog.Errorln("hex DecodeString error")
		return "", hexDecodeErr
	}

	// for X25519(profile A), q (The number of elements in the field Fq) = 2^255 - 19
	// len(pubkey) is therefore ceil((log2q)/8+1) = 32octets
	ProfileAPubKeyLen := 32
	if len(s) < ProfileAPubKeyLen+ProfileAMacLen {
		logger.UeauLog.Errorln("len of input data is too short!")
		return "", fmt.Errorf("suci input too short\n")
	}

	decryptMac := s[len(s)-ProfileAMacLen:]
	decryptPublicKey := s[:ProfileAPubKeyLen]
	decryptCipherText := s[ProfileAPubKeyLen : len(s)-ProfileAMacLen]
	// fmt.Printf("dePub: %x\ndeCiph: %x\ndeMac: %x\n", decryptPublicKey, decryptCipherText, decryptMac)

	// test data from TS33.501 Annex C.4
	// aHNPriv, _ := hex.DecodeString("c53c2208b61860b06c62e5406a7b330c2b577aa5558981510d128247d38bd1d")
	aHNPriv, _ := hex.DecodeString(udm_context.GetUdmProfileAHNPrivateKey())
	decryptSharedKey, _ := curve25519.X25519(aHNPriv, []byte(decryptPublicKey))
	// fmt.Printf("deShared: %x\n", decryptSharedKey)

	kdfKey := AnsiX963KDF(decryptSharedKey, decryptPublicKey, ProfileAEncKeyLen, ProfileAMacKeyLen, ProfileAHashLen)
	decryptEncKey := kdfKey[:ProfileAEncKeyLen]
	decryptIcb := kdfKey[ProfileAEncKeyLen : ProfileAEncKeyLen+ProfileAIcbLen]
	decryptMacKey := kdfKey[len(kdfKey)-ProfileAMacKeyLen:]
	// fmt.Printf("\ndeEncKey(size%d): %x\ndeMacKey: %x\ndeIcb: %x\n", len(decryptEncKey), decryptEncKey, decryptMacKey, decryptIcb)

	decryptMacTag := HmacSha256(decryptCipherText, decryptMacKey, ProfileAMacLen)
	if bytes.Equal(decryptMacTag, decryptMac) {
		logger.UeauLog.Infoln("decryption MAC match")
	} else {
		logger.UeauLog.Errorln("decryption MAC failed")
		return "", fmt.Errorf("decryption MAC failed\n")
	}

	decryptPlainText := Aes128ctr(decryptCipherText, decryptEncKey, decryptIcb)

	return hex.EncodeToString(decryptPlainText), nil
}

func profileB(input string) (string, error) {
	logger.UeauLog.Infoln("SuciToSupi Profile B")
	s, hexDecodeErr := hex.DecodeString(input)
	if hexDecodeErr != nil {
		logger.UeauLog.Errorln("hex DecodeString error")
		return "", hexDecodeErr
	}

	var ProfileBPubKeyLen int // p256, module q = 2^256 - 2^224 + 2^192 + 2^96 - 1
	var uncompressed bool
	if s[0] == 0x02 || s[0] == 0x03 {
		ProfileBPubKeyLen = 33 // ceil(log(2, q)/8) + 1 = 33
		uncompressed = false
	} else if s[0] == 0x04 {
		ProfileBPubKeyLen = 65 // 2*ceil(log(2, q)/8) + 1 = 65
		uncompressed = true
	} else {
		logger.UeauLog.Errorln("input error")
		return "", fmt.Errorf("suci input error\n")
	}

	// fmt.Printf("len:%d %d\n", len(s), ProfileBPubKeyLen + ProfileBMacLen)
	if len(s) < ProfileBPubKeyLen+ProfileBMacLen {
		logger.UeauLog.Errorln("len of input data is too short!")
		return "", fmt.Errorf("suci input too short\n")
	}
	decryptPublicKey := s[:ProfileBPubKeyLen]
	decryptMac := s[len(s)-ProfileBMacLen:]
	decryptCipherText := s[ProfileBPubKeyLen : len(s)-ProfileBMacLen]
	// fmt.Printf("dePub: %x\ndeCiph: %x\ndeMac: %x\n", decryptPublicKey, decryptCipherText, decryptMac)

	// test data from TS33.501 Annex C.4
	// bHNPriv, _ := hex.DecodeString("F1AB1074477EBCC7F554EA1C5FC368B1616730155E0041AC447D6301975FECDA")
	bHNPriv, _ := hex.DecodeString(udm_context.GetUdmProfileBHNPrivateKey())

	var xUncompressed, yUncompressed *big.Int
	if uncompressed {
		xUncompressed = new(big.Int).SetBytes(decryptPublicKey[1:(ProfileBPubKeyLen/2 + 1)])
		yUncompressed = new(big.Int).SetBytes(decryptPublicKey[(ProfileBPubKeyLen/2 + 1):])
	} else {
		xUncompressed, yUncompressed = uncompressKey(decryptPublicKey, bHNPriv)
		if xUncompressed == nil || yUncompressed == nil {
			logger.UeauLog.Errorln("Uncompressed key has invalid point")
			return "", fmt.Errorf("Key uncompression error\n")
		}
	}
	// fmt.Printf("xUncom: %x\nyUncom: %x\n", xUncompressed, yUncompressed)

	// x-coordinate is the shared key
	decryptSharedKey, _ := elliptic.P256().ScalarMult(xUncompressed, yUncompressed, bHNPriv)
	// fmt.Printf("deShared: %x\n", decryptSharedKey.Bytes())

	decryptPublicKeyForKDF := decryptPublicKey
	if uncompressed {
		decryptPublicKeyForKDF = CompressKey(decryptPublicKey, yUncompressed)
	}

	kdfKey := AnsiX963KDF(decryptSharedKey.Bytes(), decryptPublicKeyForKDF, ProfileBEncKeyLen, ProfileBMacKeyLen, ProfileBHashLen)
	// fmt.Printf("kdfKey: %x\n", kdfKey)
	decryptEncKey := kdfKey[:ProfileBEncKeyLen]
	decryptIcb := kdfKey[ProfileBEncKeyLen : ProfileBEncKeyLen+ProfileBIcbLen]
	decryptMacKey := kdfKey[len(kdfKey)-ProfileBMacKeyLen:]
	// fmt.Printf("\ndeEncKey(size%d): %x\ndeMacKey: %x\ndeIcb: %x\n", len(decryptEncKey), decryptEncKey, decryptMacKey, decryptIcb)

	decryptMacTag := HmacSha256(decryptCipherText, decryptMacKey, ProfileBMacLen)
	if bytes.Equal(decryptMacTag, decryptMac) {
		logger.UeauLog.Infoln("decryption MAC match")
	} else {
		logger.UeauLog.Errorln("decryption MAC failed")
		return "", fmt.Errorf("decryption MAC failed\n")
	}

	decryptPlainText := Aes128ctr(decryptCipherText, decryptEncKey, decryptIcb)

	return hex.EncodeToString(decryptPlainText), nil
}

// suci-0(SUPI type)-mcc-mnc-routingIndentifier-protectionScheme-homeNetworkPublicKeyIdentifier-schemeOutput
const supiTypePlace = 1
const mccPlace = 2
const mncPlace = 3
const schemePlace = 5

const typeIMSI = "0"
const imsiPrefix = "imsi-"
const profileAScheme = "1"
const profileBScheme = "2"

func ToSupi(suci string) (string, error) {
	suciPart := strings.Split(suci, "-")
	logger.UeauLog.Infof("suciPart %s\n", suciPart)

	suciPrefix := suciPart[0]
	if suciPrefix == "imsi" || suciPrefix == "nai" {
		logger.UeauLog.Infof("Got supi\n")
		return suci, nil

	} else if suciPrefix == "suci" {
		if len(suciPart) < 6 {
			logger.UeauLog.Errorf("Suci with wrong format\n")
			return suci, fmt.Errorf("Suci with wrong format\n")
		}

	} else {
		logger.UeauLog.Errorf("Unknown suciPrefix\n")
		return suci, fmt.Errorf("Unknown suciPrefix\n")
	}

	logger.UeauLog.Infof("scheme %s\n", suciPart[schemePlace])
	scheme := suciPart[schemePlace]
	mccMnc := suciPart[mccPlace] + suciPart[mncPlace]

	supiPrefix := imsiPrefix
	if suciPrefix == "suci" && suciPart[supiTypePlace] == typeIMSI {
		supiPrefix = imsiPrefix
		logger.UeauLog.Infof("SUPI type is IMSI\n")
	}

	if scheme == profileAScheme {
		profileAResult, err := profileA(suciPart[len(suciPart)-1])
		if err != nil {
			return "", err
		} else {
			return supiPrefix + mccMnc + profileAResult, nil
		}
	} else if scheme == profileBScheme {
		profileBResult, err := profileB(suciPart[len(suciPart)-1])
		if err != nil {
			return "", err
		} else {
			return supiPrefix + mccMnc + profileBResult, nil
		}
	} else { // NULL scheme
		return supiPrefix + mccMnc + suciPart[len(suciPart)-1], nil
	}
}
