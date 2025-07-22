package suci

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/bits"
	"regexp"
	"strconv"
	"strings"

	"github.com/free5gc/udm/internal/logger"
)

// suci-0(SUPI type: IMSI)-mcc-mnc-routingIndicator-protectionScheme-homeNetworkPublicKeyID-schemeOutput.
// TODO: suci-1(SUPI type: NAI)-homeNetworkID-routingIndicator-protectionScheme-homeNetworkPublicKeyID-schemeOutput.

const (
	PrefixIMSI     = "imsi-"
	PrefixSUCI     = "suci"
	SupiTypeIMSI   = "0"
	NullScheme     = "0"
	ProfileAScheme = "1"
	ProfileBScheme = "2"
)

var (
	// Network and identification patterns.
	// Mobile Country Code; 3 digits
	mccRegex = `(?P<mcc>\d{3})`
	// Mobile Network Code; 2 or 3 digits
	mncRegex = `(?P<mnc>\d{2,3})`

	// MCC-MNC
	imsiTypeRegex = fmt.Sprintf("(?P<imsiType>0-%s-%s)", mccRegex, mncRegex)

	// The Home Network Identifier consists of a string of
	// characters with a variable length representing a domain name
	// as specified in Section 2.2 of RFC 7542
	naiTypeRegex = "(?P<naiType>1-.*)"

	// SUPI type; 0 = IMSI, 1 = NAI (for n3gpp)
	supiTypeRegex = fmt.Sprintf("(?P<supi_type>%s|%s)",
		imsiTypeRegex,
		naiTypeRegex)

	// Routing Indicator, used by the AUSF to find the appropriate UDM when SUCI is encrypted 1-4 digits
	routingIndicatorRegex = `(?P<routing_indicator>\d{1,4})`
	// Protection Scheme ID; 0 = NULL Scheme (unencrypted), 1 = Profile A, 2 = Profile B
	protectionSchemeRegex = `(?P<protection_scheme_id>(?:[0-2]))`
	// Public Key ID; 1-255
	publicKeyIDRegex = `(?P<public_key_id>(?:\d{1,2}|1\d{2}|2[0-4]\d|25[0-5]))`
	// Scheme Output; unbounded hex string (safe from ReDoS due to bounded length of SUCI)
	schemeOutputRegex = `(?P<scheme_output>[A-Fa-f0-9]+)`
	// Subscription Concealed Identifier (SUCI) Encrypted SUPI as sent by the UE to the AMF; 3GPP TS 29.503 - Annex C
	suciRegex = regexp.MustCompile(fmt.Sprintf("^suci-%s-%s-%s-%s-%s$",
		supiTypeRegex,
		routingIndicatorRegex,
		protectionSchemeRegex,
		publicKeyIDRegex,
		schemeOutputRegex,
	))
)

type Suci struct {
	SupiType         string // 0 for IMSI, 1 for NAI
	Mcc              string // 3 digits
	Mnc              string // 2-3 digits
	HomeNetworkId    string // variable-length string
	RoutingIndicator string // 1-4 digits
	ProtectionScheme string // 0-2
	PublicKeyID      string // 1-255
	SchemeOutput     string // hex string
}

func parseSuci(input string) *Suci {
	matches := suciRegex.FindStringSubmatch(input)
	if matches == nil || len(matches) != 10 {
		return nil
	}

	// The indices correspond to the order of the regex groups in the pattern
	return &Suci{
		SupiType:         matches[1], // First capture group
		Mcc:              matches[3], // Third capture group
		Mnc:              matches[4], // Fourth capture group
		HomeNetworkId:    matches[5], // Fifth capture group
		RoutingIndicator: matches[6], // Sixth capture group
		ProtectionScheme: matches[7], // Seventh capture group
		PublicKeyID:      matches[8], // Eighth capture group
		SchemeOutput:     matches[9], // Ninth capture group
	}
}

type SuciProfile struct {
	ProtectionScheme string `yaml:"ProtectionScheme,omitempty"`
	PrivateKey       string `yaml:"PrivateKey,omitempty"`
	PublicKey        string `yaml:"PublicKey,omitempty"`
}

// profile A.
const (
	ProfileAMacKeyLen = 32 // octets
	ProfileAEncKeyLen = 16 // octets
	ProfileAIcbLen    = 16 // octets
	ProfileAMacLen    = 8  // octets
	ProfileAHashLen   = 32 // octets
)

// profile B.
const (
	ProfileBMacKeyLen = 32 // octets
	ProfileBEncKeyLen = 16 // octets
	ProfileBIcbLen    = 16 // octets
	ProfileBMacLen    = 8  // octets
	ProfileBHashLen   = 32 // octets
)

func HmacSha256(input, macKey []byte, macLen int) ([]byte, error) {
	h := hmac.New(sha256.New, macKey)
	if _, err := h.Write(input); err != nil {
		return nil, fmt.Errorf("HMAC SHA256 error: %w", err)
	}
	macVal := h.Sum(nil)
	return macVal[:macLen], nil
}

func Aes128ctr(input, encKey, icb []byte) ([]byte, error) {
	output := make([]byte, len(input))
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, fmt.Errorf("AES128 CTR error: %w", err)
	}
	stream := cipher.NewCTR(block, icb)
	stream.XORKeyStream(output, input)
	return output, nil
}

func AnsiX963KDF(sharedKey, publicKey []byte, encKeyLen, macKeyLen, hashLen int) []byte {
	var counter uint32 = 1
	var kdfKey []byte
	kdfRounds := int(math.Ceil(float64(encKeyLen+macKeyLen) / float64(hashLen)))
	for i := 0; i < kdfRounds; i++ {
		counterBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(counterBytes, counter)
		tmpK := sha256.Sum256(append(append(sharedKey, counterBytes...), publicKey...))
		kdfKey = append(kdfKey, tmpK[:]...)
		counter++
	}
	return kdfKey
}

func swapNibbles(input []byte) []byte {
	output := make([]byte, len(input))
	for i, b := range input {
		output[i] = bits.RotateLeft8(b, 4)
	}
	return output
}

func calcSchemeResult(decryptPlainText []byte, supiType string) string {
	var result string
	if supiType == SupiTypeIMSI {
		result = hex.EncodeToString(swapNibbles(decryptPlainText))
		if len(result) > 0 && result[len(result)-1] == 'f' {
			result = result[:len(result)-1]
		}
	} else {
		result = hex.EncodeToString(decryptPlainText)
	}
	return result
}

func decryptWithKdf(sharedKey, kdfPubKey, cipherText, providedMac []byte,
	encKeyLen, macKeyLen, hashLen, icbLen, macLen int,
) ([]byte, error) {
	kdfKey := AnsiX963KDF(sharedKey, kdfPubKey, encKeyLen, macKeyLen, hashLen)
	encKey := kdfKey[:encKeyLen]
	icb := kdfKey[encKeyLen : encKeyLen+icbLen]
	macKey := kdfKey[len(kdfKey)-macKeyLen:]

	computedMac, err := HmacSha256(cipherText, macKey, macLen)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(computedMac, providedMac) {
		return nil, fmt.Errorf("decryption MAC failed")
	}
	logger.SuciLog.Infoln("decryption MAC match")

	return Aes128ctr(cipherText, encKey, icb)
}

func ecdhX25519(privateKeyHex string, peerPubKey []byte) ([]byte, error) {
	aHNPrivBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode X25519 private key: %w", err)
	}
	x25519Curve := ecdh.X25519()
	priv, err := x25519Curve.NewPrivateKey(aHNPrivBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X25519 private key: %w", err)
	}
	pub, err := x25519Curve.NewPublicKey(peerPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X25519 public key: %w", err)
	}
	return priv.ECDH(pub)
}

var ErrorPublicKeyUnmarshalling = fmt.Errorf("failed to unmarshal uncompressed public key")

func ecdhP256(privateKeyHex string, transmittedPubKey []byte) (sharedKey, kdfPubKey []byte, err error) {
	bHNPrivBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode P-256 private key: %w", err)
	}
	p256Curve := ecdh.P256()
	priv, err := p256Curve.NewPrivateKey(bHNPrivBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse P-256 private key: %w", err)
	}

	var pubKeyForECDH []byte
	switch transmittedPubKey[0] {
	case 0x02, 0x03:
		// Compressed format
		x, y := elliptic.UnmarshalCompressed(elliptic.P256(), transmittedPubKey)
		if x == nil || y == nil {
			return nil, nil, fmt.Errorf("failed to uncompress public key")
		}
		pubKeyForECDH = elliptic.Marshal(elliptic.P256(), x, y)
		kdfPubKey = transmittedPubKey

	case 0x04:
		// Uncompressed format.
		pubKeyForECDH = transmittedPubKey

		// For KDF, we need the compressed form.
		x, y := elliptic.Unmarshal(elliptic.P256(), transmittedPubKey)
		if x == nil || y == nil {
			return nil, nil, ErrorPublicKeyUnmarshalling
		}
		kdfPubKey = elliptic.MarshalCompressed(elliptic.P256(), x, y)
	default:
		return nil, nil, fmt.Errorf("unknown public key format")
	}

	pub, err := p256Curve.NewPublicKey(pubKeyForECDH)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create P-256 public key: %w", err)
	}

	sharedKey, err = priv.ECDH(pub)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compute ECDH: %w", err)
	}

	return sharedKey, kdfPubKey, nil
}

func profileA(input, supiType, privateKey string) (string, error) {
	logger.SuciLog.Infoln("SuciToSupi Profile A")

	s, err := hex.DecodeString(input)
	if err != nil {
		logger.SuciLog.Errorln("hex DecodeString error:", err)
		return "", err
	}

	const ProfileAPubKeyLen = 32
	if len(s) < ProfileAPubKeyLen+ProfileAMacLen {
		return "", fmt.Errorf("suci input too short")
	}

	peerPubKey := s[:ProfileAPubKeyLen]
	cipherText := s[ProfileAPubKeyLen : len(s)-ProfileAMacLen]
	providedMac := s[len(s)-ProfileAMacLen:]

	sharedKey, err := ecdhX25519(privateKey, peerPubKey)
	if err != nil {
		return "", err
	}

	plainText, err := decryptWithKdf(sharedKey, peerPubKey, cipherText, providedMac,
		ProfileAEncKeyLen, ProfileAMacKeyLen, ProfileAHashLen, ProfileAIcbLen, ProfileAMacLen)
	if err != nil {
		return "", err
	}
	return calcSchemeResult(plainText, supiType), nil
}

func profileB(input, supiType, privateKey string) (string, error) {
	logger.SuciLog.Infoln("SuciToSupi Profile B")

	s, err := hex.DecodeString(input)
	if err != nil || len(s) < 1 {
		return "", fmt.Errorf("hex DecodeString error: %w", err)
	}

	var ProfileBPubKeyLen int
	switch s[0] {
	case 0x02, 0x03:
		ProfileBPubKeyLen = 33
	case 0x04:
		ProfileBPubKeyLen = 65
	default:
		return "", fmt.Errorf("suci input error: unknown public key format")
	}

	if len(s) < ProfileBPubKeyLen+ProfileBMacLen {
		return "", fmt.Errorf("suci input too short")
	}

	transmittedPubKey := s[:ProfileBPubKeyLen]
	cipherText := s[ProfileBPubKeyLen : len(s)-ProfileBMacLen]
	providedMac := s[len(s)-ProfileBMacLen:]

	sharedKey, kdfPubKey, err := ecdhP256(privateKey, transmittedPubKey)
	if err != nil {
		return "", err
	}

	plainText, err := decryptWithKdf(sharedKey, kdfPubKey, cipherText, providedMac,
		ProfileBEncKeyLen, ProfileBMacKeyLen, ProfileBHashLen, ProfileBIcbLen, ProfileBMacLen)
	if err != nil {
		return "", err
	}
	return calcSchemeResult(plainText, supiType), nil
}

func ToSupi(suci string, suciProfiles []SuciProfile) (string, error) {
	parsedSuci := parseSuci(suci)
	if parsedSuci == nil {
		if strings.HasPrefix(suci, "imsi-") || strings.HasPrefix(suci, "nai-") {
			logger.SuciLog.Infof("Got supi\n")
			return suci, nil
		}
		return "", fmt.Errorf("unknown suci [%s]", suci)
	}

	logger.SuciLog.Infof("scheme %s", parsedSuci.ProtectionScheme)
	scheme := parsedSuci.ProtectionScheme
	mccMnc := parsedSuci.Mcc + parsedSuci.Mnc
	supiPrefix := PrefixIMSI

	if !strings.HasPrefix(parsedSuci.SupiType, SupiTypeIMSI) {
		logger.SuciLog.Infof("SUPI type is NAI")
		return "", fmt.Errorf("unsupported suciType NAI")
	}
	logger.SuciLog.Infof("SUPI type is IMSI")

	if scheme == NullScheme {
		return supiPrefix + mccMnc + parsedSuci.SchemeOutput, nil
	}

	keyIndex, err := strconv.Atoi(parsedSuci.PublicKeyID)
	if err != nil {
		return "", fmt.Errorf("parse HNPublicKeyID error: %w", err)
	}
	if keyIndex < 1 || keyIndex > len(suciProfiles) {
		return "", fmt.Errorf("keyIndex (%d) out of range (%d)", keyIndex, len(suciProfiles))
	}

	profile := suciProfiles[keyIndex-1]
	if scheme != profile.ProtectionScheme {
		return "", fmt.Errorf("protect Scheme mismatch [%s:%s]", scheme, profile.ProtectionScheme)
	}

	switch scheme {
	case ProfileAScheme:
		result, err := profileA(parsedSuci.SchemeOutput, SupiTypeIMSI, profile.PrivateKey)
		if err != nil {
			return "", err
		}
		return supiPrefix + mccMnc + result, nil
	case ProfileBScheme:
		result, err := profileB(parsedSuci.SchemeOutput, SupiTypeIMSI, profile.PrivateKey)
		if err != nil {
			return "", err
		}
		return supiPrefix + mccMnc + result, nil
	default:
		return "", fmt.Errorf("protect Scheme (%s) is not supported", scheme)
	}
}
