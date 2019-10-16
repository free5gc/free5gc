package ausf_producer

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/bronze1man/radius"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	// Nudm_UEAU "free5gc/lib/Nudm_UEAuthentication"
	"free5gc/lib/UeauCommon"
	"free5gc/lib/openapi/models"
	"free5gc/src/ausf/ausf_context"
	"free5gc/src/ausf/ausf_handler/ausf_message"
	"free5gc/src/ausf/logger"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func HandleEapAuthComfirmRequest(respChan chan ausf_message.HandlerResponseMessage, id string, body models.EapSession) {
	var response models.EapSession

	if !ausf_context.CheckIfAusfUeContextExists(id) {
		logger.EapAuthComfirmLog.Infoln("SUPI does not exist, confirmation failed")
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusBadRequest, problemDetails)
		return
	}
	ausfCurrentContext := ausf_context.GetAusfUeContext(id)
	servingNetworkName := ausfCurrentContext.ServingNetworkName
	eapPayload, _ := base64.StdEncoding.DecodeString(body.EapPayload)
	// fmt.Printf("eapPayload = %x\n", eapPayload)

	eapGoPkt := gopacket.NewPacket(eapPayload, layers.LayerTypeEAP, gopacket.Default)
	eapLayer := eapGoPkt.Layer(layers.LayerTypeEAP)
	eapContent, _ := eapLayer.(*layers.EAP)
	// fmt.Printf("d.Code=%x,\nd.Identitifier=%x,\nd.Type=%x,\nd.Data=%x\n", eapContent.Code, eapContent.Id, eapContent.Type, eapContent.TypeData)

	if eapContent.Code != layers.EAPCodeResponse {
		logConfirmFailureAndInformUDM(id, models.AuthType_EAP_AKA_PRIME, servingNetworkName, "eap packet code error", ausfCurrentContext.UdmUeauUrl)
		ausfCurrentContext.AuthStatus = models.AuthResult_FAILURE
		response.AuthResult = models.AuthResult_ONGOING
		failEapAkaNoti := ConstructFailEapAkaNotification(eapContent.Id)
		response.EapPayload = failEapAkaNoti
		ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, response)
		return
	}

	switch ausfCurrentContext.AuthStatus {
	case models.AuthResult_ONGOING:
		response.KSeaf = ausfCurrentContext.Kseaf
		response.Supi = id
		Kautn := ausfCurrentContext.K_aut
		XRES := ausfCurrentContext.XRES
		RES, decodeOK := decodeResMac(eapContent.TypeData, eapContent.Contents, Kautn)
		if !decodeOK {
			ausfCurrentContext.AuthStatus = models.AuthResult_FAILURE
			response.AuthResult = models.AuthResult_ONGOING
			logConfirmFailureAndInformUDM(id, models.AuthType_EAP_AKA_PRIME, servingNetworkName, "eap packet decode error", ausfCurrentContext.UdmUeauUrl)
			failEapAkaNoti := ConstructFailEapAkaNotification(eapContent.Id)
			response.EapPayload = failEapAkaNoti

		} else if XRES == string(RES) { // decodeOK && XRES == res, auth success
			logger.EapAuthComfirmLog.Infoln("Correct RES value, EAP-AKA' auth succeed")
			response.AuthResult = models.AuthResult_SUCCESS
			eapSuccPkt := ConstructEapNoTypePkt(radius.EapCodeSuccess, eapContent.Id)
			response.EapPayload = eapSuccPkt
			udmUrl := ausfCurrentContext.UdmUeauUrl
			if sendErr := sendAuthResultToUDM(id, models.AuthType_EAP_AKA_PRIME, true, servingNetworkName, udmUrl); sendErr != nil {
				logger.EapAuthComfirmLog.Infoln(sendErr.Error())
				var problemDetails models.ProblemDetails
				problemDetails.Cause = "UPSTREAM_SERVER_ERROR"
				ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusInternalServerError, problemDetails)
				return
			}
			ausfCurrentContext.AuthStatus = models.AuthResult_SUCCESS

		} else { // decodeOK but XRES != res, auth failed
			// fmt.Printf("XRES = %x\nstring(RES) = %x\n", XRES, RES)
			ausfCurrentContext.AuthStatus = models.AuthResult_FAILURE
			response.AuthResult = models.AuthResult_ONGOING
			logConfirmFailureAndInformUDM(id, models.AuthType_EAP_AKA_PRIME, servingNetworkName, "Wrong RES value, EAP-AKA' auth failed", ausfCurrentContext.UdmUeauUrl)
			failEapAkaNoti := ConstructFailEapAkaNotification(eapContent.Id)
			response.EapPayload = failEapAkaNoti
		}

	case models.AuthResult_FAILURE:
		eapFailPkt := ConstructEapNoTypePkt(radius.EapCodeFailure, uint8(eapPayload[1]))
		response.EapPayload = eapFailPkt
		response.AuthResult = models.AuthResult_FAILURE
	}

	ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, response)
}

func HandleAuth5gAkaComfirmRequest(respChan chan ausf_message.HandlerResponseMessage, id string, body models.ConfirmationData) {
	var response models.ConfirmationDataResponse
	success := false
	response.AuthResult = models.AuthResult_FAILURE

	if !ausf_context.CheckIfAusfUeContextExists(id) {
		logger.Auth5gAkaComfirmLog.Infoln("SUPI does not exist, confirmation failed")
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "USER_NOT_FOUND"
		ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusBadRequest, problemDetails)
		return
	}

	ausfCurrentContext := ausf_context.GetAusfUeContext(id)
	servingNetworkName := ausfCurrentContext.ServingNetworkName

	// Compare the received RES* with the stored XRES*
	// fmt.Printf("res*: %x\nXres*: %x\n", body.ResStar, ausfCurrentContext.XresStar)
	if strings.Compare(body.ResStar, ausfCurrentContext.XresStar) == 0 {
		ausfCurrentContext.AuthStatus = models.AuthResult_SUCCESS
		response.AuthResult = models.AuthResult_SUCCESS
		success = true
		logger.Auth5gAkaComfirmLog.Infoln("5G AKA confirmation succeeded")
		response.Kseaf = ausfCurrentContext.Kseaf
	} else {
		ausfCurrentContext.AuthStatus = models.AuthResult_FAILURE
		response.AuthResult = models.AuthResult_FAILURE
		logConfirmFailureAndInformUDM(id, models.AuthType__5_G_AKA, servingNetworkName, "5G AKA confirmation failed", ausfCurrentContext.UdmUeauUrl)
	}

	if sendErr := sendAuthResultToUDM(id, models.AuthType__5_G_AKA, success, servingNetworkName, ausfCurrentContext.UdmUeauUrl); sendErr != nil {
		logger.Auth5gAkaComfirmLog.Infoln(sendErr.Error())
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "UPSTREAM_SERVER_ERROR"
		ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusInternalServerError, problemDetails)

		return
	}

	response.Supi = id
	ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusOK, response)
}

func HandleUeAuthPostRequest(respChan chan ausf_message.HandlerResponseMessage, body models.AuthenticationInfo) {
	var response models.UeAuthenticationCtx
	var authInfoReq models.AuthenticationInfoRequest

	supiOrSuci := body.SupiOrSuci
	// fmt.Println("Got supi:", supiOrSuci)

	// check if SEAF is authorized to use the serving network name as in 33501 clause 6.1.2
	snName := body.ServingNetworkName
	servingNetworkAuthorized := ausf_context.IsServingNetworkAuthorized(snName)
	if !servingNetworkAuthorized {
		var problemDetails models.ProblemDetails
		problemDetails.Cause = "SERVING_NETWORK_NOT_AUTHORIZED"
		logger.UeAuthPostLog.Infoln("403 forbidden: serving network NOT AUTHORIZED")
		ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusForbidden, problemDetails)
		return
	}
	logger.UeAuthPostLog.Infoln("Serving network authorized")

	response.ServingNetworkName = snName
	authInfoReq.ServingNetworkName = snName
	self := ausf_context.GetSelf()
	authInfoReq.AusfInstanceId = self.GetSelfID()

	udmUrl := getUdmUrl(self.NrfUri)
	client := createClientToUdmUeau(udmUrl)
	authInfoResult, _, err := client.GenerateAuthDataApi.GenerateAuthData(context.Background(), supiOrSuci, authInfoReq)
	if err != nil {
		logger.UeAuthPostLog.Infoln(err.Error())
		var problemDetails models.ProblemDetails
		if authInfoResult.AuthenticationVector == nil {
			problemDetails.Cause = "AV_GENERATION_PROBLEM"
		} else {
			problemDetails.Cause = "UPSTREAM_SERVER_ERROR"
		}
		ausf_message.SendHttpResponseMessage(respChan, nil, http.StatusInternalServerError, problemDetails)
		return
	}

	ueid := authInfoResult.Supi
	ausfUeContext := ausf_context.NewAusfUeContext(ueid)
	ausf_context.AddAusfUeContextToPool(ausfUeContext)
	ausfUeContext.ServingNetworkName = snName
	ausfUeContext.AuthStatus = models.AuthResult_ONGOING
	ausfUeContext.UdmUeauUrl = udmUrl

	locationURI := self.Url + "/nausf-auth/v1/ue-authentications/" + ueid
	putLink := locationURI
	if authInfoResult.AuthType == models.AuthType__5_G_AKA {
		logger.UeAuthPostLog.Infoln("Use 5G AKA auth method")
		putLink += "/5g-aka-confirmation"

		// Derive HXRES* from XRES*
		concat := authInfoResult.AuthenticationVector.Rand + authInfoResult.AuthenticationVector.XresStar
		hxresStarBytes, _ := hex.DecodeString(concat)
		hxresStarAll := sha256.Sum256(hxresStarBytes)
		hxresStar := hex.EncodeToString(hxresStarAll[16:]) // last 128 bits
		// fmt.Printf("5G AV Rand: %s,  XresStar: %s, Autn: %s\n", authInfoResult.AuthenticationVector.Rand, authInfoResult.AuthenticationVector.XresStar, authInfoResult.AuthenticationVector.Autn)
		// fmt.Printf("hxresStar = %x\n", hxresStar) // 231b3f2e0a8a5082c19fdd6735888c4a

		// Derive Kseaf from Kausf
		// Test data
		// P0 := "internet"
		Kausf := authInfoResult.AuthenticationVector.Kausf
		KausfDecode, _ := hex.DecodeString(Kausf)
		P0 := []byte(snName)
		Kseaf := UeauCommon.GetKDFValue(KausfDecode, UeauCommon.FC_FOR_KSEAF_DERIVATION, P0, UeauCommon.KDFLen(P0))
		ausfUeContext.XresStar = authInfoResult.AuthenticationVector.XresStar
		ausfUeContext.Kausf = Kausf
		ausfUeContext.Kseaf = hex.EncodeToString(Kseaf)

		var av5gAka models.Av5gAka
		av5gAka.Rand = authInfoResult.AuthenticationVector.Rand
		av5gAka.Autn = authInfoResult.AuthenticationVector.Autn
		av5gAka.HxresStar = hxresStar

		response.Var5gAuthData = av5gAka

	} else if authInfoResult.AuthType == models.AuthType_EAP_AKA_PRIME {
		logger.UeAuthPostLog.Infoln("Use EAP-AKA' auth method")
		putLink += "/eap-session"

		identity := ueid
		ikPrime := authInfoResult.AuthenticationVector.IkPrime
		ckPrime := authInfoResult.AuthenticationVector.CkPrime
		RAND := authInfoResult.AuthenticationVector.Rand
		AUTN := authInfoResult.AuthenticationVector.Autn
		XRES := authInfoResult.AuthenticationVector.Xres
		ausfUeContext.XRES = XRES

		// Test data
		// identity = "0555444333222111"

		K_encr, K_aut, K_re, MSK, EMSK := eapAkaPrimePrf(ikPrime, ckPrime, identity)
		_, _, _, _, _ = K_encr, K_aut, K_re, MSK, EMSK
		ausfUeContext.K_aut = K_aut
		Kausf := EMSK[0:32]
		ausfUeContext.Kausf = Kausf
		KausfDecode, _ := hex.DecodeString(Kausf)
		P0 := []byte(snName)
		Kseaf := UeauCommon.GetKDFValue(KausfDecode, UeauCommon.FC_FOR_KSEAF_DERIVATION, P0, UeauCommon.KDFLen(P0))
		ausfUeContext.Kseaf = hex.EncodeToString(Kseaf)

		var eapPkt radius.EapPacket
		var randIdentifier int
		rand.Seed(time.Now().Unix())

		eapPkt.Code = radius.EapCode(1)
		randIdentifier = rand.Intn(256)
		eapPkt.Identifier = uint8(randIdentifier)
		eapPkt.Type = radius.EapType(50) // accroding to RFC5448 6.1
		atRand, _ := EapEncodeAttribute("AT_RAND", RAND)
		atAutn, _ := EapEncodeAttribute("AT_AUTN", AUTN)
		atKdf, _ := EapEncodeAttribute("AT_KDF", snName)
		atKdfInput, _ := EapEncodeAttribute("AT_KDF_INPUT", snName)
		atMAC, _ := EapEncodeAttribute("AT_MAC", "")

		dataArrayBeforeMAC := atRand + atAutn + atMAC + atKdf + atKdfInput
		eapPkt.Data = []byte(dataArrayBeforeMAC)
		encodedPktBeforeMAC := eapPkt.Encode()

		MACvalue := CalculateAtMAC([]byte(K_aut), encodedPktBeforeMAC)
		// fmt.Printf("MAC value = %x\n", MACvalue)
		atMacNum := fmt.Sprintf("%02x", ausf_context.AT_MAC_ATTRIBUTE)
		atMACfirstRow, _ := hex.DecodeString(atMacNum + "05" + "0000")
		wholeAtMAC := append(atMACfirstRow, MACvalue...)

		atMAC = string(wholeAtMAC)
		dataArrayAfterMAC := atRand + atAutn + atMAC + atKdf + atKdfInput

		eapPkt.Data = []byte(dataArrayAfterMAC)
		encodedPktAfterMAC := eapPkt.Encode()
		response.Var5gAuthData = base64.StdEncoding.EncodeToString(encodedPktAfterMAC)

		// fmt.Printf("p.Code=%x,\np.Identitifier=%x,\np.Type=%x,\np.Data=%x\n", byte(eapPkt.Code), byte(eapPkt.Identifier), byte(eapPkt.Type), eapPkt.Data)
		// fmt.Printf("encodedPktAfterMAC: %x\n", encodedPktAfterMAC)
		// fmt.Printf("Var5gAuthData: %x\n", response.Var5gAuthData)
	}

	var linksValue = models.LinksValueSchema{Href: putLink}
	response.Links = make(map[string]models.LinksValueSchema)
	response.Links["link"] = linksValue
	response.AuthType = authInfoResult.AuthType

	respHeader := make(http.Header)
	respHeader.Set("Location", locationURI)
	ausf_message.SendHttpResponseMessage(respChan, respHeader, http.StatusCreated, response)
}
