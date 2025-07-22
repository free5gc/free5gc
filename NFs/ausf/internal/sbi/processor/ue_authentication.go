package processor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bronze1man/radius"
	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	ausf_context "github.com/free5gc/ausf/internal/context"
	"github.com/free5gc/ausf/internal/logger"
	"github.com/free5gc/ausf/pkg/factory"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/ueauth"
)

func (p *Processor) HandleEapAuthComfirmRequest(c *gin.Context, eapSession models.EapSession, eapSessionId string) {
	logger.Auth5gAkaLog.Infof("EapAuthComfirmRequest")

	p.EapAuthComfirmRequestProcedure(c, eapSession, eapSessionId)
}

func (p *Processor) EapAuthComfirmRequestProcedure(
	c *gin.Context,
	updateEapSession models.EapSession,
	eapSessionID string,
) {
	var eapSession models.EapSession

	if !ausf_context.CheckIfSuciSupiPairExists(eapSessionID) {
		logger.AuthELog.Infoln("supiSuciPair does not exist, confirmation failed")
		problemDetails := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "USER_NOT_FOUND",
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	currentSupi := ausf_context.GetSupiFromSuciSupiMap(eapSessionID)
	if !ausf_context.CheckIfAusfUeContextExists(currentSupi) {
		logger.AuthELog.Infoln("SUPI does not exist, confirmation failed")
		problemDetails := models.ProblemDetails{
			Status: http.StatusNotFound,
			Cause:  "USER_NOT_FOUND",
		}
		c.JSON(int(problemDetails.Status), problemDetails)
		return
	}

	ausfCurrentContext := ausf_context.GetAusfUeContext(currentSupi)
	servingNetworkName := ausfCurrentContext.ServingNetworkName

	if ausfCurrentContext.AuthStatus == models.AusfUeAuthenticationAuthResult_FAILURE {
		logger.AuthELog.Warnf("Authentication failed with status: %s", ausfCurrentContext.AuthStatus)
		eapFailPkt := ConstructEapNoTypePkt(radius.EapCodeFailure, 0)
		eapSession.EapPayload = eapFailPkt
		eapSession.AuthResult = models.AusfUeAuthenticationAuthResult_FAILURE
		c.JSON(http.StatusUnauthorized, eapSession)
		return
	}

	var eapPayload []byte
	if eapPayloadTmp, err := base64.StdEncoding.DecodeString(updateEapSession.EapPayload); err != nil {
		logger.AuthELog.Warnf("EAP Payload decode failed: %+v", err)
	} else {
		eapPayload = eapPayloadTmp
	}

	eapGoPkt := gopacket.NewPacket(eapPayload, layers.LayerTypeEAP, gopacket.Default)
	eapLayer := eapGoPkt.Layer(layers.LayerTypeEAP)
	eapContent, _ := eapLayer.(*layers.EAP)
	eapOK := true
	var eapErrStr string

	if eapContent.Code != layers.EAPCodeResponse {
		eapOK = false
		eapErrStr = "eap packet code error"
	} else if eapContent.Type != ausf_context.EAP_AKA_PRIME_TYPENUM {
		eapOK = false
		eapErrStr = "eap packet type error"
	} else if decodeEapAkaPrimePkt, err := decodeEapAkaPrime(eapContent.Contents); err != nil {
		logger.AuthELog.Warnf("EAP-AKA' decode failed: %+v", err)
		eapOK = false
		eapErrStr = "eap packet error"
	} else {
		switch decodeEapAkaPrimePkt.Subtype {
		case ausf_context.AKA_CHALLENGE_SUBTYPE:
			K_autStr := ausfCurrentContext.K_aut
			var K_aut []byte
			if K_autTmp, err := hex.DecodeString(K_autStr); err != nil {
				logger.AuthELog.Warnf("K_aut decode error: %+v", err)
			} else {
				K_aut = K_autTmp
			}
			XMAC := CalculateAtMAC(K_aut, decodeEapAkaPrimePkt.MACInput)
			MAC := decodeEapAkaPrimePkt.Attributes[ausf_context.AT_MAC_ATTRIBUTE].Value
			XRES := ausfCurrentContext.XRES
			RES := hex.EncodeToString(decodeEapAkaPrimePkt.Attributes[ausf_context.AT_RES_ATTRIBUTE].Value)

			if !bytes.Equal(MAC, XMAC) {
				eapOK = false
				eapErrStr = "EAP-AKA' integrity check fail"
			} else if XRES == RES {
				logger.AuthELog.Infoln("Correct RES value, EAP-AKA' auth succeed")
				eapSession.KSeaf = ausfCurrentContext.Kseaf
				eapSession.Supi = currentSupi
				eapSession.AuthResult = models.AusfUeAuthenticationAuthResult_SUCCESS
				eapSuccPkt := ConstructEapNoTypePkt(radius.EapCodeSuccess, eapContent.Id)
				eapSession.EapPayload = eapSuccPkt
				udmUrl := ausfCurrentContext.UdmUeauUrl
				if sendErr := p.Consumer().SendAuthResultToUDM(
					eapSessionID,
					models.UdmUeauAuthType_EAP_AKA_PRIME,
					true,
					servingNetworkName,
					udmUrl); sendErr != nil {
					logger.AuthELog.Infoln(sendErr.Error())
					problemDetails := models.ProblemDetails{
						Cause: "UPSTREAM_SERVER_ERROR",
					}
					c.JSON(http.StatusInternalServerError, problemDetails)
					return
				}
				ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_SUCCESS
			} else {
				eapOK = false
				eapErrStr = "Wrong RES value, EAP-AKA' auth failed"
			}
		case ausf_context.AKA_AUTHENTICATION_REJECT_SUBTYPE:
			ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_FAILURE
		case ausf_context.AKA_SYNCHRONIZATION_FAILURE_SUBTYPE:
			logger.AuthELog.Warnf("EAP-AKA' synchronziation failure")
			if ausfCurrentContext.Resynced {
				eapOK = false
				eapErrStr = "2 consecutive Synch Failure, terminate authentication procedure"
			} else {
				var authInfo models.AuthenticationInfo
				AUTS := decodeEapAkaPrimePkt.Attributes[ausf_context.AT_AUTS_ATTRIBUTE].Value
				resynchronizationInfo := &models.ResynchronizationInfo{
					Auts: hex.EncodeToString(AUTS[:]),
				}
				authInfo.SupiOrSuci = eapSessionID
				authInfo.ServingNetworkName = servingNetworkName
				authInfo.ResynchronizationInfo = resynchronizationInfo
				p.UeAuthPostRequestProcedure(c, authInfo)
				return
			}
		case ausf_context.AKA_NOTIFICATION_SUBTYPE:
			ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_FAILURE
		case ausf_context.AKA_CLIENT_ERROR_SUBTYPE:
			logger.AuthELog.Warnf("EAP-AKA' failure: receive client-error")
			ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_FAILURE
		default:
			ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_FAILURE
		}
	}

	if !eapOK {
		logger.AuthELog.Warnf("EAP-AKA' failure: %s", eapErrStr)
		if sendErr := p.Consumer().SendAuthResultToUDM(eapSessionID, models.UdmUeauAuthType_EAP_AKA_PRIME,
			false, servingNetworkName, ausfCurrentContext.UdmUeauUrl); sendErr != nil {
			logger.AuthELog.Infoln(sendErr.Error())
			problemDetails := models.ProblemDetails{
				Status: http.StatusInternalServerError,
				Cause:  "UPSTREAM_SERVER_ERROR",
			}
			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		}

		ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_FAILURE
		eapSession.AuthResult = models.AusfUeAuthenticationAuthResult_ONGOING
		failEapAkaNoti := ConstructFailEapAkaNotification(eapContent.Id)
		eapSession.EapPayload = failEapAkaNoti
		self := ausf_context.GetSelf()
		linkUrl := self.Url + factory.AusfAuthResUriPrefix + "/ue-authentications/" + eapSessionID + "/eap-session"
		linksValue := models.Link{Href: linkUrl}
		eapSession.Links = make(map[string][]models.Link)
		eapSession.Links["eap-session"] = []models.Link{linksValue}
	} else if ausfCurrentContext.AuthStatus == models.AusfUeAuthenticationAuthResult_FAILURE {
		if sendErr := p.Consumer().SendAuthResultToUDM(eapSessionID, models.UdmUeauAuthType_EAP_AKA_PRIME, false,
			servingNetworkName, ausfCurrentContext.UdmUeauUrl); sendErr != nil {
			logger.AuthELog.Infoln(sendErr.Error())
			var problemDetails models.ProblemDetails
			problemDetails.Status = http.StatusInternalServerError
			problemDetails.Cause = "UPSTREAM_SERVER_ERROR"
			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		}

		eapFailPkt := ConstructEapNoTypePkt(radius.EapCodeFailure, eapPayload[1])
		eapSession.EapPayload = eapFailPkt
		eapSession.AuthResult = models.AusfUeAuthenticationAuthResult_FAILURE
	}

	c.JSON(http.StatusOK, eapSession)
}

func (p *Processor) HandleUeAuthPostRequest(c *gin.Context, authenticationInfo models.AuthenticationInfo) {
	logger.UeAuthLog.Infof("HandleUeAuthPostRequest")
	p.UeAuthPostRequestProcedure(c, authenticationInfo)
}

func (p *Processor) UeAuthPostRequestProcedure(c *gin.Context, updateAuthenticationInfo models.AuthenticationInfo) {
	var responseBody models.UeAuthenticationCtx
	var authInfoReq models.AuthenticationInfoRequest

	supiOrSuci := updateAuthenticationInfo.SupiOrSuci

	snName := updateAuthenticationInfo.ServingNetworkName
	servingNetworkAuthorized := ausf_context.IsServingNetworkAuthorized(snName)
	if !servingNetworkAuthorized {
		problemDetails := models.ProblemDetails{
			Cause:  "SERVING_NETWORK_NOT_AUTHORIZED",
			Status: http.StatusForbidden,
		}
		logger.UeAuthLog.Infoln("403 forbidden: serving network NOT AUTHORIZED")
		c.JSON(http.StatusForbidden, problemDetails)
		return
	}
	logger.UeAuthLog.Infoln("Serving network authorized")

	responseBody.ServingNetworkName = snName
	authInfoReq.ServingNetworkName = snName
	self := ausf_context.GetSelf()
	authInfoReq.AusfInstanceId = self.GetSelfID()

	var lastEapID uint8
	if updateAuthenticationInfo.ResynchronizationInfo != nil {
		logger.UeAuthLog.Warningln("Auts: ", updateAuthenticationInfo.ResynchronizationInfo.Auts)
		ausfCurrentSupi := ausf_context.GetSupiFromSuciSupiMap(supiOrSuci)
		logger.UeAuthLog.Warningln(ausfCurrentSupi)
		ausfCurrentContext := ausf_context.GetAusfUeContext(ausfCurrentSupi)
		logger.UeAuthLog.Warningln(ausfCurrentContext.Rand)
		if updateAuthenticationInfo.ResynchronizationInfo.Rand == "" {
			updateAuthenticationInfo.ResynchronizationInfo.Rand = ausfCurrentContext.Rand
		}
		logger.UeAuthLog.Warningln("Rand: ", updateAuthenticationInfo.ResynchronizationInfo.Rand)
		authInfoReq.ResynchronizationInfo = updateAuthenticationInfo.ResynchronizationInfo
		lastEapID = ausfCurrentContext.EapID
	}

	udmUrl := p.Consumer().GetUdmUrl(self.NrfUri)

	result, pd, err := p.Consumer().GenerateAuthDataApi(udmUrl, supiOrSuci, authInfoReq)
	if err != nil {
		logger.UeAuthLog.Infof("GenerateAuthDataApi error: %+v", err)
		c.JSON(http.StatusInternalServerError, pd)
		return
	}
	authInfoResult := *result

	ueid := authInfoResult.Supi
	ausfUeContext := ausf_context.NewAusfUeContext(ueid)
	ausfUeContext.ServingNetworkName = snName
	ausfUeContext.AuthStatus = models.AusfUeAuthenticationAuthResult_ONGOING
	ausfUeContext.UdmUeauUrl = udmUrl
	ausf_context.AddAusfUeContextToPool(ausfUeContext)

	logger.UeAuthLog.Infof("Add SuciSupiPair (%s, %s) to map.\n", supiOrSuci, ueid)
	ausf_context.AddSuciSupiPairToMap(supiOrSuci, ueid)

	locationURI := self.Url + factory.AusfAuthResUriPrefix + "/ue-authentications/" + supiOrSuci
	putLink := locationURI
	switch authInfoResult.AuthType {
	case models.UdmUeauAuthType__5_G_AKA:
		logger.UeAuthLog.Infoln("Use 5G AKA auth method")
		putLink += "/5g-aka-confirmation"

		// Derive HXRES* from XRES*
		concat := authInfoResult.AuthenticationVector.Rand + authInfoResult.AuthenticationVector.XresStar
		var hxresStarBytes []byte
		if bytes, err := hex.DecodeString(concat); err != nil {
			logger.Auth5gAkaLog.Errorf("decode concat error: %+v", err)
			problemDetails := models.ProblemDetails{
				Title:  "Concat Decode Problem",
				Cause:  "CONCAT_DECODE_PROBLEM",
				Detail: err.Error(),
				Status: http.StatusInternalServerError,
			}
			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		} else {
			hxresStarBytes = bytes
		}
		hxresStarAll := sha256.Sum256(hxresStarBytes)
		hxresStar := hex.EncodeToString(hxresStarAll[16:]) // last 128 bits
		logger.Auth5gAkaLog.Infof("XresStar = %x\n", authInfoResult.AuthenticationVector.XresStar)

		// Derive Kseaf from Kausf
		Kausf := authInfoResult.AuthenticationVector.Kausf
		var KausfDecode []byte
		if ausfDecode, err := hex.DecodeString(Kausf); err != nil {
			logger.Auth5gAkaLog.Errorf("decode Kausf failed: %+v", err)
			problemDetails := models.ProblemDetails{
				Title:  "Kausf Decode Problem",
				Cause:  "KAUSF_DECODE_PROBLEM",
				Detail: err.Error(),
				Status: http.StatusInternalServerError,
			}
			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		} else {
			KausfDecode = ausfDecode
		}
		P0 := []byte(snName)
		Kseaf, err := ueauth.GetKDFValue(KausfDecode, ueauth.FC_FOR_KSEAF_DERIVATION, P0, ueauth.KDFLen(P0))
		if err != nil {
			logger.Auth5gAkaLog.Errorf("GetKDFValue failed: %+v", err)
			problemDetails := models.ProblemDetails{
				Title:  "Kseaf Derivation Problem",
				Cause:  "KSEAF_DERIVATION_PROBLEM",
				Detail: err.Error(),
				Status: http.StatusInternalServerError,
			}
			c.JSON(http.StatusInternalServerError, problemDetails)
			return
		}
		ausfUeContext.XresStar = authInfoResult.AuthenticationVector.XresStar
		ausfUeContext.Kausf = Kausf
		ausfUeContext.Kseaf = hex.EncodeToString(Kseaf)
		ausfUeContext.Rand = authInfoResult.AuthenticationVector.Rand

		var av5gAka models.Av5gAka
		av5gAka.Rand = authInfoResult.AuthenticationVector.Rand
		av5gAka.Autn = authInfoResult.AuthenticationVector.Autn
		av5gAka.HxresStar = hxresStar
		responseBody.Var5gAuthData = av5gAka

		linksValue := models.Link{Href: putLink}
		responseBody.Links = make(map[string][]models.Link)
		responseBody.Links["5g-aka"] = []models.Link{linksValue}
	case models.UdmUeauAuthType_EAP_AKA_PRIME:
		logger.UeAuthLog.Infoln("Use EAP-AKA' auth method")
		putLink += "/eap-session"

		var identity string
		// TODO support more SUPI type
		if ueid[:4] == "imsi" {
			if !self.EapAkaSupiImsiPrefix {
				// 33.501 v15.9.0 or later
				identity = ueid[5:]
			} else {
				// 33.501 v15.8.0 or earlier
				identity = ueid
			}
		}
		ikPrime := authInfoResult.AuthenticationVector.IkPrime
		ckPrime := authInfoResult.AuthenticationVector.CkPrime
		RAND := authInfoResult.AuthenticationVector.Rand
		AUTN := authInfoResult.AuthenticationVector.Autn
		XRES := authInfoResult.AuthenticationVector.Xres
		ausfUeContext.XRES = XRES

		ausfUeContext.Rand = authInfoResult.AuthenticationVector.Rand

		_, K_aut, _, _, EMSK := eapAkaPrimePrf(ikPrime, ckPrime, identity)
		logger.AuthELog.Tracef("K_aut: %x", K_aut)
		ausfUeContext.K_aut = hex.EncodeToString(K_aut)
		Kausf := EMSK[0:32]
		ausfUeContext.Kausf = hex.EncodeToString(Kausf)
		P0 := []byte(snName)
		Kseaf, err := ueauth.GetKDFValue(Kausf, ueauth.FC_FOR_KSEAF_DERIVATION, P0, ueauth.KDFLen(P0))
		if err != nil {
			logger.AuthELog.Errorf("GetKDFValue failed: %+v", err)
		}
		ausfUeContext.Kseaf = hex.EncodeToString(Kseaf)

		var eapPkt radius.EapPacket
		eapPkt.Code = radius.EapCode(1)
		if updateAuthenticationInfo.ResynchronizationInfo == nil {
			src := rand.NewSource(time.Now().UnixNano())
			r := rand.New(src)
			randIdentifier := r.Intn(256)
			ausfUeContext.EapID = uint8(randIdentifier)
		} else {
			ausfUeContext.EapID = lastEapID + 1
		}
		eapPkt.Identifier = ausfUeContext.EapID
		eapPkt.Type = radius.EapType(50) // according to RFC5448 6.1

		var eapAKAHdr, atRand, atAutn, atKdf, atKdfInput, atMAC string
		eapAKAHdrBytes := make([]byte, 3) // RFC4187 8.1
		eapAKAHdrBytes[0] = ausf_context.AKA_CHALLENGE_SUBTYPE
		eapAKAHdr = string(eapAKAHdrBytes)
		if atRandTmp, err := EapEncodeAttribute("AT_RAND", RAND); err != nil {
			logger.AuthELog.Errorf("EAP encode RAND failed: %+v", err)
		} else {
			atRand = atRandTmp
		}
		if atAutnTmp, err := EapEncodeAttribute("AT_AUTN", AUTN); err != nil {
			logger.AuthELog.Errorf("EAP encode AUTN failed: %+v", err)
		} else {
			atAutn = atAutnTmp
		}
		if atKdfTmp, err := EapEncodeAttribute("AT_KDF", snName); err != nil {
			logger.AuthELog.Errorf("EAP encode KDF failed: %+v", err)
		} else {
			atKdf = atKdfTmp
		}
		if atKdfInputTmp, err := EapEncodeAttribute("AT_KDF_INPUT", snName); err != nil {
			logger.AuthELog.Errorf("EAP encode KDF failed: %+v", err)
		} else {
			atKdfInput = atKdfInputTmp
		}
		if atMACTmp, err := EapEncodeAttribute("AT_MAC", ""); err != nil {
			logger.AuthELog.Errorf("EAP encode MAC failed: %+v", err)
		} else {
			atMAC = atMACTmp
		}

		dataArrayBeforeMAC := eapAKAHdr + atRand + atAutn + atKdf + atKdfInput + atMAC
		eapPkt.Data = []byte(dataArrayBeforeMAC)
		encodedPktBeforeMAC := eapPkt.Encode()

		MacValue := CalculateAtMAC(K_aut, encodedPktBeforeMAC)
		atMAC = atMAC[:4] + string(MacValue)

		dataArrayAfterMAC := eapAKAHdr + atRand + atAutn + atKdf + atKdfInput + atMAC

		eapPkt.Data = []byte(dataArrayAfterMAC)
		encodedPktAfterMAC := eapPkt.Encode()
		responseBody.Var5gAuthData = base64.StdEncoding.EncodeToString(encodedPktAfterMAC)

		linksValue := models.Link{Href: putLink}
		responseBody.Links = make(map[string][]models.Link)
		responseBody.Links["eap-session"] = []models.Link{linksValue}
	}

	responseBody.AuthType = models.AusfUeAuthenticationAuthType(authInfoResult.AuthType)

	c.Header("Location", locationURI)
	c.JSON(http.StatusCreated, responseBody)
}

func (p *Processor) HandleAuth5gAkaComfirmRequest(
	c *gin.Context,
	confirmationData models.ConfirmationData,
	confirmationDataResponseId string,
) {
	logger.Auth5gAkaLog.Infof("Auth5gAkaComfirmRequest")
	p.Auth5gAkaComfirmRequestProcedure(c, confirmationData, confirmationDataResponseId)
}

func (p *Processor) Auth5gAkaComfirmRequestProcedure(c *gin.Context, updateConfirmationData models.ConfirmationData,
	ConfirmationDataResponseID string,
) {
	var confirmDataRsp models.ConfirmationDataResponse
	success := false
	confirmDataRsp.AuthResult = models.AusfUeAuthenticationAuthResult_FAILURE

	if !ausf_context.CheckIfSuciSupiPairExists(ConfirmationDataResponseID) {
		logger.Auth5gAkaLog.Infof("supiSuciPair does not exist, confirmation failed (queried by %s)\n",
			ConfirmationDataResponseID)
		problemDetails := models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	currentSupi := ausf_context.GetSupiFromSuciSupiMap(ConfirmationDataResponseID)
	if !ausf_context.CheckIfAusfUeContextExists(currentSupi) {
		logger.Auth5gAkaLog.Infof("SUPI does not exist, confirmation failed (queried by %s)\n", currentSupi)
		problemDetails := models.ProblemDetails{
			Cause:  "USER_NOT_FOUND",
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, problemDetails)
		return
	}

	ausfCurrentContext := ausf_context.GetAusfUeContext(currentSupi)
	servingNetworkName := ausfCurrentContext.ServingNetworkName

	// Compare the received RES* with the stored XRES*
	logger.Auth5gAkaLog.Infof("res*: %x\nXres*: %x\n", updateConfirmationData.ResStar, ausfCurrentContext.XresStar)
	if strings.EqualFold(updateConfirmationData.ResStar, ausfCurrentContext.XresStar) {
		ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_SUCCESS
		confirmDataRsp.AuthResult = models.AusfUeAuthenticationAuthResult_SUCCESS
		success = true
		logger.Auth5gAkaLog.Infoln("5G AKA confirmation succeeded")
		confirmDataRsp.Supi = currentSupi
		confirmDataRsp.Kseaf = ausfCurrentContext.Kseaf
	} else {
		ausfCurrentContext.AuthStatus = models.AusfUeAuthenticationAuthResult_FAILURE
		confirmDataRsp.AuthResult = models.AusfUeAuthenticationAuthResult_FAILURE
		p.logConfirmFailureAndInformUDM(ConfirmationDataResponseID, models.AusfUeAuthenticationAuthType__5_G_AKA,
			servingNetworkName, "5G AKA confirmation failed", ausfCurrentContext.UdmUeauUrl)
	}

	if sendErr := p.Consumer().SendAuthResultToUDM(currentSupi, models.UdmUeauAuthType__5_G_AKA, success,
		servingNetworkName, ausfCurrentContext.UdmUeauUrl); sendErr != nil {
		logger.Auth5gAkaLog.Infoln(sendErr.Error())
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "UPSTREAM_SERVER_ERROR",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
		return
	}

	c.JSON(http.StatusOK, confirmDataRsp)
}

func KDF5gAka(param ...string) hash.Hash {
	s := param[0]
	s += param[1]
	if p0len, err := strconv.Atoi(param[2]); err != nil {
		logger.AuthELog.Warnf("atoi failed: %+v", err)
	} else {
		s += strconv.FormatInt(int64(p0len), 16)
	}
	h := hmac.New(sha256.New, []byte(s))

	return h
}

func intToByteArray(i int) []byte {
	r := make([]byte, 2)
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
		logger.AuthELog.Errorln(err.Error())
	}
	sum := h.Sum(nil)
	return sum[:16]
}

func EapEncodeAttribute(attributeType string, data string) (string, error) {
	var attribute string
	var length int

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
		logger.AuthELog.Errorf("UNKNOWN attributeType %s\n", attributeType)
		return "", nil
	}

	if r, err := hex.DecodeString(attribute); err != nil {
		return "", err
	} else {
		return string(r), nil
	}
}

func eapAkaPrimePrf(ikPrime string, ckPrime string, identity string) ([]byte, []byte, []byte, []byte, []byte) {
	keyAp := ikPrime + ckPrime

	var key []byte
	if keyTmp, err := hex.DecodeString(keyAp); err != nil {
		logger.AuthELog.Warnf("Decode key AP failed: %+v", err)
	} else {
		key = keyTmp
	}
	sBase := []byte("EAP-AKA'" + identity)

	MK := []byte("")
	prev := []byte("")
	prfRounds := 208/32 + 1
	for i := 0; i < prfRounds; i++ {
		// Create a new HMAC by defining the hash type and the key (as byte array)
		h := hmac.New(sha256.New, key)

		hexNum := (byte)(i + 1)
		ap := append(sBase, hexNum)
		s := append(prev, ap...)

		// Write Data to it
		if _, err := h.Write(s); err != nil {
			logger.AuthELog.Errorln(err.Error())
		}

		// Get result
		sha := h.Sum(nil)
		MK = append(MK, sha...)
		prev = sha
	}

	K_encr := MK[0:16]  // 0..127
	K_aut := MK[16:48]  // 128..383
	K_re := MK[48:80]   // 384..639
	MSK := MK[80:144]   // 640..1151
	EMSK := MK[144:208] // 1152..1663
	return K_encr, K_aut, K_re, MSK, EMSK
}

func decodeEapAkaPrime(eapPkt []byte) (*ausf_context.EapAkaPrimePkt, error) {
	var decodePkt ausf_context.EapAkaPrimePkt
	var attrLen int
	var decodeAttr ausf_context.EapAkaPrimeAttribute
	attributes := make(map[uint8]ausf_context.EapAkaPrimeAttribute)
	data := eapPkt[5:]
	decodePkt.Subtype = data[0]
	dataLen := len(data)

	// decode attributes
	for i := 3; i < dataLen; i += attrLen {
		attrType := data[i]
		attrLen = int(data[i+1]) * 4
		if attrLen == 0 {
			return nil, fmt.Errorf("attribute length equal to zero")
		}
		if i+attrLen > dataLen {
			return nil, fmt.Errorf("packet length out of range")
		}
		switch attrType {
		case ausf_context.AT_RES_ATTRIBUTE:
			logger.AuthELog.Tracef("Decoding AT_RES\n")
			accLen := int(data[i+3] >> 3)
			if accLen > 16 || accLen < 4 || accLen+4 > attrLen {
				return nil, fmt.Errorf("attribute AT_RES decode err")
			}

			decodeAttr.Type = attrType
			decodeAttr.Length = data[i+1]
			decodeAttr.Value = data[i+4 : i+4+accLen]
			attributes[attrType] = decodeAttr
		case ausf_context.AT_MAC_ATTRIBUTE:
			logger.AuthELog.Tracef("Decoding AT_MAC\n")
			if attrLen != 20 {
				return nil, fmt.Errorf("attribute AT_MAC decode err")
			}
			decodeAttr.Type = attrType
			decodeAttr.Length = data[i+1]
			Mac := make([]byte, attrLen-4)
			copy(Mac, data[i+4:i+attrLen])
			decodeAttr.Value = Mac
			attributes[attrType] = decodeAttr

			// clean AT_MAC value for integrity check later
			zeros := make([]byte, attrLen-4)
			copy(data[i+4:i+attrLen], zeros)
			decodePkt.MACInput = eapPkt
		case ausf_context.AT_KDF_ATTRIBUTE:
			logger.AuthELog.Tracef("Decoding AT_KDF\n")
			if attrLen != 4 {
				return nil, fmt.Errorf("attribute AT_KDF decode err")
			}
			decodeAttr.Type = attrType
			decodeAttr.Length = data[i+1]
			decodeAttr.Value = data[i+2 : i+attrLen]
			attributes[attrType] = decodeAttr
		case ausf_context.AT_AUTS_ATTRIBUTE:
			logger.AuthELog.Tracef("Decoding AT_AUTS\n")
			if attrLen != 16 {
				return nil, fmt.Errorf("attribute AT_AUTS decode err")
			}
			decodeAttr.Type = attrType
			decodeAttr.Length = data[i+1]
			decodeAttr.Value = data[i+2 : i+attrLen]
			attributes[attrType] = decodeAttr
		case ausf_context.AT_CLIENT_ERROR_CODE_ATTRIBUTE:
			logger.AuthELog.Tracef("Decoding AT_CLIENT_ERROR_CODE\n")
			if attrLen != 4 {
				return nil, fmt.Errorf("attribute AT_CLIENT_ERROR_CODE decode err")
			}
			decodeAttr.Type = attrType
			decodeAttr.Length = data[i+1]
			decodeAttr.Value = data[i+2 : i+attrLen]
			attributes[attrType] = decodeAttr
		default:
			logger.AuthELog.Tracef("attribute type %x skipped\n", attrType)
		}
	}

	switch decodePkt.Subtype {
	case ausf_context.AKA_CHALLENGE_SUBTYPE:
		logger.AuthELog.Tracef("Subtype AKA-Challenge\n")
		if _, ok := attributes[ausf_context.AT_RES_ATTRIBUTE]; !ok {
			return nil, fmt.Errorf("AKA-Challenge attributes error")
		} else if _, ok := attributes[ausf_context.AT_MAC_ATTRIBUTE]; !ok {
			return nil, fmt.Errorf("AKA-Challenge attributes error")
		}
	case ausf_context.AKA_AUTHENTICATION_REJECT_SUBTYPE:
		logger.AuthELog.Tracef("Subtype AKA-Authentication-Reject\n")
		if len(attributes) != 0 {
			return nil, fmt.Errorf("AKA-Authentication-Reject attributes error")
		}
	case ausf_context.AKA_SYNCHRONIZATION_FAILURE_SUBTYPE:
		logger.AuthELog.Tracef("Subtype AKA-Synchronization-Failure\n")
		if len(attributes) != 2 {
			return nil, fmt.Errorf("AKA-Synchornization-Failure attributes error")
		} else if _, ok := attributes[ausf_context.AT_AUTS_ATTRIBUTE]; !ok {
			return nil, fmt.Errorf("AKA-Synchornization-Failure attributes error")
		} else if _, ok := attributes[ausf_context.AT_KDF_ATTRIBUTE]; !ok {
			return nil, fmt.Errorf("AKA-Synchornization-Failure attributes error")
		} else if kdfVal := attributes[ausf_context.AT_KDF_ATTRIBUTE].Value; kdfVal[0] != 0 || kdfVal[1] != 1 {
			return nil, fmt.Errorf("AKA-Synchornization-Failure attributes error")
		}
	case ausf_context.AKA_NOTIFICATION_SUBTYPE:
		logger.AuthELog.Tracef("Subtype AKA-Notification\n")
	case ausf_context.AKA_CLIENT_ERROR_SUBTYPE:
		logger.AuthELog.Tracef("Subtype AKA-Client-Error\n")
		if len(attributes) != 1 {
			return nil, fmt.Errorf("AKA-Client-Error attributes error")
		} else if _, ok := attributes[ausf_context.AT_CLIENT_ERROR_CODE_ATTRIBUTE]; !ok {
			return nil, fmt.Errorf("AKA-Client-Error attributes error")
		}
	default:
		logger.AuthELog.Tracef("subtype %x skipped\n", decodePkt.Subtype)
	}

	decodePkt.Attributes = attributes

	return &decodePkt, nil
}

func ConstructFailEapAkaNotification(oldPktId uint8) string {
	var eapPkt radius.EapPacket
	eapPkt.Code = radius.EapCodeRequest
	eapPkt.Identifier = oldPktId + 1
	eapPkt.Type = ausf_context.EAP_AKA_PRIME_TYPENUM

	eapAkaHdrBytes := make([]byte, 3)
	eapAkaHdrBytes[0] = ausf_context.AKA_NOTIFICATION_SUBTYPE

	attrNum := fmt.Sprintf("%02x", ausf_context.AT_NOTIFICATION_ATTRIBUTE)
	attribute := attrNum + "01" + "4000"
	var attrHex []byte
	if attrHexTmp, err := hex.DecodeString(attribute); err != nil {
		logger.AuthELog.Warnf("Decode attribute failed: %+v", err)
	} else {
		attrHex = attrHexTmp
	}

	eapPkt.Data = append(eapAkaHdrBytes, attrHex...)
	eapPktEncode := eapPkt.Encode()
	return base64.StdEncoding.EncodeToString(eapPktEncode)
}

func ConstructEapNoTypePkt(code radius.EapCode, pktID uint8) string {
	b := make([]byte, 4)
	b[0] = byte(code)
	b[1] = pktID
	binary.BigEndian.PutUint16(b[2:4], uint16(4))
	return base64.StdEncoding.EncodeToString(b)
}

func (p *Processor) logConfirmFailureAndInformUDM(
	id string, authType models.AusfUeAuthenticationAuthType, servingNetworkName, errStr, udmUrl string,
) {
	udmAuthType := models.UdmUeauAuthType(authType)

	switch authType {
	case models.AusfUeAuthenticationAuthType__5_G_AKA:
		logger.Auth5gAkaLog.Infoln(servingNetworkName, errStr)
		if sendErr := p.Consumer().SendAuthResultToUDM(id, udmAuthType, false, "", udmUrl); sendErr != nil {
			logger.Auth5gAkaLog.Infoln(sendErr.Error())
		}
	case models.AusfUeAuthenticationAuthType_EAP_AKA_PRIME:
		logger.AuthELog.Infoln(errStr)
		if sendErr := p.Consumer().SendAuthResultToUDM(id, udmAuthType, false, "", udmUrl); sendErr != nil {
			logger.AuthELog.Infoln(sendErr.Error())
		}
	}
}
