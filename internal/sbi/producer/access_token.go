package producer

import (
	"crypto/x509"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"

	nrf_context "github.com/free5gc/nrf/internal/context"
	"github.com/free5gc/nrf/internal/logger"
	"github.com/free5gc/nrf/pkg/factory"
	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/util/httpwrapper"
	"github.com/free5gc/util/mongoapi"
)

func HandleAccessTokenRequest(request *httpwrapper.Request) *httpwrapper.Response {
	// Param of AccessTokenRsp
	logger.AccTokenLog.Infoln("Handle AccessTokenRequest")

	accessTokenReq := request.Body.(models.AccessTokenReq)

	response, errResponse := AccessTokenProcedure(accessTokenReq)

	if response != nil {
		// status code is based on SPEC, and option headers
		return httpwrapper.NewResponse(http.StatusOK, nil, response)
	} else if errResponse != nil {
		return httpwrapper.NewResponse(http.StatusBadRequest, nil, errResponse)
	}
	problemDetails := &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}
	return httpwrapper.NewResponse(http.StatusForbidden, nil, problemDetails)
}

func AccessTokenProcedure(request models.AccessTokenReq) (
	*models.AccessTokenRsp, *models.AccessTokenErr,
) {
	logger.AccTokenLog.Infoln("In AccessTokenProcedure")

	var expiration int32 = 1000
	scope := request.Scope
	tokenType := "Bearer"
	now := int32(time.Now().Unix())

	errResponse := AccessTokenScopeCheck(request)
	if errResponse != nil {
		return nil, errResponse
	}

	// Create AccessToken
	nrfCtx := nrf_context.GetSelf()
	accessTokenClaims := models.AccessTokenClaims{
		Iss:            nrfCtx.Nrf_NfInstanceID,    // NF instance id of the NRF
		Sub:            request.NfInstanceId,       // nfInstanceId of service consumer
		Aud:            request.TargetNfInstanceId, // nfInstanceId of service producer
		Scope:          request.Scope,              // TODO: the name of the NF services for which the
		Exp:            now + expiration,           // access_token is authorized for use
		StandardClaims: jwt.StandardClaims{},
	}
	accessTokenClaims.IssuedAt = int64(now)

	// Use NRF private key to sign AccessToken
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), accessTokenClaims)
	accessToken, err := token.SignedString(nrfCtx.NrfPrivKey)
	if err != nil {
		logger.AccTokenLog.Warnln("Signed string error: ", err)
		return nil, &models.AccessTokenErr{
			Error: "invalid_request",
		}
	}

	response := &models.AccessTokenRsp{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiration,
		Scope:       scope,
	}
	return response, nil
}

func AccessTokenScopeCheck(req models.AccessTokenReq) *models.AccessTokenErr {
	// Check with nf profile
	collName := "NfProfile"
	reqGrantType := req.GrantType
	reqNfType := strings.ToUpper(string(req.NfType))
	reqTargetNfType := strings.ToUpper(string(req.TargetNfType))
	reqNfInstanceId := req.NfInstanceId

	if reqGrantType != "client_credentials" {
		return &models.AccessTokenErr{
			Error: "unsupported_grant_type",
		}
	}

	if reqNfType == "" || reqTargetNfType == "" || reqNfInstanceId == "" {
		return &models.AccessTokenErr{
			Error: "invalid_request",
		}
	}

	filter := bson.M{"nfInstanceId": reqNfInstanceId}
	consumerNfInfo, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.AccTokenLog.Errorln("mongoapi RestfulAPIGetOne error: " + err.Error())
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	nfProfile := models.NfProfile{}
	err = mapstructure.Decode(consumerNfInfo, &nfProfile)
	if err != nil {
		logger.AccTokenLog.Errorln("Certificate verify error: " + err.Error())
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	if strings.ToUpper(string(nfProfile.NfType)) != reqNfType {
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	// Verify NF's certificate with root certificate
	roots := x509.NewCertPool()
	nrfCtx := nrf_context.GetSelf()
	roots.AddCert(nrfCtx.RootCert)

	nfCert, err := openapi.ParseCertFromPEM(
		openapi.GetNFCertPath(factory.NrfConfig.GetCertBasePath(), reqNfType))
	if err != nil {
		logger.AccTokenLog.Errorln("NF Certificate get error: " + err.Error())
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	opts := x509.VerifyOptions{
		Roots:   roots,
		DNSName: reqNfType,
	}
	if _, err = nfCert.Verify(opts); err != nil {
		logger.AccTokenLog.Errorln("Certificate verify error: " + err.Error())
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	uri := nfCert.URIs[0]
	id := strings.Split(uri.Opaque, ":")[1]
	if id != reqNfInstanceId {
		logger.AccTokenLog.Errorln("Certificate verify error: NF Instance Id mismatch (Expected ID: " +
			reqNfInstanceId + " Received ID: " + id + ")")
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	// Check scope
	if reqTargetNfType == "NRF" {
		return nil
	}
	filter = bson.M{"nfType": reqTargetNfType}
	producerNfInfo, err := mongoapi.RestfulAPIGetOne(collName, filter)
	if err != nil {
		logger.AccTokenLog.Errorln("mongoapi.RestfulApiGetOne error: " + err.Error())
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	nfProfile = models.NfProfile{}
	err = mapstructure.Decode(producerNfInfo, &nfProfile)
	nfServices := *nfProfile.NfServices
	if err != nil {
		logger.AccTokenLog.Errorln("Certificate verify error: " + err.Error())
		return &models.AccessTokenErr{
			Error: "invalid_client",
		}
	}

	scopes := strings.Split(req.Scope, " ")

	for _, reqNfService := range scopes {
		found := false
		for _, nfService := range nfServices {
			if string(nfService.ServiceName) == reqNfService {
				if len(nfService.AllowedNfTypes) == 0 {
					found = true
					break
				} else {
					for _, nfType := range nfService.AllowedNfTypes {
						if string(nfType) == reqNfType {
							found = true
							break
						}
					}
					break
				}
			}
		}
		if !found {
			logger.AccTokenLog.Errorln("Certificate verify error: Request out of scope (" + reqNfService + ")")
			return &models.AccessTokenErr{
				Error: "invalid_scope",
			}
		}
	}

	return nil
}
