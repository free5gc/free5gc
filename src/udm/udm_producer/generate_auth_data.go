package udm_producer

import (
	"context"
	"encoding/hex"

	// "fmt"
	"github.com/antihax/optional"
	// "free5gc/lib/CommonConsumerTestData/UDM/TestGenAuthData"
	"free5gc/lib/Nudr_DataRepository"
	"free5gc/lib/UeauCommon"
	"free5gc/lib/milenage"
	"free5gc/lib/openapi/common"
	"free5gc/lib/openapi/models"
	"free5gc/src/udm/logger"
	"free5gc/src/udm/udm_handler/udm_message"
	"math/rand"
	"net/http"
	"time"
)

func HandleGenerateAuthData(respChan chan udm_message.HandlerResponseMessage, supiOrSuci string, body models.AuthenticationInfoRequest) {
	var response models.AuthenticationInfoResult
	var problemDetails models.ProblemDetails
	rand.Seed(time.Now().UnixNano())

	client := createUDMClientToUDR(supiOrSuci, false)
	authSubs, _, err := client.AuthenticationDataDocumentApi.QueryAuthSubsData(context.Background(), supiOrSuci, nil)
	if err != nil {
		logger.UeauLog.Errorln("Return from UDR QueryAuthSubsData error")
		problemDetails.Cause = "AUTHENTICATION_REJECTED"
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, response)
	}

	/*
		K, RAND, CK, IK: 128 bits (16 bytes) (hex len = 32)
		SQN, AK: 48 bits (6 bytes) (hex len = 12) TS33.102 - 6.3.2
		AMF: 16 bits (2 bytes) (hex len = 4) TS33.102 - Annex H
	*/

	has_K, has_OP, has_OPC := false, false, false
	var K_str, OP_str, OPC_str string
	K, OP, OPC := make([]byte, 16), make([]byte, 16), make([]byte, 16)

	if authSubs.PermanentKey != nil {
		K_str = authSubs.PermanentKey.PermanentKeyValue
		K, _ = hex.DecodeString(K_str)
		has_K = true
	} else {
		logger.UeauLog.Errorln("Nil PermanentKey")
		problemDetails.Cause = "AUTHENTICATION_REJECTED"
		udm_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, response)
	}

	if authSubs.Milenage != nil {
		if authSubs.Milenage.Op != nil {
			OP_str = authSubs.Milenage.Op.OpValue
			OP, _ = hex.DecodeString(OP_str)
			has_OP = true
		} else {
			logger.UeauLog.Infoln("Nil Op")
		}
	} else {
		logger.UeauLog.Infoln("Nil Milenage")
	}

	if authSubs.Opc != nil {
		OPC_str = authSubs.Opc.OpcValue
		OPC, _ = hex.DecodeString(OPC_str)
		has_OPC = true
	} else {
		logger.UeauLog.Infoln("Nil Opc")
	}

	if !has_OPC {
		if has_K && has_OP {
			milenage.GenerateOPC(K, OP, OPC)
		} else {
			logger.UeauLog.Errorln("Unable to derive OPc")
			problemDetails.Cause = "AUTHENTICATION_REJECTED"
			udm_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, response)
		}
	}

	SQN_str := authSubs.SequenceNumber
	SQN, _ := hex.DecodeString(SQN_str)
	// fmt.Printf("K=%x\nSQN=%x\nOP=%x\nOPC=%x\n", K, SQN, OP, OPC)

	RAND := make([]byte, 16)
	rand.Read(RAND)
	AMF, _ := hex.DecodeString("8000")
	// fmt.Printf("RAND=%x\nAMF=%x\n", RAND, AMF)

	// for test
	// RAND, _ = hex.DecodeString(TestGenAuthData.MilenageTestSet19.RAND)
	// AMF, _ = hex.DecodeString(TestGenAuthData.MilenageTestSet19.AMF)
	// fmt.Printf("For test: RAND=%x, AMF=%x\n", RAND, AMF)

	// Run milenage
	MAC_A, MAC_S := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)

	// Generate MAC_A, MAC_S
	milenage.F1_Test(OPC, K, RAND, SQN, AMF, MAC_A, MAC_S)

	// Generate RES, CK, IK, AK, AKstar
	// RES == XRES (expected RES) for server
	milenage.F2345_Test(OPC, K, RAND, RES, CK, IK, AK, AKstar)
	// fmt.Printf("milenage RES = %s\n", hex.EncodeToString(RES))

	// Generate AUTN
	// fmt.Printf("SQN=%x\nAK =%x\n", SQN, AK)
	// fmt.Printf("AMF=%x, MAC_A=%x\n", AMF, MAC_A)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(SQN); i++ {
		SQNxorAK[i] = SQN[i] ^ AK[i]
	}
	// fmt.Printf("SQN xor AK = %x\n", SQNxorAK)
	AUTN := append(append(SQNxorAK, AMF...), MAC_A...)
	// fmt.Printf("AUTN = %x\n", AUTN)

	var av models.AuthenticationVector
	if authSubs.AuthenticationMethod == models.AuthMethod__5_G_AKA {
		response.AuthType = models.AuthType__5_G_AKA

		// derive XRES*
		key := append(CK, IK...)
		FC := UeauCommon.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
		P0 := []byte(body.ServingNetworkName)
		P1 := RAND
		P2 := RES

		kdfVal_for_xresStar := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1), P2, UeauCommon.KDFLen(P2))
		xresStar := kdfVal_for_xresStar[len(kdfVal_for_xresStar)/2:]
		// fmt.Printf("xresStar = %x\n", xresStar)

		// derive Kausf
		FC = UeauCommon.FC_FOR_KAUSF_DERIVATION
		P0 = []byte(body.ServingNetworkName)
		P1 = SQNxorAK
		kdfVal_for_Kausf := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1))
		// fmt.Printf("Kausf = %x\n", kdfVal_for_Kausf)

		// Fill in rand, xresStar, autn, kausf
		av.Rand = hex.EncodeToString(RAND)
		av.XresStar = hex.EncodeToString(xresStar)
		av.Autn = hex.EncodeToString(AUTN)
		av.Kausf = hex.EncodeToString(kdfVal_for_Kausf)

	} else { // EAP-AKA'
		response.AuthType = models.AuthType_EAP_AKA_PRIME

		// derive CK' and IK'
		key := append(CK, IK...)
		FC := UeauCommon.FC_FOR_CK_PRIME_IK_PRIME_DERIVATION
		P0 := []byte(body.ServingNetworkName)
		P1 := SQNxorAK
		kdfVal := UeauCommon.GetKDFValue(key, FC, P0, UeauCommon.KDFLen(P0), P1, UeauCommon.KDFLen(P1))
		// fmt.Printf("kdfVal = %x (len = %d)\n", kdfVal, len(kdfVal))

		// For TS 35.208 test set 19 & RFC 5448 test vector 1
		// CK': 0093 962d 0dd8 4aa5 684b 045c 9edf fa04
		// IK': ccfc 230c a74f cc96 c0a5 d611 64f5 a76

		ckPrime := kdfVal[:len(kdfVal)/2]
		ikPrime := kdfVal[len(kdfVal)/2:]
		// fmt.Printf("ckPrime: %x\nikPrime: %x\n", ckPrime, ikPrime)

		// Fill in rand, xres, autn, ckPrime, ikPrime
		av.Rand = hex.EncodeToString(RAND)
		av.Xres = hex.EncodeToString(RES)
		av.Autn = hex.EncodeToString(AUTN)
		av.CkPrime = hex.EncodeToString(ckPrime)
		av.IkPrime = hex.EncodeToString(ikPrime)
	}
	response.AuthenticationVector = &av
	response.Supi = supiOrSuci
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, response)
}

func HandleConfirmAuthData(respChan chan udm_message.HandlerResponseMessage, supi string, body models.AuthEvent) {
	var createAuthParam Nudr_DataRepository.CreateAuthenticationStatusParamOpts
	optInterface := optional.NewInterface(body)
	createAuthParam.AuthEvent = optInterface

	client := createUDMClientToUDR(supi, false)
	resp, err := client.AuthenticationStatusDocumentApi.CreateAuthenticationStatus(context.Background(), supi, &createAuthParam)
	if err != nil {
		logger.UeauLog.Errorln("[ConfirmAuth] ", err.Error())
		var problemDetails models.ProblemDetails
		problemDetails.Cause = err.(common.GenericOpenAPIError).Model().(models.ProblemDetails).Cause
		udm_message.SendHttpResponseMessage(respChan, nil, resp.StatusCode, problemDetails)
	}
	udm_message.SendHttpResponseMessage(respChan, nil, http.StatusCreated, nil)
}
