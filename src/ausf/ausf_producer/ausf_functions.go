package ausf_producer

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/antihax/optional"
	"github.com/bronze1man/radius"
	"free5gc/lib/Nnrf_NFDiscovery"
	Nudm_UEAU "free5gc/lib/Nudm_UEAuthentication"
	"free5gc/lib/openapi/models"
	"free5gc/src/ausf/ausf_consumer"
	"free5gc/src/ausf/ausf_context"
	"free5gc/src/ausf/logger"
	"hash"
	"strconv"
	"time"
)

func KDF5gAka(param ...string) hash.Hash {
	s := param[0]
	s += param[1]
	p0len, _ := strconv.Atoi(param[2])
	s += strconv.FormatInt(int64(p0len), 16)
	h := hmac.New(sha256.New, []byte(s))

	// Test data
	/* s2, _ := hex.DecodeString("4a656665")
	h := hmac.New(sha256.New, s2) */

	return h
}

func intToByteArray(i int) []byte {
	var r = make([]byte, 2)
	binary.BigEndian.PutUint16(r, uint16(i))
	return r
}

func padZeros(byteArray []byte, size int) []byte {
	l := len(byteArray)
	if l == size {
		return byteArray
	}
	r := make([]byte, size)
	copy(r[size-l:], byteArray)
	return r
}

func CalculateAtMAC(key []byte, input []byte) []byte {
	// keyed with K_aut
	h := hmac.New(sha256.New, key)
	if _, err := h.Write(input); err != nil {
		logger.EapAuthComfirmLog.Errorln(err.Error())
	}
	sha := string(h.Sum(nil))
	// fmt.Printf("[CalculateAtMAC] input: %x, key = %x\n sha: %x\n sha[:16]: %x\n", input, key, sha, sha[:16])
	return []byte(sha[:16])
}

func EapEncodeAttribute(attributeType string, data string) (returnStr string, err error) {
	var attribute string
	var length int
	var r []byte

	switch attributeType {
	case "AT_RAND":
		length = len(data)/8 + 1
		if length != 5 {
			return "", fmt.Errorf("[eapEncodeAttribute] AT_RAND Length Error")
		}
		attrNum := fmt.Sprintf("%02x", ausf_context.AT_RAND_ATTRIBUTE)
		attribute = attrNum + "05" + "0000" + data

	case "AT_AUTN":
		length = len(data)/8 + 1
		if length != 5 {
			return "", fmt.Errorf("[eapEncodeAttribute] AT_AUTN Length Error")
		}
		attrNum := fmt.Sprintf("%02x", ausf_context.AT_AUTN_ATTRIBUTE)
		attribute = attrNum + "05" + "0000" + data

	case "AT_KDF_INPUT":
		var byteName []byte
		nLength := len(data)
		length := (nLength+3)/4 + 1
		b := make([]byte, length*4)
		byteNameLength := intToByteArray(nLength)
		byteName = []byte(data)
		pad := padZeros(byteName, (length-1)*4)
		b[0] = 23
		b[1] = byte(length)
		copy(b[2:4], byteNameLength)
		copy(b[4:], pad)
		// fmt.Printf("AT_KDF_INPUT: %x\n", b[:])
		return string(b[:]), nil

	case "AT_KDF":
		// Value 1 default key derivation function for EAP-AKA'
		attrNum := fmt.Sprintf("%02x", ausf_context.AT_KDF_ATTRIBUTE)
		attribute = attrNum + "01" + "0001"

	case "AT_MAC":
		// Pad MAC value with 16 bytes of 0 since this is just for the calculation of MAC
		attrNum := fmt.Sprintf("%02x", ausf_context.AT_MAC_ATTRIBUTE)
		attribute = attrNum + "05" + "0000" + "00000000000000000000000000000000"

	case "AT_RES":
		var byteName []byte
		nLength := len(data)
		length := (nLength+3)/4 + 1
		b := make([]byte, length*4)
		byteNameLength := intToByteArray(nLength)
		byteName = []byte(data)
		pad := padZeros(byteName, (length-1)*4)
		b[0] = 3
		b[1] = byte(length)
		copy(b[2:4], byteNameLength)
		copy(b[4:], pad)
		return string(b[:]), nil

	default:
		logger.EapAuthComfirmLog.Errorf("UNKNOWN attributeType %s\n", attributeType)
		return "", nil
	}

	r, _ = hex.DecodeString(attribute)
	// fmt.Printf("%s: %x\n", attributeType, r)
	return string(r), nil
}

func eapAkaPrimePrf(ikPrime string, ckPrime string, identity string) (K_encr string, K_aut string, K_re string, MSK string, EMSK string) {
	keyAp := ikPrime + ckPrime

	// Test data
	// key, _ := hex.DecodeString("0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b")

	key, _ := hex.DecodeString(keyAp)
	sBase := []byte("EAP-AKA'" + identity)

	// Test data
	// sBase := []byte("Hi There")

	MK := ""
	prev := []byte("")
	_ = prev
	prfRounds := 208/32 + 1
	for i := 0; i < prfRounds; i++ {
		// Create a new HMAC by defining the hash type and the key (as byte array)
		h := hmac.New(sha256.New, key)

		hexNum := string(i + 1)
		ap := append(sBase, hexNum...)
		s := append(prev, ap...)

		// Write Data to it
		if _, err := h.Write([]byte(s)); err != nil {
			logger.EapAuthComfirmLog.Errorln(err.Error())
		}

		// Get result and encode as hexadecimal string
		sha := string(h.Sum(nil))
		MK += sha
		prev = []byte(sha)
		// fmt.Printf("MK(len %d): %s\n", len(MK), MK)
	}

	K_encr = MK[0:16]  // 0..127
	K_aut = MK[16:48]  // 128..383
	K_re = MK[48:80]   // 384..639
	MSK = MK[80:144]   // 640..1151
	EMSK = MK[144:208] // 1152..1663
	// fmt.Printf(" K_encr: %x\n K_aut: %x\n K_re: %x\n MSK: %x\n EMSK: %x\n", K_encr, K_aut, K_re, MSK, EMSK)
	return K_encr, K_aut, K_re, MSK, EMSK
}

func checkMACintegrity(offset int, expectedMacValue []byte, packet []byte, Kautn string) bool {
	eapDecode, decodeErr := radius.EapDecode(packet)
	if decodeErr != nil {
		logger.EapAuthComfirmLog.Infoln(decodeErr.Error())
	}
	zeroBytes, _ := hex.DecodeString("00000000000000000000000000000000")
	copy(eapDecode.Data[offset+4:offset+20], zeroBytes)
	encodeAfter := eapDecode.Encode()
	// fmt.Printf("check pkt with MAC val 0: %x, key = %x\n", encodeAfter, Kautn)
	MACvalue := CalculateAtMAC([]byte(Kautn), encodeAfter)
	// fmt.Printf("MAC value = %x\nExpected %x\n", MACvalue, expectedMacValue)

	if bytes.Equal(MACvalue, expectedMacValue) {
		return true
	} else {
		return false
	}
}

func decodeResMac(packetData []byte, wholePacket []byte, Kautn string) (RES []byte, success bool) {
	detectRes := false
	detectMac := false
	macCorrect := false
	dataArray := packetData
	var attributeLength int
	var attributeType int
	// fmt.Printf("packet's data: %x\n", dataArray)

	for i := 0; i < len(dataArray); i += attributeLength {
		attributeLength = int(uint(dataArray[1+i])) * 4
		attributeType = int(uint(dataArray[0+i]))

		if attributeType == ausf_context.AT_RES_ATTRIBUTE {
			logger.EapAuthComfirmLog.Infoln("Detect AT_RES attribute")
			detectRes = true
			resLength := int(uint(dataArray[3+i]) | uint(dataArray[2+i])<<8)
			RES = dataArray[4+i : 4+i+attributeLength-4]
			byteRes := padZeros(RES, resLength)
			RES = byteRes

		} else if attributeType == ausf_context.AT_MAC_ATTRIBUTE {
			logger.EapAuthComfirmLog.Infoln("Detect AT_MAC attribute")
			detectMac = true
			macStr := string(dataArray[4+i : 20+i])
			if checkMACintegrity(i, []byte(macStr), wholePacket, Kautn) {
				logger.EapAuthComfirmLog.Infoln("check MAC integrity succeed")
				macCorrect = true
			} else {
				logger.EapAuthComfirmLog.Infoln("check MAC integrity failed")
			}

		} else {
			logger.EapAuthComfirmLog.Infof("Detect unknown attribute with type %d\n", attributeType)
		}

	}

	if detectRes && detectMac && macCorrect {
		return RES, true
	}
	return nil, false
}

func ConstructFailEapAkaNotification(oldPktId uint8) string {
	var eapPkt radius.EapPacket
	eapPkt.Code = radius.EapCodeRequest
	eapPkt.Identifier = oldPktId + 1
	eapPkt.Type = ausf_context.EAP_AKA_PRIME_TYPENUM
	attrNum := fmt.Sprintf("%02x", ausf_context.AT_NOTIFICATION_ATTRIBUTE)
	// S bit = 0, P bit = 1 (0100 0000 0000 0000 = 16384)
	attribute := attrNum + "01" + "4000"
	attrHex, _ := hex.DecodeString(attribute)
	eapPkt.Data = attrHex
	eapPktEncode := eapPkt.Encode()
	return base64.StdEncoding.EncodeToString(eapPktEncode)
}

func ConstructEapNoTypePkt(code radius.EapCode, pktID uint8) string {
	b := make([]byte, 4)
	b[0] = byte(code)
	b[1] = byte(pktID)
	binary.BigEndian.PutUint16(b[2:4], uint16(4))
	// fmt.Printf("NoTypePkt: %x\n", b)
	return base64.StdEncoding.EncodeToString(b)
}

func getUdmUrl(nrfUri string) string {
	udmUrl := "https://localhost:29503" // default
	nfDiscoverParam := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		ServiceNames: optional.NewInterface([]models.ServiceName{models.ServiceName_NUDM_UEAU}),
	}
	res, err := ausf_consumer.SendSearchNFInstances(nrfUri, models.NfType_UDM, models.NfType_AUSF, nfDiscoverParam)
	if err != nil {
		logger.UeAuthPostLog.Errorln("[Search UDM UEAU] ", err.Error())
	} else if len(res.NfInstances) > 0 {
		udmInstance := res.NfInstances[0]
		if len(udmInstance.Ipv4Addresses) > 0 && udmInstance.NfServices != nil {
			ueauService := (*udmInstance.NfServices)[0]
			ueauEndPoint := (*ueauService.IpEndPoints)[0]
			udmUrl = string(ueauService.Scheme) + "://" + ueauEndPoint.Ipv4Address + ":" + strconv.Itoa(int(ueauEndPoint.Port))
		}
	} else {
		logger.UeAuthPostLog.Errorln("[Search UDM UEAU] len(NfInstances) = 0")
	}
	return udmUrl
}

func createClientToUdmUeau(udmUrl string) *Nudm_UEAU.APIClient {
	cfg := Nudm_UEAU.NewConfiguration()
	cfg.SetBasePath(udmUrl)
	clientAPI := Nudm_UEAU.NewAPIClient(cfg)
	return clientAPI
}

func sendAuthResultToUDM(id string, authType models.AuthType, success bool, servingNetworkName, udmUrl string) error {
	timeNow := time.Now()
	timePtr := &timeNow

	var authEvent models.AuthEvent
	authEvent.TimeStamp = timePtr
	authEvent.AuthType = authType
	authEvent.Success = success
	authEvent.ServingNetworkName = servingNetworkName

	client := createClientToUdmUeau(udmUrl)
	_, _, confirmAuthErr := client.ConfirmAuthApi.ConfirmAuth(context.Background(), id, authEvent)
	return confirmAuthErr
}

func logConfirmFailureAndInformUDM(id string, authType models.AuthType, servingNetworkName, errStr, udmUrl string) {
	if authType == models.AuthType__5_G_AKA {
		logger.Auth5gAkaComfirmLog.Infoln(errStr)
		if sendErr := sendAuthResultToUDM(id, authType, false, "", udmUrl); sendErr != nil {
			logger.Auth5gAkaComfirmLog.Infoln(sendErr.Error())
		}
	} else if authType == models.AuthType_EAP_AKA_PRIME {
		logger.EapAuthComfirmLog.Infoln(errStr)
		if sendErr := sendAuthResultToUDM(id, authType, false, "", udmUrl); sendErr != nil {
			logger.EapAuthComfirmLog.Infoln(sendErr.Error())
		}
	}
}
