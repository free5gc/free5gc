package UEAuthentication

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/Nudm_UEAuthentication"
	"free5gc/lib/http2_util"
	"free5gc/lib/openapi/models"
	"free5gc/lib/path_util"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_context"
	"free5gc/src/udm/udm_handler"
	"free5gc/src/udm/udm_util"
	"math/big"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/curve25519"
)

func generateProfileAEphemeralKey() ([]byte, []byte) {
	privKey25519 := make([]byte, 32)
	_, _ = rand.Read(privKey25519)
	pubKey25519, _ := curve25519.X25519(privKey25519, curve25519.Basepoint)
	fmt.Printf("[EphemeralKey_profileA] privKey25519: %x\npubKey25519: %x\n", privKey25519, pubKey25519)
	return privKey25519, pubKey25519
}

func generateProfileBEphemeralKey() ([]byte, *big.Int, *big.Int) {
	p256 := elliptic.P256()
	priv, x, y, _ := elliptic.GenerateKey(p256, rand.Reader)
	fmt.Printf("[EphemeralKey_profileB] Private key: %x\nX of Public key: %x\nY of Public key: %x\n", priv, x, y)
	return priv, x, y
}

func generateTestData(sharedKey, pubEph, plaintext []byte, profileScheme int) []byte {
	var encKeyLen, macKeyLen, hashLen, icbLen, macLen int
	if profileScheme == 1 {
		encKeyLen = udm_util.ProfileAEncKeyLen
		macKeyLen = udm_util.ProfileAMacKeyLen
		hashLen = udm_util.ProfileAHashLen
		icbLen = udm_util.ProfileAIcbLen
		macLen = udm_util.ProfileAMacLen
	} else if profileScheme == 2 {
		encKeyLen = udm_util.ProfileBEncKeyLen
		macKeyLen = udm_util.ProfileBMacKeyLen
		hashLen = udm_util.ProfileBHashLen
		icbLen = udm_util.ProfileBIcbLen
		macLen = udm_util.ProfileBMacLen
	} else {
		return plaintext
	}

	kdfKey := udm_util.AnsiX963KDF(sharedKey, pubEph, encKeyLen, macKeyLen, hashLen)
	encKey := kdfKey[:encKeyLen]
	icb := kdfKey[encKeyLen : encKeyLen+icbLen]
	macKey := kdfKey[len(kdfKey)-macKeyLen:]
	// fmt.Printf("kdfKey: %x\nencKey: %x\nmacKey: %x\nicb: %x\n", kdfKey, encKey, macKey, icb)

	ciphertext := udm_util.Aes128ctr(plaintext, encKey, icb)
	macTag := udm_util.HmacSha256(ciphertext, macKey, macLen)
	// fmt.Printf("plain: %x\ncipher: %x\nmacTag: %x\n", plaintext, ciphertext, macTag)

	output := append(append(pubEph, ciphertext...), macTag...)
	fmt.Printf("scheme output: %x\n", output)
	return output
}

const profileAScheme = 1
const profileBScheme = 2

func TestUeAuthenticationsPost(t *testing.T) {
	go func() { // udm server
		router := gin.Default()
		AddService(router)

		udmLogPath := path_util.Gofree5gcPath("free5gc/udmsslkey.log")
		udmPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.pem")
		udmKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udm.key")

		server, err := http2_util.NewServer(":29503", udmLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udmPemPath, udmKeyPath))
			assert.True(t, err == nil)
		}
	}()
	// udm_util.testInitUdmConfig()
	udm_context.TestInit()
	go udm_handler.Handle()

	go func() { // fake udr server
		router := gin.Default()

		router.GET("/nudr-dr/v1/subscription-data/:ueId/authentication-data/authentication-subscription", func(c *gin.Context) {
			ueId := c.Param("ueId")
			fmt.Println("ueId: ", ueId)
			var authSubs models.AuthenticationSubscription
			var pk models.PermanentKey
			var opc models.Opc
			var var_milenage models.Milenage
			var op models.Op

			pk.PermanentKeyValue = TestGenAuthData.MilenageTestSet19.K
			opc.OpcValue = TestGenAuthData.MilenageTestSet19.OPC
			op.OpValue = TestGenAuthData.MilenageTestSet19.OP
			var_milenage.Op = &op

			authSubs.PermanentKey = &pk
			authSubs.Opc = &opc
			authSubs.Milenage = &var_milenage
			authSubs.SequenceNumber = TestGenAuthData.MilenageTestSet19.SQN
			authSubs.AuthenticationMethod = models.AuthMethod__5_G_AKA
			// authSubs.AuthenticationMethod = models.AuthMethod_EAP_AKA_PRIME

			c.JSON(http.StatusOK, authSubs)
		})

		udrLogPath := path_util.Gofree5gcPath("free5gc/udrsslkey.log")
		udrPemPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.pem")
		udrKeyPath := path_util.Gofree5gcPath("free5gc/support/TLS/udr.key")

		server, err := http2_util.NewServer(":29504", udrLogPath, router)
		if err == nil && server != nil {
			logger.InitLog.Infoln(server.ListenAndServeTLS(udrPemPath, udrKeyPath))
			assert.True(t, err == nil)
		}
	}()

	var authInfoReq models.AuthenticationInfoRequest
	authInfoReq.ServingNetworkName = TestGenAuthData.TestGenAuthDataTable[TestGenAuthData.SUCCESS_CASE].ServingNetworkName

	fmt.Printf("\n==========")
	fmt.Printf("Please modify variable \"profileScheme\" to test different schemes")
	fmt.Printf("==========\n\n")
	// Please modify HERE(profileScheme) to test different schemes!
	// profile: 0=>NULL scheme, 1=>Profile A, 2=>Profile B
	profileScheme := 1
	var testData []byte
	var supiOrSuci string
	// fill in the suci you want as plaintext here
	// test data from TS33.501 Annex C.4
	plaintext, _ := hex.DecodeString("00012080f6")
	// plaintext, _ := hex.DecodeString("aabb8787")

	// suci-0(SUPI type)-mcc-mnc-routingIndentifier-protectionScheme-homeNetworkPublicKeyIdentifier-schemeOutput
	suciTestPrefix := "suci-0-274-012-0001-" + strconv.Itoa(profileScheme) + "-01-"
	if profileScheme == profileAScheme {
		privEphProfileA, pubEphProfileA := generateProfileAEphemeralKey()
		pubHNProfileA, _ := hex.DecodeString(udm_context.GetUdmProfileAHNPublicKey())

		// test data from TS33.501 Annex C.4
		// privEphProfileA, _ = hex.DecodeString("c80949f13ebe61af4ebdbd293ea4f942696b9e815d7e8f0096bbf6ed7de62256")
		// pubEphProfileA, _ = hex.DecodeString("b2e92f836055a255837debf850b528997ce0201cb82adfe4be1f587d07d8457d")
		// pubHNProfileA, _ = hex.DecodeString("5a8d38864820197c3394b92613b20b91633cbd897119273bf8e4a6f4eec0a650")

		sharedKey, _ := curve25519.X25519(privEphProfileA, pubHNProfileA)
		fmt.Printf("[profileA] shared key: %x\n", sharedKey)

		testData = generateTestData(sharedKey, pubEphProfileA, plaintext, profileScheme)
		supiOrSuci = suciTestPrefix + hex.EncodeToString(testData)
	} else if profileScheme == profileBScheme {
		pubHNProfileBstr := udm_context.GetUdmProfileBHNPublicKey()
		privEphProfileB, pubEphProfileBx, pubEphProfileBy := generateProfileBEphemeralKey()
		pubEphProfileBstr := "04" + hex.EncodeToString(pubEphProfileBx.Bytes()) + hex.EncodeToString(pubEphProfileBy.Bytes())

		// test data from TS33.501 Annex C.4
		// privEphProfileB, _ := hex.DecodeString("99798858A1DC6A2C68637149A4B1DBFD1FDFF5ADDD62A2142F06699ED7602529")
		// pubEphProfileBxStr := "9AAB8376597021E855679A9778EA0B67396E68C66DF32C0F41E9ACCA2DA9B9D1"
		// pubEphProfileByStr := "F44EA1C87AA7478B954537BDE79951E748A43294A4Fde4CF86EAFF1789C9C81F"
		// pubEphProfileBy, _ := new(big.Int).SetString(pubEphProfileByStr, 16)
		// pubEphProfileBstr := "04" + pubEphProfileBxStr + pubEphProfileByStr
		// pubEphProfileBstr := "039AAB8376597021E855679A9778EA0B67396E68C66DF32C0F41E9ACCA2DA9B9D1"
		// pubEphProfileB, _ := hex.DecodeString(pubEphProfileBstr)
		// pubHNProfileBstr := "0472DA71976234CE833A6907425867B82E074D44EF907DFB4B3E21C1C2256EBCD15A7DED52FCBB097A4ED250E036C7B9C8C7004C4EEDC4F068CD7BF8D3F900E3B4"

		pubHNProfileB, _ := hex.DecodeString(pubHNProfileBstr)
		pubEphProfileBuncom, _ := hex.DecodeString(pubEphProfileBstr)
		pubEphProfileB := udm_util.CompressKey(pubEphProfileBuncom, pubEphProfileBy)
		// fmt.Printf("pubEphB: %x\npubHNB: %x\n", pubEphProfileB, pubHNProfileB)

		pubHNProfileBx := pubHNProfileB[1:33]
		pubHNProfileBy := pubHNProfileB[33:65]
		HNBx := new(big.Int).SetBytes(pubHNProfileBx)
		HNBy := new(big.Int).SetBytes(pubHNProfileBy)
		sharedKey, _ := elliptic.P256().ScalarMult(HNBx, HNBy, privEphProfileB)
		// fmt.Printf("HNBx = %x\nHNBy = %x\n", HNBx, HNBy)
		fmt.Printf("[Profile B] sharedKey = %x\n", sharedKey)

		testData = generateTestData(sharedKey.Bytes(), pubEphProfileB, plaintext, profileScheme)
		supiOrSuci = suciTestPrefix + hex.EncodeToString(testData)
	} else { // NULL scheme
		fmt.Printf("[NULL scheme]\n")
		testData = generateTestData(nil, nil, plaintext, profileScheme)
		supiOrSuci = suciTestPrefix + hex.EncodeToString(testData)
	}
	// fmt.Printf("testData = %x\n", testData)

	// for the profile test data from TS33501 Annex C.4:
	// ProfileASuci             = "0-123-45-0002-1-17-b2e92f836055a255837debf850b528997ce0201cb82adfe4be1f587d07d8457dcb02352410cddd9e730ef3fa87"
	// ProfileBSuciCompressed   = "0-123-45-0002-2-17-039AAB8376597021E855679A9778EA0B67396E68C66DF32C0F41E9ACCA2DA9B9D146A33FC2716AC7DAE96AA30A4D"
	// ProfileBSuciUncompressed = "0-123-45-0002-2-17-049AAB8376597021E855679A9778EA0B67396E68C66DF32C0F41E9ACCA2DA9B9D1D1F44EA1C87AA7478B954537BDE79951E748A43294A4F4CF86EAFF1789C9C81F46A33FC2716AC7DAE96AA30A4D"
	fmt.Printf("supiOrSuci = %s\n", supiOrSuci)

	cfg := Nudm_UEAuthentication.NewConfiguration()
	cfg.SetBasePath("https://localhost:29503")
	client := Nudm_UEAuthentication.NewAPIClient(cfg)

	authInfoRes, resp, err := client.GenerateAuthDataApi.GenerateAuthData(context.TODO(), supiOrSuci, authInfoReq)
	fmt.Println("=====")
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		fmt.Println("resp: ", resp)
		fmt.Println("authInfoRes: ", authInfoRes)
	}

	switch authInfoRes.AuthType {
	case models.AuthType__5_G_AKA:
		fmt.Printf("auth type: 5G AKA\n")
		// rand, xresStar, autn, kausf
		av := authInfoRes.AuthenticationVector
		if av != nil {
			fmt.Printf("rand: %s\nxresStar: %s\nautn: %s\nkausf: %s\n", av.Rand, av.XresStar, av.Autn, av.Kausf)
		} else {
			fmt.Printf("nil av\n")
		}

	case models.AuthType_EAP_AKA_PRIME:
		fmt.Printf("auth type: EAP-AKA'\n")
		// rand, xres, autn, ckPrime, ikPrime
		av := authInfoRes.AuthenticationVector
		if av != nil {
			fmt.Printf("rand: %s\nxres: %s\nautn: %s\nCK': %s\nIK': %s", av.Rand, av.Xres, av.Autn, av.CkPrime, av.IkPrime)
		} else {
			fmt.Printf("nil av\n")
		}

	default:
		fmt.Println("authInfoRes authType error")
	}
}
