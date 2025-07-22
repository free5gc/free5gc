package processor

import (
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	Nudr_DataRepository "github.com/free5gc/openapi/udr/DataRepository"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/udm/internal/util"
	"github.com/free5gc/udm/pkg/suci"
	"github.com/free5gc/util/ueauth"
)

const (
	SqnMAx    int64 = 0xFFFFFFFFFFFF
	ind       int64 = 32
	keyStrLen int   = 32
	opStrLen  int   = 32
	opcStrLen int   = 32
)

const (
	authenticationRejected string = "AUTHENTICATION_REJECTED"
	resyncAMF              string = "0000"
)

func (p *Processor) aucSQN(opc, k, auts, rand []byte) ([]byte, []byte) {
	AK, SQNms := make([]byte, 6), make([]byte, 6)
	macS := make([]byte, 8)
	ConcSQNms := auts[:6]
	AMF, err := hex.DecodeString(resyncAMF)
	if err != nil {
		return nil, nil
	}

	logger.UeauLog.Tracef("aucSQN: ConcSQNms=[%x]", ConcSQNms)

	err = util.MilenageF2345(opc, k, rand, nil, nil, nil, nil, AK)
	if err != nil {
		logger.UeauLog.Errorln("aucSQN milenage F2345 err:", err)
	}

	for i := 0; i < 6; i++ {
		SQNms[i] = AK[i] ^ ConcSQNms[i]
	}

	logger.UeauLog.Tracef("aucSQN: opc=[%x], k=[%x], rand=[%x], AMF=[%x], SQNms=[%x]\n", opc, k, rand, AMF, SQNms)
	// The AMF used to calculate MAC-S assumes a dummy value of all zeros
	err = util.MilenageF1(opc, k, rand, SQNms, AMF, nil, macS)
	if err != nil {
		logger.UeauLog.Errorln("aucSQN milenage F1 err:", err)
	}
	logger.UeauLog.Tracef("aucSQN: macS=[%x]\n", macS)
	return SQNms, macS
}

func (p *Processor) strictHex(ss string, n int) string {
	l := len(ss)
	if l < n {
		return strings.Repeat("0", n-l) + ss
	} else {
		return ss[l-n : l]
	}
}

func (p *Processor) ConfirmAuthDataProcedure(c *gin.Context,
	authEvent models.AuthEvent,
	supi string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	var createAuthStatusRequest Nudr_DataRepository.CreateAuthenticationStatusRequest
	createAuthStatusRequest.AuthEvent = &authEvent
	createAuthStatusRequest.UeId = &supi

	client, err := p.Consumer().CreateUDMClientToUDR(supi)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	_, err = client.AuthenticationStatusDocumentApi.CreateAuthenticationStatus(
		ctx, &createAuthStatusRequest)
	if err != nil {
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			return
		}
		logger.UeauLog.Errorln("ConfirmAuth err:", err.Error())
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	// AuthEvent in response body is optional
	c.JSON(http.StatusCreated, gin.H{})
}

func (p *Processor) GenerateAuthDataProcedure(
	c *gin.Context,
	authInfoRequest models.AuthenticationInfoRequest,
	supiOrSuci string,
) {
	ctx, pd, err := p.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NrfNfManagementNfType_UDR)
	if err != nil {
		c.JSON(int(pd.Status), pd)
		return
	}
	logger.UeauLog.Traceln("In GenerateAuthDataProcedure")

	response := &models.UdmUeauAuthenticationInfoResult{}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	supi, err := suci.ToSupi(supiOrSuci, p.Context().SuciProfiles)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
			Detail: err.Error(),
		}

		logger.UeauLog.Errorln("suciToSupi error: ", err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	logger.UeauLog.Tracef("supi conversion => [%s]", supi)

	client, err := p.Consumer().CreateUDMClientToUDR(supi)
	if err != nil {
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}
	var queryAuthSubsDataRequest Nudr_DataRepository.QueryAuthSubsDataRequest
	queryAuthSubsDataRequest.UeId = &supi

	authSubs, err := client.AuthenticationDataDocumentApi.QueryAuthSubsData(ctx, &queryAuthSubsDataRequest)
	if err != nil {
		logger.ProcLog.Errorf("Error on QueryAuthSubsData: %+v", err)
		apiError, ok := err.(openapi.GenericOpenAPIError)
		if ok {
			c.JSON(apiError.ErrorStatus, apiError.RawBody)
			switch apiError.ErrorStatus {
			case http.StatusNotFound:
				logger.UeauLog.Warnf("Return from UDR QueryAuthSubsData error")
			default:
				logger.UeauLog.Errorln("Return from UDR QueryAuthSubsData error")
			}
			return
		}
		problemDetails := openapi.ProblemDetailsSystemFailure(err.Error())
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	/*
		K, RAND, CK, IK: 128 bits (16 bytes) (hex len = 32)
		SQN, AK: 48 bits (6 bytes) (hex len = 12) TS33.102 - 6.3.2
		AMF: 16 bits (2 bytes) (hex len = 4) TS33.102 - Annex H
	*/

	hasOPC := false
	var kStr, opcStr string
	var k, op, opc []byte
	if authSubs.AuthenticationSubscription.EncPermanentKey != "" {
		kStr = authSubs.AuthenticationSubscription.EncPermanentKey
		if len(kStr) == keyStrLen {
			k, err = hex.DecodeString(kStr)
			if err != nil {
				logger.UeauLog.Errorln("err:", err)
			}
		} else {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  authenticationRejected,
				Detail: "len(kStr) != keyStrLen",
			}

			logger.UeauLog.Errorln("kStr length is ", len(kStr))
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}
	} else {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
			Detail: "EncPermanentKey == ''",
		}

		logger.UeauLog.Errorln("Nil PermanentKey")
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	if authSubs.AuthenticationSubscription.EncOpcKey != "" {
		opcStr = authSubs.AuthenticationSubscription.EncOpcKey
		if len(opcStr) == opcStrLen {
			opc, err = hex.DecodeString(opcStr)
			if err != nil {
				logger.UeauLog.Errorln("err:", err)
			} else {
				hasOPC = true
			}
		} else {
			logger.UeauLog.Errorln("opcStr length is ", len(opcStr))
		}
	} else {
		logger.UeauLog.Infoln("Nil Opc")
	}

	if !hasOPC {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	sqnStr := p.strictHex(authSubs.AuthenticationSubscription.SequenceNumber.Sqn, 12)
	logger.UeauLog.Traceln("sqnStr", sqnStr)
	sqn, err := hex.DecodeString(sqnStr)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
			Detail: err.Error(),
		}

		logger.UeauLog.Errorln("err:", err)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	logger.UeauLog.Tracef("K=[%x], sqn=[%x], OP=[%x], OPC=[%x]", k, sqn, op, opc)

	RAND := make([]byte, 16)
	_, err = cryptoRand.Read(RAND)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
			Detail: err.Error(),
		}

		logger.UeauLog.Errorln("err:", err)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	amfStr := p.strictHex(authSubs.AuthenticationSubscription.AuthenticationManagementField, 4)
	logger.UeauLog.Traceln("amfStr", amfStr)
	AMF, err := hex.DecodeString(amfStr)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
			Detail: err.Error(),
		}

		logger.UeauLog.Errorln("err:", err)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	logger.UeauLog.Tracef("RAND=[%x], AMF=[%x]", RAND, AMF)

	// re-synchronization
	if authInfoRequest.ResynchronizationInfo != nil {
		logger.UeauLog.Infof("Authentication re-synchronization")

		Auts, deCodeErr := hex.DecodeString(authInfoRequest.ResynchronizationInfo.Auts)
		if deCodeErr != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  authenticationRejected,
				Detail: deCodeErr.Error(),
			}

			logger.UeauLog.Errorln("err:", deCodeErr)
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}

		randHex, deCodeErr := hex.DecodeString(authInfoRequest.ResynchronizationInfo.Rand)
		if deCodeErr != nil {
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  authenticationRejected,
				Detail: deCodeErr.Error(),
			}

			logger.UeauLog.Errorln("err:", deCodeErr)
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}

		SQNms, macS := p.aucSQN(opc, k, Auts, randHex)
		if reflect.DeepEqual(macS, Auts[6:]) {
			_, err = cryptoRand.Read(RAND)
			if err != nil {
				problemDetails := &models.ProblemDetails{
					Status: http.StatusForbidden,
					Cause:  authenticationRejected,
					Detail: err.Error(),
				}

				logger.UeauLog.Errorln("err:", err)
				c.JSON(int(problemDetails.Status), problemDetails)
				return
			}

			// increment sqn authSubs.SequenceNumber
			bigSQN := big.NewInt(0)
			sqnStr = hex.EncodeToString(SQNms)
			logger.UeauLog.Tracef("SQNstr=[%s]", sqnStr)
			bigSQN.SetString(sqnStr, 16)

			bigInc := big.NewInt(ind + 1)

			bigP := big.NewInt(SqnMAx)
			bigSQN = bigInc.Add(bigSQN, bigInc)
			bigSQN = bigSQN.Mod(bigSQN, bigP)
			sqnStr = fmt.Sprintf("%x", bigSQN)
			sqnStr = p.strictHex(sqnStr, 12)
		} else {
			logger.UeauLog.Errorf("Re-Sync MAC failed for UE with identity supiOrSuci=[%s], resolvedSupi=[%s]", supiOrSuci, supi)
			logger.UeauLog.Errorln("MACS ", macS)
			logger.UeauLog.Errorln("Auts[6:] ", Auts[6:])
			logger.UeauLog.Errorln("Sqn ", SQNms)
			problemDetails := &models.ProblemDetails{
				Status: http.StatusForbidden,
				Cause:  "modification is rejected",
			}
			c.JSON(int(problemDetails.Status), problemDetails)
			return
		}
	}

	// increment sqn
	bigSQN := big.NewInt(0)
	sqn, err = hex.DecodeString(sqnStr)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  authenticationRejected,
			Detail: err.Error(),
		}

		logger.UeauLog.Errorln("err:", err)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	bigSQN.SetString(sqnStr, 16)

	bigInc := big.NewInt(1)
	bigSQN = bigInc.Add(bigSQN, bigInc)

	SQNheStr := fmt.Sprintf("%x", bigSQN)
	SQNheStr = p.strictHex(SQNheStr, 12)
	patchItemArray := []models.PatchItem{
		{
			Op:   models.PatchOperation_REPLACE,
			Path: "/sequenceNumber",
			Value: models.SequenceNumber{
				Sqn: SQNheStr,
			},
		},
	}

	logger.ProcLog.Infoln("ModifyAuthenticationSubscriptionRequest: ", patchItemArray)

	var modifyAuthenticationSubscriptionRequest Nudr_DataRepository.ModifyAuthenticationSubscriptionRequest
	modifyAuthenticationSubscriptionRequest.UeId = &supi
	modifyAuthenticationSubscriptionRequest.PatchItem = patchItemArray
	_, err = client.AuthenticationSubscriptionDocumentApi.ModifyAuthenticationSubscription(
		ctx, &modifyAuthenticationSubscriptionRequest)
	if err != nil {
		problemDetails := &models.ProblemDetails{
			Status: http.StatusForbidden,
			Cause:  "modification is rejected ",
			Detail: err.Error(),
		}

		logger.UeauLog.Errorln("update sqn error:", err)
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	// Run milenage
	macA, macS := make([]byte, 8), make([]byte, 8)
	CK, IK := make([]byte, 16), make([]byte, 16)
	RES := make([]byte, 8)
	AK, AKstar := make([]byte, 6), make([]byte, 6)

	// Generate macA, macS
	err = util.MilenageF1(opc, k, RAND, sqn, AMF, macA, macS)
	if err != nil {
		logger.UeauLog.Errorln("milenage F1 err:", err)
	}

	// Generate RES, CK, IK, AK, AKstar
	// RES == XRES (expected RES) for server
	err = util.MilenageF2345(opc, k, RAND, RES, CK, IK, AK, AKstar)
	if err != nil {
		logger.UeauLog.Errorln("milenage F2345 err:", err)
	}
	logger.UeauLog.Tracef("milenage RES=[%s]", hex.EncodeToString(RES))

	// Generate AUTN
	logger.UeauLog.Tracef("SQN=[%x], AK=[%x]", sqn, AK)
	logger.UeauLog.Tracef("AMF=[%x], macA=[%x]", AMF, macA)
	SQNxorAK := make([]byte, 6)
	for i := 0; i < len(sqn); i++ {
		SQNxorAK[i] = sqn[i] ^ AK[i]
	}
	logger.UeauLog.Tracef("SQN xor AK=[%x]", SQNxorAK)
	AUTN := append(append(SQNxorAK, AMF...), macA...)
	logger.UeauLog.Tracef("AUTN=[%x]", AUTN)

	var av models.AuthenticationVector
	if authSubs.AuthenticationSubscription.AuthenticationMethod == models.AuthMethod__5_G_AKA {
		response.AuthType = models.UdmUeauAuthType__5_G_AKA

		// derive XRES*
		key := append(CK, IK...)
		FC := ueauth.FC_FOR_RES_STAR_XRES_STAR_DERIVATION
		P0 := []byte(authInfoRequest.ServingNetworkName)
		P1 := RAND
		P2 := RES

		kdfValForXresStar, err := ueauth.GetKDFValue(
			key, FC, P0, ueauth.KDFLen(P0), P1, ueauth.KDFLen(P1), P2, ueauth.KDFLen(P2))
		if err != nil {
			logger.UeauLog.Errorf("Get kdfValForXresStar err: %+v", err)
		}
		xresStar := kdfValForXresStar[len(kdfValForXresStar)/2:]
		logger.UeauLog.Tracef("xresStar=[%x]", xresStar)

		// derive Kausf
		FC = ueauth.FC_FOR_KAUSF_DERIVATION
		P0 = []byte(authInfoRequest.ServingNetworkName)
		P1 = SQNxorAK
		kdfValForKausf, err := ueauth.GetKDFValue(key, FC, P0, ueauth.KDFLen(P0), P1, ueauth.KDFLen(P1))
		if err != nil {
			logger.UeauLog.Errorf("Get kdfValForKausf err: %+v", err)
		}
		logger.UeauLog.Tracef("Kausf=[%x]", kdfValForKausf)

		// Fill in rand, xresStar, autn, kausf
		av.Rand = hex.EncodeToString(RAND)
		av.XresStar = hex.EncodeToString(xresStar)
		av.Autn = hex.EncodeToString(AUTN)
		av.Kausf = hex.EncodeToString(kdfValForKausf)
		av.AvType = models.AvType__5_G_HE_AKA
	} else { // EAP-AKA'
		response.AuthType = models.UdmUeauAuthType_EAP_AKA_PRIME
		// derive CK' and IK'
		key := append(CK, IK...)
		FC := ueauth.FC_FOR_CK_PRIME_IK_PRIME_DERIVATION
		P0 := []byte(authInfoRequest.ServingNetworkName)
		P1 := SQNxorAK
		kdfVal, err := ueauth.GetKDFValue(key, FC, P0, ueauth.KDFLen(P0), P1, ueauth.KDFLen(P1))
		if err != nil {
			logger.UeauLog.Errorf("Get kdfVal err: %+v", err)
		}
		logger.UeauLog.Tracef("kdfVal=[%x] (len=%d)", kdfVal, len(kdfVal))

		// For TS 35.208 test set 19 & RFC 5448 test vector 1
		// CK': 0093 962d 0dd8 4aa5 684b 045c 9edf fa04
		// IK': ccfc 230c a74f cc96 c0a5 d611 64f5 a76

		ckPrime := kdfVal[:len(kdfVal)/2]
		ikPrime := kdfVal[len(kdfVal)/2:]
		logger.UeauLog.Tracef("ckPrime=[%x], kPrime=[%x]", ckPrime, ikPrime)

		// Fill in rand, xres, autn, ckPrime, ikPrime
		av.Rand = hex.EncodeToString(RAND)
		av.Xres = hex.EncodeToString(RES)
		av.Autn = hex.EncodeToString(AUTN)
		av.CkPrime = hex.EncodeToString(ckPrime)
		av.IkPrime = hex.EncodeToString(ikPrime)
		av.AvType = models.AvType_EAP_AKA_PRIME
	}

	response.AuthenticationVector = &av
	response.Supi = supi
	c.JSON(http.StatusOK, response)
}
